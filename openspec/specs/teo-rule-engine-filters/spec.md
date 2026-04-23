### Requirement: filters input parameter for DescribeRules
The `tencentcloud_teo_rule_engine` resource SHALL accept an optional `filters` parameter of type list, where each element is an object with `name` (string, required) and `values` (set of strings, required). When provided, these filters SHALL be passed to the `DescribeRules` API request as the `Filters` field. When not provided, the existing default behavior (hardcoded `rule-id` filter) SHALL be preserved.

#### Scenario: User specifies custom filters
- **WHEN** user configures `filters` with `name = "rule-id"` and `values = ["rule-xxx"]` in the `tencentcloud_teo_rule_engine` resource
- **THEN** the `DescribeRules` API SHALL be called with the user-specified filters in the `Filters` request field

#### Scenario: User does not specify filters
- **WHEN** user does not configure `filters` in the `tencentcloud_teo_rule_engine` resource
- **THEN** the `DescribeRules` API SHALL be called with the default `rule-id` filter (existing behavior preserved)

### Requirement: rule_items computed output from DescribeRules response
The `tencentcloud_teo_rule_engine` resource SHALL expose a computed `rule_items` parameter that maps to the `RuleItems` field in the `DescribeRules` API response. Each `rule_items` element SHALL contain `rule_id` (string), `rule_name` (string), `status` (string), `rules` (list, same structure as the existing `rules` schema), `rule_priority` (integer), and `tags` (set of strings).

#### Scenario: DescribeRules returns rule items
- **WHEN** the `DescribeRules` API returns a non-empty `RuleItems` list in the response
- **THEN** the `rule_items` computed parameter SHALL be populated with the full list of rule items, mapping each `RuleItem` field to the corresponding Terraform schema field

#### Scenario: DescribeRules returns empty rule items
- **WHEN** the `DescribeRules` API returns an empty or nil `RuleItems` list in the response
- **THEN** the `rule_items` computed parameter SHALL be set to an empty list

### Requirement: Backward compatibility
Adding the `filters` and `rule_items` parameters SHALL NOT break any existing Terraform configurations. The `filters` parameter SHALL be Optional, and the `rule_items` parameter SHALL be Computed. Existing configurations that do not use these parameters SHALL continue to work identically to before the change.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing `tencentcloud_teo_rule_engine` configuration that does not include `filters` or reference `rule_items`
- **THEN** the resource SHALL behave identically to before the change, with the default `rule-id` filter used in the read operation
