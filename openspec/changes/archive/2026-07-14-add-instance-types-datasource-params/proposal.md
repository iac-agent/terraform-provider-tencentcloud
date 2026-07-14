## Why

The `tencentcloud_instance_types` data source currently only passes `Zones`, `Memory`, and `DiskTypes` parameters to the `DescribeDiskConfigQuota` API call through indirect derivation — `Zones` comes from `availability_zone`, `Memory` from `memory_size`, and `DiskTypes` from the nested `cbs_filter.disk_types`. This indirect mapping limits user control when querying disk configuration quotas, as users cannot independently specify zones, memory sizes, or disk types that differ from instance type query parameters. Adding these as explicit top-level input parameters gives users direct control over the `DescribeDiskConfigQuota` API call parameters.

## What Changes

- Add `disk_types` as a top-level optional input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, mapped to `request.DiskTypes` of the `DescribeDiskConfigQuota` API
- Add `zones` as a top-level optional input parameter (TypeList of TypeString) to the `tencentcloud_instance_types` data source schema, mapped to `request.Zones` of the `DescribeDiskConfigQuota` API
- Add `memory` as a top-level optional input parameter (TypeInt) to the `tencentcloud_instance_types` data source schema, mapped to `request.Memory` of the `DescribeDiskConfigQuota` API
- Update the `dataSourceTencentCloudInstanceTypesRead` function to pass these new top-level parameters to the CBS `DescribeDiskConfigQuota` call when `cbs_filter` is present, with top-level parameters taking priority over derived values

## Capabilities

### New Capabilities
- `datasource-cvm-instance-types-params`: Adds top-level input parameters (disk_types, zones, memory) for controlling the DescribeDiskConfigQuota API call in the instance_types data source

### Modified Capabilities

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` — Add 3 new optional input parameters to schema and update read function logic
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` — Update documentation with new parameters and examples
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` — Add unit tests for new parameters

### Affected APIs
- `DescribeDiskConfigQuota` (CBS API, package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312`) — Already supports `DiskTypes`, `Zones`, and `Memory` as request parameters

### Breaking Changes
None — all new parameters are optional and backward compatible. Existing configurations continue to work without modification.

### Dependencies
None — uses existing `DescribeDiskConfigQuota` API call that already supports these parameters.
