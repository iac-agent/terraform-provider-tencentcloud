## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently does not expose the `filters` input parameter and `rules` output parameter from the `DescribeL7AccRules` API. Adding these parameters allows users to filter L7 acceleration rules by specific criteria (e.g., rule-id) and access the full rules list as a computed output, enabling more flexible resource querying and better visibility into rule configurations.

## What Changes

- Add `filters` as an optional input parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to the `DescribeL7AccRules` API's `Filters` field. This allows users to specify filter conditions when reading the resource.
- Add `rules` as a computed output parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to the `DescribeL7AccRules` API's `Rules` response field. This exposes the full list of L7 acceleration rules returned by the API.

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-filters-and-rules`: Add `filters` input parameter and `rules` computed output parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource, enabling rule filtering and full rules list visibility.

### Modified Capabilities

## Impact

- Resource file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` - Schema definition and Read function modifications
- Service file: `tencentcloud/services/teo/service_tencentcloud_teo.go` - Update `DescribeTeoL7AccRuleById` to support `filters` parameter
- Test file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` - Add test cases for new parameters
- Documentation: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` - Update with new parameter documentation
- Cloud API: `DescribeL7AccRules` (teo/v20220901) - Already supports `Filters` input and `Rules` output in the SDK
