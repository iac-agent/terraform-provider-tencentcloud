## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_l7_acc_rules` resource SHALL expose a Terraform schema with the following parameters:
- `zone_id` (Required, ForceNew, TypeString): The TEO zone ID under which rules are managed.
- `rules` (Required, TypeList): The list of L7 acceleration rules. Each element SHALL contain `rule_name`, `status`, `description`, and `branches` sub-fields. The `rule_id` and `rule_priority` sub-fields SHALL be Computed and set during Read.
- `filters` (Optional, TypeList): Filter conditions for the Describe API, each with `name` and `values` sub-fields.
- `rule_ids` (Computed, TypeList): The list of rule IDs returned by the Create API.

#### Scenario: Schema validation requires zone_id and rules
- **WHEN** a user defines a `tencentcloud_teo_l7_acc_rules` resource without `zone_id` or `rules`
- **THEN** Terraform SHALL report a validation error during plan

#### Scenario: zone_id is ForceNew
- **WHEN** a user changes the `zone_id` of an existing `tencentcloud_teo_l7_acc_rules` resource
- **THEN** Terraform SHALL plan to destroy and recreate the resource

### Requirement: Create L7 Acceleration Rules
The resource SHALL create L7 acceleration rules by calling the `CreateL7AccRules` API with the provided `zone_id` and `rules` list. On success, the resource SHALL populate `rule_ids` from the API response and call Read to refresh state.

#### Scenario: Successful batch creation
- **WHEN** a user applies a `tencentcloud_teo_l7_acc_rules` resource with valid `zone_id` and multiple `rules`
- **THEN** the system SHALL call `CreateL7AccRules` with all rules, set the resource ID to `zone_id`, and populate `rule_ids` from the response

#### Scenario: Create with empty rules list
- **WHEN** a user applies a `tencentcloud_teo_l7_acc_rules` resource with an empty `rules` list
- **THEN** the system SHALL return an error indicating that rules are required

#### Scenario: Create API returns empty response
- **WHEN** the `CreateL7AccRules` API call returns a response with nil Response or empty RuleIds
- **THEN** the system SHALL return a `NonRetryableError` to prevent writing an empty resource ID

### Requirement: Read L7 Acceleration Rules
The resource SHALL read L7 acceleration rules by calling the `DescribeL7AccRules` API with the `zone_id` and optional `filters`. The API SHALL use pagination with `Limit=1000` and `Offset=0`. The response rules SHALL be mapped to the `rules` attribute.

#### Scenario: Successful read
- **WHEN** the Read function is called for an existing resource with a valid `zone_id`
- **THEN** the system SHALL call `DescribeL7AccRules`, parse the response, and set `rules`, `rule_ids`, and `zone_id` in the state

#### Scenario: Rules not found
- **WHEN** the `DescribeL7AccRules` API returns an empty Rules list
- **THEN** the system SHALL log the event and return a `NonRetryableError` to trigger retry; it SHALL NOT call `d.SetId("")` to avoid data loss from transient API issues

#### Scenario: Read with filters
- **WHEN** the Read function is called with `filters` configured in the resource state
- **THEN** the system SHALL pass the filters to the `DescribeL7AccRules` API call

### Requirement: Update L7 Acceleration Rules
The resource SHALL update L7 acceleration rules by computing a diff between old and new `rules` lists. For rules that exist in both but differ, the system SHALL call `ModifyL7AccRule`. For new rules, the system SHALL call `CreateL7AccRules`. For removed rules, the system SHALL call `DeleteL7AccRules`.

#### Scenario: Rule modification
- **WHEN** a user modifies the `rule_name` or `status` of an existing rule in the `rules` list
- **THEN** the system SHALL call `ModifyL7AccRule` with the updated rule including its `RuleId`

#### Scenario: Rule addition
- **WHEN** a user adds a new rule to the `rules` list
- **THEN** the system SHALL call `CreateL7AccRules` to create the new rule

#### Scenario: Rule removal
- **WHEN** a user removes a rule from the `rules` list
- **THEN** the system SHALL call `DeleteL7AccRules` with the removed rule's ID

#### Scenario: No changes
- **WHEN** the `rules` list has not changed compared to the current state
- **THEN** the system SHALL skip all API calls and proceed to Read

### Requirement: Delete L7 Acceleration Rules
The resource SHALL delete all L7 acceleration rules by calling the `DeleteL7AccRules` API with all `rule_ids` from the current state.

#### Scenario: Successful batch deletion
- **WHEN** a user destroys a `tencentcloud_teo_l7_acc_rules` resource
- **THEN** the system SHALL call `DeleteL7AccRules` with all `rule_ids` from the state, then call `d.SetId("")`

#### Scenario: Delete with no rules
- **WHEN** a destroy is triggered but `rule_ids` is empty
- **THEN** the system SHALL call `d.SetId("")` directly without calling the API

### Requirement: Resource Import
The resource SHALL support import using the `zone_id` as the import identifier. During import, the Read function SHALL populate all rules under the zone.

#### Scenario: Import by zone_id
- **WHEN** a user runs `terraform import tencentcloud_teo_l7_acc_rules.example zone-xxx`
- **THEN** the system SHALL set the resource ID to `zone-xxx`, call `DescribeL7AccRules` to fetch all rules, and populate the state

### Requirement: Provider Registration
The resource SHALL be registered in `tencentcloud/provider.go` with the resource name `tencentcloud_teo_l7_acc_rules` and the factory function `ResourceTencentCloudTeoL7AccRules`.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the `tencentcloud_teo_l7_acc_rules` resource SHALL be available for use in Terraform configurations