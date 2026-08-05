## Why

腾讯云 EdgeOne (TEO) 推理服务的 API Token 管理目前只能通过控制台或 API 手动操作，无法通过 Terraform 进行基础设施即代码管理。用户需要创建、查询和删除推理 API Token 以使用 TEO 推理服务，而 Terraform Provider 缺少对应的资源支持。腾讯云已提供 `CreateInferenceAPIToken`、`DescribeInferenceAPITokens` 和 `DeleteInferenceAPIToken` 三个 API，SDK 已包含在 vendor 中，现在是添加此资源的最佳时机。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_inference_api_token_v7` 用于管理 TEO 推理 API Token
- 支持创建推理 API Token（CreateInferenceAPIToken API）
- 支持查询推理 API Token（DescribeInferenceAPITokens API）
- 支持删除推理 API Token（DeleteInferenceAPIToken API）
- 在 `tencentcloud/provider.go` 中注册新资源
- 添加完整的文档和单元测试

## Capabilities

### New Capabilities
- `teo-inference-api-token`: 管理 TEO 推理 API Token 的 Terraform 资源，支持创建、读取和删除操作。该资源为 CRD 类型（无 Update 接口），所有用户可配置字段均为 ForceNew。

### Modified Capabilities
<!-- 无现有 capability 需要修改 -->

## Impact

**新增文件:**
- `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7.go` - 资源实现
- `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7_test.go` - 资源单元测试
- `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7.md` - 资源文档

**修改文件:**
- `tencentcloud/provider.go` - 注册新资源
- `tencentcloud/provider.md` - 注册新资源

**依赖:**
- 依赖 `tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的 `CreateInferenceAPIToken`、`DescribeInferenceAPITokens`、`DeleteInferenceAPIToken` API（已在 vendor 中）