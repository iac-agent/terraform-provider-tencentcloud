## 1. Schema Definition

- [x] 1.1 Add `inquiry_type` optional string parameter to the `tencentcloud_instance_types` datasource schema in `data_source_tc_instance_types.go`, with description and valid values (`INQUIRY_CBS_CONFIG`, `INQUIRY_CVM_CONFIG`)
- [x] 1.2 Add `disk_charge_type` optional string parameter to the `tencentcloud_instance_types` datasource schema in `data_source_tc_instance_types.go`, with description and valid values (`PREPAID`, `POSTPAID_BY_HOUR`)

## 2. Read Function Implementation

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead`, read the `inquiry_type` parameter value and pass it to the CBS service's `DescribeDiskConfigQuota` call (defaulting to `INQUIRY_CVM_CONFIG` when not specified)
- [x] 2.2 In `dataSourceTencentCloudInstanceTypesRead`, read the top-level `disk_charge_type` parameter value and pass it to the CBS service's `DescribeDiskConfigQuota` call, with top-level value taking precedence over `cbs_filter.disk_charge_type`

## 3. CBS Service Layer Update

- [x] 3.1 Update `CbsService.DescribeDiskConfigQuota` in `tencentcloud/services/cbs/service_tencentcloud_cbs.go` to accept `inquiry_type` from the parameter map instead of hardcoding `INQUIRY_CVM_CONFIG`
- [x] 3.2 Update `CbsService.DescribeDiskConfigQuota` to handle optional `disk_charge_type` from the parameter map (only set `request.DiskChargeType` when the value is provided)

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cvm/data_source_tc_instance_types.md` to document the new `inquiry_type` and `disk_charge_type` parameters with example usage

## 5. Testing

- [x] 5.1 Add unit tests in `data_source_tc_instance_types_test.go` to verify the new `inquiry_type` and `disk_charge_type` parameters are correctly passed to the `DescribeDiskConfigQuota` API
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
