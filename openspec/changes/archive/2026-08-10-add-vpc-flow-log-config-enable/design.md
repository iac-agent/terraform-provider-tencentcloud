## Context

VPC 流日志（Flow Log）的启用/停用状态目前无法通过 Terraform 声明式管理。腾讯云 VPC API 提供了 `EnableFlowLogs` 和 `DisableFlowLogs` 两个独立接口来控制流日志的启停状态，同时 `DescribeFlowLogs` 接口的返回结果中包含 `Enable` 字段表示当前启停状态。

该资源 (`tencentcloud_vpc_flow_log_config`) 专注于管理已有流日志的配置状态（启用/停用），不涉及流日志本身的创建和删除。流日志实例通过其他方式（如控制台或 `CreateFlowLog` API）创建后，由本资源配置其启停状态。

### 涉及的云 API

| API | 用途 | 关键参数 |
|-----|------|---------|
| `EnableFlowLogs` | 启动流日志 | `FlowLogIds []*string` |
| `DisableFlowLogs` | 停止流日志 | `FlowLogIds []*string` |
| `DescribeFlowLogs` | 查询流日志状态 | `FlowLogId *string`，返回 `FlowLog[].Enable *bool` |

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_vpc_flow_log_config` 资源，提供 `enable` 参数控制流日志启停
- 支持完整的 CRUD 生命周期：Create（设置启停状态）、Read（查询启停状态）、Update（修改启停状态）、Delete（从 state 中移除）
- 支持 Terraform import，允许导入已有的流日志配置

**Non-Goals:**
- 不涉及流日志实例本身的创建和删除（由 `tencentcloud_vpc_flow_log` 或其他方式管理）
- 不修改已有的流日志相关资源和数据源
- 不新增 `FlowLogName` 或其他描述性参数

## Decisions

### 1. 资源类型：RESOURCE_KIND_GENERAL

使用 RESOURCE_KIND_GENERAL 类型，管理资源的完整 CRUD 生命周期。`flow_log_id` 作为 Terraform 资源的唯一标识符（即 `d.SetId(flowLogId)`）。

### 2. Schema 设计

```go
"flow_log_id": {
    Type:        schema.TypeString,
    Required:    true,
    ForceNew:    true,
    Description: "Flow log ID.",
},
"enable": {
    Type:        schema.TypeBool,
    Required:    true,
    Description: "Whether to enable the flow log. true: enable, false: disable.",
},
```

- `flow_log_id`：必填，ForceNew（修改 ID 会销毁重建资源）
- `enable`：必填，控制启停状态

### 3. CRUD 操作映射

- **Create**：根据 `enable` 参数值，调用 `EnableFlowLogs`（enable=true）或 `DisableFlowLogs`（enable=false），将 `flow_log_id` 设置到 Terraform state
- **Read**：调用 `DescribeFlowLogs`，将返回的 `FlowLog[0].Enable` 同步到 state；若查询结果为空，设置 `d.SetId("")`
- **Update**：若 `enable` 值变更，调用对应的 `EnableFlowLogs` 或 `DisableFlowLogs` API
- **Delete**：仅从 Terraform state 中移除（`d.SetId("")`），不调用任何云 API（避免误删底层流日志实例）

### 4. API 调用策略

- 所有云 API 调用使用 `tccommon.ReadRetryTimeout` 作为超时时间
- 使用 `retry.Retry()` 包装 API 调用，错误通过 `tccommon.RetryError()` 返回
- `EnableFlowLogs` 和 `DisableFlowLogs` 的 `FlowLogIds` 参数传入 `[]*string{&flowLogId}`
- `DescribeFlowLogs` 使用 `FlowLogId` 参数精确查询单个流日志

### 5. Import 支持

支持通过 `terraform import tencentcloud_vpc_flow_log_config.foo flow_log_id` 导入已有流日志配置。导入时 `enable` 参数从 `DescribeFlowLogs` API 读取。

## Risks / Trade-offs

- **[风险] Delete 操作不调用云 API**：用户执行 `terraform destroy` 时不会删除底层流日志实例，仅移除 Terraform state。这是符合预期的设计，因为本资源只管理配置状态。若用户需要删除流日志实例本身，应使用其他方式。
- **[风险] EnableFlowLogs/DisableFlowLogs 是异步操作**：云 API 调用成功后，状态可能存在短暂延迟。Read 操作中已包含重试逻辑，能够应对最终一致性问题。
- **[权衡] EnableFlowLogs/DisableFlowLogs 接受批量 ID**：云 API 接受 `FlowLogIds` 数组，但本资源仅传入单个 `flow_log_id`，简化了 Terraform 资源与 API 的映射关系。
