## Context

TEO 推理服务（Inference Service）是 EdgeOne 提供的边缘推理服务，允许用户在边缘节点部署模型推理容器。云 API 已提供完整的推理服务管理接口，需要将其封装为 Terraform 资源。

当前 TEO 产品在 `tencentcloud/services/teo/` 下已有多个资源实现，遵循统一的代码模式：使用 `helper.Retry()` 处理重试、使用 `tccommon.LogElapsed()` 和 `tccommon.InconsistentCheck()` 处理错误。

### 云 API 分析

| 接口 | 用途 | 关键特征 |
|------|------|----------|
| `CreateInferenceService` | 创建 | 入参含 ZoneId、Name、ListenPort、Containers、ResourceConfig、RequestPaths、Description；出参返回 ServiceId |
| `DescribeInferenceServices` | 查询 | 支持按 ZoneId + Filters(service-id) 过滤，分页 Offset/Limit，排序 Order/Direction |
| `ModifyInferenceService` | 修改 | Containers 使用 `InferenceContainerConfigForModify`（与 Create 的 `InferenceContainerConfig` 结构相同但类型不同）；ResourceConfig 使用 `InferenceResourceConfigForModify`（无 HardwareSpec 字段） |
| `OperateInferenceService` | 操作 | 支持 Stop/Resume/Delete 三种操作 |

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_inference_service_v1` 资源的 CRUD 完整生命周期
- 支持推理服务的 Stop/Resume 操作（通过 `operation` 字段）
- 支持推理服务的容器配置、资源配置、请求路径等全部参数
- 正确处理 Create 和 Modify 接口中 Container/ResourceConfig 结构体类型差异

**Non-Goals:**
- 不实现 DATA_SOURCE 类型的数据源查询
- 不实现推理服务的部署日志查询（DescribeInferenceServiceDeploymentLogs）
- 不实现推理服务的监控数据查询（DescribeInferenceServiceMonitorData）
- 不改变现有 TEO 资源的任何行为

## Decisions

### 1. 资源 ID 设计

**决策**: 使用 `service_id`（CreateInferenceService 返回的 ServiceId）作为 Terraform 资源 ID。

**理由**: ServiceId 是唯一标识推理服务的字段，无需联合 ID。DescribeInferenceServices 支持按 service-id 过滤，可以直接通过单一 ID 查询。

### 2. Schema 参数设计

**决策**: 资源 Schema 定义以下参数：

| 参数名 | 类型 | 特性 | 说明 |
|--------|------|------|------|
| `zone_id` | TypeString | Required, ForceNew | 站点 ID，创建后不可变更 |
| `name` | TypeString | Required, ForceNew | 推理服务名称，创建后不可变更 |
| `listen_port` | TypeInt | Required | 监听端口 1-65535 |
| `containers` | TypeList | Required | 容器配置列表（嵌套结构） |
| `resource_config` | TypeList | Required | 资源配置（嵌套结构） |
| `request_paths` | TypeSet | Optional | 请求路径列表 |
| `description` | TypeString | Optional | 描述信息 |
| `operation` | TypeString | Optional | 操作类型：Stop/Resume（不持久化到 state） |
| `service_id` | TypeString | Computed | 服务 ID（由 API 返回） |
| `status` | TypeString | Computed | 服务状态 |
| `inference_url` | TypeString | Computed | 推理访问地址 |
| `create_time` | TypeString | Computed | 创建时间 |
| `update_time` | TypeString | Computed | 更新时间 |

**理由**: 参数映射直接对应云 API 入参。`operation` 字段用于触发 Stop/Resume，不作为持久化状态。

### 3. Containers 结构体类型差异处理

**决策**: 在 Terraform Schema 中统一使用一套嵌套结构，在 Create 时转换为 `InferenceContainerConfig`，在 Modify 时转换为 `InferenceContainerConfigForModify`。

**理由**: 两个结构体字段完全相同（ImageType、TcrRepositoryConfig、StartupCommand、EnvironmentVariables），仅在 SDK 中的类型名不同。统一 Schema 可以避免用户配置两份相同内容。

### 4. ResourceConfig 中 HardwareSpec 的处理

**决策**: Schema 中包含 `hardware_spec` 字段，仅在 Create 时传递，Modify 时不传递。

**理由**: ModifyInferenceService 的 `InferenceResourceConfigForModify` 没有 HardwareSpec 字段（硬件规格创建后不可修改），因此 Update 时跳过该字段。

### 5. Delete 实现

**决策**: 调用 `OperateInferenceService` 接口，设置 `Operation = "Delete"`。

**理由**: 云 API 未提供独立的 DeleteInferenceService 接口，删除操作通过 OperateInferenceService 统一处理。

### 6. Read 实现

**决策**: 调用 `DescribeInferenceServices`，使用 `ZoneId` + `Filters`（service-id 过滤，精确匹配）查询，分页 Limit 设为最大值 200。

**理由**: 查询接口是列表接口，需要通过 filter 精确匹配单个服务。为减少分页请求，Limit 设为最大值。

### 7. Update 实现

**决策**: 调用 `ModifyInferenceService`，传递所有可修改字段（ListenPort、RequestPaths、Containers、ResourceConfig、Description）。检测 `operation` 字段变化，若为 Stop/Resume 则调用 `OperateInferenceService`。

**理由**: ModifyInferenceService 支持全量更新，无需比对差异。Stop/Resume 操作与配置修改分离，使用独立的 API。

### 8. 异步操作处理

**决策**: OperateInferenceService（Delete/Stop/Resume）为异步操作，调用后需要轮询 DescribeInferenceServices 检查状态变更。

**理由**: 操作接口不返回操作结果，需要通过查询接口确认状态变更是否完成。

## Risks / Trade-offs

- **[风险] Create/Modify 接口中 Containers 和 ResourceConfig 使用不同结构体类型** → 缓解：在代码中显式创建不同类型的结构体实例，确保类型安全
- **[风险] OperateInferenceService 为异步接口，状态轮询可能超时** → 缓解：使用 `helper.Retry()` 并设置合理的超时时间（`tccommon.ReadRetryTimeout`）
- **[风险] `operation` 字段用于 Stop/Resume 时，若用户同时修改配置和 operation，需要确定执行顺序** → 缓解：先执行 ModifyInferenceService 修改配置，再执行 OperateInferenceService 变更状态
- **[权衡] ResourcesConfig 中 HardwareSpec 创建后不可修改** → 在 Schema 描述中说明该字段仅在创建时生效