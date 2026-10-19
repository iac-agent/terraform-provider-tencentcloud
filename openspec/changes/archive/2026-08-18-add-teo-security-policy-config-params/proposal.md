## Why

`tencentcloud_teo_security_policy_config` 资源当前已实现了 TEO 安全策略的 CRUD 管理，但其 schema 与云 API `DescribeSecurityPolicy` / `ModifySecurityPolicy` 的参数映射尚未在 openspec 规范中沉淀。为了让安全策略配置（尤其是速率限制 `RateLimitConfig` 下的用户规则、模板、智能客户端过滤、托管定制规则，以及 `DescribeSecurityPolicy` 的查询入参 `zone_id` / `entity` / `host` / `template_id`）的契约可被追溯与回归，需要为该资源建立完整的 spec 规范，记录现有实现所支持的参数与行为。

## What Changes

- 为 `tencentcloud_teo_security_policy_config` 资源新增（规范化记录）顶层查询参数：
  - `zone_id`（必填，ForceNew）：站点 ID，对应 `DescribeSecurityPolicyRequest.ZoneId`
  - `entity`（可选）：安全策略类型，对应 `DescribeSecurityPolicyRequest.Entity`
  - `host`（可选）：域名，对应 `DescribeSecurityPolicyRequest.Host`
  - `template_id`（可选）：策略模板 ID，对应 `DescribeSecurityPolicyRequest.TemplateId`
- 规范化记录 `security_config.rate_limit_config` 下的速率限制参数（对应 `ModifySecurityPolicyRequest.SecurityConfig.RateLimitConfig`）：
  - `switch`：速率限制总开关
  - `rate_limit_user_rules`：用户自定义速率限制规则列表，每条规则包含 `threshold`、`period`、`rule_name`、`action`、`punish_time`、`punish_time_unit`、`rule_status`、`acl_conditions`（含 `match_from`、`match_param`、`operator`、`match_content`）、`rule_priority`、`rule_id`（Computed）、`freq_fields`、`update_time`（Computed）、`freq_scope`、`name`、`custom_response_id`、`response_code`、`redirect_url`
  - `rate_limit_template`：速率限制模板，含 `mode`、`action` 及 `rate_limit_template_detail`（Computed，含 `mode`、`id`、`action`、`punish_time`、`threshold`、`period`）
  - `rate_limit_intelligence`：智能客户端过滤，含 `switch`、`action`、`rule_id`（Computed）
  - `rate_limit_customizes`：托管定制规则列表，字段结构与 `rate_limit_user_rules` 一致
- 上述字段均为对现有资源实现的参数契约化记录，不引入破坏性变更。

## Capabilities

### New Capabilities
- `teo-security-policy-config-resource`: 管理 TEO 安全策略配置的 Terraform 资源，覆盖站点级/域名级/模板级安全策略的创建、读取、更新、删除全生命周期，重点规范化 `zone_id`/`entity`/`host`/`template_id` 查询入参及 `security_config.rate_limit_config` 速率限制子配置的参数契约。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`：现有资源已实现本次规范化的全部参数，无需改动业务代码
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`：现有测试覆盖 CRUD，无需新增
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`：现有文档已包含示例
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.DescribeSecurityPolicyRequest` / `ModifySecurityPolicyRequest` 及 `RateLimitConfig` / `RateLimitUserRule` / `RateLimitTemplate` / `RateLimitTemplateDetail` / `RateLimitIntelligence` / `AclCondition` 结构体，无需变更 vendor。
- 向后兼容：本次为纯参数契约化记录，无 schema 变更，不影响已有 state 与 TF 配置。
- 文档：`website/docs/` 由 `make doc` 自动生成流程读取 `.md` 文件，无需手动改动。
