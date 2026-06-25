## ADDED Requirements

### Requirement: TEO Function Replica resource CRUD
The system SHALL provide a `tencentcloud_teo_function_replica` Terraform resource that manages the full lifecycle of EdgeOne function replicas via CreateFunctionReplica, DescribeFunctionReplicas, ModifyFunctionReplica, and DeleteFunctionReplica APIs.

#### Scenario: Create a function replica
- **WHEN** a user defines a `tencentcloud_teo_function_replica` resource with `zone_id`, `function_id`, `replica_name`, `content`, and optional `remark`
- **THEN** the provider calls `CreateFunctionReplica` API and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: Read a function replica
- **WHEN** the provider reads an existing `tencentcloud_teo_function_replica` resource
- **THEN** it calls `DescribeFunctionReplicas` with `Filters: [{Name: "replica-name", Values: [replica_name]}]` and populates `content`, `remark`, `created_on`, `modified_on` from the matched result

#### Scenario: Update a function replica
- **WHEN** a user modifies `content` or `remark` of an existing `tencentcloud_teo_function_replica` resource
- **THEN** the provider calls `ModifyFunctionReplica` API with the updated values

#### Scenario: Delete a function replica
- **WHEN** a user destroys a `tencentcloud_teo_function_replica` resource
- **THEN** the provider calls `DeleteFunctionReplica` API with `zone_id`, `function_id`, and `ReplicaNames: [replica_name]`

#### Scenario: Resource not found on read
- **WHEN** the provider reads a `tencentcloud_teo_function_replica` resource and `DescribeFunctionReplicas` returns no matching replica
- **THEN** the provider sets the resource ID to empty string, marking it as removed from state

### Requirement: TEO Function Replica resource import
The system SHALL support importing existing function replicas using the composite ID format `zone_id#function_id#replica_name`.

#### Scenario: Import a function replica
- **WHEN** a user runs `terraform import tencentcloud_teo_function_replica.foo zone-xxx#function-yyy#my-replica`
- **THEN** the provider reads the function replica state via `DescribeFunctionReplicas` and populates all computed fields

### Requirement: TEO Function Replica force-new on identity fields
The system SHALL mark `zone_id`, `function_id`, and `replica_name` as ForceNew, causing resource recreation when any of these fields change.

#### Scenario: Replica name change triggers recreation
- **WHEN** a user changes `replica_name` of an existing `tencentcloud_teo_function_replica`
- **THEN** Terraform destroys the old replica and creates a new one with the new name