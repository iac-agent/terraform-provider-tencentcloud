# config-list-aggregators-datasource Specification

## Purpose
TBD - created by archiving change add-config-list-aggregators-datasource. Update Purpose after archive.
## Requirements
### Requirement: 查询账号组列表数据源

系统 SHALL 提供数据源 `tencentcloud_config_list_aggregators`，通过云 API `ListAggregators`（config v20220802）查询腾讯云配置审计账号组列表，并暴露以下 computed 属性：`total`、`items`（含 `name`、`description`、`owner_uin`、`create_time`、`account_count`、`type`、`account_group_id`、`aggregator_status`、`member_name`）以及 `result_output_file`。

#### Scenario: 成功查询账号组列表

- **WHEN** 用户引用 `data "tencentcloud_config_list_aggregators"` 且云 API 正常返回
- **THEN** 数据源输出 `total` 等于 `Response.Total`，`items` 列表逐项映射 `Response.Items` 中的 `Name/Description/OwnerUin/CreateTime/AccountCount/Type/AccountGroupId/AggregatorStatus/MemberName` 字段

#### Scenario: API 响应为空

- **WHEN** 云 API 调用返回 `Response` 为 nil
- **THEN** 重试逻辑按 `tccommon.ReadRetryTimeout` 处理，且不清空数据源 id（不调用 `d.SetId("")`）

#### Scenario: 输出到文件

- **WHEN** 用户设置 `result_output_file`
- **THEN** 查询结果写入指定文件

