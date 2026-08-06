## Why

The `tencentcloud_teo_zone` resource currently does not expose the `TotalCount` field returned by the `DescribeZones` API. This field — representing the total number of zones matching the query — is useful for users who need to know the total count of zones in their account without resorting to the separate data source. Adding this computed field enhances the resource's observability.

## What Changes

- Add a new computed `total_count` (TypeInt) field to the `tencentcloud_teo_zone` resource schema, sourced from `DescribeZonesResponseParams.TotalCount` in the `DescribeZones` API.

## Capabilities

### New Capabilities
- `teo-zone-totalcount`: Expose the `TotalCount` field from the `DescribeZones` API response as a computed `total_count` attribute on the `tencentcloud_teo_zone` resource.

### Modified Capabilities
<!-- No existing capabilities are modified at the spec level. -->

## Impact

- **Affected code**: `tencentcloud/services/teo/resource_tc_teo_zone.go` (schema + Read function), `tencentcloud/services/teo/service_tencentcloud_teo.go` (service layer to return TotalCount), `tencentcloud/services/teo/resource_tc_teo_zone.md` (documentation)
- **Backward compatibility**: Fully backward compatible — only a new computed field is added
- **Dependencies**: No new dependencies; uses existing `tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK