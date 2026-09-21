## Why

腾讯云 TEO（边缘安全加速平台）SDK 已提供边缘函数副本（Function Replica）的完整 CRUD API（`CreateFunctionReplica` / `DescribeFunctionReplicas` / `ModifyFunctionReplica` / `DeleteFunctionReplica`，包 `teov20220901`）。当前 provider 中的 `tencentcloud_teo_function_replica` 资源仅覆盖基础参数（`zone_id`、`function_id`、`replica_name`、`content`、`remark`），未暴露查询侧参数（`sort_by`、`sort_order`、`filters`）与删除侧参数（`replica_names`），也不回填 `created_on` / `modified_on` 计算属性。需要新增 `tencentcloud_teo_function_replica_v5` 资源（RESOURCE_KIND_GENERAL），完整对齐云 API 能力，使用户可以声明式管理边缘函数副本的整个生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_function_replica_v5`，代码文件 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.go`：
  - **Create**：调用 `CreateFunctionReplica`，传入 `ZoneId`、`FunctionId`、`ReplicaName`、`Content`（必填）与 `Remark`（可选）；成功后以 `zone_id#function_id#replica_name` 复合 ID 写入 state，再调用 Read 回填
  - **Read**：调用 `DescribeFunctionReplicas`，传入 `ZoneId`、`FunctionId`（来自复合 ID），透传 `SortBy`、`SortOrder`、`Filters`（用户配置的 filters 额外追加内部 `replica-name` 精确过滤条件用于定位本资源），`Limit` 取 API 注释标注的最大值 200；从 `FunctionReplicas` 列表中按 `replica_name` 精确匹配后平铺回填 `content`、`remark`、`created_on`、`modified_on`
  - **Update**：当 `content` / `remark` 变更时调用 `ModifyFunctionReplica`，传入 `ZoneId`、`FunctionId`、`ReplicaName`（来自复合 ID）与变更后的 `Content`、`Remark`
  - **Delete**：调用 `DeleteFunctionReplica`，传入 `ZoneId`、`FunctionId` 与 `ReplicaNames`；`ReplicaNames` 优先取 schema 参数 `replica_names`，未配置时回退为复合 ID 中的 `replica_name` 单元素列表
- Schema 按 Describe 列表响应**平铺**（遵循规范：不引入 `function_replicas` 嵌套包装层），新增 `created_on` / `modified_on` Computed 字段
- 复合 ID 使用 `tccommon.FILED_SP`（`#`）分隔，支持 `terraform import`（ImportStatePassthrough）
- 在 `tencentcloud/provider.go` 的 ResourcesMap 注册资源，`tencentcloud/provider.md` 资源列表追加资源名
- 新增资源文档 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.md`（一句话描述 + Example Usage + Import 示例，说明联合 ID 用法；不含 Argument/Attribute Reference，由工具自动生成）
- 新增单元测试 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5_test.go`，使用 gomonkey mock 云 API，仅做业务逻辑单测（不使用 terraform 测试套件）

非破坏性变更：现有 `tencentcloud_teo_function_replica` 资源及其 spec（`teo-function-replica`）完全不受影响，本次为纯增量新增。

## Capabilities

### New Capabilities
- `teo-function-replica-v5`: 新资源 `tencentcloud_teo_function_replica_v5` 的完整生命周期管理，包括 schema 定义（含查询/删除侧参数与 Computed 时间字段）、Create/Read/Update/Delete 四个操作、复合 ID 与 import、provider 注册与文档

### Modified Capabilities
<!-- 无：现有 capability（含 teo-function-replica）的需求均不发生变化，本次为纯新增资源 -->

## Impact

- 代码：
  - 新增 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.go`
  - 新增 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5_test.go`
  - 新增 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.md`
  - 修改 `tencentcloud/provider.go`（ResourcesMap 中注册 `tencentcloud_teo_function_replica_v5`）
  - 修改 `tencentcloud/provider.md`（资源清单追加 `tencentcloud_teo_function_replica_v5`）
- 依赖：使用已 vendored 的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，四个云 API 方法（`CreateFunctionReplicaWithContext` / `DescribeFunctionReplicasWithContext` / `ModifyFunctionReplicaWithContext` / `DeleteFunctionReplicaWithContext`）均已存在，无需变更 vendor
- 向后兼容：纯新增资源，不影响任何现有资源、TF 配置与 state
- 文档：`website/` 目录文档由收尾阶段 `make doc` 自动生成，不在本变更内手工修改
