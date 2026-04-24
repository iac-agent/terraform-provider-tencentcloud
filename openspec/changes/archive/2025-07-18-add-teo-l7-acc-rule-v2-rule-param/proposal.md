## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently maps individual schema fields (status, rule_name, description, branches) to the `Rule` parameter of the `ModifyL7AccRule` API, but does not expose the `rule` parameter as a first-class schema field. Adding a `rule` parameter to the Terraform schema allows users to manage the rule configuration as a structured nested block, providing a more intuitive and consistent interface that aligns with the cloud API structure.

## What Changes

- Add a new `rule` parameter (TypeList, MaxItems: 1) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to `request.Rule` in the `ModifyL7AccRule` API
- The `rule` parameter wraps the `RuleEngineItem` struct fields: `status`, `rule_name`, `description`, `branches`
- Update the Create, Read, Update, and Delete functions to support the new `rule` parameter
- Update the existing top-level fields (`status`, `rule_name`, `description`, `branches`) to remain for backward compatibility
- Add unit tests for the new parameter

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-rule-param`: Adds the `rule` nested block parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource, enabling structured rule configuration aligned with the ModifyL7AccRule API's `Rule` field

### Modified Capabilities

## Impact

- **Affected files**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`
- **API**: `ModifyL7AccRule` (teo/v20220901) - the `Rule` field is already supported by the API; this change exposes it in the Terraform schema
- **Backward compatibility**: Existing top-level fields (status, rule_name, description, branches) remain unchanged; the new `rule` parameter is Optional, so existing configurations continue to work
