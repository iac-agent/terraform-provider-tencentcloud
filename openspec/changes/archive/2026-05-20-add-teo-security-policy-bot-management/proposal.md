## Why

The `tencentcloud_teo_security_policy_config` resource's `security_policy` block is missing the `bot_management` parameter, which is available in the cloud API's `SecurityPolicy` struct (returned by `DescribeSecurityPolicy` and accepted by `ModifySecurityPolicy`). Bot management is a critical security feature for TEO (TencentCloud EdgeOne) that allows users to configure Bot detection, custom Bot rules, browser impersonation detection, and client attestation rules. Without this parameter, users cannot manage Bot protection through Terraform, limiting the completeness of infrastructure-as-code coverage for TEO security policies.

## What Changes

- Add a new `bot_management` parameter (TypeList, MaxItems: 1, Optional) to the `security_policy` block of `tencentcloud_teo_security_policy_config`
- The `bot_management` parameter will include the following sub-fields:
  - `enabled` (string, Required): Whether Bot management is enabled (on/off)
  - `custom_rules` (TypeList, MaxItems: 1, Optional): Bot custom rules configuration
    - `rules` (TypeList, Optional): List of Bot custom rules, each containing `id`, `name`, `enabled`, `priority`, `condition`, and `action` (using `SecurityWeightedAction` with weighted actions)
  - `basic_bot_settings` (TypeList, MaxItems: 1, Optional): Basic Bot settings affecting all domains
    - `source_idc` (TypeList, MaxItems: 1, Optional): IDC source IP configuration
    - `search_engine_bots` (TypeList, MaxItems: 1, Optional): Search engine bot configuration
    - `known_bot_categories` (TypeList, MaxItems: 1, Optional): Known Bot category configuration
    - `ip_reputation` (TypeList, MaxItems: 1, Optional): IP threat intelligence configuration
    - `bot_intelligence` (TypeList, MaxItems: 1, Optional): Bot intelligence analysis configuration
  - `client_attestation_rules` (TypeList, MaxItems: 1, Optional): Client attestation rules (in beta)
    - `rules` (TypeList, Optional): List of client attestation rules
  - `browser_impersonation_detection` (TypeList, MaxItems: 1, Optional): Browser impersonation detection rules
    - `rules` (TypeList, Optional): List of browser impersonation detection rules
- Add corresponding Read/Create/Update logic to handle the new `bot_management` parameter
- Add unit tests for the new parameter
- Update the resource markdown documentation

## Capabilities

### New Capabilities
- `teo-security-policy-bot-management`: Add Bot management configuration support to the `tencentcloud_teo_security_policy_config` resource's `security_policy` block, including custom rules, basic bot settings, client attestation rules, and browser impersonation detection.

### Modified Capabilities
<!-- No existing spec-level behavior changes -->

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` - Add schema definition, Read/Create/Update/Delete logic for `bot_management`
- **Tests**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` - Add unit tests
- **Docs**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` - Update markdown documentation
- **APIs**: Uses existing `DescribeSecurityPolicy` and `ModifySecurityPolicy` cloud API endpoints
- **Backward Compatibility**: Fully backward compatible - adding an Optional parameter does not affect existing configurations
