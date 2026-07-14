## Why

The `tencentcloud_instance_types` data source currently hardcodes the `InquiryType` parameter as `"INQUIRY_CVM_CONFIG"` and derives `DiskChargeType` and `InstanceFamilies` from existing CVM instance attributes when calling the `DescribeDiskConfigQuota` API. Users need the ability to directly specify these three parameters (`InquiryType`, `DiskChargeType`, `InstanceFamilies`) as input parameters for the CBS filter, enabling more flexible queries such as querying cloud disk configurations independently (`INQUIRY_CBS_CONFIG`) or filtering by specific disk charge types and instance families without relying on CVM instance attributes.

## What Changes

- Add `inquiry_type` as an optional input parameter to the `cbs_filter` nested schema of the `tencentcloud_instance_types` data source, allowing users to specify the query category (`INQUIRY_CBS_CONFIG` or `INQUIRY_CVM_CONFIG`)
- Add `instance_families` as an optional input parameter to the `cbs_filter` nested schema, allowing users to filter by specific instance family names instead of using the instance type's family
- Update the CBS service `DescribeDiskConfigQuota` method call to pass these new parameters from the data source input instead of hardcoding them
- Update data source documentation with new parameters and examples

## Capabilities

### New Capabilities
- `datasource-cvm-instance-types-cbs-input-params`: Adds optional input parameters (inquiry_type, disk_charge_type, instance_families) to the cbs_filter of the tencentcloud_instance_types data source for more flexible DescribeDiskConfigQuota API queries

### Modified Capabilities
<!-- No existing spec-level requirements are changing -->

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new parameters to cbs_filter schema and pass them to CBS service
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Update DescribeDiskConfigQuota method to accept inquiry_type, disk_charge_type, and instance_families from input params (currently hardcoded)
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new parameters
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for new parameters

### Breaking Changes
None - all changes are additive (new optional parameters within existing cbs_filter nested schema)

### Dependencies
None - uses existing `DescribeDiskConfigQuota` API call with newly exposed parameters
