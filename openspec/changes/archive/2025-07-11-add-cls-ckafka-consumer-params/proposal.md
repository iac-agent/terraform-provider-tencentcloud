## Why

The `tencentcloud_cls_ckafka_consumer` resource is missing several parameters that are already supported by the CLS cloud API (CreateConsumer, ModifyConsumer, DescribeConsumer). Users cannot configure the `effective` (delivery task enablement), `role_arn` (cross-account role), `external_id` (external ID for role), and `advanced_config` (advanced Kafka delivery settings including partition hash) parameters through Terraform, which limits their ability to fully manage CLS CKafka consumer delivery tasks.

## What Changes

- Add `effective` (TypeBool, Optional) parameter to the resource schema, settable in update (ModifyConsumer) and readable from DescribeConsumer response
- Add `role_arn` (TypeString, Optional) parameter to the resource schema, settable in create (CreateConsumer) and update (ModifyConsumer)
- Add `external_id` (TypeString, Optional) parameter to the resource schema, settable in create (CreateConsumer) and update (ModifyConsumer)
- Add `advanced_config` (TypeList, MaxItems 1, Optional) parameter to the resource schema with nested fields `partition_hash_status` (TypeBool) and `partition_fields` (TypeSet of strings), settable in create (CreateConsumer) and update (ModifyConsumer)
- Update CRUD functions (create, read, update) to handle the new parameters
- Update unit tests to cover new parameter logic
- Update resource documentation (.md file)

## Capabilities

### New Capabilities
- `cls-ckafka-consumer-advanced-params`: Adds effective, role_arn, external_id, and advanced_config parameters to the tencentcloud_cls_ckafka_consumer resource, enabling full configuration of CLS CKafka consumer delivery tasks

### Modified Capabilities


## Impact

- `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go` - Schema and CRUD function modifications
- `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer_test.go` - Unit test additions
- `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.md` - Documentation updates
- Cloud API dependencies: cls/v20201016 (CreateConsumer, ModifyConsumer, DescribeConsumer) - already available in vendor
