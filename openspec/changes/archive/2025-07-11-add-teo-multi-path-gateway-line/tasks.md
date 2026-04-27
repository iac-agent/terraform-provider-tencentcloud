## 1. 资源 Schema 与 CRUD 函数实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.go` 文件，定义 `ResourceTencentCloudTeoMultiPathGatewayLine()` 资源函数，包含 Schema 定义（zone_id Required ForceNew, gateway_id Required ForceNew, line_type Required, line_address Required, proxy_id Optional, rule_id Optional, line_id Computed）和 Import 支持
- [x] 1.2 实现 `resourceTencentCloudTeoMultiPathGatewayLineCreate` 函数：调用 CreateMultiPathGatewayLine 接口，传入 zone_id/gateway_id/line_type/line_address/proxy_id/rule_id，获取返回的 line_id，设置复合 ID（zone_id + FIELD_SP + gateway_id + FIELD_SP + line_id），调用 Read 刷新状态
- [x] 1.3 实现 `resourceTencentCloudTeoMultiPathGatewayLineRead` 函数：从 d.Get() 获取 zone_id/gateway_id/line_id，调用 DescribeMultiPathGatewayLine 接口，解析返回的 Line 对象（LineId/LineType/LineAddress/ProxyId/RuleId），映射到 Terraform schema 字段；资源不存在时清除 state；使用 helper.Retry() + tccommon.ReadRetryTimeout 进行重试
- [x] 1.4 实现 `resourceTencentCloudTeoMultiPathGatewayLineUpdate` 函数：从 d.Get() 获取 zone_id/gateway_id/line_id，调用 ModifyMultiPathGatewayLine 接口传入 line_type/line_address/proxy_id/rule_id，调用 Read 刷新状态
- [x] 1.5 实现 `resourceTencentCloudTeoMultiPathGatewayLineDelete` 函数：从 d.Get() 获取 zone_id/gateway_id/line_id，调用 DeleteMultiPathGatewayLine 接口；使用 helper.Retry() + tccommon.ReadRetryTimeout 进行重试

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_multi_path_gateway_line` 资源，参考 tencentcloud_igtm_strategy 的注册方式
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_multi_path_gateway_line` 资源条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.md` 资源样例文档，包含一句话描述、Example Usage（custom 类型线路和 proxy 类型线路示例）、Import 说明

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line_test.go`，使用 gomonkey mock 云 API 调用，编写 Create/Read/Update/Delete 操作的单元测试用例
- [x] 4.2 运行 `go test`（带 `-gcflags=all=-l` 参数）验证所有单元测试通过
