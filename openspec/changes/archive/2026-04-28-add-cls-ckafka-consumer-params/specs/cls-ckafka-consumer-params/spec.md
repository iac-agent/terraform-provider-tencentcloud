## ADDED Requirements

### Requirement: cls_ckafka_consumer 新增 effective 参数
资源 `tencentcloud_cls_ckafka_consumer` SHALL 支持 `effective` 参数（TypeBool, Optional, Computed），用于控制投递任务是否生效。在 ModifyConsumer 请求中作为入参传递，在 DescribeConsumer 响应中作为出参读取。

#### Scenario: 创建资源时不指定 effective
- **WHEN** 用户创建 cls_ckafka_consumer 资源且未指定 effective 参数
- **THEN** CreateConsumer 请求中不传递 Effective 字段，资源创建成功后 Read 函数从 DescribeConsumer 响应中获取 Effective 值并设置到 state

#### Scenario: 更新资源时指定 effective 为 true
- **WHEN** 用户更新 cls_ckafka_consumer 资源并将 effective 设置为 true
- **THEN** ModifyConsumer 请求中传递 Effective=true，资源更新成功

#### Scenario: Read 函数读取 effective 值
- **WHEN** DescribeConsumer 响应中 Effective 字段非 nil
- **THEN** Read 函数将 Effective 值设置到 state 的 effective 字段

### Requirement: cls_ckafka_consumer 新增 role_arn 参数
资源 `tencentcloud_cls_ckafka_consumer` SHALL 支持 `role_arn` 参数（TypeString, Optional），用于设置角色访问描述名。在 CreateConsumer 和 ModifyConsumer 请求中作为入参传递，DescribeConsumer 响应不返回此字段。

#### Scenario: 创建资源时指定 role_arn
- **WHEN** 用户创建 cls_ckafka_consumer 资源并指定 role_arn
- **THEN** CreateConsumer 请求中传递 RoleArn 字段，值与用户指定的一致

#### Scenario: 更新资源时修改 role_arn
- **WHEN** 用户更新 cls_ckafka_consumer 资源并修改 role_arn 值
- **THEN** ModifyConsumer 请求中传递新的 RoleArn 字段值

#### Scenario: Read 函数不覆盖 role_arn
- **WHEN** 执行 terraform refresh 或 Read 函数被调用
- **THEN** role_arn 的值保持 state 中的值不变（因 DescribeConsumer 不返回此字段）

### Requirement: cls_ckafka_consumer 新增 external_id 参数
资源 `tencentcloud_cls_ckafka_consumer` SHALL 支持 `external_id` 参数（TypeString, Optional），用于设置外部ID。在 CreateConsumer 和 ModifyConsumer 请求中作为入参传递，DescribeConsumer 响应不返回此字段。

#### Scenario: 创建资源时指定 external_id
- **WHEN** 用户创建 cls_ckafka_consumer 资源并指定 external_id
- **THEN** CreateConsumer 请求中传递 ExternalId 字段，值与用户指定的一致

#### Scenario: 更新资源时修改 external_id
- **WHEN** 用户更新 cls_ckafka_consumer 资源并修改 external_id 值
- **THEN** ModifyConsumer 请求中传递新的 ExternalId 字段值

#### Scenario: Read 函数不覆盖 external_id
- **WHEN** 执行 terraform refresh 或 Read 函数被调用
- **THEN** external_id 的值保持 state 中的值不变（因 DescribeConsumer 不返回此字段）

### Requirement: cls_ckafka_consumer 新增 advanced_config 参数
资源 `tencentcloud_cls_ckafka_consumer` SHALL 支持 `advanced_config` 参数（TypeList, MaxItems:1, Optional），用于设置高级配置项。在 CreateConsumer 和 ModifyConsumer 请求中作为入参传递，DescribeConsumer 响应不返回此字段。

`advanced_config` 包含以下子字段：
- `partition_hash_status`（TypeBool, Optional）：CKafka分区hash状态，默认 false
- `partition_fields`（TypeList of TypeString, Optional）：需要计算hash的字段列表，最大支持5个字段

#### Scenario: 创建资源时指定 advanced_config
- **WHEN** 用户创建 cls_ckafka_consumer 资源并指定 advanced_config 块
- **THEN** CreateConsumer 请求中传递 AdvancedConfig 字段，其中 PartitionHashStatus 和 PartitionFields 值与用户指定的一致

#### Scenario: 更新资源时修改 advanced_config
- **WHEN** 用户更新 cls_ckafka_consumer 资源并修改 advanced_config 块中的字段
- **THEN** ModifyConsumer 请求中传递新的 AdvancedConfig 字段值

#### Scenario: Read 函数不覆盖 advanced_config
- **WHEN** 执行 terraform refresh 或 Read 函数被调用
- **THEN** advanced_config 的值保持 state 中的值不变（因 DescribeConsumer 不返回此字段）

### Requirement: Update 函数检测新参数变更
资源 Update 函数 SHALL 将 `effective`、`role_arn`、`external_id`、`advanced_config` 加入 mutableArgs 列表，以检测这些参数的变更并触发 ModifyConsumer 调用。

#### Scenario: effective 变更触发更新
- **WHEN** 用户修改 effective 参数值
- **THEN** Update 函数检测到变更，构建 ModifyConsumer 请求并传递新的 Effective 值

#### Scenario: role_arn 变更触发更新
- **WHEN** 用户修改 role_arn 参数值
- **THEN** Update 函数检测到变更，构建 ModifyConsumer 请求并传递新的 RoleArn 值

#### Scenario: external_id 变更触发更新
- **WHEN** 用户修改 external_id 参数值
- **THEN** Update 函数检测到变更，构建 ModifyConsumer 请求并传递新的 ExternalId 值

#### Scenario: advanced_config 变更触发更新
- **WHEN** 用户修改 advanced_config 参数值
- **THEN** Update 函数检测到变更，构建 ModifyConsumer 请求并传递新的 AdvancedConfig 值
