# teo-l7-acc-rules-resource Specification

## Purpose
TBD - created by archiving change add-teo-l7-acc-rules. Update Purpose after archive.
## Requirements
### Requirement: Resource schema definition
The `tencentcloud_teo_l7_acc_rules` resource SHALL define a schema with the following fields:
- `zone_id` (Required, ForceNew, TypeString): The zone ID
- `rules` (Optional, TypeList): List of rule configurations, each containing:
  - `rule_name` (Optional, TypeString): Rule name (max 255 characters)
  - `status` (Optional, TypeString): Rule status, values: `enable` or `disable`
  - `description` (Optional, TypeList of TypeString): Rule annotations
  - `branches` (Optional, TypeList): Sub-rule branches using `TencentTeoL7RuleBranchBasicInfo` schema
  - `rule_id` (Computed, TypeString): Rule ID assigned by the API
  - `rule_priority` (Computed, TypeInt): Rule priority assigned by the API

#### Scenario: Schema registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_l7_acc_rules` SHALL be registered in `tencentcloud/provider.go`

### Requirement: Resource creation via CreateL7AccRules
The Create function SHALL call `CreateL7AccRules` API with `ZoneId` and `Rules` parameters, then set the resource ID to `ZoneId#RuleId`.

#### Scenario: Successful rule creation
- **WHEN** a Terraform configuration specifies `zone_id` and at least one `rules` block
- **THEN** the Create function SHALL call `CreateL7AccRules` and return the first `RuleId` from the response
- **AND** the resource ID SHALL be set to `ZoneId#RuleId`

#### Scenario: Empty rule response
- **WHEN** `CreateL7AccRules` returns an empty `RuleIds` list
- **THEN** the Create function SHALL return a `NonRetryableError` error

### Requirement: Resource read via DescribeL7AccRules
The Read function SHALL use `DescribeL7AccRules` with `rule-id` filter to fetch the current state of the rule.

#### Scenario: Rule exists
- **WHEN** the Read function is called with a valid resource ID
- **THEN** the function SHALL parse the `ZoneId` and `RuleId` from the ID
- **AND** call `DescribeL7AccRules` with `rule-id` filter
- **AND** set all schema fields from the API response

#### Scenario: Rule not found
- **WHEN** `DescribeL7AccRules` returns an empty `Rules` list
- **THEN** the Read function SHALL log a warning and call `d.SetId("")` to remove the resource from state

### Requirement: Resource update via ModifyL7AccRule
The Update function SHALL call `ModifyL7AccRule` when rule content changes (status, rule_name, description, branches).

#### Scenario: Rule content update
- **WHEN** any of `status`, `rule_name`, `description`, or `branches` changes
- **THEN** the Update function SHALL call `ModifyL7AccRule` with the updated rule content
- **AND** the `RuleId` SHALL be included in the request

### Requirement: Resource deletion via DeleteL7AccRules
The Delete function SHALL call `DeleteL7AccRules` with the rule's `RuleId`.

#### Scenario: Rule deletion
- **WHEN** the Terraform resource is destroyed
- **THEN** the Delete function SHALL call `DeleteL7AccRules` with `ZoneId` and `RuleIds` containing the rule ID
- **AND** the deletion SHALL use retry with `tccommon.WriteRetryTimeout`

### Requirement: Error handling and retry
All API calls SHALL use `resource.Retry` with appropriate timeout values and proper error wrapping.

#### Scenario: Transient API error
- **WHEN** a Cloud API call returns a transient error
- **THEN** the function SHALL retry using `tccommon.RetryError` wrapping
- **AND** the retry timeout SHALL be `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations

### Requirement: Documentation
Each resource SHALL have a corresponding `.md` documentation file in `tencentcloud/services/teo/`.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_teo_l7_acc_rules.md` file SHALL exist with example usage and import instructions

