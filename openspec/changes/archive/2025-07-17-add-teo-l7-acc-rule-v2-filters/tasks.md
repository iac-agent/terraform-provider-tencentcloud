## 1. Schema Definition

- [x] 1.1 Add `filters` parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema in `resource_tc_teo_l7_acc_rule_v2.go`. The parameter should be `TypeList`, `Optional`, with `Elem` being a `schema.Resource` containing `name` (TypeString, Required) and `values` (TypeSet of TypeString, Required) fields, matching the pattern used in other teo data sources (e.g., `data_source_tc_teo_plans`).

## 2. Service Layer Update

- [x] 2.1 Update the `DescribeTeoL7AccRuleById` function in `service_tencentcloud_teo.go` to accept an additional `filters` parameter of type `[]*teov20220901.Filter`, which will be appended to the existing `rule-id` filter when constructing the `DescribeL7AccRules` request.
- [x] 2.2 Ensure the service function preserves backward compatibility: when no additional filters are provided, behavior remains unchanged (only the `rule-id` filter is applied).

## 3. CRUD Function Update

- [x] 3.1 Update `ResourceTencentCloudTeoL7AccRuleV2Read` in `resource_tc_teo_l7_acc_rule_v2.go` to read the `filters` parameter from the schema and pass it to the updated `DescribeTeoL7AccRuleById` service function.
- [x] 3.2 Add logic to convert the `filters` schema data (TypeList of maps with `name` and `values`) to `[]*teov20220901.Filter` before passing to the service function.

## 4. Unit Tests

- [x] 4.1 Add unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` to verify the `filters` parameter is correctly converted and passed to the `DescribeL7AccRules` API request during the Read operation, using gomonkey mocks for the cloud API client.
- [x] 4.2 Add unit test to verify backward compatibility: when `filters` is not specified, the Read function still works correctly with only the internal `rule-id` filter.
- [x] 4.3 Run unit tests with `go test -gcflags=all=-l` to ensure all tests pass.

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to include the `filters` parameter description and example usage.
