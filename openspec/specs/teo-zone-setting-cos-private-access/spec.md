## ADDED Requirements

### Requirement: cos_private_access parameter in origin block
The `tencentcloud_teo_zone_setting` resource SHALL include a `cos_private_access` field of type `string` within the `origin` block. The field SHALL be optional and computed. Valid values SHALL be `"on"` (private access) and `"off"` (public access). The field description SHALL be: "Whether to enable private access for the origin server bucket when the origin is a Tencent Cloud COS bucket. Valid values: `on` for private access, `off` for public access."

#### Scenario: User sets cos_private_access to on
- **WHEN** a user configures `origin.cos_private_access = "on"` in their Terraform configuration
- **THEN** the resource SHALL send `CosPrivateAccess: "on"` in the `ModifyZoneSetting` API request's `Origin` field

#### Scenario: User sets cos_private_access to off
- **WHEN** a user configures `origin.cos_private_access = "off"` in their Terraform configuration
- **THEN** the resource SHALL send `CosPrivateAccess: "off"` in the `ModifyZoneSetting` API request's `Origin` field

#### Scenario: User does not specify cos_private_access
- **WHEN** a user does not include `cos_private_access` in the `origin` block
- **THEN** the resource SHALL NOT send `CosPrivateAccess` in the `ModifyZoneSetting` API request, and the field value SHALL be read from the `DescribeZoneSetting` API response

#### Scenario: Reading cos_private_access from API response
- **WHEN** the `DescribeZoneSetting` API returns `ZoneSetting.Origin.CosPrivateAccess` with a non-nil value
- **THEN** the resource SHALL set `cos_private_access` in the Terraform state to the returned value

#### Scenario: Reading cos_private_access when API returns nil
- **WHEN** the `DescribeZoneSetting` API returns `ZoneSetting.Origin.CosPrivateAccess` as nil
- **THEN** the resource SHALL NOT set `cos_private_access` in the Terraform state, consistent with existing origin field handling
