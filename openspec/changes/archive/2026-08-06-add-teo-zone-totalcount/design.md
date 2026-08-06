## Context

The `tencentcloud_teo_zone` resource's Read function calls `DescribeTeoZoneById` (service layer), which internally invokes the `DescribeZones` API. The `DescribeZonesResponseParams` includes a `TotalCount` field (`*int64`) representing the total number of zones matching the query filter. Currently this field is ignored. The change exposes it as a computed attribute on the resource.

## Goals / Non-Goals

**Goals:**
- Add a new computed `total_count` (TypeInt) field to the `tencentcloud_teo_zone` resource schema
- Populate `total_count` from `DescribeZonesResponseParams.TotalCount` during the Read operation
- Maintain full backward compatibility

**Non-Goals:**
- Do not add `total_count` to the Create/Update/Delete paths
- Do not modify the zone creation or update API calls
- Do not change any existing schema fields

## Decisions

### Decision 1: Modify `DescribeTeoZoneById` to return `TotalCount`

**Chosen**: Modify the service method signature to return `(totalCount int64, zone *teo.Zone, err error)`.

**Rationale**: The `DescribeZones` API returns `TotalCount` at the response level (not per-zone). The simplest approach is to capture it from the first API response page in the pagination loop and return it alongside the zone. This avoids duplicating the entire pagination logic in the Read function.

**Alternatives considered**:
- **Call `DescribeZones` directly in the Read function**: Rejected — would duplicate the pagination logic already in `DescribeTeoZoneById`.
- **Add a separate service method**: Rejected — overkill for a single field.

### Decision 2: Capture `TotalCount` from the first page response

**Chosen**: Read `TotalCount` from the first successful `DescribeZones` response within the pagination loop.

**Rationale**: `TotalCount` is consistent across all pages of the same query. The first page already contains the correct value.

### Decision 3: Schema field as `TypeInt` with `Computed: true`

**Chosen**: `"total_count": { Type: schema.TypeInt, Computed: true }`.

**Rationale**: `TotalCount` is `*int64` in the SDK. In Terraform, `TypeInt` is the correct mapping. The field is read-only (computed) since it is an API response value, not user-configurable.

## Risks / Trade-offs

- **Risk**: `TotalCount` may be nil in the API response → **Mitigation**: Only set `total_count` when `response.Response.TotalCount != nil`; otherwise skip setting the field.
- **Risk**: The `DescribeTeoZoneById` method is used by other callers → **Mitigation**: Check all callers and update them. Currently only `resourceTencentCloudTeoZoneRead` calls this method.