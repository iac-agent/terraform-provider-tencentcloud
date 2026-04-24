## 1. Service Layer

- [x] 1.1 Add `DescribeTeoL7AccRuleByFilters` method to `TeoService` in `tencentcloud/services/teo/service_tencentcloud_teo.go` that accepts `zoneId` and `filters` parameters, constructs a `DescribeL7AccRules` request with `ZoneId`, `Filters`, and `Limit` (set to 1000), and returns the `DescribeL7AccRulesResponseParams`

## 2. Schema Definition

- [x] 2.1 Add `filters` schema field to `ResourceTencentCloudTeoL7AccRuleV2` in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`: TypeList, Optional, with nested `name` (TypeString, Required) and `values` (TypeList of TypeString, Required) fields
- [x] 2.2 Add `rules` schema field to `ResourceTencentCloudTeoL7AccRuleV2`: TypeList, Computed, with nested fields `status` (TypeString, Computed), `rule_id` (TypeString, Computed), `rule_name` (TypeString, Computed), `description` (TypeList of TypeString, Computed), `branches` (TypeList, Computed, same nested schema as top-level `branches`), `rule_priority` (TypeInt, Computed)

## 3. CRUD Function Updates

- [x] 3.1 Update `ResourceTencentCloudTeoL7AccRuleV2Read` to read `filters` from schema, call `DescribeTeoL7AccRuleByFilters` when filters are specified, and flatten `rules` from the response into the computed `rules` field
- [x] 3.2 Add helper function `flattenTeoL7AccRules` to convert `[]*teov20220901.RuleEngineItem` to the Terraform `rules` schema format

## 4. Unit Tests

- [x] 4.1 Add unit tests for the new `filters` and `rules` parameters in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey mock approach
- [x] 4.2 Run `go test -gcflags=all=-l` on the test file to verify all tests pass

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` to include the new `filters` and `rules` parameters in the example usage
