## Why

The existing `tencentcloud_teo_l7_acc_rule` (v1) resource uses the legacy rule engine API which has been deprecated by EdgeOne. While `tencentcloud_teo_l7_acc_rule_v2` (v2) provides access to the new L7 acceleration rules API, it only manages a single rule at a time, making it inefficient for managing multiple rules in bulk. The new `tencentcloud_teo_l7_acc_rules` resource leverages the new `CreateL7AccRules`/`DeleteL7AccRules` batch APIs to manage multiple rules together, enabling Terraform users to efficiently manage entire rule sets for a zone.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_l7_acc_rules` for managing multiple L7 acceleration rules in batch
- Use the new Cloud API: `CreateL7AccRules` (batch create), `DescribeL7AccRules` (query with pagination), `ModifyL7AccRule` (single rule update), `ModifyL7AccRulePriority` (priority ordering), `DeleteL7AccRules` (batch delete)
- Resource ID uses `ZoneId#RuleId` compound key with `tccommon.FILED_SP` separator
- Schema supports nested `rules` block containing `rule_name`, `status`, `description`, `branches`, with computed `rule_id` and `rule_priority` fields
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## Capabilities

### New Capabilities
- `teo-l7-acc-rules-resource`: Terraform resource for batch management of EdgeOne L7 acceleration rules using the new Cloud API (CreateL7AccRules, DescribeL7AccRules, ModifyL7AccRule, ModifyL7AccRulePriority, DeleteL7AccRules)

### Modified Capabilities
<!-- None - this is a new resource that does not modify existing capabilities -->

## Impact

- New file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go` - resource implementation
- New file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go` - unit tests
- New file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md` - documentation
- Modified: `tencentcloud/provider.go` - resource registration
- Modified: `tencentcloud/provider.md` - resource documentation entry
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already vendored)