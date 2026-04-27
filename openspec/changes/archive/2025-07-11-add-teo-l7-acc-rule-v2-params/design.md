## Context

`tencentcloud_teo_l7_acc_rule_v2` 是 TEO（边缘安全加速平台）的七层加速规则 V2 资源，管理单个规则的完整生命周期（CRUD）。当前资源的 Read 操作通过 `DescribeTeoL7AccRules` API 查询规则，使用 `zone_id` 和 `rule_id` 构造内部过滤条件（`rule-id` filter），但未将 `filters` 参数暴露为用户可配置的 Terraform schema 字段。

### 当前实现

- **Read 操作**: `service.DescribeTeoL7AccRuleById(ctx, zoneId, ruleId)` 内部构造 `rule-id` 过滤条件
- **Schema**: 无 `filters` 字段，用户无法自定义查询过滤条件
- **SDK**: `DescribeL7AccRulesRequest.Filters` 为 `[]*Filter`，支持 `rule-id` 过滤

### SDK Filter 结构

```go
type Filter struct {
    Name   *string   `json:"Name,omitnil,omitempty"`
    Values []*string `json:"Values,omitnil,omitempty"`
}
```

### 约束
- 必须保持向后兼容：新增字段必须为 Optional
- 不能修改已有资源的 schema 行为
- vendor 模式管理依赖，SDK 已包含 Filter 支持

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源中添加 `filters` 可选参数
- 将 `filters` 参数传递给 `DescribeL7AccRules` API 的 `Filters` 字段
- 保持与现有资源的完全向后兼容

**Non-Goals:**
- 不修改 Create/Update/Delete 操作（filters 仅用于 Read）
- 不改变资源的 ID 格式（仍为 `zone_id#rule_id`）
- 不修改 v1 资源（`tencentcloud_teo_l7_acc_rule`）
- 不改变 DescribeTeoL7AccRuleById 的现有行为（当 filters 未指定时仍使用默认过滤条件）

## Decisions

### Decision 1: filters 作为 Optional 参数

**选择**: 在资源 schema 中添加 `filters` 作为 Optional 字段，不参与 Create/Update/Delete 操作

**理由**:
- `filters` 仅用于 DescribeL7AccRules 的查询过滤，属于 Read 操作范畴
- 作为 Optional 字段确保向后兼容
- 当用户未指定 filters 时，Read 操作使用现有的内部过滤逻辑（通过 rule-id 过滤）

**备选方案**:
- 将 filters 作为 Computed 字段：不合理，因为 filters 是查询输入参数，不是输出
- 将 filters 作为 Required 字段：破坏向后兼容性

### Decision 2: filters 的 Schema 结构

**选择**: TypeList，包含 `name`（String, Required）和 `values`（TypeList of String, Required）子结构

**理由**:
- 与 SDK `Filter` 结构一致（Name + Values）
- 与腾讯云 Terraform Provider 中其他资源/数据源的 filters 模式一致
- TypeList 保持过滤条件的顺序

### Decision 3: filters 在 Read 操作中的使用方式

**选择**: 当用户指定 filters 时，将用户指定的 filters 传递给 DescribeL7AccRules API；当用户未指定时，使用现有的内部过滤逻辑（通过 rule-id 过滤）

**理由**:
- 保持向后兼容：未指定 filters 时行为不变
- 提供灵活性：用户可以自定义过滤条件

### Decision 4: service 层方法更新

**选择**: 更新 `DescribeTeoL7AccRuleById` 方法签名，增加可选 filters 参数；或者新增一个方法支持自定义 filters

**理由**:
- 最小化变更：在现有方法上增加可选参数
- 保持现有行为：当 filters 为空时使用默认过滤逻辑

## Risks / Trade-offs

- **[Risk] filters 与 rule-id 过滤冲突** → 当用户指定 filters 时，不再自动添加 rule-id 过滤条件，而是完全使用用户指定的 filters。用户需要自行确保 filters 包含正确的过滤条件以查询到目标规则。
- **[Risk] filters 参数可能被滥用** → filters 仅用于查询过滤，不影响资源的创建、更新或删除。在文档中明确说明 filters 的用途和限制。
- **[Trade-off] 增加了 schema 复杂度** → 添加了新的可选参数，但这是为了提供更灵活的查询能力。
