## ADDED Requirements

### Requirement: Expose template review status

The `tencentcloud_ses_template` resource SHALL expose the email template's review status as a computed output attribute `template_status` of type integer.

The values SHALL be:
- `0`: Approved (template has passed review)
- `1`: Pending review (template is awaiting review)
- `2`: Rejected (template was rejected during review)

#### Scenario: Template status is read from API response
- **WHEN** the resource Read function is called after a template is created or updated
- **THEN** the `template_status` attribute SHALL be populated with the `TemplateStatus` value from the `GetEmailTemplate` API response

#### Scenario: Template status is not available
- **WHEN** the `TemplateStatus` field in the `GetEmailTemplate` API response is nil
- **THEN** the `template_status` attribute SHALL NOT be set (remain at zero value or not present in Terraform state)

#### Scenario: Template status is displayed in terraform state
- **WHEN** `terraform state show tencentcloud_ses_template.example` is run
- **THEN** the output SHALL include `template_status` with one of the values 0, 1, or 2