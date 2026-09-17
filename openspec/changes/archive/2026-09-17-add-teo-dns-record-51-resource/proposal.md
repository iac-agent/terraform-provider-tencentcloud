## Why

腾讯云 TEO（边缘安全加速平台，EdgeOne）需要以新的通用资源 `tencentcloud_teo_dns_record_51` 管理 NS 接入模式站点下的 DNS 记录全生命周期（增删改查）。虽然 provider 中已存在 `tencentcloud_teo_dns_record` 资源，但本次需求要求按最新的云 API 参数映射（含 `Location`、`Weight`、`Priority` 等新字段语义以及 `DescribeDnsRecords` 的 `AdvancedFilter` 过滤能力）生成一个独立的 RESOURCE_KIND_GENERAL 资源 `tencentcloud_teo_dns_record_51`，以承载新的 schema 契约，避免破坏既有资源的 schema 与 state。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_51`（RESOURCE_KIND_GENERAL），管理 TEO 站点下的一条 DNS 记录（创建 / 查询 / 修改 / 删除 / 导入），代码文件为 `tencentcloud/services/teo/resource_tc_teo_dns_record_51.go`。
- 资源 schema 字段（顶层平铺，不引入额外嵌套层）：
  - `zone_id`（Required，ForceNew，TypeString）：站点 ID。
  - `name`（Required，TypeString）：DNS 记录名（中文/韩文/日文需转 punycode）。
  - `type`（Required，TypeString）：DNS 记录类型（A / AAAA / MX / CNAME / TXT / NS / CAA / SRV）。
  - `content`（Required，TypeString）：DNS 记录内容。
  - `location`（Optional+Computed，TypeString）：解析线路，默认 Default。
  - `ttl`（Optional+Computed，TypeInt）：缓存时间，60~86400 秒，默认 300。
  - `weight`（Optional+Computed，TypeInt）：记录权重，-1~100，默认 -1（不设置）。
  - `priority`（Optional+Computed，TypeInt）：MX 优先级，0~50，默认 0。
  - `record_id`（Computed，TypeString）：DNS 记录 ID。
  - `status`（Computed，TypeString）：解析状态（enable / disable）。
  - `created_on`（Computed，TypeString）：创建时间。
  - `modified_on`（Computed，TypeString）：修改时间。
- CRUD 接口映射（均使用 vendor 中已存在的 `teov20220901` SDK）：
  - Create → `CreateDnsRecord`（入参 ZoneId/Name/Type/Content/Location/TTL/Weight/Priority；出参 RecordId）。
  - Read → `DescribeDnsRecords`（服务层按 `id` 过滤器 + `Limit=1000`（云 API 注释标注的最大值）分页轮询定位单条记录，读取平铺的所有出参字段）。
  - Update → `ModifyDnsRecords`（单条 `DnsRecords` 元素，仅传可变字段 ZoneId/RecordId/Name/Type/Location/Content/TTL/Weight/Priority；`ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` 在 ModifyDnsRecords 中为出参字段，不作为入参传入）。
  - Delete → `DeleteDnsRecords`（ZoneId + RecordIds）。
- 资源 ID 使用复合 ID：`zoneId#recordId`（`tccommon.FILED_SP` 分隔），Read/Update/Delete 从 `d.Id()` 解析；支持 `terraform import`。
- 所有云 API 调用均以 `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` 包裹 `resource.Retry` 重试；失败时用 `tccommon.RetryError()` 包装返回。
- 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `tencentcloud_teo_dns_record_51`，并在 `tencentcloud/provider.md` 中补充注册说明。
- 新增资源文档 `tencentcloud/services/teo/resource_tc_teo_dns_record_51.md`（一句话描述 + Example Usage + Import），供 `make doc` 生成 website 文档。
- 新增单元测试 `tencentcloud/services/teo/resource_tc_teo_dns_record_51_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件），覆盖 Create/Read/Update/Delete 业务逻辑。

## Capabilities

### New Capabilities
- `teo-dns-record-51-resource`: 新资源 `tencentcloud_teo_dns_record_51` 的全生命周期管理能力，包括 schema 定义、Create（CreateDnsRecord）、Read（DescribeDnsRecords 按记录 ID 过滤定位）、Update（ModifyDnsRecords 单条修改）、Delete（DeleteDnsRecords）、复合 ID（zoneId#recordId）解析与 terraform import 支持。

### Modified Capabilities
<!-- 无：本次仅新增独立资源，不修改任何既有资源（含 tencentcloud_teo_dns_record）的 spec 行为 -->

## Impact

- **新增代码**：
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_51.go`（schema + CRUD，代码风格严格参考 `tencentcloud_igtm_strategy` 资源）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_51.md`（资源文档示例）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_51_test.go`（gomonkey mock 单元测试）
- **修改代码**：
  - `tencentcloud/provider.go`：在 `ResourcesMap` 中注册 `"tencentcloud_teo_dns_record_51": teo.ResourceTencentCloudTeoDnsRecord51()`
  - `tencentcloud/provider.md`：在 TEO 分类下补充资源名注册说明
- **API 依赖**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 中的 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords`（均已存在于 vendor，无需变更 vendor）。
- **无破坏性变更**：纯新增资源，与既有 `tencentcloud_teo_dns_record` 并存，互不影响 state 与配置。
- **文档**：website 文档由收尾阶段 `make doc` 自动生成，禁止手改 `website/` 目录。
