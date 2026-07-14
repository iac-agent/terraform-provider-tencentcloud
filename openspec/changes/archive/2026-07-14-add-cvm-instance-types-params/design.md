## Context

The `tencentcloud_instance_types` datasource is a RESOURCE_KIND_DATASOURCE in the CVM service that queries instance type configurations. When users need CBS (Cloud Block Storage) disk configuration information alongside instance types, the datasource calls `DescribeDiskConfigQuota` from the CBS API via the `cbs_filter` parameter.

Currently:
- `InquiryType` is hardcoded to `INQUIRY_CVM_CONFIG` in the CBS service's `DescribeDiskConfigQuota` function, meaning users cannot query pure cloud disk configurations (`INQUIRY_CBS_CONFIG`).
- `DiskChargeType` is only accessible as a nested parameter inside `cbs_filter.disk_charge_type`, requiring users to use the nested block structure even if they only want to specify the payment model.

The CBS API `DescribeDiskConfigQuotaRequest` supports both `InquiryType` and `DiskChargeType` as direct request parameters, making them suitable as top-level datasource parameters.

## Goals / Non-Goals

**Goals:**
- Add `inquiry_type` as an optional top-level parameter to the `tencentcloud_instance_types` datasource schema
- Add `disk_charge_type` as an optional top-level parameter to the `tencentcloud_instance_types` datasource schema
- Pass these parameters to the `DescribeDiskConfigQuota` API call when `cbs_filter` is present
- Maintain backward compatibility by defaulting `inquiry_type` to `INQUIRY_CVM_CONFIG`
- Update the datasource documentation (.md file) with new parameter descriptions and examples

**Non-Goals:**
- Removing or modifying the existing `cbs_filter.disk_charge_type` nested parameter (it should remain for backward compatibility)
- Changing the `DescribeZoneInstanceConfigInfos` API call or its parameters
- Adding any response/output fields to the datasource

## Decisions

**Decision 1: Top-level parameters vs. extending cbs_filter**

- **Choice**: Add `inquiry_type` and `disk_charge_type` as top-level parameters
- **Rationale**: The requirement explicitly specifies these as top-level SchemaName parameters (`InquiryType` and `DiskChargeType`). The `cbs_filter` block already contains `disk_charge_type`, but the new top-level `disk_charge_type` parameter will be the primary way to pass this value. If both the top-level parameter and `cbs_filter.disk_charge_type` are set, the top-level value will take precedence.
- **Alternative considered**: Extending the `cbs_filter` block to include `inquiry_type` — rejected because the requirement specifies top-level schema names.

**Decision 2: Default value for inquiry_type**

- **Choice**: Default to `INQUIRY_CVM_CONFIG`
- **Rationale**: This matches the current hardcoded behavior, ensuring backward compatibility. Existing configurations that don't specify `inquiry_type` will continue to get the same results.

**Decision 3: Handling overlap between top-level disk_charge_type and cbs_filter.disk_charge_type**

- **Choice**: Top-level `disk_charge_type` takes precedence over `cbs_filter.disk_charge_type`
- **Rationale**: The top-level parameter is the newer, more direct way to set this value. If both are specified, the top-level value should be used to avoid ambiguity. The `cbs_filter.disk_charge_type` remains for backward compatibility but is effectively superseded.

## Risks / Trade-offs

- **[Risk] Parameter overlap confusion**: Users might be confused about whether to use the top-level `disk_charge_type` or `cbs_filter.disk_charge_type`. → **Mitigation**: Document clearly that the top-level parameter is preferred, and the `cbs_filter` nested version is maintained for backward compatibility.
- **[Risk] Breaking change if cbs_filter.disk_charge_type is removed**: → **Mitigation**: Do NOT remove the `cbs_filter.disk_charge_type` parameter; keep it for backward compatibility.
