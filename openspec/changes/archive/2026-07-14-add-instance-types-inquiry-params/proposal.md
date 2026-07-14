## Why

The `tencentcloud_instance_types` data source currently hardcodes `InquiryType` to `"INQUIRY_CVM_CONFIG"` when calling the `DescribeDiskConfigQuota` CBS API via `cbs_filter`. Users who want to query only cloud disk configuration lists (without instance pairing) cannot do so because the inquiry type is not configurable. Adding `InquiryType` as an input parameter in `cbs_filter` allows users to control the query category for disk configuration.

## What Changes

- Add `inquiry_type` field to the `cbs_filter` block of the `tencentcloud_instance_types` data source schema, allowing users to specify the query category (`INQUIRY_CBS_CONFIG` or `INQUIRY_CVM_CONFIG`) for the `DescribeDiskConfigQuota` API call
- Update the CBS service `DescribeDiskConfigQuota` function to accept `inquiry_type` from the data source instead of hardcoding it
- Update documentation with the new parameter and usage examples

## Capabilities

### New Capabilities
- `instance-types-inquiry-type`: Adds configurable `inquiry_type` input parameter to the `cbs_filter` block of `tencentcloud_instance_types` data source, enabling users to control the `DescribeDiskConfigQuota` API query category

### Modified Capabilities

## Impact

- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add `inquiry_type` field to `cbs_filter` schema and pass it to CBS service
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Modify `DescribeDiskConfigQuota` to accept `inquiry_type` parameter instead of hardcoding
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new parameter
- Backward compatible: `inquiry_type` will default to `"INQUIRY_CVM_CONFIG"` to maintain existing behavior
