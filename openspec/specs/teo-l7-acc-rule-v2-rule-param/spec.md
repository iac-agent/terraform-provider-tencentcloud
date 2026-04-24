## Requirements

### Requirement: rule parameter schema definition
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL include a `rule` parameter of type TypeList with MaxItems: 1 and Optional:true. The `rule` block SHALL contain the following nested fields:
- `status` (TypeString, Optional): Rule status, values: `enable` or `disable`
- `rule_name` (TypeString, Optional): Rule name, max 255 characters
- `description` (TypeList of TypeString, Optional): Rule annotations
- `branches` (TypeList, Optional): Sub-rule branches, using the same `TencentTeoL7RuleBranchBasicInfo` schema as the existing `branches` field

#### Scenario: User specifies rule block in resource configuration
- **WHEN** a user configures the `rule` block with `status`, `rule_name`, `description`, and `branches` sub-fields
- **THEN** the resource SHALL accept the configuration without error and use the `rule` block values for API calls

#### Scenario: User omits rule block
- **WHEN** a user does not configure the `rule` block and instead uses top-level fields (`status`, `rule_name`, `description`, `branches`)
- **THEN** the resource SHALL behave exactly as before, using the top-level fields for API calls

### Requirement: rule parameter in Create operation
When creating a `tencentcloud_teo_l7_acc_rule_v2` resource, if the `rule` block is specified, the Create function SHALL construct the `RuleEngineItem` from the `rule` block's sub-fields and pass it to the `CreateL7AccRules` API request.

#### Scenario: Create with rule block
- **WHEN** a user creates a resource with the `rule` block specified containing `status`, `rule_name`, `description`, and `branches`
- **THEN** the Create function SHALL construct the `RuleEngineItem` from the `rule` block values and set it in the `Rules` field of `CreateL7AccRulesRequest`

#### Scenario: Create without rule block
- **WHEN** a user creates a resource without the `rule` block but with top-level fields
- **THEN** the Create function SHALL construct the `RuleEngineItem` from the top-level fields as it currently does

### Requirement: rule parameter in Update operation
When updating a `tencentcloud_teo_l7_acc_rule_v2` resource, if the `rule` block has changes, the Update function SHALL construct the `Rule` field of `ModifyL7AccRuleRequest` from the `rule` block's sub-fields.

#### Scenario: Update with rule block changes
- **WHEN** a user updates the `rule` block (e.g., changes `status` or `branches` inside the `rule` block)
- **THEN** the Update function SHALL detect the change via `d.HasChange("rule")`, construct the `Rule` from the `rule` block, and call `ModifyL7AccRule`

#### Scenario: Update with top-level field changes (no rule block)
- **WHEN** a user updates top-level fields (`status`, `rule_name`, `description`, `branches`) without a `rule` block
- **THEN** the Update function SHALL detect changes via `d.HasChange` for the individual fields and construct the `Rule` from top-level fields as before

### Requirement: rule parameter in Read operation
When reading a `tencentcloud_teo_l7_acc_rule_v2` resource, the Read function SHALL populate the `rule` block from the API response if the `rule` block is present in the configuration.

#### Scenario: Read populates rule block
- **WHEN** a Read operation is performed and the resource has a `rule` block configured
- **THEN** the Read function SHALL set the `rule` block fields (`status`, `rule_name`, `description`, `branches`) from the API response

#### Scenario: Read populates top-level fields
- **WHEN** a Read operation is performed and the resource does NOT have a `rule` block configured
- **THEN** the Read function SHALL set the top-level fields from the API response as it currently does

### Requirement: rule parameter unit tests
Unit tests SHALL be added for the new `rule` parameter in the `resource_tc_teo_l7_acc_rule_v2_test.go` file using gomonkey mock approach.

#### Scenario: Unit test for Create with rule block
- **WHEN** the Create function is called with a `rule` block in the schema
- **THEN** the test SHALL verify that the `RuleEngineItem` is correctly constructed from the `rule` block values

#### Scenario: Unit test for Update with rule block
- **WHEN** the Update function is called with changes in the `rule` block
- **THEN** the test SHALL verify that the `ModifyL7AccRule` request is correctly constructed from the `rule` block values
