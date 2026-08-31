## Context

The `tencentcloud_teo_zone` resource (`resource_tc_teo_zone.go`) currently manages TEO zone lifecycle via the `CreateZone`, `DescribeZones`, `ModifyZone`, `ModifyZoneStatus`, `ModifyZoneWorkMode`, and `DeleteZone` APIs (all v20220901). The resource already has a functional schema with `zone_id`, `zone_name`, `type`, `alias_zone_name`, `area`, `plan_id`, `paused`, `status`, `ownership_verification` (partial — only `dns_verification`), `name_servers`, `tags`, and `work_mode_infos`.

The `CreateZone` API request struct (`CreateZoneRequest`) in the SDK exposes `AllowDuplicates *bool` and `JumpStart *bool` fields that are not currently mapped to any Terraform schema attribute. The `CreateZone` API response's `OwnershipVerification` struct contains `DnsVerification`, `FileVerification`, and `NsVerification` sub-structs, but the current Terraform schema only exposes `dns_verification`.

## Goals / Non-Goals

**Goals:**
- Add `allow_duplicates` (TypeBool, Optional, ForceNew) to the resource schema, passing it to `CreateZone`
- Add `jump_start` (TypeBool, Optional, ForceNew) to the resource schema, passing it to `CreateZone`
- Extend `ownership_verification` block with `file_verification` sub-block (`path`, `content`) and `ns_verification` sub-block (`name_servers` list), all Computed
- Populate new `ownership_verification` sub-blocks in both `Create` (from `CreateZone` response) and `Read` (from `DescribeZones` response)

**Non-Goals:**
- No changes to `ModifyZone`, `DeleteZone`, or `ModifyZoneStatus` logic beyond what already exists
- No new API calls beyond what's already used
- No changes to the `tags` handling or `work_mode_infos` handling
- No modification of existing parameter types or behaviors

## Decisions

### Decision 1: `allow_duplicates` and `jump_start` as ForceNew

These parameters control creation-time behavior and cannot be changed after a zone is created (the `ModifyZone` API does not accept them). Marking them as `ForceNew: true` ensures Terraform destroys and recreates the resource if the user changes these values, which is the correct behavior.

### Decision 2: `file_verification` and `ns_verification` as Computed-only sub-blocks

These fields are returned by the API as part of the `OwnershipVerification` response. They are informational/read-only for the user — they guide the verification process but cannot be set by the user. Marking them as `Computed: true` without `Optional` is appropriate.

**Alternative considered**: Omitting these fields entirely. Rejected because they provide useful information for users automating DNS/file-based verification steps.

### Decision 3: Schema structure for `ownership_verification`

The existing `ownership_verification` block contains `dns_verification` as a `TypeList` with `MaxItems: 1`. Following the same pattern:
- `file_verification`: TypeList, Computed, MaxItems: 1, containing `path` (TypeString, Computed) and `content` (TypeString, Computed)
- `ns_verification`: TypeList, Computed, MaxItems: 1, containing `name_servers` (TypeList of TypeString, Computed)

This matches the existing `dns_verification` pattern and the SDK struct definitions.

### Decision 4: Populate from both Create and Read paths

The `ownership_verification` data is returned in both `CreateZoneResponse` and `DescribeZones` response. The `Create` method already calls `resourceTencentCloudTeoZoneRead` at the end, so the Read path will handle population. However, to be safe, we also populate in the Create handler's post-processing.

## Risks / Trade-offs

- **Risk**: If the API returns `nil` for `FileVerification` or `NsVerification` (e.g., for NS access type where file verification is not applicable), the code must handle this gracefully. → **Mitigation**: Always check for `nil` before accessing sub-fields. Set empty list `[]interface{}{}` when nil, consistent with existing `dns_verification` handling.
- **Risk**: The `ownership_verification` block is also on the Zone struct (deprecated, historical fields). The `CNAMEDetail` and `NSDetail` structs also have their own `OwnershipVerification`. → **Mitigation**: Only populate from the top-level `Zone.OwnershipVerification` for now, matching the existing pattern. The detail-level verifications are out of scope.