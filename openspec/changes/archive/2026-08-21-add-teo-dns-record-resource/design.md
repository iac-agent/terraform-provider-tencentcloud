## Context

TEO（EdgeOne）提供 DNS 记录管理能力，SDK `teov20220901` 暴露了 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords` 四个接口。当前 provider 已存在一个 `tencentcloud_teo_dns_record` 资源封装了同一组接口（并额外使用 `ModifyDnsRecordsStatus` 支持启停状态）。

本次需要新增一个独立的 `tencentcloud_teo_dns_record_v3` 资源，差异在于：

- `status`、`created_on`、`modified_on` 全部作为纯 `Computed` 只读字段，**不**引入 `ModifyDnsRecordsStatus`，避免用户通过 Terraform 直接修改云端只读状态；
- 资源名以 `_v3` 后缀独立存在，与既有 `tencentcloud_teo_dns_record` 共存，互不影响。

关键约束（来自 vendor 下 SDK 源码确认）：

- `CreateDnsRecordRequest` 入参：`ZoneId`、`Name`、`Type`、`Content`、`Location`、`TTL`、`Weight`、`Priority`；出参：`RecordId`。
- `DescribeDnsRecordsRequest` 入参：`ZoneId`、`Offset`、`Limit`（上限 1000）、`Filters`（`[]*AdvancedFilter`，其中 `Name`/`Values`/`Fuzzy`）、`SortBy`、`SortOrder`、`Match`；出参：`TotalCount`、`DnsRecords`（`[]*DnsRecord`）。
- `DnsRecord` 结构体字段：`ZoneId`、`RecordId`、`Name`、`Type`、`Location`、`Content`、`TTL`、`Weight`、`Priority`、`Status`、`CreatedOn`、`ModifiedOn`。其中 `ZoneId`、`Status`、`CreatedOn`、`ModifiedOn` 在 SDK 注释中明确标注为“仅做出参使用，在 ModifyDnsRecords 不可作为入参使用，如有传此参数，会忽略”。
- `ModifyDnsRecordsRequest` 入参：`ZoneId`、`DnsRecords`（`[]*DnsRecord`，一次最多 100 条）。
- `DeleteDnsRecordsRequest` 入参：`ZoneId`、`RecordIds`（`[]*string`，上限 1000）。

这些接口均为同步接口（无异步状态等待要求），无需额外 Timeouts 轮询。

## Goals / Non-Goals

**Goals:**

- 新增 `tencentcloud_teo_dns_record_v3` 资源，提供完整 CRUD 生命周期；
- schema 与云 API 字段 1:1 对应，只暴露云 API 支持的字段；
- 资源 ID 采用 `zone_id#record_id` 复合 ID，支持 import；
- 复用已存在的 service 层 helper `DescribeTeoDnsRecordById`，避免重复实现；
- 通过 gomonkey mock 单元测试覆盖业务逻辑分支。

**Non-Goals:**

- 不修改既有 `tencentcloud_teo_dns_record` 资源；
- 不支持启停状态（不调用 `ModifyDnsRecordsStatus`），`status` 仅作 Computed 输出；
- 不暴露 `DescribeDnsRecords` 的分页（`Offset`/`Limit`）与排序（`SortBy`/`SortOrder`/`Match`）参数给用户；
- 不将 `filters` 做成用户可配置的 schema 字段（Read 内部按 `id` 过滤即可）。

## Decisions

### Decision 1: 资源 ID 使用 `zone_id#record_id` 复合 ID

**选择**：Create 成功后，`d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))`，Read/Update/Delete 中 `strings.Split(d.Id(), tccommon.FILED_SP)` 拆出 `zoneId`、`recordId`。

**理由**：
- `record_id` 仅在站点（`zone_id`）内唯一，需要 `zone_id` 共同定位；
- 与既有 `tencentcloud_teo_dns_record` 及 provider 大多数复合 ID 资源保持一致，import 示例清晰（`{zoneId}#{recordId}`）。

### Decision 2: Read 复用既有 service helper `DescribeTeoDnsRecordById`

**选择**：Read 调用 `service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)`（已存在于 `service_tencentcloud_teo.go`），该 helper 内部用 `DescribeDnsRecords` + `Filters=[{Name:"id", Values:[recordId]}]` 查询。

**理由**：
- 避免重复实现相同的查询逻辑，保持 service 层单一实现；
- helper 已处理 `resource.Retry(tccommon.ReadRetryTimeout)` 与空结果返回 `nil`。

### Decision 3: `status`/`created_on`/`modified_on` 为纯 Computed

**选择**：这三个字段仅声明 `Computed: true`，Update 的 `mutableArgs` 仅包含 `name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`。

**理由**：
- SDK 明确 `Status`/`CreatedOn`/`ModifiedOn` 在 `ModifyDnsRecords` 中仅做出参使用，不可作为入参；
- 启停状态属于云端只读语义，不应通过本资源驱动（区别于既有 `tencentcloud_teo_dns_record`）。

### Decision 4: Update 仅下发可变更字段 + `RecordId` 定位

**选择**：Update 构造 `ModifyDnsRecordsRequest`，设置 `ZoneId` 与单元素 `DnsRecords`，其中 `DnsRecord` 仅设置 `RecordId`（定位）以及发生变化的可变更字段（`Name`/`Type`/`Content`/`Location`/`TTL`/`Weight`/`Priority`）。**不**设置 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn`（SDK 明确忽略，避免误导）。

**理由**：
- `ModifyDnsRecords` 是批量接口，单条记录场景传一个元素即可；
- 遵循 SDK 注释，避免下发被忽略的字段造成混淆。

### Decision 5: `zone_id` 为 ForceNew，其余可变更

**选择**：仅 `zone_id` 设置 `ForceNew: true`；`name`、`type`、`content`、`location`、`ttl`、`weight`、`priority` 均可在 Update 中变更。

**理由**：
- 记录不能跨站点移动，`zone_id` 变更应触发重建；
- `ModifyDnsRecords` 支持修改 `Name`/`Type`/`Content`/`Location`/`TTL`/`Weight`/`Priority`，因此这些字段不必 ForceNew。

### Decision 6: 单元测试使用 gomonkey mock，不依赖真实云账号

**选择**：新增 `resource_tc_teo_dns_record_v3_test.go`，用 gomonkey mock `CreateDnsRecordWithContext`、`DescribeDnsRecords`、`ModifyDnsRecordsWithContext`、`DeleteDnsRecordsWithContext`，仅做业务逻辑单测（创建/读取/更新/删除分支）。

**理由**：符合“新增资源不使用 Terraform 验收测试套件，而用 gomonkey mock 云 API”的要求，测试可在无 TENCENTCLOUD_SECRET_ID/KEY 环境运行。

## Risks / Trade-offs

- **Risk**：与既有 `tencentcloud_teo_dns_record` 存在功能重叠，用户可能困惑该选哪个 → **Mitigation**：在文档与 proposal 中明确 `_v3` 的差异（无状态启停能力、`status` 纯只读）。
- **Risk**：`DescribeDnsRecords` 按 `id` 过滤存在最终一致性，创建后立即 Read 可能短暂返回空 → **Mitigation**：沿用 provider 统一 `resource.Retry` + `ReadRetryTimeout`，空结果返回后由上层重试/刷新收敛。
- **Risk**：`ModifyDnsRecords` 为批量接口，若未来云端对单条修改语义有差异 → **Mitigation**：当前单元素调用即等价单条修改，SDK 无单条修改接口，风险可控。
- **Trade-off**：`ttl`/`weight`/`priority` 采用 Optional+Computed，云端存在默认值（300/-1/0），Read 回填可能与未显式配置产生短暂 plan diff → 可接受，Terraform 会在 refresh 后收敛。

## Migration Plan

- 纯新增资源，无 state 迁移需求；
- 资源注册：在 `provider.go` ResourcesMap 追加 `"tencentcloud_teo_dns_record_v3": teo.ResourceTencentCloudTeoDnsRecordV3()`，并同步更新 `provider.md`；
- 文档：新增 `resource_tc_teo_dns_record_v3.md`，由 `make doc` 生成 website docs；
- 回滚：删除注册项与新增文件即可，不影响任何既有资源。

## Open Questions

- 无
