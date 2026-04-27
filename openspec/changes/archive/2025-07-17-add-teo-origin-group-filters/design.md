## Context

The `tencentcloud_teo_origin_group` resource manages TEO (TencentCloud EdgeOne) origin groups via the Terraform provider. The resource currently supports full CRUD operations using the `CreateOriginGroup`, `DescribeOriginGroup`, `ModifyOriginGroup`, and `DeleteOriginGroup` APIs.

The `DescribeOriginGroup` API supports an `AdvancedFilter` mechanism via its `Filters` request parameter, which allows filtering by `origin-group-id` and `origin-group-name` with optional fuzzy matching. However, the current Terraform resource does not expose this `filters` parameter in its schema, and the service method `DescribeTeoOriginGroupById` hardcodes only the `origin-group-id` filter internally.

The cloud API's `AdvancedFilter` struct contains:
- `Name` (*string): The field name to filter on (e.g., `origin-group-id`, `origin-group-name`)
- `Values` ([]*string): The filter values
- `Fuzzy` (*bool): Whether to enable fuzzy matching

Current resource composite ID format: `zoneId:groupId` (using `tccommon.FILED_SP` separator).

## Goals / Non-Goals

**Goals:**
- Add `filters` parameter to the `tencentcloud_teo_origin_group` resource schema as an Optional+Computed TypeList field
- Each filter item maps to the `AdvancedFilter` struct with `name`, `values`, and `fuzzy` sub-fields
- Update the read operation to pass user-provided `filters` to the `DescribeOriginGroup` API request
- Update the service method to accept and use `filters` from the resource
- Maintain full backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Changing existing schema fields or their behavior
- Adding `filters` to create, update, or delete operations (filters are only relevant for the read/query operation)
- Creating a new datasource for origin groups (this change is for the existing resource only)

## Decisions

### Decision 1: `filters` as Optional+Computed TypeList

**Choice**: Add `filters` as `TypeList` with `Optional: true, Computed: true` and nested schema containing `name`, `values`, `fuzzy`.

**Rationale**: The `filters` parameter is optional in the cloud API. Making it `Computed: true` ensures that if no filters are specified by the user, the resource read operation still works correctly (the service method will use default `origin-group-id` filter). Making it `Optional` allows users to provide custom filter configurations.

**Alternative considered**: Making `filters` a `TypeSet` to avoid order sensitivity. However, the cloud API's `AdvancedFilter` is an ordered list, and using `TypeList` is consistent with similar filter implementations in other resources.

### Decision 2: Service method signature update

**Choice**: Modify `DescribeTeoOriginGroupById` to accept `zoneId` and `filters` parameters, and pass them to the API request.

**Rationale**: The current service method doesn't pass `ZoneId` to the API request, even though the API marks it as required. Adding `zoneId` as a parameter fixes this gap and ensures the API call is complete. The filters parameter allows user-configured filters to be passed through.

**Alternative considered**: Creating a separate service method for filtered queries. This would add unnecessary code duplication for a simple parameter addition.

### Decision 3: Backward compatibility with default filter behavior

**Choice**: In the read function, if `filters` is not set by the user, the service method will still use the default `origin-group-id` filter (current behavior) to look up the origin group by its ID.

**Rationale**: This ensures that existing Terraform configurations that don't specify `filters` continue to work exactly as before. The default filter behavior matches the current hardcoded filter in `DescribeTeoOriginGroupById`.

## Risks / Trade-offs

- [Risk] Adding `filters` to a resource (not datasource) is unconventional → Mitigation: The parameter is Optional+Computed and defaults to the existing behavior, so it doesn't change the standard resource lifecycle. Users who need custom filtering can opt in.
- [Risk] The `DescribeOriginGroup` API returns a list of origin groups, and the service method currently picks the first match → Mitigation: This is the existing behavior and is not changed by this addition. The default filter by `origin-group-id` ensures a unique match.
- [Risk] Changing the service method signature could affect other callers → Mitigation: `DescribeTeoOriginGroupById` is only called from the resource's read function, so the signature change is safe.
