## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently does not support the `filters` input parameter for the `DescribeL7AccRules` API, nor does it expose the `rules` computed output from the response. Adding these parameters enables users to filter L7 acceleration rules by custom criteria (e.g., rule-id) and access the full list of matched rules from the read operation.

## What Changes

- Add `filters` parameter (Optional, TypeList) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to `request.Filters` in the `DescribeL7AccRules` API. Each filter item contains `name` (string) and `values` (list of strings), corresponding to the `Filter` struct in the teo SDK.
- Add `rules` parameter (Computed, TypeList) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to `response.Response.Rules` in the `DescribeL7AccRules` API. Each rule item reflects the `RuleEngineItem` struct fields.
- Update the `ResourceTencentCloudTeoL7AccRuleV2Read` function to set `filters` in the `DescribeL7AccRules` request and flatten `rules` from the response.
- Update the `DescribeTeoL7AccRuleById` service method (or create a new method) to accept filters.
- Update unit tests to cover the new parameters.
- Update the `.md` documentation file with the new parameters.

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-filters`: Add `filters` input parameter and `rules` computed output parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource for filtering and exposing L7 acceleration rules.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` - Schema and CRUD function changes
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - Service layer method changes for DescribeL7AccRules with filters
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` - Unit test additions
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` - Documentation update
- Cloud API: `DescribeL7AccRules` (teo/v20220901) - New input parameters `Filters` and new output field `Rules`
