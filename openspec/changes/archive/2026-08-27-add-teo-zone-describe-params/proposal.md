## Why

The `tencentcloud_teo_zone` resource currently calls the `DescribeZones` API internally with hardcoded pagination parameters (`Offset=0`, `Limit=20`). Users cannot control pagination nor retrieve the total count of zones matching their filters. Exposing `Offset`, `Limit`, and `TotalCount` as Terraform schema parameters gives users the ability to control pagination behavior and retrieve the total zone count, which is useful for large-scale deployments and automation scenarios.

## What Changes

- Add `offset` (Optional, TypeInt) to `tencentcloud_teo_zone` resource schema — maps to `DescribeZones` request `Offset`
- Add `limit` (Optional, TypeInt) to `tencentcloud_teo_zone` resource schema — maps to `DescribeZones` request `Limit`
- Add `total_count` (Computed, TypeInt) to `tencentcloud_teo_zone` resource schema — maps to `DescribeZones` response `TotalCount`

## Capabilities

### New Capabilities
- `teo-zone-describe-pagination`: Expose `Offset` and `Limit` as optional input parameters and `TotalCount` as a computed output parameter for the `tencentcloud_teo_zone` resource's `DescribeZones` API call.

### Modified Capabilities
<!-- None — this is a new capability addition, no existing specs are modified. -->

## Impact

- **Affected files**:
  - `tencentcloud/services/teo/resource_tc_teo_zone.go`: Add `Offset`, `Limit`, `TotalCount` to schema and set `TotalCount` in Read method
  - `tencentcloud/services/teo/resource_tc_teo_zone_test.go`: Add test cases for new parameters
  - `tencentcloud/services/teo/resource_tc_teo_zone.md`: Update documentation with new parameters
- **API**: `DescribeZones` — already supports `Offset`, `Limit`, `TotalCount` in the vendor SDK (`v20220901`)
- **Backward compatibility**: Fully backward compatible — all new parameters are Optional or Computed, no existing fields are modified
- **Dependencies**: No new dependencies required; vendor SDK already contains the necessary fields