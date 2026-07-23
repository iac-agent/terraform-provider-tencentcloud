## Context

The `tencentcloud_teo_l7_acc_setting` resource manages L7 acceleration settings for a TEO zone. Currently, the `zone_name` field is `Computed: true` only, meaning it is populated from the cloud API response but cannot be specified by users in their Terraform configuration.

The Cloud API `DescribeL7AccSetting` returns `ZoneConfigParameters` which includes `ZoneName` — this is already used in the Read function to populate `zone_name`. The `ModifyL7AccSetting` API does not accept `ZoneName` (only `ZoneId` and `ZoneConfig`), so the field remains purely computed from the API side.

The change is simple: add `Optional: true` to the `zone_name` schema definition so users can optionally include it in their configuration. The field retains `Computed: true` so the authoritative value always comes from the API.

## Goals / Non-Goals

**Goals:**
- Allow users to optionally specify `zone_name` in the `tencentcloud_teo_l7_acc_setting` resource configuration
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- No changes to the API call logic (Create/Read/Update/Delete remain unchanged)
- No changes to the `zone_name` value source (still read from API response)
- No changes to the `.md` documentation file (the generated docs will reflect the schema change)

## Decisions

**Decision: Change `zone_name` from `Computed: true` to `Optional: true, Computed: true`**

Rationale: This is the standard Terraform pattern for a field that is optionally specified by the user but whose authoritative value comes from the API. Adding `Optional: true` does not break any existing configuration because all existing `Computed` fields continue to work as before.

**Alternative considered: Keep as-is**

Rejected because the user request specifically asks for this capability. The change is minimal and safe.

## Risks / Trade-offs

- **Risk**: Users may specify a `zone_name` that differs from the actual zone name → **Mitigation**: The field is `Computed: true`, so Terraform will always use the API-returned value, overwriting any user-specified value during refresh. This is standard Terraform behavior for `Optional` + `Computed` fields.
- **Risk**: Potential confusion about whether `zone_name` can be used to identify the zone → **Mitigation**: The `zone_id` field remains the primary identifier (`Required: true, ForceNew: true`). The description of `zone_name` already says "Zone name." which is clear enough.