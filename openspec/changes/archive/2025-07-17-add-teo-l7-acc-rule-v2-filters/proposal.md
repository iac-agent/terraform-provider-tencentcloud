## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently uses hardcoded filter logic (filtering by `rule-id`) in the service layer when calling `DescribeL7AccRules`. Exposing the `filters` parameter in the Terraform schema allows users to customize filtering conditions when reading rules, providing more flexibility and consistency with the cloud API interface.

## What Changes

- Add a new `filters` parameter (TypeList, Optional) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, corresponding to the `Filters` field in the `DescribeL7AccRules` API request. Each filter item contains `name` and `values` fields, matching the cloud API's `Filter` struct.

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-filters`: Add `filters` parameter to the teo_l7_acc_rule_v2 resource to support custom filtering when querying L7 acceleration rules via the DescribeL7AccRules API.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`: Add `filters` schema definition and update Read function to pass filters from schema to the DescribeL7AccRules request
- `tencentcloud/services/teo/service_tencentcloud_teo.go`: Update `DescribeTeoL7AccRuleById` to accept and pass through user-specified filters
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`: Add unit tests for the new `filters` parameter
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`: Update documentation with `filters` parameter example
