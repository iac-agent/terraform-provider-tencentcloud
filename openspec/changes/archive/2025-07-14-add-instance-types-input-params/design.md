## Context

The `tencentcloud_instance_types` data source (`data_source_tc_instance_types.go`) currently supports two filtering mechanisms:

1. **Top-level parameters**: `cpu_core_count`, `gpu_core_count`, `memory_size`, `availability_zone`, `exclude_sold_out` for basic filtering
2. **`filter` block**: Maps to `Filters` parameter of `DescribeZoneInstanceConfigInfos` CVM API, supporting filter names like `zone`, `instance-family`, `instance-type`, `instance-charge-type`, `sort-keys`
3. **`cbs_filter` block**: Maps to `DescribeDiskConfigQuota` CBS API, supporting `disk_types`, `disk_charge_type`, `disk_usage`

The current data source read flow:
1. Calls `DescribeZoneInstanceConfigInfos` (CVM) via `DescribeInstancesSellTypeByFilter` to get instance type quota items
2. Filters results by `cpu_core_count`, `gpu_core_count`, `memory_size`, `exclude_sold_out` in-memory
3. If `cbs_filter` is provided, calls `DescribeDiskConfigQuota` (CBS) for each instance type result to get CBS config

### API Constraints

- **`DescribeZoneInstanceConfigInfos`** (CVM): Only has `Filters []*Filter` as request parameter. There is no direct `InstanceFamilies` field. Instance family filtering is done via a filter with `Name: "instance-family"`.
- **`DescribeDiskConfigQuota`** (CBS): Has direct `InstanceFamilies []*string` and `DiskTypes []*string` request fields. Currently, `InstanceFamilies` is auto-populated from each instance type's family, and `DiskTypes` is passed from `cbs_filter.disk_types`.

## Goals / Non-Goals

**Goals:**
- Add `InstanceFamilies` as a top-level Optional input parameter (TypeList of TypeString) to simplify instance family filtering
- Add `DiskTypes` as a top-level Optional input parameter (TypeList of TypeString) to simplify disk type filtering
- Maintain full backward compatibility with existing configurations

**Non-Goals:**
- Removing or deprecating the existing `filter` or `cbs_filter` blocks
- Changing the output schema or computed fields
- Modifying the `DescribeInstancesSellTypeByFilter` service function signature

## Decisions

### Decision 1: InstanceFamilies → Filter Translation for DescribeZoneInstanceConfigInfos

**Choice**: Translate `InstanceFamilies` values into `instance-family` filter values, appended to the existing filter map.

**Rationale**: The `DescribeZoneInstanceConfigInfos` API does not have a direct `InstanceFamilies` field. The only way to filter by instance family is through the `Filters` parameter with `Name: "instance-family"`. This is consistent with how `availability_zone` is already handled (translated to `zone` filter).

**Alternative considered**: Adding `InstanceFamilies` as a separate API call - rejected because it would add unnecessary complexity and duplicate results.

### Decision 2: InstanceFamilies → Direct Parameter for DescribeDiskConfigQuota

**Choice**: When `InstanceFamilies` is provided as a top-level parameter, pass it directly to the `DescribeDiskConfigQuota` API's `InstanceFamilies` field instead of auto-populating from each instance type's family.

**Rationale**: The `DescribeDiskConfigQuota` API supports `InstanceFamilies` as a direct request field. If the user provides `InstanceFamilies`, it should be used to filter the CBS disk config quota results. When `InstanceFamilies` is not provided, fall back to the existing behavior of auto-populating from each instance type's family.

### Decision 3: DiskTypes → Direct Parameter for DescribeDiskConfigQuota

**Choice**: When `DiskTypes` is provided as a top-level parameter, pass it directly to the `DescribeDiskConfigQuota` API's `DiskTypes` field. When `cbs_filter` is not provided but `DiskTypes` is, still trigger the CBS config query.

**Rationale**: The `DescribeDiskConfigQuota` API supports `DiskTypes` as a direct request field. Currently `disk_types` is only available inside `cbs_filter`. Making it a top-level parameter simplifies the most common use case.

### Decision 4: Conflict Handling with Existing Filter Mechanisms

**Choice**: `InstanceFamilies` top-level parameter SHALL NOT conflict with the `filter` block's `instance-family` filter. If both are provided, the `InstanceFamilies` values SHALL be merged with any existing `instance-family` filter values. Similarly, `DiskTypes` top-level parameter SHALL NOT conflict with `cbs_filter.disk_types`. If both are provided, the top-level `DiskTypes` SHALL take precedence.

**Rationale**: Avoiding conflicts simplifies the user experience. Merging for `InstanceFamilies` and taking precedence for `DiskTypes` provides predictable behavior.

## Risks / Trade-offs

- **[Risk] Duplicate filter specification**: Users might specify `InstanceFamilies` and also include `instance-family` in the `filter` block, leading to potentially confusing behavior → Mitigation: Merge values from both sources; document the behavior clearly.
- **[Risk] CBS query triggered without cbs_filter**: Adding `DiskTypes` as a top-level parameter could trigger CBS queries even without `cbs_filter` → Mitigation: Only trigger CBS query when `DiskTypes` is provided, consistent with existing `cbs_filter` trigger behavior.
