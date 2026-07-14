## Context

The `tencentcloud_instance_types` data source uses two APIs:
1. `DescribeZoneInstanceConfigInfos` (CVM) to fetch instance type configurations
2. `DescribeDiskConfigQuota` (CBS) to fetch disk configuration quotas when `cbs_filter` is provided

Currently, when calling `DescribeDiskConfigQuota` via `cbs_filter`, the `InquiryType` parameter is hardcoded to `"INQUIRY_CVM_CONFIG"` in the CBS service layer (`service_tencentcloud_cbs.go` line 979). This prevents users from querying only cloud disk configuration lists without instance pairing constraints.

The `DescribeDiskConfigQuota` CBS API supports `InquiryType` values:
- `INQUIRY_CBS_CONFIG`: Query cloud disk configuration list only
- `INQUIRY_CVM_CONFIG`: Query cloud disk and instance pairing configuration list (current default)

The `cbs_filter` block already has `disk_charge_type` and `disk_usage` fields but lacks `inquiry_type`.

## Goals / Non-Goals

**Goals:**
- Add `inquiry_type` as an optional input parameter within the `cbs_filter` block
- Allow users to specify whether they want `INQUIRY_CBS_CONFIG` or `INQUIRY_CVM_CONFIG` as the query category
- Maintain backward compatibility by defaulting `inquiry_type` to `INQUIRY_CVM_CONFIG`
- Update the CBS service `DescribeDiskConfigQuota` function to accept `inquiry_type` from the caller

**Non-Goals:**
- Changing any other parameters in `cbs_filter` or the data source schema
- Adding `InquiryType` or `DiskChargeType` as top-level data source parameters (they belong inside `cbs_filter`)
- Modifying the `DescribeZoneInstanceConfigInfos` CVM API call behavior
- Adding new computed/output fields to the data source

## Decisions

### Decision 1: Place `inquiry_type` inside `cbs_filter` block

**Rationale**: The `InquiryType` parameter is specific to the `DescribeDiskConfigQuota` CBS API call, which is only invoked when `cbs_filter` is provided. Placing it inside `cbs_filter` keeps related parameters together and maintains logical grouping. This follows the existing pattern where `disk_charge_type`, `disk_types`, and `disk_usage` are already inside `cbs_filter`.

**Alternative considered**: Adding `inquiry_type` as a top-level data source parameter was rejected because it's only relevant when `cbs_filter` is used, and a top-level parameter would be confusing without `cbs_filter`.

### Decision 2: Default `inquiry_type` to `INQUIRY_CVM_CONFIG`

**Rationale**: The current hardcoded value is `INQUIRY_CVM_CONFIG`, so defaulting to this value ensures backward compatibility. Existing Terraform configurations that don't specify `inquiry_type` will continue to work exactly as before.

### Decision 3: Modify CBS service to accept `inquiry_type` dynamically

**Rationale**: Currently `DescribeDiskConfigQuota` in `service_tencentcloud_cbs.go` hardcodes `request.InquiryType = helper.String("INQUIRY_CVM_CONFIG")`. We need to change this to accept the value from the caller's `cbsFilterParams` map, so the data source can pass the user-specified value.

## Risks / Trade-offs

- **[Risk] Existing CBS callers may rely on hardcoded InquiryType** → Mitigation: The default value remains `INQUIRY_CVM_CONFIG`, so callers that don't pass `inquiry_type` in the map will get the same behavior. We need to handle the case where `inquiry_type` is not in the map by using the default value.
- **[Risk] Invalid InquiryType values** → Mitigation: The Terraform schema description should clearly document the valid values. The CBS API will reject invalid values and return an error, which will be surfaced to the user.
