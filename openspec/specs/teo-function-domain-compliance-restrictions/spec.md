# teo-function-domain-compliance-restrictions Specification

## Purpose
定义 `tencentcloud_teo_function` 资源对边缘函数默认域名合规限制信息（`domain_compliance_restrictions`）的暴露能力：作为纯 Computed 出参承接 `DescribeFunctions` 返回的 `DomainComplianceRestrictions` 列表（含 `reason` 限制原因、`region` 限制地区码），仅在 Read 阶段回填 state，不参与 Create/Update 请求，并以 gomonkey mock 单元测试覆盖回填逻辑。
## Requirements
### Requirement: teo_function 资源暴露默认域名合规限制参数
系统 SHALL 在 `tencentcloud_teo_function` 资源 schema 顶层定义 `domain_compliance_restrictions` 参数（`TypeList`、`Computed: true`、`Elem: &schema.Resource{...}`），用于承接 `DescribeFunctions` 出参 `response.Functions[0].DomainComplianceRestrictions` 列表。该参数 SHALL NOT 可被用户配置（无 Optional/Required），且元素结构 SHALL 包含两个 Computed 子字段：
- `reason`（TypeString）：下发访问限制的原因，枚举值为 `ICP_RECORD_REQUIRED`（未备案）、`GOVERNMENT_ORDER`（政府指令）
- `region`（TypeString）：限制访问地区的具体国家/地区码，使用 ISO 3166 国家/地区代码标准

#### Scenario: Schema 定义嵌套 Computed 列表参数
- **WHEN** `ResourceTencentCloudTeoFunction()` 的 schema 被定义
- **THEN** `domain_compliance_restrictions` SHALL 为 `schema.TypeList` 且 `Computed: true`，其 `Elem` SHALL 为包含 `reason`、`region` 两个 `schema.TypeString` Computed 子字段的 `schema.Resource`

#### Scenario: 参数对用户只读
- **WHEN** 用户在 Terraform 配置中为 `tencentcloud_teo_function` 声明 `domain_compliance_restrictions` 的值
- **THEN** terraform plan SHALL 报错拒绝该配置（Computed 字段不可写）

### Requirement: Read 阶段回填合规限制列表
系统 SHALL 在 `resourceTencentCloudTeoFunctionRead` 中，当 `DescribeTeoFunctionById` 返回的 `Function.DomainComplianceRestrictions` 非 nil 且长度大于 0 时，遍历该列表构造 `[]map[string]interface{}`（每个元素在对应子字段非 nil 时写入 `reason`/`region` 键），并整体 `d.Set("domain_compliance_restrictions", ...)` 写入 state；当列表为 nil 或为空时 SHALL NOT 调用 `d.Set`，state 中该字段保持为空列表。

#### Scenario: 云端返回多个合规限制
- **WHEN** `DescribeFunctions` 返回的 `Functions[0].DomainComplianceRestrictions` 包含两个元素 `[{Reason: "ICP_RECORD_REQUIRED", Region: "CN"}, {Reason: "GOVERNMENT_ORDER", Region: "US"}]`
- **THEN** state 中 `domain_compliance_restrictions` SHALL 为包含两个元素的列表，第一项 `reason = "ICP_RECORD_REQUIRED"`、`region = "CN"`，第二项 `reason = "GOVERNMENT_ORDER"`、`region = "US"`

#### Scenario: 云端未返回合规限制
- **WHEN** `DescribeFunctions` 返回的 `Functions[0].DomainComplianceRestrictions` 为 nil 或空列表
- **THEN** Read SHALL 不调用 `d.Set("domain_compliance_restrictions", ...)`，state 中该字段为空列表且不产生错误

#### Scenario: 列表元素子字段为 nil
- **WHEN** `DomainComplianceRestrictions` 中某元素的 `Reason` 或 `Region` 字段为 nil
- **THEN** 构造 map 时 SHALL 跳过该 nil 子字段（不向 map 写入该键），整体列表仍正常 set

#### Scenario: 已有配置无 plan drift
- **WHEN** 用户在 Provider 升级前已存在的 `tencentcloud_teo_function` 配置（未引用新参数）执行 `terraform plan`
- **THEN** plan SHALL 不因新增 `domain_compliance_restrictions` 字段而显示任何变更（纯 Computed 字段向后兼容）

### Requirement: 合规限制不参与写入接口调用
系统 SHALL NOT 在 `resourceTencentCloudTeoFunctionCreate`（`CreateFunction` 请求）或 `resourceTencentCloudTeoFunctionUpdate`（`ModifyFunction` 请求）中传递 `DomainComplianceRestrictions` 相关参数；`domain_compliance_restrictions` SHALL NOT 出现在 Update 的 `mutableArgs` 或 `immutableArgs` 数组中。

#### Scenario: Create 与 Update 不受影响
- **WHEN** 用户创建或更新 `tencentcloud_teo_function` 资源
- **THEN** 发出的 `CreateFunction` / `ModifyFunction` 请求 SHALL 与变更前完全一致，不包含合规限制相关入参

### Requirement: 单元测试覆盖合规限制回填逻辑
系统 SHALL 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中使用 gomonkey mock 云 API 的方式补充单元测试用例，覆盖"存在合规限制信息"与"不存在合规限制信息（nil/空列表）"两类 Read 回填场景，且 SHALL NOT 依赖真实云环境。

#### Scenario: mock 存在合规限制的读取
- **WHEN** 测试 mock `DescribeFunctions` 返回的 `Function` 携带非空 `DomainComplianceRestrictions`
- **THEN** 断言 Read 后资源 state 中 `domain_compliance_restrictions` 的内容与 mock 数据一致

#### Scenario: mock 不存在合规限制的读取
- **WHEN** 测试 mock `DescribeFunctions` 返回的 `Function` 的 `DomainComplianceRestrictions` 为 nil
- **THEN** 断言 Read 成功且 state 中 `domain_compliance_restrictions` 为空列表

