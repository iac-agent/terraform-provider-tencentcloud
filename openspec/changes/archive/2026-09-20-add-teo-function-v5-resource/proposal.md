## Why

TEO（EdgeOne）边缘函数是边缘计算的核心能力，当前 provider 已有 `tencentcloud_teo_function` 资源，但其 Read 仅覆盖基础字段（function_id、name、remark、content、domain、create_time、update_time），未同步云 API 新增的 `DomainComplianceRestrictions`（默认域名因合规问题产生的地区访问限制列表）。为完整管理边缘函数生命周期并暴露合规限制信息，需要新增 `tencentcloud_teo_function_v5` 资源，复用同一套 CRUD 接口（CreateFunction / DescribeFunctions / ModifyFunction / DeleteFunction）并在 Read 中补充 `domain_compliance_restrictions` 字段。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_function_v5`（RESOURCE_KIND_GENERAL），文件 `tencentcloud/services/teo/resource_tc_teo_function_v5.go`，完整实现 Create / Read / Update / Delete 及 Import。
- Schema 入参：`zone_id`（必填，ForceNew）、`name`（必填）、`content`（必填）、`remark`（可选）。
- Schema 计算字段：`function_id`、`domain`、`domain_compliance_restrictions`（列表，含 `reason`、`region`）、`create_time`、`update_time`。
- 复合 ID：`zone_id#function_id`（使用 `tccommon.FILED_SP` 分隔）。
- Create 为异步接口：拿到 `function_id` 后先 `d.SetId()`，再通过 `DescribeFunctions` 轮询直到 `Domain` 字段返回（参考现有 `resourceTeoFunctionCreateStateRefreshFunc`）。
- Update 使用 `ModifyFunction`，仅 `remark`、`content` 可变；`name` 为不可变参数（immutableArgs）。
- Delete 使用 `DeleteFunction`。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。
- 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoFunctionV5ById` 服务方法（复用 `DescribeFunctions`，返回包含 `DomainComplianceRestrictions` 的 `Function` 对象）。
- 补充单元测试文件 `resource_tc_teo_function_v5_test.go`（使用 gomonkey mock 云 API，不使用 TF ACC 测试套件）。
- 补充文档 `tencentcloud/services/teo/resource_tc_teo_function_v5.md`。

## Capabilities

### New Capabilities
- `teo-function-v5-resource`: 新增 `tencentcloud_teo_function_v5` 资源，管理 EdgeOne 边缘函数完整生命周期，并在 Read 中暴露 `domain_compliance_restrictions` 合规限制信息。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增代码文件：`tencentcloud/services/teo/resource_tc_teo_function_v5.go`、`resource_tc_teo_function_v5_test.go`、`resource_tc_teo_function_v5.md`。
- 修改文件：`tencentcloud/provider.go`（注册资源）、`tencentcloud/provider.md`（文档）、`tencentcloud/services/teo/service_tencentcloud_teo.go`（新增 Describe 服务方法）。
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，已存在于 vendor。
- 向后兼容：新增资源，不影响现有 `tencentcloud_teo_function` 资源。
