## Context

`tencentcloud_teo_function` 是一个 RESOURCE_KIND_GENERAL 类型的 Terraform 资源，管理 TEO（边缘安全加速平台）的边缘函数。当前资源的 Read 操作通过 `DescribeFunctions` API 查询单个函数，但仅使用了 `ZoneId` 和 `FunctionIds` 两个入参，且出参仅映射了单个函数的字段。

`DescribeFunctions` API 还支持 `Filters` 入参（按 name、remark 模糊匹配过滤），以及 `Functions` 出参返回函数列表。当前资源未暴露这些参数，用户无法通过 Terraform 进行批量查询或条件过滤。

当前资源 Schema 已有 `zone_id`（Required, ForceNew）字段，本次变更只需新增 `function_ids`、`filters`、`functions` 三个参数。

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_teo_function` 资源新增 `function_ids` 参数，支持按函数 ID 列表过滤查询
- 为 `tencentcloud_teo_function` 资源新增 `filters` 参数，支持按 name、remark 条件过滤查询
- 为 `tencentcloud_teo_function` 资源新增 `functions` 计算属性，返回 DescribeFunctions API 查询到的函数列表
- 保持向后兼容，所有新增参数均为 Optional 或 Computed，不影响现有 Terraform 配置

**Non-Goals:**
- 不创建新的数据源（data source），仅在现有资源中新增参数
- 不修改 Create/Update/Delete 流程
- 不修改 `zone_id` 参数的现有行为

## Decisions

### Decision 1: 参数 Schema 设计

**选择**: 新增 3 个参数到现有资源 Schema

| 参数 | 类型 | 属性 | 对应 API 字段 | 说明 |
|------|------|------|--------------|------|
| `function_ids` | `TypeList` of `TypeString` | Optional | `request.FunctionIds` | 按函数 ID 列表过滤，用于批量查询 |
| `filters` | `TypeList` of `TypeMap` / Object | Optional | `request.Filters` | 过滤条件，支持 name 和 remark 字段 |
| `functions` | `TypeList` of `TypeMap` / Object | Computed | `response.Functions` | 返回查询到的函数详情列表 |

**理由**: 所有新增参数均为 Optional 或 Computed，不会破坏现有配置。`function_ids` 和 `filters` 作为查询条件仅在 Read 时使用，`functions` 作为计算属性返回查询结果。

### Decision 2: filters 参数结构设计

**选择**: 使用 `TypeList` 嵌套 `TypeSet` of `TypeString` 结构

```go
"filters": {
    Type:     schema.TypeList,
    Optional: true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "values": {
                Type:     schema.TypeSet,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
                Set:      schema.HashString,
            },
        },
    },
}
```

**理由**: 与 TEO 其他资源（如 `resource_tc_teo_function_rule`）中 Filter 的使用方式保持一致，符合 SDK 中 `Filter` 结构体的定义（Name + Values）。

### Decision 3: functions 参数结构设计

**选择**: 使用 `TypeList` 嵌套 Object 结构

```go
"functions": {
    Type:     schema.TypeList,
    Computed: true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "function_id": {Type: schema.TypeString, Computed: true},
            "zone_id":     {Type: schema.TypeString, Computed: true},
            "name":        {Type: schema.TypeString, Computed: true},
            "remark":      {Type: schema.TypeString, Computed: true},
            "content":     {Type: schema.TypeString, Computed: true},
            "domain":      {Type: schema.TypeString, Computed: true},
            "create_time": {Type: schema.TypeString, Computed: true},
            "update_time": {Type: schema.TypeString, Computed: true},
        },
    },
}
```

**理由**: `functions` 出参对应 SDK 的 `Function` 结构体，包含 `FunctionId`、`ZoneId`、`Name`、`Remark`、`Content`、`Domain`、`CreateTime`、`UpdateTime` 字段，需要完整映射。

### Decision 4: Read 函数修改策略

**选择**: 修改现有 Read 函数，在调用 `DescribeFunctions` 时传入新增的 `function_ids` 和 `filters` 参数，并处理 `functions` 出参

**理由**: 当前 Read 函数通过 `service.DescribeTeoFunctionById` 服务层方法查询，该方法固定传 `ZoneId` + `FunctionIds`（单个 ID）。为支持新增的 `function_ids`（多个 ID）和 `filters` 参数，需要扩展 Read 函数直接调用 `DescribeFunctions` API，或在服务层新增方法。

**实现方案**: 在 Read 函数中，保持现有的通过 ID 拆分获取 `zoneId` 和 `functionId` 的逻辑不变，新增对 `function_ids` 和 `filters` 参数的处理。当用户设置了这些参数时，将其传入 DescribeFunctions 请求中。对于 `functions` 出参，将其映射到 computed 属性中。

## Risks / Trade-offs

- **[向后兼容性风险]** → 新增参数均为 Optional 或 Computed，不会破坏现有配置。Read 函数的核心逻辑（通过 ID 拆分获取单个函数信息）保持不变
- **[查询性能风险]** → 当 `function_ids` 列表很大或 `filters` 范围过宽时，可能返回大量数据。DescribeFunctions API 有分页支持（Offset/Limit），但本资源是通用资源类型，不需要实现完整分页
- **[functions 与顶层字段重叠]** → `functions` 列表中的字段与资源顶层字段（name、remark、content 等）有重叠，但这是查询结果列表，与单资源字段含义不同
