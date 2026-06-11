## ADDED Requirements

### Requirement: Resource schema definition
The resource `tencentcloud_teo_function_replica` SHALL define the following schema fields:

**Required fields:**
- `zone_id` (TypeString, ForceNew): Site ID
- `function_id` (TypeString, ForceNew): Function ID
- `replica_name` (TypeString, ForceNew): Edge function replica name (1-50 chars, a-z/0-9/-, no consecutive/leading/trailing hyphens, unique per FunctionId)
- `content` (TypeString): Edge function replica content (JavaScript, max 5MB)

**Optional fields:**
- `remark` (TypeString): Edge function replica description (max 50 chars)

**Computed fields:**
- `created_on` (TypeString): Creation time
- `modified_on` (TypeString): Last modification time

The resource SHALL use composite ID `zone_id#function_id#replica_name` with `tccommon.FILED_SP` as separator.

#### Scenario: Schema defines all required and optional fields
- **WHEN** the resource schema is initialized
- **THEN** it SHALL contain `zone_id`, `function_id`, `replica_name`, `content` as required fields, `remark` as optional, and `created_on`, `modified_on` as computed fields

#### Scenario: Composite ID uses three-part format
- **WHEN** a resource is created
- **THEN** the resource ID SHALL be `zone_id#function_id#replica_name` joined by `tccommon.FILED_SP`

### Requirement: Create function replica
The resource SHALL call `CreateFunctionReplica` API to create an edge function replica. After the API call, the resource SHALL set its composite ID and call Read to refresh state.

#### Scenario: Successful creation
- **WHEN** `resourceTencentCloudTeoFunctionReplicaCreate` is called with valid `zone_id`, `function_id`, `replica_name`, and `content`
- **THEN** it SHALL call `CreateFunctionReplica` with the provided parameters
- **AND** set the composite ID to `zone_id#function_id#replica_name`
- **AND** call Read to refresh state

#### Scenario: API returns nil response
- **WHEN** `CreateFunctionReplica` returns a nil response
- **THEN** the create function SHALL return `NonRetryableError` with a descriptive message

### Requirement: Read function replica
The resource SHALL call `DescribeFunctionReplicas` API with `AdvancedFilter` filtering by `replica-name` to read the current state of a function replica. If the replica is not found, the resource SHALL be removed from state.

#### Scenario: Successful read
- **WHEN** `resourceTencentCloudTeoFunctionReplicaRead` is called
- **THEN** it SHALL split the composite ID into `zone_id`, `function_id`, `replica_name`
- **AND** call `DescribeFunctionReplicas` with `ZoneId`, `FunctionId`, and `AdvancedFilter{Name: "replica-name", Values: [replica_name], Fuzzy: false}`
- **AND** set all schema fields from the first matching `FunctionReplica` in the response

#### Scenario: Replica not found
- **WHEN** `DescribeFunctionReplicas` returns an empty list for the given replica name
- **THEN** it SHALL log `[CRUD] teo_function_replica id=<id>` and set `d.SetId("")` to remove from state

#### Scenario: Nil response field check
- **WHEN** a field in the `FunctionReplica` response struct is nil
- **THEN** the Read function SHALL skip calling `d.Set()` for that field

### Requirement: Update function replica
The resource SHALL call `ModifyFunctionReplica` API to update the `content` and `remark` fields of a function replica. The `zone_id`, `function_id`, and `replica_name` fields SHALL be immutable.

#### Scenario: Content or remark changed
- **WHEN** `resourceTencentCloudTeoFunctionReplicaUpdate` is called and `content` or `remark` has changed
- **THEN** it SHALL call `ModifyFunctionReplica` with `ZoneId`, `FunctionId`, `ReplicaName`, and the changed fields
- **AND** call Read to refresh state

#### Scenario: Immutable field change attempted
- **WHEN** an immutable field (`zone_id`, `function_id`, `replica_name`) has changed
- **THEN** the update function SHALL NOT call the Modify API (these are ForceNew and handled by Terraform's recreation)

### Requirement: Delete function replica
The resource SHALL call `DeleteFunctionReplica` API to delete a function replica. It SHALL pass the `replica_name` as a single-element list to `ReplicaNames`.

#### Scenario: Successful deletion
- **WHEN** `resourceTencentCloudTeoFunctionReplicaDelete` is called
- **THEN** it SHALL split the composite ID into `zone_id`, `function_id`, `replica_name`
- **AND** call `DeleteFunctionReplica` with `ZoneId`, `FunctionId`, and `ReplicaNames: [replica_name]`

### Requirement: Resource import support
The resource SHALL support import via the composite ID `zone_id#function_id#replica_name`.

#### Scenario: Import with valid composite ID
- **WHEN** `terraform import` is called with ID `zone-123#func-456#my-replica`
- **THEN** the resource SHALL split the ID and call Read to populate state

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` as `tencentcloud_teo_function_replica` and listed in `tencentcloud/provider.md`.

#### Scenario: Resource is registered
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_function_replica` SHALL be available in the `ResourcesMap`

### Requirement: Retry and error handling
All API calls in CRUD functions SHALL use `resource.Retry` with appropriate timeout (`WriteRetryTimeout` for Create/Update/Delete, `ReadRetryTimeout` for Read via service layer). Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Transient API error
- **WHEN** an API call fails with a retryable error (e.g., network error, rate limit)
- **THEN** the retry mechanism SHALL attempt the call again within the timeout period

#### Scenario: Non-retryable API error
- **WHEN** an API call fails with a non-retryable error
- **THEN** the function SHALL return the error immediately

### Requirement: Resource documentation
A markdown documentation file SHALL be created at `tencentcloud/services/teo/resource_tc_teo_function_replica.md` following the provider documentation format.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `.md` file SHALL exist with a one-line description mentioning TEO, example usage, and import section

### Requirement: Unit tests
Unit tests SHALL be written using gomonkey mocks for all CRUD functions in `resource_tc_teo_function_replica_test.go`.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** `go test -gcflags=all=-l` is run on the test file
- **THEN** tests for Create, Read, Update, and Delete SHALL pass using mocked API responses
