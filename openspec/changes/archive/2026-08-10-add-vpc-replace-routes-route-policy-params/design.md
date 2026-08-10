## Context

`tencentcloud_vpc_replace_routes_with_route_policy_config` 是一个 RESOURCE_KIND_GENERAL 资源，通过 `ReplaceRoutesWithRoutePolicy` 接口执行路由策略替换，通过 `DescribeRouteTables` 接口读取路由表状态。

当前状态：
- 资源 schema 已包含 `route_table_id`（必填，ForceNew）和 `routes`（必填）两个字段。
- Read 方法调用 `VpcService.DescribeRouteTables(ctx, routeTableId, "", "", map[string]string{}, nil, "")`，仅通过 `route-table-id` filter 查询，未使用云API提供的 `NeedRouterInfo`、`Filters.Name`、`Filters.Values` 参数，也未将 `TotalCount` 出参回填。
- service 层 `DescribeRouteTables` 方法签名固定为 `(ctx, routeTableId, routeTableName, vpcId, tags, associationMain, tagKey)`，不支持 `NeedRouterInfo`、`Name`、`Values`。

云API能力（`DescribeRouteTablesRequest`）：
- `Filters []*Filter`：每个 `Filter` 含 `Name *string` 与 `Values []*string`。
- `NeedRouterInfo *bool`：是否需要获取路由策略信息。
- 响应 `TotalCount *uint64`。

约束：
- 必须保持向后兼容，新增字段均为 Optional/Computed。
- service 层 `DescribeRouteTables` 被多个资源复用，签名扩展需保持对现有调用方的兼容。

## Goals / Non-Goals

**Goals:**
- 为资源 schema 新增 `need_router_info`（Optional bool）、`name`（Optional string）、`values`（Optional list）、`total_count`（Computed int）4 个字段。
- 扩展 service 层 `DescribeRouteTables` 方法，支持透传 `needRouterInfo`、`name`、`values`。
- 在资源 Read 方法中将新增入参透传给 service 层，并将 `total_count` 写回 state。

**Non-Goals:**
- 不修改 `ReplaceRoutesWithRoutePolicy` 接口的调用逻辑（`route_table_id`、`routes` 已存在且满足需求）。
- 不修改 Create/Update/Delete 方法的核心业务逻辑。
- 不变更资源类型（仍为 RESOURCE_KIND_GENERAL）。

## Decisions

### Decision 1: service 层 DescribeRouteTables 方法签名扩展方式
**选择**：通过新增独立方法或扩展现有方法签名增加 `needRouterInfo *bool`、`name string`、`values []*string` 参数。

**理由**：现有 `DescribeRouteTables` 被 `tencentcloud_vpc_route_table` 等多个资源复用。为最小化影响，采用扩展签名方式，新增参数放在参数列表末尾，现有调用方传零值即可保持兼容。`name` 与 `values` 组合构成一个 `Filter`（与现有 `route-table-id` 等 filter 并列追加），`needRouterInfo` 直接赋值到 `request.NeedRouterInfo`。

**备选方案**：引入 options struct 传参。否决原因：与现有代码风格（多参数平铺）不一致，且本次仅新增 3 个参数，平铺更直观。

### Decision 2: name/values 的 schema 类型
**选择**：`name` 为 `schema.TypeString`（Optional），`values` 为 `schema.TypeList`（Optional，Elem 为 TypeString）。

**理由**：云API `Filter.Name` 是单值 string，`Filter.Values` 是 `[]*string` 列表。二者独立暴露为顶层参数，符合"列表展开、平铺到顶层"的规范。Read 时将 `name`+`values` 组装为一个 `Filter` 追加到 `request.Filters`。

### Decision 3: total_count 作为 Computed 字段
**选择**：`total_count` 为 `schema.TypeInt`，`Computed: true`。

**理由**：`TotalCount` 是云API响应出参，用户不应手动设置。Computed 字段会在 Read 时由 `d.Set("total_count", ...)` 写入 state。

### Decision 4: Read 方法空响应处理
**选择**：当 `DescribeRouteTables` 返回空时，先 `log.Printf("[CRUD] ...")` 保留 id 现场，再 `d.SetId("")`。

**理由**：遵循项目规范，避免日志中无法定位是哪一次调用导致 id 被清空。

## Risks / Trade-offs

- **[Risk] service 层方法签名变更影响其他调用方** → 通过将新参数置于末尾、现有调用方传零值保持兼容；实现阶段需全局检查所有调用点。
- **[Risk] `name`/`values` 组合语义模糊** → 约定二者需配合使用（`name` 指定 filter 维度，`values` 指定取值），单独设置 `values` 而 `name` 为空时不追加 filter，避免构造出不完整的 Filter。
- **[Risk] `total_count` 在分页场景下的含义** → service 层 `DescribeRouteTables` 内部分页聚合，`total_count` 取首次响应的 `TotalCount`，反映全量匹配数而非单页数量。
