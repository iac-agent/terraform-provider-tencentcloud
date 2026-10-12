# tag-category-param Specification

## Purpose
为 `tencentcloud_tag` 资源新增 `Category` Computed 参数，用于展示标签类型（Custom/System/All）。该参数由 `GetTags` 接口返回，在 Read 流程中回填至 Terraform state，用户不可配置，不影响创建与删除行为。
## Requirements
### Requirement: tencentcloud_tag 资源暴露 Category 参数

`tencentcloud_tag` 资源的 schema 中 SHALL 包含 `Category` 参数（Type: String, Computed），用于展示标签的类型属性。该参数 MUST NOT 出现在 `CreateTag` 请求中（云 API 创建接口无此入参），用户不可通过配置文件设置该参数。

参数取值 SHALL 为以下之一（由 `GetTags` 接口返回）：
- `Custom`：自定义标签
- `System`：系统标签
- `All`：全部标签

#### Scenario: schema 中存在 Category 参数

- **WHEN** 用户查看 `tencentcloud_tag` 资源 schema
- **THEN** schema 中存在名为 `Category` 的 String 类型 Computed 参数，且该参数不是 Required、不是 Optional、不是 ForceNew

#### Scenario: 创建资源时不传 Category

- **WHEN** 用户通过 `terraform apply` 创建 `tencentcloud_tag` 资源
- **THEN** `CreateTag` 请求体中不包含 `Category` 字段，创建行为与变更前完全一致

### Requirement: Read 流程回填 Category 字段

`tencentcloud_tag` 资源的 Read 方法 SHALL 调用现有服务层方法 `DescribeTagResourceById`（内部通过 `GetTags` API 获取标签），并当返回的 `tag.Tag` 结构体的 `Category` 字段非 nil 时，MUST 将该值通过 `d.Set("Category", ...)` 写入 Terraform state。

当 `Category` 字段为 nil 时（云 API 可能返回 null），Read 方法 MUST 跳过 `d.Set("Category", ...)` 调用且不产生错误。

#### Scenario: GetTags 返回非空 Category

- **WHEN** `GetTags` 响应中匹配 tagKey/tagValue 的 `Tag` 元素包含非空 `Category` 字段（如 `Custom`）
- **THEN** Terraform state 中 `Category` 属性被设置为该值（`Custom`）

#### Scenario: GetTags 返回空 Category

- **WHEN** `GetTags` 响应中匹配元素的 `Category` 字段为 null
- **THEN** Read 方法跳过 `Category` 的设置，不向 state 写入该字段，且 Read 不返回错误

#### Scenario: 现有标签的属性读取不受影响

- **WHEN** 对已存在的 `tencentcloud_tag` 资源执行 `terraform plan`（refresh）
- **THEN** `tag_key` 与 `tag_value` 的读取行为保持不变，`Category` 作为新增 Computed 字段被自动填充，plan 不产生破坏性变更（无 force replacement）

