## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO L7 acceleration rules. Currently, the resource exposes individual top-level schema fields (`status`, `rule_name`, `description`, `branches`) that are internally assembled into a `RuleEngineItem` struct for API calls. The `ModifyL7AccRule` API accepts `ZoneId` and `Rule` as input parameters. While `zone_id` is already a schema field, the `Rule` parameter does not have a corresponding schema field — the individual sub-fields are flattened at the top level.

The cloud API `ModifyL7AccRule` (package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) accepts:
- `ZoneId` (*string) — already mapped to `zone_id` in the schema
- `Rule` (*RuleEngineItem) — needs to be exposed as a `rule` schema parameter

The `RuleEngineItem` struct contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, `RulePriority`.

## Goals / Non-Goals

**Goals:**
- Add a `rule` parameter (TypeList, MaxItems: 1, Optional) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema
- The `rule` parameter maps to `request.Rule` in the `ModifyL7AccRule` API
- The `rule` block contains sub-fields: `status`, `rule_name`, `description`, `branches`
- Update Create/Read/Update functions to support the new `rule` parameter
- Maintain backward compatibility with existing top-level fields
- Add unit tests for the new parameter

**Non-Goals:**
- Removing or deprecating existing top-level schema fields (`status`, `rule_name`, `description`, `branches`)
- Changing the `ModifyL7AccRule` API behavior
- Modifying any other resources or data sources

## Decisions

### Decision 1: `rule` parameter as TypeList with MaxItems: 1

The `rule` parameter will be a `TypeList` with `MaxItems: 1` containing a nested `schema.Resource` with sub-fields matching `RuleEngineItem`. This follows the existing pattern in the codebase for nested configuration blocks (e.g., `branches`).

**Alternative considered**: TypeMap — rejected because the sub-fields have different types (string, list of strings, nested list) which cannot be represented in a TypeMap.

### Decision 2: Coexistence with top-level fields

The `rule` parameter and the existing top-level fields (`status`, `rule_name`, `description`, `branches`) will coexist. When the `rule` block is specified, its values will be used for constructing the `RuleEngineItem`. When `rule` is not specified, the existing top-level fields will be used as before.

**Alternative considered**: Deprecate top-level fields — rejected because it would break backward compatibility.

### Decision 3: Read function populates both

The Read function will continue to populate the existing top-level fields from the API response. It will also populate the `rule` block if the user has configured it.

## Risks / Trade-offs

- **Duplicate configuration paths**: Users can specify rule properties via either the `rule` block or the top-level fields. If both are specified, the `rule` block takes precedence. → Mitigation: Document the behavior clearly in the resource description and .md file.
- **State consistency**: When reading, both the top-level fields and the `rule` block will be populated. If a user only uses top-level fields, the `rule` block will remain empty, and vice versa. → Mitigation: This follows standard Terraform patterns for optional nested blocks.
