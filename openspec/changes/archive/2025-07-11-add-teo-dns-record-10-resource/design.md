## Context

TEO（EdgeOne）产品需要新增 `tencentcloud_teo_dns_record_10` Terraform 资源，用于管理 DNS 记录的完整 CRUD 生命周期。已有的 `tencentcloud_teo_dns_record` 资源提供了参考实现，新资源将遵循相同的架构模式，但作为独立的资源实现。

当前状态：
- 已有 `tencentcloud_teo_dns_record` 资源，使用 CreateDnsRecord/DescribeDnsRecords/ModifyDnsRecords/DeleteDnsRecords 接口
- 云 API SDK 中所有四个 CRUD 接口均已可用
- DnsRecord 模型包含 ZoneId、RecordId、Name、Type、Content、Location、TTL、Weight、Priority、Status、CreatedOn、ModifiedOn 字段

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_10` 资源，支持完整的 Create/Read/Update/Delete 操作
- 使用复合 ID（zone_id + FILED_SP + record_id）标识资源
- 支持 Terraform Import 功能
- 遵循现有 `tencentcloud_teo_dns_record` 资源的代码风格和架构模式
- 参考 `tencentcloud_igtm_strategy` 资源的代码生成样式风格

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_dns_record` 资源
- 不实现 DNS 记录状态的启用/停用（ModifyDnsRecordsStatus），该功能由现有资源处理
- 不添加除 CRUD 之外的额外功能

## Decisions

### 1. 资源 ID 设计
**决策**: 使用复合 ID 格式 `zone_id#record_id`（FILED_SP 分隔符）
**理由**: Read/Delete 接口需要 zone_id 和 record_id 两个参数，复合 ID 可在 Read/Delete 方法中解析出两个参数。与现有 `tencentcloud_teo_dns_record` 资源保持一致。

### 2. Create 接口参数映射
**决策**: CreateDnsRecord 接口入参全部映射为 Terraform Schema 参数
- ZoneId → zone_id (Required, ForceNew)
- Name → name (Required)
- Type → type (Required)
- Content → content (Required)
- Location → location (Optional, Computed)
- TTL → ttl (Optional, Computed)
- Weight → weight (Optional, Computed)
- Priority → priority (Optional, Computed)

**理由**: 与云 API 接口入参完全对应。zone_id 为 ForceNew 因为变更站点等同于创建新资源。出参 RecordId 存储为 record_id（Computed 字段），作为复合 ID 的一部分。

### 3. Read 接口实现
**决策**: 通过 DescribeDnsRecords 接口 + Filter（id=recordId）查询单条记录
**理由**: 没有单条记录查询接口，需要使用列表接口配合过滤条件。在 service 层实现 DescribeTeoDnsRecord10ById 方法，与现有 DescribeTeoDnsRecordById 方法类似。

### 4. Update 接口实现
**决策**: 使用 ModifyDnsRecords 接口，构建 DnsRecord 对象传入 RecordId 和可修改字段
**理由**: ModifyDnsRecords 接受 DnsRecord 列表，每次只修改一条记录。可修改字段包括 name、type、content、location、ttl、weight、priority。

### 5. Delete 接口实现
**决策**: 使用 DeleteDnsRecords 接口，传入 ZoneId 和 RecordIds
**理由**: DeleteDnsRecords 接受 RecordIds 列表，删除单条记录时传入 [recordId]。

### 6. 额外 Computed 字段
**决策**: 从 Read 响应中额外读取 status、created_on、modified_on 字段
**理由**: DnsRecord 模型返回这些字段，作为 Computed 属性供用户查看，但不参与 Create/Update 操作。

### 7. 单元测试策略
**决策**: 使用 gomonkey mock 方式进行单元测试
**理由**: 按照规范要求，新增资源使用 mock（gomonkey）方法对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试。

## Risks / Trade-offs

- [风险] DnsRecord 模型中 ZoneId 字段在 ModifyDnsRecords 中不能作为入参 → 缓解：Update 方法中不设置 DnsRecord 的 ZoneId 字段
- [风险] DescribeDnsRecords 是列表接口，可能返回多条记录 → 缓解：使用 id 过滤条件精确匹配，只取第一条结果
- [风险] 与现有 tencentcloud_teo_dns_record 资源功能重叠 → 缓解：作为独立资源，不影响现有资源的兼容性
