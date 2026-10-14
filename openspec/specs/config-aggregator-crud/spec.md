# config-aggregator-crud Specification

## Purpose

Terraform 资源 `tencentcloud_config_aggregator`（Config 账号组）的创建、查询、编辑、删除能力规范，覆盖与云 API `CreateAggregator` / `DescribeAggregator` / `UpdateAggregator` / `DeleteAggregators`（config/v20220802）的交互约定，以及资源注册与文档要求。

## Requirements

### Requirement: 创建账号组（Create）

TF 资源 `tencentcloud_config_aggregator` 创建时 MUST 调用云 API `CreateAggregator`（config/v20220802），传入 `name`、`description`、`type`、可选 `aggregator_accounts`（member_uin/member_name），成功后从 `Response.AccountGroupId` 获取 id 并以 `<owner_uin><FILED_SP><account_group_id>` 联合格式 `d.SetId()`。

- **Response/AccountGroupId 为 nil 或空字符串时返回 NonRetryableError，并打印 logId**。
- retry 使用 `tccommon.WriteRetryTimeout`，retry 块内仅调用 API。

#### Scenario: 成功创建账号组

- **WHEN** 用户 apply 一份含 `tencentcloud_config_aggregator` 的配置（type 为 RD 或 CUSTOM）
- **THEN** 调用 CreateAggregator 创建账号组，资源 id 为 `<owner_uin>;<account_group_id>` 格式，Read 回填 name/description/type/aggregator_accounts/account_group_id/aggregator_status。

### Requirement: 查询账号组（Read）

Read MUST 使用联合 id 拆出 `owner_uin`、`account_group_id`，调用 `DescribeAggregator`，将 `Name`、`Description`、`Type`、`AggregatorAccounts`、`AggregatorStatus` 回填至 state；所有字段 set 前判 nil。

- 资源不存在（response 为空）时先打印 `log.Printf("[CRUD] config_aggregator id=%s", d.Id())` 再 `d.SetId("")`。

#### Scenario: 资源已被手动删除

- **WHEN** Read 时 DescribeAggregator 返回空
- **THEN** 打印保留现场的日志，清空 id，下一次 plan 重新创建。

### Requirement: 编辑账号组（Update）

Update MUST 调用 `UpdateAggregator`，传入 name、description、aggregator_accounts（来自 schema）及从 id 拆出的 account_group_id、owner_uin。

#### Scenario: 修改账号组名称与成员

- **WHEN** 修改 `name` 或 `aggregator_accounts`
- **THEN** 调用 UpdateAggregator 并随后 Read 刷新 state。

### Requirement: 删除账号组（Delete）

Delete MUST 调用 `DeleteAggregators`，传入从 id 拆出的 account_group_id 与 owner_uin。

#### Scenario: terraform destroy

- **WHEN** 用户执行 destroy
- **THEN** 调用 DeleteAggregators 删除账号组，资源从 state 移除。

### Requirement: 资源注册与文档

- SHALL 在 `tencentcloud/provider.go` 注册 `"tencentcloud_config_aggregator": config.ResourceTencentCloudConfigAggregator()`，同步更新 `tencentcloud/provider.md`。
- 新增 `resource_tc_config_aggregator.md`：一句话描述（含 config 产品名）、Example Usage、Import（联合 id）。

#### Scenario: provider 注册生效

- **WHEN** `terraform providers` / plan 中引用 `tencentcloud_config_aggregator`
- **THEN** 资源可被识别并创建。
