## 1. Resource Schema Definition

- [x] 1.1 Create `resource_tc_teo_inference_service_v1.go` with `ResourceTencentCloudTeoInferenceServiceV1` schema function defining all parameters: `zone_id`, `name`, `listen_port`, `containers`, `resource_config`, `request_paths`, `description`, `operation`, `service_id`, `status`, `inference_url`, `create_time`, `update_time`
- [x] 1.2 Define nested schema for `containers` (TypeList) with: `image_type`, `tcr_repository_config` (nested: `tcr_type`, `image`, `registry_id`, `region_name`), `startup_command`, `environment_variables` (nested: `key`, `value`)
- [x] 1.3 Define nested schema for `resource_config` (TypeList) with: `scaling_mode`, `hardware_spec`, `auto_scaling_config` (nested: `min_instance_count`, `scaling_policies`), `manual_instance_config` (nested: `fixed_instance_count`), `concurrency`
- [x] 1.4 Add Timeouts block to schema for asynchronous operations (Create/Delete/Update)

## 2. Create Operation

- [x] 2.1 Implement `resourceTencentCloudTeoInferenceServiceV1Create` function calling `CreateInferenceService` API
- [x] 2.2 Map Terraform schema fields to `CreateInferenceServiceRequest` parameters (ZoneId, Name, ListenPort, Containers as `InferenceContainerConfig`, ResourceConfig as `InferenceResourceConfig`, RequestPaths, Description)
- [x] 2.3 Handle response: extract `ServiceId` from response and set as `d.SetId()`
- [x] 2.4 Add retry logic with `tccommon.ReadRetryTimeout` and `tccommon.RetryError()`
- [x] 2.5 Validate response is not nil and ServiceId is not empty before setting ID

## 3. Read Operation

- [x] 3.1 Implement `resourceTencentCloudTeoInferenceServiceV1Read` function calling `DescribeInferenceServices` API
- [x] 3.2 Build request with `ZoneId` from state and `Filters` (service-id filter with exact match, Fuzzy=false)
- [x] 3.3 Set pagination `Offset=0` and `Limit=200` (max value)
- [x] 3.4 Handle empty response: log and call `d.SetId("")`
- [x] 3.5 Map response fields back to schema: ServiceId, Name, ListenPort, Containers, ResourceConfig, RequestPaths, Description, Status, InferenceURL, CreateTime, UpdateTime
- [x] 3.6 Add retry logic with `tccommon.ReadRetryTimeout`

## 4. Update Operation

- [x] 4.1 Implement `resourceTencentCloudTeoInferenceServiceV1Update` function
- [x] 4.2 Call `ModifyInferenceService` with mutable fields: ListenPort, RequestPaths, Containers (as `InferenceContainerConfigForModify`), ResourceConfig (as `InferenceResourceConfigForModify` without HardwareSpec), Description
- [x] 4.3 Handle `operation` field: if changed to "Stop" or "Resume", call `OperateInferenceService` after config modification
- [x] 4.4 Add retry logic for both API calls
- [x] 4.5 After update, call Read to refresh state

## 5. Delete Operation

- [x] 5.1 Implement `resourceTencentCloudTeoInferenceServiceV1Delete` function calling `OperateInferenceService` with `Operation="Delete"`
- [x] 5.2 Add retry logic with `tccommon.ReadRetryTimeout`
- [x] 5.3 Poll `DescribeInferenceServices` after delete to confirm service removal (asynchronous operation handling)

## 6. Provider Registration

- [x] 6.1 Register `tencentcloud_teo_inference_service_v1` resource in `tencentcloud/provider.go`
- [x] 6.2 Add resource entry in `tencentcloud/provider.md` documentation

## 7. Documentation

- [x] 7.1 Create `resource_tc_teo_inference_service_v1.md` with Example Usage section showing full configuration
- [x] 7.2 Add Import section documenting import by service_id

## 8. Unit Tests

- [x] 8.1 Create `resource_tc_teo_inference_service_v1_test.go` with unit tests using gomonkey mocks
- [x] 8.2 Test Create, Read, Update, Delete, and Import scenarios

## 9. Finalization

- [x] 9.1 Run `gofmt` on all changed Go files
- [x] 9.2 Run `make doc` to generate website documentation
- [x] 9.3 Create changelog entry in `.changelog/` directory