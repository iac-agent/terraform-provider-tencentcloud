## ADDED Requirements

### Requirement: Resource schema definition
The system SHALL define a Terraform resource `tencentcloud_cls_search_view_v2` with the following schema fields:

- `view_id` (Computed, string): 视图ID，由云API返回
- `logset_id` (Required, ForceNew, string): 日志集ID，仅创建时可设置
- `logset_region` (Required, ForceNew, string): 日志集所属地域，仅创建时可设置
- `view_name` (Required, string): 视图名称，最大255字符
- `view_type` (Required, string): 视图类型，枚举值：log、metric
- `topics` (Required, list of map): 查询视图中包含的主题，最大10个，每个主题包含 region、logset_id、topic_id
- `view_id_prefix` (Optional, ForceNew, string): 自定义视图ID前缀，仅创建时可设置
- `description` (Optional, string): 描述信息
- `create_time` (Computed, string): 创建时间
- `update_time` (Computed, string): 更新时间

#### Scenario: Schema validation for required fields
- **WHEN** a user creates a `tencentcloud_cls_search_view_v2` resource without specifying `logset_id`
- **THEN** Terraform SHALL return a validation error indicating the field is required

#### Scenario: ForceNew behavior on logset_id change
- **WHEN** a user changes `logset_id` on an existing resource
- **THEN** Terraform SHALL destroy and recreate the resource

#### Scenario: ForceNew behavior on logset_region change
- **WHEN** a user changes `logset_region` on an existing resource
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: Create operation
The system SHALL implement the Create operation by calling `CreateSearchView` API, mapping all create-time schema fields to the corresponding API request parameters, and storing the returned `ViewId` as the resource ID.

#### Scenario: Successful creation
- **WHEN** a user creates a `tencentcloud_cls_search_view_v2` resource with valid parameters
- **THEN** the system SHALL call `CreateSearchView` with the provided parameters
- **AND** store the returned `ViewId` as the Terraform resource ID

#### Scenario: Creation with all optional fields
- **WHEN** a user creates a resource with `view_id_prefix` and `description` specified
- **THEN** the system SHALL include these fields in the `CreateSearchView` request

### Requirement: Read operation
The system SHALL implement the Read operation by calling `DescribeSearchViews` API with a Filter on `viewId` equal to the resource's view_id, and map the returned `SearchViewInfo` fields back to the Terraform state.

#### Scenario: Resource exists
- **WHEN** the Read operation is called for an existing resource
- **THEN** the system SHALL call `DescribeSearchViews` with Filter key `viewId` and value matching the resource ID
- **AND** populate all schema fields from the first matching `SearchViewInfo` record

#### Scenario: Resource not found
- **WHEN** the Read operation is called and `DescribeSearchViews` returns no matching record
- **THEN** the system SHALL set the resource ID to empty string to signal resource deletion

### Requirement: Update operation
The system SHALL implement the Update operation by calling `ModifySearchView` API, mapping updatable schema fields (view_name, view_type, topics, description) to the API request parameters along with view_id.

#### Scenario: Successful update
- **WHEN** a user updates `view_name` or `view_type` or `topics` or `description`
- **THEN** the system SHALL call `ModifySearchView` with `ViewId` and the changed fields

#### Scenario: Fields that trigger ForceNew
- **WHEN** a user changes `logset_id`, `logset_region`, or `view_id_prefix`
- **THEN** Terraform SHALL NOT call ModifySearchView but instead destroy and recreate the resource

### Requirement: Delete operation
The system SHALL implement the Delete operation by calling `DeleteSearchView` API with `ViewId` parameter.

#### Scenario: Successful deletion
- **WHEN** a user destroys a `tencentcloud_cls_search_view_v2` resource
- **THEN** the system SHALL call `DeleteSearchView` with the resource's `view_id`

### Requirement: Provider registration
The system SHALL register `tencentcloud_cls_search_view_v2` resource in `provider.go` and document it in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** a user writes a Terraform configuration using `tencentcloud_cls_search_view_v2`
- **THEN** the Terraform provider SHALL recognize and manage this resource type

### Requirement: Unit tests
The system SHALL include unit tests for the resource using mock (gomonkey) approach, covering Create, Read, Update, and Delete operations.

#### Scenario: Unit test coverage
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** all Create, Read, Update, Delete functions SHALL be tested with mocked cloud API responses

### Requirement: Resource documentation
The system SHALL generate a `.md` documentation file for the resource following the project documentation conventions.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_cls_search_view_v2.md` file SHALL exist in the CLS service directory with description, example usage, and import section
