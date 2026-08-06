## Context

The `tencentcloud_teo_zone` resource's Read method calls `DescribeTeoZoneById`, which internally uses the `DescribeZones` API to look up a zone by its ID. Currently, the `DescribeTeoZoneById` method hardcodes pagination parameters (`offset=0`, `limit=20`) and internally loops through all pages. The `DescribeZones` API supports `Offset` and `Limit` parameters for pagination control, but these are not exposed to Terraform users.

The `DescribeZones` API's `Offset` and `Limit` fields are defined in the vendor SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go` as `*int64` typed fields on `DescribeZonesRequest`.

## Goals / Non-Goals

**Goals:**
- Add `offset` (Optional, TypeInt) parameter to the `tencentcloud_teo_zone` resource schema
- Add `limit` (Optional, TypeInt) parameter to the `tencentcloud_teo_zone` resource schema
- Pass these parameters to the `DescribeTeoZoneById` service method, which forwards them to the `DescribeZones` API call
- Maintain full backward compatibility

**Non-Goals:**
- Changing the internal pagination loop logic (the loop still works as before when offset/limit are not specified)
- Adding these parameters to the datasource `tencentcloud_teo_zones`
- Modifying Create, Update, or Delete operations

## Decisions

### Decision 1: Add parameters to the resource schema, not the datasource

**Rationale**: The user requirement specifies adding these to the `tencentcloud_teo_zone` resource. The datasource `tencentcloud_teo_zones` already has its own internal pagination logic that properly handles all pages.

### Decision 2: Modify `DescribeTeoZoneById` signature to accept offset/limit

**Rationale**: The cleanest approach is to add `offset` and `limit` parameters to the `DescribeTeoZoneById` service method signature. When the user provides these values, the method uses them instead of the hardcoded defaults. The method signature changes from:
```go
func (me *TeoService) DescribeTeoZoneById(ctx context.Context, zoneId string) (ret *teo.Zone, errRet error)
```
to:
```go
func (me *TeoService) DescribeTeoZoneById(ctx context.Context, zoneId string, offset int64, limit int64) (ret *teo.Zone, errRet error)
```

When `offset` and `limit` are non-zero, the method uses them directly instead of the hardcoded values. When zero, the method falls back to the existing default behavior (`offset=0`, `limit=20`).

**Alternatives considered**:
- Passing values via context: More complex, harder to trace, and not idiomatic for this codebase
- Creating a separate method: Unnecessary duplication given the minimal change

### Decision 3: Schema fields are Optional with no default values

**Rationale**: Both `offset` and `limit` are `Optional` with `TypeInt`. When not set by the user, the service method uses its existing hardcoded defaults (`offset=0`, `limit=20`), ensuring backward compatibility.

## Risks / Trade-offs

- **Risk**: The `DescribeTeoZoneById` method signature change could break other callers. → **Mitigation**: Check all callers of `DescribeTeoZoneById` in the codebase and update them to pass `0, 0` (which triggers the default behavior).

- **Risk**: For a single-zone lookup (by ID), offset/limit pagination has limited practical use. → **Mitigation**: These parameters are optional; users who don't need them simply don't set them. The parameters provide API surface consistency.