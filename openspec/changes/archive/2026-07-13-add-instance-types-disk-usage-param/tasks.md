## 1. Schema Definition

- [x] 1.1 Add `DiskUsage` top-level computed field (TypeString) to the `tencentcloud_instance_types` data source schema in `tencentcloud/services/cvm/data_source_tc_instance_types.go`
- [x] 1.2 Add description for the `DiskUsage` field: "Cloud disk type. Value range: SYSTEM_DISK (system disk), DATA_DISK (data disk). Only populated when cbs_filter is provided."

## 2. Data Mapping Implementation

- [x] 2.1 In `dataSourceTencentCloudInstanceTypesRead` function, set the top-level `DiskUsage` computed field value from `cbs_filter.disk_usage` input when `cbs_filter` is provided
- [x] 2.2 Ensure `DiskUsage` remains empty/null when `cbs_filter` is not provided

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/cvm/data_source_tc_instance_types.md` to add the new `DiskUsage` computed output field with description and example usage

## 4. Testing

- [x] 4.1 Add unit test in `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` to verify `DiskUsage` is correctly populated when `cbs_filter` with `disk_usage` is provided
- [x] 4.2 Add unit test to verify `DiskUsage` is empty/null when `cbs_filter` is not provided
- [x] 4.3 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
