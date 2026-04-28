## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go` 的 Schema 中新增 `effective` 参数（TypeBool, Optional, Computed, Description: "投递任务是否生效。"）
- [x] 1.2 在 Schema 中新增 `role_arn` 参数（TypeString, Optional, Description: "角色访问描述名。"）
- [x] 1.3 在 Schema 中新增 `external_id` 参数（TypeString, Optional, Description: "外部ID。"）
- [x] 1.4 在 Schema 中新增 `advanced_config` 参数（TypeList, MaxItems:1, Optional, Description: "高级配置项。"），包含子字段 `partition_hash_status`（TypeBool, Optional, Description: "Ckafka分区hash状态，默认false。"）和 `partition_fields`（TypeList of TypeString, Optional, Description: "需要计算hash的字段列表，最大支持5个字段。"）

## 2. Create 函数

- [x] 2.1 在 `resourceTencentCloudClsCkafkaConsumerCreate` 函数中，添加 `role_arn` 参数的处理逻辑：使用 `d.GetOk("role_arn")` 获取值并设置到 `request.RoleArn`
- [x] 2.2 在 Create 函数中，添加 `external_id` 参数的处理逻辑：使用 `d.GetOk("external_id")` 获取值并设置到 `request.ExternalId`
- [x] 2.3 在 Create 函数中，添加 `advanced_config` 参数的处理逻辑：使用 `helper.InterfacesHeadMap(d, "advanced_config")` 获取嵌套结构，构建 `cls.AdvancedConsumerConfiguration`，设置 `PartitionHashStatus` 和 `PartitionFields`，然后赋值到 `request.AdvancedConfig`

## 3. Read 函数

- [x] 3.1 在 `resourceTencentCloudClsCkafkaConsumerRead` 函数中，添加 `effective` 参数的读取逻辑：检查 `ckafkaConsumer.Effective != nil`，若非 nil 则 `d.Set("effective", ckafkaConsumer.Effective)`

## 4. Update 函数

- [x] 4.1 在 `resourceTencentCloudClsCkafkaConsumerUpdate` 函数的 mutableArgs 列表中，添加 `"effective"`、`"role_arn"`、`"external_id"`、`"advanced_config"` 四个参数
- [x] 4.2 在 Update 函数的 needChange 代码块中，添加 `effective` 参数的处理逻辑：使用 `d.GetOkExists("effective")` 获取值并设置到 `request.Effective`
- [x] 4.3 在 Update 函数中，添加 `role_arn` 参数的处理逻辑：使用 `d.GetOk("role_arn")` 获取值并设置到 `request.RoleArn`
- [x] 4.4 在 Update 函数中，添加 `external_id` 参数的处理逻辑：使用 `d.GetOk("external_id")` 获取值并设置到 `request.ExternalId`
- [x] 4.5 在 Update 函数中，添加 `advanced_config` 参数的处理逻辑：使用 `helper.InterfacesHeadMap(d, "advanced_config")` 获取嵌套结构，构建 `cls.AdvancedConsumerConfiguration`，设置 `PartitionHashStatus` 和 `PartitionFields`，然后赋值到 `request.AdvancedConfig`

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer_test.go` 中，补充新增参数的单元测试用例，覆盖 Create/Read/Update 中新增参数的处理逻辑，使用 gomonkey mock 云 API

## 6. 文档

- [x] 6.1 更新 `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.md` 文档，添加新增参数的 Example Usage
