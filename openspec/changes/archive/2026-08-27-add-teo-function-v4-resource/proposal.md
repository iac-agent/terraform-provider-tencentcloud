## Why

腾讯云边缘函数（TEO Edge Function）是该产品下用于在边缘节点执行 JavaScript 代码的核心能力。现有 provider 已存在 `tencentcloud_teo_function` 资源，但本次需要新增一个独立的 `tencentcloud_teo_function_v4` 资源，以便在不影响既有 `tencentcloud_teo_function` 行为与 state 的前提下，按统一的 `RESOURCE_KIND_GENERAL` 代码风格（参考 `tencentcloud_igtm_strategy`）重新提供边缘函数的完整生命周期管理（创建 / 查询 / 修改 / 删除 / 导入）。腾讯云 TEO SDK 的 `teo/v20220901` 包已内置 `CreateFunction`、`DescribeFunctions`、`ModifyFunction`、`DeleteFunction` 四个同步接口，且已存在于当前 vendor 中，可行性已确认。

## What Changes

- 新增资源 `tencentcloud_teo_function_v4`，对应文件 `tencentcloud/services/teo/resource_tc_teo_function_v4.go`，实现 `RESOURCE_KIND_GENERAL` 的标准 CRUD：
  - Create：调用 `CreateFunction`（入参 `ZoneId`、`Name`、`Content`、`Remark`），返回 `FunctionId` 作为资源标识的一部分。
  - Read：调用 `DescribeFunctions`（按 `ZoneId` + `FunctionIds` 精确查询），回填 `function_id`、`zone_id`、`name`、`remark`、`content`、`domain`、`create_time`、`update_time`。
  - Update：调用 `ModifyFunction`（入参 `ZoneId`、`FunctionId`、`Remark`、`Content`）。
  - Delete：调用 `DeleteFunction`（入参 `ZoneId`、`FunctionId`）。
- 资源采用复合 ID `zone_id#function_id`（使用 `tccommon.FILED_SP` 分隔），与现有 `tencentcloud_teo_function` 保持一致，导入时需使用该联合 ID。
- `name` 字段为 `ForceNew`（`ModifyFunction` 不支持修改 `name`），更新阶段将其加入 `immutableArgs` 校验；`remark`、`content` 为可更新字段。
- Read 阶段对 `DescribeFunctions` 返回的拼接 `name`（形如 `my-func-zone-2qtuhspy7cr6-1310708577`）进行拆分处理，仅将原始 `name` 设置到 state，复用既有 `ParseTeoFunctionOriginalName` 辅助函数的逻辑（避免 plan 阶段误报变更）。
- 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_function_v4`，并在 `tencentcloud/provider.md` 中追加注册说明。
- 新增服务层查询函数 `TeoService.DescribeTeoFunctionV4ById`（按 `FunctionIds` 精确查询，复用既有 `DescribeFunctions` 调用，并使用 `tccommon.ReadRetryTimeout` 做 retry）。
- 新增单元测试文件 `resource_tc_teo_function_v4_test.go`，使用 gomonkey 对云 API 进行 mock，仅做业务逻辑的单测（不使用 terraform 测试套件）。
- 新增文档 `resource_tc_teo_function_v4.md`（参考 `resource_tc_teo_function.md` / `resource_tc_igtm_strategy.md` 风格）。

## Capabilities

### New Capabilities
- `teo-function-v4-resource`: 通过 `tencentcloud_teo_function_v4` 资源管理 TEO 边缘函数（Edge Function）的完整生命周期（创建 / 查询 / 修改 / 删除 / 导入），schema 与 `CreateFunction`/`DescribeFunctions`/`ModifyFunction`/`DeleteFunction` 四个云 API 入参严格对齐，包含 `ParseTeoFunctionOriginalName` 的 name 拆分行为。

### Modified Capabilities
<!-- 无：本次仅新增资源，不修改任何已有 capability 的需求级行为。 -->

## Impact

- **新增代码**:
  - `tencentcloud/services/teo/resource_tc_teo_function_v4.go`（schema + CRUD + name 拆分逻辑，单文件，代码风格参考 `tencentcloud_igtm_strategy`）。
  - `tencentcloud/services/teo/resource_tc_teo_function_v4.md`（资源文档 + import 示例）。
  - `tencentcloud/services/teo/resource_tc_teo_function_v4_test.go`（gomonkey mock 单元测试）。
- **修改代码**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`: 新增 `DescribeTeoFunctionV4ById` 查询辅助函数。
  - `tencentcloud/provider.go`: 在 `ResourcesMap` 中注册 `tencentcloud_teo_function_v4`。
  - `tencentcloud/provider.md`: 追加该资源的注册说明。
- **使用的云 API**: `CreateFunction`、`DescribeFunctions`、`ModifyFunction`、`DeleteFunction`（均已存在于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`，均为同步接口，无需轮询异步任务）。
- **无破坏性变更**: 纯新增资源，不影响既有 `tencentcloud_teo_function` 资源及其 state。
- **无需 SDK 升级**: 所需接口均已 vendor。
