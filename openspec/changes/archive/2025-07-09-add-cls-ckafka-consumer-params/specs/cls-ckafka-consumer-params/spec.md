## ADDED Requirements

### Requirement: Resource supports effective parameter
The resource SHALL support an `effective` (bool, Optional) parameter that controls whether the CKafka consumer delivery task is effective. This parameter SHALL be sent in the ModifyConsumer API request when changed and SHALL be read from the DescribeConsumer API response.

#### Scenario: Set effective to true
- **WHEN** user sets `effective = true` in the resource configuration
- **THEN** the update function SHALL include `Effective = true` in the ModifyConsumer request

#### Scenario: Read effective from API
- **WHEN** the read function calls DescribeConsumer and the response contains `Effective`
- **THEN** the `effective` field SHALL be set in the Terraform state from the API response

#### Scenario: Effective defaults to nil
- **WHEN** user does not set `effective` in the resource configuration
- **THEN** the create function SHALL NOT send the `Effective` field in the CreateConsumer request

### Requirement: Resource supports role_arn parameter
The resource SHALL support a `role_arn` (string, Optional) parameter that specifies the role ARN for cross-account CKafka access. This parameter SHALL be sent in CreateConsumer and ModifyConsumer API requests. Since DescribeConsumer does not return this field, it is a write-only parameter managed through Terraform state.

#### Scenario: Set role_arn during creation
- **WHEN** user sets `role_arn` in the resource configuration
- **THEN** the create function SHALL include `RoleArn` in the CreateConsumer request

#### Scenario: Update role_arn
- **WHEN** user changes `role_arn` in the resource configuration
- **THEN** the update function SHALL include `RoleArn` in the ModifyConsumer request

#### Scenario: role_arn not set
- **WHEN** user does not set `role_arn` in the resource configuration
- **THEN** neither CreateConsumer nor ModifyConsumer SHALL include the `RoleArn` field

### Requirement: Resource supports external_id parameter
The resource SHALL support an `external_id` (string, Optional) parameter that specifies the external ID for role assumption. This parameter SHALL be sent in CreateConsumer and ModifyConsumer API requests. Since DescribeConsumer does not return this field, it is a write-only parameter managed through Terraform state.

#### Scenario: Set external_id during creation
- **WHEN** user sets `external_id` in the resource configuration
- **THEN** the create function SHALL include `ExternalId` in the CreateConsumer request

#### Scenario: Update external_id
- **WHEN** user changes `external_id` in the resource configuration
- **THEN** the update function SHALL include `ExternalId` in the ModifyConsumer request

#### Scenario: external_id not set
- **WHEN** user does not set `external_id` in the resource configuration
- **THEN** neither CreateConsumer nor ModifyConsumer SHALL include the `ExternalId` field

### Requirement: Resource supports advanced_config parameter
The resource SHALL support an `advanced_config` (TypeList, MaxItems:1, Optional) parameter that specifies advanced consumer configuration. The parameter SHALL contain sub-fields: `partition_hash_status` (bool, Optional) and `partition_fields` (TypeSet of string, Optional). This parameter SHALL be mapped to the `AdvancedConsumerConfiguration` cloud API struct and sent in CreateConsumer and ModifyConsumer API requests. Since DescribeConsumer does not return this field, it is a write-only parameter managed through Terraform state.

#### Scenario: Set advanced_config during creation
- **WHEN** user configures `advanced_config` with `partition_hash_status = true` and `partition_fields = ["field1", "field2"]`
- **THEN** the create function SHALL map these values to an `AdvancedConsumerConfiguration` struct and include it in the CreateConsumer request

#### Scenario: Update advanced_config
- **WHEN** user changes `advanced_config` in the resource configuration
- **THEN** the update function SHALL include the updated `AdvancedConfig` in the ModifyConsumer request

#### Scenario: advanced_config not set
- **WHEN** user does not set `advanced_config` in the resource configuration
- **THEN** neither CreateConsumer nor ModifyConsumer SHALL include the `AdvancedConfig` field

#### Scenario: partition_fields with max 5 fields
- **WHEN** user sets `partition_fields` with up to 5 field names
- **THEN** the fields SHALL be mapped to the `PartitionFields` array in the `AdvancedConsumerConfiguration` struct
