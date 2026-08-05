## Context

EdgeOne (TEO) has deprecated the legacy rule engine API and introduced new L7 acceleration rules APIs:
- `CreateL7AccRules` (batch create) / `DeleteL7AccRules` (batch delete)
- `DescribeL7AccRules` (query with pagination and filters)
- `ModifyL7AccRule` (update single rule by RuleId)
- `ModifyL7AccRulePriority` (reorder rule priority)

The existing `tencentcloud_teo_l7_acc_rule` (v1) uses the legacy API and is deprecated. `tencentcloud_teo_l7_acc_rule_v2` (v2) manages a single rule per resource instance. The new `tencentcloud_teo_l7_acc_rules` (plural) will manage multiple rules as a batch in a single Terraform resource, using the `rules` list block.

The `RuleEngineItem` struct is shared across all these APIs and contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches` ([]*RuleBranch), `RulePriority`.

## Goals / Non-Goals

**Goals:**
- Create a Terraform resource `tencentcloud_teo_l7_acc_rules` that manages multiple L7 acceleration rules in a single zone
- Use `CreateL7AccRules` for batch creation, `DeleteL7AccRules` for batch deletion
- Use `ModifyL7AccRule` for individual rule updates, `ModifyL7AccRulePriority` for priority reordering
- Use `DescribeL7AccRules` for reading state
- Follow the existing `tencentcloud_teo_l7_acc_rule_v2` patterns for schema structure (reuse `TencentTeoL7RuleBranchBasicInfo`)

**Non-Goals:**
- Do NOT modify existing `tencentcloud_teo_l7_acc_rule` or `tencentcloud_teo_l7_acc_rule_v2` resources
- Do NOT support the legacy rule engine API
- Do NOT change the `RuleEngineItem` or `RuleBranch` Go struct definitions (they come from the SDK)

## Decisions

### Decision 1: Resource ID format
Use `ZoneId#RuleId` compound ID with `tccommon.FILED_SP` separator, matching the v2 pattern. This allows the Read/Update/Delete functions to identify both the zone and the specific rule.

### Decision 2: Schema layout
Use a flat schema with `zone_id` at the top level and a `rules` TypeList block containing individual rule objects (`rule_name`, `status`, `description`, `branches`). Each rule's `rule_id` and `rule_priority` are computed. This matches the batch nature of the API.

### Decision 3: Create flow
Use `CreateL7AccRules` with the full `Rules` slice. After creation, extract the first `RuleId` from the response and set the resource ID. This allows creating multiple rules in one API call.

### Decision 4: Read flow
Use `DescribeL7AccRules` with `rule-id` filter to fetch the specific rule by its ID. This leverages the existing `DescribeTeoL7AccRuleById` service function.

### Decision 5: Update flow
Use `ModifyL7AccRule` for changes to rule content (status, rule_name, description, branches). Use `ModifyL7AccRulePriority` if the priority ordering needs to change (but since priorities are computed, this is handled externally).

### Decision 6: Delete flow
Use `DeleteL7AccRules` with the single `RuleId` to delete the rule.

### Decision 7: Reuse existing helper functions
Reuse `TencentTeoL7RuleBranchBasicInfo` from `resource_tc_teo_l7_acc_rule_extension.go` for branch schema, and `DescribeTeoL7AccRuleById` from `service_tencentcloud_teo.go` for reading.

### Decision 8: No ForceNew on zone_id change
`zone_id` is set as `ForceNew: true` since changing the zone requires recreating the resource.

## Risks / Trade-offs

- **Risk**: The batch API `CreateL7AccRules` supports creating multiple rules, but the Terraform resource manages one rule per instance. → **Mitigation**: Always use single-element slices for `Rules` in the API call.
- **Risk**: `ModifyL7AccRulePriority` requires the full ordered list of rule IDs for the zone. → **Mitigation**: This is a separate concern; the resource does not manage priority ordering directly (handled by `tencentcloud_teo_l7_acc_rule_priority_operation`).
- **Risk**: API may return empty results for newly created rules before they propagate. → **Mitigation**: Use `tccommon.ReadRetryTimeout` retry pattern in the Read function, as done in v2.