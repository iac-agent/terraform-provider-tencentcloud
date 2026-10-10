## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_function.go` 的 `ResourceTencentCloudTeoFunction()` Schema 中，于 `domain` 与 `create_time` 字段之间新增 `domain_compliance_restrictions` 参数：`Type: schema.TypeList`、`Computed: true`、`Elem: &schema.Resource{ Schema: map[string]*schema.Schema{...} }`，子字段 `reason` 与 `region` 均为 `Type: schema.TypeString`、`Computed: true`
- [x] 1.2 为新参数及子字段编写 Description：外层说明"Region access restriction list of the function default domain due to compliance issues. Valid restriction reasons: `ICP_RECORD_REQUIRED` (ICP filing not obtained), `GOVERNMENT_ORDER` (government order)."；`reason` 说明限制原因枚举值；`region` 说明使用 ISO 3166 国家/地区代码标准
- [x] 1.3 确认新参数不含 `Optional`/`Required`，不进入 Update 的 `mutableArgs` / `immutableArgs` 数组（纯 Computed 出参，Create/Update/Delete 请求体零改动）

## 2. Read 回填逻辑

- [x] 2.1 在 `resourceTencentCloudTeoFunctionRead` 中（`respData.Domain` 回填之后、`respData.CreateTime` 回填之前）新增：当 `respData.DomainComplianceRestrictions != nil && len(respData.DomainComplianceRestrictions) > 0` 时，遍历构造 `[]map[string]interface{}`，每个元素仅在 `item.Reason != nil` 时写 `map["reason"] = *item.Reason`、仅在 `item.Region != nil` 时写 `map["region"] = *item.Region`，最后 `_ = d.Set("domain_compliance_restrictions", restrictionsList)`（参考 design.md 中摘录的 `work_mode_infos` 回填模式）
- [x] 2.2 确认回填前的 `respData == nil` 判空、现有字段回填顺序与逻辑均未被改动，资源 ID（`zoneId#functionId`）解析逻辑保持不变

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中按既有 mock（gomonkey）风格新增用例：mock `DescribeTeoFunctionById`（或底层 `DescribeFunctions` API）返回携带两个 `ComplianceRestriction`（`{Reason: "ICP_RECORD_REQUIRED", Region: "CN"}`、`{Reason: "GOVERNMENT_ORDER", Region: "US"}`）的 `Function`，断言 Read 后 state 中 `domain_compliance_restrictions` 内容与 mock 数据一致
- [x] 3.2 新增用例：mock 返回的 `Function.DomainComplianceRestrictions` 为 nil（空），断言 Read 成功且 `domain_compliance_restrictions` 为空列表
- [x] 3.3 保证现有测试用例不受影响（不修改其行为），新增用例仅使用 mock，不依赖真实云凭证

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_function.md` 中补充 `domain_compliance_restrictions` 的说明或在示例下方展示该 Computed 字段的输出结构（说明 reason/region 的取值含义）
- [x] 4.2 收尾阶段执行 `make doc` 重新生成 `website/docs/` 下对应资源文档（禁止手改 website/ 目录；禁止在本阶段之前执行 make doc）——apply 阶段不执行 `make doc`，该操作由收尾阶段的 tfpacer-finalize skill 统一执行；`tencentcloud/services/teo/resource_tc_teo_function.md`（make doc 的文档源）已于任务 4.1 更新完毕

## 5. 验证

- [x] 5.1 检查 `vendor` 中 `teov20220901.ComplianceRestriction` 字段名拼写（`Reason`/`Region`）与代码引用一致，确保可正确编译（本阶段不执行 go build / go vet）
- [x] 5.2 对照 openspec/specs 变更后的 spec（`teo-function-domain-compliance-restrictions`）中各 Scenario 逐条核对实现与测试覆盖情况
- [x] 5.3 确认 `provider.go` / `provider.md` 无需改动（资源 `tencentcloud_teo_function` 已注册）
