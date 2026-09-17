## Context

腾讯云 TEO（EdgeOne）NS 接入模式站点下，DNS 记录是最常用的资源类型之一。provider 中已有 `tencentcloud_teo_dns_record` 资源（文件 `tencentcloud/services/teo/resource_tc_teo_dns_record.go`），本次需求按新的参数映射契约新增一个**独立**的 RESOURCE_KIND_GENERAL 资源 `tencentcloud_teo_dns_record_51`，两者并存、互不影响。

云 API 能力核对（全部来自 vendor 中 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）：

- `CreateDnsRecordRequest`：入参 `ZoneId`、`Name`、`Type`、`Content`、`Location`、`TTL`、`Weight`、`Priority`（均为 `*string`/`*int64`，omitnil）。出参 `CreateDnsRecordResponseParams.RecordId *string`。
- `DescribeDnsRecordsRequest`：入参 `ZoneId`（必填）、`Offset`、`Limit`（注释标注上限 1000）、`Filters []*AdvancedFilter`（`Name`/`Values`/`Fuzzy`）、`SortBy`、`SortOrder`、`Match`。出参 `TotalCount` 与 `DnsRecords []*DnsRecord`。
- `ModifyDnsRecordsRequest`：入参 `ZoneId`、`DnsRecords []*DnsRecord`（一次最多 100 条；`DnsRecord` 结构中 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` **仅做出参**，ModifyDnsRecords 传入会被忽略）。
- `DeleteDnsRecordsRequest`：入参 `ZoneId`、`RecordIds []*string`（上限 1000）。

现有基础设施工具备：
- `TeoService.DescribeTeoDnsRecordById(ctx, zoneId, recordId)`（`service_tencentcloud_teo.go:1796`）已实现按 `id` 过滤器查询单条记录，内部已含 `resource.Retry(tccommon.ReadRetryTimeout)` + `tccommon.RetryError` 重试。
- `connectivity.TencentCloudClient.UseTeoV20220901Client()` 客户端绑定已就绪。
- provider.go `ResourcesMap` 已注册大量 teo 资源，新增注册为纯加法。

约束：
- 不得修改既有 `tencentcloud_teo_dns_record` 资源（向后兼容硬约束）。
- 异步接口才需要调用后轮询；本组 4 个接口（CreateDnsRecord/DescribeDnsRecords/ModifyDnsRecords/DeleteDnsRecords）在 vendor client.go 注释中均**未标注为异步接口**，无需任务 ID 轮询。
- 代码风格严格参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（RESOURCE_KIND_GENERAL 标杆资源）。
- vendor 模式管理依赖，本次 4 个 API 均已 vendored，无需 go.mod 变更。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_51` 资源，完整管理 TEO DNS 记录的 Create/Read/Update/Delete/Import 生命周期。
- schema 字段顶层平铺（不引入 `xxx_set`/`xxx_list` 包裹层），与云 API 字段一一对应。
- 复合 ID `zoneId#recordId`（`tccommon.FILED_SP`），支持 import。
- 所有 SDK 调用走 `resource.Retry` + `tccommon.ReadRetryTimeout`/`WriteRetryTimeout`，失败用 `tccommon.RetryError()` 包装。
- gomonkey mock 单元测试覆盖业务逻辑（不执行 terraform 验收套件）。
- 代码可编译：CRUD 入参与云 API 接口入参严格对齐（Update 不传 ModifyDnsRecords 会忽略的出参字段）。

**Non-Goals:**
- 不修改/废弃既有 `tencentcloud_teo_dns_record` 资源。
- 不实现 DNS 记录批量启停（`ModifyDnsRecordsStatus`，不在本次参数映射内）。
- 不实现数据源（datasource）形态。
- 不新增 `_extension.go` 文件（非必须）。
- 不手改 `website/` 文档（由收尾阶段 `make doc` 生成）。
- 不在本阶段执行 `go build`/`go vet`/`gofmt`/`make doc`（由后续流程执行）。

## Decisions

### D1: 独立新资源，而非修改既有 `tencentcloud_teo_dns_record`

**选择**：新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_51.go`，导出 `ResourceTencentCloudTeoDnsRecord51()`。
**备选**：扩展现有资源 schema。
**理由**：需求明确要求新增 `tencentcloud_teo_dns_record_51`（RESOURCE_KIND_GENERAL），且向后兼容硬约束禁止改动既有资源 schema；新资源与旧资源并存，state 隔离，风险最小。

### D2: schema 顶层平铺，Optional+Computed 用于有服务端默认值的字段

字段清单（顺序与 Create 入参一致）：`zone_id`(Required,ForceNew)、`name`(Required)、`type`(Required)、`content`(Required)、`location`(Optional+Computed)、`ttl`(Optional+Computed)、`weight`(Optional+Computed)、`priority`(Optional+Computed)、`record_id`(Computed)、`status`(Computed)、`created_on`(Computed)、`modified_on`(Computed)。

- `zone_id` 设 ForceNew：记录属于特定站点，跨站点迁移无意义，且 Delete/Update API 均需 ZoneId。
- `location`/`ttl`/`weight`/`priority` 设 Optional+Computed：云 API 有默认值（Default/300/-1/0），未配置时由 Read 回填，避免 plan 抖动。
- `status` 仅 Computed：启停需要 `ModifyDnsRecordsStatus`，不在本次参数映射内，因此只读展示（云 API 注释明确 `Status` 在 ModifyDnsRecords 中作为入参会被忽略，传了也无效）。
- 不引入额外包裹层，符合"禁止 xxx_set/xxx_list 嵌套"的规则。

### D3: 复合 ID = `zoneId#recordId`，用 `tccommon.FILED_SP` 分隔

Create 成功后 `d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))`；Read/Update/Delete 中 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析（len!=2 报错 `id is broken,%s`）。import 使用 `terraform import tencentcloud_teo_dns_record_51.xxx zone-xxx#record-xxx`，文档中说明联合 id。

### D4: Read 复用服务层 `DescribeTeoDnsRecordById`，分页参数取云 API 标注最大值

服务层方法按 `id` 过滤器（`AdvancedFilter{Name:"id", Values:[recordId]}`）调用 `DescribeDnsRecords`，资源 Read 拿到 `*DnsRecord` 后逐字段判 nil 再 `d.Set(...)`。

**关键决策——Limit 取值**：`DescribeDnsRecordsRequest.Limit` 注释标注"默认值 20，上限 1000"。按规则"若查询接口中有分页字段，则给定值应该是云API接口注释中标注的最大值"，服务层查询将 `Limit` 设为 **1000**。当前 `DescribeTeoDnsRecordById` 未显式设置 Limit（走 API 默认 20），且以 `id` 精确过滤单条记录足够。为保证 id 查询的确定性并符合规则，新资源 Read 将在服务层新增/调整查询时设置 `Limit=1000`（新方法 `DescribeTeoDnsRecord51ById` 或在现有实现基础上为 id 过滤路径补充 Limit），不影响既有调用方语义。

具体落地：在 `service_tencentcloud_teo.go` 新增 `DescribeTeoDnsRecord51ById(ctx, zoneId, recordId)`，内部构造 `AdvancedFilter`（`Name=id`、`Values=[recordId]`）、`Limit=1000`、`Offset=0`，`resource.Retry(tccommon.ReadRetryTimeout)` 内调用 `DescribeDnsRecords`，返回首个匹配元素。既有 `DescribeTeoDnsRecordById` 保持不变（向后兼容）。

**读空处理**：若 `respData == nil`，按规则先 `log.Printf("[CRUD] teo_dns_record_51 id=%s", d.Id())` 保留现场，再 `d.SetId("")` 返回 nil（资源已删除场景）。

### D5: Update 仅调用 `ModifyDnsRecords`，只传可变字段

`mutableArgs = ["name", "type", "content", "location", "ttl", "weight", "priority"]`，任一 `d.HasChange` 即触发。构造单个 `teov20220901.DnsRecord` 元素：
- 传入 `RecordId`（从 `d.Id()` 解析）+ 上述可变字段。
- **不传** `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` 到 DnsRecord 元素——云 API 注释明确这些字段在 ModifyDnsRecords 中仅做出参、传入被忽略；请求顶层 `ZoneId` 才是站点定位参数。
- `request.DnsRecords = []*teov20220901.DnsRecord{dnsRecord}`（单条，低于 100 条上限）。
- `resource.Retry(tccommon.WriteRetryTimeout)` 包裹 `ModifyDnsRecordsWithContext`。
- 若无字段变化则跳过 API 调用，直接进入 Read 回刷。

### D6: Create 的空返回防御

按规则，Create 调用完成后必须检查：`response == nil` / `response.Response == nil` / `response.Response.RecordId == nil` / `*RecordId == ""`，任一命中则返回 `resource.NonRetryableError`（在 retry 块内检查，避免写入空 id）。检查前先 `log.Printf("[DEBUG]%s api[%s] ..., logId, d.Id()` 便于排障（Create 阶段 d.Id() 为空，打印 logId 与请求上下文）。`d.SetId(...)` 放在 retry 成功之后（retry 块外），随后调用 Read。

### D7: Delete 直接调用 `DeleteDnsRecords`

`ZoneId` + `RecordIds=[recordId]`（单条，远低于 1000 上限），`resource.Retry(tccommon.WriteRetryTimeout)` 包裹 `DeleteDnsRecordsWithContext`，检查 `response == nil || response.Response == nil` 返回 NonRetryableError。删除成功后返回 nil（无需 SetId("")，由 Terraform 框架处理）。

### D8: 非异步接口，无需 Read 轮询

vendor client.go 中四个接口均无"异步接口"标注，`ModifyDnsRecords`/`CreateDnsRecord`/`DeleteDnsRecords` 响应中也没有任务 ID 字段（仅 RequestId），因此不做调用后轮询；Create/Update 结束时调用一次 Read 回刷状态即可满足一致性要求。

### D9: 单元测试使用 gomonkey mock

新资源属于新增 terraform 资源，按规则使用 mock（gomonkey）方式对云 API 做处理，仅测试业务代码逻辑：
- mock `UseTeoV20220901Client` 返回空 `*teov20220901.Client`，再 ApplyMethodFunc mock `CreateDnsRecordWithContext`/`DescribeDnsRecords`/`ModifyDnsRecordsWithContext`/`DeleteDnsRecordsWithContext`。
- 用例：Create 成功（返回 RecordId，断言复合 ID 与 state 字段）、Create 空 RecordId 返回错误、Read 命中/未命中（未命中 SetId("")）、Update 触发 Modify（断言请求体字段）、Update 无变化不调 API、Delete 成功、复合 ID 解析错误分支。
- 参考 `resource_tc_teo_bind_security_template_test.go` 中 `mockMeta`/`gomonkey` 的既有写法，不执行 `go test`。

### D10: provider 注册与文档

- `tencentcloud/provider.go`：`"tencentcloud_teo_dns_record_51": teo.ResourceTencentCloudTeoDnsRecord51()`，放在既有 `tencentcloud_teo_dns_record` 注册项附近保持 teo 命名空间连续。
- `tencentcloud/provider.md`：在 TEO Resource 分类下按字母序插入 `tencentcloud_teo_dns_record_51`。
- `resource_tc_teo_dns_record_51.md`：一句话描述（带云产品名称 TEO，"Provides a resource to ..."）+ Example Usage（覆盖必填与可选字段）+ Import（说明使用联合 id `zoneId#recordId`）。不添加 Argument Reference / Attribute Reference（工具自动生成）。

## Risks / Trade-offs

- [Risk] `DescribeDnsRecords` 以 `id` 过滤器查询，`Filters.Values` 上限 20——单值精确过滤不受影响 → Mitigation: 服务层固定传单元素 `Values`。
- [Risk] 与既有 `tencentcloud_teo_dns_record` 资源管理同一云对象时可能产生双写竞争 → Mitigation: 文档不引导混用；两者 state 独立，terraform 层面不会自动冲突；属用户操作约束。
- [Risk] `Limit=1000` 为云 API 注释标注上限，若云端未来下调会报 InvalidParameter → Mitigation: 错误经 RetryError 包装透出，便于定位；当前值与注释一致。
- [Trade-off] `status` 只读（不接 `ModifyDnsRecordsStatus`）：该接口不在本次参数映射内，避免超出需求范围；用户如需启停仍可用旧资源。
- [Trade-off] 服务层新增 `DescribeTeoDnsRecord51ById` 而非复用旧方法：多一个近似方法，但保证旧调用方零影响且新方法显式 `Limit=1000` 符合分页规则。

## Migration Plan

- 纯新增资源，无 state 迁移、无回滚需求；如需回退，移除 provider 注册与三个新文件即可。
- 部署顺序：实现资源代码与测试 → 注册 provider → 补充 .md → （后续流程）gofmt / make doc / changelog。

## Open Questions

无：4 个云 API 均已在 vendor 中核对，字段映射与需求描述完全一致。
