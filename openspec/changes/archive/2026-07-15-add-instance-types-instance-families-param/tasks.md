## 1. Schema Definition

- [x] 1.1 Add `instance_families` field (TypeList of TypeString, Optional) to the `cbs_filter` nested schema in `tencentcloud/services/cvm/data_source_tc_instance_types.go`

## 2. Read Function Logic

- [x] 2.1 Modify `dataSourceTencentCloudInstanceTypesRead` in `data_source_tc_instance_types.go` to extract `instance_families` from the `cbs_filter` block and pass it to `DescribeDiskConfigQuota` via `cbsFilterParams`
- [x] 2.2 Modify `DescribeDiskConfigQuota` service method in `tencentcloud/services/cbs/service_tencentcloud_cbs.go` to accept `instance_families` from `cvmInfo` map and pass it as `request.InstanceFamilies` when provided, otherwise fall back to using the single `family` value
- [x] 2.3 Update the `cbs_filter` parameter extraction logic in `dataSourceTencentCloudInstanceTypesRead` to pass `instance_families` through `cbsFilterParams`, with priority: user-provided `instance_families` overrides the derived single `family`

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/cvm/data_source_tc_instance_types.md` to document the new `instance_families` parameter in the `cbs_filter` block with description and usage example

## 4. Testing

- [x] 4.1 Add unit test in `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` to verify the `instance_families` parameter is correctly extracted and passed to `DescribeDiskConfigQuota`
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` on the affected test file to verify correctness
