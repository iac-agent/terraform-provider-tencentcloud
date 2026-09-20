## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoFunctionV6ById(ctx, zoneId, functionId) (*teo.Function, error)` 方法，构建 `DescribeFunctionsRequest`（设置 ZoneId 与 FunctionIds），使用 `tccommon.ReadRetryTimeout` retry 调用 `DescribeFunctions`，若 `response.Response == nil || len(response.Response.Functions) < 1` 返回 nil，否则返回 `Functions[0]`

## 2. 资源代码实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v6.go`，定义 `ResourceTencentCloudTeoFunctionV6()` 函数，包含 schema（zone_id/name/content/remark 输入字段，function_id/domain/domain_compliance_restrictions/create_time/update_time 计算字段）与 Importer，参照 `tencentcloud_igtm_strategy` 代码风格
- [x] 2.2 实现 `resourceTencentCloudTeoFunctionV6Create`：构建 CreateFunctionRequest（ZoneId/Name/Content/Remark）→ WriteRetryTimeout retry 调用 → 检查 FunctionId 为空返回 NonRetryableError → `d.SetId(zoneId#functionId)` → 异步轮询 DescribeFunctions 直到 Domain 非空（Delay 10s, MinTimeout 3s, Timeout 600s）→ 调用 Read
- [x] 2.3 实现 `resourceTencentCloudTeoFunctionV6Read`：解析复合 ID → 调用 DescribeTeoFunctionV6ById → 若空先 `log.Printf("[CRUD] teo_function_v6 id=%s", d.Id())` 再 `d.SetId("")` → 非空逐字段判断 nil 后 d.Set（含 domain_compliance_restrictions 嵌套列表）
- [x] 2.4 实现 `resourceTencentCloudTeoFunctionV6Update`：解析复合 ID → 检查 immutableArgs(name) → 若 remark/content 有变更构建 ModifyFunctionRequest（ZoneId/FunctionId/Remark/Content）→ WriteRetryTimeout retry 调用 → 调用 Read
- [x] 2.5 实现 `resourceTencentCloudTeoFunctionV6Delete`：解析复合 ID → 构建 DeleteFunctionRequest（ZoneId/FunctionId）→ WriteRetryTimeout retry 调用
- [x] 2.6 实现异步轮询 state refresh 函数 `resourceTeoFunctionV6CreateStateRefreshFunc`，调用 DescribeFunctions 检查 Functions[0].Domain 是否非空

## 3. 单元测试实现

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v6_test.go`，使用 gomonkey mock 云 API（CreateFunctionWithContext/DescribeFunctions/ModifyFunctionWithContext/DeleteFunctionWithContext），编写 Create/Read/Update/Delete 业务逻辑单元测试用例，禁止使用 terraform 测试套件

## 4. 文档与注册

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v6.md`，包含一句话描述（提及 EdgeOne/TEO）、Example Usage（HCL）、Import 部分（说明使用复合 id zoneId#functionId）
- [x] 4.2 在 `tencentcloud/provider.go` 的 resources map 中注册 `"tencentcloud_teo_function_v6": teo.ResourceTencentCloudTeoFunctionV6()`
- [x] 4.3 在 `tencentcloud/provider.md` 资源列表中新增 `tencentcloud_teo_function_v6`

## 5. 验证

- [x] 5.1 确认生成的 Go 代码可正确构建（检查 import、函数签名、字段映射与云 API struct 一致），不执行 go build
- [x] 5.2 确认 CRUD 参数映射正确：CreateFunction 入参（ZoneId/Name/Content/Remark）、ModifyFunction 入参（ZoneId/FunctionId/Remark/Content）、DeleteFunction 入参（ZoneId/FunctionId）均与云 API struct 字段一致