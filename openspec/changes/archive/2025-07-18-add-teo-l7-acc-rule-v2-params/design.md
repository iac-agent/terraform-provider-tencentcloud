## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO (TencentCloud EdgeOne) L7 acceleration rules. It currently supports CRUD operations via the `CreateL7AccRules`, `DescribeL7AccRules`, `ModifyL7AccRule`, and `DeleteL7AccRules` APIs.

The existing resource schema already includes `zone_id` (Required/ForceNew), `status`, `rule_name`, `description`, `branches` (all Optional), and `rule_id`, `rule_priority` (Computed). The Read function calls `DescribeTeoL7AccRuleById` which internally uses `Filters` with `rule-id` to fetch a specific rule, but `Filters` is not exposed as a user-configurable schema parameter. The `Rules` response is processed to extract individual fields but the full list is not exposed as a computed output.

The cloud API `DescribeL7AccRules` supports:
- **Input**: `ZoneId` (string), `Filters` ([]*Filter with Name/Values), `Limit` (int64), `Offset` (int64)
- **Output**: `TotalCount` (int64), `Rules` ([]*RuleEngineItem), `RequestId` (string)

The `Filter` type has `Name` (string) and `Values` ([]string), supporting filter key `rule-id`.

## Goals / Non-Goals

**Goals:**
- Add `filters` as an optional input parameter to allow users to specify filter conditions when reading the resource
- Add `rules` as a computed output parameter to expose the full list of L7 acceleration rules from the `DescribeL7AccRules` API response
- Maintain full backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Not changing the existing schema fields or their behavior
- Not adding `Limit`/`Offset` pagination parameters (these are internal implementation details)
- Not modifying the Create/Update/Delete operations

## Decisions

### Decision 1: `filters` schema design as TypeList with nested fields

The `filters` parameter will be defined as `TypeList` with `MaxItems: 20` (matching the API's `Filters.Values` upper limit of 20), containing nested resources with `name` (Required, TypeString) and `values` (Required, TypeList of TypeString) fields. This maps directly to the cloud API's `Filter` struct.

**Rationale**: This follows the existing pattern used in other TencentCloud Terraform resources (e.g., `tencentcloud_igtm_strategy`). The `TypeList` with nested resource approach provides clear structure and validation.

### Decision 2: `rules` schema design as computed TypeList

The `rules` parameter will be defined as `Computed: true, TypeList` containing nested resources that mirror the `RuleEngineItem` structure (with `status`, `rule_id`, `rule_name`, `description`, `branches`, `rule_priority` fields).

**Rationale**: Exposing the full rules list as a computed parameter allows users to access all rules returned by the API, which is particularly useful when `filters` are used to query specific rules. The nested structure reuses the existing `TencentTeoL7RuleBranchBasicInfo` helper for the `branches` field.

### Decision 3: Pass `filters` to the service Read method

The `DescribeTeoL7AccRuleById` service method will be updated to accept an optional `filters` parameter. When `filters` are provided in the Terraform configuration, they will be passed to the API call. When not provided, the existing behavior (filtering by `rule-id` from the composite ID) will be preserved.

**Rationale**: This maintains backward compatibility while enabling the new filtering capability. The service method already builds `Filters` internally, so extending it to accept external filters is a natural extension.

### Decision 4: Flattening `rules` response into schema

The Read function will flatten the `Rules` response from `DescribeL7AccRules` into the `rules` computed parameter. Each `RuleEngineItem` will be mapped to a flat map structure.

**Rationale**: Consistent with how other computed list parameters are handled in the provider.

## Risks / Trade-offs

- **[Risk] Adding `filters` to a resource may confuse users** → The `filters` parameter is primarily useful for data sources, but since the requirement explicitly asks for it in the resource, it will be documented clearly as an optional parameter that influences the Read query behavior.
- **[Risk] `rules` output may contain redundant data** → The `rules` computed parameter may overlap with existing individual fields (`status`, `rule_name`, etc.) when no filters are applied. This is acceptable since `rules` provides the full list while individual fields provide the single-rule convenience.
- **[Trade-off] Backward compatibility** → Both new parameters are additive (one Optional, one Computed), so existing configurations will continue to work without any changes.
