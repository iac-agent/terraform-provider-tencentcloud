## Context

The `tencentcloud_instance_types` data source queries instance type configurations using the CVM `DescribeZoneInstanceConfigInfos` API and optionally queries CBS disk configurations using the CBS `DescribeDiskConfigQuota` API when a `cbs_filter` block is provided.

Currently, the `DescribeDiskConfigQuota` call in the CBS service layer has the following limitations:
1. `InquiryType` is hardcoded to `"INQUIRY_CVM_CONFIG"` - users cannot switch to `"INQUIRY_CBS_CONFIG"` mode to query CBS-only configurations without CVM instance coupling
2. `InstanceFamilies` is derived from the instance type query results' `family` field - users cannot specify custom instance families for CBS queries independently of the instance type results
3. `DiskChargeType` already exists as a field in the `cbs_filter` block, so it is already configurable

The CBS SDK's `DescribeDiskConfigQuotaRequest` structure already supports all three parameters (`InquiryType`, `DiskChargeType`, `InstanceFamilies`), so no SDK update is needed.

### Current Architecture
- Data source: `tencentcloud/services/cvm/data_source_tc_instance_types.go`
- CBS service: `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - `DescribeDiskConfigQuota` method accepts a `map[string]interface{}`
- The `cbs_filter` block schema currently has: `disk_types`, `disk_charge_type`, `disk_usage`
- The CBS service method currently hardcodes `InquiryType` and derives `InstanceFamilies` from the map

### Key Constraints
- Must maintain backward compatibility - existing configurations must work unchanged
- Default `InquiryType` should remain `"INQUIRY_CVM_CONFIG"` when not specified by user
- `InstanceFamilies` when not specified should continue to be derived from instance type results
- No vendor SDK update required - CBS SDK already has these fields

## Goals / Non-Goals

**Goals:**
- Add `InquiryType` and `InstanceFamilies` parameters to the `cbs_filter` schema block
- Allow users to override the hardcoded `InquiryType` value with `"INQUIRY_CBS_CONFIG"`
- Allow users to specify `InstanceFamilies` directly instead of relying on instance type results
- Maintain backward compatibility with default values

**Non-Goals:**
- Changing the `DescribeZoneInstanceConfigInfos` API behavior or parameters
- Restructuring the `cbs_filter` block beyond adding new optional fields
- Adding new top-level parameters outside the `cbs_filter` block
- Modifying the CBS SDK or updating vendor dependencies

## Decisions

### Decision 1: Add parameters to `cbs_filter` block rather than as top-level fields

**Choice**: Add `InquiryType` and `InstanceFamilies` as fields within the existing `cbs_filter` nested schema block.

**Rationale**: All three parameters (`InquiryType`, `DiskChargeType`, `InstanceFamilies`) are input parameters for the `DescribeDiskConfigQuota` CBS API, which is only invoked when `cbs_filter` is present. Placing them in the `cbs_filter` block maintains logical grouping - they all control the CBS disk configuration query behavior. This is consistent with how `disk_charge_type` already lives in `cbs_filter`.

**Alternatives considered**:
- Add as top-level schema fields: This would create confusion since these parameters only affect CBS queries, not instance type queries. Top-level parameters would imply they affect both APIs.
- Create a separate `cbs_params` block: This would add unnecessary complexity and could conflict with `cbs_filter` semantics.

### Decision 2: Keep `DiskChargeType` as existing field in `cbs_filter`

**Choice**: `DiskChargeType` already exists as `disk_charge_type` in the `cbs_filter` block. No schema change needed for this field.

**Rationale**: The `cbs_filter` block already has a `disk_charge_type` field (TypeString, Optional) that maps to `request.DiskChargeType` in the `DescribeDiskConfigQuota` call. This field is already fully functional and maps correctly to the CBS API. The requirement lists `DiskChargeType` as a new parameter, but since it already exists and is working, we should document it as an existing capability rather than add a duplicate.

### Decision 3: Default values for backward compatibility

**Choice**:
- `InquiryType`: When not specified by the user, default to `"INQUIRY_CVM_CONFIG"` (current hardcoded behavior)
- `InstanceFamilies`: When not specified by the user, derive from instance type results' `family` field (current behavior)

**Rationale**: This ensures that existing Terraform configurations that don't specify these new parameters continue to work identically. Users only need to set these parameters when they want different CBS query behavior.

### Decision 4: Modify CBS service method to accept explicit parameters

**Choice**: Modify the `DescribeDiskConfigQuota` method signature to accept `InquiryType` and `InstanceFamilies` as optional parameters alongside the existing `cvmInfo` map.

**Rationale**: Currently the method derives these values from the `cvmInfo` map. With the new parameters, we need to support both the default behavior (derive from map) and user-specified override values. The simplest approach is to pass these as optional string parameters that, when empty, fall back to the current derivation logic.

**Alternatives considered**:
- Add all new parameters to the `cvmInfo` map: This would mix user-provided override values with derived values in the same map, making it unclear which values are user-provided vs derived. It could also cause conflicts if both derived and user-provided values exist for the same key.
- Create a new method variant: This would add unnecessary code duplication.

## Risks / Trade-offs

- [Risk] Users might accidentally set `InquiryType` to `"INQUIRY_CBS_CONFIG"` without understanding it changes the CBS query mode → Mitigation: Documentation clearly explains both modes and their effects
- [Risk] Adding `InstanceFamilies` as an override might confuse users who expect it to filter instance type results rather than CBS results → Mitigation: Documentation and schema description clearly state this parameter only affects CBS disk configuration queries
- [Trade-off] Keeping `InquiryType` and `InstanceFamilies` inside `cbs_filter` limits discoverability since they're nested, but this is consistent with the existing pattern and keeps related parameters grouped logically
