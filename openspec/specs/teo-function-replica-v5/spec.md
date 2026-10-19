# teo-function-replica-v5 Specification

## Purpose
TBD - created by archiving change add-teo-function-replica-v5. Update Purpose after archive.
## Requirements
### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_function_replica_v5` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `function_id` (Required, ForceNew, TypeString): 函数 ID
- `replica_name` (Required, ForceNew, TypeString): 边缘函数副本名称，1-50 个字符，允许 a-z、0-9、-，同一 FunctionId 下唯一
- `content` (Required, TypeString): 边缘函数副本内容（JavaScript 代码，最大 5MB）
- `remark` (Optional, TypeString): 边缘函数副本描述，最大 50 个字符
- `sort_by` (Optional, TypeString): 排序依据，默认 created-on（创建时间）
- `sort_order` (Optional, TypeString): 列表排序方式，asc（升序）/ desc（降序），默认 asc
- `filters` (Optional, TypeList of AdvancedFilter): 过滤条件列表，元素结构为 `name`（Required, TypeString）、`values`（Required, TypeSet of TypeString）、`fuzzy`（Optional, TypeBool）
- `replica_names` (Optional, TypeList of TypeString): 删除时传入的副本名称列表；未配置时默认删除资源自身对应的副本
- `created_on` (Computed, TypeString): 边缘函数副本创建时间
- `modified_on` (Computed, TypeString): 边缘函数副本更新时间

The resource SHALL NOT introduce a nested `function_replicas` list wrapper; fields of the Describe list element SHALL be flattened to the top level of the schema.

#### Scenario: Schema defines all CRUD fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include zone_id, function_id, replica_name, content, remark, sort_by, sort_order, filters, replica_names, created_on, and modified_on with correct types and optionality

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** zone_id, function_id, or replica_name is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Computed time fields are read-only
- **WHEN** the DescribeFunctionReplicas API returns a replica with CreatedOn and ModifiedOn populated
- **THEN** the Read function SHALL set created_on and modified_on in the state without accepting user input for these fields

### Requirement: Create TEO function replica v5
The system SHALL create an edge function replica by calling `CreateFunctionReplica` API with ZoneId, FunctionId, ReplicaName, Content, and optionally Remark. Upon success, the resource ID SHALL be set to the composite `zone_id#function_id#replica_name` (joined by `tccommon.FILED_SP`), followed by a Read to populate the state.

#### Scenario: Successful creation with all parameters
- **WHEN** user provides zone_id, function_id, replica_name, content, and remark in the Terraform configuration
- **THEN** the system calls CreateFunctionReplica with ZoneId, FunctionId, ReplicaName, Content, and Remark, and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: Successful creation without optional remark
- **WHEN** user provides zone_id, function_id, replica_name, and content without remark
- **THEN** the system calls CreateFunctionReplica without Remark and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: API returns error during creation
- **WHEN** the CreateFunctionReplica API returns an error
- **THEN** the system SHALL retry via resource.Retry(tccommon.WriteRetryTimeout) with tccommon.RetryError wrapping and return the error if retries are exhausted

#### Scenario: Create response is empty
- **WHEN** the CreateFunctionReplica API succeeds but the response or its Response field is nil
- **THEN** the system SHALL return resource.NonRetryableError instead of setting the resource ID

### Requirement: Read TEO function replica v5
The system SHALL read the current state of an edge function replica by calling `DescribeFunctionReplicas` API with ZoneId and FunctionId parsed from the composite ID, the user-configured SortBy, SortOrder, and Filters passed through, an additional internal filter `{Name: "replica-name", Values: [replica_name]}` appended to locate the resource, and Limit set to 200 (the maximum value documented in the API). The system SHALL then match the replica whose ReplicaName equals the replica_name from the composite ID and set content, remark, created_on, and modified_on from the matched replica.

#### Scenario: Successful read of existing replica
- **WHEN** the resource exists and DescribeFunctionReplicas returns a matching replica
- **THEN** the system SHALL set content, remark, created_on, and modified_on from the matched FunctionReplica fields, each guarded by a nil check

#### Scenario: Replica not found during read
- **WHEN** the response is empty, the response Response field is nil, or no replica matches the replica_name from the composite ID
- **THEN** the system SHALL first log `[CRUD] teo_function_replica_v5 id=<id>` to preserve context, then remove the resource from state (d.SetId(""))

#### Scenario: Read uses maximum page size and internal filter
- **WHEN** reading the resource state
- **THEN** the system SHALL set Limit to 200 and append the internal replica-name filter to the request Filters

#### Scenario: User-configured query parameters are passed through
- **WHEN** the user configures sort_by, sort_order, or filters in the Terraform configuration
- **THEN** the Read function SHALL pass them to the DescribeFunctionReplicas request as SortBy, SortOrder, and Filters respectively

#### Scenario: Fuzzy user filter returning multiple replicas
- **WHEN** a user-configured fuzzy filter causes DescribeFunctionReplicas to return multiple replicas
- **THEN** the system SHALL select only the replica whose ReplicaName exactly equals the replica_name from the composite ID

### Requirement: Update TEO function replica v5
The system SHALL allow users to update the content and/or remark of an existing edge function replica by calling `ModifyFunctionReplica` API with ZoneId, FunctionId, ReplicaName parsed from the composite ID, and the new Content and/or Remark values. The update SHALL only be triggered when content or remark has changed.

#### Scenario: Update content only
- **WHEN** user changes the content field in the Terraform configuration
- **THEN** the system calls ModifyFunctionReplica with ZoneId, FunctionId, ReplicaName, and the new Content

#### Scenario: Update remark only
- **WHEN** user changes the remark field in the Terraform configuration
- **THEN** the system calls ModifyFunctionReplica with ZoneId, FunctionId, ReplicaName, and the new Remark

#### Scenario: Update both content and remark
- **WHEN** user changes both content and remark fields in the Terraform configuration
- **THEN** the system calls ModifyFunctionReplica with both new values in a single request

#### Scenario: No mutable field changed
- **WHEN** neither content nor remark is changed in the Terraform configuration
- **THEN** the system SHALL NOT call ModifyFunctionReplica and SHALL proceed directly to Read

#### Scenario: API returns error during update
- **WHEN** the ModifyFunctionReplica API returns an error
- **THEN** the system SHALL retry via resource.Retry(tccommon.WriteRetryTimeout) with tccommon.RetryError wrapping and return the error if retries are exhausted

### Requirement: Delete TEO function replica v5
The system SHALL delete an edge function replica by calling `DeleteFunctionReplica` API with ZoneId and FunctionId parsed from the composite ID, and ReplicaNames populated from the user-configured `replica_names` parameter when set, or a single-element list containing the replica_name from the composite ID when not set.

#### Scenario: Successful deletion with default replica name
- **WHEN** the resource is destroyed and replica_names is not configured
- **THEN** the system calls DeleteFunctionReplica with ReplicaNames containing the single replica_name from the composite ID

#### Scenario: Successful deletion with explicit replica_names
- **WHEN** the resource is destroyed and replica_names is configured with a list of names
- **THEN** the system calls DeleteFunctionReplica with ReplicaNames containing the configured list

#### Scenario: API returns error during deletion
- **WHEN** the DeleteFunctionReplica API returns an error
- **THEN** the system SHALL retry via resource.Retry(tccommon.WriteRetryTimeout) with tccommon.RetryError wrapping and return the error if retries are exhausted

### Requirement: Import TEO function replica v5
The system SHALL support importing an existing edge function replica using the composite ID format `zone_id#function_id#replica_name`, and the resource documentation SHALL describe this composite ID usage in the Import section.

#### Scenario: Successful import with composite ID
- **WHEN** user runs terraform import with ID in format `zone_id#function_id#replica_name`
- **THEN** the system SHALL parse the composite ID and call Read to populate the state

#### Scenario: Invalid import ID format
- **WHEN** user provides an import ID that does not contain exactly 3 parts separated by `#`
- **THEN** the Read/Update/Delete functions SHALL return an error indicating the ID is broken

### Requirement: Provider registration and documentation
The system SHALL register the `tencentcloud_teo_function_replica_v5` resource in `tencentcloud/provider.go` ResourcesMap and append `tencentcloud_teo_function_replica_v5` to the resource list in `tencentcloud/provider.md`. The system SHALL provide a resource documentation file `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.md` containing a one-sentence description mentioning the TEO product, an Example Usage section, and an Import section (without manually written Argument/Attribute Reference sections, which are generated by tooling).

#### Scenario: Resource is available after provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_function_replica_v5` SHALL be available for use in Terraform configurations

#### Scenario: Resource documentation exists
- **WHEN** the documentation file is generated
- **THEN** it SHALL contain a one-sentence description starting with "Provides a resource to create a TEO ...", an Example Usage HCL block, and an Import section explaining the `zone_id#function_id#replica_name` composite ID

### Requirement: Unit tests with mocked cloud APIs
The system SHALL provide unit tests in `tencentcloud/services/teo/resource_tc_teo_function_replica_v5_test.go` using gomonkey to mock the CreateFunctionReplica, DescribeFunctionReplicas, ModifyFunctionReplica, and DeleteFunctionReplica APIs, covering Create, Read (found and not found), Update (with and without changes), and Delete (default and explicit replica_names) scenarios without using the terraform test suite.

#### Scenario: Create test with mocked APIs
- **WHEN** the Create test runs with CreateFunctionReplica and DescribeFunctionReplicas mocked
- **THEN** the test SHALL verify the request parameters and that the resource ID equals `zone_id#function_id#replica_name`

#### Scenario: Read not found test
- **WHEN** the Read test mocks DescribeFunctionReplicas returning an empty replica list
- **THEN** the test SHALL verify the resource ID is cleared to an empty string

#### Scenario: Update test only fires Modify on change
- **WHEN** the Update test runs with content changed and ModifyFunctionReplica mocked
- **THEN** the test SHALL verify ModifyFunctionReplica is invoked with the expected parameters

#### Scenario: Delete test with replica_names
- **WHEN** the Delete test runs with replica_names configured and DeleteFunctionReplica mocked
- **THEN** the test SHALL verify ReplicaNames equals the configured list

