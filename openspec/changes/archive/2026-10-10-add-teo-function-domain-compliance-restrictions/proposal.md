## Why

腾讯云 TEO 边缘函数（Edge Function）的默认域名（`Domain`）在部分地区会因合规问题被限制访问（如未备案、政府指令），SDK 的 `DescribeFunctions` 接口已在 `Function` 结构体中新增 `DomainComplianceRestrictions` 出参，返回"边缘函数默认域名因合规问题产生的地区访问限制列表"（每项含 `Reason` 限制原因、`Region` 被限制的国家/地区码）。当前 `tencentcloud_teo_function` 资源未暴露该信息，用户无法通过 Terraform state 获知函数默认域名在哪些地区因何种原因被限制访问，只能到控制台或调用 SDK 查询，破坏了以 Terraform state 作为单一事实来源的声明式运维体验。

## What Changes

- 为 `tencentcloud_teo_function` 资源新增一个 **Computed** 嵌套列表参数 `domain_compliance_restrictions`（TypeList of Resource），用于承载 `DescribeFunctions` 出参 `response.Functions[0].DomainComplianceRestrictions`：
  - 子字段 `reason`（Computed, TypeString）：下发访问限制的原因。枚举值：`ICP_RECORD_REQUIRED`（未备案）、`GOVERNMENT_ORDER`（政府指令）
  - 子字段 `region`（Computed, TypeString）：限制访问地区的具体国家/地区码，使用 ISO 3166 国家/地区代码标准（如 `CN`）
- 在 `resourceTencentCloudTeoFunctionRead` 中新增对该字段的回填逻辑：仅当 `DomainComplianceRestrictions` 非 nil 且长度大于 0 时构造 `[]map[string]interface{}` 并 `d.Set`，列表为空时不设置（保持 state 中为空列表）
- 在 `resource_tc_teo_function_test.go` 中补充 mock（gomonkey）单元测试用例，覆盖存在限制信息与不存在限制信息两种场景
- 更新资源示例文档 `resource_tc_teo_function.md`，展示新参数（由 `make doc` 收尾阶段同步生成 website 文档）

非破坏性变更：仅新增 Computed 只读字段，不改变任何现有字段的 schema 与行为，已有 TF 配置与 state 升级后无 plan drift；Create/Update/Delete 逻辑不受影响（云 API 的创建/修改接口均不接收该参数，属于纯出参）。

## Capabilities

### New Capabilities
- `teo-function-domain-compliance-restrictions`: `tencentcloud_teo_function` 资源暴露函数默认域名的合规访问限制信息（`domain_compliance_restrictions` 只读嵌套列表，含 `reason`/`region` 子字段），由 Read 阶段从 `DescribeFunctions` 出参回填

### Modified Capabilities
<!-- 无：本次仅新增一个 Computed 出参字段，不改变任何已有 spec 的需求级行为 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_function.go`（Schema 新增 `domain_compliance_restrictions` 嵌套 Computed 列表；Read 方法新增回填逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_function_test.go`（新增 gomonkey mock 单测用例）
  - `tencentcloud/services/teo/resource_tc_teo_function.md`（示例文档补充新参数说明）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.ComplianceRestriction`（`Reason`/`Region` 字段）与 `Function.DomainComplianceRestrictions`，无需变更 vendor
- 云 API：仅涉及 `DescribeFunctions`（`teo/v20220901`）的出参读取，不涉及任何写入接口参数
- 向后兼容：纯新增 Computed 字段，对现有配置零影响；不需要修改 `provider.go` / `provider.md`（资源已注册）
