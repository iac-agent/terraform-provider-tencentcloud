## Context

The `tencentcloud_teo_zone_setting` resource manages EdgeOne zone settings. The `origin` block already exists in the schema with three fields (`origins`, `backup_origins`, `origin_pull_protocol`), but the cloud API's `Origin` struct also includes a `CosPrivateAccess` field that is not yet exposed in Terraform. The cloud SDK (`teo/v20220901`) already contains this field in both the request (`ModifyZoneSettingRequest.Origin.CosPrivateAccess`) and response (`ZoneSetting.Origin.CosPrivateAccess`).

## Goals / Non-Goals

**Goals:**
- Add `cos_private_access` field to the `origin` block of `tencentcloud_teo_zone_setting`
- Maintain backward compatibility — existing TF configurations without `cos_private_access` must continue to work
- Support both read and write (create/update) paths for the new field

**Non-Goals:**
- No changes to the `origin` block's existing fields (`origins`, `backup_origins`, `origin_pull_protocol`)
- No changes to other blocks in the resource
- No new resource or data source creation

## Decisions

1. **Field type: `schema.TypeString` with `Optional: true, Computed: true`**
   - The cloud API defines `CosPrivateAccess` as `*string` with values `"on"` or `"off"`.
   - `Computed: true` ensures backward compatibility — if the user doesn't specify it, the value is read from the API response.

2. **No `ValidateFunc` for value constraints**
   - Following the existing pattern in this resource (e.g., `origin_pull_protocol` also has valid values but no validation function), we keep consistent and do not add validation.

3. **Read path: Map `respData.Origin.CosPrivateAccess` to `originMap["cos_private_access"]`**
   - Follow the existing nil-check pattern used for other `Origin` fields.

4. **Update path: Read `cos_private_access` from `originMap` and set `origin.CosPrivateAccess`**
   - Follow the existing pattern for `origin_pull_protocol` — read from the map and set via `helper.String()`.

## Risks / Trade-offs

- [Backward compatibility] → Mitigated by using `Optional: true, Computed: true`. Existing configurations that don't specify `cos_private_access` will not be affected.
- [API default behavior] → If the API returns `nil` for `CosPrivateAccess`, the field won't be set in the state, which is consistent with the existing handling of other optional fields in this resource.
