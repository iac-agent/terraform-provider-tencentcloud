## 1. Schema 定义

- [x] 1.1 在 `resource_tc_teo_l7_acc_rule_v2.go` 的 Schema 中新增 `rule` 参数（TypeList, MaxItems=1, Optional），嵌套字段包含 `rule_id`（Computed）、`status`、`rule_name`、`description`、`branches`，branches 复用 `TencentTeoL7RuleBranchBasicInfo` 函数

## 2. CRUD 函数修改

- [x] 2.1 修改 `ResourceTencentCloudTeoL7AccRuleV2Create` 函数，支持从 `rule` 参数构建 `RuleEngineItem`，若 `rule` 未指定则继续使用现有的顶层字段（status、rule_name、description、branches）
- [x] 2.2 修改 `ResourceTencentCloudTeoL7AccRuleV2Read` 函数，在读取 API 响应后填充 `rule` 参数的嵌套字段值
- [x] 2.3 修改 `ResourceTencentCloudTeoL7AccRuleV2Update` 函数，支持从 `rule` 参数构建 `RuleEngineItem` 并设置到 `ModifyL7AccRule` 请求的 `Rule` 字段，若 `rule` 未指定则继续使用现有的顶层字段

## 3. 文档更新

- [x] 3.1 更新 `resource_tc_teo_l7_acc_rule_v2.md` 示例文件，添加 `rule` 参数的使用示例

## 4. 单元测试

- [x] 4.1 在 `resource_tc_teo_l7_acc_rule_v2_test.go` 中补充 `rule` 参数相关的单元测试用例，使用 gomonkey mock 云 API
- [x] 4.2 运行单元测试验证通过（使用 `go test -gcflags=all=-l` 参数）
