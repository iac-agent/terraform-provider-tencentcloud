## 1. Schema Definition

- [x] 1.1 Add `inquiry_type` field (TypeString, Optional) to the `cbs_filter` nested schema in `data_source_tc_instance_types.go`, with description listing valid values `INQUIRY_CBS_CONFIG` and `INQUIRY_CVM_CONFIG`

## 2. Data Source Read Logic

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead`, extract `inquiry_type` from `cbs_filter` map and add it to `cbsFilterParams` when provided by user
- [x] 2.2 Default `inquiry_type` to `"INQUIRY_CVM_CONFIG"` in `cbsFilterParams` when user does not specify it, to maintain backward compatibility

## 3. CBS Service Layer Update

- [x] 3.1 Modify `DescribeDiskConfigQuota` function in `service_tencentcloud_cbs.go` to accept `inquiry_type` from `cvmInfo` map instead of hardcoding `request.InquiryType = helper.String("INQUIRY_CVM_CONFIG")`
- [x] 3.2 Handle the case where `inquiry_type` is not present in `cvmInfo` by defaulting to `"INQUIRY_CVM_CONFIG"`

## 4. Documentation

- [x] 4.1 Update `data_source_tc_instance_types.md` to include `inquiry_type` parameter description and example usage in the `cbs_filter` block

## 5. Testing

- [x] 5.1 Add unit test case in `data_source_tc_instance_types_test.go` verifying that `inquiry_type` is correctly passed to `DescribeDiskConfigQuota` when specified
- [x] 5.2 Add unit test case verifying backward compatibility when `inquiry_type` is not specified (defaults to `INQUIRY_CVM_CONFIG`)
