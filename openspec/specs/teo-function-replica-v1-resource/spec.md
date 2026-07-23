## Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_function_replica_v1` SHALL define the following schema fields:

**Required fields:**
- `zone_id` (string, ForceNew): 站点 ID
- `function_id` (string, ForceNew): 函数 ID
- `replica_name` (string, ForceNew): 边缘函数副本名称
- `content` (string): 边缘函数副本内容
- `remark` (string): 边缘函数副本描述，最大支持 50 个字符

**Computed fields:**
- `create_time` (string): 创建时间
- `update_time` (string): 更新时间

The resource ID SHALL be composed of `zone_id#function_id#replica_name` using `tccommon.FILED_SP` as separator.

#### Scenario: Create resource with all required fields
- **WHEN** user creates a `tencentcloud_teo_function_replica_v1` resource with zone_id, function_id, replica_name, content, and remark
- **THEN** the resource SHALL be created with a composite ID of `zone_id#function_id#replica_name`

#### Scenario: Missing remark field
- **WHEN** user attempts to create a `tencentcloud_teo_function_replica_v1` resource without providing remark
- **THEN** Terraform SHALL reject the configuration with a required field error

### Requirement: Resource Create
The Create function SHALL call `CreateFunctionReplica` API with ZoneId, FunctionId, ReplicaName, Content, and Remark parameters. After successful creation, the resource ID SHALL be set to `zone_id#function_id#replica_name`.

The Create function SHALL check that the API response is not empty and that key fields are not nil. If empty, it SHALL return a NonRetryableError.

#### Scenario: Successful creation
- **WHEN** CreateFunctionReplica API returns successfully
- **THEN** the resource SHALL set its ID to `zone_id#function_id#replica_name` and call Read to refresh state

#### Scenario: API returns empty response
- **WHEN** CreateFunctionReplica API returns nil response or empty fields
- **THEN** the resource SHALL return a NonRetryableError

### Requirement: Resource Read
The Read function SHALL call `DescribeFunctionReplicas` API with ZoneId and FunctionId, using Filters with replica-name to locate the specific replica. If the replica is found in the response list, it SHALL set all schema fields. If the response is empty or the replica is not found, it SHALL log the id and then call `d.SetId("")`.

#### Scenario: Replica exists
- **WHEN** DescribeFunctionReplicas returns a list containing the matching replica
- **THEN** the resource SHALL set zone_id, function_id, replica_name, content, remark, create_time, and update_time from the response

#### Scenario: Replica not found
- **WHEN** DescribeFunctionReplicas returns an empty list or the replica is not in the list
- **THEN** the resource SHALL log `[CRUD] function_replica_v1 id=%s` with the current ID and then call `d.SetId("")`

### Requirement: Resource Update
The Update function SHALL call `ModifyFunctionReplica` API when content or remark fields change. The `zone_id`, `function_id`, and `replica_name` fields SHALL be immutable (ForceNew or checked in update).

#### Scenario: Update content or remark
- **WHEN** user changes the content or remark field
- **THEN** the resource SHALL call ModifyFunctionReplica with ZoneId, FunctionId, ReplicaName, and the updated Content/Remark

#### Scenario: Attempt to change immutable fields
- **WHEN** user attempts to change zone_id, function_id, or replica_name
- **THEN** Terraform SHALL force resource recreation (ForceNew) or return an error

### Requirement: Resource Delete
The Delete function SHALL call `DeleteFunctionReplica` API with ZoneId, FunctionId, and ReplicaNames (as a single-element list containing the replica_name).

#### Scenario: Successful deletion
- **WHEN** DeleteFunctionReplica API returns successfully
- **THEN** the resource SHALL be removed from state

#### Scenario: API error during deletion
- **WHEN** DeleteFunctionReplica API returns an error
- **THEN** the resource SHALL retry using tccommon.WriteRetryTimeout and return the error if retries are exhausted

### Requirement: Provider Registration
The resource `tencentcloud_teo_function_replica_v1` SHALL be registered in `tencentcloud/provider.go` with the function `teo.ResourceTencentCloudTeoFunctionReplicaV1()` and listed in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_function_replica_v1` SHALL be available for use in Terraform configurations

### Requirement: Documentation
The resource SHALL have a corresponding `.md` documentation file at `tencentcloud/services/teo/resource_tc_teo_function_replica_v1.md` with a one-line description, example usage, and import section.

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** a markdown file with description, example usage, and import instructions SHALL exist

### Requirement: Unit Tests
The resource SHALL have unit tests at `tencentcloud/services/teo/resource_tc_teo_function_replica_v1_test.go` using gomonkey to mock cloud API calls, covering Create, Read, Update, and Delete operations.

#### Scenario: Unit tests pass
- **WHEN** running `go test -gcflags=all=-l` on the test file
- **THEN** all test cases for Create, Read, Update, and Delete SHALL pass
