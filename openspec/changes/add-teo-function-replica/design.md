## Context

Terraform Provider for TencentCloud needs to support managing TEO (TencentCloud EdgeOne) edge function replicas. Currently, only the parent `tencentcloud_teo_function` resource exists. Edge function replicas allow users to deploy different versions of edge function code under the same function, accessible via the `EO-Function-Replica-Name` request header.

The TEO SDK (`teo/v20220901`) provides four synchronous APIs for function replica management: `CreateFunctionReplica`, `DescribeFunctionReplicas`, `ModifyFunctionReplica`, and `DeleteFunctionReplica`. None of these are async.

Since `CreateFunctionReplica` does not return a unique ID (the replica is identified by the combination of `ZoneId + FunctionId + ReplicaName`), the Terraform resource will use a 3-part composite ID: `zone_id#function_id#replica_name` using `tccommon.FILED_SP` as separator.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_teo_function_replica` resource with full CRUD lifecycle
- Use composite ID `zone_id#function_id#replica_name` for resource identification
- Follow existing code patterns from `tencentcloud_teo_function` and `tencentcloud_igtm_strategy`
- Support import via 3-part composite ID
- Add unit tests with gomonkey mocks

**Non-Goals:**
- Adding a data source for function replicas (not in scope)
- Supporting batch replica operations (each Terraform resource manages one replica)
- Modifying existing `tencentcloud_teo_function` resource

## Decisions

### Decision 1: Composite ID with 3 fields
- **Choice**: `zone_id#function_id#replica_name`
- **Rationale**: `CreateFunctionReplica` returns no ID; the replica is uniquely identified by the trio of ZoneId + FunctionId + ReplicaName. The `#` separator (`tccommon.FILED_SP`) is the established convention in this provider.
- **Alternative**: Could use `zone_id#function_id` as ID and store `replica_name` as a separate ForceNew field, but this makes import harder and doesn't reflect the true identity of the resource.

### Decision 2: Read uses DescribeFunctionReplicas with AdvancedFilter
- **Choice**: Use `AdvancedFilter` with `replica-name` filter to query the specific replica by name
- **Rationale**: `DescribeFunctionReplicas` is a list API and there is no single-item Describe API. Filtering by `replica-name` with `Fuzzy=false` returns only the matching replica. Set `Limit=200` (API maximum) to ensure all results are returned.
- **Alternative**: Could fetch all replicas and search client-side, but filtering at API level is more efficient.

### Decision 3: Schema design - zone_id, function_id, replica_name as Required/ForceNew
- **Choice**: `zone_id`, `function_id`, `replica_name` are all `Required: true, ForceNew: true` since they form the composite ID and cannot be changed without recreating the resource.
- **Rationale**: `ModifyFunctionReplica` uses `ReplicaName` as an identifier (not a mutable field), so the replica name cannot be changed after creation. ZoneId and FunctionId are also immutable identifiers.
- **Note**: `ModifyFunctionReplica` can update `content` and `remark` fields.

### Decision 4: Delete uses ReplicaNames list field
- **Choice**: Pass `replica_name` (singular) as a single-element list to `ReplicaNames` in `DeleteFunctionReplica`
- **Rationale**: The delete API accepts `ReplicaNames []*string` for batch deletion, but our Terraform resource manages one replica at a time, so we wrap the single name in a list.

### Decision 5: Filters parameter as computed-only for DescribeFunctionReplicas read
- **Choice**: `sort_by`, `sort_order`, and `filters` are NOT exposed in the resource schema since they are query parameters for the list API, not attributes of the replica itself.
- **Rationale**: These parameters are only relevant for listing/searching, not for managing a single replica resource. The resource only needs `zone_id`, `function_id`, `replica_name`, `content`, `remark` as schema fields, plus `created_on` and `modified_on` as computed fields.

## Risks / Trade-offs

- **[Risk] No single-item Describe API** → Mitigation: Use `DescribeFunctionReplicas` with `AdvancedFilter` by replica-name. If the API is temporarily unavailable and returns an empty list, the Read retry mechanism handles this.
- **[Risk] Replica name character restrictions** → Mitigation: Schema validation via `ValidateFunc` is not strictly required; the cloud API will return an error for invalid names. The name must be 1-50 chars, a-z/0-9/-, no consecutive/leading/trailing hyphens.
- **[Risk] Concurrent modification** → Mitigation: Terraform's state locking prevents concurrent modifications to the same resource.
