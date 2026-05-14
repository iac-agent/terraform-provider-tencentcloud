## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘安全加速平台。当前已存在 `tencentcloud_teo_dns_record` 资源，使用 `CreateDnsRecord`/`DescribeDnsRecords`/`ModifyDnsRecords`/`DeleteDnsRecords` 等 API 管理 DNS 记录。本次需新增 `tencentcloud_teo_dns_record_20` 资源，作为新版本实现。

现有代码位于:
- 资源: `tencentcloud/services/teo/resource_tc_teo_dns_record.go`
- 服务: `tencentcloud/services/teo/service_tencentcloud_teo.go`

云 API (TEO v20220901) 接口:
- `CreateDnsRecord`: 创建 DNS 记录，入参 ZoneId/Name/Type/Content/Location/TTL/Weight/Priority，出参 RecordId
- `DescribeDnsRecords`: 查询 DNS 记录列表，入参 ZoneId/Offset/Limit/Filters/SortBy/SortOrder/Match，出参 TotalCount/DnsRecords
- `ModifyDnsRecords`: 批量修改 DNS 记录，入参 ZoneId/DnsRecords[]，其中 DnsRecord 结构体包含 RecordId/Name/Type/Content/Location/TTL/Weight/Priority
- `ModifyDnsRecordsStatus`: 修改记录状态，入参 ZoneId/RecordsToEnable/RecordsToDisable
- `DeleteDnsRecords`: 批量删除 DNS 记录，入参 ZoneId/RecordIds

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_20` Terraform 资源，支持 TEO DNS 记录完整 CRUD 生命周期
- 支持创建 DNS 记录并设置 zone_id、name、type、content、location、ttl、weight、priority 参数
- 支持读取 DNS 记录并同步 status、created_on、modified_on 等计算属性
- 支持更新 DNS 记录（通过 ModifyDnsRecords 修改记录属性，通过 ModifyDnsRecordsStatus 修改启用/停用状态）
- 支持删除 DNS 记录
- 支持 import 导入已存在的 DNS 记录
- 使用 zone_id + record_id 作为复合 ID
- 在 provider.go 和 provider.md 中注册新资源
- 新增 .md 文档和单元测试

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_dns_record` 资源的任何代码
- 不支持批量创建/管理多条 DNS 记录
- 不新增 _extension.go 文件

## Decisions

### 1. 资源命名与文件组织
**决策**: 新资源命名为 `tencentcloud_teo_dns_record_20`，文件为 `resource_tc_teo_dns_record_20.go`
**理由**: 遵循需求指定的命名规范，与现有 `tencentcloud_teo_dns_record` 资源共存但不互相影响

### 2. 复合 ID 设计
**决策**: 使用 `zone_id#record_id` 作为复合 ID，使用 `tccommon.FILED_SP` 作为分隔符
**理由**: 与现有 `tencentcloud_teo_dns_record` 资源保持一致，支持通过 zone_id + record_id 唯一定位一条 DNS 记录

### 3. Schema 设计
**决策**:
- 必填字段(Required): zone_id (ForceNew), name, type, content
- 可选字段(Optional): location, ttl, weight, priority, status
- 计算字段(Computed): record_id, created_on, modified_on
**理由**: 与 CreateDnsRecord API 入参对应，zone_id 为 ForceNew 因为创建后不可更改。status 字段虽由 ModifyDnsRecordsStatus 单独管理，但在 Terraform 层面统一为资源属性

### 4. 更新操作拆分
**决策**: Update 方法中分两步处理：
1. 当 name/type/content/location/ttl/weight/priority 变更时，调用 ModifyDnsRecords
2. 当 status 变更时，调用 ModifyDnsRecordsStatus
**理由**: 与现有 `tencentcloud_teo_dns_record` 实现一致，云 API 对属性修改和状态修改使用了不同的接口

### 5. Read 方法
**决策**: 在 service 层新增 `DescribeTeoDnsRecord20ById` 方法，通过 DescribeDnsRecords API + Filter(id=recordId) 查询特定记录
**理由**: DescribeDnsRecords 是列表查询 API，需通过 Filter 过滤获取特定记录，与现有实现模式一致

### 6. 测试策略
**决策**: 使用 gomonkey mock 方式编写单元测试
**理由**: 新增资源按照要求使用 mock 方式进行业务逻辑单元测试，不依赖 Terraform 测试套件

## Risks / Trade-offs

- [风险] 现有 `tencentcloud_teo_dns_record` 已存在，用户可能混淆两个资源 → 通过文档说明这是新版本实现，命名 _20 后缀明确区分
- [风险] DescribeDnsRecords 为列表查询接口，需通过 Filter 精确定位 → 使用 id 精确过滤，限制 Limit=1，确保查询效率
- [风险] DnsRecord 结构体中 ZoneId/Status/CreatedOn/ModifiedOn 字段在 ModifyDnsRecords 中仅做出参使用 → 在 Update 方法中不将这些字段设置为 ModifyDnsRecords 的入参
