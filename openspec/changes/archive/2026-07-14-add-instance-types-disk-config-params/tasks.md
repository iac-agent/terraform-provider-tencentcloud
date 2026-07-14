## 1. Schema Definition

- [x] 1.1 Add `disk_types` optional input field (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, with description listing valid values: CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD
- [x] 1.2 Add `zones` optional input field (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, with description explaining it overrides the zones parameter for DescribeDiskConfigQuota
- [x] 1.3 Add `memory` optional input field (TypeInt) to the `tencentcloud_instance_types` data source schema, with description explaining it overrides the memory parameter for DescribeDiskConfigQuota (unit: GB)

## 2. Read Function Implementation

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead`, extract `disk_types`, `zones`, and `memory` values from the ResourceData when provided
- [x] 2.2 In the CBS filter section, pass `disk_types_override`, `zones_override`, and `memory_override` keys to the `filterParams` map when the corresponding top-level parameters are provided
- [x] 2.3 For `disk_types`: when top-level `disk_types` is provided, use it as override; otherwise fall back to `cbs_filter.disk_types`
- [x] 2.4 For `zones`: when top-level `zones` is provided, use it as override; otherwise fall back to instance type's `availability_zone`
- [x] 2.5 For `memory`: when top-level `memory` is provided, use it as override; otherwise fall back to instance type's `memory_size`

## 3. CBS Service Layer Update

- [x] 3.1 In `CbsService.DescribeDiskConfigQuota` (service_tencentcloud_cbs.go), check for `disk_types_override` key in cvmInfo map; when present, use it for `request.DiskTypes` instead of `cvmInfo["disk_types"]`
- [x] 3.2 In `CbsService.DescribeDiskConfigQuota`, check for `zones_override` key in cvmInfo map; when present, use it for `request.Zones` instead of deriving from `cvmInfo["availability_zone"]`
- [x] 3.3 In `CbsService.DescribeDiskConfigQuota`, check for `memory_override` key in cvmInfo map; when present, use it for `request.Memory` instead of deriving from `cvmInfo["memory_size"]`

## 4. Documentation

- [x] 4.1 Update `data_source_tc_instance_types.md` to include the three new parameters with descriptions and usage examples
- [x] 4.2 Add example showing usage of `disk_types`, `zones`, and `memory` parameters alongside `cbs_filter`

## 5. Testing

- [x] 5.1 Add unit test cases in `data_source_tc_instance_types_test.go` for the new `disk_types` parameter
- [x] 5.2 Add unit test cases in `data_source_tc_instance_types_test.go` for the new `zones` parameter
- [x] 5.3 Add unit test cases in `data_source_tc_instance_types_test.go` for the new `memory` parameter
- [x] 5.4 Add unit test verifying backward compatibility when new parameters are not provided
- [x] 5.5 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass
