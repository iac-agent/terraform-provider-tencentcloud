## ADDED Requirements

### Requirement: Create TEO function replica v3
The system SHALL allow users to create a TEO edge function replica by specifying `zone_id`, `function_id`, `replica_name`, `content` (all required) and optionally `remark`, by calling the `CreateFunctionReplica` API with retry using `tccommon.WriteRetryTimeout`. Upon success, the resource ID SHALL be set to the composite `zone_id#function_id#replica_name` joined by `tccommon.FILED_SP` (`#`).

#### Scenario: Successful creation with all parameters
- **WHEN** user provides zone_id, function_id, replica_name, content, and remark in the Terraform configuration
- **THEN** the system calls CreateFunctionReplica API with all parameters and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: Successful creation without optional remark
- **WHEN** user provides zone_id, function_id, replica_name, and content without remark
- **THEN** the system calls CreateFunctionReplica API without remark and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: Create API returns empty response
- **WHEN** the CreateFunctionReplica API returns nil response or nil Response body
- **THEN** the system SHALL return a NonRetryableError instead of setting an empty ID

#### Scenario: API returns error during creation
- **WHEN** the CreateFunctionReplica API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError and return the error if retries are exhausted

### Requirement: Read TEO function replica v3
The system SHALL read the current state of a TEO edge function replica by calling the `DescribeFunctionReplicas` API with zone_id, function_id (parsed from the resource ID), `Limit` set to 200 (the maximum value documented in the cloud API), and a filter with `Name=replica-name`, `Values=[replica_name]` to locate the specific replica, then set content, remark, created_on, modified_on, and the flattened `function_replicas` computed list from the matching replica.

#### Scenario: Successful read of existing replica
- **WHEN** the resource exists and DescribeFunctionReplicas returns a matching replica (exact replica_name match)
- **THEN** the system SHALL set zone_id, function_id, replica_name from the parsed ID, and content, remark, created_on, modified_on, function_replicas from the response, each guarded by a nil check before calling d.Set

#### Scenario: Replica not found during read
- **WHEN** DescribeFunctionReplicas returns no matching replica for the given replica_name, or the response is empty
- **THEN** the system SHALL first log `[CRUD] ... id=<resource id>` to preserve the incident context, then remove the resource from state (d.SetId(""))

#### Scenario: Read maps user query parameters to request
- **WHEN** the user configured sort_by, sort_order, or filters in the resource configuration
- **THEN** the Read function SHALL map them to the request's SortBy, SortOrder, and Filters fields respectively, while still filtering by replica-name to locate the single replica

#### Scenario: API returns error during read
- **WHEN** the DescribeFunctionReplicas API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError with tccommon.ReadRetryTimeout

### Requirement: Update TEO function replica v3
The system SHALL allow users to update `content` and/or `remark` of an existing TEO edge function replica by calling the `ModifyFunctionReplica` API with zone_id, function_id, replica_name (parsed from the resource ID) and the changed mutable fields, with retry using `tccommon.WriteRetryTimeout`. Changes to sort_by, sort_order, or filters SHALL NOT trigger a Modify call because they are query-only parameters not accepted by the Modify API.

#### Scenario: Update content only
- **WHEN** user changes the content field in Terraform configuration
- **THEN** the system calls ModifyFunctionReplica API with the new content value

#### Scenario: Update remark only
- **WHEN** user changes the remark field in Terraform configuration
- **THEN** the system calls ModifyFunctionReplica API with the new remark value

#### Scenario: Query-only parameter change does not trigger Modify
- **WHEN** user changes only sort_by, sort_order, or filters
- **THEN** the system SHALL skip the ModifyFunctionReplica API call and re-read the state

#### Scenario: API returns error during update
- **WHEN** the ModifyFunctionReplica API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError and return the error if retries are exhausted

### Requirement: Delete TEO function replica v3
The system SHALL delete a TEO edge function replica by calling the `DeleteFunctionReplica` API with zone_id, function_id (parsed from the resource ID), and the `replica_names` list from the resource configuration mapped to `ReplicaNames`, with retry using `tccommon.WriteRetryTimeout`. If `replica_names` is not set in configuration, the system SHALL fall back to a single-element list containing the replica_name parsed from the resource ID.

#### Scenario: Successful deletion with configured replica_names
- **WHEN** the resource is destroyed and replica_names is configured
- **THEN** the system calls DeleteFunctionReplica API with ReplicaNames containing the configured values

#### Scenario: Deletion falls back to ID-parsed replica name
- **WHEN** replica_names is not configured in state
- **THEN** the system calls DeleteFunctionReplica API with ReplicaNames containing only the replica_name parsed from the composite ID

#### Scenario: API returns error during deletion
- **WHEN** the DeleteFunctionReplica API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError and return the error if retries are exhausted

### Requirement: Schema covers all mapped API parameters
The resource schema SHALL include: `zone_id`, `function_id`, `replica_name` (Required, ForceNew, String); `content` (Required, String); `remark` (Optional, String); `sort_by`, `sort_order` (Optional, String); `filters` (Optional, list of objects with `name` Required String, `values` Required list of String, `fuzzy` Optional Bool); `replica_names` (Required, list of String); and computed attributes `created_on`, `modified_on` (String) plus `function_replicas` (flattened computed list of objects with function_id, replica_name, content, remark, created_on, modified_on). The schema SHALL NOT wrap the replica list in an extra nested container layer.

#### Scenario: Schema field coverage
- **WHEN** the resource schema is defined
- **THEN** all fields listed in the API parameter mapping SHALL be present with the specified types and required/optional/computed attributes

#### Scenario: Flattened computed list
- **WHEN** the DescribeFunctionReplicas response returns FunctionReplicas
- **THEN** the system SHALL set the top-level `function_replicas` computed list with each element's fields flattened, without an additional xxx_set/xxx_list wrapper layer

### Requirement: Import TEO function replica v3
The system SHALL support importing an existing TEO edge function replica using the composite ID format `zone_id#function_id#replica_name` via ImportStatePassthrough.

#### Scenario: Successful import with composite ID
- **WHEN** user runs terraform import with ID in format `zone_id#function_id#replica_name`
- **THEN** the system SHALL pass the ID through to state and call Read to populate remaining fields

#### Scenario: Invalid import ID format
- **WHEN** user provides an import ID that does not contain exactly 3 parts separated by `#`
- **THEN** the system SHALL return an error indicating the expected format

### Requirement: ForceNew on identity fields
The system SHALL force resource recreation when `zone_id`, `function_id`, or `replica_name` are changed, as the ModifyFunctionReplica API does not support modifying these identity fields (replica_name is the locator; zone_id and function_id are ownership identifiers).

#### Scenario: Change replica_name triggers recreation
- **WHEN** user changes the replica_name in Terraform configuration
- **THEN** Terraform SHALL plan to destroy the old resource and create a new one

#### Scenario: Change zone_id triggers recreation
- **WHEN** user changes the zone_id in Terraform configuration
- **THEN** Terraform SHALL plan to destroy the old resource and create a new one

### Requirement: Provider registration for function replica v3
The system SHALL register the `tencentcloud_teo_function_replica_v3` resource in `tencentcloud/provider.go` and list it in `tencentcloud/provider.md`.

#### Scenario: Resource is available after provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_function_replica_v3` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey mocks
The system SHALL provide unit tests in `resource_tc_teo_function_replica_v3_test.go` using gomonkey to mock the cloud API client methods (CreateFunctionReplicaWithContext, DescribeFunctionReplicasWithContext, ModifyFunctionReplicaWithContext, DeleteFunctionReplicaWithContext), covering Create, Read (found and not-found), Update, and Delete operations, without using the Terraform acceptance test suite.

#### Scenario: Create unit test
- **WHEN** the Create unit test runs with mocked CreateFunctionReplicaWithContext and DescribeFunctionReplicasWithContext
- **THEN** the test asserts request parameters, successful completion, and the composite ID format

#### Scenario: Read not-found unit test
- **WHEN** the mocked DescribeFunctionReplicasWithContext returns an empty replica list
- **THEN** the test asserts the resource ID is cleared

#### Scenario: Update and Delete unit tests
- **WHEN** the Update/Delete unit tests run with mocked Modify/Delete API methods
- **THEN** the tests assert the Modify request fields and the ReplicaNames parameter passed to DeleteFunctionReplica
