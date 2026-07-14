## Why

The `tencentcloud_instance_types` data source currently does not allow users to directly pass `DiskTypes`, `Zones`, and `Memory` parameters to the `DescribeDiskConfigQuota` API. When querying CBS disk configurations via the `cbs_filter` block, the `Zones` and `Memory` parameters are automatically derived from each instance type's attributes (availability_zone and memory_size from the `DescribeZoneInstanceConfigInfos` response), and `DiskTypes` can only be specified inside the `cbs_filter` block. This limits flexibility when users want to query disk configurations for specific zones, memory sizes, or disk types that may not match the instance type results.

Adding these three parameters as top-level optional input fields allows users to override or explicitly specify the `DiskTypes`, `Zones`, and `Memory` values passed to the `DescribeDiskConfigQuota` API, providing more precise control over the CBS disk configuration query.

## What Changes

- Add `disk_types` as an optional top-level input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, mapped to `request.DiskTypes` of the `DescribeDiskConfigQuota` API
- Add `zones` as an optional top-level input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, mapped to `request.Zones` of the `DescribeDiskConfigQuota` API
- Add `memory` as an optional top-level input parameter (TypeInt) to the `tencentcloud_instance_types` data source schema, mapped to `request.Memory` of the `DescribeDiskConfigQuota` API
- Update the Read function to pass these new parameters to `DescribeDiskConfigQuota` when provided, falling back to the existing behavior (derived from instance type attributes) when not specified

## Capabilities

### New Capabilities
- `instance-types-disk-config-query-params`: Adds explicit input parameters (disk_types, zones, memory) for controlling the DescribeDiskConfigQuota API query in the instance_types data source

### Modified Capabilities
_(None - no existing spec-level behavior changes)_

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add 3 new optional input fields to schema and update Read function to pass them to DescribeDiskConfigQuota
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new parameters
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Update `DescribeDiskConfigQuota` method to accept optional disk_types, zones, memory overrides
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for new parameters

### Breaking Changes
None - all new fields are optional and backward compatible

### Dependencies
None - uses existing `DescribeDiskConfigQuota` API
