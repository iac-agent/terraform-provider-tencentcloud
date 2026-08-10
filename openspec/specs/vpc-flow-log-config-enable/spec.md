## ADDED Requirements

### Requirement: Flow log enable/disable configuration
The system SHALL provide a Terraform resource `tencentcloud_vpc_flow_log_config` that allows users to manage the enable/disable state of an existing VPC flow log via the `enable` parameter.

#### Scenario: Enable a flow log
- **WHEN** user creates a `tencentcloud_vpc_flow_log_config` resource with `enable = true` and a valid `flow_log_id`
- **THEN** the system calls `EnableFlowLogs` API with the specified `flow_log_id`
- **AND** the flow log is enabled

#### Scenario: Disable a flow log
- **WHEN** user creates a `tencentcloud_vpc_flow_log_config` resource with `enable = false` and a valid `flow_log_id`
- **THEN** the system calls `DisableFlowLogs` API with the specified `flow_log_id`
- **AND** the flow log is disabled

#### Scenario: Read flow log enable state
- **WHEN** the system reads the state of a `tencentcloud_vpc_flow_log_config` resource
- **THEN** the system calls `DescribeFlowLogs` API with the `flow_log_id`
- **AND** the `enable` parameter in Terraform state is synchronized with the `Enable` field from the API response

#### Scenario: Update flow log enable state from true to false
- **WHEN** user updates `enable` from `true` to `false`
- **THEN** the system calls `DisableFlowLogs` API with the `flow_log_id`
- **AND** the flow log is disabled

#### Scenario: Update flow log enable state from false to true
- **WHEN** user updates `enable` from `false` to `true`
- **THEN** the system calls `EnableFlowLogs` API with the `flow_log_id`
- **AND** the flow log is enabled

#### Scenario: Delete flow log config (remove from state)
- **WHEN** user destroys a `tencentcloud_vpc_flow_log_config` resource
- **THEN** the system removes the resource from Terraform state without calling any cloud API
- **AND** the underlying flow log instance is NOT deleted

#### Scenario: Import existing flow log config
- **WHEN** user imports a flow log config using `terraform import tencentcloud_vpc_flow_log_config.foo <flow_log_id>`
- **THEN** the system calls `DescribeFlowLogs` API to read the current enable state
- **AND** the resource is added to Terraform state with the correct `enable` value

#### Scenario: Handle flow log not found during read
- **WHEN** `DescribeFlowLogs` API returns an empty result for the specified `flow_log_id`
- **THEN** the system sets the resource ID to empty string (removing it from state)
- **AND** logs the event for troubleshooting
