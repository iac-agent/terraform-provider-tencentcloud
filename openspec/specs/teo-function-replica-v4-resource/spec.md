# teo-function-replica-v4-resource Specification

## Purpose
TBD - created by archiving change add-teo-function-replica-v4-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_teo_function_replica_v4` that manages a single Tencent Cloud TEO edge function replica per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the teo resource map (adjacent to the existing `tencentcloud_teo_function_replica` entry), and MUST be listed in the TEO Resources section of `tencentcloud/provider.md` so that `make doc` generates the website documentation.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration referencing `resource "tencentcloud_teo_function_replica_v4" "<name>"`
- **THEN** Terraform resolves the resource type without an "unknown resource" error

#### Scenario: Existing v1 resource remains unaffected
- **WHEN** the new v4 resource is added to the provider
- **THEN** the existing `tencentcloud_teo_function_replica` resource registration, schema, and behavior SHALL remain unchanged

#### Scenario: Provider compiles
- **WHEN** the codebase is built
- **THEN** the build succeeds with no compilation errors related to the new resource

### Requirement: Resource schema definition
The resource schema SHALL define the following fields:
- `zone_id` (TypeString, Required, ForceNew): Zone ID
- `function_id` (TypeString, Required, ForceNew): Function ID
- `replica_name` (TypeString, Required, ForceNew): Edge function replica name (1-50 chars, a-z/0-9/-, unique under the same FunctionId)
- `content` (TypeString, Required): Edge function replica content (JavaScript code, max 5MB)
- `remark` (TypeString, Optional): Edge function replica description (max 50 chars)
- `created_on` (TypeString, Computed): Replica creation time, hydrated from DescribeFunctionReplicas
- `modified_on` (TypeString, Computed): Replica update time, hydrated from DescribeFunctionReplicas

The schema SHALL NOT expose query-helper parameters of DescribeFunctionReplicas (sort_by, sort_order, filters, offset, limit) as resource fields.

#### Scenario: Schema mirrors API parameters
- **WHEN** a developer inspects the resource schema
- **THEN** every CreateFunctionReplica input field (ZoneId, FunctionId, ReplicaName, Content, Remark) maps to a schema field, and every FunctionReplica response field (CreatedOn, ModifiedOn) is exposed as Computed

#### Scenario: Identity fields force recreation
- **WHEN** `zone_id`, `function_id`, or `replica_name` is changed in the Terraform configuration
- **THEN** Terraform SHALL plan destroy-and-recreate (ForceNew)

#### Scenario: Mutable fields support in-place update
- **WHEN** `content` or `remark` is changed in the Terraform configuration
- **THEN** Terraform SHALL plan an in-place update without recreation

### Requirement: Composite resource ID
The resource ID SHALL be the 3-segment composite `<zone_id>#<function_id>#<replica_name>` joined with `tccommon.FILED_SP`. The resource SHALL support `terraform import` with this composite ID (ImportStatePassthrough).

#### Scenario: Create sets the composite ID
- **WHEN** CreateFunctionReplica succeeds
- **THEN** the resource calls `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))` outside the retry block

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_teo_function_replica_v4.x zone-xxx#ef-yyy#replica-zzz`
- **THEN** the resource state is hydrated from DescribeFunctionReplicas using the parsed 3-tuple

#### Scenario: Malformed ID
- **WHEN** the resource ID does not contain exactly 3 segments
- **THEN** the CRUD functions SHALL return a descriptive error (e.g. "id is broken") before any SDK call

### Requirement: Create operation
On Create, the resource SHALL call `CreateFunctionReplica` with ZoneId, FunctionId, ReplicaName, Content, and (if configured) Remark, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with errors wrapped via `tccommon.RetryError`. When the response or its Response field is nil, the resource SHALL return a NonRetryableError instead of proceeding. After a successful call, the resource SHALL set the composite ID and invoke Read.

#### Scenario: Successful creation with all parameters
- **WHEN** the user provides zone_id, function_id, replica_name, content, and remark
- **THEN** the system calls CreateFunctionReplica with all five parameters and sets the composite ID

#### Scenario: Successful creation without optional remark
- **WHEN** the user omits remark
- **THEN** the Remark field of the request remains unset

#### Scenario: Nil response treated as failure
- **WHEN** CreateFunctionReplica returns a nil result or nil Response
- **THEN** the system returns a NonRetryableError rather than setting the resource ID

#### Scenario: API error retried
- **WHEN** CreateFunctionReplica returns a retryable error
- **THEN** the system retries within tccommon.WriteRetryTimeout and returns the wrapped error if retries are exhausted

### Requirement: Read operation
On Read, the resource SHALL call `DescribeFunctionReplicas` with ZoneId and FunctionId parsed from the ID, a Filters entry of Name="replica-name" and Values=[replicaName], and Limit=200 (the documented maximum), wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. The resource SHALL locate the replica by exact match on ReplicaName within the response list. When found, it SHALL set zone_id, function_id, replica_name, and (after nil checks) content, remark, created_on, and modified_on. When the response is empty or no matching replica exists, it SHALL first log `[CRUD]` with the current ID and then call `d.SetId("")`, returning nil.

#### Scenario: Successful read populates all fields
- **WHEN** DescribeFunctionReplicas returns a replica whose ReplicaName exactly matches
- **THEN** the resource sets content, remark, created_on, and modified_on from the response (each guarded by a nil check)

#### Scenario: Replica not found clears state
- **WHEN** the response contains no replica with a matching ReplicaName
- **THEN** the resource logs the `[CRUD]` line including the current ID, calls `d.SetId("")`, and returns nil

#### Scenario: Read uses maximum page size and replica-name filter
- **WHEN** the DescribeFunctionReplicas request is built
- **THEN** Limit is 200 and Filters contains Name="replica-name" with Values=[replicaName]

### Requirement: Update operation
On Update, the resource SHALL call `ModifyFunctionReplica` with ZoneId, FunctionId, and ReplicaName parsed from the ID, plus Content and Remark (when configured), wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. The Modify call SHALL be skipped when neither `content` nor `remark` has changed. After a successful update, the resource SHALL invoke Read.

#### Scenario: Update content and remark
- **WHEN** both content and remark change
- **THEN** the system calls ModifyFunctionReplica once with both new values

#### Scenario: Update single field
- **WHEN** only content (or only remark) changes
- **THEN** the system calls ModifyFunctionReplica with the updated value

#### Scenario: No-op update
- **WHEN** no mutable field (content, remark) has changed
- **THEN** the system skips the ModifyFunctionReplica call and directly invokes Read

### Requirement: Delete operation
On Delete, the resource SHALL call `DeleteFunctionReplica` with ZoneId and FunctionId parsed from the ID, and ReplicaNames containing exactly one element (the replica_name parsed from the ID), wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Successful deletion
- **WHEN** the resource is destroyed
- **THEN** the system calls DeleteFunctionReplica with ReplicaNames containing the single replica_name

#### Scenario: API error retried
- **WHEN** DeleteFunctionReplica returns a retryable error
- **THEN** the system retries within tccommon.WriteRetryTimeout and returns the wrapped error if retries are exhausted

### Requirement: Unit tests
The system SHALL provide unit tests in `resource_tc_teo_function_replica_v4_test.go` using gomonkey to mock the cloud API client methods (not the Terraform acceptance test suite), covering Create, Read (found and not-found), Update, and Delete operations.

#### Scenario: Create test
- **WHEN** the Create test runs with mocked CreateFunctionReplicaWithContext and DescribeFunctionReplicasWithContext
- **THEN** the request parameters are asserted and the composite ID is set to `zone#function#replica`

#### Scenario: Read test asserts computed fields
- **WHEN** the Read test runs against a mocked replica
- **THEN** content, remark, created_on, and modified_on are populated into the resource data

#### Scenario: Read-not-found test
- **WHEN** the mocked DescribeFunctionReplicasWithContext returns an empty list
- **THEN** the resource ID is cleared (d.Id() becomes "")

#### Scenario: Update test
- **WHEN** the Update test runs with changed content and remark
- **THEN** ModifyFunctionReplica is invoked with the updated values

#### Scenario: Delete test
- **WHEN** the Delete test runs
- **THEN** DeleteFunctionReplica is invoked with a single-element ReplicaNames list

### Requirement: Resource documentation
The system SHALL provide `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.md` containing a one-line description mentioning TEO, an Example Usage block, and an Import section documenting the 3-segment composite ID format. The file SHALL NOT contain manually written Argument Reference or Attribute Reference sections (generated by `make doc`).

#### Scenario: Documentation file exists with required sections
- **WHEN** the resource is implemented
- **THEN** the .md file exists with description, example usage, and import sections showing `zone_id#function_id#replica_name` format

