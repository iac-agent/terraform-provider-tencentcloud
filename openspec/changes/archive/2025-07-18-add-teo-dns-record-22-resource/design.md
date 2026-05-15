## Context

当前 Terraform Provider 中已存在 `tencentcloud_teo_dns_record` 资源，用于管理 TEO (EdgeOne) DNS 记录的 CRUD 操作。现需要新增 `tencentcloud_teo_dns_record_22` 资源，该资源使用相同的云 API 接口（CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、ModifyDnsRecordsStatus、DeleteDnsRecords），但作为独立的资源实现。

云 API 关键信息：
- `CreateDnsRecord`: 创建 DNS 记录，入参 ZoneId/Name/Type/Content/Location/TTL/Weight/Priority，出参 RecordId
- `DescribeDnsRecords`: 查询 DNS 记录列表，入参 ZoneId/Offset/Limit/Filters/SortBy/SortOrder/Match，出参 TotalCount/DnsRecords
- `ModifyDnsRecords`: 批量修改 DNS 记录，入参 ZoneId/DnsRecords（DnsRecord 包含 RecordId/Name/Type/Content/Location/TTL/Weight/Priority）
- `ModifyDnsRecordsStatus`: 修改 DNS 记录状态，入参 ZoneId/RecordsToEnable/RecordsToDisable
- `DeleteDnsRecords`: 批量删除 DNS 记录，入参 ZoneId/RecordIds
- DnsRecord 结构体还包含只读字段：ZoneId、Status、CreatedOn、ModifiedOn（在 ModifyDnsRecords 中不能作为入参）

现有资源参考：`tencentcloud_teo_dns_record`（位于 `tencentcloud/services/teo/resource_tc_teo_dns_record.go`），使用 zone_id + record_id 的联合 ID 格式。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_22` 资源，支持 DNS 记录的完整 CRUD 生命周期管理
- 支持 Create（创建）、Read（读取）、Update（更新）、Delete（删除）操作
- 支持导入（Import）
- 在 provider.go 中注册新资源
- 生成对应的 .md 文档
- 编写单元测试（使用 gomonkey mock 方式）

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_dns_record` 资源的行为
- 不新增数据源（data source）
- 不支持异步操作轮询（所有接口均为同步接口）

## Decisions

1. **资源 Schema 设计**：参考现有 `tencentcloud_teo_dns_record` 资源的 Schema 设计，保持参数一致性
   - `zone_id`: Required, ForceNew - 站点 ID
   - `name`: Required - DNS 记录名称
   - `type`: Required - DNS 记录类型
   - `content`: Required - DNS 记录内容
   - `location`: Optional, Computed - 解析线路
   - `ttl`: Optional, Computed - 缓存时间
   - `weight`: Optional, Computed - 权重
   - `priority`: Optional, Computed - MX 优先级
   - `status`: Optional, Computed - 解析状态
   - `created_on`: Computed - 创建时间
   - `modified_on`: Computed - 修改时间
   - `record_id`: Computed - DNS 记录 ID

2. **复合 ID 格式**：使用 `zone_id#record_id` 的联合 ID 格式（使用 `tccommon.FILED_SP` 分隔符），与现有 `tencentcloud_teo_dns_record` 资源保持一致

3. **Update 逻辑**：分两部分
   - 可变参数变更（name/type/content/location/ttl/weight/priority）：调用 ModifyDnsRecords
   - 状态变更（status）：调用 ModifyDnsRecordsStatus

4. **Read 逻辑**：调用 DescribeDnsRecords，通过 Filters 按 RecordId 过滤获取单条记录

5. **测试方式**：使用 gomonkey mock 方式编写单元测试，不使用 terraform 测试套件

## Risks / Trade-offs

- **命名冲突风险**：`tencentcloud_teo_dns_record_22` 中的 "22" 后缀可能造成用户困惑 → 在文档中明确说明资源的用途
- **DescribeDnsRecords 无单条查询接口**：DescribeDnsRecords 返回的是列表，需要通过 Filters 过滤 RecordId 来获取单条记录 → 在 service 层封装 DescribeTeoDnsRecord22ById 方法，自动处理分页和过滤逻辑
