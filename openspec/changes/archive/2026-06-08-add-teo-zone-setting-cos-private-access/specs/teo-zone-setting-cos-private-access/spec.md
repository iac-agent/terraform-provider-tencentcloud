## ADDED Requirements

### Requirement: CosPrivateAccess parameter in origin block
The `tencentcloud_teo_zone_setting` resource SHALL expose a `cos_private_access` parameter within the `origin` nested block. This parameter SHALL be of type string, optional, and computed. Valid values SHALL be `on` (private access) and `off` (public access). The parameter SHALL map to the `CosPrivateAccess` field in the cloud API's `Origin` struct.

#### Scenario: Read CosPrivateAccess from DescribeZoneSetting
- **WHEN** the resource is read and the API response contains `ZoneSetting.Origin.CosPrivateAccess` with a non-nil value
- **THEN** the value SHALL be set to the `cos_private_access` field in the `origin` block of the Terraform state

#### Scenario: Read CosPrivateAccess when nil in response
- **WHEN** the resource is read and the API response has `ZoneSetting.Origin.CosPrivateAccess` as nil
- **THEN** the `cos_private_access` field SHALL NOT be set (skipped), consistent with existing nil-check patterns

#### Scenario: Update CosPrivateAccess via ModifyZoneSetting
- **WHEN** the `cos_private_access` field is changed in the Terraform configuration
- **THEN** the value SHALL be included in the `Origin` struct of the `ModifyZoneSetting` API request

#### Scenario: CosPrivateAccess not specified in configuration
- **WHEN** a Terraform configuration does not include `cos_private_access`
- **THEN** the existing behavior SHALL be preserved (no value sent in update request, value read from API on refresh)
