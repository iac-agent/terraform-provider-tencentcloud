## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_function.go` 的 `ResourceTencentCloudTeoFunction()` schema 中新增顶层 `domain_compliance_restrictions` 计算属性（`Type: schema.TypeList`、`Computed: true`），`Elem` 为 `&schema.Resource{...}`，其中包含两个计算子字段 `reason`（TypeString，映射 `ComplianceRestriction.Reason`，枚举值 `ICP_RECORD_REQUIRED`/`GOVERNMENT_ORDER`）和 `region`（TypeString，映射 `ComplianceRestriction.Region`，ISO 3166 国家/地区码），Description 参考 vendored SDK `teov20220901.ComplianceRestriction` 字段注释编写
- [x] 1.2 确认无需修改 Create/Update/Delete 逻辑：`DomainComplianceRestrictions` 仅为 `DescribeFunctions` 出参，`CreateFunction`/`ModifyFunction`/`DeleteFunction` 请求中不存在对应参数，不将其加入 `immutableArgs`/`mutableArgs`

## 2. Read 函数

- [x] 2.1 在 `resourceTencentCloudTeoFunctionRead` 中（`respData` 非空分支内），新增对 `respData.DomainComplianceRestrictions` 的处理：先判断列表非 nil 且长度大于 0，再遍历每个 `ComplianceRestriction` 元素，对 `Reason`、`Region` 分别做非 nil 判断后写入 `map[string]interface{}`（键为 `reason`/`region`），追加到列表后调用 `d.Set("domain_compliance_restrictions", restrictionsList)`，模式与 `resource_tc_teo_origin_group.go` 中 `references` 块一致

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中使用 gomonkey 新增 Read 场景测试：mock teo client 的 `DescribeFunctions`（参考 `resource_tc_teo_function_replica_test.go` 的 mockMeta/ptrString 模式），使返回的 `Function.DomainComplianceRestrictions` 包含至少 2 个条目，调用 `res.Read(d, meta)` 后断言 state 中 `domain_compliance_restrictions` 的 `reason`/`region` 值与 mock 数据一致
- [x] 3.2 新增 `DomainComplianceRestrictions` 为 nil/空列表场景的测试：断言 Read 不报错、不设置该属性（state 中为空）
- [x] 3.3 新增子字段为 nil 的容错场景测试：mock 返回条目 `Reason` 或 `Region` 为 nil 时，Read 正常完成且对应键不写入（读回零值空字符串，与 `origin_group` references 测试的 nil 断言模式一致）
- [x] 3.4 保证既有测试（`TestAccTencentCloudTeoFunctionResource_basic`、`TestParseTeoFunctionOriginalName` 等）不受影响

## 4. 文档同步

- [x] 4.1 检查 `tencentcloud/services/teo/resource_tc_teo_function.md`：保持一句话描述、Example Usage、Import 三段结构不变，如示例需体现新计算属性可在 Example 中补充输出说明；不手动添加 Argument/Attribute Reference（由 `make doc` 自动生成）
- [x] 4.2 确认 `tencentcloud/provider.go`/`provider.md` 无需变更（`tencentcloud_teo_function` 已注册）

## 5. 代码正确性检查

- [x] 5.1 对照 vendored SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）复核：`Function.DomainComplianceRestrictions` 类型为 `[]*ComplianceRestriction`，`ComplianceRestriction` 仅含 `Reason`/`Region` 两个 `*string` 字段，schema 与 Read 映射与之完全一致
- [x] 5.2 检查所有函数返回的 error 均已处理（`d.Set` 结果用 `_ =` 赋值），确保生成的 Go 代码可正确编译（不执行 go build/go vet，由后续流程验证）

## 6. 收尾（由 tfpacer-finalize skill 执行，非本阶段操作）

- [x] 6.1 `gofmt` 格式化变更的 Go 文件（由 tfpacer-finalize skill 执行）
- [x] 6.2 `make doc` 生成 website/docs 文档与 `.changelog` 文件（由 tfpacer-finalize skill 执行）
