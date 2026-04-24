## Requirements

### Requirement: rule parameter in tencentcloud_teo_l7_acc_rule_v2 resource
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL include a `rule` parameter of type TypeList with MaxItems=1, marked as Optional. This parameter maps to the `Rule` field (type `*RuleEngineItem`) of the ModifyL7AccRule API.

The `rule` parameter SHALL contain the following nested fields:
- `rule_id` (TypeString, Computed): Rule ID, unique identifier of the rule
- `status` (TypeString, Optional): Rule status, values: `enable` or `disable`
- `rule_name` (TypeString, Optional): Rule name, length limit 255 characters
- `description` (TypeList of TypeString, Optional): Rule annotations
- `branches` (TypeList, Optional): Sub-rule branches, using the same schema as `TencentTeoL7RuleBranchBasicInfo`

#### Scenario: Create resource with rule parameter
- **WHEN** a user creates a `tencentcloud_teo_l7_acc_rule_v2` resource with the `rule` parameter specified containing `status`, `rule_name`, `description`, and `branches`
- **THEN** the Create function SHALL construct a `RuleEngineItem` from the `rule` parameter and send it in the `CreateL7AccRules` API request

#### Scenario: Update resource with rule parameter
- **WHEN** a user updates a `tencentcloud_teo_l7_acc_rule_v2` resource with the `rule` parameter specified
- **THEN** the Update function SHALL construct a `RuleEngineItem` from the `rule` parameter and set it as the `Rule` field in the `ModifyL7AccRule` API request

#### Scenario: Read resource with rule parameter
- **WHEN** the Read function fetches the current state of a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the function SHALL populate the `rule` parameter from the API response's `RuleEngineItem` data, including `rule_id`, `status`, `rule_name`, `description`, and `branches`

#### Scenario: Backward compatibility with existing top-level fields
- **WHEN** a user has an existing `tencentcloud_teo_l7_acc_rule_v2` configuration using top-level `status`, `rule_name`, `description`, and `branches` fields without the `rule` parameter
- **THEN** the resource SHALL continue to function correctly without requiring any configuration changes

#### Scenario: rule parameter takes precedence in Update
- **WHEN** a user specifies both the `rule` parameter and individual top-level fields (`status`, `rule_name`, `description`, `branches`) in an update
- **THEN** the Update function SHALL use the `rule` parameter to construct the `RuleEngineItem` for the `ModifyL7AccRule` API request
