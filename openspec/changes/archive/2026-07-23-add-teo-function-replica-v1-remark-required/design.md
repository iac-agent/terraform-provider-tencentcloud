## Context

TEO（EdgeOne）边缘函数副本（Function Replica）是边缘函数的版本管理能力。当前 Terraform Provider 中已实现 `tencentcloud_teo_function` 资源，但尚未实现函数副本资源。云 API 已提供完整的 CRUD 接口：
- `CreateFunctionReplica`：创建函数副本（参数：ZoneId, FunctionId, ReplicaName, Content, Remark）
- `DescribeFunctionReplicas`：查询函数副本列表（参数：ZoneId, FunctionId, Offset, Limit, SortBy, SortOrder, Filters）
- `ModifyFunctionReplica`：修改函数副本（参数：ZoneId, FunctionId, ReplicaName, Content, Remark）
- `DeleteFunctionReplica`：删除函数副本（参数：ZoneId, FunctionId, ReplicaNames）

DescribeFunctionReplicas 返回 `FunctionReplica` 结构体，包含：FunctionId, ReplicaName, Content, Remark, CreatedOn, ModifiedOn。

DeleteFunctionReplica 使用 ReplicaNames（列表）进行删除，而不是单个名称。这意味着该资源删除时需要将 ReplicaName 放入列表中传入。

关键约束：`remark` 字段需设为 Required（必填），区别于 `tencentcloud_teo_function` 中 remark 为 Optional 的设计。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_replica_v1` 资源，实现完整 CRUD
- `remark` 字段设为 Required（必填）
- 使用 `zone_id` + `function_id` + `replica_name` 作为复合 ID（因 DescribeFunctionReplicas 返回列表需通过 ReplicaName 定位具体副本）
- 在 provider.go 和 provider.md 中注册该资源
- 生成对应的 .md 文档
- 编写单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_function` 资源
- 不新增数据源
- 不实现 import 功能（因 DeleteFunctionReplica 使用列表删除，且副本的复合 ID 已足够定位）

## Decisions

### 1. 复合 ID 设计：zone_id + function_id + replica_name

**选择**：使用 `zone_id#function_id#replica_name` 作为 Terraform Resource ID

**理由**：
- DescribeFunctionReplicas 需通过 ZoneId + FunctionId 查询，再通过 ReplicaName 在返回列表中定位
- ReplicaName 在同一 FunctionId 下唯一
- 三字段组合可唯一定位一个函数副本

**替代方案**：使用 zone_id + function_id 作为 ID，但无法在返回列表中精确定位具体副本

### 2. remark 设为 Required

**选择**：将 remark 字段在 schema 中设为 Required: true

**理由**：需求明确要求 remark 为必填，确保每个副本都有描述信息

### 3. 不支持 Import

**选择**：不添加 Importer

**理由**：用户可通过指定 zone_id、function_id、replica_name 直接创建资源。如需 import 也可后续添加。

### 4. Delete API 适配

**选择**：DeleteFunctionReplica 的 ReplicaNames 参数为列表类型，删除时将 replica_name 作为单元素列表传入

**理由**：虽然 API 支持批量删除，但 Terraform 资源每次只管理一个副本

### 5. Read 逻辑：从列表中过滤

**选择**：调用 DescribeFunctionReplicas 后，通过 ReplicaName 在返回列表中匹配目标副本

**理由**：API 无单个副本查询接口，只能从列表中过滤。使用 Filters 的 replica-name 进行精确过滤。

## Risks / Trade-offs

- [API 返回列表可能为空] → 在 Read 中检查 response 和列表是否为空，若为空则 d.SetId("") 标记资源已删除
- [remark 必填可能导致用户使用不便] → 这是明确需求，用户必须提供副本描述
- [DeleteFunctionReplica 使用列表删除，可能误删] → 只传入单个 replica_name，避免批量删除风险
