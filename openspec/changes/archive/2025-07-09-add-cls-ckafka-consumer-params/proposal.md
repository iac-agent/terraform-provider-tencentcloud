## Why

The `tencentcloud_cls_ckafka_consumer` resource is missing several parameters that the cloud API already supports. Specifically, the `effective`, `role_arn`, `external_id`, and `advanced_config` parameters are available in the CreateConsumer and ModifyConsumer API interfaces but are not yet exposed in the Terraform resource schema. Adding these parameters enables users to fully configure CLS CKafka consumer delivery tasks, including controlling task effectiveness, cross-account access via role ARN, and advanced delivery options like partition hashing.

## What Changes

- Add `effective` (bool, Optional) parameter to the resource schema: controls whether the consumer delivery task is effective, set via ModifyConsumer and returned by DescribeConsumer
- Add `role_arn` (string, Optional) parameter to the resource schema: specifies the role ARN for cross-account CKafka access, passed to CreateConsumer and ModifyConsumer
- Add `external_id` (string, Optional) parameter to the resource schema: specifies the external ID for role assumption, passed to CreateConsumer and ModifyConsumer
- Add `advanced_config` (TypeList, MaxItems:1, Optional) parameter to the resource schema: specifies advanced consumer configuration including partition hash status and partition fields, passed to CreateConsumer and ModifyConsumer
- Update create, read, and update resource functions to handle new parameters
- Update the mutable args list in update function to include new parameters
- Update unit tests to cover new parameter logic
- Update the resource .md documentation file

## Capabilities

### New Capabilities
- `cls-ckafka-consumer-params`: Adds effective, role_arn, external_id, and advanced_config parameters to the tencentcloud_cls_ckafka_consumer resource

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go` - schema, create, read, update functions
- **Tests**: `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer_test.go` - unit tests for new parameters
- **Docs**: `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.md` - documentation update
- **APIs**: CreateConsumer, ModifyConsumer, DescribeConsumer (cls v20201016)
- **Backward Compatibility**: All new parameters are Optional, so existing configurations remain valid
