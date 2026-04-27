## 1. Schema 定义与资源工厂函数

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.go`，定义 `ResourceTencentCloudTeoMultiPathGatewayLine()` 工厂函数和完整 schema，包含 zone_id(Required,ForceNew)、gateway_id(Required,ForceNew)、line_type(Required)、line_address(Required)、proxy_id(Optional)、rule_id(Optional)、line_id(Computed) 字段
- [x] 1.2 添加 Importer 支持，使用 `schema.ImportStatePassthrough`

## 2. CRUD 函数实现

- [x] 2.1 实现 `resourceTencentCloudTeoMultiPathGatewayLineCreate` 函数：调用 `CreateMultiPathGatewayLine` API，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，成功后设置复合 ID `{zone_id}#{gateway_id}#{line_id}`，并调用 Read 同步状态
- [x] 2.2 实现 `resourceTencentCloudTeoMultiPathGatewayLineRead` 函数：从复合 ID 解析 zone_id/gateway_id/line_id，调用 `DescribeMultiPathGatewayLine` API 查询详情，将响应字段写回 schema，资源不存在时设置 ID 为空
- [x] 2.3 实现 `resourceTencentCloudTeoMultiPathGatewayLineUpdate` 函数：检测 line_type/line_address/proxy_id/rule_id 变更，调用 `ModifyMultiPathGatewayLine` API，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹，完成后调用 Read 同步状态
- [x] 2.4 实现 `resourceTencentCloudTeoMultiPathGatewayLineDelete` 函数：从复合 ID 解析参数，调用 `DeleteMultiPathGatewayLine` API，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中添加 `"tencentcloud_teo_multi_path_gateway_line": teo.ResourceTencentCloudTeoMultiPathGatewayLine()` 注册条目
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO 资源部分添加 `tencentcloud_teo_multi_path_gateway_line` 条目

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line_test.go`，使用 gomonkey mock 方式编写 Create/Read/Update/Delete 操作的单元测试
- [x] 4.2 运行 `go test` 确保单元测试通过

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.md` 示例文件，包含一句话描述、Example Usage HCL 示例、Import 说明（使用 `zone_id#gateway_id#line_id` 格式）
