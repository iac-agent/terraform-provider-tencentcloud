## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently maps individual fields (status, rule_name, description, branches) as top-level schema parameters, but the ModifyL7AccRule cloud API now accepts a `Rule` parameter that encapsulates the entire rule configuration as a single `RuleEngineItem` object. Adding a `rule` parameter to the Terraform resource allows users to pass the complete rule configuration in a single structured block, aligning with the cloud API's interface design and enabling more declarative resource management.

## What Changes

- Add a new `rule` parameter (TypeList, maxItems=1) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to the `Rule` field of the ModifyL7AccRule API's `RuleEngineItem` type.
- The `rule` parameter will contain nested fields matching the `RuleEngineItem` structure: `rule_id`, `status`, `rule_name`, `description`, and `branches`.
- Update the resource's Create, Read, and Update functions to handle the new `rule` parameter alongside the existing individual fields.
- Maintain backward compatibility: existing configurations using individual top-level fields (status, rule_name, description, branches) must continue to work.

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-rule-param`: Add `rule` parameter to the tencentcloud_teo_l7_acc_rule_v2 resource, mapping to the ModifyL7AccRule API's Rule field (RuleEngineItem type).

### Modified Capabilities

## Impact

- Affected files: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`
- Cloud API: ModifyL7AccRule (teo/v20220901)
- Backward compatibility: Existing configurations using individual top-level fields must remain functional
