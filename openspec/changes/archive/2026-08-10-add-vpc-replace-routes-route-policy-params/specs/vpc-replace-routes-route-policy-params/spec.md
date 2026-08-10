## ADDED Requirements

### Requirement: Resource schema exposes need_router_info parameter
The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL expose an optional boolean parameter `need_router_info` that maps to the `NeedRouterInfo` input of the `DescribeRouteTables` cloud API, controlling whether route policy information is fetched.

#### Scenario: User configures need_router_info
- **WHEN** the user sets `need_router_info = false` on the resource
- **THEN** the Read method SHALL pass `NeedRouterInfo=false` to the `DescribeRouteTables` request so that route policy info is not fetched

#### Scenario: User omits need_router_info
- **WHEN** the user does not set `need_router_info`
- **THEN** the resource SHALL not set `NeedRouterInfo` on the `DescribeRouteTables` request, preserving the cloud API default behavior

### Requirement: Resource schema exposes name filter parameter
The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL expose an optional string parameter `name` that maps to the `Filters.Name` input of the `DescribeRouteTables` cloud API, allowing the user to filter by attribute name.

#### Scenario: User configures name filter
- **WHEN** the user sets `name = "route-table-name"` on the resource
- **THEN** the Read method SHALL add a filter with `Name="route-table-name"` to the `DescribeRouteTables` request

#### Scenario: User omits name filter
- **WHEN** the user does not set `name`
- **THEN** the resource SHALL not add a name-based filter to the `DescribeRouteTables` request

### Requirement: Resource schema exposes values filter parameter
The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL expose an optional list parameter `values` that maps to the `Filters.Values` input of the `DescribeRouteTables` cloud API, allowing the user to filter by attribute values.

#### Scenario: User configures values filter
- **WHEN** the user sets `values = ["rtb-xxx"]` on the resource
- **THEN** the Read method SHALL add the corresponding values to the filter on the `DescribeRouteTables` request

#### Scenario: User omits values filter
- **WHEN** the user does not set `values`
- **THEN** the resource SHALL not add values to the filter on the `DescribeRouteTables` request

### Requirement: Resource schema exposes total_count output
The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL expose a computed integer output parameter `total_count` that maps to the `TotalCount` field of the `DescribeRouteTables` response, representing the total number of matching instances.

#### Scenario: Read returns total count
- **WHEN** the Read method calls `DescribeRouteTables` and the response contains `TotalCount`
- **THEN** the resource SHALL set `total_count` in the state to the returned value

#### Scenario: Read returns empty response
- **WHEN** the Read method calls `DescribeRouteTables` and the response is empty
- **THEN** the resource SHALL not set `total_count` and SHALL clear the resource id

### Requirement: Service layer supports new filter parameters
The `VpcService.DescribeRouteTables` method SHALL accept `needRouterInfo *bool`, `name string`, and `values []*string` parameters and propagate them to the `DescribeRouteTables` cloud API request.

#### Scenario: Service layer passes need_router_info
- **WHEN** the caller passes a non-nil `needRouterInfo`
- **THEN** the service layer SHALL set `request.NeedRouterInfo` to the given value

#### Scenario: Service layer passes name and values
- **WHEN** the caller passes a non-empty `name` and `values`
- **THEN** the service layer SHALL construct a `Filter` with the given `Name` and `Values` and append it to `request.Filters`

### Requirement: Backward compatibility preserved
All new parameters (`need_router_info`, `name`, `values`, `total_count`) SHALL be optional or computed, and SHALL NOT break existing Terraform configurations or state.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing configuration that does not include the new parameters
- **THEN** the resource SHALL continue to function as before without requiring changes
