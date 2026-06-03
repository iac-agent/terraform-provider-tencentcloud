## Context

Terraform Provider for TencentCloud has existing CVM image datasources (`tencentcloud_image` and `tencentcloud_images`) that use legacy naming conventions and lack several fields from the latest DescribeImages API response. The new `tencentcloud_cvm_image_4` datasource follows the current naming convention (`tencentcloud_<product>_<resource>`) and exposes the complete set of Image struct fields.

The CVM DescribeImages API supports querying images by `ImageIds`, `Filters`, `Offset`, `Limit`, and `InstanceType`. The API returns an `ImageSet` array with each Image containing 20 fields.

The existing CVM service layer already has a `DescribeImagesByFilter` method that handles pagination, but it uses a `map[string][]string` filter approach rather than the SDK-native `[]*Filter` approach. The new datasource needs to support both `ImageIds` (direct ID-based query) and `Filters` (SDK Filter struct-based query).

## Goals / Non-Goals

**Goals:**
- Create a new datasource `tencentcloud_cvm_image_4` that calls DescribeImages API
- Support input parameters: `image_ids`, `filters`, `instance_type`
- Expose the full `image_set` output with all Image struct fields including tags, license_type, image_family, image_deprecated, cdc_cache_status (fields missing from the legacy datasource)
- Follow the RESOURCE_KIND_DATASOURCE pattern as `tencentcloud_igtm_instance_list`
- Register the datasource in `provider.go` and `provider.md`
- Create documentation in `.md` file

**Non-Goals:**
- Do not modify or deprecate existing `tencentcloud_image` or `tencentcloud_images` datasources
- Do not add client-side filtering (like image_name_regex or os_name fuzzy match) - those are legacy patterns
- Do not call CBS service for snapshot details - use Snapshot fields directly from CVM API response

## Decisions

### 1. Direct SDK API call vs. service layer reuse
**Decision**: Call the DescribeImages API directly in the datasource read function using `paramMap` approach, with a new service method `DescribeCvmImage4ByFilter` that accepts `paramMap`.

**Rationale**: The existing `DescribeImagesByFilter` uses `map[string][]string` for filters which doesn't support `ImageIds` directly. The new method needs to handle both `ImageIds` and `Filters` as separate request parameters. Following the `igtm_instance_list` pattern with `paramMap` provides more flexibility.

### 2. Snapshot fields
**Decision**: Use CVM Snapshot struct fields directly (snapshot_id, disk_usage, disk_size) without calling CBS service.

**Rationale**: The CVM Snapshot struct has 3 fields (SnapshotId, DiskUsage, DiskSize). The legacy datasource fetches additional `snapshot_name` from CBS, but this adds unnecessary dependency and complexity. The new datasource keeps it simple and uses only what CVM API provides.

### 3. Filter schema structure
**Decision**: Use TypeList with nested `name` and `values` fields (similar to `igtm_instance_list`), matching the CVM SDK Filter struct.

**Rationale**: The CVM SDK `Filter` struct has `Name` (string) and `Values` ([]*string). The `image_ids` parameter is separate from filters per the API specification (ImageIds and Filters cannot be specified simultaneously).

### 4. Resource ID generation
**Decision**: Use `helper.BuildToken()` for the datasource ID.

**Rationale**: Following the `igtm_instance_list` pattern for list-type datasources that don't have a single natural key.

## Risks / Trade-offs

- [API constraint] ImageIds and Filters are mutually exclusive in the DescribeImages API → Document this in schema descriptions and validate at Terraform level
- [Compatibility] The new datasource coexists with legacy datasources → No migration path needed; users can adopt the new datasource at their own pace
- [Field completeness] SnapshotSet only has 3 fields from CVM API (no snapshot_name) → Acceptable trade-off for simplicity; users needing snapshot_name can use the CBS datasource
