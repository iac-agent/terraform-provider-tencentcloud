# teo-l4-proxy Specification

## Requirements

### Requirement: L4 Proxy Area Field is Required
The `area` field of the `tencentcloud_teo_l4_proxy` resource SHALL be a required field. Users MUST provide the `area` value when creating or configuring an L4 proxy instance. The field SHALL accept one of three values: `mainland`, `overseas`, or `global`.

#### Scenario: User creates L4 proxy with area specified
- **WHEN** user defines a `tencentcloud_teo_l4_proxy` resource with a valid `area` value (e.g., `"overseas"`)
- **THEN** the resource is created successfully with the specified acceleration zone

#### Scenario: User creates L4 proxy without area
- **WHEN** user defines a `tencentcloud_teo_l4_proxy` resource without specifying the `area` field
- **THEN** Terraform validation fails with a missing required field error during `plan` or `apply`

#### Scenario: User attempts to update area on existing resource
- **WHEN** user modifies the `area` field on an existing `tencentcloud_teo_l4_proxy` resource
- **THEN** the update operation fails because `area` is immutable and cannot be changed after creation