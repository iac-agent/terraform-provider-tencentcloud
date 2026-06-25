# teo-function-replica-resource Specification

## Purpose
EdgeOne Function Replica resource management via Terraform provider for tencentcloud.

## Requirements
### Requirement: TFR-001 - TEO Function Replica resource CRUD
**Priority**: High  
**Type**: Functional

The system SHALL provide a `tencentcloud_teo_function_replica` Terraform resource that manages the full lifecycle of EdgeOne function replicas via CreateFunctionReplica, DescribeFunctionReplicas, ModifyFunctionReplica, and DeleteFunctionReplica APIs.

**Acceptance Criteria**:
- Resource supports Create, Read, Update, Delete operations
- Resource ID uses composite format `zone_id#function_id#replica_name`
- Schema includes `zone_id`, `function_id`, `replica_name`, `content`, `remark`, `created_on`, `modified_on`

#### Scenario: Create a function replica
**Given** a user defines a `tencentcloud_teo_function_replica` resource  
**When** `terraform apply` is executed  
**Then** the provider calls `CreateFunctionReplica` API  
**And** sets the resource ID to `zone_id#function_id#replica_name`  
**And** the replica is created successfully

#### Scenario: Read a function replica
**Given** an existing `tencentcloud_teo_function_replica` resource  
**When** Terraform reads the resource state  
**Then** it calls `DescribeFunctionReplicas` with `Filters: [{Name: "replica-name", Values: [replica_name]}]`  
**And** populates `content`, `remark`, `created_on`, `modified_on` from the matched result

#### Scenario: Update a function replica
**Given** an existing `tencentcloud_teo_function_replica` resource  
**When** a user modifies `content` or `remark`  
**Then** the provider calls `ModifyFunctionReplica` API with the updated values  
**And** the changes are applied successfully

#### Scenario: Delete a function replica
**Given** an existing `tencentcloud_teo_function_replica` resource  
**When** a user destroys the resource  
**Then** the provider calls `DeleteFunctionReplica` API with `zone_id`, `function_id`, and `ReplicaNames: [replica_name]`  
**And** the replica is deleted successfully

#### Scenario: Resource not found on read
**Given** a `tencentcloud_teo_function_replica` resource that no longer exists  
**When** Terraform reads the resource state  
**Then** `DescribeFunctionReplicas` returns no matching replica  
**And** the provider sets the resource ID to empty string  
**And** marks the resource as removed from state

---

### Requirement: TFR-002 - TEO Function Replica resource import
**Priority**: High  
**Type**: Functional

The system SHALL support importing existing function replicas using the composite ID format `zone_id#function_id#replica_name`.

**Acceptance Criteria**:
- Import works via `terraform import tencentcloud_teo_function_replica.foo zone-xxx#function-yyy#my-replica`
- After import, replica state is read via `DescribeFunctionReplicas`
- All computed fields are populated

#### Scenario: Import a function replica
**Given** a function replica exists in Tencent Cloud  
**When** user runs `terraform import tencentcloud_teo_function_replica.foo zone-xxx#function-yyy#my-replica`  
**Then** the provider reads the function replica state via `DescribeFunctionReplicas`  
**And** populates all computed fields  
**And** `terraform plan` shows no changes

---

### Requirement: TFR-003 - TEO Function Replica force-new on identity fields
**Priority**: High  
**Type**: Functional

The system SHALL mark `zone_id`, `function_id`, and `replica_name` as ForceNew, causing resource recreation when any of these fields change.

**Acceptance Criteria**:
- `zone_id` is marked ForceNew
- `function_id` is marked ForceNew
- `replica_name` is marked ForceNew
- Changing any identity field triggers resource recreation

#### Scenario: Identity field change triggers recreation
**Given** an existing `tencentcloud_teo_function_replica` resource  
**When** a user changes `replica_name`  
**Then** Terraform destroys the old replica  
**And** creates a new one with the new name