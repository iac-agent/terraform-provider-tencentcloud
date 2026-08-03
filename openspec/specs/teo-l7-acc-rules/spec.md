# teo-l7-acc-rules Specification

## Purpose
TBD - created by archiving change add-teo-l7-acc-rules. Update Purpose after archive.
## Requirements
### Requirement: Resource Schema Definition
The `tencentcloud_teo_l7_acc_rules` resource SHALL expose the following schema:
- `zone_id` (Required, ForceNew, TypeString): The zone ID
- `rules` (Optional, Computed, TypeList): List of rule objects, each containing:
  - `rule_id` (Computed, TypeString): Rule ID
  - `rule_name` (Optional, TypeString): Rule name
  - `status` (Optional, TypeString): Rule status (`enable` / `disable`)
  - `description` (Optional, TypeList of TypeString): Rule descriptions
  - `rule_priority` (Computed, TypeInt): Rule priority
  - `branches` (Optional, TypeList): Sub-rule branches
- `rule_ids` (Computed, TypeList of TypeString): List of rule IDs returned by Create API

#### Scenario: Schema defines all required and optional fields
- **WHEN** provider is initialized
- **THEN** the resource schema contains `zone_id`, `rules`, and `rule_ids` fields with correct type and attribute flags

### Requirement: Create Operation
The resource SHALL create L7 acceleration rules by calling `CreateL7AccRules` API.

#### Scenario: Successful creation of rules
- **WHEN** a Terraform configuration with `zone_id` and `rules` list is applied
- **THEN** the `CreateL7AccRules` API is called with the provided rules
- **AND** the resource ID is set to the `zone_id` value
- **AND** the `rule_ids` field is populated from the API response
- **AND** the resource state is refreshed by calling the Read operation

#### Scenario: Create with empty rules list
- **WHEN** a Terraform configuration with `zone_id` but empty `rules` list is applied
- **THEN** no API call is made
- **AND** the resource ID is set to the `zone_id` value

### Requirement: Read Operation
The resource SHALL read L7 acceleration rules by calling `DescribeL7AccRules` API with pagination support.

#### Scenario: Successful read of rules
- **WHEN** the Read operation is invoked with a valid `zone_id`
- **THEN** `DescribeL7AccRules` API is called with `Limit=1000`
- **AND** the `rules` field is populated with the returned rule list
- **AND** the `rule_ids` field is populated with all rule IDs

#### Scenario: Read when zone has no rules
- **WHEN** `DescribeL7AccRules` returns empty rules list
- **THEN** `rules` field is set to empty list
- **AND** the resource ID remains unchanged

#### Scenario: Read when zone does not exist
- **WHEN** the zone associated with the resource ID no longer exists
- **THEN** the resource ID is cleared (`d.SetId("")`)
- **AND** a warning log is emitted

### Requirement: Update Operation
The resource SHALL reconcile the `rules` field by comparing the desired state with the actual state, and applying the necessary API calls.

#### Scenario: Update when rules are modified
- **WHEN** the `rules` field in the Terraform configuration is changed
- **THEN** the diff between desired and actual rules is computed
- **AND** `ModifyL7AccRule` API is called for each modified rule that has a `rule_id`
- **AND** `CreateL7AccRules` API is called for new rules without `rule_id`
- **AND** `DeleteL7AccRules` API is called for removed rules
- **AND** the Read operation is called to refresh the state

#### Scenario: Update when only rule attributes change
- **WHEN** `rule_name`, `status`, `description`, or `branches` is changed on an existing rule
- **THEN** `ModifyL7AccRule` API is called with the updated `RuleEngineItem`

### Requirement: Delete Operation
The resource SHALL delete L7 acceleration rules by calling `DeleteL7AccRules` API.

#### Scenario: Successful deletion of all rules
- **WHEN** the Delete operation is invoked
- **THEN** all current `rule_ids` are collected from state
- **AND** `DeleteL7AccRules` API is called with all rule IDs
- **AND** the Read operation is called to confirm deletion

#### Scenario: Delete when resource has no rules
- **WHEN** the resource has no rules in state
- **THEN** the Delete operation succeeds without calling the API

### Requirement: Import Support
The resource SHALL support `terraform import` using `zone_id` as the import identifier.

#### Scenario: Import existing zone rules
- **WHEN** `terraform import tencentcloud_teo_l7_acc_rules.example zone-xxx` is executed
- **THEN** the resource is imported with the given `zone_id` as its ID
- **AND** the Read operation populates the `rules` field from the API

