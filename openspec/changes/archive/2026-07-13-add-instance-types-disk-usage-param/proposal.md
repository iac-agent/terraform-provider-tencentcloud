## Why

The `tencentcloud_instance_types` data source calls the `DescribeDiskConfigQuota` API (via CBS service) to retrieve disk configuration information when `cbs_filter` is provided. While `disk_usage` is already available as an input parameter in `cbs_filter` and as a computed field inside `cbs_configs`, users need `DiskUsage` exposed as a top-level computed output parameter at the data source level, so they can directly reference the disk usage type (SYSTEM_DISK or DATA_DISK) used in the CBS query without having to navigate into the nested `cbs_configs` structure.

## What Changes

- Add a new computed output parameter `DiskUsage` (TypeString) to the `tencentcloud_instance_types` data source schema, at the top level of the resource
- Map `DiskUsage` from the `DescribeDiskConfigQuota` API response (`DiskConfig.DiskUsage`) to this new top-level field
- Update the data source documentation (`data_source_tc_instance_types.md`) with the new `DiskUsage` field

## Capabilities

### New Capabilities
- `instance-types-disk-usage-output`: Add DiskUsage as a top-level computed output parameter to the tencentcloud_instance_types data source, exposing the disk usage type (SYSTEM_DISK/DATA_DISK) from the DescribeDiskConfigQuota API response

### Modified Capabilities
(No existing capabilities are being modified at the spec level)

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new `DiskUsage` computed field to schema and mapping logic in read function
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new field description
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for the new field

### Breaking Changes
None - the new `DiskUsage` field is an optional computed field, fully backward compatible

### Dependencies
- Uses existing `DescribeDiskConfigQuota` API call via CBS service
- The `DiskUsage` field already exists in the CBS SDK `DiskConfig` struct (`*string` type)
