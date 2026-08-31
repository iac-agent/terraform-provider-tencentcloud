## Why

The `tencentcloud_teo_zone` resource currently lacks several parameters that are available in the underlying TEO Cloud API (`CreateZone`, `DescribeZones`, `ModifyZone`). Specifically, the `AllowDuplicates` and `JumpStart` flags for `CreateZone` are not exposed, and the `ownership_verification` block is incomplete — missing `file_verification` and `ns_verification` sub-blocks that the API returns. Adding these parameters enables users to control duplicate site access behavior, skip DNS record scanning during creation, and access complete ownership verification information.

## What Changes

- **New parameter `allow_duplicates`** (TypeBool, Optional, ForceNew): Controls whether duplicate zone access is allowed. Defaults to false. Passed to `CreateZone` API.
- **New parameter `jump_start`** (TypeBool, Optional, ForceNew): Controls whether to skip existing DNS record scanning during zone creation. Defaults to false. Passed to `CreateZone` API.
- **Extend `ownership_verification` block**: Add `file_verification` sub-block with `path` and `content` fields (Computed), and `ns_verification` sub-block with `name_servers` field list (Computed). These are populated from the `CreateZone` response and `DescribeZones` response.

## Capabilities

### New Capabilities
- `teo-zone-create-params`: Expose `AllowDuplicates` and `JumpStart` parameters on `tencentcloud_teo_zone` resource for the `CreateZone` API call, and extend `ownership_verification` with `file_verification` and `ns_verification` sub-blocks.

### Modified Capabilities
<!-- None - all changes are additive new capabilities, no existing requirements change -->

## Impact

- **Affected code**: `tencentcloud/services/teo/resource_tc_teo_zone.go` (schema + Create/Read methods), `tencentcloud/services/teo/resource_tc_teo_zone_extension.go` (if needed), `tencentcloud/services/teo/resource_tc_teo_zone_test.go` (unit tests), `website/docs/r/teo_zone.html.markdown` (regenerated via `make doc`)
- **API dependencies**: `CreateZone` (v20220901), `DescribeZones` (v20220901) — already used by the resource
- **Backward compatibility**: All changes are additive (new Optional/Computed fields only), no breaking changes