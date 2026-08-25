## Requirements

### Requirement: Origin group reference exposes alias zone name
The `tencentcloud_teo_origin_group` resource SHALL expose the `alias_zone_name` computed attribute within the `references` block, sourced from the `OriginGroupReference.AliasZoneName` field returned by the `DescribeOriginGroup` API.

#### Scenario: AliasZoneName is returned by the API
- **WHEN** the `DescribeOriginGroup` API returns `OriginGroupReference` entries with `AliasZoneName` set to a non-nil string value
- **THEN** the Terraform state SHALL contain the corresponding `alias_zone_name` value in the `references` block

#### Scenario: AliasZoneName is nil in the API response
- **WHEN** the `DescribeOriginGroup` API returns `OriginGroupReference` entries with `AliasZoneName` set to nil
- **THEN** the Terraform state SHALL omit the `alias_zone_name` field from that reference entry (or leave it as empty string)