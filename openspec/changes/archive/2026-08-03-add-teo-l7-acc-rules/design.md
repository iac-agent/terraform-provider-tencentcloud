## Context

TEO (EdgeOne) 产品已有 `tencentcloud_teo_l7_acc_rule`（旧版，使用 ImportZoneConfig 接口）和 `tencentcloud_teo_l7_acc_rule_v2`（新版，使用新 CRUD 接口）两个单规则管理资源。本次新增 `tencentcloud_teo_l7_acc_rules`（复数）资源，支持在单个 Terraform 资源中批量管理多条七层加速规则，使用新的 L7AccRules 系列 CRUD 接口。

SDK 中已有完整的四个接口：
- `CreateL7AccRules`：批量创建规则，入参为 `Rules` 列表，返回 `RuleIds` 列表
- `DescribeL7AccRules`：查询规则，支持 `Filters` 过滤（可按 `rule-id` 过滤）
- `ModifyL7AccRule`：修改单条规则（注意：是单数接口，每次只能修改一条）
- `DeleteL7AccRules`：批量删除规则，入参为 `RuleIds` 列表

## Goals / Non-Goals

**Goals:**
- 创建 `tencentcloud_teo_l7_acc_rules` 资源，实现完整的 CRUD 生命周期管理
- 使用 `zone_id` 作为资源 ID（批量资源，一个 zone 下所有规则由该资源统一管理）
- 支持 `rules` 字段批量管理多条规则，每条规则包含 `RuleEngineItem` 的所有字段
- 支持 Import 功能

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_l7_acc_rule` 和 `tencentcloud_teo_l7_acc_rule_v2` 资源
- 不涉及规则优先级批量调整（`ModifyL7AccRulePriority` 是独立接口）

## Decisions

### 1. 资源 ID 设计
- **决策**: 使用 `zone_id` 作为资源 ID（`d.SetId(zoneId)`）
- **理由**: 该资源是批量资源，管理一个 zone 下的所有规则集合。`zone_id` 唯一标识一个规则集合，无需拼接其他字段。
- **替代方案**: 使用 zone_id + 所有 rule_ids 拼接，但 rule_ids 数量可变会导致 ID 过长且不稳定。

### 2. Schema 设计
- **决策**: 顶层 `zone_id`（Required, ForceNew），顶层 `rules` 列表（Optional, Computed），`rule_ids` 作为 Computed 字段
- **理由**: 
  - `zone_id` 设置 ForceNew 因为规则集合绑定 zone，不支持跨 zone 迁移
  - `rules` 是 Optional（创建时提供）和 Computed（Read 时回填），每条规则内部结构参考 `RuleEngineItem`
  - `rule_ids` 在 Create 后从 API 返回填充，仅 Computed

### 3. Create 流程
- **决策**: 调用 `CreateL7AccRules` 批量创建 `rules` 列表中的所有规则
- **流程**: 构建 `Rules []*RuleEngineItem` → 调用 API → 获取 `RuleIds` → 设置 `zone_id` 为 ID → 调用 Read 回填

### 4. Read 流程
- **决策**: 调用 `DescribeL7AccRules` 查询所有规则，分页获取全量数据
- **流程**: 复用已有的 `DescribeTeoL7AccRuleById` 服务方法（传入空 `ruleId` 获取全部规则），设置 `Limit` 为最大值（1000）以避免分页遗漏

### 5. Update 流程
- **决策**: 检测 `rules` 列表变更，对变更的规则调用 `ModifyL7AccRule`（每次调用只修改一条规则）
- **流程**: 
  1. 获取新旧 `rules` 列表
  2. 对比找出新增、修改、删除的规则
  3. 对新增的规则调用 `CreateL7AccRules`
  4. 对修改的规则调用 `ModifyL7AccRule`（需要 `rule_id`）
  5. 对删除的规则调用 `DeleteL7AccRules`（需要 `rule_ids`）
- **风险**: 部分更新失败时状态不一致，需要 Read 回填以暴露 drift

### 6. Delete 流程
- **决策**: 调用 `DeleteL7AccRules` 删除所有规则
- **流程**: 收集当前所有 `rule_ids` → 调用 API 批量删除 → 调用 Read 确认

### 7. Import 支持
- **决策**: 支持 `terraform import`，使用 `zone_id` 作为导入 ID
- **理由**: 资源 ID 就是 `zone_id`，使用 `schema.ImportStatePassthrough` 即可

## Risks / Trade-offs

- [风险] ModifyL7AccRule 是单数接口，批量更新需要逐条调用，性能较差 → 缓解：Terraform 资源更新场景下规则变更通常数量有限
- [风险] 部分 API 调用失败时，已成功的规则状态与 Terraform state 不一致 → 缓解：最终调用 Read 回填，利用 Terraform 的 drift detection 机制
- [风险] `DescribeL7AccRules` 默认分页 limit=20，可能遗漏规则 → 缓解：设置 `Limit=1000`（API 最大值）