## ADDED Requirements

### Requirement: bot_management_lite schema definition
The resource `tencentcloud_teo_web_security_template` SHALL include an optional `bot_management_lite` field within the `security_policy` block. The field SHALL be of type TypeList with MaxItems of 1.

The `bot_management_lite` block SHALL contain:
- `captcha_page_challenge` (TypeList, MaxItems: 1, Optional): CAPTCHA page challenge configuration
  - `enabled` (TypeString, Optional): Whether CAPTCHA page challenge is enabled, valid values: "on", "off"
- `ai_crawler_detection` (TypeList, MaxItems: 1, Optional): AI crawler detection configuration
  - `enabled` (TypeString, Optional): Whether AI crawler detection is enabled, valid values: "on", "off"
  - `action` (TypeList, MaxItems: 1, Optional): Execution action for AI crawler detection, following the SecurityAction schema pattern with sub-fields: name, deny_action_parameters, redirect_action_parameters, allow_action_parameters, challenge_action_parameters, block_ip_action_parameters, return_custom_page_action_parameters

#### Scenario: Creating a resource with bot_management_lite
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource with `security_policy.bot_management_lite` configured
- **THEN** the CreateWebSecurityTemplate API SHALL be called with `SecurityPolicy.BotManagementLite` populated including CAPTCHAPageChallenge and AICrawlerDetection settings

#### Scenario: Creating a resource without bot_management_lite
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource without `security_policy.bot_management_lite`
- **THEN** the resource SHALL be created successfully without error, and the BotManagementLite field SHALL NOT be sent in the API request

### Requirement: bot_management_lite read operation
The Read operation SHALL read the `BotManagementLite` field from the `DescribeWebSecurityTemplate` API response and flatten it into the Terraform state.

#### Scenario: Reading bot_management_lite from API response
- **WHEN** the DescribeWebSecurityTemplate API returns a non-nil `SecurityPolicy.BotManagementLite`
- **THEN** the `bot_management_lite` field SHALL be populated in the Terraform state with `captcha_page_challenge` and `ai_crawler_detection` sub-fields

#### Scenario: Reading when BotManagementLite is nil
- **WHEN** the DescribeWebSecurityTemplate API returns a nil `SecurityPolicy.BotManagementLite`
- **THEN** the `bot_management_lite` field SHALL NOT be set in the Terraform state (nil check before flattening)

### Requirement: bot_management_lite update operation
The Update operation SHALL include `BotManagementLite` in the `ModifyWebSecurityTemplate` API request when the `security_policy` field has changes.

#### Scenario: Updating bot_management_lite configuration
- **WHEN** a user updates the `bot_management_lite` field in an existing resource
- **THEN** the ModifyWebSecurityTemplate API SHALL be called with the updated `SecurityPolicy.BotManagementLite` value

#### Scenario: Removing bot_management_lite configuration
- **WHEN** a user removes the `bot_management_lite` field from an existing resource configuration
- **THEN** the ModifyWebSecurityTemplate API SHALL be called without BotManagementLite set (empty/nil)

### Requirement: bot_management_lite backward compatibility
The addition of `bot_management_lite` SHALL be fully backward compatible. Existing Terraform configurations without this field MUST continue to work without any changes.

#### Scenario: Existing configuration without bot_management_lite
- **WHEN** a user applies an existing Terraform configuration that does not include `bot_management_lite`
- **THEN** the resource SHALL continue to function identically to before the parameter was added
