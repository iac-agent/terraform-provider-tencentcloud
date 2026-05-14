## Context

TEO（EdgeOne）是腾讯云的边缘安全加速平台，支持通过云API管理 DNS 记录。当前已存在 `tencentcloud_teo_dns_record` 资源，但需要新增 `tencentcloud_teo_dns_record_18` 资源以提供更新的实现。新资源基于以下云API接口：

- `CreateDnsRecord`：创建单条 DNS 记录，入参为扁平字段（ZoneId, Name, Type, Content, Location, TTL, Weight, Priority），出参返回 RecordId
- `DescribeDnsRecords`：查询 DNS 记录列表，支持分页和过滤，出参返回 DnsRecords 列表
- `ModifyDnsRecords`：批量修改 DNS 记录，入参包含 DnsRecords 子结构体列表
- `DeleteDnsRecords`：批量删除 DNS 记录，入参包含 RecordIds 列表

现有 TeoService 中已有 `DescribeTeoDnsRecordById` 方法，可直接复用。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_18` 资源的完整 CRUD 操作
- 遵循 RESOURCE_KIND_GENERAL 模式，参考现有 `tencentcloud_igtm_strategy` 和 `tencentcloud_teo_dns_record` 资源
- 支持资源导入（Import），使用 zone_id 和 record_id 作为复合 ID
- 复用现有的 TeoService 服务层方法
- 提供基于 gomonkey 的单元测试

**Non-Goals:**
- 不修改已有的 `tencentcloud_teo_dns_record` 资源
- 不实现 DescribeDnsRecords 的 filters、sort_by、sort_order、match 等查询参数作为资源字段（这些是查询辅助参数，不属于资源本身属性）
- 不支持异步操作轮询（所有接口均为同步接口）

## Decisions

### 1. 资源 Schema 设计

**决策**：Schema 字段与 CreateDnsRecord 入参对齐，zone_id 设为 Required + ForceNew（不可变），name/type/content 设为 Required（可更新），location/ttl/weight/priority 设为 Optional + Computed（可更新），record_id/status/created_on/modified_on 设为 Computed（只读）。

**理由**：zone_id 是资源的父级标识，创建后不可变更，因此 ForceNew: true。name/type/content 是 DNS 记录的核心属性，必须提供且可更新。location/ttl/weight/priority 有服务端默认值，设为 Computed 以便读取默认值。record_id 由创建接口返回，status/created_on/modified_on 由查询接口返回。

### 2. 复合 ID 格式

**决策**：使用 `zone_id` + `FILED_SP` + `record_id` 作为复合 ID。

**理由**：与现有 `tencentcloud_teo_dns_record` 资源保持一致，且 DescribeDnsRecords 接口需要 ZoneId 作为查询条件，RecordId 在 DnsRecord 子结构中返回，两者缺一不可。

### 3. Read 操作实现

**决策**：复用 TeoService 中现有的 `DescribeTeoDnsRecordById` 方法，该方法内部调用 DescribeDnsRecords 接口并按 RecordId 过滤返回单条记录。

**理由**：避免重复代码，保持与现有资源的一致性。

### 4. Update 操作实现

**决策**：调用 ModifyDnsRecords 接口，构造包含单条 DnsRecord 的请求体（包含 RecordId 和所有可更新字段）。需要检测 name/type/content/location/ttl/weight/priority 字段的变化。

**理由**：ModifyDnsRecords 接口支持批量修改，但对于单个 Terraform 资源，只需修改一条记录。所有可更新字段都需要在修改请求中提供（即使未变更），以确保 API 调用的一致性。

### 5. Delete 操作实现

**决策**：调用 DeleteDnsRecords 接口，传入 RecordIds 列表（包含单条 record_id）。

**理由**：DeleteDnsRecords 接口支持批量删除，但单个资源只需删除一条记录。

## Risks / Trade-offs

- [风险] DescribeDnsRecords 是列表接口，通过 RecordId 过滤单条记录可能存在性能问题 → 缓解：现有服务层已实现该方法，且 API 支持按 ID 过滤，性能影响可控
- [风险] ModifyDnsRecords 是批量接口，单条记录更新时需要构造 DnsRecord 子结构体 → 缓解：只需包含一条记录，实现简单
- [风险] 新资源名称 `tencentcloud_teo_dns_record_18` 与现有 `tencentcloud_teo_dns_record` 存在命名相似性，可能引起用户混淆 → 缓解：在文档中明确说明新资源的用途
