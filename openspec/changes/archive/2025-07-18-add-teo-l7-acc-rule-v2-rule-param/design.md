## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource currently maps the `ModifyL7AccRule` API's `Rule` parameter (of type `*RuleEngineItem`) by decomposing it into individual top-level Terraform schema fields: `status`, `rule_name`, `description`, and `branches`. The `ZoneId` field is also a top-level `zone_id` parameter.

The ModifyL7AccRule API expects two parameters:
- `ZoneId` (*string) - already mapped to `zone_id` (ForceNew, Required)
- `Rule` (*RuleEngineItem) - currently constructed from individual top-level fields but not exposed as a single schema parameter

The `RuleEngineItem` struct contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, `RulePriority`.

Current state: The update function creates a `RuleEngineItem` from the top-level schema fields and sets `request.Rule = rule`. There is no `rule` parameter in the Terraform schema.

## Goals / Non-Goals

**Goals:**
- Add a new `rule` parameter (TypeList, MaxItems=1, Optional) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema that maps to the `Rule` field of the ModifyL7AccRule API
- The `rule` parameter will contain nested fields matching the `RuleEngineItem` structure: `rule_id`, `status`, `rule_name`, `description`, and `branches`
- Maintain full backward compatibility with existing configurations using individual top-level fields
- Update Create, Read, and Update functions to handle the new `rule` parameter

**Non-Goals:**
- Removing or deprecating existing top-level fields (status, rule_name, description, branches) - these remain for backward compatibility
- Modifying the `zone_id` parameter - it already exists and is correctly mapped
- Changes to other TEO resources or data sources

## Decisions

1. **New `rule` parameter as TypeList with MaxItems=1**: The `Rule` field in ModifyL7AccRule is a single `*RuleEngineItem`, not a list. However, following Terraform convention for nested blocks in this codebase, we use TypeList with MaxItems=1 to represent a single object. This matches the pattern used for `branches` in the same resource.

2. **Coexistence with existing top-level fields**: The new `rule` parameter will coexist with the existing top-level fields (`status`, `rule_name`, `description`, `branches`). In the Update function, if the `rule` parameter is specified, it will be used to populate the `RuleEngineItem`. Otherwise, the existing individual fields will be used. This maintains backward compatibility.

3. **Nested schema for `rule`**: The `rule` parameter's nested schema will include `rule_id` (Computed), `status`, `rule_name`, `description`, and `branches` - matching the `RuleEngineItem` structure. The `branches` sub-field will reuse the existing `TencentTeoL7RuleBranchBasicInfo` helper function.

4. **Read function**: When reading the resource state, the `rule` parameter will be populated from the API response alongside the existing top-level fields.

## Risks / Trade-offs

- [Schema duplication] The same data can be represented both via top-level fields and via the `rule` block → Mitigation: In the Update function, the `rule` block takes precedence if specified; otherwise individual fields are used. Documentation should clarify the preferred approach.
- [State migration] Existing state files will not have the `rule` field populated → Mitigation: The `rule` field is Optional, so existing state files are valid. The Read function will populate it on the next refresh.
