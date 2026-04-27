## ADDED Requirements

### Requirement: Resource schema includes effective parameter
The `tencentcloud_cls_ckafka_consumer` resource SHALL include an `effective` parameter of type TypeBool, marked as Optional and Computed. This parameter controls whether the delivery task is effective. It is settable via ModifyConsumer and readable from DescribeConsumer response.

#### Scenario: User sets effective to true
- **WHEN** user configures `effective = true` in the resource
- **THEN** the update function SHALL call ModifyConsumer with Effective=true after resource creation

#### Scenario: User does not set effective
- **WHEN** user does not configure `effective` in the resource
- **THEN** the read function SHALL still populate `effective` from the DescribeConsumer response

#### Scenario: Read populates effective from API
- **WHEN** the read function calls DescribeConsumer
- **THEN** if the response contains Effective, it SHALL be set in the Terraform state

### Requirement: Resource schema includes role_arn parameter
The `tencentcloud_cls_ckafka_consumer` resource SHALL include a `role_arn` parameter of type TypeString, marked as Optional. This parameter specifies the role access description name for cross-account access. It is settable via CreateConsumer and ModifyConsumer.

#### Scenario: User sets role_arn in create
- **WHEN** user configures `role_arn = "qcs::cam::uin/123456789:roleName/MyRole"` in the resource
- **THEN** the create function SHALL pass RoleArn to the CreateConsumer API request

#### Scenario: User updates role_arn
- **WHEN** user changes `role_arn` value in the resource configuration
- **THEN** the update function SHALL pass the new RoleArn to the ModifyConsumer API request

#### Scenario: User does not set role_arn
- **WHEN** user does not configure `role_arn` in the resource
- **THEN** the create and update functions SHALL NOT include RoleArn in the API request

### Requirement: Resource schema includes external_id parameter
The `tencentcloud_cls_ckafka_consumer` resource SHALL include an `external_id` parameter of type TypeString, marked as Optional. This parameter specifies the external ID for role assumption. It is settable via CreateConsumer and ModifyConsumer.

#### Scenario: User sets external_id in create
- **WHEN** user configures `external_id = "my-external-id"` in the resource
- **THEN** the create function SHALL pass ExternalId to the CreateConsumer API request

#### Scenario: User updates external_id
- **WHEN** user changes `external_id` value in the resource configuration
- **THEN** the update function SHALL pass the new ExternalId to the ModifyConsumer API request

#### Scenario: User does not set external_id
- **WHEN** user does not configure `external_id` in the resource
- **THEN** the create and update functions SHALL NOT include ExternalId in the API request

### Requirement: Resource schema includes advanced_config parameter
The `tencentcloud_cls_ckafka_consumer` resource SHALL include an `advanced_config` parameter of type TypeList with MaxItems 1, marked as Optional. It contains nested fields `partition_hash_status` (TypeBool, Optional) and `partition_fields` (TypeSet of TypeString, Optional). It is settable via CreateConsumer and ModifyConsumer.

#### Scenario: User configures advanced_config with all fields
- **WHEN** user configures `advanced_config { partition_hash_status = true partition_fields = ["__SOURCE__", "__HOSTNAME__"] }` in the resource
- **THEN** the create function SHALL pass AdvancedConfig with PartitionHashStatus=true and PartitionFields=["__SOURCE__", "__HOSTNAME__"] to the CreateConsumer API request

#### Scenario: User configures advanced_config with only partition_hash_status
- **WHEN** user configures `advanced_config { partition_hash_status = false }` in the resource
- **THEN** the create function SHALL pass AdvancedConfig with PartitionHashStatus=false and no PartitionFields to the CreateConsumer API request

#### Scenario: User updates advanced_config
- **WHEN** user changes `advanced_config` in the resource configuration
- **THEN** the update function SHALL pass the updated AdvancedConfig to the ModifyConsumer API request

#### Scenario: User does not set advanced_config
- **WHEN** user does not configure `advanced_config` in the resource
- **THEN** the create and update functions SHALL NOT include AdvancedConfig in the API request

### Requirement: Update function handles new mutable parameters
The update function SHALL include `effective`, `role_arn`, `external_id`, and `advanced_config` in the mutableArgs list. When any of these parameters change, the update function SHALL call ModifyConsumer with all current parameter values.

#### Scenario: Only effective changes
- **WHEN** only `effective` parameter changes in the resource configuration
- **THEN** the update function SHALL detect the change and call ModifyConsumer with TopicId and Effective, along with other current parameter values

#### Scenario: Multiple new parameters change together
- **WHEN** both `role_arn` and `advanced_config` change in the resource configuration
- **THEN** the update function SHALL call ModifyConsumer once with all updated parameters

### Requirement: Unit tests cover new parameter logic
Unit tests SHALL be added using gomonkey mock approach to verify the new parameter handling logic in create, read, and update functions.

#### Scenario: Create with all new parameters
- **WHEN** create function is called with effective, role_arn, external_id, and advanced_config
- **THEN** the function SHALL correctly map all parameters to the CreateConsumer API request

#### Scenario: Read populates effective from API response
- **WHEN** read function is called and DescribeConsumer returns Effective=true
- **THEN** the function SHALL set `effective` to true in the Terraform state
