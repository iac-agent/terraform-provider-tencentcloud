## 1. Schema Definition

- [x] 1.1 Add `inquiry_type` optional string field to `cbs_filter` nested schema in `data_source_tc_instance_types.go`, with description indicating valid values: `INQUIRY_CBS_CONFIG` and `INQUIRY_CVM_CONFIG`
- [x] 1.2 Add `instance_families` optional list-of-strings field to `cbs_filter` nested schema in `data_source_tc_instance_types.go`, with description indicating it overrides the instance type's family for CBS filtering

## 2. Data Source Read Function Update

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead`, extract `inquiry_type` from `cbs_filter` input and add it to `cbsFilterParams` map
- [x] 2.2 In `dataSourceTencentCloudInstanceTypesRead`, extract `instance_families` from `cbs_filter` input and add it to `cbsFilterParams` map
- [x] 2.3 In `dataSourceTencentCloudInstanceTypesRead`, update the CBS filter parameter handling so that when `inquiry_type` is not provided, default `"INQUIRY_CVM_CONFIG"` is used (matching current hardcoded behavior)
- [x] 2.4 In `dataSourceTencentCloudInstanceTypesRead`, update the CBS filter parameter handling so that when `instance_families` is not provided, the instance type's `family` field is used (matching current behavior)

## 3. CBS Service Layer Update

- [x] 3.1 In `tencentcloud/services/cbs/service_tencentcloud_cbs.go`, update the `DescribeDiskConfigQuota` method to accept and pass `inquiry_type` from `cvmInfo` map to `request.InquiryType`, with fallback to `"INQUIRY_CVM_CONFIG"` when not provided
- [x] 3.2 In `tencentcloud/services/cbs/service_tencentcloud_cbs.go`, update the `DescribeDiskConfigQuota` method to accept and pass `instance_families` from `cvmInfo` map to `request.InstanceFamilies`, with fallback to using the instance type's `family` when not provided

## 4. Documentation Update

- [x] 4.1 Update `data_source_tc_instance_types.md` to add `inquiry_type` and `instance_families` parameters in the `cbs_filter` block description
- [x] 4.2 Add example usage in `data_source_tc_instance_types.md` showing how to use `inquiry_type` and `instance_families` parameters

## 5. Testing

- [x] 5.1 Add unit test cases in `data_source_tc_instance_types_test.go` verifying that `inquiry_type` defaults to `"INQUIRY_CVM_CONFIG"` when not provided
- [x] 5.2 Add unit test cases in `data_source_tc_instance_types_test.go` verifying that `instance_families` defaults to instance type's family when not provided
- [x] 5.3 Add unit test cases in `data_source_tc_instance_types_test.go` verifying that user-provided `inquiry_type` overrides the default
- [x] 5.4 Add unit test cases in `data_source_tc_instance_types_test.go` verifying that user-provided `instance_families` overrides the instance type's family
