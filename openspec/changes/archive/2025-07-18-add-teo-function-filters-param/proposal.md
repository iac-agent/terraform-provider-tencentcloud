## Why

The `tencentcloud_teo_function` resource currently uses the `DescribeFunctions` API only with `ZoneId` and `FunctionIds` parameters for the READ operation. However, the API also supports a `Filters` parameter that allows filtering functions by name (fuzzy match) and remark (fuzzy match). Adding the `filters` parameter to the resource schema enables users to leverage this API capability within the resource definition, providing more flexible query options when reading TEO function state.

## What Changes

- Add a new `filters` Optional schema parameter to the `tencentcloud_teo_function` resource
  - `filters` is a TypeList of nested blocks, each containing `name` (string) and `values` (list of string)
  - Supported filter names: `name` (function name fuzzy match), `remark` (function description fuzzy match)
  - When `filters` is specified, it will be passed to the `DescribeFunctions` API during the READ operation

## Capabilities

### New Capabilities
- `teo-function-filters`: Adds `filters` parameter support to the `tencentcloud_teo_function` resource, mapping to the `Filters` field in the `DescribeFunctions` API request

### Modified Capabilities

## Impact

- **Files Modified**:
  - `tencentcloud/services/teo/resource_tc_teo_function.go` - Add `filters` schema field and pass it to the service layer
  - `tencentcloud/services/teo/service_tencentcloud_teo.go` - Update `DescribeTeoFunctionById` to accept and pass `filters` parameter
  - `tencentcloud/services/teo/resource_tc_teo_function.md` - Update resource documentation with `filters` parameter
  - `tencentcloud/services/teo/resource_tc_teo_function_test.go` - Add unit tests for the new `filters` parameter
- **Backward Compatible**: The new `filters` parameter is Optional; existing configurations continue working without changes
- **API Compatibility**: The `Filters` field is already supported in the `DescribeFunctionsRequest` SDK struct; no SDK upgrade required
