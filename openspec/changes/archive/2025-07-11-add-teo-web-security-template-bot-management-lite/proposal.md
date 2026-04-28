## Why

The TEO WebSecurityTemplate resource currently lacks support for the `BotManagementLite` parameter in the `security_policy` block. The cloud API has added this field (as `SecurityPolicy.BotManagementLite`) in the SDK, which provides basic Bot management configuration including CAPTCHA page challenge and AI crawler detection for security policy templates. This feature needs to be exposed in the Terraform resource so users can configure basic Bot management settings in their TEO web security templates.

## What Changes

- Add `bot_management_lite` field (TypeList, MaxItems: 1, Optional) to the `security_policy` block in `tencentcloud_teo_web_security_template` resource
- The `bot_management_lite` field contains two sub-fields:
  - `captcha_page_challenge`: CAPTCHA page challenge configuration (TypeList, MaxItems: 1, Optional)
    - `enabled`: Whether CAPTCHA page challenge is enabled, values: on/off (TypeString, Optional)
  - `ai_crawler_detection`: AI crawler detection configuration (TypeList, MaxItems: 1, Optional)
    - `enabled`: Whether AI crawler detection is enabled, values: on/off (TypeString, Optional)
    - `action`: Execution action for AI crawler detection (TypeList, MaxItems: 1, Optional) - reuses existing SecurityAction schema pattern
- Update Create method to include BotManagementLite in CreateWebSecurityTemplate request
- Update Read method to read BotManagementLite from DescribeWebSecurityTemplate response
- Update Update method to include BotManagementLite in ModifyWebSecurityTemplate request
- Update unit tests to cover the new parameter

## Capabilities

### New Capabilities
- `bot-management-lite`: Add BotManagementLite parameter support to tencentcloud_teo_web_security_template resource, enabling basic Bot management configuration (CAPTCHA page challenge and AI crawler detection) in security policy templates

### Modified Capabilities


## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_web_security_template.go` - Add schema definition, flatten/expand logic for bot_management_lite
- **Code**: `tencentcloud/services/teo/resource_tc_teo_web_security_template_test.go` - Add unit tests for the new parameter
- **Code**: `tencentcloud/services/teo/resource_tc_teo_web_security_template.md` - Update documentation with bot_management_lite example
- **API**: Uses existing `SecurityPolicy.BotManagementLite` field in CreateWebSecurityTemplate, ModifyWebSecurityTemplate, and DescribeWebSecurityTemplate APIs
- **SDK**: Uses `BotManagementLite`, `CAPTCHAPageChallenge`, `AICrawlerDetection` types from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- **Backward Compatibility**: New optional field - fully backward compatible
