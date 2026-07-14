## 1. Schema Changes

- [x] 1.1 Add `inquiry_type` optional string field to the `cbs_filter` nested schema in `data_source_tc_instance_types.go`, with description indicating valid values `INQUIRY_CBS_CONFIG` and `INQUIRY_CVM_CONFIG`
- [x] 1.2 Add `instance_families` optional list-of-strings field to the `cbs_filter` nested schema in `data_source_tc_instance_types.go`, with description indicating it specifies instance family names for CBS configuration filtering

## 2. Data Source Read Logic

- [x] 2.1 Update `dataSourceTencentCloudInstanceTypesRead` to extract `inquiry_type` from `cbs_filter` block when present, and pass it to CBS service method
- [x] 2.2 Update `dataSourceTencentCloudInstanceTypesRead` to extract `instance_families` from `cbs_filter` block when present, and pass it to CBS service method
- [x] 2.3 Ensure backward compatibility: when `inquiry_type` is not specified, default to `"INQUIRY_CVM_CONFIG"`; when `instance_families` is not specified, continue deriving from instance type results' `family` field

## 3. CBS Service Layer

- [x] 3.1 Modify `DescribeDiskConfigQuota` method in `service_tencentcloud_cbs.go` to accept `inquiryType` and `instanceFamilies` as optional override parameters alongside the existing `cvmInfo` map
- [x] 3.2 Update `DescribeDiskConfigQuota` implementation to use user-provided `inquiryType` when specified, otherwise default to `"INQUIRY_CVM_CONFIG"`
- [x] 3.3 Update `DescribeDiskConfigQuota` implementation to use user-provided `instanceFamilies` when specified, otherwise derive from `cvmInfo["family"]` as currently done

## 4. Documentation

- [x] 4.1 Update `data_source_tc_instance_types.md` to document the new `inquiry_type` parameter within `cbs_filter`, including valid values and default behavior
- [x] 4.2 Update `data_source_tc_instance_types.md` to document the new `instance_families` parameter within `cbs_filter`, including its purpose and default behavior
- [x] 4.3 Add usage example in `data_source_tc_instance_types.md` showing `cbs_filter` with the new `inquiry_type` and `instance_families` parameters

## 5. Testing

- [x] 5.1 Add unit test case in `data_source_tc_instance_types_unit_test.go` verifying `inquiry_type` is correctly extracted from `cbs_filter` and passed to CBS service
- [x] 5.2 Add unit test case in `data_source_tc_instance_types_unit_test.go` verifying `instance_families` is correctly extracted from `cbs_filter` and passed to CBS service
- [x] 5.3 Add unit test case verifying backward compatibility: when new parameters are not specified, default values are used correctly
- [x] 5.4 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass

## 6. Code Quality

- [x] 6.1 Verify no compilation errors in the modified Go files
- [x] 6.2 Verify backward compatibility: existing configurations without new parameters still work identically
