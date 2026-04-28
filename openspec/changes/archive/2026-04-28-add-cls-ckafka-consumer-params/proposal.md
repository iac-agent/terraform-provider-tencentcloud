## Why

当前 `tencentcloud_cls_ckafka_consumer` 资源缺少云API已支持的高级配置参数（`effective`、`role_arn`、`external_id`、`advanced_config`），导致用户无法通过 Terraform 管理消费投递任务是否生效、角色访问描述、外部ID以及CKafka分区Hash等高级配置。

## What Changes

- 为 `tencentcloud_cls_ckafka_consumer` 资源新增 `effective` 参数（TypeBool），用于控制投递任务是否生效，对应 ModifyConsumer 的入参和 DescribeConsumer 的出参
- 为 `tencentcloud_cls_ckafka_consumer` 资源新增 `role_arn` 参数（TypeString），用于设置角色访问描述名，对应 CreateConsumer 和 ModifyConsumer 的入参
- 为 `tencentcloud_cls_ckafka_consumer` 资源新增 `external_id` 参数（TypeString），用于设置外部ID，对应 CreateConsumer 和 ModifyConsumer 的入参
- 为 `tencentcloud_cls_ckafka_consumer` 资源新增 `advanced_config` 参数（TypeList），用于设置高级配置项，对应 CreateConsumer 和 ModifyConsumer 的入参
  - 子字段 `partition_hash_status`（TypeBool）：CKafka分区hash状态
  - 子字段 `partition_fields`（TypeList of TypeString）：需要计算hash的字段列表
- 在 Create 函数中添加 `role_arn`、`external_id`、`advanced_config` 参数的处理逻辑
- 在 Update 函数中添加 `effective`、`role_arn`、`external_id`、`advanced_config` 参数的处理逻辑
- 在 Read 函数中添加 `effective` 参数的读取逻辑（DescribeConsumer 仅返回 effective，不返回 role_arn/external_id/advanced_config）

## Capabilities

### New Capabilities

- `cls-ckafka-consumer-params`: 为 cls_ckafka_consumer 资源新增 effective、role_arn、external_id、advanced_config 四个参数

### Modified Capabilities

## Impact

- 修改文件：`tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go`（Schema 和 CRUD 函数）
- 修改文件：`tencentcloud/services/cls/resource_tc_cls_ckafka_consumer_test.go`（单元测试）
- 修改文件：`tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.md`（文档）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016` 中 `AdvancedConsumerConfiguration` 结构体
- 云API接口：CreateConsumer、ModifyConsumer、DescribeConsumer
