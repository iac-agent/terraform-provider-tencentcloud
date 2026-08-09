## Why

The `tencentcloud_teo_origin_group` resource currently does not support the `filters` parameter in its read operation when calling the `DescribeOriginGroup` API. The API provides an `AdvancedFilter` mechanism that allows filtering by `origin-group-id` and `origin-group-name` with optional fuzzy matching. Adding the `filters` parameter to the resource schema enables users to configure filtering behavior for origin group queries, improving the resource's query capabilities and aligning the Terraform resource with the full API surface.

## What Changes

- Add a new `filters` parameter (TypeList, Optional, Computed) to the `tencentcloud_teo_origin_group` resource schema
- Each filter item contains `name` (string), `values` (list of string), and `fuzzy` (bool) sub-fields, mapping to the `AdvancedFilter` struct in the cloud API
- Update the `resourceTencentCloudTeoOriginGroupRead` function to pass `filters` to the `DescribeOriginGroup` API request
- Update the `DescribeTeoOriginGroupById` service method to accept and pass `filters` parameter
- Update unit tests to cover the new `filters` parameter
- Update the resource documentation (.md file)

## Capabilities

### New Capabilities
- `origin-group-filters`: Add filters parameter to tencentcloud_teo_origin_group resource for filtering DescribeOriginGroup API queries

### Modified Capabilities


## Impact

- **Resource file**: `tencentcloud/services/teo/resource_tc_teo_origin_group.go` - schema definition, read function
- **Service file**: `tencentcloud/services/teo/service_tencentcloud_teo.go` - `DescribeTeoOriginGroupById` method signature and implementation
- **Test file**: `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go` - unit test coverage for filters
- **Documentation**: `tencentcloud/services/teo/resource_tc_teo_origin_group.md` - add filters parameter documentation
- **Cloud API**: `DescribeOriginGroup` API in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` - already supports `Filters` field (AdvancedFilter type)
