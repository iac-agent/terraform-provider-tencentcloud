## Context

现有 Terraform 资源 `tencentcloud_cls_ckafka_consumer` 位于 `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go`，当前已支持 `topic_id`、`need_content`、`content`、`ckafka`、`compression` 五个参数。云 API（CreateConsumer / ModifyConsumer / DescribeConsumer）已支持更多参数（`effective`、`role_arn`、`external_id`、`advanced_config`），但 Terraform 资源尚未暴露这些参数。

当前 DescribeConsumer 响应仅返回 `Effective`、`NeedContent`、`Content`、`Ckafka`、`Compression`，不返回 `RoleArn`、`ExternalId`、`AdvancedConfig`。因此 `role_arn`、`external_id`、`advanced_config` 三个参数在 Read 函数中无法从 API 恢复，属于只写参数（Write-only），Terraform 会在 refresh 时依赖本地 state 中的值。

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_cls_ckafka_consumer` 资源新增 `effective`、`role_arn`、`external_id`、`advanced_config` 四个 Optional 参数
- 在 Create 函数中支持 `role_arn`、`external_id`、`advanced_config` 入参
- 在 Update 函数中支持 `effective`、`role_arn`、`external_id`、`advanced_config` 变更
- 在 Read 函数中从 DescribeConsumer 响应读取 `effective` 字段
- 在 Update 的 mutableArgs 中添加新参数以检测变更
- 保持向后兼容：所有新增参数均为 Optional，不影响现有配置

**Non-Goals:**
- 不修改已有参数的 schema（不改变 Required/Optional/ForceNew 等属性）
- 不修改资源的 ID 逻辑（仍使用 topic_id）
- 不修改 provider.go 的注册逻辑（资源已注册）
- 不修改 DescribeClsCkafkaConsumerById 服务层函数（其已返回完整的 DescribeConsumerResponseParams）

## Decisions

### 1. `effective` 参数设计为 Optional + Computed

**决策**: `effective` 参数设为 Optional 且 Computed，因为 DescribeConsumer 响应中会返回该字段值。

**替代方案**: 仅设为 Optional 不设 Computed — 但这会导致首次 import 或 read 时无法从 API 获取该值，产生 state 不一致。

### 2. `role_arn`、`external_id` 设为 Optional（不设 Computed）

**决策**: 这两个参数仅存在于 CreateConsumer 和 ModifyConsumer 的请求中，DescribeConsumer 不返回。设为 Optional，不设 Computed。

**影响**: 执行 `terraform refresh` 后，state 中的 `role_arn` 和 `external_id` 不会被 API 响应覆盖，保持用户配置值。这是可接受的，因为这些参数主要用于配置投递权限，不需要从服务端回读。

### 3. `advanced_config` 设计为 TypeList（MaxItems:1）嵌套结构

**决策**: `advanced_config` 映射为 `AdvancedConsumerConfiguration` 结构体，包含 `partition_hash_status`（TypeBool）和 `partition_fields`（TypeList of TypeString）两个子字段。采用 TypeList MaxItems:1 的方式，与现有 `content`、`ckafka` 参数保持一致的风格。

**替代方案**: 将子字段平铺到顶层 — 但这不符合 Terraform 嵌套结构的惯例，且与云 API 的结构映射不清晰。

### 4. `advanced_config` 不设 Computed

**决策**: DescribeConsumer 不返回 `AdvancedConfig`，因此该参数不设 Computed。Read 函数中不从 API 响应设置该字段。

### 5. Update 函数中将新参数加入 mutableArgs

**决策**: 将 `effective`、`role_arn`、`external_id`、`advanced_config` 加入 mutableArgs 列表，使 Update 函数能检测这些参数的变更。

## Risks / Trade-offs

- **[只写参数 state 不一致]**: `role_arn`、`external_id`、`advanced_config` 为只写参数，DescribeConsumer 不返回这些字段。如果用户在控制台修改了这些值，Terraform 无法检测到 drift。风险较低，因为这些参数通常不会在控制台被修改。缓解措施：在文档中说明这些参数为只写参数。
- **[向后兼容]**: 所有新增参数均为 Optional，不影响现有配置和 state。风险极低。
