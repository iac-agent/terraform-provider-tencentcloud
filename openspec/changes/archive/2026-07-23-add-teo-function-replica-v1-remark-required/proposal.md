## Why

资源 `tencentcloud_teo_function_replica_v1` 目前尚未在 Terraform Provider 中实现。TEO（EdgeOne）边缘函数副本（Function Replica）是边缘函数的版本管理能力，用户需要通过 Terraform 管理函数副本的创建、读取、更新和删除。根据需求，`remark`（副本描述）字段应设为必填（Required），确保每个副本创建时都带有描述信息。

## What Changes

- 新增资源 `tencentcloud_teo_function_replica_v1`，实现完整的 CRUD 逻辑：
  - Create: 调用 `CreateFunctionReplica` API 创建函数副本
  - Read: 调用 `DescribeFunctionReplicas` API 读取函数副本详情
  - Update: 调用 `ModifyFunctionReplica` API 更新函数副本
  - Delete: 调用 `DeleteFunctionReplica` API 删除函数副本
- `remark` 字段设为 **Required**（必填），区别于 `tencentcloud_teo_function` 资源中 `remark` 为 Optional 的设计
- 使用 `zone_id` + `function_id` + `replica_name` 作为复合 ID
- 在 `provider.go` 和 `provider.md` 中注册该资源
- 生成对应的 `.md` 文档和单元测试

## Capabilities

### New Capabilities
- `teo-function-replica-v1-resource`: 新增 TEO 边缘函数副本资源 tencentcloud_teo_function_replica_v1，支持 CRUD 操作，remark 字段为必填

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_function_replica_v1.go`、对应的 `_test.go`、`.md` 文件
- 修改文件：`tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（文档更新）
- 依赖云 API：`CreateFunctionReplica`、`DescribeFunctionReplicas`、`ModifyFunctionReplica`、`DeleteFunctionReplica`（TEO v20220901）
- 向后兼容：纯新增资源，不影响现有资源配置和 state
