# Proposal: add-config-list-aggregators-datasource

## Why

腾讯云配置审计（Config, v20220802）产品当前缺少查询"账号组列表"的 Terraform 数据源。云 API `ListAggregators` 已在 vendor SDK 中可用，需要新增数据源 `tencentcloud_config_list_aggregators`，以便用户在 Terraform 中查询账号组（Aggregator）的名称、描述、创建者 UIN、成员数量、状态等信息。

## What Changes

新增 RESOURCE_KIND_DATASOURCE 类型资源 `tencentcloud_config_list_aggregators`：

- 新增数据源实现文件 `tencentcloud/services/config/data_source_tc_config_list_aggregators.go`，调用云 API `ListAggregators`（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`）。
- 在 `ConfigService` 中新增服务层查询方法（带分页、retry）。
- 输出字段（云 API 出参 → SchemaName 映射，按列表项字段平铺到 schema 顶层）：
  - `response.Response.Total` → `total`
  - `response.Response.Items` → 各字段平铺（列表型字段 `items`，元素字段如下）
  - `Items[].Name` → `name`
  - `Items[].Description` → `description`
  - `Items[].OwnerUin` → `owner_uin`
  - `Items[].CreateTime` → `create_time`
  - `Items[].AccountCount` → `account_count`
  - `Items[].Type` → `type`
  - `Items[].AccountGroupId` → `account_group_id`
  - `Items[].AggregatorStatus` → `aggregator_status`
  - `Items[].MemberName` → `member_name`
  - 另按惯例增加 `result_output_file` 可选输出字段。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册数据源 `tencentcloud_config_list_aggregators`。
- 新增文档 `data_source_tc_config_list_aggregators.md`（"Use this data source to query ..."格式）。
- 新增单元测试 `data_source_tc_config_list_aggregators_test.go`（使用 gomonkey mock 云 API，不使用 terraform 测试套件）。

## Capabilities

### New Capabilities

- `config-list-aggregators-datasource`: 查询腾讯云配置审计账号组列表的数据源。

## Impact

- 仅新增文件与注册代码，不修改既有资源行为。
- 涉及文件：`tencentcloud/services/config/data_source_tc_config_list_aggregators.go`（新增）、`service_tencentcloud_config.go`（新增方法）、`tencentcloud/provider.go`、`tencentcloud/provider.md`、对应 `.md` 与 `_test.go`。
