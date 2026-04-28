## Context

The `tencentcloud_teo_web_security_template` resource manages TEO (TencentCloud EdgeOne) web security policy templates. The resource already supports `security_policy` with sub-fields: `custom_rules`, `managed_rules`, `http_ddos_protection`, `rate_limiting_rules`, `exception_rules`, and `bot_management`.

The TEO SDK (v1.3.88) has added a new field `BotManagementLite` to the `SecurityPolicy` struct, providing basic Bot management configuration for security policy templates. This is a lighter-weight alternative to the full `BotManagement` configuration, specifically designed for template usage with two key features: CAPTCHA page challenge and AI crawler detection.

## Goals / Non-Goals

**Goals:**
- Add `bot_management_lite` parameter to the `security_policy` block of `tencentcloud_teo_web_security_template` resource
- Support full CRUD for the new parameter (Create, Read, Update via existing APIs)
- Maintain backward compatibility - existing configurations without `bot_management_lite` must continue to work
- Follow existing code patterns in the resource file

**Non-Goals:**
- Modifying the existing `bot_management` parameter or any other existing parameters
- Adding `DefaultDenySecurityActionParameters` (another missing field, out of scope)
- Changing the cloud API behavior

## Decisions

1. **Schema design**: `bot_management_lite` will be a TypeList with MaxItems:1 and Optional, following the same pattern as other `security_policy` sub-fields (e.g., `bot_management`, `custom_rules`).

2. **Sub-fields structure**:
   - `captcha_page_challenge` (TypeList, MaxItems:1, Optional): Contains `enabled` (TypeString, Optional)
   - `ai_crawler_detection` (TypeList, MaxItems:1, Optional): Contains `enabled` (TypeString, Optional) and `action` (TypeList, MaxItems:1, Optional) using the existing SecurityAction schema pattern

3. **Action schema reuse**: The `action` field within `ai_crawler_detection` will follow the same SecurityAction pattern already used extensively in the resource (with deny_action_parameters, redirect_action_parameters, etc.), since the SDK type `AICrawlerDetection.Action` is of type `*SecurityAction`.

4. **Create/Update handling**: The `bot_management_lite` field will be processed alongside other security_policy sub-fields in both the Create and Update functions, setting `SecurityPolicy.BotManagementLite`.

5. **Read handling**: The Read function will flatten `SecurityPolicy.BotManagementLite` back to the Terraform state, checking for nil before processing (following existing nil-check patterns).

6. **No mutableArgs change needed**: `security_policy` is already in the `mutableArgs` list, so `bot_management_lite` changes will be detected as part of `security_policy` changes.

## Risks / Trade-offs

- [Nil response handling] → The BotManagementLite field may be nil in the API response when not configured. Must add nil checks before flattening, consistent with existing patterns in the resource.
- [Action schema complexity] → The SecurityAction type has many sub-fields (deny_action_parameters, redirect_action_parameters, etc.). For BotManagementLite, the AI crawler detection only supports a subset of actions (Deny, Monitor, Allow, Challenge). We will still define the full SecurityAction schema to stay consistent with the existing pattern, but the API will only accept certain action names.
