# Proposal: add-teo-zone-v2-resource

## Why

EdgeOne（TEO）站点管理能力近年来持续演进：云API `CreateZone` 新增了 `PlanId`（绑定套餐激活）、`AllowDuplicates`（同名站点接入）、`JumpStart`（跳过 DNS 记录扫描）等参数，`ModifyZone` 新增了 `VanityNameServers`（自定义 NS）、`ZoneName`（无域名接入切换 CNAME）等参数，`DescribeZones` 也返回了按接入类型区分的详情（NSDetail/CNAMEDetail/DNSPodDetail）、计费资源列表（Resources）、版本管理工作模式（WorkModeInfos）等大量新字段。现有资源 `tencentcloud_teo_zone` schema 陈旧（如 `type` 必填、缺少 `plan_id`/`allow_duplicates`/`jump_start`/`vanity_name_servers` 等参数），且受硬约束"不能修改已有资源的 schema（除非只新增 Optional 字段）"限制，无法完整暴露新版站点 CRUD 能力。因此需要新增 `tencentcloud_teo_zone_v2` 资源，基于最新的 `CreateZone`/`DescribeZones`/`ModifyZone`/`DeleteZone` 接口完整管理站点生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_zone_v2`（RESOURCE_KIND_GENERAL），文件为 `tencentcloud/services/teo/resource_tc_teo_zone_v2.go`，实现站点（Zone）的完整 CRUD：
  - Create：调用 `CreateZone`（支持 `type`、`zone_name`、`area`、`plan_id`、`alias_zone_name`、`tags`、`allow_duplicates`、`jump_start` 入参；创建成功后立即 `d.SetId(zoneId)`），同时把创建返回的归属权验证信息（`ownership_verification` 含 dns/file/ns 验证详情）写入 state。
  - Read：调用 `DescribeZones`（以 `zone-id` 过滤、按云API注释标注的最大分页值 `Limit=100` 拉取），将 Zone 全量字段（含 `alias_zone_name`、`area`、`type`、`tags`、计费资源列表 `resources`、按接入类型返回的 `ns_detail`/`cname_detail`/`dns_pod_detail`、`status`、`cname_status`、`active_status`、`lock_status`、`paused`、`work_mode_infos` 等）写入 state；若资源不存在，先打印 `log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())` 再 `d.SetId("")`。
  - Update：调用 `ModifyZone`（支持 `type`、`vanity_name_servers`（`switch`/`servers`）、`alias_zone_name`、`area`、`zone_name` 可变参数；`zone_id` 从 `d.Id()` 取值）。
  - Delete：调用 `DeleteZone`（`zone_id` 从 `d.Id()` 取值）。
- 资源 id 采用单字段 `ZoneId`（无联合 id）。
- 新增单元测试 `tencentcloud/services/teo/resource_tc_teo_zone_v2_test.go`，使用 gomonkey mock 云API（不使用 terraform 测试套件）。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册 `tencentcloud_teo_zone_v2`。
- 新增资源文档 `tencentcloud/services/teo/resource_tc_teo_zone_v2.md`（含一句话描述、Example Usage、Import 部分；不手写 Argument/Attribute Reference）。

不修改现有 `tencentcloud_teo_zone` 资源（保持向后兼容）。

## Capabilities

### New Capabilities
- `teo-zone-v2-resource`: 基于 teo/v20220901 云API `CreateZone`/`DescribeZones`/`ModifyZone`/`DeleteZone` 的 EdgeOne 站点（Zone）资源管理能力，覆盖站点创建（含套餐绑定、同名接入、标签、归属权验证信息透出）、查询（含计费资源、接入类型详情、版本管理工作模式）、修改（含自定义 NS）、删除全生命周期。

### Modified Capabilities

（无。现有能力规格不变，`tencentcloud_teo_zone` 保持原样。）

## Impact

- 新增代码：
  - `tencentcloud/services/teo/resource_tc_teo_zone_v2.go`（资源实现）
  - `tencentcloud/services/teo/resource_tc_teo_zone_v2_test.go`（gomonkey 单测）
  - `tencentcloud/services/teo/resource_tc_teo_zone_v2.md`（资源文档）
- 修改代码：
  - `tencentcloud/provider.go`：注册 `tencentcloud_teo_zone_v2` 资源
  - `tencentcloud/provider.md`：在 TencentCloud EdgeOne(TEO) Resource 列表中追加 `tencentcloud_teo_zone_v2`
- 依赖：仅使用 vendor 下已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，无新增第三方依赖；无需生成 `_extension.go` 文件。
- 兼容性：纯新增资源，不影响任何现有 TF 配置与 state；现有 `tencentcloud_teo_zone` 不变。
- 云API接口均为同步接口（无 FlowId/TaskId 返回），无异步轮询诉求。
