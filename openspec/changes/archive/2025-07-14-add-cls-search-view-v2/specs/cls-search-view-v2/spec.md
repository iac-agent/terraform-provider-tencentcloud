## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_cls_search_view_v2` with the following schema fields:
- `logset_id` (TypeString, Required, ForceNew): 日志集 ID
- `logset_region` (TypeString, Required, ForceNew): 日志集所属地域
- `view_name` (TypeString, Required): 视图名称，最大 255 字符，不能包含 "|" 字符
- `view_type` (TypeString, Required): 视图类型，枚举值：log（日志主题）、metric（指标主题）
- `topics` (TypeList, Required): 查询视图中包含的主题列表，最大 10 个主题，嵌套块包含 region/logset_id/topic_id
- `description` (TypeString, Optional): 描述信息
- `view_id_prefix` (TypeString, Optional, ForceNew): 自定义查询视图 ID 前缀
- `view_id` (TypeString, Computed): 视图 ID，由云 API 生成

#### Scenario: Schema fields match cloud API parameters
- **WHEN** the resource schema is defined
- **THEN** all CreateSearchView request parameters are mapped to schema fields with correct types and Required/Optional attributes

#### Scenario: Immutable fields are ForceNew
- **WHEN** logset_id, logset_region, or view_id_prefix is changed after creation
- **THEN** Terraform SHALL force resource recreation

### Requirement: Resource Create
The system SHALL implement resource creation by calling CreateSearchView API, mapping all schema fields to request parameters, and setting the resource ID from the response ViewId.

#### Scenario: Successful creation
- **WHEN** `tencentcloud_cls_search_view_v2` resource is created with valid parameters
- **THEN** the system SHALL call CreateSearchView API with logset_id, logset_region, view_name, view_type, topics, description, and view_id_prefix
- **AND** set the resource ID to the returned ViewId
- **AND** set view_id in state to the returned ViewId

#### Scenario: Creation with retry on error
- **WHEN** CreateSearchView API call fails
- **THEN** the system SHALL retry using tccommon.ReadRetryTimeout with tccommon.RetryError wrapping

### Requirement: Resource Read
The system SHALL implement resource reading by calling DescribeSearchViews API with a Filter on viewId to fetch current resource state.

#### Scenario: Resource exists
- **WHEN** reading an existing resource
- **THEN** the system SHALL call DescribeSearchViews with Filter key "viewId" and value of the current view_id
- **AND** populate all schema fields from the matching SearchViewInfo entry

#### Scenario: Resource not found
- **WHEN** DescribeSearchViews returns empty Infos for the given viewId
- **THEN** the system SHALL set d.SetId("") to indicate the resource no longer exists

#### Scenario: Read with retry on error
- **WHEN** DescribeSearchViews API call fails
- **THEN** the system SHALL retry using tccommon.ReadRetryTimeout with tccommon.RetryError wrapping

### Requirement: Resource Update
The system SHALL implement resource update by calling ModifySearchView API, mapping updatable schema fields (view_name, view_type, topics, description) to request parameters.

#### Scenario: Successful update
- **WHEN** updatable fields (view_name, view_type, topics, description) are changed
- **THEN** the system SHALL call ModifySearchView API with view_id, view_name, view_type, topics, and description
- **AND** call Read to refresh state after update

#### Scenario: Update with retry on error
- **WHEN** ModifySearchView API call fails
- **THEN** the system SHALL retry using tccommon.ReadRetryTimeout with tccommon.RetryError wrapping

### Requirement: Resource Delete
The system SHALL implement resource deletion by calling DeleteSearchView API with the ViewId.

#### Scenario: Successful deletion
- **WHEN** `tencentcloud_cls_search_view_v2` resource is destroyed
- **THEN** the system SHALL call DeleteSearchView API with view_id

#### Scenario: Delete with retry on error
- **WHEN** DeleteSearchView API call fails
- **THEN** the system SHALL retry using tccommon.ReadRetryTimeout with tccommon.RetryError wrapping

### Requirement: Resource Import
The system SHALL support importing existing CLS SearchView resources by ViewId.

#### Scenario: Import by ViewId
- **WHEN** `terraform import` is called for this resource
- **THEN** the system SHALL accept ViewId as the import ID and populate all fields via Read

### Requirement: Provider Registration
The system SHALL register `tencentcloud_cls_search_view_v2` resource in provider.go and provider.md.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cls_search_view_v2` SHALL be available as a resource type

### Requirement: Unit Tests
The system SHALL include unit tests using gomonkey mock approach for all CRUD operations.

#### Scenario: CRUD operations tested with mock
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** Create, Read, Update, Delete operations SHALL be tested using gomonkey mocks for cloud API calls

### Requirement: Documentation
The system SHALL include a .md documentation file for the resource.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a resource_tc_cls_search_view_v2.md file SHALL exist with description, example usage, and import sections
