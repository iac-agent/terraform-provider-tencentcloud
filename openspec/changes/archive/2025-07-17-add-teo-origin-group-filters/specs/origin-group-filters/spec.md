## ADDED Requirements

### Requirement: Filters parameter in origin group resource schema
The `tencentcloud_teo_origin_group` resource SHALL include a `filters` parameter of type `TypeList` with `Optional: true` and `Computed: true`. Each filter item SHALL contain three sub-fields:
- `name` (TypeString, Required): The field name to filter on. Valid values include `origin-group-id` and `origin-group-name`.
- `values` (TypeList of TypeString, Required): The filter values for the specified field.
- `fuzzy` (TypeBool, Optional): Whether to enable fuzzy matching for the filter. Defaults to false.

#### Scenario: Resource created without filters
- **WHEN** a user creates a `tencentcloud_teo_origin_group` resource without specifying the `filters` parameter
- **THEN** the resource SHALL use the default filter behavior (filter by `origin-group-id`) during the read operation, maintaining backward compatibility

#### Scenario: Resource created with filters
- **WHEN** a user creates a `tencentcloud_teo_origin_group` resource with `filters` specified
- **THEN** the resource SHALL pass the user-provided filters to the `DescribeOriginGroup` API request during the read operation

#### Scenario: Filter with origin-group-name and fuzzy matching
- **WHEN** a user specifies a filter with `name = "origin-group-name"`, `values = ["my-group"]`, and `fuzzy = true`
- **THEN** the resource SHALL pass the filter to the API with fuzzy matching enabled for the origin group name

### Requirement: Service method accepts filters and zone_id parameters
The `DescribeTeoOriginGroupById` service method SHALL accept `zoneId` and `filters` parameters in addition to the existing `originGroupId` parameter. The method SHALL pass `ZoneId` and `Filters` to the `DescribeOriginGroup` API request. When no custom filters are provided, the method SHALL use the default `origin-group-id` filter.

#### Scenario: Service method with zone_id and default filter
- **WHEN** the read operation calls `DescribeTeoOriginGroupById` with a `zoneId` and `originGroupId` but no custom `filters`
- **THEN** the service method SHALL set `ZoneId` on the request and use the default `origin-group-id` filter with the provided `originGroupId`

#### Scenario: Service method with custom filters
- **WHEN** the read operation calls `DescribeTeoOriginGroupById` with custom `filters`
- **THEN** the service method SHALL set `ZoneId` on the request and use the provided `filters` instead of the default `origin-group-id` filter

### Requirement: Read function passes filters from schema to service method
The `resourceTencentCloudTeoOriginGroupRead` function SHALL read the `filters` parameter from the resource schema and pass it to the `DescribeTeoOriginGroupById` service method. The function SHALL also pass `zone_id` from the schema to the service method.

#### Scenario: Read with filters from schema
- **WHEN** the read function is called and the resource has `filters` configured in the schema
- **THEN** the filters SHALL be converted from schema format to `[]*teo.AdvancedFilter` and passed to the service method

#### Scenario: Read without filters in schema
- **WHEN** the read function is called and the resource does not have `filters` configured
- **THEN** no custom filters SHALL be passed to the service method, and the default filter behavior SHALL be used
