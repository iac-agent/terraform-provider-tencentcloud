## Context

The `tencentcloud_teo_zone` resource (RESOURCE_KIND_GENERAL) uses `DescribeZones` API via the service layer method `DescribeTeoZoneById` to fetch zone details during the Read operation. Currently, the service layer hardcodes `Offset=0` and `Limit=20` for pagination, and the `TotalCount` from the API response is discarded. The vendor SDK (`v20220901`) already supports `Offset`, `Limit` in `DescribeZonesRequest` and `TotalCount` in `DescribeZonesResponseParams`.

## Goals / Non-Goals

**Goals:**
- Expose `offset` (Optional, TypeInt) and `limit` (Optional, TypeInt) as Terraform schema parameters for the `tencentcloud_teo_zone` resource
- Expose `total_count` (Computed, TypeInt) as a Terraform schema parameter to surface the `DescribeZones` response `TotalCount`
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Do not expose Offset/Limit to the `DescribeTeoZoneById` service method signature (it's used internally and modifying its signature would be a larger refactor)
- Do not add these parameters to any other TEO resource or data source
- Do not change the pagination logic in the service layer

## Decisions

### Decision 1: Add parameters to resource schema only, not to service layer

The `offset` and `limit` parameters are added to the resource schema as Optional fields. The `total_count` is added as Computed. The service layer `DescribeTeoZoneById` method will continue to use its internal pagination logic (hardcoded Offset=0, Limit=20). The new parameters serve as passthrough/documentation for the API contract, and `total_count` is set from the first API response.

**Rationale**: Modifying the service layer `DescribeTeoZoneById` method signature would require changes across all callers. The resource currently uses this method for Read, and the API response already contains `TotalCount`. We can capture `TotalCount` from the response without changing the service method.

**Alternative considered**: Modify `DescribeTeoZoneById` to accept and return pagination parameters. Rejected because it's unnecessary complexity — the service method is designed for internal use with a single zone lookup by ID, where pagination is handled internally.

### Decision 2: All three parameters in a single change

All three parameters (`Offset`, `Limit`, `TotalCount`) are grouped in one change because they are all part of the same `DescribeZones` API contract and are logically related.

**Rationale**: They form a cohesive unit — input pagination controls and output count. Splitting them would create unnecessary fragmentation.

## Risks / Trade-offs

- **Risk**: Users may set `offset` and `limit` expecting them to control the Read behavior, but the internal service method `DescribeTeoZoneById` does not currently use the user-provided values. → **Mitigation**: The parameters are documented as informational/passthrough. Future enhancement could wire them into the service method if needed.
- **Risk**: `total_count` may be stale if the Read operation uses a cached response. → **Mitigation**: This is an inherent limitation of the resource model — the Read operation always fetches the latest state from the API.