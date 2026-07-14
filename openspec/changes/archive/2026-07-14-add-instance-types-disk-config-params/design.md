## Context

The `tencentcloud_instance_types` data source queries two APIs:
1. **DescribeZoneInstanceConfigInfos** (CVM SDK) - Returns instance type configurations with CPU, memory, GPU, zone, family, status, etc.
2. **DescribeDiskConfigQuota** (CBS SDK) - Returns CBS disk configurations compatible with specific instance types. This API is only called when the `cbs_filter` block is provided.

Currently, when `DescribeDiskConfigQuota` is called, the `Zones` and `Memory` parameters are derived from each instance type's `availability_zone` and `memory_size` from the `DescribeZoneInstanceConfigInfos` response. The `DiskTypes` parameter is sourced from `cbs_filter.disk_types`. Users cannot override these values independently.

The `DescribeDiskConfigQuota` API (in CBS SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312`) accepts the following request fields:
- `InquiryType` (string) - Currently hardcoded to `INQUIRY_CVM_CONFIG`
- `DiskChargeType` (string) - From `cbs_filter.disk_charge_type`
- `InstanceFamilies` ([]string) - Derived from instance type's `family`
- `DiskTypes` ([]string) - Currently only from `cbs_filter.disk_types`
- `Zones` ([]string) - Currently derived from instance type's `availability_zone`
- `Memory` (uint64) - Currently derived from instance type's `memory_size`
- `DiskUsage` (string) - From `cbs_filter.disk_usage`
- `CPU` (uint64) - Derived from instance type's `cpu_core_count`

## Goals / Non-Goals

**Goals:**
- Add `disk_types`, `zones`, and `memory` as top-level optional input parameters to the `tencentcloud_instance_types` data source
- When these new parameters are provided, they override the values derived from instance type attributes in the `DescribeDiskConfigQuota` API call
- Maintain full backward compatibility - existing configurations without these new parameters continue to work identically

**Non-Goals:**
- Adding other `DescribeDiskConfigQuota` parameters (CPU, InstanceFamilies, DedicatedClusterId) as top-level parameters - these are not part of the requirement
- Changing the `DescribeZoneInstanceConfigInfos` API call or its filtering behavior
- Changing the `cbs_filter` block structure or its existing parameters

## Decisions

### Decision 1: Parameter placement - top-level vs. inside cbs_filter

**Choice**: Add `disk_types`, `zones`, and `memory` as **top-level** parameters in the data source schema.

**Rationale**: The `cbs_filter` block is a nested structure focused on CBS-specific filtering (disk_charge_type, disk_usage). The `zones` and `memory` parameters are already top-level concepts in the data source (there's already `availability_zone` and `memory_size` at the top level). Adding them as top-level parameters provides a cleaner user experience and allows these values to influence the `DescribeDiskConfigQuota` call independently of the `cbs_filter` block.

**Alternative considered**: Placing all three inside `cbs_filter` - rejected because `zones` and `memory` are not CBS-specific concepts and would create confusion with the existing `availability_zone` and `memory_size` top-level parameters.

### Decision 2: Parameter naming

**Choice**: Use `disk_types`, `zones`, and `memory` as the schema names.

**Rationale**: These names directly correspond to the `DescribeDiskConfigQuota` API's `DiskTypes`, `Zones`, and `Memory` request fields, following the Terraform convention of using snake_case for schema names. Note that `zones` is a list (allowing multiple zones) which differs from the existing `availability_zone` (single string). `memory` is an integer matching the API's `Memory` field (uint64).

### Decision 3: Override vs. merge behavior

**Choice**: When a new top-level parameter is provided, it **overrides** the value derived from instance type attributes in the `DescribeDiskConfigQuota` call. When not provided, the existing behavior (deriving from instance type attributes) continues.

**Rationale**: Override behavior is simpler and more predictable. If a user explicitly specifies `zones = ["ap-guangzhou-1", "ap-guangzhou-2"]`, they want the CBS query to use those zones, not the zone from each instance type result. The merge/combination approach would create ambiguity about which values take precedence.

### Decision 4: CBS service method changes

**Choice**: Modify `CbsService.DescribeDiskConfigQuota` to accept optional override values for `DiskTypes`, `Zones`, and `Memory` in the `cvmInfo` map, using the top-level parameter values when present and falling back to instance-type-derived values when absent.

**Rationale**: The existing method signature uses `map[string]interface{}` for `cvmInfo`, which naturally supports optional key-value pairs. Adding new keys (`"disk_types_override"`, `"zones_override"`, `"memory_override"`) to the map avoids changing the method signature and keeps the change minimal.

## Risks / Trade-offs

- **[Risk] Parameter name collision with existing fields** → The existing `memory_size` (TypeInt) is a different parameter from the new `memory` (TypeInt). While `memory_size` is used for filtering `DescribeZoneInstanceConfigInfos` results, `memory` is used for overriding the `DescribeDiskConfigQuota` query. Documentation must clearly distinguish these. → Mitigation: Clear documentation explaining that `memory_size` filters instance types while `memory` overrides the CBS disk config query parameter.

- **[Risk] zones vs availability_zone confusion** → The existing `availability_zone` (TypeString) accepts a single zone, while the new `zones` (TypeList of TypeString) accepts multiple zones for the CBS query. → Mitigation: Documentation clearly explains that `availability_zone` filters the CVM instance types query, while `zones` overrides the CBS disk config query zones parameter.

- **[Risk] disk_types duplication with cbs_filter.disk_types** → Both `disk_types` (top-level) and `cbs_filter.disk_types` (nested) serve similar purposes for the `DescribeDiskConfigQuota` API. When both are provided, the top-level `disk_types` should override `cbs_filter.disk_types`. → Mitigation: The top-level `disk_types` takes precedence when both are specified. This is documented clearly.
