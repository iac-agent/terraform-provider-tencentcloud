## ADDED Requirements

### Requirement: function_ids parameter for tencentcloud_teo_function
The `tencentcloud_teo_function` resource SHALL support an optional `function_ids` parameter of type list of string, which maps to the `FunctionIds` field of the `DescribeFunctions` API request. When specified, the Read operation SHALL pass this list to the DescribeFunctions API to filter results by function IDs.

#### Scenario: User specifies function_ids to query multiple functions
- **WHEN** user sets `function_ids` in the tencentcloud_teo_function resource configuration
- **THEN** the Read operation SHALL pass the `function_ids` values as `FunctionIds` in the DescribeFunctions request

#### Scenario: User does not specify function_ids
- **WHEN** user does not set `function_ids` in the resource configuration
- **THEN** the Read operation SHALL NOT include `FunctionIds` in the DescribeFunctions request (unless derived from resource ID)

### Requirement: filters parameter for tencentcloud_teo_function
The `tencentcloud_teo_function` resource SHALL support an optional `filters` parameter of type list of object, which maps to the `Filters` field of the `DescribeFunctions` API request. Each filter object SHALL contain a `name` field (required string) and a `values` field (optional set of strings). The supported filter names are `name` (function name fuzzy match) and `remark` (function description fuzzy match).

#### Scenario: User specifies filters to query functions by name
- **WHEN** user sets `filters` with `name = "name"` and `values = ["test-func"]` in the resource configuration
- **THEN** the Read operation SHALL pass the filter as `Filters` in the DescribeFunctions request

#### Scenario: User specifies filters to query functions by remark
- **WHEN** user sets `filters` with `name = "remark"` and `values = ["production"]` in the resource configuration
- **THEN** the Read operation SHALL pass the filter as `Filters` in the DescribeFunctions request

#### Scenario: User does not specify filters
- **WHEN** user does not set `filters` in the resource configuration
- **THEN** the Read operation SHALL NOT include `Filters` in the DescribeFunctions request

### Requirement: functions computed attribute for tencentcloud_teo_function
The `tencentcloud_teo_function` resource SHALL support a computed `functions` attribute of type list of object, which maps to the `Functions` field of the `DescribeFunctions` API response. Each function object SHALL contain the following computed fields: `function_id`, `zone_id`, `name`, `remark`, `content`, `domain`, `create_time`, `update_time`.

#### Scenario: Read operation returns functions list
- **WHEN** the DescribeFunctions API returns a list of functions in the response
- **THEN** the Read operation SHALL map each function to the `functions` computed attribute, including all fields (function_id, zone_id, name, remark, content, domain, create_time, update_time)

#### Scenario: DescribeFunctions returns empty functions list
- **WHEN** the DescribeFunctions API returns an empty functions list
- **THEN** the `functions` attribute SHALL be set to an empty list

### Requirement: Backward compatibility of existing tencentcloud_teo_function behavior
The existing behavior of `tencentcloud_teo_function` resource SHALL NOT be affected by the new parameters. All existing fields (zone_id, function_id, name, remark, content, domain, create_time, update_time) SHALL continue to work as before. The new `function_ids`, `filters`, and `functions` parameters SHALL be optional/computed and SHALL NOT require any changes to existing Terraform configurations.

#### Scenario: Existing configuration without new parameters
- **WHEN** user applies an existing tencentcloud_teo_function configuration that does not include `function_ids`, `filters`, or `functions`
- **THEN** the resource SHALL behave identically to before, with the existing Read logic preserving the single-function query via zone_id and function_id from the composite resource ID
