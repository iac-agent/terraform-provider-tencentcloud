## Why

腾讯云 EdgeOne (TEO) 边缘函数副本用于版本管理与灰度发布（通过请求头 `EO-Function-Replica-Name` 访问特定副本）。现有 `tencentcloud_teo_function_replica` 资源为早期实现，未暴露 `created_on`/`modified_on` 等回读字段，且缺少对查询接口排序/过滤参数的完整覆盖。云 API（`teov20220901`）已提供完整的 `CreateFunctionReplica`/`DescribeFunctionReplicas`/`ModifyFunctionReplica`/`DeleteFunctionReplica` 接口，本次新增 `tencentcloud_teo_function_replica_v4` 资源（遵循 `tencentcloud_teo_l7_acc_rule_v2` 的版本化命名惯例，不破坏现有 v1 资源），为用户提供覆盖完整参数映射的边缘函数副本管理能力。

## What Changes

- 新增 Terraform RESOURCE_KIND_GENERAL 资源 `tencentcloud_teo_function_replica_v4`，实现边缘函数副本的完整 CRUD 生命周期：
  - Create → `CreateFunctionReplica`（ZoneId / FunctionId / ReplicaName / Content / Remark）
  - Read → `DescribeFunctionReplicas`（按 replica-name 过滤，Limit=200，回读 content / remark / created_on / modified_on）
  - Update → `ModifyFunctionReplica`（ZoneId / FunctionId / ReplicaName / Content / Remark）
  - Delete → `DeleteFunctionReplica`（ZoneId / FunctionId / ReplicaNames）
- 资源通过 `zone_id#function_id#replica_name` 三段联合 ID 唯一标识（分隔符 `tccommon.FILED_SP`），支持 `terraform import`
- `zone_id` / `function_id` / `replica_name` 设为 ForceNew（Modify 接口不支持改名，`replica_name` 是定位副本的标识）
- 新增 `created_on` / `modified_on` Computed 字段（与 v1 资源的关键差异）
- 保留 `tencentcloud_teo_function_replica`（v1）资源不动，保证向后兼容
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源

## Capabilities

### New Capabilities
- `teo-function-replica-v4-resource`: 管理 TEO 边缘函数副本（v4）资源的完整 CRUD 生命周期，包括创建、读取（含 created_on/modified_on 回读）、更新和删除边缘函数副本

### Modified Capabilities
<!-- 无：本次为纯新增资源，不修改任何现有 capability -->

## Impact

- 新增文件:
  - `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.go`
  - `tencentcloud/services/teo/resource_tc_teo_function_replica_v4_test.go`
  - `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.md`
- 修改文件:
  - `tencentcloud/provider.go`（注册 `tencentcloud_teo_function_replica_v4`）
  - `tencentcloud/provider.md`（TEO Resources 列表新增条目，供 `make doc` 生成 website 文档）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中，无需变更）
- 兼容性: 纯新增资源，不影响现有 `tencentcloud_teo_function_replica`（v1）资源的行为与 state
- 接口均为同步接口（无 TaskId/FlowId 返回），无需异步轮询
