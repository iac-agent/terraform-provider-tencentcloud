## Context

The `tencentcloud_instance_types` data source (`data_source_tc_instance_types.go`) provides CVM instance type information and optional CBS disk configuration quota details through the `cbs_filter` nested block. Currently, when `cbs_filter` is specified, the data source calls the CBS `DescribeDiskConfigQuota` API for each instance type returned by `DescribeZoneInstanceConfigInfos`, passing the single instance family (`family`) from each result item.

The CBS SDK's `DescribeDiskConfigQuota` API (in `cbs/v20170312/models.go`) supports the `InstanceFamilies` field as a `[]*string` parameter, allowing multiple instance families to be queried simultaneously. However, the current Terraform implementation only passes a single instance family derived from the instance type result, and there is no way for users to specify their own list of instance families to filter by.

Current `cbs_filter` schema fields:
- `disk_types` (TypeList of TypeString) - already mapped to `request.DiskTypes`
- `disk_charge_type` (TypeString) - already mapped to `request.DiskChargeType`
- `disk_usage` (TypeString) - already mapped to `request.DiskUsage`

The `InstanceFamilies` parameter is missing from `cbs_filter`, and the `DescribeDiskConfigQuota` service method (`tencentcloud/services/cbs/service_tencentcloud_cbs.go`) currently derives `InstanceFamilies` from a single `family` string in `cvmInfo`, which comes from the instance type query result rather than user input.

## Goals / Non-Goals

**Goals:**
- Add `instance_families` field to the `cbs_filter` nested schema in `tencentcloud_instance_types` data source
- Modify the `DescribeDiskConfigQuota` service method to accept and pass `instance_families` as a list of strings to the API request
- When `instance_families` is provided in `cbs_filter`, use it instead of the single `family` from instance type results
- Maintain backward compatibility: when `instance_families` is not provided, the existing behavior (using `family` from instance type results) should remain unchanged

**Non-Goals:**
- Changing the `DescribeZoneInstanceConfigInfos` API call parameters (no changes to Filters)
- Adding new computed/output fields to the data source
- Modifying any other data sources or resources
- Changing the `disk_charge_type` field behavior (it already exists in `cbs_filter`)

## Decisions

1. **Add `instance_families` to `cbs_filter` as TypeList of TypeString**
   - Rationale: Matches the CBS SDK's `InstanceFamilies []*string` type. Using TypeList allows users to specify multiple instance families.
   - Alternative considered: TypeSet (would enforce uniqueness but loses order) - rejected because the CBS API accepts a list, not a set.

2. **Modify `DescribeDiskConfigQuota` service method to accept `instance_families` as a separate parameter**
   - Rationale: The current method signature `DescribeDiskConfigQuota(ctx context.Context, cvmInfo map[string]interface{})` uses a generic map. Adding `instance_families` support within this map maintains the current pattern and allows the method to differentiate between user-provided instance families vs. the single `family` derived from instance type results.
   - Alternative considered: Change method signature to accept explicit parameters - rejected as it would require broader refactoring across the codebase.

3. **Priority logic: user-provided `instance_families` overrides single `family`**
   - Rationale: When the user explicitly specifies `instance_families`, their input should take precedence over the automatically derived `family` from instance type results. If `instance_families` is not specified, fall back to the existing `family` behavior.
   - Alternative considered: Always combine user `instance_families` with derived `family` - rejected because it would add duplicates and create confusing behavior.

## Risks / Trade-offs

- **[Backward compatibility]** → Mitigation: The `instance_families` field is Optional with no default. When not specified, the existing `family`-based behavior is preserved. No breaking changes to existing Terraform configurations.
- **[Empty list handling]** → Mitigation: If `instance_families` is specified as an empty list, it should be treated as "not specified" and fall back to the single `family` behavior.
