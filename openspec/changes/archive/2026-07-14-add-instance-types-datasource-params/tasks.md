## 1. Schema Definition

- [x] 1.1 Add `disk_types` optional input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, with description listing valid values: CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD
- [x] 1.2 Add `zones` optional input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, with description explaining it specifies availability zones for disk config quota query
- [x] 1.3 Add `memory` optional input parameter (TypeInt) to the `tencentcloud_instance_types` data source schema, with description explaining it specifies instance memory size in GB for disk config quota query

## 2. Read Function Implementation

- [x] 2.1 Update `dataSourceTencentCloudInstanceTypesRead` to read top-level `disk_types` parameter and use it to override `cbs_filter.disk_types` in the `DescribeDiskConfigQuota` call when both are present
- [x] 2.2 Update `dataSourceTencentCloudInstanceTypesRead` to read top-level `zones` parameter and use it to override the derived `availability_zone` in the `DescribeDiskConfigQuota` call when present
- [x] 2.3 Update `dataSourceTencentCloudInstanceTypesRead` to read top-level `memory` parameter and use it to override the derived `memory_size` in the `DescribeDiskConfigQuota` call when present

## 3. Documentation

- [x] 3.1 Update `data_source_tc_instance_types.md` to document the new `disk_types`, `zones`, and `memory` parameters with descriptions and example usage

## 4. Testing

- [x] 4.1 Add unit test in `data_source_tc_instance_types_test.go` for `disk_types` top-level parameter overriding `cbs_filter.disk_types`
- [x] 4.2 Add unit test for `zones` top-level parameter overriding derived `availability_zone`
- [x] 4.3 Add unit test for `memory` top-level parameter overriding derived `memory_size`
- [x] 4.4 Add unit test verifying backward compatibility when new parameters are not specified
- [x] 4.5 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
