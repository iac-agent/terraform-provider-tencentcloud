## Why

The `tencentcloud_teo_zone_setting` resource currently lacks the `cos_private_access` field within the `origin` block. Users cannot configure private access for COS origin buckets through Terraform, even though the TEO cloud API (`ModifyZoneSetting` / `DescribeZoneSetting`) already supports this parameter. Adding this field enables users to manage COS private access settings as part of their zone configuration.

## What Changes

- Add `cos_private_access` field (type: string, optional + computed) to the `origin` block in the `tencentcloud_teo_zone_setting` resource schema
- Update the resource Read function to map `Origin.CosPrivateAccess` from the API response to the Terraform state
- Update the resource Update function to send `cos_private_access` value via `Origin.CosPrivateAccess` in the `ModifyZoneSetting` request
- Update the `.md` documentation file for the resource

## Capabilities

### New Capabilities

- `teo-zone-setting-cos-private-access`: Adds `cos_private_access` parameter to the `origin` block of `tencentcloud_teo_zone_setting`, allowing Terraform to manage COS private access configuration for TEO zone settings.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_zone_setting.go` — schema definition, read, and update functions
- `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go` — unit tests for the new parameter
- `tencentcloud/services/teo/resource_tc_teo_zone_setting.md` — documentation update
