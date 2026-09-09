## ADDED Requirements

### Requirement: NAT gateway deletion protection management
The system SHALL allow users to manage deletion protection for NAT gateway instances through the `deletion_protection_enabled` parameter.

#### Scenario: Create NAT gateway with deletion protection enabled
- **WHEN** user creates a `tencentcloud_nat_gateway` resource with `deletion_protection_enabled = true`
- **THEN** the NAT gateway SHALL be created with deletion protection enabled

#### Scenario: Create NAT gateway without deletion protection
- **WHEN** user creates a `tencentcloud_nat_gateway` resource without setting `deletion_protection_enabled`
- **THEN** the NAT gateway SHALL be created with the default deletion protection state determined by the cloud API

#### Scenario: Update deletion protection from disabled to enabled
- **WHEN** user updates `deletion_protection_enabled` from `false` to `true` on an existing `tencentcloud_nat_gateway` resource
- **THEN** the NAT gateway SHALL have deletion protection enabled after the update

#### Scenario: Update deletion protection from enabled to disabled
- **WHEN** user updates `deletion_protection_enabled` from `true` to `false` on an existing `tencentcloud_nat_gateway` resource
- **THEN** the NAT gateway SHALL have deletion protection disabled after the update

#### Scenario: Read deletion protection state
- **WHEN** user reads a `tencentcloud_nat_gateway` resource
- **THEN** the `deletion_protection_enabled` attribute SHALL reflect the current deletion protection state of the NAT gateway

#### Scenario: Plan shows no diff when deletion protection is unchanged
- **WHEN** user runs `terraform plan` after `deletion_protection_enabled` has been set and no change is made
- **THEN** the plan SHALL show no changes for `deletion_protection_enabled`