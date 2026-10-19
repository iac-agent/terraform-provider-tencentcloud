## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_security_policy_config` with the following top-level schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID，对应 `DescribeSecurityPolicyRequest.ZoneId` 与 `ModifySecurityPolicyRequest.ZoneId`
- `entity` (Optional, ForceNew, TypeString): 安全策略类型，合法取值 `ZoneDefaultPolicy`、`Template`、`Host`，对应 `DescribeSecurityPolicyRequest.Entity` 与 `ModifySecurityPolicyRequest.Entity`
- `host` (Optional, ForceNew, TypeString): 域名，当 `entity` 为 `Host` 时必填，对应 `DescribeSecurityPolicyRequest.Host` 与 `ModifySecurityPolicyRequest.Host`
- `template_id` (Optional, ForceNew, TypeString): 策略模板 ID，当 `entity` 为 `Template` 时必填，对应 `DescribeSecurityPolicyRequest.TemplateId` 与 `ModifySecurityPolicyRequest.TemplateId`
- `security_config` (Optional, TypeList, MaxItems 1): 安全策略配置，对应 `ModifySecurityPolicyRequest.SecurityConfig`

#### Scenario: Schema defines top-level query parameters
- **WHEN** the resource schema is defined
- **THEN** it SHALL include `zone_id`, `entity`, `host`, `template_id`, and `security_config` fields with correct types and constraints

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** `zone_id`, `entity`, `host`, or `template_id` is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Entity validation
- **WHEN** `entity` is set to a value other than `ZoneDefaultPolicy`, `Template`, or `Host`
- **THEN** the schema validation SHALL reject the configuration

### Requirement: Composite ID Format
The system SHALL construct a composite resource ID using `tccommon.FILED_SP` as the separator, with three legal forms determined by the `entity` value:
- `ZoneDefaultPolicy`: ID = `zoneId#ZoneDefaultPolicy` (2 segments), `host` and `template_id` MUST NOT be set
- `Host`: ID = `zoneId#Host#host` (3 segments), only `host` MUST be set
- `Template`: ID = `zoneId#Template#templateId` (3 segments), only `template_id` MUST be set

#### Scenario: ZoneDefaultPolicy ID construction
- **WHEN** `entity` is `ZoneDefaultPolicy` and neither `host` nor `template_id` is set
- **THEN** the resource ID SHALL be `zoneId#ZoneDefaultPolicy`

#### Scenario: Host ID construction
- **WHEN** `entity` is `Host` and `host` is set (and `template_id` is not set)
- **THEN** the resource ID SHALL be `zoneId#Host#host`

#### Scenario: Template ID construction
- **WHEN** `entity` is `Template` and `template_id` is set (and `host` is not set)
- **THEN** the resource ID SHALL be `zoneId#Template#templateId`

#### Scenario: Illegal entity combination rejected
- **WHEN** `entity` is `ZoneDefaultPolicy` but `host` or `template_id` is set, OR `entity` is `Host` but `host` is empty or `template_id` is set, OR `entity` is `Template` but `template_id` is empty or `host` is set
- **THEN** Create SHALL return an error describing the legal combinations

### Requirement: RateLimitConfig Schema
The system SHALL define a `security_config.rate_limit_config` block (Optional, Computed, TypeList, MaxItems 1) corresponding to `SecurityConfig.RateLimitConfig`, containing:
- `switch` (Optional, Computed, TypeString): 速率限制总开关，取值 `on`/`off`，对应 `RateLimitConfig.Switch`
- `rate_limit_user_rules` (Optional, Computed, TypeList): 用户自定义速率限制规则列表，对应 `RateLimitConfig.RateLimitUserRules`（`[]*RateLimitUserRule`）
- `rate_limit_template` (Optional, Computed, TypeList, MaxItems 1): 速率限制模板，对应 `RateLimitConfig.RateLimitTemplate`
- `rate_limit_intelligence` (Optional, Computed, TypeList, MaxItems 1): 智能客户端过滤，对应 `RateLimitConfig.RateLimitIntelligence`
- `rate_limit_customizes` (Optional, Computed, TypeList): 托管定制规则列表，对应 `RateLimitConfig.RateLimitCustomizes`（`[]*RateLimitUserRule`），字段结构与 `rate_limit_user_rules` 一致

#### Scenario: RateLimitConfig fields defined
- **WHEN** the `rate_limit_config` schema is defined
- **THEN** it SHALL include `switch`, `rate_limit_user_rules`, `rate_limit_template`, `rate_limit_intelligence`, and `rate_limit_customizes` fields

#### Scenario: rate_limit_customizes reuses rate_limit_user_rules structure
- **WHEN** the `rate_limit_customizes` element schema is defined
- **THEN** it SHALL share the same field set as `rate_limit_user_rules` (both map to SDK `RateLimitUserRule`)

### Requirement: RateLimitUserRule Schema
The system SHALL define the element schema for `rate_limit_user_rules` and `rate_limit_customizes` (both map to SDK `RateLimitUserRule`) with the following fields:
- `threshold` (Optional, Computed, TypeInt): 速率限制统计阈值，范围 0-4294967294，对应 `RateLimitUserRule.Threshold`
- `period` (Optional, Computed, TypeInt): 速率限制统计周期，取值 10/20/30/40/50/60 秒，对应 `RateLimitUserRule.Period`
- `rule_name` (Optional, Computed, TypeString): 规则名，对应 `RateLimitUserRule.RuleName`
- `action` (Optional, Computed, TypeString): 处置动作，取值 `monitor`/`drop`/`redirect`/`page`/`alg`，对应 `RateLimitUserRule.Action`
- `punish_time` (Optional, Computed, TypeInt): 惩罚时长，0-2 天，对应 `RateLimitUserRule.PunishTime`
- `punish_time_unit` (Optional, Computed, TypeString): 惩罚时长单位，取值 `second`/`minutes`/`hour`，对应 `RateLimitUserRule.PunishTimeUnit`
- `rule_status` (Optional, Computed, TypeString): 规则状态，取值 `on`/`off`，对应 `RateLimitUserRule.RuleStatus`
- `acl_conditions` (Optional, Computed, TypeList): 规则详情条件列表，对应 `RateLimitUserRule.AclConditions`（`[]*AclCondition`）
- `rule_priority` (Optional, Computed, TypeInt): 规则权重，范围 0-100，对应 `RateLimitUserRule.RulePriority`
- `rule_id` (Computed, TypeInt): 规则 ID，仅出参，对应 `RateLimitUserRule.RuleID`
- `freq_fields` (Optional, Computed, TypeList of TypeString): 过滤词，取值 `sip`，对应 `RateLimitUserRule.FreqFields`
- `update_time` (Computed, TypeString): 更新时间，仅出参，对应 `RateLimitUserRule.UpdateTime`
- `freq_scope` (Optional, Computed, TypeList of TypeString): 统计范围，取值 `source_to_eo`/`client_to_eo`，对应 `RateLimitUserRule.FreqScope`
- `name` (Optional, Computed, TypeString): 自定义返回页面名称，`action` 为 `page` 时必填，对应 `RateLimitUserRule.Name`
- `custom_response_id` (Optional, Computed, TypeString): 自定义响应 ID，对应 `RateLimitUserRule.CustomResponseId`
- `response_code` (Optional, Computed, TypeInt): 自定义返回页面响应码，100-600（不含 3xx），`action` 为 `page` 时必填，对应 `RateLimitUserRule.ResponseCode`
- `redirect_url` (Optional, Computed, TypeString): 重定向地址，`action` 为 `redirect` 时必填，对应 `RateLimitUserRule.RedirectUrl`

#### Scenario: RateLimitUserRule fields defined
- **WHEN** the `rate_limit_user_rules` element schema is defined
- **THEN** it SHALL include all RateLimitUserRule fields with correct types and Computed/Optional attributes

#### Scenario: Computed output-only fields
- **WHEN** the schema defines `rule_id` and `update_time`
- **THEN** these fields SHALL be Computed-only (no Optional) since they are output-only API fields

### Requirement: AclCondition Schema
The system SHALL define the element schema for `acl_conditions` (maps to SDK `AclCondition`) with the following fields:
- `match_from` (Optional, Computed, TypeString): 匹配字段，对应 `AclCondition.MatchFrom`
- `match_param` (Optional, Computed, TypeString): 匹配参数，当 `match_from` 为 `header` 时填入 header key，对应 `AclCondition.MatchParam`
- `operator` (Optional, Computed, TypeString): 匹配关系，对应 `AclCondition.Operator`
- `match_content` (Optional, Computed, TypeString): 匹配内容，对应 `AclCondition.MatchContent`

#### Scenario: AclCondition fields defined
- **WHEN** the `acl_conditions` element schema is defined
- **THEN** it SHALL include `match_from`, `match_param`, `operator`, and `match_content` fields

### Requirement: RateLimitTemplate Schema
The system SHALL define the `rate_limit_template` block (maps to SDK `RateLimitTemplate`) with the following fields:
- `mode` (Optional, Computed, TypeString): 模板等级，取值 `sup_loose`/`loose`/`emergency`/`normal`/`strict`/`close`，对应 `RateLimitTemplate.Mode`
- `action` (Optional, Computed, TypeString): 模板处置方式，取值 `alg`/`monitor`，对应 `RateLimitTemplate.Action`
- `rate_limit_template_detail` (Computed, TypeList, MaxItems 1): 模板值详情，仅出参，对应 `RateLimitTemplate.RateLimitTemplateDetail`

#### Scenario: RateLimitTemplate fields defined
- **WHEN** the `rate_limit_template` schema is defined
- **THEN** it SHALL include `mode`, `action`, and `rate_limit_template_detail` fields

### Requirement: RateLimitTemplateDetail Schema
The system SHALL define the `rate_limit_template_detail` block (Computed-only, maps to SDK `RateLimitTemplateDetail`) with the following fields:
- `mode` (Computed, TypeString): 模板等级，对应 `RateLimitTemplateDetail.Mode`
- `id` (Computed, TypeInt): 唯一 ID，对应 `RateLimitTemplateDetail.ID`
- `action` (Computed, TypeString): 模板处置方式，对应 `RateLimitTemplateDetail.Action`
- `punish_time` (Computed, TypeInt): 惩罚时间（秒），对应 `RateLimitTemplateDetail.PunishTime`
- `threshold` (Computed, TypeInt): 统计阈值，对应 `RateLimitTemplateDetail.Threshold`
- `period` (Computed, TypeInt): 统计周期（秒），对应 `RateLimitTemplateDetail.Period`

#### Scenario: RateLimitTemplateDetail is Computed-only
- **WHEN** the `rate_limit_template_detail` schema is defined
- **THEN** all its fields SHALL be Computed-only (output-only), reflecting that `RateLimitTemplateDetail` is returned solely by the API

### Requirement: RateLimitIntelligence Schema
The system SHALL define the `rate_limit_intelligence` block (maps to SDK `RateLimitIntelligence`) with the following fields:
- `switch` (Optional, Computed, TypeString): 功能开关，取值 `on`/`off`，对应 `RateLimitIntelligence.Switch`
- `action` (Optional, Computed, TypeString): 执行动作，取值 `monitor`/`alg`，对应 `RateLimitIntelligence.Action`
- `rule_id` (Computed, TypeInt): 规则 ID，仅出参，对应 `RateLimitIntelligence.RuleId`

#### Scenario: RateLimitIntelligence fields defined
- **WHEN** the `rate_limit_intelligence` schema is defined
- **THEN** it SHALL include `switch`, `action`, and `rule_id` fields, with `rule_id` being Computed-only

### Requirement: Resource Create Operation
The system SHALL implement the Create operation by validating the `entity`/`host`/`template_id` combination, constructing the composite ID, and then delegating to the Update operation (calling `ModifySecurityPolicy` with `ZoneId`, `Entity`, `Host`, `TemplateId`, and `SecurityConfig`). No independent creation API is invoked.

#### Scenario: Create delegates to Update
- **WHEN** a valid `entity`/`host`/`template_id` combination is provided
- **THEN** the system SHALL set the composite ID and invoke the Update operation to call `ModifySecurityPolicy`

#### Scenario: ModifySecurityPolicy API failure
- **WHEN** `ModifySecurityPolicy` returns a retryable error during Create
- **THEN** the system SHALL retry using `tccommon.WriteRetryTimeout` via `resource.Retry`
- **AND** return `tccommon.RetryError` on non-retryable errors

### Requirement: Resource Read Operation
The system SHALL implement the Read operation by parsing the composite ID into `zoneId`/`entity`/`host`/`templateId`, calling `DescribeSecurityPolicy` with these parameters, and then populating the schema fields from the response. If the response is empty, the system SHALL log the ID and clear the resource ID from state.

#### Scenario: Read parses composite ID and queries
- **WHEN** Read is invoked
- **THEN** the system SHALL split the ID by `tccommon.FILED_SP`, extract `zoneId`/`entity`/`host`/`templateId`, and call `DescribeSecurityPolicy` with these parameters

#### Scenario: Read handles missing resource
- **WHEN** `DescribeSecurityPolicy` returns an empty response
- **THEN** the system SHALL log `[WARN]` with the logId and resource ID, then call `d.SetId("")` to remove the resource from state

#### Scenario: Read populates top-level fields
- **WHEN** `DescribeSecurityPolicy` returns a valid response
- **THEN** the system SHALL set `zone_id`, `entity`, `host`, and `template_id` from the parsed ID, and populate `security_config` sub-structures from the response

### Requirement: Resource Update Operation
The system SHALL implement the Update operation by reading the `security_config` block from the Terraform configuration, constructing the `SecurityConfig` (including `RateLimitConfig` with `RateLimitUserRules`, `RateLimitTemplate`, `RateLimitIntelligence`, and `RateLimitCustomizes`), and calling `ModifySecurityPolicy` with `ZoneId`, `Entity`, `Host`, `TemplateId`, and `SecurityConfig`.

#### Scenario: Update constructs RateLimitConfig
- **WHEN** the `security_config.rate_limit_config` block is present in the configuration
- **THEN** the system SHALL build a `RateLimitConfig` with `Switch`, `RateLimitUserRules`, `RateLimitTemplate`, `RateLimitIntelligence`, and `RateLimitCustomizes` from the configuration and include it in the `ModifySecurityPolicy` request

#### Scenario: ModifySecurityPolicy API failure during Update
- **WHEN** `ModifySecurityPolicy` returns a retryable error during Update
- **THEN** the system SHALL retry using `tccommon.WriteRetryTimeout` via `resource.Retry`
- **AND** return `tccommon.RetryError` on non-retryable errors

#### Scenario: Update re-reads state
- **WHEN** the `ModifySecurityPolicy` call succeeds
- **THEN** the system SHALL invoke the Read operation to refresh the state

### Requirement: Resource Delete Operation
The system SHALL implement the Delete operation by clearing the `SecurityConfig` (setting it to an empty value) and calling `ModifySecurityPolicy` to neutralize the security policy configuration.

#### Scenario: Delete neutralizes security config
- **WHEN** Delete is invoked
- **THEN** the system SHALL call `ModifySecurityPolicy` with an empty/neutralized `SecurityConfig` to remove the configuration

#### Scenario: ModifySecurityPolicy API failure during Delete
- **WHEN** `ModifySecurityPolicy` returns a retryable error during Delete
- **THEN** the system SHALL retry using `tccommon.WriteRetryTimeout` via `resource.Retry`
- **AND** return `tccommon.RetryError` on non-retryable errors

### Requirement: Resource Import
The system SHALL support importing the resource via `schema.ImportStatePassthrough`, using the composite ID as the import argument.

#### Scenario: Import via composite ID
- **WHEN** a user runs `terraform import` with a composite ID (`zoneId#entity` or `zoneId#entity#host` or `zoneId#entity#templateId`)
- **THEN** the system SHALL parse the ID and populate the state via the Read operation

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_security_policy_config.md` with a one-line description mentioning TEO, example usage, and an import section describing the composite ID format.

#### Scenario: Documentation exists
- **WHEN** the resource is defined
- **THEN** a `.md` file SHALL exist with a one-line description mentioning TEO, example usage covering the supported `entity` combinations, and an import section showing the composite ID format

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_security_policy_config_test.go` covering Create, Read, Update, and Delete operations.

#### Scenario: Unit tests cover CRUD
- **WHEN** unit tests are run
- **THEN** test cases SHALL cover Create, Read, Update, and Delete operations for the `tencentcloud_teo_security_policy_config` resource
