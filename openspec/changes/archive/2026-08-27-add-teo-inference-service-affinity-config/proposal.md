## Why

TEO 推理服务（Inference Service）的 Terraform 资源 `tencentcloud_teo_inference_service_v1` 需要支持亲和性配置（AffinityConfig），使用户能够通过 Terraform 声明式配置推理服务的会话亲和策略。当前该资源尚未创建，本次变更将创建完整的 CRUD 资源，并在 Create 和 Update 操作中支持 `AffinityConfig` 参数块，包含亲和开关、亲和方式及会话 ID 亲和配置。

## What Changes

- 创建 TEO 推理服务 Terraform 资源 `tencentcloud_teo_inference_service_v1`（RESOURCE_KIND_GENERAL），提供完整的 CRUD 生命周期管理
- 在资源 schema 中新增 `affinity_config` 嵌套参数块，包含以下子字段：
  - `switch`（TypeString，Optional）：推理服务亲和总开关，枚举值 `On`/`Off`
  - `affinity_mode`（TypeString，Optional）：推理服务亲和方式，枚举值 `SessionId`
  - `source`（TypeString，Optional）：会话 ID 参数的传递位置，枚举值 `Header`
  - `header_name`（TypeString，Optional）：传递会话 ID 的请求头名称
- Create 操作：调用 `CreateInferenceService` API 时传入 `AffinityConfig` 参数
- Update 操作：调用 `ModifyInferenceService` API 时传入 `AffinityConfig` 参数
- Delete 操作：调用 `OperateInferenceService` API 停止/删除推理服务
- Read 操作：调用 `DescribeInferenceServices` API 获取推理服务信息

## Capabilities

### New Capabilities

- `teo-inference-service-resource`: 创建 TEO 推理服务的完整 Terraform 资源，支持亲和性配置（AffinityConfig）的 CRUD 管理

### Modified Capabilities

<!-- 无现有 capability 被修改，本次为全新资源 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_inference_service_v1.go`（新建：schema 定义与 CRUD 逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_inference_service_v1_test.go`（新建：单元测试）
  - `tencentcloud/services/teo/resource_tc_teo_inference_service_v1.md`（新建：资源文档）
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`（可能修改：新增 service 层辅助函数）
  - `tencentcloud/provider.go`（修改：注册新资源）
  - `tencentcloud/provider.md`（修改：注册新资源文档）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901` 包的 API（`CreateInferenceService`、`DescribeInferenceServices`、`ModifyInferenceService`、`OperateInferenceService`），无需变更 vendor
- 向后兼容：全新资源，不影响现有配置或 state
- 文档：需同步创建 `.md` 示例文件，并通过 `make doc` 生成 website docs