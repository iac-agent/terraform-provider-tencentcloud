## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_function.go` 的 Schema 中添加 `filters` 字段定义，类型为 `schema.TypeList`，Optional，嵌套 `name`(TypeString, Required) 和 `values`(TypeList of TypeString, Required)
- [x] 1.2 在 Update 函数的 `immutableArgs` 数组中添加 `"filters"`，确保 filters 变更时返回错误

## 2. Service 层更新

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 的 `DescribeTeoFunctionById` 方法中添加 `filters` 参数支持，将传入的 filters 转换为 `[]*teo.Filter` 并设置到 `DescribeFunctionsRequest` 中

## 3. CRUD 函数更新

- [x] 3.1 在 `resourceTencentCloudTeoFunctionRead` 函数中，从 `d.Get("filters")` 获取 filters 配置，传递给 `DescribeTeoFunctionById` 服务方法
- [x] 3.2 在 `resourceTeoFunctionCreateStateRefreshFunc_0_0` 函数中，不使用 filters 参数（保持现有逻辑不变，因为创建状态轮询使用 FunctionIds 精确定位）

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中添加使用 gomock/gomonkey 的单元测试用例，验证 filters 参数在 READ 操作中正确传递给 DescribeFunctions API
- [x] 4.2 添加单元测试验证 filters 在 immutableArgs 中，更新 filters 时应返回错误
- [x] 4.3 运行 `go test` 确保所有单元测试通过

## 5. 文档更新

- [x] 5.1 更新 `tencentcloud/services/teo/resource_tc_teo_function.md`，在 Example Usage 中添加带 filters 参数的示例，并在参数说明中添加 filters 字段描述
