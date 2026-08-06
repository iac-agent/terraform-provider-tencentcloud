## ADDED Requirements

### Requirement: Zone pagination offset parameter
The `tencentcloud_teo_zone` resource SHALL support an optional `offset` parameter of type `TypeInt` that maps to the `DescribeZones` API's `Offset` field. When provided, the provider SHALL pass this value to the `DescribeZones` API call during the Read operation to control the starting offset for paginated zone queries.

#### Scenario: Offset parameter is provided
- **WHEN** a user configures `offset = 10` in the `tencentcloud_teo_zone` resource
- **THEN** the provider SHALL pass `Offset = 10` to the `DescribeZones` API call during the Read operation

#### Scenario: Offset parameter is not provided
- **WHEN** a user does not configure the `offset` parameter
- **THEN** the provider SHALL use the default offset of `0` for the `DescribeZones` API call

### Requirement: Zone pagination limit parameter
The `tencentcloud_teo_zone` resource SHALL support an optional `limit` parameter of type `TypeInt` that maps to the `DescribeZones` API's `Limit` field. When provided, the provider SHALL pass this value to the `DescribeZones` API call during the Read operation to control the page size for paginated zone queries.

#### Scenario: Limit parameter is provided
- **WHEN** a user configures `limit = 50` in the `tencentcloud_teo_zone` resource
- **THEN** the provider SHALL pass `Limit = 50` to the `DescribeZones` API call during the Read operation

#### Scenario: Limit parameter is not provided
- **WHEN** a user does not configure the `limit` parameter
- **THEN** the provider SHALL use the default limit of `20` for the `DescribeZones` API call

### Requirement: Backward compatibility
The addition of `offset` and `limit` parameters SHALL be fully backward compatible. Existing configurations that do not include these parameters MUST continue to function identically without any changes.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user has an existing `tencentcloud_teo_zone` resource configuration without `offset` or `limit` parameters
- **THEN** the provider SHALL behave exactly as before, using default pagination values