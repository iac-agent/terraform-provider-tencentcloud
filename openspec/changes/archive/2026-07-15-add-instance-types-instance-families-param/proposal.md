## Why

The `tencentcloud_instance_types` data source's `cbs_filter` block currently only supports passing a single instance family (derived from the instance type query result) to the `DescribeDiskConfigQuota` API. Users need the ability to filter disk configuration quotas by multiple instance families and specify disk charge type directly, which the current implementation does not support. Adding `InstanceFamilies` as a new parameter in `cbs_filter` will allow users to query disk configurations for multiple instance families at once.

## What Changes

- Add `instance_families` field (TypeList of TypeString) to the `cbs_filter` nested schema in the `tencentcloud_instance_types` data source
- Modify the `DescribeDiskConfigQuota` call logic in the data source read function to pass the `instance_families` parameter from `cbs_filter` to the API request, instead of only using the single `family` value from the instance type query result
- Update the data source documentation (`data_source_tc_instance_types.md`) with the new parameter and usage examples

## Capabilities

### New Capabilities
- `cvm-instance-types-instance-families`: Add `instance_families` parameter to the `cbs_filter` block of the `tencentcloud_instance_types` data source, enabling multi-instance-family filtering for disk configuration quota queries

### Modified Capabilities

## Impact

- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add `instance_families` field to `cbs_filter` schema, modify read function to pass `InstanceFamilies` to `DescribeDiskConfigQuota`
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Modify `DescribeDiskConfigQuota` method to accept `instance_families` parameter instead of deriving it from a single `family` value
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Add documentation for the new `instance_families` parameter
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add unit test for the new parameter
