## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO L7 acceleration rules. It uses the `DescribeL7AccRules` API to read rule state, which supports a `Filters` parameter for filtering results. Currently, the resource hardcodes a `rule-id` filter in the service layer (`DescribeTeoL7AccRuleById`) and does not expose the `Filters` API parameter to users. Adding `filters` as an Optional schema parameter allows users to specify custom filter conditions when the resource performs its read operation, providing parity with the cloud API.

The cloud API's `Filter` struct contains two fields:
- `Name` (*string): The field to filter by (e.g., `rule-id`)
- `Values` ([]*string): The filter values (max 20 per filter)

## Goals / Non-Goals

**Goals:**
- Add a `filters` parameter to the `tencentcloud_teo_l7_acc_rule_v2` resource schema
- Pass user-specified filters to the `DescribeL7AccRules` API request during the Read operation
- Maintain full backward compatibility — existing configurations without `filters` must work unchanged

**Non-Goals:**
- Changing the existing `rule-id` filter logic in the service layer (backward compatible behavior must be preserved)
- Adding `filters` support to Create, Update, or Delete operations (filters are only relevant for the Read/Describe operation)
- Changing the composite ID format or any other existing schema fields

## Decisions

1. **Schema definition for `filters`**: Use `TypeList` of `schema.Resource` with `name` (TypeString, Required) and `values` (TypeSet of TypeString, Required) sub-fields, matching the pattern used in other teo data sources (e.g., `data_source_tc_teo_plans`, `data_source_tc_teo_zones`).
   - Rationale: Consistency with existing patterns in the codebase.

2. **Read function behavior**: When `filters` is specified in the schema, pass the user-provided filters to the `DescribeL7AccRules` request IN ADDITION to the existing `rule-id` filter. This ensures the resource still reads its own state correctly while allowing users to specify additional filter conditions.
   - Alternative considered: Replace the hardcoded `rule-id` filter with user-provided filters — rejected because it could break the resource's ability to find its own rule.
   - Rationale: Backward compatibility and correct resource state management.

3. **Service layer update**: Update `DescribeTeoL7AccRuleById` to accept an additional `filters` parameter, or modify the Read function to construct the request with both the rule-id filter and user-specified filters.
   - Decision: Modify the Read function to append user-specified filters alongside the existing rule-id filter in the service layer call.

4. **`filters` is Optional and not Computed**: Since `filters` is an input-only parameter for the DescribeL7AccRules request and is not returned in the API response, it should be `Optional` only (not `Computed`). Users who don't specify it will get the default behavior.

## Risks / Trade-offs

- **[Risk] `filters` is unusual for a resource**: Filters are more commonly associated with data sources. Adding `filters` to a resource may confuse users. → Mitigation: Clear documentation explaining that `filters` is used to refine the query during resource read operations.
- **[Risk] Filter conflicts**: If a user specifies a filter that conflicts with the internal `rule-id` filter, the API behavior depends on the server-side logic (AND/OR). → Mitigation: Document that the `rule-id` filter is always applied, and user filters are appended.
