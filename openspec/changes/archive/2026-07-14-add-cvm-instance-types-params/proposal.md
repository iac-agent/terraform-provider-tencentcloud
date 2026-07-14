## Motivation

The `tencentcloud_instance_types` datasource currently hardcodes `InquiryType` to `INQUIRY_CVM_CONFIG` when calling the CBS API's `DescribeDiskConfigQuota`, preventing users from querying pure cloud disk configurations (`INQUIRY_CBS_CONFIG`). Additionally, `DiskChargeType` is only accessible as a nested parameter inside the `cbs_filter` block, requiring users to use the nested block structure even if they only want to specify the payment model.

## Changes

- Add `inquiry_type` as an optional top-level parameter to the `tencentcloud_instance_types` datasource schema, defaulting to `INQUIRY_CVM_CONFIG` for backward compatibility
- Add `disk_charge_type` as an optional top-level parameter to the `tencentcloud_instance_types` datasource schema, taking precedence over `cbs_filter.disk_charge_type` when both are specified
- Update `CbsService.DescribeDiskConfigQuota` to read `inquiry_type` from the parameter map instead of hardcoding it
- Update `CbsService.DescribeDiskConfigQuota` to handle optional `disk_charge_type` from the parameter map
- Update datasource documentation with new parameter descriptions and examples
- Add gomonkey unit tests for the new parameters
