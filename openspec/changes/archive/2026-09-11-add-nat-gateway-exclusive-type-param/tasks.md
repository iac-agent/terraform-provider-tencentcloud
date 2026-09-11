## 1. Schema Change

- [x] 1.1 Add `exclusive_type` field to the resource schema in `resource_tc_nat_gateway.go` (TypeString, Optional, Computed, ForceNew, with validation for allowed values: ExclusiveSmall, ExclusiveMedium1, ExclusiveLarge1)

## 2. Create Method Update

- [x] 2.1 Add `ExclusiveType` to the CreateNatGateway request in `resourceTencentCloudNatGatewayCreate` function when `exclusive_type` is set

## 3. Read Method Update

- [x] 3.1 Read `ExclusiveType` from DescribeNatGateways response and set it as `exclusive_type` in `resourceTencentCloudNatGatewayRead` function

## 4. Update Method Update

- [x] 4.1 Add logic to handle `exclusive_type` changes via `ResetNatGatewayConnection` API in `resourceTencentCloudNatGatewayUpdate` function (follow existing `max_concurrent` pattern)

## 5. Documentation

- [x] 5.1 Update the resource Markdown example file (`resource_tc_nat_gateway.md`) with the new `exclusive_type` parameter example usage

## 6. Unit Test

- [x] 6.1 Add unit test cases for the `exclusive_type` parameter in the resource test file