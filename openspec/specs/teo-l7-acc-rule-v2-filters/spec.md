## Requirements

### Requirement: filters parameter in teo_l7_acc_rule_v2 resource
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support a `filters` optional parameter of type TypeList, where each element contains `name` (String, Required) and `values` (TypeList of String, Required) sub-fields. This parameter maps to the `Filters` field of the `DescribeL7AccRules` API request.

#### Scenario: Resource created without filters
- **WHEN** a user creates a `tencentcloud_teo_l7_acc_rule_v2` resource without specifying the `filters` parameter
- **THEN** the Read operation SHALL use the default internal filtering logic (filtering by `rule-id` derived from the composite resource ID), maintaining backward compatibility

#### Scenario: Resource created with filters
- **WHEN** a user creates a `tencentcloud_teo_l7_acc_rule_v2` resource with the `filters` parameter specified
- **THEN** the Read operation SHALL pass the user-specified filters to the `DescribeL7AccRules` API's `Filters` field instead of the default `rule-id` filter

#### Scenario: filters schema structure
- **WHEN** a user specifies the `filters` parameter
- **THEN** each filter element SHALL have a `name` field (String, Required) and a `values` field (TypeList of String, Required)
- **AND** the `name` field SHALL correspond to the filter field name (e.g., `rule-id`)
- **AND** the `values` field SHALL contain the filter values

#### Scenario: filters only used in Read operation
- **WHEN** the resource performs Create, Update, or Delete operations
- **THEN** the `filters` parameter SHALL NOT be included in the respective API requests (CreateL7AccRules, ModifyL7AccRule, DeleteL7AccRules)
- **AND** only the `DescribeL7AccRules` API SHALL receive the filters parameter

#### Scenario: filters backward compatibility
- **WHEN** an existing Terraform configuration does not include the `filters` parameter
- **THEN** the resource SHALL continue to function exactly as before
- **AND** no changes to the existing state or configuration SHALL be required
