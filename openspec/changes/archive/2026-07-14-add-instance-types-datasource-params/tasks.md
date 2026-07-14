## 1. Schema Definition

- [x] 1.1 Add `InquiryType` optional parameter (TypeString) to the `tencentcloud_instance_types` data source schema with description "Query category for DescribeDiskConfigQuota. Valid values: INQUIRY_CBS_CONFIG, INQUIRY_CVM_CONFIG. Default is INQUIRY_CVM_CONFIG."
- [x] 1.2 Add `DiskChargeType` optional parameter (TypeString) to the data source schema with description "Disk payment model for DescribeDiskConfigQuota. Valid values: PREPAID, POSTPAID_BY_HOUR."
- [x] 1.3 Add `InstanceFamilies` optional parameter (TypeList of TypeString) to the data source schema with description "Instance family filters for DescribeDiskConfigQuota. When specified, overrides auto-populated instance families from instance type results."
- [x] 1.4 Add `Available` computed attribute (TypeBool) to the data source schema with description "Whether disk configurations are available from DescribeDiskConfigQuota."

## 2. Read Function Implementation

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead`, read the `InquiryType` parameter from `d.Get("inquiry_type")` and pass it to the CBS service when calling `DescribeDiskConfigQuota`
- [x] 2.2 In `dataSourceTencentCloudInstanceTypesRead`, read the `DiskChargeType` parameter from `d.Get("disk_charge_type")` and pass it to the CBS service, with top-level value taking precedence over `cbs_filter.disk_charge_type`
- [x] 2.3 In `dataSourceTencentCloudInstanceTypesRead`, read the `InstanceFamilies` parameter from `d.Get("instance_families")` and pass it to the CBS service, overriding auto-populated family when explicitly provided
- [x] 2.4 In `dataSourceTencentCloudInstanceTypesRead`, after calling `DescribeDiskConfigQuota`, compute and set the top-level `Available` attribute based on disk config availability (true if at least one DiskConfig has Available=true, false otherwise)
- [x] 2.5 Ensure backward compatibility: when new top-level parameters are not provided, existing behavior is preserved (InquiryType defaults to INQUIRY_CVM_CONFIG, InstanceFamilies auto-populated from family, DiskChargeType from cbs_filter, Available defaults to false when no CBS query)

## 3. CBS Service Layer Update

- [x] 3.1 Update `DescribeDiskConfigQuota` in `tencentcloud/services/cbs/service_tencentcloud_cbs.go` to accept `InquiryType` from the data source parameters instead of hardcoding "INQUIRY_CVM_CONFIG"
- [x] 3.2 Update `DescribeDiskConfigQuota` to handle top-level `DiskChargeType` parameter alongside the existing `cbs_filter.disk_charge_type`
- [x] 3.3 Update `DescribeDiskConfigQuota` to handle top-level `InstanceFamilies` parameter alongside auto-populated instance families
- [x] 3.4 Ensure the CBS service function properly handles cases where new parameters are not provided (backward compatibility)

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cvm/data_source_tc_instance_types.md` to add `InquiryType`, `DiskChargeType`, `InstanceFamilies`, and `Available` parameter descriptions and examples
- [x] 4.2 Add example usage showing how to use new parameters together with existing `cbs_filter`

## 5. Testing

- [x] 5.1 Add unit test in `data_source_tc_instance_types_test.go` to verify `InquiryType` parameter is correctly passed to the CBS API request
- [x] 5.2 Add unit test to verify `DiskChargeType` top-level parameter takes precedence over `cbs_filter.disk_charge_type`
- [x] 5.3 Add unit test to verify `InstanceFamilies` top-level parameter overrides auto-populated family
- [x] 5.4 Add unit test to verify `Available` computed attribute is correctly set based on DiskConfig availability
- [x] 5.5 Add unit test to verify backward compatibility when new parameters are not provided
