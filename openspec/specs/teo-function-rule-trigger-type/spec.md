## ADDED Requirements

### Requirement: Function rule trigger type configuration
The `tencentcloud_teo_function_rule` resource SHALL support a `trigger_type` parameter that allows users to configure the function selection strategy. The parameter MUST be a string value accepting one of: `direct`, `weight`, or `region`. When not specified, the API defaults to `direct`.

#### Scenario: Create function rule with trigger_type
- **WHEN** a user creates a `tencentcloud_teo_function_rule` resource with `trigger_type = "weight"`
- **THEN** the `TriggerType` field is set to `"weight"` in the `CreateFunctionRule` API request
- **AND** the resource is created successfully

#### Scenario: Create function rule without trigger_type
- **WHEN** a user creates a `tencentcloud_teo_function_rule` resource without specifying `trigger_type`
- **THEN** the `TriggerType` field is not set in the `CreateFunctionRule` API request
- **AND** the API defaults to `direct` behavior

#### Scenario: Read function rule trigger_type
- **WHEN** a `tencentcloud_teo_function_rule` resource is read via `DescribeFunctionRules`
- **THEN** the `trigger_type` attribute SHALL be populated from the API response's `FunctionRules[].TriggerType` field

#### Scenario: Update function rule trigger_type
- **WHEN** a user changes `trigger_type` from `"direct"` to `"weight"` on an existing `tencentcloud_teo_function_rule` resource
- **THEN** the `TriggerType` field is set to `"weight"` in the `ModifyFunctionRule` API request
- **AND** the resource is updated successfully