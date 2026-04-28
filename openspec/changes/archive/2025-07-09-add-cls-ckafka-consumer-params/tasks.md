## 1. Schema Definition

- [x] 1.1 Add `effective` (bool, Optional) parameter to the resource schema in `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go`
- [x] 1.2 Add `role_arn` (string, Optional) parameter to the resource schema
- [x] 1.3 Add `external_id` (string, Optional) parameter to the resource schema
- [x] 1.4 Add `advanced_config` (TypeList, MaxItems:1, Optional) parameter to the resource schema with sub-fields: `partition_hash_status` (bool, Optional) and `partition_fields` (TypeSet of string, Optional)

## 2. Create Function Update

- [x] 2.1 Add `role_arn` parameter handling in `resourceTencentCloudClsCkafkaConsumerCreate`: set `request.RoleArn` from schema
- [x] 2.2 Add `external_id` parameter handling in `resourceTencentCloudClsCkafkaConsumerCreate`: set `request.ExternalId` from schema
- [x] 2.3 Add `advanced_config` parameter handling in `resourceTencentCloudClsCkafkaConsumerCreate`: map schema fields to `cls.AdvancedConsumerConfiguration` struct and set `request.AdvancedConfig`

## 3. Read Function Update

- [x] 3.1 Add `effective` field reading in `resourceTencentCloudClsCkafkaConsumerRead`: set `effective` from DescribeConsumer response's `Effective` field

## 4. Update Function Update

- [x] 4.1 Add `effective`, `role_arn`, `external_id`, `advanced_config` to the `mutableArgs` list in `resourceTencentCloudClsCkafkaConsumerUpdate`
- [x] 4.2 Add `effective` parameter handling in the update function: set `request.Effective` from schema when changed
- [x] 4.3 Add `role_arn` parameter handling in the update function: set `request.RoleArn` from schema when changed
- [x] 4.4 Add `external_id` parameter handling in the update function: set `request.ExternalId` from schema when changed
- [x] 4.5 Add `advanced_config` parameter handling in the update function: map schema fields to `cls.AdvancedConsumerConfiguration` struct and set `request.AdvancedConfig` when changed

## 5. Documentation Update

- [x] 5.1 Update `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.md` to include new parameters in the example usage

## 6. Unit Tests

- [x] 6.1 Add unit tests for `effective` parameter in create/read/update flows using gomonkey mock
- [x] 6.2 Add unit tests for `role_arn` parameter in create/update flows using gomonkey mock
- [x] 6.3 Add unit tests for `external_id` parameter in create/update flows using gomonkey mock
- [x] 6.4 Add unit tests for `advanced_config` parameter in create/update flows using gomonkey mock
- [x] 6.5 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
