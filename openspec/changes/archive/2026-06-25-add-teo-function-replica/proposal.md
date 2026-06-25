## Why

EdgeOne (TEO) 边缘函数副本（Function Replica）是边缘函数的多版本管理能力，允许同一函数拥有多个副本（不同版本），每个副本可独立进行代码内容、描述的修改，并在函数规则中按权重或地域分发到不同副本。当前 Terraform Provider 缺少对函数副本资源的支持，用户无法通过 IaC 管理函数的多版本部署。

## What Changes

- 新增 `tencentcloud_teo_function_replica` 资源（RESOURCE_KIND_GENERAL），支持边缘函数副本的创建、读取、更新、删除和导入
- 通过 `CreateFunctionReplica` 接口创建函数副本
- 通过 `DescribeFunctionReplicas` 接口查询函数副本列表
- 通过 `ModifyFunctionReplica` 接口更新函数副本内容和描述
- 通过 `DeleteFunctionReplica` 接口删除函数副本
- 资源 ID 使用联合 ID 格式 `zone_id#function_id#replica_name`（因 Create 接口不返回副本 ID）

## Capabilities

### New Capabilities
- `teo-function-replica-resource`: 为 TEO 边缘函数提供副本管理资源，支持副本的 CRUD 和导入

### Modified Capabilities
<!-- None -->

## Impact

- **新增文件**: `tencentcloud/services/teo/resource_tc_teo_function_replica.go`（资源实现）
- **新增文件**: `tencentcloud/services/teo/resource_tc_teo_function_replica_test.go`（单元测试）
- **新增文件**: `tencentcloud/services/teo/resource_tc_teo_function_replica.md`（文档样例）
- **修改文件**: `tencentcloud/provider.go`（资源注册）
- **修改文件**: `tencentcloud/provider.md`（文档注册）
- **依赖**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（vendor 中已有 `CreateFunctionReplica`、`DescribeFunctionReplicas`、`ModifyFunctionReplica`、`DeleteFunctionReplica` 接口）