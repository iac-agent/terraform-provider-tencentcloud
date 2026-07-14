## Context

The `tencentcloud_instance_types` data source currently queries two cloud APIs:
1. `DescribeZoneInstanceConfigInfos` (CVM) - retrieves instance type configuration information
2. `DescribeDiskConfigQuota` (CBS) - retrieves CBS disk configuration quota information when `cbs_filter` is provided

The `DescribeDiskConfigQuota` API call in `tencentcloud/services/cbs/service_tencentcloud_cbs.go` currently:
- Hardcodes `InquiryType` to "INQUIRY_CVM_CONFIG"
- Derives `InstanceFamilies` from the instance type query results (auto-populated from `family` field)
- Passes `DiskChargeType` from the `cbs_filter` sub-parameter
- Only called when `cbs_filter` is provided by the user

The data source currently has `cbs_filter` as a sub-parameter block containing `disk_types`, `disk_charge_type`, and `disk_usage`. The `DiskChargeType` field already exists within `cbs_filter` but needs to also be available as a top-level parameter.

The `Available` field already exists within `cbs_configs` (nested inside `instance_types`), but needs to be exposed as a top-level computed attribute on the data source.

## Goals / Non-Goals

**Goals:**
- Add `InquiryType` as a top-level optional parameter to allow users to specify the query category for `DescribeDiskConfigQuota`
- Add `DiskChargeType` as a top-level optional parameter to allow users to specify the disk payment model directly (not only within `cbs_filter`)
- Add `InstanceFamilies` as a top-level optional parameter to allow users to filter by specific instance families when querying disk configs
- Add `Available` as a top-level computed attribute to expose the overall disk configuration availability status
- Maintain backward compatibility - all new fields are optional/computed additions

**Non-Goals:**
- Changing the existing `cbs_filter` block or its sub-parameters
- Modifying the `DescribeZoneInstanceConfigInfos` API call or its parameters
- Adding new nested schema structures beyond the top-level `Available` field

## Decisions

### Decision 1: Top-level vs nested parameter placement
**Choice**: Add `InquiryType`, `DiskChargeType`, and `InstanceFamilies` as top-level parameters on the data source, not nested within `cbs_filter`.

**Rationale**: The requirement specifies these as top-level schema parameters (SchemaName paths). While `DiskChargeType` already exists in `cbs_filter`, the requirement calls for a separate top-level parameter. The top-level parameters will be passed to `DescribeDiskConfigQuota` alongside or independently of `cbs_filter`.

### Decision 2: InquiryType default value
**Choice**: When `InquiryType` is not specified, the default behavior should remain "INQUIRY_CVM_CONFIG" to maintain backward compatibility.

**Rationale**: Current code hardcodes this value. Defaulting to the same value ensures existing configurations work unchanged.

### Decision 3: Available field semantics
**Choice**: `Available` will be a top-level computed boolean field representing whether the queried disk configuration is available. When `DescribeDiskConfigQuota` returns results, `Available` will be set based on the availability status of the disk configurations.

**Rationale**: The `Available` field in `DiskConfig` structure represents per-configuration availability. As a top-level field, it provides a quick check for overall availability. When multiple disk configs are returned, we aggregate availability - if any config is available, the top-level `Available` is true.

### Decision 4: Interaction between top-level DiskChargeType and cbs_filter disk_charge_type
**Choice**: The top-level `DiskChargeType` and `cbs_filter.disk_charge_type` serve different purposes. The top-level `DiskChargeType` is used as a filter for `DescribeDiskConfigQuota` when querying without the full `cbs_filter` context. When both are provided, the top-level value takes precedence as it's more explicit.

**Rationale**: This allows users to specify disk charge type directly without having to use the full `cbs_filter` block, providing a simpler API for common use cases.

## Risks / Trade-offs

- [Risk] Top-level `DiskChargeType` duplicates functionality of `cbs_filter.disk_charge_type` → Mitigation: Document the difference clearly; top-level parameter provides simpler access for the common case
- [Risk] `InstanceFamilies` top-level parameter may conflict with auto-populated family from instance types → Mitigation: When `InstanceFamilies` is explicitly provided, use it instead of auto-populating from instance type results
- [Risk] `Available` aggregation logic (whether all or any config must be available) → Mitigation: Use "any available = true" semantics to match typical user expectations
