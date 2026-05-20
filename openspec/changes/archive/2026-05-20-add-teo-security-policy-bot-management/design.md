## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies. The resource's `security_policy` block currently supports 5 configuration areas: `custom_rules`, `managed_rules`, `http_ddos_protection`, `rate_limiting_rules`, and `exception_rules`. However, the cloud API's `SecurityPolicy` struct (used by both `DescribeSecurityPolicy` and `ModifySecurityPolicy`) includes a `BotManagement` field that is not yet exposed in Terraform.

Bot Management is a critical security feature that allows configuring Bot detection rules, including:
- Custom Bot rules with weighted actions (e.g., 80% monitor + 20% deny)
- Basic Bot settings (IDC sources, search engine bots, known bot categories, IP reputation, bot intelligence)
- Client attestation rules (beta feature for device-based authentication)
- Browser impersonation detection (cookie validation, session tracking, client behavior detection)

The existing resource already handles complex nested structures like `http_ddos_protection` with `SecurityAction` patterns. The `BotManagement` parameter follows the same cloud API patterns.

## Goals / Non-Goals

**Goals:**
- Add `bot_management` parameter to the `security_policy` block of `tencentcloud_teo_security_policy_config`
- Support all sub-fields of `BotManagement`: `enabled`, `custom_rules`, `basic_bot_settings`, `client_attestation_rules`, `browser_impersonation_detection`
- Implement full CRUD support: the `bot_management` config is read via `DescribeSecurityPolicy` and written via `ModifySecurityPolicy`
- Maintain backward compatibility: adding an Optional parameter does not affect existing configurations
- Reuse existing `SecurityAction` schema patterns from `http_ddos_protection` where applicable

**Non-Goals:**
- Adding `BotManagementLite` parameter (separate future task)
- Adding `DefaultDenySecurityActionParameters` parameter (separate future task)
- Supporting deprecated `BlockIP` and `ReturnCustomPage` action names within the new `BotManagement` context (use `Deny` with `DenyActionParameters` instead)
- Modifying any existing schema fields or behavior

## Decisions

### 1. SecurityAction Schema Reuse within BotManagement

**Decision**: Define a shared `SecurityAction` schema block that can be reused across all action parameters within `bot_management`, similar to how `http_ddos_protection` already defines its action block.

**Rationale**: The `SecurityAction` struct is polymorphic and used in many places within BotManagement (custom rules, basic bot settings, browser impersonation detection, etc.). Reusing a consistent schema pattern reduces code duplication and maintenance burden. The action block should include:
- `name` (Required, string)
- `deny_action_parameters` (Optional, MaxItems: 1) - with `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, `stall`
- `redirect_action_parameters` (Optional, MaxItems: 1) - with `url`
- `allow_action_parameters` (Optional, MaxItems: 1) - with `min_delay_time`, `max_delay_time`
- `challenge_action_parameters` (Optional, MaxItems: 1) - with `challenge_option`, `interval`, `attester_id`

**Alternative considered**: Define separate action blocks for each BotManagement sub-component. Rejected because it leads to massive code duplication.

### 2. BotManagementCustomRule Action as WeightedAction List

**Decision**: Represent `BotManagementCustomRule.Action` (which is `[]*SecurityWeightedAction`) as a list of weighted action blocks, where each entry contains `security_action` (a single action block) and `weight` (integer).

**Rationale**: Unlike other action fields that use a single `SecurityAction`, BotManagement custom rules use `SecurityWeightedAction` which includes a weight. All weights must sum to 100. This requires a list representation rather than a single block.

### 3. BasicBotSettings Sub-Fields Structure

**Decision**: Each sub-field in `BasicBotSettings` (`source_idc`, `search_engine_bots`, `known_bot_categories`, `ip_reputation`, `bot_intelligence`) will be a TypeList with MaxItems: 1, containing:
- A `base_action` block (using the shared SecurityAction schema)
- An `action_overrides` list (for BotManagementActionOverrides)

For `ip_reputation`, add an `enabled` field and nest `ip_reputation_group` inside it.

For `bot_intelligence`, add an `enabled` field and `bot_ratings` block with action fields for each risk level.

**Rationale**: These sub-fields all follow a consistent pattern of `BaseAction` + `BotManagementActionOverrides` in the cloud API. Mapping them to similar Terraform schemas provides a consistent user experience.

### 4. Client Attestation Rules (Beta Feature)

**Decision**: Include `client_attestation_rules` in the schema but document it as a beta feature requiring a support ticket.

**Rationale**: The cloud API marks this as "内测中" (in beta), but it's still part of the API and should be available for users who have been granted access. Omitting it would make the provider incomplete.

### 5. BrowserImpersonationDetection Complex Action Structure

**Decision**: The `BrowserImpersonationDetectionRule.Action` uses `BrowserImpersonationDetectionAction` (not `SecurityAction`). This will be represented with two sub-blocks: `bot_session_validation` and `client_behavior_detection`.

**Rationale**: `BrowserImpersonationDetectionAction` is a unique struct with `BotSessionValidation` and `ClientBehaviorDetection` sub-fields. These don't follow the standard `SecurityAction` pattern and need their own schema definition.

### 6. Schema Organization

**Decision**: Add the `bot_management` schema definition inline within the existing `security_policy` block, and create helper functions for flattening/expanding the BotManagement data, similar to the existing pattern used for `http_ddos_protection`.

**Rationale**: The existing resource already has a very long schema definition. Adding inline schema keeps the code structure consistent with the rest of the resource. Helper functions for flatten/expand operations keep the Read/Create/Update functions manageable.

## Risks / Trade-offs

- **[Complexity]** → The BotManagement structure is deeply nested with many levels (up to 6-7 levels of nesting). This increases the complexity of the Terraform configuration but is unavoidable given the cloud API structure. Mitigation: Provide comprehensive documentation and examples.

- **[Schema Bloat]** → Adding BotManagement will significantly increase the size of an already large resource file. Mitigation: Use helper functions and organize code logically to maintain readability.

- **[Polymorphic SecurityAction]** → The `SecurityAction` struct is polymorphic with different parameter sets depending on the `name` value. Terraform schemas don't natively support discriminated unions. Mitigation: Make all action parameter blocks Optional and document which parameters apply for each action name.

- **[Backward Compatibility]** → Adding an Optional parameter is backward compatible by design. Existing configurations without `bot_management` will continue to work unchanged.

- **[Client Attestation Beta]** → `client_attestation_rules` is in beta. Mitigation: Document clearly that this feature requires a support ticket.
