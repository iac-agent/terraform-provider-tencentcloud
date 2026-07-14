## Context

The `tencentcloud_instance_types` data source combines two API calls:
1. `DescribeZoneInstanceConfigInfos` (CVM API) — queries available instance type configurations
2. `DescribeDiskConfigQuota` (CBS API) — queries disk configuration quotas for CBS filtering

Currently, when `cbs_filter` is provided, the data source derives `DescribeDiskConfigQuota` request parameters from the instance type query results and the nested `cbs_filter` block:
- `Zones` is derived from `availability_zone` (the instance type query zone)
- `Memory` is derived from `memory_size` (the instance type query memory)
- `DiskTypes` is derived from `cbs_filter.disk_types` (nested inside cbs_filter)

This indirect derivation means users cannot independently specify zones, memory sizes, or disk types for the disk quota query — they are constrained to use the same values as the instance type query.

The CBS `DescribeDiskConfigQuota` API (package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312`) already supports `DiskTypes` ([]*string), `Zones` ([]*string), and `Memory` (*uint64) as direct request parameters.

## Goals / Non-Goals

**Goals:**
- Add `disk_types`, `zones`, and `memory` as top-level optional input parameters to the data source schema
- These top-level parameters take priority over derived values when calling `DescribeDiskConfigQuota`
- Maintain backward compatibility — existing configurations using `cbs_filter` continue to work unchanged
- Users can now specify disk quota query parameters independently of instance type query parameters

**Non-Goals:**
- Removing or changing the existing `cbs_filter` nested block — it remains functional
- Changing the `DescribeZoneInstanceConfigInfos` API call behavior
- Adding new output/computed fields (covered by the existing `extend-instance-types-datasource-fields` change)

## Decisions

### Decision 1: Top-level parameters override derived values in DescribeDiskConfigQuota

**Rationale**: When a user explicitly provides `disk_types`, `zones`, or `memory` at the top level, these should override the corresponding derived values (from instance type results or `cbs_filter`). This gives users direct control without removing the existing automatic derivation behavior.

**Priority logic in `dataSourceTencentCloudInstanceTypesRead`**:
- For `DiskTypes`: top-level `disk_types` > `cbs_filter.disk_types` (current behavior unchanged if neither provided)
- For `Zones`: top-level `zones` > derived `availability_zone` from instance type results
- For `Memory`: top-level `memory` > derived `memory_size` from instance type results

**Alternative considered**: Making top-level parameters mutually exclusive with `cbs_filter` — rejected because it breaks backward compatibility and adds unnecessary complexity.

### Decision 2: Schema types match CBS API field types

- `disk_types`: TypeList of TypeString (matches CBS API `DiskTypes` []*string)
- `zones`: TypeList of TypeString (matches CBS API `Zones` []*string)
- `memory`: TypeInt (matches CBS API `Memory` *uint64)

**Alternative considered**: Using TypeString for `zones` (single zone) — rejected because the CBS API supports multiple zones and users may want to query disk quotas across multiple zones at once.

### Decision 3: These parameters only affect DescribeDiskConfigQuota call

The new top-level `disk_types`, `zones`, and `memory` parameters are ONLY used in the `DescribeDiskConfigQuota` API call within the CBS filter section. They do NOT affect the `DescribeZoneInstanceConfigInfos` API call or the instance type filtering logic. The existing `memory_size` and `availability_zone` parameters continue to serve their original purpose for instance type filtering.

## Risks / Trade-offs

- [Parameter naming confusion] → `memory` (new top-level) vs `memory_size` (existing top-level): `memory_size` is used for instance type filtering, while `memory` specifically controls the CBS disk quota query. The distinction is intentional to avoid breaking changes and to clarify that they serve different API calls. Documentation should clearly explain this difference.

- [Backward compatibility] → No risk: all new parameters are optional. When not specified, the existing derivation logic from `cbs_filter` and instance type results continues to work identically.

- [Duplicate parameter paths] → Users might specify `disk_types` both at top-level and inside `cbs_filter.disk_types`. The design handles this by prioritizing the top-level value, which is the most explicit user intent.
