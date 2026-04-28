## ADDED Requirements

### Requirement: Tags parameter for cls_cloud_product_log_task_v2 resource
The `tencentcloud_cls_cloud_product_log_task_v2` resource SHALL accept a `tags` parameter of type list, containing objects with `key` (string, required) and `value` (string, required) fields. The maximum number of tag entries SHALL be 10. When provided, tags SHALL be passed to the `CreateCloudProductLogCollection` API request's `Tags` field.

#### Scenario: Create resource with tags
- **WHEN** a user creates a `tencentcloud_cls_cloud_product_log_task_v2` resource with `tags` specified
- **THEN** the tags SHALL be included in the `CreateCloudProductLogCollection` API request

#### Scenario: Create resource without tags
- **WHEN** a user creates a `tencentcloud_cls_cloud_product_log_task_v2` resource without specifying `tags`
- **THEN** the `Tags` field SHALL NOT be set in the `CreateCloudProductLogCollection` API request

#### Scenario: Tags are optional and not ForceNew
- **WHEN** a user updates the `tags` parameter on an existing resource
- **THEN** the resource SHALL NOT be destroyed and recreated, and the change SHALL be noted in the state

### Requirement: is_delete_topic parameter for cls_cloud_product_log_task_v2 resource
The `tencentcloud_cls_cloud_product_log_task_v2` resource SHALL accept an `is_delete_topic` parameter of type bool, optional with default false. When provided and set to true, it SHALL be passed to the `DeleteCloudProductLogCollection` API request's `IsDeleteTopic` field.

#### Scenario: Delete resource with is_delete_topic set to true
- **WHEN** a user deletes a `tencentcloud_cls_cloud_product_log_task_v2` resource with `is_delete_topic` set to true
- **THEN** the `IsDeleteTopic` field SHALL be set to true in the `DeleteCloudProductLogCollection` API request

#### Scenario: Delete resource with is_delete_topic not set
- **WHEN** a user deletes a `tencentcloud_cls_cloud_product_log_task_v2` resource without setting `is_delete_topic`
- **THEN** the `IsDeleteTopic` field SHALL NOT be set in the `DeleteCloudProductLogCollection` API request

### Requirement: is_delete_logset parameter for cls_cloud_product_log_task_v2 resource
The `tencentcloud_cls_cloud_product_log_task_v2` resource SHALL accept an `is_delete_logset` parameter of type bool, optional with default false. When provided and set to true, it SHALL be passed to the `DeleteCloudProductLogCollection` API request's `IsDeleteLogset` field.

#### Scenario: Delete resource with is_delete_logset set to true
- **WHEN** a user deletes a `tencentcloud_cls_cloud_product_log_task_v2` resource with `is_delete_logset` set to true
- **THEN** the `IsDeleteLogset` field SHALL be set to true in the `DeleteCloudProductLogCollection` API request

#### Scenario: Delete resource with is_delete_logset not set
- **WHEN** a user deletes a `tencentcloud_cls_cloud_product_log_task_v2` resource without setting `is_delete_logset`
- **THEN** the `IsDeleteLogset` field SHALL NOT be set in the `DeleteCloudProductLogCollection` API request

### Requirement: Status computed attribute for cls_cloud_product_log_task_v2 resource
The `tencentcloud_cls_cloud_product_log_task_v2` resource SHALL expose a `status` computed attribute of type int. The status SHALL be read from the `CreateCloudProductLogCollection` API response after creation, and from the `DescribeCloudProductLogTasks` API response during read operations. Valid values: -1 (creating), 0 (creating), 1 (created), 2 (deleting), 3 (deleted).

#### Scenario: Status is set after resource creation
- **WHEN** a `tencentcloud_cls_cloud_product_log_task_v2` resource is successfully created
- **THEN** the `status` attribute SHALL be set from the `CreateCloudProductLogCollection` response's `Status` field

#### Scenario: Status is read during refresh
- **WHEN** a `tencentcloud_cls_cloud_product_log_task_v2` resource is read/refreshed
- **THEN** the `status` attribute SHALL be set from the `DescribeCloudProductLogTasks` response's task `Status` field

#### Scenario: Status is not user-configurable
- **WHEN** a user attempts to set the `status` attribute in the Terraform configuration
- **THEN** the attribute SHALL be ignored as it is computed-only
