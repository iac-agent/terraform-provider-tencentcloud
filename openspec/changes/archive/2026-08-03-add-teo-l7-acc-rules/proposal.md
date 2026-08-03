## Why

TEO (EdgeOne) 产品缺少 Terraform 资源来管理七层加速规则（L7AccRules）。当前已有 `tencentcloud_teo_l7_acc_rule`（旧版单规则管理）和 `tencentcloud_teo_l7_acc_rule_v2`（新版单规则管理），但缺少一个支持批量管理七层加速规则的资源 `tencentcloud_teo_l7_acc_rules`，该资源使用新的 CRUD 接口（CreateL7AccRules/DescribeL7AccRules/ModifyL7AccRule/DeleteL7AccRules），允许用户在单个 Terraform 资源中批量管理多条加速规则。

## What Changes

- 新增 RESOURCE_KIND_GENERAL 资源 `tencentcloud_teo_l7_acc_rules`，支持七层加速规则的批量增删改查
- 使用新的云 API 接口：`CreateL7AccRules`（批量创建）、`DescribeL7AccRules`（查询）、`ModifyL7AccRule`（单条修改）、`DeleteL7AccRules`（批量删除）
- 在 `tencentcloud/provider.go` 中注册新资源

## Capabilities

### New Capabilities
- `teo-l7-acc-rules`: TEO 七层加速规则批量管理资源，支持 CRUD 操作，使用 zone_id 作为资源 ID，通过 rules 字段管理多条规则的集合

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md`
- 修改文件：`tencentcloud/provider.go`
- 新增文件：`tencentcloud/provider.md`（如需要）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已存在于 vendor）
