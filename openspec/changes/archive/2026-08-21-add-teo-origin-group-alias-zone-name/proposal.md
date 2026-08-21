## Why

The TEO cloud API `OriginGroupReference` struct exposes an `AliasZoneName` field that carries the alias zone name of each reference entry in an origin group. The Terraform resource `tencentcloud_teo_origin_group` currently exposes `instance_type`, `instance_id`, and `instance_name` in the `references` computed block, but does not expose `alias_zone_name`, so users cannot observe which alias zone a reference belongs to — especially in cross-zone reference scenarios.

## What Changes

- Add a single new computed attribute `alias_zone_name` to the `references` nested block of the `tencentcloud_teo_origin_group` resource, sourced from `OriginGroupReference.AliasZoneName` in the `DescribeOriginGroup` API response.
- No changes are required to Create/Update/Delete operations: the new field is read-only and computed, populated only during the Read operation.
- Update the resource unit test and the `.md` example to reflect the new computed field.

## Capabilities

### New Capabilities

- `teo-origin-group-alias-zone-name`: Add the `alias_zone_name` computed attribute to the `references` block of the `tencentcloud_teo_origin_group` resource.

### Modified Capabilities

## Impact

- Code:
  - `tencentcloud/services/teo/resource_tc_teo_origin_group.go` (add `alias_zone_name` schema field in the `references` block and the corresponding Read logic)
  - `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go` (add gomonkey-based unit tests for the new computed field)
  - `tencentcloud/services/teo/resource_tc_teo_origin_group.md` (update the example)
- Dependencies: uses the already vendored `tencentcloud-sdk-go` `teo/v20220901` package (`OriginGroupReference.AliasZoneName`), no vendor changes required.
- Backward compatibility: purely additive computed field; existing configuration and state remain valid.
- Documentation: generated from the `.md` file via `make doc` (no manual `website/` edits).
