## ADDED Requirements

### Requirement: Resource supports filters parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL accept an optional `filters` parameter of type `TypeList`, where each element is a resource with `name` (TypeString, Required) and `values` (TypeSet of TypeString, Required) fields. This parameter maps to the `Filters` field in the `DescribeL7AccRules` API request.

#### Scenario: User specifies filters in resource configuration
- **WHEN** a user configures the `filters` parameter in a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the resource SHALL pass the user-specified filters to the `DescribeL7AccRules` API request during the Read operation, alongside the internal `rule-id` filter

#### Scenario: User does not specify filters
- **WHEN** a user does not configure the `filters` parameter
- **THEN** the resource SHALL behave exactly as before, using only the internal `rule-id` filter to query the rule

### Requirement: Filters schema structure matches cloud API Filter struct
The `filters` parameter SHALL use a schema structure that maps directly to the cloud API's `Filter` struct: each filter item SHALL have a `name` field (string, required) corresponding to `Filter.Name`, and a `values` field (set of strings, required) corresponding to `Filter.Values`.

#### Scenario: Filters are correctly mapped to API request
- **WHEN** the Read function constructs the `DescribeL7AccRules` request
- **THEN** each filter item from the schema SHALL be converted to a `teov20220901.Filter` struct with `Name` and `Values` fields populated from the `name` and `values` schema fields respectively

### Requirement: Backward compatibility is preserved
Adding the `filters` parameter SHALL NOT break any existing Terraform configurations or state files for the `tencentcloud_teo_l7_acc_rule_v2` resource.

#### Scenario: Existing configuration without filters
- **WHEN** a user applies an existing configuration that does not include the `filters` parameter
- **THEN** the resource SHALL create, read, update, and delete exactly as it did before this change

### Requirement: Unit tests cover filters parameter
Unit tests SHALL be added to verify the correct handling of the `filters` parameter in the Read function, including cases with and without filters specified.

#### Scenario: Unit test for filters in Read operation
- **WHEN** the Read function is called with filters specified
- **THEN** the `DescribeL7AccRules` request SHALL include the user-specified filters appended to the internal `rule-id` filter

#### Scenario: Unit test for Read operation without filters
- **WHEN** the Read function is called without filters specified
- **THEN** the `DescribeL7AccRules` request SHALL include only the internal `rule-id` filter
