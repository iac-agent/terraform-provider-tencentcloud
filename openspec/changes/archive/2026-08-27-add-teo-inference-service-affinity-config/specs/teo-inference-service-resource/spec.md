## ADDED Requirements

### Requirement: Create TEO Inference Service with Affinity Config
The system SHALL provide a Terraform resource `tencentcloud_teo_inference_service_v1` that supports creating a TEO inference service with optional affinity configuration.

#### Scenario: Create inference service with full affinity config
- **WHEN** user creates a `tencentcloud_teo_inference_service_v1` resource with `affinity_config` block containing `switch`, `affinity_mode`, `source`, and `header_name`
- **THEN** the system calls `CreateInferenceService` API with `AffinityConfig` populated and returns the created `ServiceId`

#### Scenario: Create inference service without affinity config
- **WHEN** user creates a `tencentcloud_teo_inference_service_v1` resource without `affinity_config` block
- **THEN** the system calls `CreateInferenceService` API with `AffinityConfig` set to nil and returns the created `ServiceId`

### Requirement: Read TEO Inference Service with Affinity Config
The system SHALL support reading a TEO inference service and returning its affinity configuration when available.

#### Scenario: Read inference service with affinity config
- **WHEN** user runs `terraform refresh` on a `tencentcloud_teo_inference_service_v1` resource that has affinity config set
- **THEN** the system calls `DescribeInferenceServices` API and sets `affinity_config` fields from the API response

#### Scenario: Read inference service where affinity config is nil in API response
- **WHEN** user runs `terraform refresh` on a `tencentcloud_teo_inference_service_v1` resource and the API response does not contain `AffinityConfig`
- **THEN** the system preserves the existing state values for `affinity_config` without overwriting them

### Requirement: Update TEO Inference Service Affinity Config
The system SHALL support updating the affinity configuration of an existing TEO inference service.

#### Scenario: Update affinity config switch from Off to On
- **WHEN** user updates `affinity_config.switch` from `Off` to `On` in a `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `ModifyInferenceService` API with the updated `AffinityConfig`

#### Scenario: Update affinity config with session ID settings
- **WHEN** user updates `affinity_config.source` and `affinity_config.header_name` in a `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `ModifyInferenceService` API with the updated `AffinityConfig` containing `SessionIdAffinityConfig`

### Requirement: Delete TEO Inference Service
The system SHALL support deleting a TEO inference service via the `OperateInferenceService` API.

#### Scenario: Delete inference service
- **WHEN** user runs `terraform destroy` on a `tencentcloud_teo_inference_service_v1` resource
- **THEN** the system calls `OperateInferenceService` API with `Operation="Delete"` and removes the resource from state

### Requirement: Import Existing TEO Inference Service
The system SHALL support importing an existing TEO inference service into Terraform state.

#### Scenario: Import inference service by ServiceId
- **WHEN** user runs `terraform import tencentcloud_teo_inference_service_v1.example <service-id>`
- **THEN** the system reads the service via `DescribeInferenceServices` API and populates the state with all available attributes