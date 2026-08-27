## ADDED Requirements

### Requirement: TencentCloud TEO Zone resource supports pagination parameters
The `tencentcloud_teo_zone` resource SHALL expose `offset` and `limit` as optional integer input parameters and `total_count` as a computed integer output parameter in its Terraform schema, mapping to the `DescribeZones` API's `Offset`, `Limit`, and `TotalCount` fields respectively.

#### Scenario: User provides offset and limit
- **WHEN** a user configures `offset = 0` and `limit = 50` in a `tencentcloud_teo_zone` resource
- **THEN** the Terraform plan SHALL accept these values without error

#### Scenario: User reads total_count after zone is created
- **WHEN** a `tencentcloud_teo_zone` resource is created and the Read operation calls `DescribeZones`
- **THEN** `total_count` SHALL be set to the value returned by the API's `TotalCount` field

#### Scenario: Backward compatibility with existing configurations
- **WHEN** an existing `tencentcloud_teo_zone` configuration is applied without specifying `offset` or `limit`
- **THEN** the resource SHALL continue to work unchanged, with no drift detected in terraform plan