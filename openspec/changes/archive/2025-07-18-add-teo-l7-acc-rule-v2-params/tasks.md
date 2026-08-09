## 1. Schema Definition

- [x] 1.1 Add `filters` optional input parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema in `resource_tc_teo_l7_acc_rule_v2.go`, with nested `name` (Required, TypeString) and `values` (Required, TypeList of TypeString) fields
- [x] 1.2 Add `rules` computed output parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema in `resource_tc_teo_l7_acc_rule_v2.go`, with nested fields: `status` (Computed, TypeString), `rule_id` (Computed, TypeString), `rule_name` (Computed, TypeString), `description` (Computed, TypeList of TypeString), `branches` (Computed, TypeList using TencentTeoL7RuleBranchBasicInfo), `rule_priority` (Computed, TypeInt)

## 2. Service Layer

- [x] 2.1 Update `DescribeTeoL7AccRuleById` method in `service_tencentcloud_teo.go` to accept optional `filters` parameter and pass it to the `DescribeL7AccRules` API request

## 3. CRUD Functions

- [x] 3.1 Update `ResourceTencentCloudTeoL7AccRuleV2Read` function to pass `filters` from schema to the service Read method
- [x] 3.2 Update `ResourceTencentCloudTeoL7AccRuleV2Read` function to flatten the `Rules` response into the `rules` computed schema parameter

## 4. Testing

- [x] 4.1 Add unit test cases for `filters` parameter in `resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey mock approach
- [x] 4.2 Add unit test cases for `rules` computed output parameter in `resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey mock approach

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to document the new `filters` and `rules` parameters with example usage
