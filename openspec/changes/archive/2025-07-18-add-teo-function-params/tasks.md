## 1. Schema 定义与 CRUD 函数修改

- [x] 1.1 在 `resource_tc_teo_function.go` 的 Schema 中新增 `function_ids` 参数（Optional, TypeList of TypeString）
- [x] 1.2 在 `resource_tc_teo_function.go` 的 Schema 中新增 `filters` 参数（Optional, TypeList of Object，包含 name 和 values 字段）
- [x] 1.3 在 `resource_tc_teo_function.go` 的 Schema 中新增 `functions` 参数（Computed, TypeList of Object，包含 function_id、zone_id、name、remark、content、domain、create_time、update_time 字段）
- [x] 1.4 修改 `resourceTencentCloudTeoFunctionRead` 函数，在调用 DescribeFunctions 时传入 `function_ids` 和 `filters` 参数，并处理 `functions` 出参映射

## 2. Service 层扩展

- [x] 2.1 在 `service_tencentcloud_teo.go` 中新增或扩展 DescribeTeoFunction 相关方法，支持传入 function_ids（列表）和 filters 参数

## 3. 扩展文件修改

- [x] 3.1 修改 `resource_tc_teo_function_extension.go`，添加 filters 参数的扁平化（flatten）辅助函数和 functions 列表的映射逻辑

## 4. 单元测试

- [x] 4.1 在 `resource_tc_teo_function_test.go` 中补充 `function_ids`、`filters`、`functions` 参数的单元测试用例，使用 gomonkey mock 云 API

## 5. 文档更新

- [x] 5.1 修改 `resource_tc_teo_function.md` 文件，添加 `function_ids`、`filters`、`functions` 参数的说明和示例
