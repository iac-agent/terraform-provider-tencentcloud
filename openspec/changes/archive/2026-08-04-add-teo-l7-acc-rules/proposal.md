## Why

TEO (Tencent EdgeOne) is a key cloud product that requires Terraform management for its L7 acceleration rules. Currently, there is `tencentcloud_teo_l7_acc_rule_v2` which manages individual rules, but there is no resource that manages the complete set of L7 acceleration rules under a zone. Adding `tencentcloud_teo_l7_acc_rules` will provide a unified way to manage all rules as a single resource, enabling Terraform users to define the full rule set declaratively.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_l7_acc_rules` of type RESOURCE_KIND_GENERAL to manage the full set of L7 acceleration rules under a TEO zone
- Implement CRUD operations using the following TEO cloud APIs:
  - `CreateL7AccRules` for batch rule creation
  - `DescribeL7AccRules` for querying rules by zone with optional filters
  - `ModifyL7AccRule` for single rule modification (called iteratively for each changed rule)
  - `DeleteL7AccRules` for batch rule deletion
- Support resource import using the zone ID as the resource identifier
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## Capabilities

### New Capabilities
- `teo-l7-acc-rules`: Manage the complete set of TEO L7 acceleration rules under a zone, including creating, reading, updating (by comparing and modifying individual rules), and deleting rules, with support for import via zone ID.

### Modified Capabilities
<!-- None - this is a new resource, no existing specs are modified -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go` - main resource implementation
- **New file**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go` - unit tests using gomonkey mock
- **New file**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md` - resource documentation template
- **Modified file**: `tencentcloud/provider.go` - register the new resource
- **Modified file**: `tencentcloud/provider.md` - add resource documentation entry
- **Dependencies**: Uses existing `teov20220901` SDK package from vendor, no new dependencies required
- **Reuses**: `TencentTeoL7RuleBranchBasicInfo` helper from existing `resource_tc_teo_l7_acc_rule_extension.go`