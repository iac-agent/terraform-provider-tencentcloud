## Why

EdgeOne (TEO) 提供了边缘函数（Edge Function）能力，允许用户在边缘节点运行 JavaScript 代码处理请求。目前 provider 中已有 `tencentcloud_teo_function` 资源（由 iacg 代码生成器生成），但其设计存在若干局限：仅暴露部分字段、缺少 `domain_compliance_restrictions` 等查询返回字段的 schema、且未严格遵循项目代码规范（如新建资源应使用 gomonkey mock 编写单元测试而非 terraform 测试套件）。本次新增 `tencentcloud_teo_function_v6` 资源，采用手写实现、严格遵循项目规范，完整覆盖云 API 的 CRUD 能力与查询返回字段，为用户提供一个规范化、可维护的边缘函数管理资源。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_function_v6`（RESOURCE_KIND_GENERAL），管理 EdgeOne 边缘函数的完整生命周期（创建、读取、更新、删除）。
- 新增资源代码文件 `tencentcloud/services/teo/resource_tc_teo_function_v6.go`，包含 schema 定义与 CRUD 实现。
- 新增服务层方法 `DescribeTeoFunctionV6ById`（位于 `tencentcloud/services/teo/service_tencentcloud_teo.go`），封装 `DescribeFunctions` 调用以查询单个函数。
- 新增单元测试文件 `tencentcloud/services/teo/resource_tc_teo_function_v6_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。
- 新增资源文档 `tencentcloud/services/teo/resource_tc_teo_function_v6.md`。
- 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_function_v6` 资源。
- 在 `tencentcloud/provider.md` 资源列表中新增 `tencentcloud_teo_function_v6`。

## Capabilities

### New Capabilities
- `teo-function-v6-resource`: 管理 EdgeOne 边缘函数（Function）的 Terraform 资源，覆盖创建（CreateFunction）、读取（DescribeFunctions）、更新（ModifyFunction）、删除（DeleteFunction）全生命周期，支持异步创建轮询、复合 ID 导入等能力。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **新增代码文件**:
  - `tencentcloud/services/teo/resource_tc_teo_function_v6.go`
  - `tencentcloud/services/teo/resource_tc_teo_function_v6_test.go`
  - `tencentcloud/services/teo/resource_tc_teo_function_v6.md`
- **修改代码文件**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`（新增 `DescribeTeoFunctionV6ById` 方法）
  - `tencentcloud/provider.go`（注册新资源）
  - `tencentcloud/provider.md`（新增资源名称）
- **云 API 依赖**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，使用 `CreateFunction`、`DescribeFunctions`、`ModifyFunction`、`DeleteFunction` 四个接口，均为同步接口（CreateFunction 返回 FunctionId 后需轮询 DescribeFunctions 等待函数部署生效）。
- **向后兼容**: 新增资源不影响现有 `tencentcloud_teo_function` 资源，完全向后兼容。