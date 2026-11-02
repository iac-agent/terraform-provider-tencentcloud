## Context

The `tencentcloud_teo_origin_group` resource manages TEO origin groups. The `references` computed block exposes the list of instances that reference the origin group, sourced from the `OriginGroupReference` entries returned by the `DescribeOriginGroup` API. The cloud API SDK struct `OriginGroupReference` (in package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) contains a `ZoneName` field (JsonPath: `response.OriginGroups.References.ZoneName`) that provides the zone name of each referenced instance. This field is important for cross-zone reference scenarios where origin groups are referenced by resources in different zones.

Current state of the `references` block schema:
- `instance_type` (Computed, string)
- `instance_id` (Computed, string)
- `instance_name` (Computed, string)

Cloud API `OriginGroupReference` struct relevant field:
- `ZoneName` (`*string`) - 引用站点名称 (the zone name of the referenced instance) - to be added as `zone_name`

## Goals / Non-Goals

**Goals:**
- Add `zone_name` computed attribute to the `references` nested block of `tencentcloud_teo_origin_group`.
- Ensure the Read method populates `zone_name` from the `DescribeOriginGroup` API response's `OriginGroupReference.ZoneName`.
- Maintain full backward compatibility - this is an additive computed field.

**Non-Goals:**
- Modifying Create, Update, or Delete operations (this field is read-only from the API).
- Changing the existing `references` sub-attributes (`instance_type`, `instance_id`, `instance_name`) or their behavior.
- Adding any new top-level resource parameters.

## Decisions

1. **Schema type for the new field**: Use `schema.TypeString` with `Computed: true` for `zone_name`. This matches the existing pattern in the `references` block and the cloud API field type (`*string`).

2. **Read logic placement**: Add the field mapping inside the existing `references` loop in `resourceTencentCloudTeoOriginGroupRead`, following the same nil-check pattern as the existing fields (`if references.ZoneName != nil { referencesMap["zone_name"] = references.ZoneName }`).

3. **No changes to mutableArgs**: Since `zone_name` is a computed-only field in the `references` block (not a top-level mutable field), no changes are needed to the `mutableArgs` list in the Update method.

4. **Test approach**: Use gomonkey-based mock unit tests (not Terraform acceptance tests) since this is a modification to an existing resource. Mock the `DescribeOriginGroup` client method (via the service layer) to return test data with the `ZoneName` field populated, and assert the state reflects it.

## Risks / Trade-offs

- [Risk] The cloud API may return nil for `ZoneName` → Mitigation: Follow the existing nil-check pattern before setting the field value, so nil values are skipped.
- [Risk] State migration is not needed → Mitigation: Adding a computed field is backward compatible; existing state files will simply have the field empty until the next refresh.
