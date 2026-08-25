## Context

The `tencentcloud_teo_origin_group` resource already has a `references` computed block that exposes `OriginGroupReference` fields from the `DescribeOriginGroup` API response. The existing fields include `instance_type`, `instance_id`, `instance_name`, `zone_id`, and `zone_name`. The `AliasZoneName` field is already present in the SDK's `OriginGroupReference` struct but is not yet exposed in the Terraform resource schema.

The SDK model at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go:23431` confirms `AliasZoneName *string` is a valid field returned by `DescribeOriginGroup`.

## Goals / Non-Goals

**Goals:**
- Add `alias_zone_name` as a computed-only string field to the `references` nested block of `tencentcloud_teo_origin_group`
- Populate the field from `OriginGroupReference.AliasZoneName` in the Read method
- Ensure nil-safe handling (only set when the API returns a non-nil value)

**Non-Goals:**
- No changes to Create, Update, or Delete operations
- No changes to any other resource or data source
- No breaking changes to existing configurations

## Decisions

1. **Field placement**: The new field goes inside the existing `references` nested block (not at the top level), matching the API structure where `AliasZoneName` is a field of `OriginGroupReference`.

2. **Computed-only**: The field is `Computed: true` only, with no `Optional` or `Required`. This is because the field is an output of the `DescribeOriginGroup` API and cannot be set by users.

3. **Nil-safe pattern**: The read logic follows the existing pattern of other reference fields — check `references.AliasZoneName != nil` before setting the map value. This prevents nil pointer dereference and ensures clean state when the API omits the field.

4. **No service layer changes**: The `DescribeTeoOriginGroupById` service method already returns the full `OriginGroup` struct including `References`, so no changes are needed in the service layer (`service_tencentcloud_teo.go`).

## Risks / Trade-offs

- **Risk**: API may return `AliasZoneName` as nil for some reference types. → **Mitigation**: Nil check before setting, matching existing pattern for `ZoneName` and `ZoneId`.
- **Risk**: Schema change could cause state migration issues. → **Mitigation**: New computed field with no existing state is backward compatible; Terraform ignores unknown computed fields in existing state.