## Why

The `tencentcloud_teo_zone` resource currently does not expose `Offset` and `Limit` parameters for the `DescribeZones` API call used during the Read operation. Exposing these parameters allows users to control pagination behavior when the provider internally queries zone information, improving consistency with the underlying cloud API and enabling more flexible data retrieval patterns.

## What Changes

- Add `offset` (Optional, TypeInt) parameter to the `tencentcloud_teo_zone` resource schema, mapping to `DescribeZones` API's `Offset` field
- Add `limit` (Optional, TypeInt) parameter to the `tencentcloud_teo_zone` resource schema, mapping to `DescribeZones` API's `Limit` field
- Pass these parameters from the resource schema to the `DescribeTeoZoneById` service method, which internally calls `DescribeZones`

## Capabilities

### New Capabilities
- `teo-zone-pagination-params`: Support for `Offset` and `Limit` pagination parameters in the `tencentcloud_teo_zone` resource's Read operation, allowing users to control the internal `DescribeZones` API pagination behavior

### Modified Capabilities
<!-- No existing capabilities need modification at the spec level -->

## Impact

- **Affected code**: `tencentcloud/services/teo/resource_tc_teo_zone.go` (schema and Read method), `tencentcloud/services/teo/service_tencentcloud_teo.go` (`DescribeTeoZoneById` method)
- **Affected documentation**: `tencentcloud/services/teo/resource_tc_teo_zone.md` (will be auto-generated via `make doc`)
- **Backward compatibility**: Fully backward compatible — both new parameters are Optional, and existing configurations continue to work without changes
- **Dependencies**: No new dependencies required; uses existing `tencentcloud-sdk-go/tencentcloud/teo/v20220901` package