## Why

The cloud API `DescribeOriginGroup` returns zone-level reference information through the `OriginGroupReference.ZoneName` field (JsonPath: `response.OriginGroups.References.ZoneName`), which provides the zone name of each referenced instance for the TEO origin group. Exposing this as a computed attribute in the `tencentcloud_teo_origin_group` resource gives users visibility into which zone each reference belongs to, which is essential for cross-zone reference scenarios.

## What Changes

- Add 1 new computed attribute to the `references` nested block of the `tencentcloud_teo_origin_group` resource:
  - `zone_name`: The zone name of the referenced instance (maps to `OriginGroupReference.ZoneName`)

This is a read-only computed field sourced from the `DescribeOriginGroup` API response, requiring no changes to Create/Update/Delete operations.

## Capabilities

### New Capabilities
- `teo-origin-group-reference-zone-name`: Add `zone_name` computed attribute to the `references` block of `tencentcloud_teo_origin_group` resource, sourced from `DescribeOriginGroup` API response's `OriginGroupReference.ZoneName`.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_origin_group.go`: Add the `zone_name` schema field in the `references` block and the corresponding read logic.
- `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go`: Add unit tests for the new computed field.
- `tencentcloud/services/teo/resource_tc_teo_origin_group.md`: Update example documentation.
- No breaking changes - the new field is computed and backward compatible.
