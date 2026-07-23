## Why

The `tencentcloud_teo_l7_acc_setting` resource currently has `zone_name` as a `Computed` (read-only) field only. Users cannot supply `zone_name` in their Terraform configuration, which limits flexibility. The Cloud API `DescribeL7AccSetting` response already returns `ZoneName` in the `ZoneConfigParameters` struct, so making this field optionally inputable allows users to document the zone name in their configuration while the provider still reads the authoritative value from the API.

## What Changes

- Modify `zone_name` field schema from `Computed: true` to `Optional: true, Computed: true` in `tencentcloud_teo_l7_acc_setting` resource

## Capabilities

### New Capabilities
- `teo-l7-acc-setting-zone-name-input`: Allow users to optionally specify `zone_name` in the `tencentcloud_teo_l7_acc_setting` resource configuration. The field remains `Computed` so the authoritative value still comes from the cloud API.

### Modified Capabilities
<!-- None - no existing spec requirement changes -->

## Impact

- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting.go`
- Only the Schema definition of `zone_name` field is modified (adding `Optional: true`)
- No API call changes required - the field remains computed from the API response
- Backward compatible: existing Terraform configurations continue to work unchanged