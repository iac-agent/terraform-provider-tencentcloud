## Context

The `tencentcloud_instance_types` data source currently exposes instance type information from the CVM `DescribeZoneInstanceConfigInfos` API and disk configuration information from the CBS `DescribeDiskConfigQuota` API. The `DescribeDiskConfigQuota` API returns a `DiskConfigSet` array where each `DiskConfig` element contains a `DiskUsage` field (type `*string`), indicating whether the disk is a system disk (`SYSTEM_DISK`) or data disk (`DATA_DISK`).

Currently, `disk_usage` is available as:
1. An **input parameter** in `cbs_filter` (used as a query filter for the CBS API)
2. A **computed output** inside the nested `cbs_configs` structure (each CBS config entry has `disk_usage`)

However, `DiskUsage` is not exposed as a **top-level computed output parameter** on the data source itself, which means users cannot directly reference the disk usage type without navigating into the nested `cbs_configs` structure.

The CBS SDK `DiskConfig` struct (in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312/models.go`) already defines `DiskUsage` as a `*string` field with valid values: `SYSTEM_DISK` (system disk) and `DATA_DISK` (data disk).

## Goals / Non-Goals

**Goals:**
- Add `DiskUsage` as a new top-level computed output parameter (`TypeString`) on the `tencentcloud_instance_types` data source
- Map `DiskUsage` from the `DescribeDiskConfigQuota` API response (`DiskConfig.DiskUsage`) to the new top-level field
- Maintain full backward compatibility with existing configurations

**Non-Goals:**
- Not modifying the existing `disk_usage` field in `cbs_filter` or `cbs_configs`
- Not adding any new input parameters
- Not changing the CBS service layer's `DescribeDiskConfigQuota` function
- Not adding other fields from `DescribeDiskConfigQuota` response beyond `DiskUsage`

## Decisions

### Decision 1: Add `DiskUsage` as a top-level computed field on the data source schema
**Rationale**: The requirement specifies `response.DiskUsage` → `DiskUsage` as a new output parameter. Adding it at the top level of the data source schema allows users to directly reference the disk usage type in their Terraform configurations without needing to iterate through `cbs_configs`.

**Alternative considered**: Adding `DiskUsage` inside the `instance_types` element schema. This was rejected because `DiskUsage` comes from the CBS API (not from the CVM instance type API), and it represents the overall disk usage type for the query, not a per-instance-type attribute.

**Decision**: Add `DiskUsage` as a top-level computed `TypeString` field on the data source, directly under the root schema. The value will be populated from the `cbs_filter.disk_usage` input if provided, or from the first `DiskConfig.DiskUsage` value in the CBS response.

### Decision 2: Data mapping approach
**Rationale**: Since `DiskUsage` is already available in the `cbs_filter` input parameter and echoed back in each `DiskConfig.DiskUsage` in the CBS response, the top-level `DiskUsage` computed field should be populated from the `cbs_filter.disk_usage` input value (when provided), as this directly represents the user's intent. If `cbs_filter` is not provided, the field will remain empty/null.

## Risks / Trade-offs

- **[Backward compatibility]** → Mitigation: The new `DiskUsage` field is an optional computed field, so existing configurations will not be affected. Users who don't reference the new field will see no changes in their plans.
- **[Data source semantics]** → Mitigation: `DiskUsage` is context-dependent (it only has a meaningful value when `cbs_filter` is provided). This is consistent with how `cbs_configs` already behaves - it's only populated when `cbs_filter` is used.
