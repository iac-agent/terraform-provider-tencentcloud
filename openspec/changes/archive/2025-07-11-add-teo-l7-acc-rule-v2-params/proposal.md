## Why

目前 Terraform Provider 的 `tencentcloud_teo_l7_acc_rule_v2` 资源在查询七层加速规则时，仅在内部通过 `rule-id` 构造过滤条件调用 DescribeL7AccRules API，但未将 `filters` 参数暴露为用户可配置的 Terraform schema 字段。

根据腾讯云 TEO SDK（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`）：
- **DescribeL7AccRules API**: `Filters` 参数支持按 `rule-id` 过滤规则

当前限制：
- 用户无法通过 Terraform 资源自定义查询过滤条件
- 过滤条件硬编码为 `rule-id`，缺乏灵活性
- SDK 支持该参数，但 Provider 未暴露

通过添加 `filters` 参数支持，用户可以：
- 自定义 DescribeL7AccRules 的查询过滤条件
- 通过基础设施即代码管理更完整的规则查询配置

## What Changes

在 `tencentcloud_teo_l7_acc_rule_v2` 资源中添加 `filters` 可选参数，该参数：
- **类型**: TypeList，包含 name（String）和 values（TypeList of String）的子结构
- **可选**: Optional
- **查询时**: 传递给 DescribeL7AccRules API 的 Filters 参数

### 修改文件
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` - 添加 Schema 定义和 Read 逻辑
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - 更新 DescribeTeoL7AccRuleById 方法以支持自定义 filters
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` - 添加测试用例覆盖
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` - 更新文档

### 资源 Schema 变更
```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example" {
  zone_id = "zone-3fkff38fyw8s"
  # ... 其他现有参数 ...

  filters {                    # 新增参数
    name   = "rule-id"
    values = ["rule-xxx"]
  }
}
```

### API 集成
- **查询时**: 将 `filters` 传递给 `DescribeL7AccRulesRequest.Filters`
- **创建时**: 无变更（filters 不参与创建）
- **更新时**: 无变更（filters 不参与更新）
- **删除时**: 无变更（filters 不参与删除）

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-filters`: TEO 七层加速规则 V2 资源的 filters 查询参数支持

### Modified Capabilities
（无）

## Impact

### 受影响的代码
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` - Schema 定义，Read 逻辑
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - DescribeTeoL7AccRuleById 方法支持自定义 filters
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` - 测试覆盖
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` - 文档更新

### 向后兼容性
- **完全向后兼容** - 参数为 Optional，现有配置无需修改
- **不影响现有资源** - 不修改已有资源的行为
- **SDK 已支持** - 腾讯云 SDK 已包含该字段，无需升级依赖

### API 兼容性
- `DescribeL7AccRules` API 支持 `Filters` 参数（可选）
- 参数为可选，不传递时使用现有内部构造的过滤条件

### 依赖关系
- 无新增依赖
- 使用现有的 SDK 版本（已包含该字段）
