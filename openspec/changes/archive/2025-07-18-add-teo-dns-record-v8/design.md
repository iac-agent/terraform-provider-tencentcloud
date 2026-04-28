## Context

TEO（EdgeOne）是腾讯云的边缘安全加速平台，提供 DNS 解析管理能力。当前 provider 中已存在 `tencentcloud_teo_dns_record` 资源，基于相同的云 API 接口（CreateDnsRecord/DescribeDnsRecords/ModifyDnsRecords/DeleteDnsRecords/ModifyDnsRecordsStatus）。现需新增 `tencentcloud_teo_dns_record_v8` 作为独立的新版本资源，与旧版资源并存，避免对现有用户的 TF 配置和 state 造成破坏性变更。

现有资源代码位于 `tencentcloud/services/teo/resource_tc_teo_dns_record.go`，使用复合 ID（zoneId#recordId），支持 Import，通过 `TeoService.DescribeTeoDnsRecordById` 进行查询。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_v8` 资源，提供完整的 DNS 记录 CRUD 管理
- 与现有 `tencentcloud_teo_dns_record` 资源独立并存，互不影响
- 支持所有云 API 支持的 DNS 记录字段和状态管理
- 提供完整的单元测试和文档

**Non-Goals:**
- 不修改或弃用现有 `tencentcloud_teo_dns_record` 资源
- 不更改现有资源的 schema 或行为
- 不引入新的云 API 接口

## Decisions

### 1. 资源命名和文件组织
**决策**: 新资源命名为 `tencentcloud_teo_dns_record_v8`，文件为 `resource_tc_teo_dns_record_v8.go`
**理由**: 遵循版本化资源的命名惯例（参考 `tencentcloud_teo_l7_acc_rule_v2`），作为独立新资源与旧版本并存

### 2. Schema 设计
**决策**: Schema 与现有 dns_record 资源保持一致的字段结构
- Required + ForceNew: `zone_id`
- Required: `name`, `type`, `content`
- Computed + Optional: `location`, `ttl`（TypeInt）, `weight`（TypeInt）, `priority`（TypeInt）, `status`
- Computed: `created_on`, `modified_on`

**理由**: v8 版本使用相同的云 API 接口，字段一致。整数类型字段使用 `d.GetOkExists()` 以支持零值。

### 3. 复合 ID 设计
**决策**: 使用 `zoneId#recordId` 作为复合 ID（tccommon.FILED_SP 分隔符）
**理由**: 与现有 dns_record 资源保持一致，支持通过 zoneId 和 recordId 联合定位资源

### 4. 服务层方法
**决策**: 在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecordV8ById` 方法
**理由**: 虽然查询逻辑与 `DescribeTeoDnsRecordById` 相同，但独立方法确保 v8 资源与旧版本解耦，便于后续独立演进

### 5. Update 操作的拆分
**决策**: Update 操作拆分为两部分
- 第一部分：使用 `ModifyDnsRecords` 修改 name/type/content/location/ttl/weight/priority
- 第二部分：使用 `ModifyDnsRecordsStatus` 修改 status（enable/disable）
**理由**: 与现有资源保持一致，状态修改需要独立的 API 调用

### 6. 测试策略
**决策**: 使用 gomonkey mock 方式编写单元测试，不使用 terraform 测试套件
**理由**: 新增资源使用 mock 方式进行业务逻辑的单元测试，避免依赖云 API 环境

### 7. SDK Client 选择
**决策**: 使用 `UseTeoV20220901Client()` 配合 `*WithContext` 方法
**理由**: 与现有 dns_record 资源保持一致

## Risks / Trade-offs

- **[风险] 两套 DNS 记录资源共存** → 用户可能混淆 `tencentcloud_teo_dns_record` 和 `tencentcloud_teo_dns_record_v8`，通过文档说明两者关系和使用场景来缓解
- **[风险] DescribeDnsRecords 使用模糊匹配** → 查询接口的 id 过滤器是模糊匹配，可能返回多余记录。在 `DescribeTeoDnsRecordV8ById` 中需对返回结果进行精确匹配校验
- **[权衡] v8 与旧版资源 API 调用逻辑相同** → 代码复用度低但保证了版本独立性和向后兼容
