# Proposal: add-config-aggregator-resource

## Change Information

- **Change ID**: add-config-aggregator-resource
- **Type**: New Resource
- **Product**: config（Tencent Cloud Config，配置审计）
- **Resource Kind**: RESOURCE_KIND_GENERAL（通用资源，完整 CRUD 生命周期管理）
- **Terraform Resource Name**: `tencentcloud_config_aggregator`

## Why

腾讯云配置审计（Config）提供"账号组（Aggregator）"能力，用于将多个成员账号聚合到统一账号组中进行合规管理。当前 Terraform Provider 缺少对账号组资源的管理能力，用户无法通过 IaC 方式创建、查询、编辑和删除账号组。本次变更新增 `tencentcloud_config_aggregator` 资源，补齐该能力。

## What Changes

新增资源 `tencentcloud_config_aggregator`，包含：

### 云 API 映射（包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`）

| Terraform CRUD | 云 API | 说明 |
|---|---|---|
| Create | `CreateAggregator` | 新建账号组，出参 `AccountGroupId` 作为资源 id |
| Read | `DescribeAggregator` | 账号组详情 |
| Update | `UpdateAggregator` | 编辑账号组 |
| Delete | `DeleteAggregators` | 删除账号组 |

### Terraform 参数 Schema

参数（云 API JsonPath → SchemaName）：
- `request.Name` → `name`（Required，string）
- `request.Description` → `description`（Required，string）
- `request.Type` → `type`（Required，string；枚举 `RD` 全局账号组 / `CUSTOM` 自定义账号组）
- `request.AggregatorAccounts` → `aggregator_accounts`（Optional，TypeList）
  - `MemberUin` → `member_uin`（Required，int，成员 ID）
  - `MemberName` → `member_name`（Required，string，成员名称）
- `request.OwnerUin` → `owner_uin`（Required，int，账号组创建者 ID；Describe/Update/Delete 必传，Create 无此入参）
- `response.Response.AccountGroupId` → `account_group_id`（Computed，string）
- `response.Response.AggregatorStatus` → `aggregator_status`（Computed，int，创建状态）

### 资源 ID 设计

`owner_uin` 与 `account_group_id` 为 Read/Update/Delete 接口的联合必传参数，资源 id 使用 `tccommon.FILED_SP` 拼接为 `<owner_uin><FILED_SP><account_group_id>`，支持 import。

### 行为约束

- 所有 4 个云 API 均为同步接口（无 FlowId/TaskId），无需异步轮询。
- Create 返回后必须校验 `Response` 与 `AccountGroupId` 非空，为空返回 NonRetryableError；拿到 id 后先 `d.SetId()`。
- Read 中调用 `setXX()` 前需判断 Response 字段非 nil；若资源已不存在（Response 为空）打印 `log.Printf("[CRUD] config_aggregator id=%s", d.Id())` 后再 `d.SetId("")`。
- API 调用统一使用 `tccommon.WriteRetryTimeout`/`tccommon.ReadRetryTimeout` + `resource.Retry` + `tccommon.RetryError()` 包装。

## Capabilities

### New Capabilities

- `config-aggregator-crud`: 通过 Terraform 管理配置审计账号组的全生命周期（创建、查询、编辑、删除、导入）。

## Impact

- 新增 `tencentcloud/services/config/resource_tc_config_aggregator.go` 及配套 `_test.go`（gomonkey mock 单元测试）、`.md` 文档。
- `tencentcloud/provider.go` / `tencentcloud/provider.md` 注册 `tencentcloud_config_aggregator`。
- 不影响既有资源。
