## Why

`tencentcloud_teo_function` 资源当前仅支持通过 `zone_id` + `function_id`（复合 ID）读取单个函数信息。但 `DescribeFunctions` API 支持更多查询参数（`FunctionIds`、`Filters`），用户无法通过 Terraform 利用这些参数进行批量查询或条件过滤。需要将 `DescribeFunctions` API 支持的 `function_ids`、`filters` 入参和 `functions` 出参暴露到 Terraform 资源中，以增强资源的查询能力。

## What Changes

- 为 `tencentcloud_teo_function` 资源新增 `function_ids` 参数（Optional, TypeList of String），对应 `DescribeFunctions` API 的 `FunctionIds` 入参，用于按函数 ID 列表过滤查询
- 为 `tencentcloud_teo_function` 资源新增 `filters` 参数（Optional, TypeList of Object），对应 `DescribeFunctions` API 的 `Filters` 入参，用于按条件（name、remark）过滤查询
- 为 `tencentcloud_teo_function` 资源新增 `functions` 参数（Computed, TypeList of Object），对应 `DescribeFunctions` API 的 `Functions` 出参，用于返回查询到的函数列表详情

## Capabilities

### New Capabilities
- `teo-function-query-params`: 为 tencentcloud_teo_function 资源新增 DescribeFunctions API 支持的查询参数（function_ids、filters）和返回结果（functions）

### Modified Capabilities

## Impact

- **代码变更**: 修改 `tencentcloud/services/teo/resource_tc_teo_function.go` 中的 Schema 定义和 Read 函数
- **扩展变更**: 修改 `tencentcloud/services/teo/resource_tc_teo_function_extension.go` 添加扩展逻辑
- **文档变更**: 修改 `tencentcloud/services/teo/resource_tc_teo_function.md` 添加新参数文档
- **API 依赖**: 使用已有的 `DescribeFunctions` API（`teo/v20220901` 包），无需新增 SDK 依赖
