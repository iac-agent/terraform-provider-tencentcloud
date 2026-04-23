## Context

TEO (TencentCloud EdgeOne) 已有 `tencentcloud_teo_dns_record` 资源用于管理 DNS 记录，但该资源采用旧版代码风格。本次变更新增 `tencentcloud_teo_dns_record_v2` 资源，采用新版代码风格（参考 `tencentcloud_igtm_strategy`），提供标准化的 CRUD 流程、一致的重试处理和更好的代码组织。

云 API 接口信息：
- `CreateDnsRecord`：创建 DNS 记录，入参包括 ZoneId、Name、Type、Content、Location、TTL、Weight、Priority，出参包括 RecordId
- `DescribeDnsRecords`：查询 DNS 记录列表，使用 AdvancedFilter 进行过滤，支持分页（Offset/Limit，最大 1000）
- `ModifyDnsRecords`：批量修改 DNS 记录，入参为 ZoneId 和 DnsRecords 列表，DnsRecord 中可修改字段为 RecordId、Name、Type、Content、Location、TTL、Weight、Priority
- `DeleteDnsRecords`：批量删除 DNS 记录，入参为 ZoneId 和 RecordIds 列表

DnsRecord 结构体包含的只读字段（DescribeDnsRecords 返回但 ModifyDnsRecords 忽略）：ZoneId、Status、CreatedOn、ModifiedOn。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_v2` 资源，遵循 `tencentcloud_igtm_strategy` 的代码风格
- 实现完整的 CRUD 操作（Create、Read、Update、Delete）
- 支持 Terraform Import
- 资源 ID 采用复合格式 `zone_id#record_id`
- 在 Read 操作中使用 AdvancedFilter 按 record id 过滤查询
- 在 Update 操作中仅发送变更的字段
- 补充 gomonkey 单元测试
- 生成 `.md` 文档文件

**Non-Goals:**
- 不修改或删除现有的 `tencentcloud_teo_dns_record` 资源
- 不实现 DNS 记录状态切换功能（ModifyDnsRecordsStatus API 不在本资源范围内）
- 不实现数据源（datasource）

## Decisions

### 1. 资源 ID 格式采用复合 ID：`zone_id#record_id`
**理由**：与现有 `tencentcloud_teo_dns_record` 资源保持一致的 ID 格式，使用 `tccommon.FILED_SP` 作为分隔符。在 Read/Update/Delete 中从 `d.Get()` 获取 zone_id 和 record_id，而非直接拆分 `d.Id()`。

### 2. Read 操作使用 AdvancedFilter 按 id 过滤
**理由**：`DescribeDnsRecords` API 支持 AdvancedFilter 过滤，其中 `id` 键支持按记录 ID 精确查询。这与现有 `tencentcloud_teo_dns_record` 的实现一致，使用 `DescribeTeoDnsRecordById` 服务层方法。

### 3. Update 操作使用 ModifyDnsRecords API
**理由**：`ModifyDnsRecords` API 接受 DnsRecords 列表，其中每条记录需包含 RecordId 和要修改的字段。本资源每次只管理一条记录，因此列表中只有一条 DnsRecord。只读字段（ZoneId、Status、CreatedOn、ModifiedOn）不作为 Update 的入参。

### 4. 采用新版代码风格（参考 igtm_strategy）
**理由**：
- 使用 `tccommon.NewResourceLifeCycleHandleFuncContext` 创建上下文
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行写操作重试
- 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 进行读操作重试
- 使用 `tccommon.RetryError(e)` 包装错误

### 5. zone_id 字段设为 ForceNew
**理由**：DNS 记录绑定到特定 Zone，Zone 变更意味着需要创建新记录。

### 6. 不包含 status 字段作为可写参数
**理由**：与用户需求中列出的接口参数一致，status 是 DescribeDnsRecords 返回的只读字段，本资源不管理 DNS 记录的启停状态。status 作为 computed 字段在 Read 中从 API 响应回填。

## Risks / Trade-offs

- [与旧资源共存] → 两个资源管理同一云资源，用户需明确选择使用哪个版本。旧资源保持不变，不影响现有用户。
- [DescribeDnsRecords 分页] → Read 操作使用 AdvancedFilter 按 id 过滤，通常只返回一条记录，不会触发分页问题。
- [ModifyDnsRecords 批量接口] → 本资源仅管理单条记录，每次 Update 只传一条 DnsRecord，不影响批量场景。
