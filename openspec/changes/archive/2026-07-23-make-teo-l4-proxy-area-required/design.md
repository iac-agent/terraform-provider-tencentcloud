## Context

The `tencentcloud_teo_l4_proxy` resource currently defines the `area` field as `Optional: true`. The `area` field specifies the acceleration zone for the L4 proxy instance (`mainland`, `overseas`, or `global`). This is a critical configuration parameter that must be provided at creation time and cannot be changed afterwards (it is already in the `immutableArgs` list in the Update function).

The cloud API's `CreateL4Proxy` endpoint accepts `Area` as a parameter, and in practice, users must always specify it. The current optional nature causes confusion and potential errors when the field is omitted.

## Goals / Non-Goals

**Goals:**
- Change `area` from `Optional: true` to `Required: true` in the resource schema
- Update the Create function to use direct field access since the field is now guaranteed to be present
- Update the `.md` documentation to reflect the change (though the `make doc` in finalize phase will auto-generate this)

**Non-Goals:**
- No changes to the Read, Update, or Delete functions (they already handle `area` correctly)
- No changes to the API layer or SDK
- No changes to other resources

## Decisions

### Decision 1: Change Schema from Optional to Required

**Choice**: Set `Required: true` (remove `Optional: true`) for the `area` field.

**Rationale**: The `area` field is essential for creating an L4 proxy. The cloud API documentation lists it as a key parameter. Making it required ensures users always provide it, preventing runtime errors and improving the user experience.

**Alternatives considered**:
- Keep it optional with `Default` value: Rejected because there is no sensible default - the acceleration zone choice depends on the user's specific deployment needs.
- Keep it as-is: Rejected because it leads to poor user experience.

### Decision 2: Update Create Function to Remove GetOk Guard

**Choice**: Change `d.GetOk("area")` to `d.Get("area")` in the Create function.

**Rationale**: When a field is `Required`, Terraform guarantees it will always have a value. Using `GetOk` is unnecessary and inconsistent with how other `Required` fields are handled (e.g., `zone_id` and `proxy_name`).

### Decision 3: Area Remains Immutable

**Choice**: Keep `area` in the `immutableArgs` list in the Update function.

**Rationale**: The `area` field cannot be changed after creation - this is already enforced and consistent with the cloud API behavior. The `ModifyL4Proxy` API does not accept `Area` as a parameter.

## Risks / Trade-offs

- **Breaking Change**: Existing `.tf` configurations that omit `area` will fail validation. This is an intentional breaking change that improves correctness. Users must update their configurations to include `area`.
- **State Migration**: No state migration is needed. Existing resources in state already have `area` populated from the Read function. The change only affects new `plan`/`apply` operations.