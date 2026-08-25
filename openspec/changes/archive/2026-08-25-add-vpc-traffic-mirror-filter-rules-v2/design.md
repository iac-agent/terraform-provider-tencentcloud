## Context

TencentCloud VPC traffic mirror introduced a new "NO-DIRECTION" mode (`Direction` field on the traffic mirror instance) that deprecates the old direction-based filter rules. In this mode, filter rules are managed through dedicated APIs (`CreateTrafficMirrorFilterRules`, `ModifyTrafficMirrorFilterRules`, `DeleteTrafficMirrorFilterRules`, `DescribeTrafficMirrorFilterRules`) that support rule-level priority, individual editing, and explicit direction (INGRESS/EGRESS) per rule.

The Terraform provider currently has no resource for these new filter rule APIs. This design covers adding a new `tencentcloud_vpc_traffic_mirror_filter_rules_v2` resource that manages all filter rules under a given traffic mirror instance.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_vpc_traffic_mirror_filter_rules_v2` with full CRUD support
- Support both ingress and egress filter rules as TypeList schemas
- Use the existing `tencentcloud-sdk-go/vpc/v20170312` SDK package
- Follow existing provider patterns (retry, error handling, composite IDs)

**Non-Goals:**
- Data source for filter rules (separate change)
- Modifying the existing `tencentcloud_vpc_traffic_mirror` resource
- Supporting the old direction-based filter rules (covered by existing resource)

## Decisions

### 1. Resource ID = `traffic_mirror_id`

The resource manages all filter rules under a single traffic mirror instance. Using `traffic_mirror_id` as the Terraform resource ID is the simplest and most natural approach. The Create/Read/Update/Delete APIs all operate on rules scoped to a traffic mirror instance.

**Alternatives considered**: Composite ID with rule IDs concatenated. Rejected because the API manages rules as a batch; individual rule IDs are returned by the API and can change with each update. Using `traffic_mirror_id` avoids complex diff logic.

### 2. Schema: TypeList for ingress/egress_filter_rules

Each filter rule is a nested `TypeList` with `MaxItems: 1` concept flattened into individual rule objects. The `TrafficMirrorFilter` struct from the SDK maps directly to the schema fields:

- `src_net` (TypeString, Optional)
- `dst_net` (TypeString, Optional)
- `protocol` (TypeString, Optional)
- `src_port` (TypeString, Optional)
- `dst_port` (TypeString, Optional)
- `traffic_mirror_filter_rule_id` (TypeString, Computed)
- `priority` (TypeInt, Optional) — SDK uses `*uint64`
- `action` (TypeString, Optional)
- `description` (TypeString, Optional)
- `created_time` (TypeString, Computed)

`traffic_mirror_filter_rule_id` and `created_time` are Computed because they are returned by the API after creation, not specified by the user.

### 3. CRUD operations use retry with ReadRetryTimeout

All API calls use `tccommon.ReadRetryTimeout` with `helper.Retry()` for eventual consistency. On API errors, wrap with `tccommon.RetryError()`.

### 4. Update replaces all rules (batch update)

`ModifyTrafficMirrorFilterRules` takes the full set of rules. The Terraform Update function will send the complete desired state from the configuration. This is simpler than computing diffs and aligns with the API design.

### 5. Delete sends all rule IDs

`DeleteTrafficMirrorFilterRules` requires `IngressFilterRuleIds` and `EgressFilterRuleIds` (string slices of rule IDs). The Delete function reads current state from `d.Id()` to get rule IDs, then sends them in the delete request.

## Risks / Trade-offs

- **API consistency**: The Create/Modify APIs return the full rule set with IDs, but the response may not include `CreatedTime` consistently. The Read function should handle nil fields gracefully.
- **Batch replace semantics**: Modifying a single rule requires sending the entire rule set. This is a limitation of the API design, not the Terraform resource. Documented in the resource description.
- **Priority type**: The SDK uses `*uint64` for `Priority`, while most other fields are `*string`. The schema must use `TypeInt` to match.