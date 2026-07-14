## Why

The `tencentcloud_instance_types` data source calls the `DescribeDiskConfigQuota` API to retrieve CBS disk configuration information, but currently the `InquiryType` parameter is hardcoded as "INQUIRY_CVM_CONFIG" and `InstanceFamilies` is auto-populated from instance type results. Users need the ability to specify these parameters directly to query different disk configuration categories and filter by specific instance families. Additionally, the `DiskChargeType` needs to be available as a top-level parameter alongside the existing `cbs_filter` sub-parameter. The `Available` response field should also be exposed as a top-level computed attribute to indicate disk configuration availability.

## What Changes

- Add `InquiryType` as a top-level optional parameter to the `tencentcloud_instance_types` data source, allowing users to specify the query category for `DescribeDiskConfigQuota` (valid values: INQUIRY_CBS_CONFIG, INQUIRY_CVM_CONFIG)
- Add `DiskChargeType` as a top-level optional parameter to the `tencentcloud_instance_types` data source, allowing users to specify the disk payment model for `DescribeDiskConfigQuota` (valid values: PREPAID, POSTPAID_BY_HOUR)
- Add `InstanceFamilies` as a top-level optional parameter to the `tencentcloud_instance_types` data source, allowing users to specify instance family filters for `DescribeDiskConfigQuota`
- Add `Available` as a top-level computed attribute to the `tencentcloud_instance_types` data source, exposing the disk configuration availability status from `DescribeDiskConfigQuota` response

## Capabilities

### New Capabilities
- `datasource-cvm-instance-types-params`: New input parameters (InquiryType, DiskChargeType, InstanceFamilies) and output parameter (Available) for the `tencentcloud_instance_types` data source

### Modified Capabilities

## Impact

- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new schema fields and update read logic
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Update `DescribeDiskConfigQuota` to accept new parameters from the data source
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new parameters
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for new parameters
