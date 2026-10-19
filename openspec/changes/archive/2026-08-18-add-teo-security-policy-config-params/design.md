## Context

`tencentcloud_teo_security_policy_config` 是管理 TEO（EdgeOne）安全策略的 RESOURCE_KIND_GENERAL 资源，位于 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`（约 7100 行）。该资源封装了 `DescribeSecurityPolicy`（读取）与 `ModifySecurityPolicy`（创建/更新）两个云 API，覆盖安全策略全生命周期。

当前状态（已实现）：
- 顶层 schema：`zone_id`（Required, ForceNew）、`entity`（Optional, ForceNew, 校验取值 `ZoneDefaultPolicy`/`Template`/`Host`）、`host`（Optional, ForceNew）、`template_id`（Optional, ForceNew）
- 资源 ID 采用复合 ID：`zoneId#entity`（站点级）或 `zoneId#entity#host`（域名级）或 `zoneId#entity#templateId`（模板级），使用 `tccommon.FILED_SP` 作分隔符
- `security_config.rate_limit_config` 子结构已完整实现，包含 `switch`、`rate_limit_user_rules`、`rate_limit_template`（含 `rate_limit_template_detail` Computed 出参）、`rate_limit_intelligence`、`rate_limit_customizes`
- `rate_limit_user_rules` 与 `rate_limit_customizes` 复用同一 `rateLimitUserRuleSchema()`，字段含 `threshold`/`period`/`rule_name`/`action`/`punish_time`/`punish_time_unit`/`rule_status`/`acl_conditions`（`match_from`/`match_param`/`operator`/`match_content`）/`rule_priority`/`rule_id`(Computed)/`freq_fields`/`update_time`(Computed)/`freq_scope`/`name`/`custom_response_id`/`response_code`/`redirect_url`
- Create 不调用独立创建 API，直接转调 Update（`ModifySecurityPolicy`）；Delete 通过将安全配置置空来实现
- SDK 结构体：`teov20220901.DescribeSecurityPolicyRequest`、`ModifySecurityPolicyRequest`、`SecurityConfig`、`RateLimitConfig`、`RateLimitUserRule`、`RateLimitTemplate`、`RateLimitTemplateDetail`、`RateLimitIntelligence`、`AclCondition` 均已在 vendor 中存在

约束：
- Terraform Provider 向后兼容：本次为参数契约化记录，不改动任何 schema 类型与 ForceNew 属性
- 所有字段均已实现，无需新增 vendor 依赖
- `entity` / `host` / `template_id` 为 ForceNew，变更需重建资源

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_teo_security_policy_config` 建立完整的 openspec spec 规范，覆盖顶层查询入参与 `rate_limit_config` 速率限制子配置
- 契约化记录 `zone_id` / `entity` / `host` / `template_id` 的 ForceNew 行为与复合 ID 规则
- 契约化记录 `RateLimitConfig` 各子结构（用户规则、模板、智能过滤、托管定制规则）的字段映射，便于后续回归
- 不改动现有业务代码，确保零回归

**Non-Goals:**
- 不修改现有 schema 的类型或 ForceNew 属性
- 不新增云 API 调用或 vendor 依赖
- 不重构 `rateLimitUserRuleSchema()` 的复用结构
- 不改动 `security_policy`（SecurityPolicy，表达式语法路径）分支

## Decisions

### Decision 1: 以 ADDED Requirements 形式新建 capability spec

**选择**：在 `specs/teo-security-policy-config-resource/spec.md` 下使用 `## ADDED Requirements` 全量声明资源 schema、CRUD 行为与文档要求。

**备选**：拆成多个 capability（查询入参一个、rate_limit_config 一个）。

**理由**：
- 所有参数同属一个 Terraform 资源，拆分会导致 spec 碎片化且难以对应单一资源文件
- 与现有 `add-teo-mpgw-modify-status` 等提案保持单 capability 风格一致

### Decision 2: spec 中明确复合 ID 三种形态

**选择**：在 spec 中以 scenario 描述三种合法 `entity` 组合下的 ID 构成：
- `ZoneDefaultPolicy`：`zoneId#ZoneDefaultPolicy`，不设 `host`/`template_id`
- `Host`：`zoneId#Host#host`，仅设 `host`
- `Template`：`zoneId#Template#templateId`，仅设 `template_id`

**理由**：
- 复合 ID 是该资源 import 与 Read 的关键，非法组合会在 Create 阶段直接报错
- 明确记录可避免后续误改 ForceNew 或 ID 拼接逻辑

### Decision 3: Computed 出参字段在 spec 中标注

**选择**：在 spec 中将 `rule_id`、`update_time`、`rate_limit_template_detail` 及其子字段（`mode`/`id`/`action`/`punish_time`/`threshold`/`period`）、`rate_limit_intelligence.rule_id` 标注为 Computed（仅出参）。

**理由**：
- 与现有 schema 实现一致，避免 spec 误导用户将其作为入参配置
- 便于后续 `tasks` 阶段核对 Read 回填逻辑

### Decision 4: rate_limit_customizes 复用 rate_limit_user_rules 字段结构

**选择**：spec 中说明 `rate_limit_customizes` 与 `rate_limit_user_rules` 共享相同的字段集合（对应 SDK 的 `[]*RateLimitUserRule`）。

**理由**：
- 与现有 `rateLimitUserRuleSchema()` 复用实现一致
- 避免在 spec 中重复定义，减少维护成本

## Risks / Trade-offs

- **Risk**：spec 记录的参数与未来云 API 扩展不一致 → **Mitigation**：spec 聚焦当前已实现字段，后续 API 新增字段时通过新的 change 增量更新 spec
- **Trade-off**：本次为纯文档化 change，无代码改动，看似"无产出"，但为后续回归与 onboarding 提供契约基础 → 可接受
- **Risk**：spec 中遗漏某些 Computed 字段导致 archive 时信息丢失 → **Mitigation**：design 中已逐项列出 Computed 字段，spec 中对应标注
