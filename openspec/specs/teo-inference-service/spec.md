## Requirements

### Requirement: Create inference service
The system SHALL support creating a TEO inference service via the `tencentcloud_teo_inference_service_v1` resource. The creation MUST call `CreateInferenceService` API with `zone_id`, `name`, `listen_port`, `containers`, `resource_config`, and optional `request_paths` and `description` parameters. The API returns a `service_id` which SHALL be set as the Terraform resource ID.

#### Scenario: Successful creation with required fields
- **WHEN** user defines a `tencentcloud_teo_inference_service_v1` resource with `zone_id`, `name`, `listen_port`, `containers`, and `resource_config`
- **THEN** the system calls `CreateInferenceService` API and sets the returned `service_id` as the resource ID

#### Scenario: Creation with optional fields
- **WHEN** user defines a `tencentcloud_teo_inference_service_v1` resource with `request_paths` and `description` in addition to required fields
- **THEN** the system includes all optional parameters in the API call and creates the service successfully

### Requirement: Read inference service
The system SHALL support reading a TEO inference service's current state via `DescribeInferenceServices` API. The read operation MUST filter by `zone_id` and `service-id` to retrieve the specific service. All computed fields (`status`, `inference_url`, `create_time`, `update_time`) SHALL be populated from the API response.

#### Scenario: Successful read of existing service
- **WHEN** the system reads a `tencentcloud_teo_inference_service_v1` resource with a valid `service_id`
- **THEN** the system calls `DescribeInferenceServices` with `zone_id` and `service-id` filter, and sets all schema fields from the response

#### Scenario: Read returns empty result
- **WHEN** the `DescribeInferenceServices` API returns an empty service list
- **THEN** the system logs the missing resource and sets the resource ID to empty (resource is considered removed)

### Requirement: Update inference service
The system SHALL support updating a TEO inference service's configuration via `ModifyInferenceService` API. The update operation MUST pass `zone_id`, `service_id`, and all mutable fields (`listen_port`, `request_paths`, `containers`, `resource_config`, `description`). The `resource_config` in update MUST NOT include `hardware_spec` (not supported by Modify API).

#### Scenario: Update listen port and description
- **WHEN** user modifies `listen_port` and `description` of an existing `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `ModifyInferenceService` with the updated values

#### Scenario: Update containers configuration
- **WHEN** user modifies `containers` of an existing `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system converts containers schema to `InferenceContainerConfigForModify` type and calls `ModifyInferenceService`

### Requirement: Delete inference service
The system SHALL support deleting a TEO inference service via `OperateInferenceService` API with `Operation = "Delete"`. The delete operation is asynchronous and the system MUST poll `DescribeInferenceServices` until the service is no longer found or the deletion is confirmed.

#### Scenario: Successful deletion
- **WHEN** user triggers destroy of a `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `OperateInferenceService` with `Operation = "Delete"` and polls until the service is removed

### Requirement: Operate inference service (Stop/Resume)
The system SHALL support stopping and resuming a TEO inference service via `OperateInferenceService` API with `Operation = "Stop"` or `Operation = "Resume"`. The `operation` field in the schema controls this behavior. When `operation` is set to `Stop` or `Resume`, the system MUST call `OperateInferenceService` after completing any configuration changes.

#### Scenario: Stop a running service
- **WHEN** user sets `operation` to `"Stop"` on an existing `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `OperateInferenceService` with `Operation = "Stop"` and polls until `status` becomes `"Stopped"`

#### Scenario: Resume a stopped service
- **WHEN** user sets `operation` to `"Resume"` on a stopped `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `OperateInferenceService` with `Operation = "Resume"` and polls until `status` becomes `"Running"`

### Requirement: Import existing inference service
The system SHALL support importing an existing TEO inference service into Terraform state. The import ID SHALL be the `service_id`. On import, the system MUST call `DescribeInferenceServices` to populate all resource attributes.

#### Scenario: Import by service ID
- **WHEN** user runs `terraform import tencentcloud_teo_inference_service_v1.foo <service_id>`
- **THEN** the system reads the service by `service_id` and populates the Terraform state with all attributes