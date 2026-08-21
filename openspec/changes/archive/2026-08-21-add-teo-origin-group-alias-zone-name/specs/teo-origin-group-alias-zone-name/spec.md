## ADDED Requirements

### Requirement: References block exposes alias zone name
The `references` nested block of the `tencentcloud_teo_origin_group` resource SHALL include `alias_zone_name` as a computed string attribute, sourced from the `AliasZoneName` field of the `OriginGroupReference` struct in the `DescribeOriginGroup` API response.

#### Scenario: Read populates alias zone name when API returns it
- **WHEN** the resource Read method calls `DescribeOriginGroup` and a response `OriginGroupReference` contains a non-nil `AliasZoneName` value
- **THEN** the Terraform state SHALL contain `alias_zone_name` in the corresponding `references` block entry matching the API response value

#### Scenario: Read handles nil alias zone name gracefully
- **WHEN** the resource Read method calls `DescribeOriginGroup` and a response `OriginGroupReference` has `AliasZoneName` set to nil
- **THEN** the corresponding Terraform state field SHALL NOT be set, matching the existing nil-skip pattern for other computed fields

#### Scenario: Alias zone name is computed and not user-settable
- **WHEN** a user creates or updates a `tencentcloud_teo_origin_group` resource
- **THEN** `alias_zone_name` SHALL be a computed-only field that cannot be set by the user in the Terraform configuration

#### Scenario: Backward compatibility with existing state
- **WHEN** an existing `tencentcloud_teo_origin_group` resource state is refreshed
- **THEN** the existing `instance_type`, `instance_id`, and `instance_name` fields SHALL continue to work unchanged, and `alias_zone_name` SHALL be populated without requiring any user configuration changes
