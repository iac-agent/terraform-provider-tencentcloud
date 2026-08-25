## Why

The cloud API `DescribeOriginGroup` returns zone-level reference information through the `OriginGroupReference.AliasZoneName` field (JsonPath: `response.OriginGroups.References.AliasZoneName`), which provides the alias zone name of each referenced instance for the TEO origin group. Exposing this as a computed attribute in the `tencentcloud_teo_origin_group` resource gives users visibility into the alias zone name of each reference, which is useful for cross-zone reference scenarios.

## What Changes

- Add 1 new computed attribute to the `references` nested block of the `tencentcloud_teo_origin_group` resource:
  - `alias_zone_name`: The alias zone name of the referenced instance (maps to `OriginGroupReference.AliasZoneName`)

This is a read-only computed field sourced from the `DescribeOriginGroup` API response, requiring no changes to Create/Update/Delete operations.

## Capabilities

### New Capabilities
- `teo-origin-group-reference-alias-zone-name`: Add `alias_zone_name` computed attribute to the `references` block of `tencentcloud_teo_origin_group` resource, sourced from `DescribeOriginGroup` API response's `OriginGroupReference.AliasZoneName`.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_origin_group.go`: Add the `alias_zone_name` schema field in the `references` block and the corresponding read logic.
- `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go`: Add unit tests for the new computed field.
- `tencentcloud/services/teo/resource_tc_teo_origin_group.md`: Update example documentation.
- No breaking changes - the new field is computed and backward compatible.