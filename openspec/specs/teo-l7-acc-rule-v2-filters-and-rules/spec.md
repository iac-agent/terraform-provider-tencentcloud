## ADDED Requirements

### Requirement: Resource supports filters input parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL accept an optional `filters` parameter of type list, where each element contains a `name` (string, required) and `values` (list of string, required) field. This parameter maps to the `Filters` field of the `DescribeL7AccRules` API request. When specified, the filters SHALL be passed to the `DescribeL7AccRules` API call during the Read operation.

#### Scenario: User specifies filters to query rules
- **WHEN** user configures `filters` with `name = "rule-id"` and `values = ["rule-abc123"]` in the resource
- **THEN** the Read function SHALL pass these filters to the `DescribeL7AccRules` API request

#### Scenario: User does not specify filters
- **WHEN** user does not configure `filters` in the resource
- **THEN** the Read function SHALL use the default behavior of filtering by `rule-id` extracted from the composite resource ID

#### Scenario: Multiple filters specified
- **WHEN** user configures multiple filter blocks
- **THEN** all filters SHALL be passed to the `DescribeL7AccRules` API request

### Requirement: Resource exposes rules computed output parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL expose a computed `rules` parameter of type list. Each element in the list SHALL contain `status` (string), `rule_id` (string), `rule_name` (string), `description` (list of string), `branches` (list of branch objects), and `rule_priority` (integer) fields. This parameter maps to the `Rules` field of the `DescribeL7AccRules` API response.

#### Scenario: Rules returned from API
- **WHEN** the `DescribeL7AccRules` API returns rules in its response
- **THEN** the `rules` computed parameter SHALL be populated with all rules from the API response, with each rule containing its `status`, `rule_id`, `rule_name`, `description`, `branches`, and `rule_priority` fields

#### Scenario: No rules returned from API
- **WHEN** the `DescribeL7AccRules` API returns an empty rules list
- **THEN** the `rules` computed parameter SHALL be set to an empty list

### Requirement: Backward compatibility maintained
Adding the `filters` and `rules` parameters SHALL NOT break existing Terraform configurations or state. The `filters` parameter SHALL be optional with no default value. The `rules` parameter SHALL be computed-only.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing Terraform configuration that does not include `filters` or `rules`
- **THEN** the resource SHALL behave identically to its current behavior with no errors or state changes
