## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO L7 acceleration rules via the `CreateL7AccRules`, `DescribeL7AccRules`, `ModifyL7AccRule`, and `DeleteL7AccRules` APIs. Currently, the Read operation calls `DescribeTeoL7AccRuleById` which only supports filtering by `rule-id` internally. The `DescribeL7AccRules` API also accepts a `Filters` parameter and returns a `Rules` list, but neither is exposed in the Terraform resource schema.

Current resource schema fields: `zone_id`, `status`, `rule_name`, `description`, `branches`, `rule_id` (computed), `rule_priority` (computed).

Cloud API `DescribeL7AccRules` request parameters: `ZoneId` (already used), `Filters` (not exposed), `Limit` (not exposed), `Offset` (not exposed).

Cloud API `DescribeL7AccRules` response parameters: `TotalCount`, `Rules` ([]*RuleEngineItem).

The `Filter` struct in the SDK contains: `Name` (string) and `Values` ([]*string).

The `RuleEngineItem` struct contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, `RulePriority`.

## Goals / Non-Goals

**Goals:**
- Add `filters` parameter to the resource schema to allow custom filtering when reading L7 acceleration rules
- Add `rules` computed parameter to expose the full list of matched rules from the `DescribeL7AccRules` response
- Maintain backward compatibility (both new fields are additive and optional/computed)
- Follow existing code patterns (reference `tencentcloud_igtm_strategy` for style)

**Non-Goals:**
- Do not modify existing schema fields or break existing Terraform configurations
- Do not add `Limit`/`Offset` pagination parameters (not required by this change)
- Do not modify the Create, Update, or Delete operations

## Decisions

### Decision 1: `filters` as Optional TypeList
The `filters` parameter will be defined as `Optional: true, TypeList` with nested `name` (string) and `values` (list of strings) fields, matching the `Filter` struct in the teo SDK. This allows users to specify custom filter criteria. It is optional because the existing behavior (filtering by rule-id internally) should remain the default.

### Decision 2: `rules` as Computed TypeList
The `rules` parameter will be defined as `Computed: true, TypeList` with nested fields matching `RuleEngineItem`: `status`, `rule_id`, `rule_name`, `description`, `branches`, `rule_priority`. This exposes the full list of matched rules from the read response. It is computed-only because it is populated from the API response.

### Decision 3: Service layer method update
The existing `DescribeTeoL7AccRuleById` method already handles the `Filters` field internally (for rule-id filtering). A new service method `DescribeTeoL7AccRuleByFilters` will be added to accept custom filters, while the existing method remains unchanged for backward compatibility. Alternatively, the existing method signature can be extended with an optional filters parameter. The preferred approach is to add a new method to avoid breaking existing callers.

### Decision 4: Read function update
In `ResourceTencentCloudTeoL7AccRuleV2Read`, after the existing read logic, if `filters` is set in the schema, the new service method will be called to fetch rules with custom filters, and the `rules` computed field will be populated. The existing single-rule read logic (using `zone_id` + `rule_id`) remains unchanged.

## Risks / Trade-offs

- [Risk] Adding `filters` changes the Read behavior when specified → Mitigation: `filters` is optional; when not specified, the existing behavior is preserved
- [Risk] The `rules` computed field could contain a large number of items → Mitigation: This is a computed field that only reflects the API response; users who don't need it can ignore it
- [Risk] The `filters` parameter is only meaningful for the Read operation, not for Create/Update/Delete → Mitigation: This is consistent with how filters work in data sources; the parameter is optional and does not affect other operations
