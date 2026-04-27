## 1. Schema Definition

- [x] 1.1 Add `effective` (TypeBool, Optional, Computed) parameter to the `tencentcloud_cls_ckafka_consumer` resource schema in `resource_tc_cls_ckafka_consumer.go`
- [x] 1.2 Add `role_arn` (TypeString, Optional) parameter to the resource schema
- [x] 1.3 Add `external_id` (TypeString, Optional) parameter to the resource schema
- [x] 1.4 Add `advanced_config` (TypeList, MaxItems 1, Optional) parameter with nested fields `partition_hash_status` (TypeBool, Optional) and `partition_fields` (TypeSet of TypeString, Optional) to the resource schema

## 2. Create Function Update

- [x] 2.1 Update `resourceTencentCloudClsCkafkaConsumerCreate` to pass `role_arn` to CreateConsumer request (RoleArn field)
- [x] 2.2 Update `resourceTencentCloudClsCkafkaConsumerCreate` to pass `external_id` to CreateConsumer request (ExternalId field)
- [x] 2.3 Update `resourceTencentCloudClsCkafkaConsumerCreate` to pass `advanced_config` (with PartitionHashStatus and PartitionFields) to CreateConsumer request (AdvancedConfig field)

## 3. Read Function Update

- [x] 3.1 Update `resourceTencentCloudClsCkafkaConsumerRead` to set `effective` from DescribeConsumer response (Effective field)

## 4. Update Function Update

- [x] 4.1 Add `effective`, `role_arn`, `external_id`, `advanced_config` to the `mutableArgs` list in `resourceTencentCloudClsCkafkaConsumerUpdate`
- [x] 4.2 Update `resourceTencentCloudClsCkafkaConsumerUpdate` to pass `effective` to ModifyConsumer request (Effective field)
- [x] 4.3 Update `resourceTencentCloudClsCkafkaConsumerUpdate` to pass `role_arn` to ModifyConsumer request (RoleArn field)
- [x] 4.4 Update `resourceTencentCloudClsCkafkaConsumerUpdate` to pass `external_id` to ModifyConsumer request (ExternalId field)
- [x] 4.5 Update `resourceTencentCloudClsCkafkaConsumerUpdate` to pass `advanced_config` (with PartitionHashStatus and PartitionFields) to ModifyConsumer request (AdvancedConfig field)

## 5. Unit Tests

- [x] 5.1 Add unit tests in `resource_tc_cls_ckafka_consumer_test.go` using gomonkey mock approach to verify create function correctly maps new parameters (role_arn, external_id, advanced_config) to CreateConsumer API request
- [x] 5.2 Add unit tests to verify read function correctly populates `effective` from DescribeConsumer API response
- [x] 5.3 Add unit tests to verify update function correctly maps new parameters (effective, role_arn, external_id, advanced_config) to ModifyConsumer API request
- [x] 5.4 Run unit tests with `go test -gcflags=all=-l` to ensure all tests pass

## 6. Documentation

- [x] 6.1 Update `resource_tc_cls_ckafka_consumer.md` to include examples and documentation for the new parameters (effective, role_arn, external_id, advanced_config)
