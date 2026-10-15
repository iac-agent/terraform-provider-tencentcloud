## ADDED Requirements

### Requirement: tencentcloud_config_update_aggregate_config_deliver resource registration
The provider SHALL register a new RESOURCE_KIND_CONFIG resource named
`tencentcloud_config_update_aggregate_config_deliver` for managing the Tencent Cloud Config
account-group delivery setting (投递设置).

#### Scenario: Resource is available in the provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_config_update_aggregate_config_deliver` SHALL be present in the provider
  resources map and listed in `provider.md`

### Requirement: Account group identifier is required and immutable
The resource SHALL require an `account_group_id` argument that identifies the target account group,
and this argument SHALL be immutable (force new resource).

#### Scenario: account_group_id provided
- **WHEN** a configuration specifies `account_group_id`
- **THEN** the resource SHALL use it as both the Read request input (`DescribeAggregateConfigDeliver`)
  and the Update request input (`UpdateAggregateConfigDeliver`)

#### Scenario: account_group_id changed
- **WHEN** `account_group_id` is modified in the configuration
- **THEN** Terraform SHALL plan a resource replacement

### Requirement: Read aggregate config deliver setting
The resource SHALL read the account-group delivery setting via `DescribeAggregateConfigDeliver` and
populate `deliver_name`, `target_arn`, `status`, `create_time`, `deliver_prefix`, `deliver_type`,
`deliver_uin` and `deliver_content_type`.

#### Scenario: Successful read
- **WHEN** the Read handler runs and the API returns a non-nil response
- **THEN** the resource SHALL set each non-nil response field into Terraform state and keep the id

#### Scenario: Setting no longer exists
- **WHEN** the Read handler runs and the API returns a nil response
- **THEN** the resource SHALL log the current id and set `d.SetId("")` so Terraform plans a recreate

### Requirement: Update aggregate config deliver setting
The resource SHALL update the account-group delivery setting via `UpdateAggregateConfigDeliver`,
supporting `status`, `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_uin`
and `deliver_content_type`.

#### Scenario: Successful update
- **WHEN** a mutable field changes
- **THEN** the Update handler SHALL call `UpdateAggregateConfigDeliver` with the changed fields and
  `status`/`account_group_id`, using `tccommon.WriteRetryTimeout`, then call Read to refresh state

#### Scenario: Cloud API failure
- **WHEN** `UpdateAggregateConfigDeliver` returns an error
- **THEN** the error SHALL be wrapped with `tccommon.RetryError` and returned after retries are
  exhausted

### Requirement: Create performs the first update
Because no dedicated create API exists, the resource SHALL implement Create by setting the resource
id to `account_group_id` and then invoking the Update handler.

#### Scenario: First apply
- **WHEN** the resource is created
- **THEN** `d.SetId(account_group_id)` SHALL be called before the update call, and the setting SHALL
  be created/updated in the cloud

### Requirement: Delete is a no-op
The resource SHALL implement Delete as a no-op because no delete API exists for the account-group
delivery setting.

#### Scenario: Destroy
- **WHEN** the resource is destroyed
- **THEN** the Delete handler SHALL return nil without calling any cloud API

### Requirement: Documentation and example
The resource SHALL provide an example documentation file
`resource_tc_config_update_aggregate_config_deliver.md` describing usage.

#### Scenario: Docs present
- **WHEN** `make doc` runs
- **THEN** documentation SHALL be generated for `tencentcloud_config_update_aggregate_config_deliver`
