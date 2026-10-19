## Requirements

### Requirement: Resource schema includes filters input parameter
The `tencentcloud_teo_realtime_log_delivery` resource SHALL include an Optional `filters` parameter of TypeList, containing nested schema.Resource with three fields:
- `name` (TypeString, Required): The filter field name. Available values: `task-id`, `task-name`, `entity-list`, `task-type`.
- `values` (TypeList of TypeString, Required): The filter values. Maximum 20 values per filter.
- `fuzzy` (TypeBool, Optional): Whether to enable fuzzy query.

#### Scenario: User specifies filters to query tasks
- **WHEN** user sets `filters` with name="task-type" and values=["cls"]
- **THEN** the read operation SHALL pass these filters to `DescribeRealtimeLogDeliveryTasksRequest.Filters` as `AdvancedFilter` objects

#### Scenario: User does not specify filters
- **WHEN** user does not set `filters` in the configuration
- **THEN** the read operation SHALL use the default behavior (filter by task-id from composite ID), preserving backward compatibility

### Requirement: Resource schema includes realtime_log_delivery_tasks computed output
The `tencentcloud_teo_realtime_log_delivery` resource SHALL include a Computed `realtime_log_delivery_tasks` parameter of TypeList, containing nested schema.Resource with the following computed fields from the `RealtimeLogDeliveryTask` SDK type:
- `task_id` (TypeString, Computed)
- `task_name` (TypeString, Computed)
- `delivery_status` (TypeString, Computed)
- `task_type` (TypeString, Computed)
- `entity_list` (TypeList of TypeString, Computed)
- `log_type` (TypeString, Computed)
- `area` (TypeString, Computed)
- `fields` (TypeList of TypeString, Computed)
- `custom_fields` (TypeList, Computed): nested structure matching existing `custom_fields` schema
- `delivery_conditions` (TypeList, Computed): nested structure matching existing `delivery_conditions` schema
- `sample` (TypeInt, Computed)
- `log_format` (TypeList, Computed): nested structure matching existing `log_format` schema
- `cls` (TypeList, Computed): nested structure matching existing `cls` schema
- `custom_endpoint` (TypeList, Computed): nested structure matching existing `custom_endpoint` schema
- `s3` (TypeList, Computed): nested structure matching existing `s3` schema
- `create_time` (TypeString, Computed)
- `update_time` (TypeString, Computed)

#### Scenario: Read operation populates realtime_log_delivery_tasks
- **WHEN** the resource read operation calls `DescribeRealtimeLogDeliveryTasks` API
- **THEN** the `realtime_log_delivery_tasks` field SHALL be populated with all task objects from `response.Response.RealtimeLogDeliveryTasks`

#### Scenario: Empty result set
- **WHEN** no realtime log delivery tasks match the query
- **THEN** `realtime_log_delivery_tasks` SHALL be set to an empty list

### Requirement: Service layer supports filter-based query
The TeoService SHALL provide a method to query realtime log delivery tasks with custom filters, returning the full list of matching tasks.

#### Scenario: Query with custom filters
- **WHEN** the service method is called with ZoneId and custom AdvancedFilter list
- **THEN** it SHALL pass the filters to `DescribeRealtimeLogDeliveryTasksRequest.Filters` and return all matching tasks from the response

#### Scenario: Query with pagination
- **WHEN** the API response contains more tasks than the Limit
- **THEN** the service method SHALL paginate through all results using Offset and Limit (maximum 1000) until all tasks are retrieved
