## Context

The `tencentcloud_teo_realtime_log_delivery` resource manages EdgeOne (TEO) realtime log delivery tasks. Currently, the resource schema does not expose the `Filters` parameter of the `DescribeRealtimeLogDeliveryTasks` API, nor does it expose the `RealtimeLogDeliveryTasks` list from the response. The read operation uses a fixed filter (task-id) to look up a single task by ID.

The cloud API `DescribeRealtimeLogDeliveryTasks` supports:
- **Input**: `ZoneId` (string), `Offset` (int64), `Limit` (uint64), `Filters` ([]*AdvancedFilter)
- **Output**: `TotalCount` (uint64), `RealtimeLogDeliveryTasks` ([]*RealtimeLogDeliveryTask), `RequestId` (string)

The `AdvancedFilter` type has three fields: `Name` (string), `Values` ([]*string), `Fuzzy` (*bool). Available filter names are: `task-id`, `task-name`, `entity-list`, `task-type`.

The `RealtimeLogDeliveryTask` type contains 15 fields: TaskId, TaskName, DeliveryStatus, TaskType, EntityList, LogType, Area, Fields, CustomFields, DeliveryConditions, Sample, LogFormat, CLS, CustomEndpoint, S3, CreateTime, UpdateTime.

## Goals / Non-Goals

**Goals:**
- Add `filters` as an Optional input parameter to the resource schema, mapping to the `Filters` field of `DescribeRealtimeLogDeliveryTasksRequest`
- Add `realtime_log_delivery_tasks` as a Computed output parameter to the resource schema, mapping to the `RealtimeLogDeliveryTasks` field of `DescribeRealtimeLogDeliveryTasksResponse`
- Maintain full backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Modifying existing schema fields or changing existing CRUD behavior
- Adding new CRUD operations or changing the resource lifecycle
- Changing the resource ID format or composite key structure

## Decisions

1. **`filters` parameter type**: Use `TypeList` with nested `schema.Resource` containing `name` (string, Required), `values` (list of strings, Required), and `fuzzy` (bool, Optional) fields. This maps directly to the `AdvancedFilter` SDK type.

   *Rationale*: Using `schema.Resource` for nested objects is the standard pattern in this provider (see `delivery_conditions`, `custom_fields` etc. in the same resource). The `AdvancedFilter` type has Name, Values, and Fuzzy fields, which map naturally.

2. **`realtime_log_delivery_tasks` parameter type**: Use `TypeList` with nested `schema.Resource` containing all fields from `RealtimeLogDeliveryTask` SDK type as computed sub-fields.

   *Rationale*: This follows the pattern used by other datasource/resources in the provider that expose list results. Each sub-field is Computed since this is a read-only output.

3. **`filters` is Optional**: Users should not be required to specify filters. When not specified, the existing behavior (filtering by task-id from the composite ID) is preserved.

   *Rationale*: Backward compatibility requires that existing configurations without `filters` continue to work unchanged.

4. **Read function modification**: When `filters` is specified, pass the user-provided filters to `DescribeRealtimeLogDeliveryTasksRequest.Filters`. When not specified, keep the current behavior of using task-id filter from the composite ID.

   *Rationale*: The current read logic uses a fixed AdvancedFilter with name "task-id" and value from the resource ID. We need to preserve this as the default behavior while allowing custom filters.

5. **Service layer**: Add a new service method or modify the existing `DescribeTeoRealtimeLogDeliveryById` to support accepting custom filters and returning the full list.

   *Rationale*: The existing method only returns a single task. We need a method that can return the full list when filters are provided.

## Risks / Trade-offs

- **[Breaking read behavior]** → When `filters` is not specified, the read behavior MUST remain identical to the current implementation (filter by task-id from composite ID). The new behavior only activates when `filters` is explicitly set.
- **[Large response data]** → The `realtime_log_delivery_tasks` list could be large. Mitigate by setting Limit to the API maximum (1000) when querying, and document that this is a computed field that reflects the current API state.
- **[Schema compatibility]** → Both new fields are additive (Optional/Computed) and do not change existing fields, so no state migration is needed.
