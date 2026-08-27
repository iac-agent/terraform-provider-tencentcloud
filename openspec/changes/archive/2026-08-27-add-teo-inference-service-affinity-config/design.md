## Context

TEO 推理服务（Inference Service）是腾讯云 EdgeOne 的推理服务产品，允许用户部署模型推理服务。当前 Terraform Provider 中尚未创建该资源的 Terraform 实现。

云 API 层面已提供完整的 CRUD 接口：
- `CreateInferenceService`：创建推理服务，返回 `ServiceId`
- `DescribeInferenceServices`：查询推理服务列表（支持 Filters 过滤）
- `ModifyInferenceService`：修改推理服务配置
- `OperateInferenceService`：操作推理服务（Stop/Resume/Delete）

SDK 中 `InferenceAffinityConfig` 和 `SessionIdAffinityConfig` 结构体已定义，`CreateInferenceService` 和 `ModifyInferenceService` 均支持 `AffinityConfig` 参数。`DescribeInferenceServices` 返回的 `InferenceService` 结构体未直接包含 `AffinityConfig`，但 `InferenceServiceConfig` 包含该字段（需通过单独接口获取配置详情或从 `InferenceService` 的平铺字段间接获取）。

## Goals / Non-Goals

**Goals:**
- 创建 `tencentcloud_teo_inference_service_v1` Terraform 资源（RESOURCE_KIND_GENERAL），提供完整的 CRUD 生命周期管理
- 在资源 Schema 中支持 `affinity_config` 嵌套参数块，包含 `switch`、`affinity_mode`、`source`、`header_name` 四个子字段
- Create 操作调用 `CreateInferenceService` API 时传入 `AffinityConfig`
- Update 操作调用 `ModifyInferenceService` API 时传入 `AffinityConfig`
- Read 操作通过 `DescribeInferenceServices` 获取服务信息并将 `AffinityConfig` 回写到 state
- Delete 操作调用 `OperateInferenceService`（Operation="Delete"）删除服务

**Non-Goals:**
- 不修改现有任何 TEO 资源
- 不创建数据源（DATA_SOURCE）资源
- 不实现 `InferenceServiceConfig` 以外的其他 InferenceService 子资源
- 不处理异步操作的等待逻辑（如部署状态轮询），使用 SDK 默认行为

## Decisions

### 1. 资源 Schema 设计

采用与 `tencentcloud_teo_multi_path_gateway` 一致的风格，顶层字段使用 snake_case 命名。

**核心字段：**
- `zone_id`（Required, ForceNew）：站点 ID
- `name`（Required, ForceNew）：推理服务名称
- `listen_port`（Required）：监听端口
- `containers`（Required）：容器配置列表
- `resource_config`（Required）：资源配置
- `request_paths`（Optional）：请求路径列表
- `description`（Optional）：描述信息
- `affinity_config`（Optional）：亲和性配置嵌套块

**`affinity_config` 嵌套块：**
- `switch`（TypeString, Optional）：亲和总开关，枚举值 `On`/`Off`
- `affinity_mode`（TypeString, Optional）：亲和方式，枚举值 `SessionId`
- `source`（TypeString, Optional）：会话 ID 传递位置，枚举值 `Header`
- `header_name`（TypeString, Optional）：请求头名称

**理由：** 将 `affinity_config` 设计为 `TypeList` + `MaxItems: 1` 的嵌套块，与 Terraform 最佳实践一致，便于用户以结构化方式配置亲和性参数。

### 2. ID 设计

使用 `ServiceId` 作为 Terraform 资源的唯一 ID。Create 接口返回 `ServiceId`，Read/Update/Delete 接口均使用 `ServiceId` 定位资源。

**理由：** 云 API 的 `ServiceId` 在站点内全局唯一，无需使用复合 ID。与 `tencentcloud_teo_multi_path_gateway` 等资源保持一致。

### 3. Update 逻辑

使用 `ModifyInferenceService` API 进行更新。由于该 API 不支持部分更新（所有字段需全量传入），每次 Update 需要将 schema 中的所有字段（包括未变更的字段）传递给 API。

**理由：** 云 API 设计为全量更新，Terraform 侧无法判断哪些字段被用户修改，全量传入是最安全的方式。

### 4. Delete 逻辑

使用 `OperateInferenceService` API，设置 `Operation="Delete"` 进行删除。删除后调用 `DescribeInferenceServices` 确认资源不存在。

**理由：** 云 API 未提供独立的 `DeleteInferenceService` 接口，通过 `OperateInferenceService` 的 Delete 操作实现删除。

### 5. Read 逻辑

通过 `DescribeInferenceServices` 使用 `Filters` 按 `service-id` 过滤查询单个服务。由于 `InferenceService` 响应结构体不直接包含 `AffinityConfig`（该字段在 `InferenceServiceConfig` 中），Read 操作需要处理 `AffinityConfig` 可能为 nil 的情况。

**替代方案考虑：** 使用 `DescribeInferenceServices` 获取服务基本信息，`AffinityConfig` 在创建/更新后保存在 state 中，Read 时若 API 未返回则保留 state 中的值。但此方案可能导致 state 漂移。最终决定：Read 时尽力从 API 响应中获取 `AffinityConfig`，若 API 未返回则保持 state 中已有值不变。

## Risks / Trade-offs

- **[Risk] Read 操作中 AffinityConfig 可能为 nil**：`InferenceService` 结构体可能不包含 `AffinityConfig` 字段，导致 Read 后 state 中该字段被清空。→ **Mitigation**：Read 时检查响应中 `AffinityConfig` 是否为 nil，若为 nil 则不更新 state 中对应的嵌套块。
- **[Risk] 异步操作**：Create/Update/Delete 操作可能是异步的，资源状态可能不是立即可用。→ **Mitigation**：通过 retry 机制处理最终一致性，设置合理的 Timeout。
- **[Risk] `InferenceService` 结构体不包含 `AffinityConfig`**：当前 SDK 中 `InferenceService` 未直接包含 `AffinityConfig`，Read 时无法回填。→ **Mitigation**：在 Read 方法中，若 `AffinityConfig` 为 nil，不调用 `d.Set("affinity_config", ...)`，保留 state 中的值。

## Open Questions

- `InferenceService` 响应结构体是否会在后续 SDK 版本中增加 `AffinityConfig` 字段？当前 SDK 版本为 `v1.3.170`，`InferenceService` 不直接包含该字段。