## ADDED Requirements

### Requirement: Chinese field descriptions for teo_function_rule
The `tencentcloud_teo_function_rule` resource SHALL provide Chinese-language descriptions for all Schema fields, making the resource documentation accessible to Chinese-speaking users.

#### Scenario: All field descriptions are in Chinese
- **WHEN** a user views the `tencentcloud_teo_function_rule` resource documentation via `terraform docs` or the Terraform Registry
- **THEN** all field descriptions SHALL be displayed in Chinese instead of English

#### Scenario: Backward compatibility is maintained
- **WHEN** existing Terraform configurations using `tencentcloud_teo_function_rule` are applied
- **THEN** the resource SHALL function identically, as only Description strings are modified