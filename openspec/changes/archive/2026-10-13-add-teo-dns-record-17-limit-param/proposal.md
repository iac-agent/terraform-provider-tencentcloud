## Why

腾讯云 TEO（边缘安全加速平台，EdgeOne）需要以通用资源 `tencentcloud_teo_dns_record_17` 管理 DNS 记录的全生命周期（增删改查）。当前 provider 的 master 分支中尚不存在该资源，且需求要求在调用 `DescribeDnsRecords` 查询接口时新增可选入参 `Limit`（分页查询限制数目），使用户可通过 Terraform 配置控制查询行为，避免采用 API 默认分页值（20）导致单页查询不足的问题。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_17`（RESOURCE_KIND_GENERAL），管理 TEO 站点下的一条 DNS 记录（创建 / 查询 / 修改 / 删除 / 导入），代码文件为 `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go`。
- 资源 schema 字段（顶层平铺）：
  - `zone_id`（Required，ForceNew，TypeString）：站点 ID。
  - `name`（Required，TypeString）：DNS 记录名。
  - `type`（Required，TypeString）：DNS 记录类型（A / AAAA / MX / CNAME / TXT / NS / CAA / SRV）。
  - `content`（Required，TypeString）：DNS 记录内容。
  - `location`（Optional+Computed，TypeString）：解析线路，默认 Default。
  - `ttl`（Optional+Computed，TypeInt）：缓存时间，60~86400 秒，默认 300。
  - `weight`（Optional+Computed，TypeInt）：记录权重，-1~100，默认 -1。
  - `priority`（Optional+Computed，TypeInt）：MX 优先级，0~50，默认 0。
  - `status`（Optional+Computed，TypeString）：解析状态（enable / disable），更新时经 `ModifyDnsRecordsStatus` 启停。
  - `record_id`（Computed，TypeString）：DNS 记录 ID。
  - `created_on` / `modified_on`（Computed，TypeString）：创建 / 修改时间。
  - **`Limit`（Optional，TypeInt）：本次需求唯一新增参数——`DescribeDnsRecords` 接口的分页查询限制数目（`request.Limit`），不设置时默认使用云 API 注释标注的上限值 1000，保证按记录 ID 过滤的单条查询不受默认分页（20）影响。**
- CRUD 接口映射（均使用 vendor 中已存在的 `teov20220901` SDK）：
  - Create → `CreateDnsRecord`（入参 ZoneId/Name/Type/Content/Location/TTL/Weight/Priority；出参 RecordId）。
  - Read → `DescribeDnsRecords`（服务层按 `id` AdvancedFilter 过滤定位单条记录，`request.Limit` 由新增的 `Limit` 参数控制，默认 1000）。
  - Update → `ModifyDnsRecords`（单条 `DnsRecords` 元素，仅传可变字段 Name/Type/Location/Content/TTL/Weight/Priority；`ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` 为出参字段，不作为入参传入）；`status` 变更经 `ModifyDnsRecordsStatus`。
  - Delete → `DeleteDnsRecords`（ZoneId + RecordIds）。
- 资源 ID 使用复合 ID：`zoneId#recordId`（`tccommon.FILED_SP` 分隔），Read/Update/Delete 从 `d.Id()` 解析；支持 `terraform import`（import 示例中说明需使用联合 ID）。
- 所有云 API 调用以 `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` 包裹 `resource.Retry` 重试，失败时用 `tccommon.RetryError()` 包装返回。
- 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `tencentcloud_teo_dns_record_17`，并在 `tencentcloud/provider.md` 中补充注册说明。
- 新增资源文档 `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md`（一句话描述 + Example Usage + Import）。
- 新增单元测试 `tencentcloud/services/teo/resource_tc_teo_dns_record_17_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件），覆盖 Create/Read/Update/Delete 业务逻辑及 `Limit` 参数传递。

## Capabilities

### New Capabilities
- `teo-dns-record-17-resource`: 新资源 `tencentcloud_teo_dns_record_17` 的全生命周期管理能力，包括 schema 定义（含 `DescribeDnsRecords` 查询入参 `Limit`）、Create（CreateDnsRecord）、Read（DescribeDnsRecords 按记录 ID 过滤定位、Limit 控制分页上限）、Update（ModifyDnsRecords 单条修改 + ModifyDnsRecordsStatus 状态启停）、Delete（DeleteDnsRecords）、复合 ID（zoneId#recordId）解析与 terraform import 支持。

### Modified Capabilities
<!-- 无：本次为纯新增独立资源，不修改任何既有资源（含 tencentcloud_teo_dns_record）的 spec 行为 -->

## Impact

- **新增代码**：
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go`（schema + CRUD，代码风格严格参考 `tencentcloud_igtm_strategy` 资源）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md`（资源文档示例）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_17_test.go`（gomonkey mock 单元测试）
- **修改代码**：
  - `tencentcloud/provider.go`：在 `ResourcesMap` 中注册 `"tencentcloud_teo_dns_record_17": teo.ResourceTencentCloudTeoDnsRecord17()`
  - `tencentcloud/provider.md`：在 TEO 分类下补充资源名注册说明
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`：新增服务层方法 `DescribeTeoDnsRecord17ById`（封装 DescribeDnsRecords 按记录 ID 过滤 + Limit 分页查询）
- **API 依赖**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 中的 `CreateDnsRecord`、`DescribeDnsRecords`（含 `Limit` 入参）、`ModifyDnsRecords`、`ModifyDnsRecordsStatus`、`DeleteDnsRecords`（均已存在于 vendor，无需变更 vendor）。
- **无破坏性变更**：纯新增资源，与既有 `tencentcloud_teo_dns_record` 并存，互不影响 state 与配置。
- **文档**：website 文档由收尾阶段 `make doc` 自动生成，禁止手改 `website/` 目录。
