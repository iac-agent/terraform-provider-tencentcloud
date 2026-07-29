## 1. 资源 Schema 与 CRUD 实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v3.go`，定义资源 Schema（zone_id/ForceNew, name, content, remark, function_id/Computed, domain/Computed, create_time/Computed, update_time/Computed）及 Import 支持
- [x] 1.2 实现 Create 方法：调用 CreateFunction API，检查返回的 FunctionId 非空，轮询 DescribeFunctions 直到 Domain 非空，设置联合 ID（zone_id#function_id）
- [x] 1.3 实现 Read 方法：解析联合 ID，调用 DescribeFunctions API 查询函数详情，nil 检查后设置各字段，资源不存在时先打印 log 再 SetId("")
- [x] 1.4 实现 Update 方法：name 加入 immutableArgs 检查不可变，content 和 remark 加入 mutableArgs，变更时调用 ModifyFunction API
- [x] 1.5 实现 Delete 方法：解析联合 ID，调用 DeleteFunction API
- [x] 1.6 创建 `tencentcloud/services/teo/resource_tc_teo_function_v3_extension.go`（空 package teo 文件）

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DescribeTeoFunctionV3ById` 方法，通过 ZoneId + FunctionId 调用 DescribeFunctions API 查询函数详情

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_function_v3` 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_function_v3` 资源条目

## 4. 文档与测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v3.md` 资源示例文档（包含一句话描述、Example Usage、Import 说明）
- [x] 4.2 创建 `tencentcloud/services/teo/resource_tc_teo_function_v3_test.go` 单元测试文件，使用 gomonkey mock 云 API 调用
