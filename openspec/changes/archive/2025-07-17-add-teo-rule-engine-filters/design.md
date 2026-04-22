## Context

The `tencentcloud_teo_rule_engine` resource manages EdgeOne (TEO) rule engine rules. Currently, when reading the resource state via the `DescribeRules` API, the `Filters` parameter is hardcoded in the service layer (`DescribeTeoRuleEngineById`) to filter only by `rule-id`. The API supports custom filter conditions, and the response returns `RuleItems` (a list of rule items), but neither is exposed to Terraform users.

The existing resource uses a composite ID format `zoneId#ruleId` and already has `zone_id` as a Required/ForceNew field. The `DescribeRules` API already accepts `ZoneId` and `Filters` as request parameters, and returns `ZoneId` and `RuleItems` in the response.

## Goals / Non-Goals

**Goals:**
- Add `filters` as an Optional parameter to the `tencentcloud_teo_rule_engine` resource schema, allowing users to specify custom filter conditions for the `DescribeRules` API call
- Add `rule_items` as a Computed parameter to expose the `RuleItems` list from the `DescribeRules` API response
- Maintain full backward compatibility — existing Terraform configurations must continue to work without changes
- Properly map the `Filter` type (`Name` + `Values`) and `RuleItem` type (with `RuleId`, `RuleName`, `Status`, `Rules`, `RulePriority`, `Tags`) to Terraform schema

**Non-Goals:**
- Changing the existing `zone_id` field behavior (it already exists and works correctly)
- Modifying the Create, Update, or Delete operations (only the Read operation is affected)
- Adding new CRUD API support — this change only adds parameters to the existing DescribeRules read path

## Decisions

1. **`filters` parameter as Optional with backward-compatible defaults**
   - The `filters` parameter will be Optional. When not specified by the user, the existing behavior (hardcoded `rule-id` filter) will be preserved in the service layer, ensuring backward compatibility.
   - When `filters` is specified by the user, those custom filters will be passed to the `DescribeRules` API instead of the hardcoded `rule-id` filter.
   - Alternative considered: Always require `filters` — rejected because it would break existing configurations.

2. **`filters` schema structure**
   - `filters` will be a `TypeList` of objects with `name` (TypeString, Required) and `values` (TypeSet of TypeString, Required) sub-fields, matching the cloud API `Filter` type structure.
   - The `DescribeRules` API documentation specifies that `Filters.Values` has an upper limit of 20, and the supported filter key is `rule-id`.

3. **`rule_items` parameter as Computed**
   - `rule_items` will be a Computed `TypeList` of objects, each containing `rule_id`, `rule_name`, `status`, `rules`, `rule_priority`, and `tags` — mirroring the `RuleItem` type from the cloud API.
   - The `rules` sub-field will reuse the same nested structure as the existing `rules` schema field.

4. **Service layer adjustment**
   - The `DescribeTeoRuleEngineById` service method currently constructs the filter internally. We will add a new service method or modify the existing one to accept custom filters, while keeping the default behavior when no custom filters are provided.

## Risks / Trade-offs

- **[Risk] Complex nested `rule_items` structure** → The `RuleItem` contains deeply nested `Rules` with `Conditions`, `Actions`, and `SubRules`. Flattening this into Terraform schema requires careful mapping. Mitigation: Reuse the existing `rules` schema structure for the `rules` sub-field within `rule_items`.
- **[Risk] Backward compatibility** → Adding `filters` as Optional ensures existing configurations are not broken. The `rule_items` as Computed also has no impact on existing configurations.
- **[Risk] API pagination** → The `DescribeRules` API may support pagination. Currently, the service handles a single page. Mitigation: Keep the existing pagination logic unchanged; the new parameters don't affect pagination behavior.
