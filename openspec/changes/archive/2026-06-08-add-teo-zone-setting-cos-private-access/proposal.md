## Why

The `tencentcloud_teo_zone_setting` resource currently lacks support for the `CosPrivateAccess` parameter in the `origin` block. This parameter controls whether the origin server accesses a Tencent Cloud COS bucket privately or publicly, which is essential for users who need to configure private bucket access for their TEO site origins.

## What Changes

- Add a new `cos_private_access` parameter (string, optional, computed) to the `origin` nested block of the `tencentcloud_teo_zone_setting` resource
- The parameter maps to `Origin.CosPrivateAccess` in both the `ModifyZoneSetting` and `DescribeZoneSetting` cloud API interfaces
- Valid values: `on` (private access) and `off` (public access)

## Capabilities

### New Capabilities
- `teo-zone-setting-cos-private-access`: Adds the `cos_private_access` parameter to the `origin` block of the `tencentcloud_teo_zone_setting` resource, enabling users to configure private/public access for COS origin buckets in TEO.

### Modified Capabilities

## Impact

- Affected files:
  - `tencentcloud/services/teo/resource_tc_teo_zone_setting.go` — Add `cos_private_access` to the `origin` schema, read, and update logic
  - `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go` — Add unit tests for the new parameter
  - `tencentcloud/services/teo/resource_tc_teo_zone_setting.md` — Update documentation with the new parameter
- Cloud API: Uses existing `ModifyZoneSetting` and `DescribeZoneSetting` interfaces; no new API calls required
- Backward compatible: Adding an optional/computed field does not break existing configurations
