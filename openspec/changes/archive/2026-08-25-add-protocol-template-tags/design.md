## Context

当前 `tencentcloud_protocol_template` 资源使用 VPC API 的 `CreateServiceTemplate`、`DescribeServiceTemplates`、`ModifyServiceTemplateAttribute` 和 `DeleteServiceTemplate` 接口管理协议端口模板。云 API 的 `CreateServiceTemplate` 请求支持 `Tags` 参数（类型为 `[]*Tag`），`ServiceTemplate` 响应结构包含 `TagSet` 字段（类型为 `[]*Tag`），但 Terraform 资源未利用这些字段。

Tag 结构体定义：
```go
type Tag struct {
    Key   *string `json:"Key,omitempty"`
    Value *string `json:"Value,omitempty"`
}
```

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_protocol_template` 资源新增 `Key` 和 `Value` 可选参数
- 在创建资源时支持传入 Tags 参数到 `CreateServiceTemplate` 接口
- 在读取资源时从 `DescribeServiceTemplates` 接口返回中提取 TagSet 信息并设置到 state
- 保持完全向后兼容，新增参数为可选参数

**Non-Goals:**
- 不支持批量标签操作（仅支持单个 Key/Value 对）
- 不修改 Update 方法（标签不支持修改，如需修改需重新创建资源）
- 不新增数据源（data source）支持

## Decisions

### 1. 参数设计：使用单独的 Key 和 Value 参数

**决策**：在 Terraform schema 中新增 `Key` 和 `Value` 两个独立的字符串参数，类型为 `schema.TypeString`，均为可选参数（Optional）。

**理由**：
- 云 API 的 Tags 是 `[]*Tag` 数组，但 Terraform 资源通常只需要支持单个标签
- 使用单独的 `Key` 和 `Value` 参数比使用复杂的嵌套结构更符合 Terraform 的使用习惯
- 参考其他 TencentCloud Terraform 资源的标签实现方式

**替代方案**：
- 使用 `TypeList` 或 `TypeSet` 嵌套结构支持多个标签：过于复杂，当前需求仅需支持单个标签
- 不支持标签功能：无法满足用户需求

### 2. 服务层方法修改

**决策**：修改 `CreateServiceTemplate` 方法，新增 `key` 和 `value` 参数；修改 `DescribeServiceTemplateById` 方法，返回 `TagSet` 信息。

**理由**：
- `CreateServiceTemplate` 需要在创建时传入 Tags 参数
- `DescribeServiceTemplateById` 需要返回 TagSet 以便在 Read 方法中设置到 state

### 3. 不支持 Update 操作

**决策**：标签参数仅支持在创建时设置，不支持在 Update 中修改。

**理由**：
- 云 API 的 `ModifyServiceTemplateAttribute` 接口不支持修改 Tags 参数
- 如需修改标签，用户可以通过 `terraform destroy` 和 `terraform apply` 重新创建资源
- 这符合 Terraform 的资源管理理念

## Risks / Trade-offs

**风险 1**：用户可能期望支持多个标签
- **缓解措施**：当前实现支持单个标签，如有需求可后续扩展为列表形式

**风险 2**：标签参数在创建后不可修改可能导致用户困惑
- **缓解措施**：在文档中明确说明标签仅在创建时设置，如需修改需重新创建资源

**权衡**：选择简单实现（单个 Key/Value）而非完整实现（标签列表），以降低复杂度并快速满足需求
