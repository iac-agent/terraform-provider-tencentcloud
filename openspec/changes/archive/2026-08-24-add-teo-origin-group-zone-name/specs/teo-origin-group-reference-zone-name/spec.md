## ADDED Requirements

### Requirement: References block exposes zone_name
The `references` nested block of the `tencentcloud_teo_origin_group` resource SHALL include a `zone_name` computed attribute of type string, sourced from the `ZoneName` field of the `OriginGroupReference` struct in the `DescribeOriginGroup` API response (JsonPath: `response.OriginGroups.References.ZoneName`).

#### Scenario: Read populates zone_name when API returns it
- **WHEN** the resource Read method calls `DescribeOriginGroup` and a response `OriginGroupReference` contains a non-nil `ZoneName` value
- **THEN** the corresponding `references` block entry in the Terraform state SHALL contain a `zone_name` attribute matching the API response value

#### Scenario: Read handles nil zone_name gracefully
- **WHEN** the resource Read method calls `DescribeOriginGroup` and a response `OriginGroupReference` has `ZoneName` set to nil
- **THEN** the `zone_name` attribute SHALL NOT be set in the corresponding `references` block entry (nil fields are skipped, matching the existing pattern for other computed fields)

#### Scenario: zone_name is computed and not user-settable
- **WHEN** a user creates or updates a `tencentcloud_teo_origin_group` resource
- **THEN** `zone_name` SHALL be a computed-only field that cannot be set by the user in the Terraform configuration

#### Scenario: Backward compatibility with existing state
- **WHEN** an existing `tencentcloud_teo_origin_group` resource state is refreshed
- **THEN** the existing `instance_type`, `instance_id`, and `instance_name` fields SHALL continue to work unchanged, and the new `zone_name` field SHALL be populated without requiring any user configuration changes
