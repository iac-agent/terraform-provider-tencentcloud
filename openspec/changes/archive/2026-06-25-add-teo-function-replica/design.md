## Context

EdgeOne (TEO) 边缘函数副本（Function Replica）允许同一函数存在多个副本版本，每个副本有独立的内容和描述。函数规则通过权重或地域选择不同的函数副本执行。云 API 已提供完整的 CRUD 接口：`CreateFunctionReplica`、`DescribeFunctionReplicas`、`ModifyFunctionReplica`、`DeleteFunctionReplica`。

当前 Terraform Provider 中已有 `tencentcloud_teo_function` 资源管理边缘函数，但缺少对函数副本的管理。本次新增 `tencentcloud_teo_function_replica` 资源（RESOURCE_KIND_GENERAL）来填补这一空白。

参考资源：`tencentcloud_teo_function`（位于 `tencentcloud/services/teo/resource_tc_teo_function.go`）。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_function_replica` 资源的完整 CRUD 操作
- 支持通过 `terraform import` 导入已有函数副本
- 提供 `.md` 文档样例用于 `make doc` 自动生成文档

**Non-Goals:**
- 不修改已有 `tencentcloud_teo_function` 资源
- 不提供 DATASOURCE 类型资源

## Decisions

### 1. 资源 ID 使用联合 ID

**选择**: `zone_id#function_id#replica_name`（使用 `tccommon.FILED_SP` 分隔符）。

**理由**: `CreateFunctionReplica` 接口的返回值仅包含 `RequestId`，不返回副本 ID。函数副本在 API 层面由 `ZoneId + FunctionId + ReplicaName` 三元组唯一标识。`ModifyFunctionReplica` 和 `DeleteFunctionReplica` 同样使用这三个字段定位副本。

**替代方案**: 无。API 层面没有副本 ID 字段。

### 2. Schema 设计

**必填且 ForceNew 字段**:
- `zone_id` (TypeString, Required, ForceNew) — 站点 ID，变更意味着不同站点
- `function_id` (TypeString, Required, ForceNew) — 函数 ID，变更意味着不同函数
- `replica_name` (TypeString, Required, ForceNew) — 副本名称，作为三元组标识符的一部分

**必填非 ForceNew 字段**:
- `content` (TypeString, Required) — 副本函数内容（JavaScript 代码），可通过 Modify 接口更新

**可选字段**:
- `remark` (TypeString, Optional) — 副本描述，可通过 Modify 接口更新

**计算字段** (Computed):
- `created_on` (TypeString, Computed) — 副本创建时间，来自 `DescribeFunctionReplicas` 返回的 `FunctionReplica.CreatedOn`
- `modified_on` (TypeString, Computed) — 副本更新时间，来自 `DescribeFunctionReplicas` 返回的 `FunctionReplica.ModifiedOn`

**不纳入 Schema 的查询参数**:
`DescribeFunctionReplicas` 请求中的 `sort_by`、`sort_order`、`filters` 是查询控制参数，不反映资源状态。在 Read 方法内部通过 `filters` 精确匹配 `replica_name` 来定位资源。

### 3. CRUD 实现策略

**Create**:
- 调用 `CreateFunctionReplica`，传入 `zone_id`, `function_id`, `replica_name`, `content`, `remark`
- 调用成功后，直接构造联合 ID 设置 `d.SetId(fmt.Sprintf("%s#%s#%s", zoneId, functionId, replicaName))`
- 再调用 Read 方法确认资源存在并同步状态

**Read**:
- 从 `d.Id()` 解析出 `zone_id`, `function_id`, `replica_name`
- 调用 `DescribeFunctionReplicas`，设置 `Filters: [{Name: "replica-name", Values: [replica_name]}]`
- 若返回列表为空，则 `d.SetId("")`
- 若返回列表不为空，取第一个元素设置各字段

**Update**:
- 调用 `ModifyFunctionReplica`，传入 `zone_id`, `function_id`, `replica_name`, `content`, `remark`
- 调用后执行 Read 同步状态

**Delete**:
- 调用 `DeleteFunctionReplica`，传入 `zone_id`, `function_id`, `ReplicaNames: [replica_name]`

### 4. Retry 策略

所有云 API 调用均使用 `tccommon.ReadRetryTimeout` 超时，通过 `helper.Retry` 包装，错误使用 `tccommon.RetryError` 包装后返回。Retry 块内仅执行 API 调用。

### 5. Update 不可变字段校验

由于 `ModifyFunctionReplica` 仅支持更新 `content` 和 `remark`，而 `zone_id`、`function_id`、`replica_name` 已标记 `ForceNew`，Terraform 层会自动处理这些字段变更时的重建逻辑，无需在 Update 方法中额外校验。

## Risks / Trade-offs

- **风险**: `CreateFunctionReplica` 是幂等接口，重复调用可能返回成功。若用户修改 `replica_name`（通过 ForceNew 重建），旧副本不会被自动删除，需手动清理。
  - **缓解**: 这是 Terraform ForceNew 的标准行为，与 API 设计一致。在文档中说明。
-
- **风险**: `DescribeFunctionReplicas` 通过 `replica-name` 过滤器进行模糊查询，可能匹配到多个结果（如 "foo" 可能匹配 "foo-v1" 和 "foo"）。
  - **缓解**: 在 Read 方法中对返回结果做精确匹配（`replica_name == returnedReplicaName`），若模糊匹配返回多个结果但无精确匹配，视为资源不存在。
