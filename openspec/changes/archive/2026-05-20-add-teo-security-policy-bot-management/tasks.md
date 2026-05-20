## 1. Schema Definition

- [x] 1.1 Add `bot_management` field (TypeList, MaxItems: 1, Optional) to the `security_policy` block in the resource schema, with sub-fields: `enabled` (Required, string), `custom_rules` (Optional, MaxItems: 1), `basic_bot_settings` (Optional, MaxItems: 1), `client_attestation_rules` (Optional, MaxItems: 1), `browser_impersonation_detection` (Optional, MaxItems: 1)
- [x] 1.2 Add `custom_rules` sub-schema with `rules` list, each rule containing: `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `priority` (Optional, int), `condition` (Required, string), `action` (Required, list of weighted actions with `security_action` block and `weight` field)
- [x] 1.3 Add SecurityAction shared schema block for use within bot_management, including: `name` (Required, string), `deny_action_parameters` (Optional, MaxItems: 1, with `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, `stall`), `redirect_action_parameters` (Optional, MaxItems: 1, with `url`), `allow_action_parameters` (Optional, MaxItems: 1, with `min_delay_time`, `max_delay_time`), `challenge_action_parameters` (Optional, MaxItems: 1, with `challenge_option`, `interval`, `attester_id`)
- [x] 1.4 Add `basic_bot_settings` sub-schema with: `source_idc` (Optional, MaxItems: 1), `search_engine_bots` (Optional, MaxItems: 1), `known_bot_categories` (Optional, MaxItems: 1), `ip_reputation` (Optional, MaxItems: 1), `bot_intelligence` (Optional, MaxItems: 1)
- [x] 1.5 Add shared `base_action_and_overrides` schema pattern for `source_idc`, `search_engine_bots`, `known_bot_categories` sub-fields (each with `base_action` SecurityAction block and `action_overrides` list with `ids` and `action`)
- [x] 1.6 Add `ip_reputation` sub-schema with `enabled` (Optional, string) and `ip_reputation_group` (Optional, MaxItems: 1, with `base_action` and `action_overrides`)
- [x] 1.7 Add `bot_intelligence` sub-schema with `enabled` (Optional, string), `id` (Computed, string), `bot_ratings` (Optional, MaxItems: 1, with `high_risk_bot_requests_action`, `likely_bot_requests_action`, `verified_bot_requests_action`, `human_requests_action` - all SecurityAction blocks)
- [x] 1.8 Add `client_attestation_rules` sub-schema with `rules` list, each containing: `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `priority` (Optional, int), `condition` (Required, string), `attester_id` (Required, string), `device_profiles` (Optional, list), `invalid_attestation_action` (Optional, SecurityAction block)
- [x] 1.9 Add `device_profiles` sub-schema with: `client_type` (Required, string), `high_risk_min_score` (Optional, int), `high_risk_request_action` (Optional, SecurityAction block), `medium_risk_min_score` (Optional, int), `medium_risk_request_action` (Optional, SecurityAction block)
- [x] 1.10 Add `browser_impersonation_detection` sub-schema with `rules` list, each containing: `id` (Computed, string), `name` (Required, string), `enabled` (Required, string), `condition` (Required, string), `action` (Required, MaxItems: 1)
- [x] 1.11 Add `browser_impersonation_detection_action` sub-schema with `bot_session_validation` (Optional, MaxItems: 1) and `client_behavior_detection` (Optional, MaxItems: 1)
- [x] 1.12 Add `bot_session_validation` sub-schema with: `issue_new_bot_session_cookie` (Optional, string), `max_new_session_trigger_config` (Optional, MaxItems: 1, with `max_new_session_count_interval` and `max_new_session_count_threshold`), `session_expired_action` (Optional, SecurityAction), `session_invalid_action` (Optional, SecurityAction), `session_rate_control` (Optional, MaxItems: 1, with `enabled`, `high_rate_session_action`, `mid_rate_session_action`, `low_rate_session_action`)
- [x] 1.13 Add `client_behavior_detection` sub-schema with: `crypto_challenge_intensity` (Optional, string), `crypto_challenge_delay_before` (Optional, string), `max_challenge_count_interval` (Optional, string), `max_challenge_count_threshold` (Optional, int), `challenge_not_finished_action` (Optional, SecurityAction), `challenge_timeout_action` (Optional, SecurityAction), `bot_client_action` (Optional, SecurityAction)

## 2. Flatten Functions (Read)

- [x] 2.1 Implement `flattenBotManagement` function to convert `BotManagement` SDK struct to Terraform state map
- [x] 2.2 Implement `flattenBotManagementCustomRules` function to convert `BotManagementCustomRules` to state map
- [x] 2.3 Implement `flattenSecurityWeightedAction` function to convert `SecurityWeightedAction` list to state list
- [x] 2.4 Implement `flattenBasicBotSettings` function to convert `BasicBotSettings` to state map
- [x] 2.5 Implement `flattenBotManagementBaseActionAndOverrides` function for shared pattern in SourceIDC, SearchEngineBots, KnownBotCategories, IPReputationGroup
- [x] 2.6 Implement `flattenIPReputation` and `flattenBotIntelligence` functions
- [x] 2.7 Implement `flattenBotRatings` function
- [x] 2.8 Implement `flattenClientAttestationRules` function with DeviceProfiles and InvalidAttestationAction
- [x] 2.9 Implement `flattenBrowserImpersonationDetection` function with BotSessionValidation and ClientBehaviorDetection
- [x] 2.10 Implement `flattenSecurityAction` reusable function for all SecurityAction fields within bot_management (supporting DenyActionParameters, RedirectActionParameters, AllowActionParameters, ChallengeActionParameters)
- [x] 2.11 Add `bot_management` flatten call in the resource Read function, after existing security_policy flatten logic

## 3. Expand Functions (Create/Update)

- [x] 3.1 Implement `expandBotManagement` function to convert Terraform configuration to `BotManagement` SDK struct
- [x] 3.2 Implement `expandBotManagementCustomRules` function to convert custom rules config to SDK struct
- [x] 3.3 Implement `expandSecurityWeightedAction` function to convert weighted action list config to SDK struct
- [x] 3.4 Implement `expandBasicBotSettings` function to convert basic bot settings config to SDK struct
- [x] 3.5 Implement `expandBotManagementBaseActionAndOverrides` function for shared pattern
- [x] 3.6 Implement `expandIPReputation` and `expandBotIntelligence` functions
- [x] 3.7 Implement `expandBotRatings` function
- [x] 3.8 Implement `expandClientAttestationRules` function with DeviceProfiles and InvalidAttestationAction
- [x] 3.9 Implement `expandBrowserImpersonationDetection` function with BotSessionValidation and ClientBehaviorDetection
- [x] 3.10 Implement `expandSecurityAction` reusable function for all SecurityAction fields within bot_management
- [x] 3.11 Add `bot_management` expand call in the resource Create/Update function, setting `SecurityPolicy.BotManagement` in the ModifySecurityPolicy request

## 4. Unit Tests

- [x] 4.1 Add unit tests for `flattenBotManagement` function covering all sub-fields
- [x] 4.2 Add unit tests for `expandBotManagement` function covering all sub-fields
- [x] 4.3 Add unit tests for SecurityAction flatten/expand within bot_management context
- [x] 4.4 Add unit tests for SecurityWeightedAction flatten/expand
- [x] 4.5 Add unit tests for BasicBotSettings flatten/expand
- [x] 4.6 Add unit tests for BrowserImpersonationDetection flatten/expand (including BotSessionValidation and ClientBehaviorDetection)
- [x] 4.7 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` to add bot_management example usage in the security_policy block
