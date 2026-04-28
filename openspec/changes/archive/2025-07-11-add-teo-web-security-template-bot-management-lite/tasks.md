## 1. Schema Definition

- [x] 1.1 Add `bot_management_lite` field (TypeList, MaxItems:1, Optional) to `security_policy` schema in `resource_tc_teo_web_security_template.go`, with sub-fields `captcha_page_challenge` and `ai_crawler_detection`
- [x] 1.2 Add `captcha_page_challenge` sub-field (TypeList, MaxItems:1, Optional) with `enabled` (TypeString, Optional) field
- [x] 1.3 Add `ai_crawler_detection` sub-field (TypeList, MaxItems:1, Optional) with `enabled` (TypeString, Optional) and `action` (TypeList, MaxItems:1, Optional) fields, where `action` follows existing SecurityAction schema pattern

## 2. Create Function

- [x] 2.1 Add BotManagementLite processing logic in `resourceTencentCloudTeoWebSecurityTemplateCreate` function, after the existing bot_management processing block, to populate `SecurityPolicy.BotManagementLite` from the Terraform configuration
- [x] 2.2 Add CAPTCHAPageChallenge processing: read `captcha_page_challenge` from botManagementLiteMap and set `BotManagementLite.CAPTCHAPageChallenge`
- [x] 2.3 Add AICrawlerDetection processing: read `ai_crawler_detection` from botManagementLiteMap and set `BotManagementLite.AICrawlerDetection` including Enabled and Action fields

## 3. Read Function

- [x] 3.1 Add BotManagementLite flattening logic in `resourceTencentCloudTeoWebSecurityTemplateRead` function, after the existing bot_management flattening block
- [x] 3.2 Add nil check for `SecurityPolicy.BotManagementLite` before flattening
- [x] 3.3 Flatten CAPTCHAPageChallenge: map `Enabled` field to `captcha_page_challenge.enabled`
- [x] 3.4 Flatten AICrawlerDetection: map `Enabled` field to `ai_crawler_detection.enabled` and flatten `Action` using existing SecurityAction flattening pattern
- [x] 3.5 Set flattened botManagementLiteMap into securityPolicyMap

## 4. Update Function

- [x] 4.1 Add BotManagementLite processing logic in `resourceTencentCloudTeoWebSecurityTemplateUpdate` function, after the existing bot_management processing block, to populate `SecurityPolicy.BotManagementLite` from the Terraform configuration for ModifyWebSecurityTemplate request
- [x] 4.2 Ensure the update logic mirrors the create logic for BotManagementLite field population

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_web_security_template.md` to add `bot_management_lite` example usage within `security_policy` block

## 6. Unit Tests

- [x] 6.1 Add unit test for BotManagementLite schema in `resource_tc_teo_web_security_template_test.go` covering create/read/update with bot_management_lite configuration
- [x] 6.2 Add unit test for BotManagementLite with captcha_page_challenge enabled
- [x] 6.3 Add unit test for BotManagementLite with ai_crawler_detection enabled and action configured
- [x] 6.4 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
