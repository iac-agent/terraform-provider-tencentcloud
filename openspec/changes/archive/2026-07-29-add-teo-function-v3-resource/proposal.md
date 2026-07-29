## Why

TencentCloud EdgeOne (TEO) 边缘函数需要新增一个 v3 版本的 Terraform 资源 `tencentcloud_teo_function_v3`，以支持对边缘函数的完整生命周期管理（创建、读取、更新、删除）。虽然已存在 `tencentcloud_teo_function` 资源，但 v3 版本将提供更规范的实现和更好的代码质量。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_function_v3`，管理 TEO 边缘函数的完整 CRUD 生命周期
- 使用 `CreateFunction` API 创建边缘函数，入参包括 ZoneId、Name、Content、Remark，出参返回 FunctionId
- 使用 `DescribeFunctions` API 读取边缘函数详情，通过 ZoneId + FunctionId 查询
- 使用 `ModifyFunction` API 更新边缘函数，支持修改 Content 和 Remark
- 使用 `DeleteFunction` API 删除边缘函数
- 资源 ID 使用 zone_id 和 function_id 的联合 ID（以 `tccommon.FILED_SP` 分隔）
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 生成对应的 `.md` 文档和单元测试文件

## Capabilities

### New Capabilities
- `teo-function-v3-resource`: 管理 TEO 边缘函数的完整 CRUD 生命周期资源，支持创建、读取、更新、删除边缘函数，使用联合 ID（zone_id + function_id）

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_function_v3.go`、`tencentcloud/services/teo/resource_tc_teo_function_v3_extension.go`、`tencentcloud/services/teo/resource_tc_teo_function_v3_test.go`、`tencentcloud/services/teo/resource_tc_teo_function_v3.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（添加资源文档）
- 依赖云 API：`teo/v20220901` 包中的 CreateFunction、DescribeFunctions、ModifyFunction、DeleteFunction
