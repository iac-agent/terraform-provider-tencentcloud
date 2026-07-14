## Context

The `tencentcloud_instance_types` data source queries CVM instance type configurations via the `DescribeZoneInstanceConfigInfos` API and optionally queries CBS disk configurations via the `DescribeDiskConfigQuota` API when `cbs_filter` is provided. 

Currently, the `cbs_filter` nested schema only supports 3 parameters (`disk_types`, `disk_charge_type`, `disk_usage`). When calling `DescribeDiskConfigQuota`, the code:
- Hardcodes `InquiryType` to `"INQUIRY_CVM_CONFIG"` 
- Derives `InstanceFamilies` from the CVM instance type's `family` field
- Derives `DiskChargeType` from the `cbs_filter.disk_charge_type` input

The `DescribeDiskConfigQuota` API (in CBS SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312`) supports these as direct input parameters:
- `InquiryType` (*string): Query category - `INQUIRY_CBS_CONFIG` for disk configs alone, `INQUIRY_CVM_CONFIG` for disk+instance combos
- `DiskChargeType` (*string): Payment model - `PREPAID` or `POSTPAID_BY_HOUR`
- `InstanceFamilies` ([]*string): Instance family names for filtering

## Goals / Non-Goals

**Goals:**
- Allow users to specify `inquiry_type`, `disk_charge_type`, and `instance_families` directly as input parameters in the `cbs_filter` block
- Pass user-provided values to the `DescribeDiskConfigQuota` API call, falling back to current defaults when not provided
- Maintain backward compatibility with existing configurations

**Non-Goals:**
- Changing the default behavior when these parameters are not specified (defaults remain the same)
- Adding any other new parameters beyond these three
- Modifying the `DescribeZoneInstanceConfigInfos` API call behavior

## Decisions

### Decision 1: Add parameters to existing `cbs_filter` nested schema
**Rationale**: The three new parameters (`inquiry_type`, `instance_families`, `disk_charge_type`) are all related to the CBS disk configuration query. Adding them to the existing `cbs_filter` block keeps related parameters grouped logically. The `disk_charge_type` parameter already exists in `cbs_filter` - it will need its behavior adjusted so user input overrides the current behavior.

**Alternative considered**: Adding separate top-level parameters - rejected because these parameters only apply when CBS filtering is active, and grouping them in `cbs_filter` is more natural.

### Decision 2: `disk_charge_type` already exists in `cbs_filter` - adjust its pass-through
**Rationale**: `disk_charge_type` is already an optional parameter in the `cbs_filter` schema. Currently, its value is passed directly to the `DescribeDiskConfigQuota` request. No schema change is needed for this parameter - the current behavior already allows user input. However, the default/fallback behavior needs review: if the user doesn't provide `disk_charge_type`, the current code passes it from `cbsFilterParams`, which already handles the case where it's not set.

**Note**: Since `disk_charge_type` already exists as an input parameter in `cbs_filter`, we only need to add `inquiry_type` and `instance_families` as new schema fields.

### Decision 3: `inquiry_type` defaults to `INQUIRY_CVM_CONFIG` when not specified
**Rationale**: The current hardcoded value is `INQUIRY_CVM_CONFIG`. Keeping this as the default maintains backward compatibility. Users can override it to `INQUIRY_CBS_CONFIG` for querying disk configurations independently.

### Decision 4: `instance_families` overrides CVM instance family when provided
**Rationale**: Currently, `InstanceFamilies` is derived from each instance type's `family` field. When the user provides `instance_families` in `cbs_filter`, it should override this derivation, allowing filtering by specific families regardless of the instance type's family.

## Risks / Trade-offs

- **[Risk] Behavior change for `instance_families`**: If user provides `instance_families`, the per-instance-type family override may cause unexpected CBS results → **Mitigation**: Document clearly that `instance_families` overrides the instance type's family; when not provided, current behavior (using instance type family) remains unchanged
- **[Risk] Default `inquiry_type` is hardcoded**: Changing from hardcoded to default parameter may cause subtle differences → **Mitigation**: Default value matches current hardcoded value (`INQUIRY_CVM_CONFIG`), so behavior is identical when parameter is not specified
- **[Trade-off] `disk_charge_type` is already in schema**: Since it already exists, we don't add it as a new field, but we should ensure its value properly reaches the API call in all cases
