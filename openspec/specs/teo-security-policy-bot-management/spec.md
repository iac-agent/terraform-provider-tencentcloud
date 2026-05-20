## ADDED Requirements

### Requirement: Bot management parameter in security_policy block
The `security_policy` block of `tencentcloud_teo_security_policy_config` SHALL include an optional `bot_management` parameter (TypeList, MaxItems: 1) that maps to the `SecurityPolicy.BotManagement` field in the cloud API.

#### Scenario: Bot management parameter is optional
- **WHEN** a user creates or updates a `tencentcloud_teo_security_policy_config` resource without specifying `bot_management`
- **THEN** the resource SHALL be created/updated successfully without any BotManagement configuration changes
- **AND** existing BotManagement configuration on the cloud side SHALL remain unchanged (per API semantics: unspecified parameters are preserved)

#### Scenario: Bot management parameter is specified
- **WHEN** a user specifies the `bot_management` parameter in the `security_policy` block
- **THEN** the resource SHALL call `ModifySecurityPolicy` with the `SecurityPolicy.BotManagement` field populated
- **AND** the resource SHALL read back the BotManagement configuration via `DescribeSecurityPolicy`

### Requirement: Bot management enabled field
The `bot_management` block SHALL include an `enabled` field (TypeString, Required) with valid values `on` or `off`.

#### Scenario: Enable bot management
- **WHEN** a user sets `enabled` to `on`
- **THEN** the BotManagement feature SHALL be enabled for the security policy

#### Scenario: Disable bot management
- **WHEN** a user sets `enabled` to `off`
- **THEN** the BotManagement feature SHALL be disabled for the security policy

### Requirement: Bot management custom rules
The `bot_management` block SHALL include a `custom_rules` parameter (TypeList, MaxItems: 1, Optional) containing a `rules` list of Bot custom rule definitions.

#### Scenario: Configure bot custom rules
- **WHEN** a user specifies `custom_rules` with `rules`
- **THEN** each rule SHALL support `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `priority` (Optional, int), `condition` (Required, string), and `action` (Required, list of weighted actions)
- **AND** each action entry SHALL include `security_action` (a single action block with `name` and optional action parameters) and `weight` (Required, int, range 10-100, must be multiples of 10)
- **AND** the sum of all `weight` values in a rule's action list SHALL equal 100

#### Scenario: Bot custom rule IDs for add/modify/delete
- **WHEN** a user adds a new rule with `id` empty or unspecified
- **THEN** a new Bot custom rule SHALL be created
- **WHEN** a user specifies an existing `id`
- **THEN** the existing Bot custom rule SHALL be updated
- **WHEN** an existing rule ID is not included in the rules list
- **THEN** the existing Bot custom rule SHALL be deleted

### Requirement: Bot management basic bot settings
The `bot_management` block SHALL include a `basic_bot_settings` parameter (TypeList, MaxItems: 1, Optional) for configuring base-level Bot management that applies to all domains.

#### Scenario: Configure basic bot settings sub-fields
- **WHEN** a user specifies `basic_bot_settings`
- **THEN** it SHALL support the following sub-fields, each as TypeList with MaxItems: 1:
  - `source_idc`: IDC source IP configuration with `base_action` and `action_overrides`
  - `search_engine_bots`: Search engine bot configuration with `base_action` and `action_overrides`
  - `known_bot_categories`: Known bot category configuration with `base_action` and `action_overrides`
  - `ip_reputation`: IP reputation configuration with `enabled`, and `ip_reputation_group` containing `base_action` and `action_overrides`
  - `bot_intelligence`: Bot intelligence configuration with `enabled`, `id` (Computed), and `bot_ratings` containing action fields for `high_risk_bot_requests_action`, `likely_bot_requests_action`, `verified_bot_requests_action`, `human_requests_action`

#### Scenario: Action overrides for basic bot settings
- **WHEN** a user specifies `action_overrides` within a basic bot settings sub-field
- **THEN** each override SHALL include `ids` (list of string rule IDs) and `action` (a SecurityAction block)

### Requirement: Bot management client attestation rules
The `bot_management` block SHALL include a `client_attestation_rules` parameter (TypeList, MaxItems: 1, Optional) for configuring client attestation (beta feature).

#### Scenario: Configure client attestation rules
- **WHEN** a user specifies `client_attestation_rules` with `rules`
- **THEN** each rule SHALL support `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `priority` (Optional, int), `condition` (Required, string), `attester_id` (Required, string), `device_profiles` (Optional, list), and `invalid_attestation_action` (Optional, SecurityAction block)
- **AND** each `device_profiles` entry SHALL support `client_type` (Required, string), `high_risk_min_score` (Optional, int), `high_risk_request_action` (Optional, SecurityAction block), `medium_risk_min_score` (Optional, int), `medium_risk_request_action` (Optional, SecurityAction block)

### Requirement: Bot management browser impersonation detection
The `bot_management` block SHALL include a `browser_impersonation_detection` parameter (TypeList, MaxItems: 1, Optional) for configuring browser impersonation detection rules.

#### Scenario: Configure browser impersonation detection rules
- **WHEN** a user specifies `browser_impersonation_detection` with `rules`
- **THEN** each rule SHALL support `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `condition` (Required, string), and `action` (Required, block)
- **AND** the `action` block SHALL contain `bot_session_validation` (Optional) and `client_behavior_detection` (Optional)

#### Scenario: Bot session validation configuration
- **WHEN** a user specifies `bot_session_validation`
- **THEN** it SHALL support `issue_new_bot_session_cookie` (Optional, string), `max_new_session_trigger_config` (Optional, block with `max_new_session_count_interval` and `max_new_session_count_threshold`), `session_expired_action` (Optional, SecurityAction block), `session_invalid_action` (Optional, SecurityAction block), and `session_rate_control` (Optional, block with `enabled`, `high_rate_session_action`, `mid_rate_session_action`, `low_rate_session_action`)

#### Scenario: Client behavior detection configuration
- **WHEN** a user specifies `client_behavior_detection`
- **THEN** it SHALL support `crypto_challenge_intensity` (Optional, string), `crypto_challenge_delay_before` (Optional, string), `max_challenge_count_interval` (Optional, string), `max_challenge_count_threshold` (Optional, int), `challenge_not_finished_action` (Optional, SecurityAction block), `challenge_timeout_action` (Optional, SecurityAction block), `bot_client_action` (Optional, SecurityAction block)

### Requirement: SecurityAction schema for bot management
The SecurityAction blocks within `bot_management` SHALL follow a consistent schema pattern supporting the `name` field and applicable action parameters.

#### Scenario: Deny action with parameters
- **WHEN** a SecurityAction `name` is `Deny`
- **THEN** `deny_action_parameters` SHALL be available with `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, `stall`

#### Scenario: Redirect action with parameters
- **WHEN** a SecurityAction `name` is `Redirect`
- **THEN** `redirect_action_parameters` SHALL be available with `url`

#### Scenario: Allow action with parameters
- **WHEN** a SecurityAction `name` is `Allow`
- **THEN** `allow_action_parameters` SHALL be available with `min_delay_time` and `max_delay_time`

#### Scenario: Challenge action with parameters
- **WHEN** a SecurityAction `name` is `Challenge`
- **THEN** `challenge_action_parameters` SHALL be available with `challenge_option`, `interval`, `attester_id`

### Requirement: Read operation for bot management
The resource Read operation SHALL read the `BotManagement` field from `DescribeSecurityPolicy` response and flatten it into the `bot_management` attribute in the Terraform state.

#### Scenario: Read bot management from cloud API
- **WHEN** the `DescribeSecurityPolicy` response contains a non-nil `SecurityPolicy.BotManagement`
- **THEN** the resource SHALL populate the `bot_management` attribute in the Terraform state with all nested fields
- **WHEN** the `SecurityPolicy.BotManagement` is nil
- **THEN** the resource SHALL set `bot_management` to an empty list

### Requirement: Create/Update operation for bot management
The resource Create and Update operations SHALL include the `BotManagement` field in the `ModifySecurityPolicy` request when `bot_management` is specified in the configuration.

#### Scenario: Create with bot management
- **WHEN** a user creates the resource with `bot_management` specified
- **THEN** the `ModifySecurityPolicy` request SHALL include the expanded `BotManagement` configuration in the `SecurityPolicy` field

#### Scenario: Update bot management configuration
- **WHEN** a user updates the `bot_management` configuration
- **THEN** the `ModifySecurityPolicy` request SHALL include the updated `BotManagement` configuration in the `SecurityPolicy` field

### Requirement: Backward compatibility
The addition of the `bot_management` parameter SHALL be fully backward compatible with existing configurations.

#### Scenario: Existing configurations without bot_management
- **WHEN** an existing `tencentcloud_teo_security_policy_config` resource does not include `bot_management`
- **THEN** the resource SHALL continue to function without any changes to its behavior
- **AND** the existing schema fields SHALL remain unchanged
