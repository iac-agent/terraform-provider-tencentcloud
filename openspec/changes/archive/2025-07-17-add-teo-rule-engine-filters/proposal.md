## Why

The `tencentcloud_teo_rule_engine` resource currently hardcodes the `Filters` parameter in the `DescribeRules` API call (only filtering by `rule-id`). Users need the ability to specify custom filter conditions when querying rules, and also need access to the full `RuleItems` list returned by the `DescribeRules` API response. The `zone_id` is already available in the schema but should also be populated from the response for consistency.

## What Changes

- Add `filters` parameter (Optional, TypeList) to the `tencentcloud_teo_rule_engine` resource schema, mapping to the `Filters` input of the `DescribeRules` API. This allows users to specify custom filter conditions (e.g., `rule-id`) instead of relying on the hardcoded internal filter.
- Add `rule_items` parameter (Computed, TypeList) to the resource schema, mapping to the `RuleItems` output of the `DescribeRules` API response. This exposes the full list of rule items returned by the API.

## Capabilities

### New Capabilities
- `teo-rule-engine-filters`: Adds `filters` input parameter and `rule_items` computed output parameter to the `tencentcloud_teo_rule_engine` resource, enabling custom filter conditions and exposing the full rule items list from the DescribeRules API.

### Modified Capabilities

## Impact

- **Resource file**: `tencentcloud/services/teo/resource_tc_teo_rule_engine.go` — schema definition, read function, and service layer modifications
- **Service file**: `tencentcloud/services/teo/service_tencentcloud_teo.go` — `DescribeTeoRuleEngineById` may need adjustments to support custom filters
- **Test file**: `tencentcloud/services/teo/resource_tc_teo_rule_engine_test.go` — unit tests for new parameters
- **Doc file**: `tencentcloud/services/teo/resource_tc_teo_rule_engine.md` — documentation for new parameters
- **Cloud API**: `DescribeRules` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
