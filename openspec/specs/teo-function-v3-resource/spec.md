## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_function_v3` resource SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): 站点 ID
- `name` (TypeString, Required): 函数名称，只能包含小写字母、数字、连字符，以数字或字母开头，以数字或字母结尾，最大支持 30 个字符
- `content` (TypeString, Required): 函数内容，当前仅支持 JavaScript 代码，最大支持 5MB 大小
- `remark` (TypeString, Optional): 函数描述，最大支持 60 个字符
- `function_id` (TypeString, Computed): 函数 ID
- `domain` (TypeString, Computed): 函数默认域名
- `create_time` (TypeString, Computed): 创建时间
- `update_time` (TypeString, Computed): 修改时间

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is initialized
- **THEN** all required fields (zone_id, name, content) are present with correct types and Required flag
- **AND** all optional fields (remark) are present with Optional flag
- **AND** all computed fields (function_id, domain, create_time, update_time) are present with Computed flag
- **AND** zone_id has ForceNew set to true

### Requirement: Resource creation with CreateFunction API
The resource SHALL create an edge function by calling the CreateFunction API with ZoneId, Name, Content, and Remark parameters. After creation, it SHALL poll the DescribeFunctions API until the function's Domain field is non-empty, then set the resource ID to `zone_id + FILED_SP + function_id`.

#### Scenario: Successful creation with all required parameters
- **WHEN** resource is created with zone_id, name, and content
- **THEN** CreateFunction API is called with ZoneId, Name, and Content
- **AND** the resource ID is set to `zone_id#function_id` format

#### Scenario: Successful creation with optional remark
- **WHEN** resource is created with zone_id, name, content, and remark
- **THEN** CreateFunction API is called with ZoneId, Name, Content, and Remark
- **AND** the resource ID is set to `zone_id#function_id` format

#### Scenario: CreateFunction returns empty FunctionId
- **WHEN** CreateFunction API returns a response with empty FunctionId
- **THEN** a NonRetryableError SHALL be returned

#### Scenario: CreateFunction is asynchronous
- **WHEN** CreateFunction API returns successfully
- **THEN** the system SHALL poll DescribeFunctions API until the Domain field is non-empty
- **AND** polling timeout is 600 seconds

### Requirement: Resource read with DescribeFunctions API
The resource SHALL read the current state by calling DescribeFunctions API with ZoneId and FunctionIds. It SHALL parse the composite ID to extract zone_id and function_id. If the function is not found, it SHALL log the id and set d.SetId("").

#### Scenario: Successful read
- **WHEN** resource read is called with a valid composite ID
- **THEN** DescribeFunctions API is called with ZoneId and FunctionIds
- **AND** all computed fields are set from the API response

#### Scenario: Resource not found
- **WHEN** DescribeFunctions API returns an empty function list
- **THEN** the system SHALL log the id with `log.Printf("[CRUD] teo_function_v3 id=%s", d.Id())`
- **AND** set d.SetId("") to remove the resource from state

#### Scenario: Nil check before setting fields
- **WHEN** the API response contains nil fields
- **THEN** the system SHALL skip calling d.Set() for those nil fields

### Requirement: Resource update with ModifyFunction API
The resource SHALL update the edge function by calling ModifyFunction API. The `name` field SHALL be immutable - if changed, the system SHALL return an error. The `content` and `remark` fields SHALL be mutable.

#### Scenario: Update content field
- **WHEN** content field is changed
- **THEN** ModifyFunction API is called with ZoneId, FunctionId, and Content

#### Scenario: Update remark field
- **WHEN** remark field is changed
- **THEN** ModifyFunction API is called with ZoneId, FunctionId, and Remark

#### Scenario: Attempt to change immutable name field
- **WHEN** name field is changed
- **THEN** an error SHALL be returned with message indicating the argument cannot be changed

### Requirement: Resource deletion with DeleteFunction API
The resource SHALL delete the edge function by calling DeleteFunction API with ZoneId and FunctionId.

#### Scenario: Successful deletion
- **WHEN** resource delete is called with a valid composite ID
- **THEN** DeleteFunction API is called with ZoneId and FunctionId

### Requirement: Resource import support
The resource SHALL support import using the composite ID format `zone_id#function_id`.

#### Scenario: Successful import
- **WHEN** terraform import is called with `zone_id#function_id`
- **THEN** the resource is imported with all fields populated from the DescribeFunctions API response

### Requirement: Resource registration in provider
The resource SHALL be registered in `provider.go` with the name `tencentcloud_teo_function_v3` and documented in `provider.md`.

#### Scenario: Resource is registered
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_function_v3` resource is available for use

### Requirement: Unit tests with gomonkey
The resource SHALL have unit tests using gomonkey for mocking cloud API calls, not using terraform test suite.

#### Scenario: Test create function
- **WHEN** the create function is tested
- **THEN** gomonkey is used to mock the CreateFunction API call

#### Scenario: Test read function
- **WHEN** the read function is tested
- **THEN** gomonkey is used to mock the DescribeFunctions API call

#### Scenario: Test update function
- **WHEN** the update function is tested
- **THEN** gomonkey is used to mock the ModifyFunction API call

#### Scenario: Test delete function
- **WHEN** the delete function is tested
- **THEN** gomonkey is used to mock the DeleteFunction API call
