## ADDED Requirements

### Requirement: tencentcloud_teo_zone resource exposes total_count
The `tencentcloud_teo_zone` resource SHALL expose a computed `total_count` attribute of type `TypeInt` that reflects the `TotalCount` value returned by the `DescribeZones` API. The value SHALL be populated during the Read operation and SHALL be read-only.

#### Scenario: total_count is populated on read
- **WHEN** the `tencentcloud_teo_zone` resource performs a Read operation and the `DescribeZones` API returns a non-nil `TotalCount`
- **THEN** the resource SHALL set `total_count` to the value of `response.Response.TotalCount`

#### Scenario: total_count is nil in API response
- **WHEN** the `tencentcloud_teo_zone` resource performs a Read operation and the `DescribeZones` API returns a nil `TotalCount`
- **THEN** the resource SHALL NOT set `total_count` (leaving it as-is or zero-value)

#### Scenario: total_count is not user-configurable
- **WHEN** a user defines a `tencentcloud_teo_zone` resource in Terraform configuration
- **THEN** the `total_count` field SHALL NOT be accepted as a user input; it SHALL be computed-only