## ADDED Requirements

### Requirement: teo function replica v2 resource schema
`tencentcloud_teo_function_replica_v2` 资源 SHALL 在顶层平铺定义以下参数：`zone_id`（Required，ForceNew，字符串）、`function_id`（Required，ForceNew，字符串）、`replica_name`（Required，ForceNew，字符串）、`content`（Required，字符串，可更新）、`remark`（Optional，字符串，可更新），以及 Computed 输出字段 `created_on`、`modified_on`。资源 SHALL NOT 定义 `sort_by`/`sort_order`/`filters` 等查询辅助参数，也 SHALL NOT 定义 `function_replicas` 列表包装层。

#### Scenario: schema 定义完整
- **WHEN** 检查 `ResourceTencentCloudTeoFunctionReplicaV2()` 的 Schema
- **THEN** 存在 zone_id/function_id/replica_name（Required+ForceNew）、content（Required）、remark（Optional）、created_on/modified_on（Computed），且无任何列表包装或查询辅助字段

### Requirement: create teo function replica v2
创建 `tencentcloud_teo_function_replica_v2` 时，SHALL 调用 `CreateFunctionReplica`（WriteRetryTimeout retry 包装），请求填充 ZoneId/FunctionId/ReplicaName/Content/Remark。调用成功后 SHALL 以 `zone_id#function_id#replica_name`（tccommon.FILED_SP 分隔）设置资源 ID，随后调用 Read 同步状态；若 API 返回 Response 为 nil SHALL 返回 NonRetryableError。

#### Scenario: 创建成功
- **WHEN** 用户 apply 一个含 zone_id/function_id/replica_name/content/remark 的资源
- **THEN** 调用 CreateFunctionReplica 成功后 d.Id() 为 `zone_id#function_id#replica_name`，且 Read 将 content/remark/created_on/modified_on 写入 state

#### Scenario: 创建接口返回空响应
- **WHEN** CreateFunctionReplica 返回的 Response 为 nil
- **THEN** Create 返回 NonRetryableError，不写入空 ID

### Requirement: read teo function replica v2
读取 `tencentcloud_teo_function_replica_v2` 时，SHALL 从复合 ID 解析 zone_id/function_id/replica_name，调用 `DescribeFunctionReplicas`（ReadRetryTimeout retry 包装，`Limit=200`，Filter `replica-name` 精确值为副本名），在返回列表中按 ReplicaName 匹配目标副本；若响应为空或未匹配到，SHALL 先打印含资源 ID 的日志保留现场，再 `d.SetId("")` 并正常返回。SHALL 仅在响应字段非 nil 时调用 `d.Set`。

#### Scenario: 读取成功
- **WHEN** DescribeFunctionReplicas 返回包含同名副本的列表
- **THEN** content/remark/created_on/modified_on/zone_id/function_id/replica_name 被写入 state，ID 保持不变

#### Scenario: 副本不存在
- **WHEN** DescribeFunctionReplicas 返回空列表或不包含同名副本
- **THEN** 日志中保留原 ID 信息，d.SetId("") 被执行，Read 正常返回（无 error）

#### Scenario: ID 非法
- **WHEN** d.Id() 按分隔符 split 后长度不为 3
- **THEN** 返回 "id is broken" 类错误

### Requirement: update teo function replica v2
更新 `tencentcloud_teo_function_replica_v2` 时，SHALL 仅允许 `content`/`remark` 变更：当二者任一 `d.HasChange` 时调用 `ModifyFunctionReplica`（WriteRetryTimeout retry 包装，从复合 ID 解析 ZoneId/FunctionId/ReplicaName）。若 `zone_id`/`function_id`/`replica_name` 任一发生变更，SHALL 返回 error 提示该字段不可变。更新完成后 SHALL 调用 Read 同步状态。

#### Scenario: 更新 content 与 remark
- **WHEN** 用户修改 content 和/或 remark 后 apply
- **THEN** 调用 ModifyFunctionReplica 传入 ZoneId/FunctionId/ReplicaName/Content/Remark，随后 Read 刷新 state

#### Scenario: 修改不可变字段报错
- **WHEN** 用户修改 zone_id/function_id/replica_name 中的任一字段后 apply
- **THEN** Update 返回 error，说明该字段不可变更（immutable），不调用 ModifyFunctionReplica

#### Scenario: 无可变更字段
- **WHEN** content/remark 均无变化
- **THEN** 不调用 ModifyFunctionReplica，直接 Read 刷新状态

### Requirement: delete teo function replica v2
删除 `tencentcloud_teo_function_replica_v2` 时，SHALL 从复合 ID 解析三元组，调用 `DeleteFunctionReplica`（WriteRetryTimeout retry 包装），`ReplicaNames` 传入单个副本名。

#### Scenario: 删除成功
- **WHEN** 用户 destroy 该资源
- **THEN** DeleteFunctionReplica 被调用且 ReplicaNames 仅含目标副本名，Delete 正常返回

### Requirement: import teo function replica v2
`tencentcloud_teo_function_replica_v2` SHALL 支持 `terraform import`（ImportStatePassthrough），使用 `zone_id#function_id#replica_name` 复合 ID；md 文档 Import 部分 SHALL 说明该复合 ID 用法。

#### Scenario: 导入已有副本
- **WHEN** 用户执行 `terraform import tencentcloud_teo_function_replica_v2.example zone-xxx#ef-xxx#replica-name`
- **THEN** 资源被导入 state，后续 Read 按复合 ID 定位并刷新字段

### Requirement: provider registration for teo function replica v2
provider SHALL 在 `tencentcloud/provider.go` 注册 `"tencentcloud_teo_function_replica_v2"`（映射 `teo.ResourceTencentCloudTeoFunctionReplicaV2()`），并在 `tencentcloud/provider.md` 的 Resource 列表中追加该资源名。

#### Scenario: 注册完成
- **WHEN** 检查 provider.go 与 provider.md
- **THEN** 两个文件中均可找到 tencentcloud_teo_function_replica_v2 的注册条目

### Requirement: unit tests for teo function replica v2
SHALL 新增 `resource_tc_teo_function_replica_v2_test.go`，使用 gomonkey mock 云 API 客户端（不使用 terraform 测试套件），覆盖 Create 成功、Read 成功、Read 未找到（SetId 清空）、Update 成功、Update 不可变字段报错、Delete 成功场景，且保证代码可编译。

#### Scenario: 单元测试覆盖 CRUD
- **WHEN** 运行该测试文件中的用例（构建层面）
- **THEN** 所有 mock 场景断言请求参数与 state 结果，测试代码可正确构建执行

### Requirement: documentation for teo function replica v2
SHALL 新增 `tencentcloud/services/teo/resource_tc_teo_function_replica_v2.md`：一句话描述（含 TEO 产品名，"Provides a resource to ..."格式）、Example Usage（引用 teo zone/function 资源的 hcl 示例）、Import 部分说明复合 ID；SHALL NOT 手写 Argument Reference / Attribute Reference 部分（由工具自动生成）。

#### Scenario: 文档格式合规
- **WHEN** 检查 md 文件内容
- **THEN** 包含一句话描述、Example Usage、Import（含复合 ID 说明），且无 Argument/Attribute Reference 章节
