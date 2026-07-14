## 1. Schema Definition

- [x] 1.1 Add `InstanceFamilies` schema field (TypeList of TypeString, Optional) to `tencentcloud_instance_types` data source in `data_source_tc_instance_types.go`
- [x] 1.2 Add `DiskTypes` schema field (TypeList of TypeString, Optional) to `tencentcloud_instance_types` data source in `data_source_tc_instance_types.go`

## 2. Read Function Implementation

- [x] 2.1 Update `dataSourceTencentCloudInstanceTypesRead` to read `InstanceFamilies` parameter and translate it to `instance-family` filter values in the filterMap for `DescribeZoneInstanceConfigInfos`
- [x] 2.2 Handle merge of `InstanceFamilies` values with existing `filter` block `instance-family` values in the filterMap
- [x] 2.3 Update `dataSourceTencentCloudInstanceTypesRead` to read `DiskTypes` parameter and use it to trigger CBS config queries when `DiskTypes` is provided
- [x] 2.4 When both `DiskTypes` top-level parameter and `cbs_filter.disk_types` are provided, use `DiskTypes` values with precedence over `cbs_filter.disk_types`
- [x] 2.5 When `InstanceFamilies` is provided and CBS config queries are triggered, pass `InstanceFamilies` values directly to `DescribeDiskConfigQuota` API's `InstanceFamilies` field instead of auto-populating from instance type family

## 3. CBS Service Layer Update

- [x] 3.1 Update `DescribeDiskConfigQuota` in `tencentcloud/services/cbs/service_tencentcloud_cbs.go` to accept optional `InstanceFamilies` parameter that overrides the auto-populated family value when provided
- [x] 3.2 Ensure `DescribeDiskConfigQuota` can handle both the auto-populated family (single value) and user-specified `InstanceFamilies` (multiple values) scenarios

## 4. Documentation

- [x] 4.1 Update `data_source_tc_instance_types.md` to include `InstanceFamilies` and `DiskTypes` parameters with descriptions and usage examples
- [x] 4.2 Add example showing how to use `InstanceFamilies` to filter by instance family
- [x] 4.3 Add example showing how to use `DiskTypes` to filter CBS config by disk type

## 5. Testing

- [x] 5.1 Add unit test cases in `data_source_tc_instance_types_unit_test.go` to verify `InstanceFamilies` parameter correctly filters instance types via `instance-family` filter
- [x] 5.2 Add unit test cases in `data_source_tc_instance_types_unit_test.go` to verify `DiskTypes` parameter correctly triggers CBS config queries and passes values to `DescribeDiskConfigQuota`
- [x] 5.3 Add unit test cases to verify backward compatibility when `InstanceFamilies` and `DiskTypes` are not provided
- [x] 5.4 Add unit test cases to verify merge behavior when `InstanceFamilies` and `filter` block `instance-family` are both provided
- [x] 5.5 Add unit test cases to verify precedence behavior when `DiskTypes` and `cbs_filter.disk_types` are both provided

## 6. Code Quality

- [x] 6.1 Run `go test` with `-gcflags=all=-l` on modified test files to verify unit tests pass
