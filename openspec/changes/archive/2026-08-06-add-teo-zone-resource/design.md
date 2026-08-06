## Context

TEO (Tencent EdgeOne) zones are the top-level resource in the EdgeOne platform, representing a site that can be accessed via CNAME, NS, or no-domain modes. The `tencentcloud_teo_zone` resource has already been implemented in the codebase (`resource_tc_teo_zone.go`, `resource_tc_teo_zone_extension.go`), with supporting service layer methods in `service_tencentcloud_teo.go`. This design document formalizes the architecture decisions made during implementation.

The resource follows the standard Terraform Plugin SDK v2 patterns used throughout the tencentcloud provider, leveraging the `tencentcloud-sdk-go` v20220901 TEO client for all cloud API calls.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource for full lifecycle management of TEO zones
- Support all zone access types: CNAME (`partial`), NS (`full`), no-domain (`noDomainAccess`), DNSPod (`dnsPodAccess`), and AI (`ai`)
- Support zone configuration updates: type, area, alias zone name, vanity name servers
- Support zone pause/resume via `ModifyZoneStatus` API
- Support configuration group work mode management via `ModifyZoneWorkMode` API
- Support tag management via the standard TencentCloud tag service

**Non-Goals:**
- Zone settings management (handled by separate `tencentcloud_teo_zone_setting` resource)
- DNS record management (handled by separate `tencentcloud_teo_dns_record` resources)
- Acceleration domain management (handled by separate `tencentcloud_teo_acceleration_domain` resource)

## Decisions

### 1. Zone ID as Terraform Resource ID

**Decision**: Use `ZoneId` from `CreateZone` response as the Terraform resource ID.

**Rationale**: `ZoneId` is the unique identifier for TEO zones across all cloud APIs. It follows the format `zone-xxxxxxxx` and is used consistently in `DescribeZones`, `ModifyZone`, `DeleteZone`, and all other zone-related APIs. Using a single identifier simplifies the resource lifecycle and enables direct `ImportStatePassthrough`.

### 2. DescribeZones Filtering for Read Operations

**Decision**: Use `DescribeZones` API with `zone-id` filter for reading individual zone details, rather than relying on a hypothetical `DescribeZone` (singular) API.

**Rationale**: The TEO API only provides `DescribeZones` (plural) which returns a paginated list. The `DescribeTeoZone` and `DescribeTeoZoneById` service methods filter by `zone-id` to retrieve exactly one zone. This is the standard pattern used across other TEO resources.

### 3. Post-Create Polling for Zone Readiness

**Decision**: After creating a zone, poll `DescribeTeoZone` until the zone status transitions from `pending` to a stable state.

**Rationale**: Zone creation is asynchronous. The `CreateZone` API returns immediately, but the zone may be in `pending` status initially. The extension function `resourceTencentCloudTeoZoneCreatePostHandleResponse0` polls with `6 * ReadRetryTimeout` to wait for the zone to become active.

### 4. Separate Update Paths for Different Concerns

**Decision**: Split the update logic into three separate code paths:
1. `ModifyZone` for type, alias_zone_name, area changes
2. `ModifyZoneStatus` for paused state changes
3. `ModifyZoneWorkMode` for work_mode_infos changes

**Rationale**: The TEO API separates these concerns into different endpoints. Grouping them by API call avoids unnecessary API requests and follows the principle of only calling what changed.

### 5. Pre-Delete Pause Check

**Decision**: Before deleting a zone, check if it is paused. If not paused, call `ModifyZoneStatus` to pause it first.

**Rationale**: The TEO `DeleteZone` API requires the zone to be in a paused state before deletion. The extension function `resourceTencentCloudTeoZoneDeletePostFillRequest0` handles this automatically.

### 6. Tags via Standard Tag Service

**Decision**: Use the common `svctag` package (`tencentcloud/services/tag`) for tag management, consistent with all other resources in the provider.

**Rationale**: This ensures consistent tag behavior across all TencentCloud resources and leverages the existing `ModifyTags` / `DescribeResourceTags` infrastructure.

### 7. Plan ID as ForceNew

**Decision**: `plan_id` is marked as `ForceNew: true` in the schema.

**Rationale**: The plan binding is set at creation time via `CreateZone` API. Changing the plan for an existing zone requires a different API (`BindZoneToPlan`) and represents a fundamentally different operation. Marking it as `ForceNew` ensures Terraform recreates the resource when the plan changes.

## Risks / Trade-offs

- **[Risk] Zone deletion is irreversible** → Mitigation: The resource requires the zone to be paused before deletion, providing a safety check. Users should be aware that deleted zones cannot be recovered.
- **[Risk] Async zone creation may timeout** → Mitigation: The polling uses `6 * ReadRetryTimeout` (a generous timeout). If the zone takes longer than this to become active, the creation will fail with a clear error message.
- **[Trade-off] `zone_name` is ForceNew** → The zone name (domain) is fundamental to the zone identity. Changing it would effectively create a new zone. This is standard Terraform behavior for resources where the name is part of the identity.