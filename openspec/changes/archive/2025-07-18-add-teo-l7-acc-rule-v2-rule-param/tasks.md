## 1. Schema Definition

- [x] 1.1 Add `rule` parameter (TypeList, MaxItems: 1, Optional) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema in `resource_tc_teo_l7_acc_rule_v2.go`, with nested schema containing `status` (TypeString, Optional), `rule_name` (TypeString, Optional), `description` (TypeList of TypeString, Optional), and `branches` (TypeList using `TencentTeoL7RuleBranchBasicInfo`, Optional) sub-fields

## 2. CRUD Function Updates

- [x] 2.1 Update `ResourceTencentCloudTeoL7AccRuleV2Create` to construct `RuleEngineItem` from the `rule` block when it is specified, falling back to top-level fields when it is not
- [x] 2.2 Update `ResourceTencentCloudTeoL7AccRuleV2Read` to populate the `rule` block fields from the API response when the `rule` block is configured
- [x] 2.3 Update `ResourceTencentCloudTeoL7AccRuleV2Update` to detect changes in the `rule` block via `d.HasChange("rule")` and construct the `ModifyL7AccRuleRequest.Rule` from the `rule` block when changed, while preserving existing top-level field change detection
- [x] 2.4 Add helper functions to flatten `RuleEngineItem` into the `rule` block schema and to expand the `rule` block schema into `RuleEngineItem`

## 3. Documentation

- [x] 3.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to include the `rule` parameter in the example usage and parameter descriptions

## 4. Unit Tests

- [x] 4.1 Add unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` for the Create function with `rule` block using gomonkey mock approach
- [x] 4.2 Add unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` for the Update function with `rule` block changes using gomonkey mock approach
- [x] 4.3 Add unit tests for the helper functions (flatten/expand `rule` block) in `resource_tc_teo_l7_acc_rule_v2_test.go`
