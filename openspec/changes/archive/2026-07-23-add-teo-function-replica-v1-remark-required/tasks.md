## 1. 资源 Schema 与 CRUD 实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_function_replica_v1.go`，定义 Schema（zone_id、function_id、replica_name 为 Required+ForceNew；content 为 Required；remark 为 Required；create_time、update_time 为 Computed）
- [x] 1.2 实现 Create 函数：调用 CreateFunctionReplica API，设置复合 ID `zone_id#function_id#replica_name`，检查返回值非空
- [x] 1.3 实现 Read 函数：调用 DescribeFunctionReplicas API，通过 Filters 的 replica-name 过滤定位副本，设置各字段，处理空响应
- [x] 1.4 实现 Update 函数：检测 content 和 remark 变更，调用 ModifyFunctionReplica API，检查不可变字段变更
- [x] 1.5 实现 Delete 函数：调用 DeleteFunctionReplica API，将 replica_name 作为单元素列表传入 ReplicaNames

## 2. 服务层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DescribeTeoFunctionReplicaV1ByFilter` 方法，调用 DescribeFunctionReplicas API 并通过 replica-name 过滤定位具体副本

## 3. Provider 注册与文档

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_function_replica_v1` 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_function_replica_v1` 资源条目
- [x] 3.3 创建资源文档 `tencentcloud/services/teo/resource_tc_teo_function_replica_v1.md`，包含一句话描述、Example Usage、Import 部分

## 4. 单元测试

- [x] 4.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_function_replica_v1_test.go`，使用 gomonkey mock 云 API，覆盖 Create、Read、Update、Delete 场景
- [x] 4.2 运行 `go test -gcflags=all=-l` 验证所有单元测试通过
