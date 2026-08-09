## ADDED Requirements

### Requirement: template_id schema field
The `tencentcloud_teo_web_security_template` resource SHALL include a `template_id` field in its schema definition. The field SHALL be of type `schema.TypeString`, with `Optional: true` and `Computed: true`. The description SHALL indicate it is the policy template ID.

#### Scenario: template_id is populated after resource creation
- **WHEN** a `tencentcloud_teo_web_security_template` resource is created
- **THEN** the `template_id` field SHALL be set from the `TemplateId` value returned in the `CreateWebSecurityTemplate` API response

#### Scenario: template_id is populated during resource read
- **WHEN** the `tencentcloud_teo_web_security_template` resource is read (refresh)
- **THEN** the `template_id` field SHALL be set from the second component of the composite resource ID (zoneId#templateId)

#### Scenario: template_id is immutable
- **WHEN** a user attempts to change the `template_id` value in a Terraform configuration
- **THEN** the update function SHALL return an error indicating the argument cannot be changed

#### Scenario: template_id not specified by user
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource without specifying `template_id`
- **THEN** the resource SHALL be created successfully and `template_id` SHALL be computed from the API response
