## Context

The `tencentcloud_teo_zone_setting` resource manages site-level acceleration settings for TencentCloud EdgeOne (TEO). The resource's `origin` block currently supports `origins`, `backup_origins`, and `origin_pull_protocol` parameters. The cloud API's `Origin` struct also exposes a `CosPrivateAccess` field that controls whether the origin accesses a COS bucket privately, but this field is not yet exposed in the Terraform resource.

The `CosPrivateAccess` field is present in both the `ModifyZoneSetting` and `DescribeZoneSetting` API interfaces within the `Origin` struct, making it straightforward to add to the existing resource.

## Goals / Non-Goals

**Goals:**
- Add `cos_private_access` parameter to the `origin` nested block of `tencentcloud_teo_zone_setting`
- Ensure the parameter is properly read from `DescribeZoneSetting` and sent in `ModifyZoneSetting`
- Maintain backward compatibility by making the parameter optional and computed

**Non-Goals:**
- Adding any other new parameters to this resource
- Modifying the structure or behavior of existing parameters
- Creating new resources or data sources

## Decisions

1. **Parameter placement**: Add `cos_private_access` as a field within the existing `origin` TypeList block, since the cloud API maps `CosPrivateAccess` to the `Origin` struct. This follows the existing pattern where all origin-related fields are grouped under the `origin` block.

2. **Schema configuration**: Use `Optional: true, Computed: true` to maintain backward compatibility. Existing configurations without this parameter will continue to work, and the value will be populated from the API response during read.

3. **Type**: Use `schema.TypeString` since the cloud API field is `*string` with values `on` and `off`.

4. **Read logic**: In the `resourceTencentCloudTeoZoneSettingRead` function, add nil-check for `respData.Origin.CosPrivateAccess` before setting the value, consistent with the existing pattern for other origin fields.

5. **Update logic**: In the `resourceTencentCloudTeoZoneSettingUpdate` function, add `cos_private_access` mapping from the Terraform schema to the `teo.Origin` struct's `CosPrivateAccess` field, consistent with how `origin_pull_protocol` is handled.

## Risks / Trade-offs

- [Risk] Existing users may see a diff on next `terraform plan` if their COS origin is configured with a non-default `CosPrivateAccess` value that wasn't previously managed → Mitigation: Using `Computed: true` ensures the value is read from the API and stored in state, preventing unexpected diffs.
- [Risk] The parameter is only meaningful when the origin is a COS bucket → Mitigation: The cloud API handles this validation; the Terraform parameter simply passes the value through.
