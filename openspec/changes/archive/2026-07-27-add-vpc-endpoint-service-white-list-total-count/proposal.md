## Why

The `tencentcloud_vpc_end_point_service_white_list` resource currently does not expose the total count of matching white list entries returned by the `DescribeVpcEndPointServiceWhiteList` API. Adding a computed `total_count` field allows users to reference the total number of white list records that match the query, improving observability and enabling downstream configuration logic based on the count.

## What Changes

- Add a new computed field `total_count` (TypeInt) to the `tencentcloud_vpc_end_point_service_white_list` resource schema, mapped to the `TotalCount` field of the `DescribeVpcEndPointServiceWhiteList` API response.
- Modify the VPC service layer method `DescribeVpcEndPointServiceWhiteListById` to return the `TotalCount` value from the API response alongside the existing `*vpc.VpcEndPointServiceUser` record, so the resource Read function can populate the new field.
- Update the resource Read function to set `total_count` when the returned total count is non-nil.
- Add unit tests covering the new field using gomonkey mocks.
- Update the resource documentation (`resource_tc_vpc_end_point_service_white_list.md`).

## Capabilities

### New Capabilities
- `vpc-end-point-service-white-list`: Manages the `tencentcloud_vpc_end_point_service_white_list` resource, including its schema fields and CRUD behavior against the TencentCloud VPC EndPointServiceWhiteList APIs.

### Modified Capabilities
<!-- None. No existing spec-level behavior is being changed. -->

## Impact

- **Code**:
  - `tencentcloud/services/pls/resource_tc_vpc_end_point_service_white_list.go` (schema + read)
  - `tencentcloud/services/vpc/service_tencentcloud_vpc.go` (`DescribeVpcEndPointServiceWhiteListById` signature/return value)
  - `tencentcloud/services/pls/resource_tc_vpc_end_point_service_white_list_test.go` (unit tests)
- **APIs**: TencentCloud VPC `DescribeVpcEndPointServiceWhiteList` (response field `TotalCount`, already supported by the SDK in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312/models.go`).
- **Compatibility**: The new field is Computed-only and backward compatible; existing configurations and state are unaffected.
