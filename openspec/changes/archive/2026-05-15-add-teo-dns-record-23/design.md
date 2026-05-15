## Context

Terraform Provider for TencentCloud 已有 `tencentcloud_teo_dns_record` 和 `tencentcloud_teo_dns_record_21` 两个 TEO DNS 记录资源实现。现需新增 `tencentcloud_teo_dns_record_23` 资源，基于相同的云 API 接口（CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords），遵循现有代码模式。

当前状态：
- 已有 `resource_tc_teo_dns_record_21.go` 作为最近的参考实现
- 已有 `TeoService.DescribeTeoDnsRecordById` 服务层方法可复用
- 云 API SDK 位于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_23` 资源，支持 DNS 记录的完整 CRUD 生命周期
- 遵循 `resource_tc_teo_dns_record_21.go` 的代码模式和风格
- 支持 zone_id、name、type、content、location、ttl、weight、priority 等参数
- 支持 import 功能（使用 zone_id#record_id 复合 ID）
- 包含单元测试和文档

**Non-Goals:**
- 不修改现有的 `tencentcloud_teo_dns_record` 或 `tencentcloud_teo_dns_record_21` 资源
- 不支持 status 字段的更新（status 为只读 computed 字段）
- 不新增服务层方法（复用已有的 DescribeTeoDnsRecordById）

## Decisions

### 1. 资源 ID 格式
**决策**: 使用 `zone_id` + `tccommon.FILED_SP` + `record_id` 作为复合 ID

**理由**: 与现有 teo_dns_record_21 资源保持一致。DescribeDnsRecords 接口需要 zone_id 来过滤查询，且 record_id 需要作为唯一标识。

### 2. Read 方法复用服务层
**决策**: 复用已有的 `TeoService.DescribeTeoDnsRecordById()` 方法

**理由**: 该方法已通过 DescribeDnsRecords 接口 + Filter(id=recordId) 实现了精确查询，无需重复实现。

### 3. Update 方法使用 ModifyDnsRecords
**决策**: 当 name、type、content、location、ttl、weight、priority 任一字段变更时，调用 ModifyDnsRecords 接口

**理由**: ModifyDnsRecords 接口接受 DnsRecord 数组，包含 RecordId 和所有可修改字段。status 字段为 computed 只读，不参与更新逻辑。

### 4. Schema 设计
**决策**: 参照 teo_dns_record_21 的 schema 设计：
- `zone_id`: Required, ForceNew
- `name`: Required
- `type`: Required
- `content`: Required
- `location`: Optional + Computed
- `ttl`: Optional + Computed
- `weight`: Optional + Computed
- `priority`: Optional + Computed
- `record_id`: Computed
- `status`: Computed
- `created_on`: Computed
- `modified_on`: Computed

**理由**: 与 teo_dns_record_21 保持一致，确保向后兼容性和一致性。

### 5. 测试策略
**决策**: 使用 gomonkey mock 方式编写单元测试

**理由**: 按照项目规范，新增资源不使用 Terraform 测试套件，而是使用 mock（gomonkey）进行业务逻辑单元测试。

## Risks / Trade-offs

- [云 API 兼容性] → 已验证 vendor 中云 API 接口支持所需的 CRUD 操作，风险可控
- [与现有资源的差异] → teo_dns_record_23 与 teo_dns_record_21 功能基本一致，仅名称不同，这是为了满足版本迭代需求
