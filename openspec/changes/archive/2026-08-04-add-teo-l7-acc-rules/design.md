## Context

The TEO (Tencent EdgeOne) product already has `tencentcloud_teo_l7_acc_rule_v2` which manages individual L7 acceleration rules (one rule per resource instance). However, there is no resource that manages the complete set of L7 acceleration rules under a zone as a single declarative block. This new resource `tencentcloud_teo_l7_acc_rules` fills that gap by managing all rules under a zone at once.

The existing `resource_tc_teo_l7_acc_rule_v2.go` and `resource_tc_teo_l7_acc_rule_extension.go` provide pattern references for schema construction, cloud API calls, and the `RuleEngineItem`/`RuleBranch` data structures.

## Goals / Non-Goals

**Goals:**
- Implement a RESOURCE_KIND_GENERAL Terraform resource `tencentcloud_teo_l7_acc_rules` that manages the full set of L7 acceleration rules under a TEO zone
- Support CRUD: Create (batch), Read (list with filters), Update (diff-based: create/update/delete individual rules), Delete (batch)
- Support resource import using zone ID
- Use the existing `teov20220901` SDK (already in vendor) for all API calls
- Reuse existing `RuleEngineItem` and `RuleBranch` schema helpers from `resource_tc_teo_l7_acc_rule_extension.go`

**Non-Goals:**
- No modification to existing `tencentcloud_teo_l7_acc_rule_v2` resource
- No new API dependencies beyond what is already in vendor
- No support for rule priority reordering (handled by a separate resource `tencentcloud_teo_l7_acc_rule_priority_operation`)

## Decisions

### Decision 1: Resource ID is `zone_id`
- **Choice**: Use `zone_id` as the sole resource ID (not a composite ID)
- **Rationale**: This resource manages all rules under a zone. The zone is the top-level container. Using `zone_id` as the ID is cleaner than a composite ID since there is exactly one rule set per zone.
- **Alternative considered**: `zone_id#rule_count` or similar composite — rejected because it's unnecessarily complex and the zone uniquely identifies the managed scope.

### Decision 2: Schema Design
- **Choice**: `zone_id` (Required, ForceNew), `rules` (Required, TypeList of rule objects), `rule_ids` (Computed, TypeList), `filters` (Optional, TypeList for the Describe API)
- **Rationale**: `zone_id` is ForceNew because changing zones would require re-creating rules. `rules` is Required because the resource is meaningless without rules. `rule_ids` is Computed as it's derived from the Create response. `filters` is Optional to allow filtering during Read.
- The `rules` list uses the same `RuleEngineItem` structure as `resource_tc_teo_l7_acc_rule_v2` (status, rule_name, description, branches sub-fields), ensuring consistency.

### Decision 3: Create Strategy
- **Choice**: Call `CreateL7AccRules` with all rules in a single batch call
- **Rationale**: The API supports batch creation. The response returns `RuleIds` in the same order as the input rules, which is stored in `rule_ids` computed attribute.
- After creation, call Read to populate `rule_ids` and verify state.

### Decision 4: Read Strategy
- **Choice**: Call `DescribeL7AccRules` with `ZoneId` (required) and optional `Filters`, using pagination (Limit=1000, Offset=0)
- **Rationale**: The API returns all rules under the zone. If `filters` is configured, pass it to the request. Set `Limit` to the maximum allowed value (1000) to minimize pagination calls.
- If the response is empty (no rules found), log the event and do NOT call `d.SetId("")` — instead return a `NonRetryableError` to let the outer retry handle it. Only set `d.SetId("")` when the API clearly indicates the resource is gone (not just empty, which could be a transient issue).

### Decision 5: Update Strategy
- **Choice**: Diff-based approach: compare old `rules` (from state) with new `rules` (from config), then:
  - For rules that exist in both but differ: call `ModifyL7AccRule` with the updated rule (including its `RuleId`)
  - For rules that are new (not in old): call `CreateL7AccRules` to create them
  - For rules that are removed (not in new): call `DeleteL7AccRules` to delete them
- **Rationale**: The `ModifyL7AccRule` API only supports single-rule modification. Incremental diff-based updates are the standard approach for list-managing resources.
- For matching rules between old and new, use `rule_name` as the matching key (since `RuleId` is only available after creation). If `rule_name` is not unique, use a combination of fields.

### Decision 6: Delete Strategy
- **Choice**: Call `DeleteL7AccRules` with all `rule_ids` from the current state
- **Rationale**: The API supports batch deletion. Collect all `rule_ids` from the Read response and pass them to the Delete API.

### Decision 7: Import Strategy
- **Choice**: Import by `zone_id` only
- **Rationale**: Since the resource manages all rules under a zone, the zone ID is sufficient to identify the resource. During import, the Read function will fetch all rules and populate the state.

## Risks / Trade-offs

- **[Risk] Update complexity**: The diff-based update approach requires matching rules between old and new state. If `rule_name` is not unique, matching may be ambiguous. → **Mitigation**: Use `rule_name` as the primary matching key; document that rule names should be unique within a zone for this resource.
- **[Risk] Partial update failure**: If the update process fails midway (e.g., after creating some rules but before deleting others), the state may be inconsistent. → **Mitigation**: Terraform's standard error handling will preserve the old state; the user can re-run `terraform apply` to retry.
- **[Risk] Large rule sets**: With many rules, the diff and update process could be slow. → **Mitigation**: Use batch API calls where possible (Create/Delete); only Modify is single-rule. The Limit=1000 for Read covers most practical scenarios.
- **[Trade-off] Different from `l7_acc_rule_v2`**: Users now have two ways to manage rules (individual vs. full set). → **Mitigation**: Document clearly that these are different resources for different use cases; the `_v2` resource is for managing individual rules, while `_rules` is for managing the complete set.