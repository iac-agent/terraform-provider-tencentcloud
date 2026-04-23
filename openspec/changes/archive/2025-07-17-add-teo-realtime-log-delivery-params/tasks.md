## 1. Schema Definition

- [x] 1.1 Add `filters` Optional TypeList parameter to `ResourceTencentCloudTeoRealtimeLogDelivery` schema with nested `name` (Required, TypeString), `values` (Required, TypeList of TypeString), `fuzzy` (Optional, TypeBool) fields
- [x] 1.2 Add `realtime_log_delivery_tasks` Computed TypeList parameter to `ResourceTencentCloudTeoRealtimeLogDelivery` schema with nested computed fields: task_id, task_name, delivery_status, task_type, entity_list, log_type, area, fields, custom_fields, delivery_conditions, sample, log_format, cls, custom_endpoint, s3, create_time, update_time

## 2. Service Layer

- [x] 2.1 Add `DescribeTeoRealtimeLogDeliveryTasksByFilters` method to `TeoService` in `service_tencentcloud_teo.go` that accepts ZoneId and custom AdvancedFilter list, supports pagination (Limit=1000), and returns the full list of RealtimeLogDeliveryTask objects

## 3. CRUD Function Updates

- [x] 3.1 Update `resourceTencentCloudTeoRealtimeLogDeliveryRead` function to populate `filters` from user configuration when set, and to populate `realtime_log_delivery_tasks` computed output from the API response
- [x] 3.2 When `filters` is not set in configuration, preserve existing read behavior (filter by task-id from composite ID)

## 4. Unit Tests

- [x] 4.1 Add unit tests for `filters` schema parameter in `resource_tc_teo_realtime_log_delivery_test.go` using gomonkey mock
- [x] 4.2 Add unit tests for `realtime_log_delivery_tasks` computed output in `resource_tc_teo_realtime_log_delivery_test.go` using gomonkey mock
- [x] 4.3 Add unit tests for `DescribeTeoRealtimeLogDeliveryTasksByFilters` service method with pagination using gomonkey mock
- [x] 4.4 Run `go test` to verify all new unit tests pass

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_realtime_log_delivery.md` example file to include `filters` and `realtime_log_delivery_tasks` parameters
