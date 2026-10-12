# Change: 为 tencentcloud_tag 资源新增 Category 参数

## Why

Terraform Provider 的 `tencentcloud_tag` 资源目前仅暴露 `tag_key` 和 `tag_value` 两个参数。腾讯云标签服务 `GetTags` 接口返回的 `Tag` 结构体中已包含 `Category` 字段（标签类型：Custom 自定义标签 / System 系统标签 / All 全部标签），但 Provider 未将该字段暴露给用户，导致用户无法通过 Terraform 感知标签的类型属性（自定义标签或系统标签），状态中缺少该信息也不利于对标签资产的审计与区分。

云 API SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813/models.go`）中 `Tag` 结构体已支持 `Category` 字段，Provider 侧只需在 Read 流程中回填即可，成本低、完全向后兼容。

**本次变更仅新增一个参数 `Category`**：
- 需求中列出的 `CreateTag` 入参 `TagKey`/`TagValue` 与 `GetTags` 出参 `TagKey`/`TagValue`，经核对资源现有 schema（`tencentcloud/services/tag/resource_tc_tag.go`），`tag_key` 与 `tag_value` 参数**已存在**（Required + ForceNew，且 Read 中已回填），无需重复新增。
- 唯一需要新增的参数为 `GetTags` 出参 `Tag.Category`（JsonPath: `response.Tags.Category`）→ Terraform 参数 `Category`。

## What Changes

在 `tencentcloud_tag` 资源中新增 `Category` 参数：

- **参数名**: `Category`（遵循需求指定的 SchemaName，保留原始大小写）
- **类型**: String
- **属性**: Computed（云 API 出参，仅由 `GetTags` 响应回填，用户不可设置）
- **取值**: `Custom`（自定义标签）、`System`（系统标签）、`All`（全部标签）
- **回填逻辑**: 资源 Read 时，`DescribeTagResourceById` 返回的 `tag.Tag` 结构体中的 `Category` 字段非 nil 时，`d.Set("Category", ...)` 写入状态；为 nil 时（SDK 注释标注"此字段可能返回 null"）跳过设置，不报错

同步更新：

- `tencentcloud/services/tag/resource_tc_tag.go` — schema 新增 `Category` 字段，Read 中回填
- `tencentcloud/services/tag/resource_tc_tag.md` — 资源文档补充 `Category` 参数说明
- `tencentcloud/services/tag/resource_tc_tag_test.go` — 补充测试用例（沿用现有 terraform 测试套件风格，修改现有资源）

不涉及 Create/Delete 流程变更：`CreateTag` 接口无 `Category` 入参（已在 vendor SDK 中核对），标签创建后类型由云平台决定，参数仅作读取展示。

## Capabilities

### New Capabilities

- `tag-category-param`: `tencentcloud_tag` 资源的 `Category` 标签类型参数（Computed），定义 Read 流程中从 `GetTags` 响应回填 `Category` 的行为要求

### Modified Capabilities

（无 — `openspec/specs/` 下不存在 tag 相关既有规范，本次为新增能力规范）

## Impact

### 受影响的代码

- `tencentcloud/services/tag/resource_tc_tag.go` — Schema 新增 `Category` 参数（Computed）；`resourceTencentCloudTagResourceRead` 中回填 `Category`
- `tencentcloud/services/tag/resource_tc_tag_test.go` — 补充覆盖 `Category` 的测试用例
- `tencentcloud/services/tag/resource_tc_tag.md` — 文档补充参数说明与示例输出

无需改动：

- `service_tencentcloud_tag.go` — `DescribeTagResourceById` 返回 `*tag.Tag` 已包含 `Category` 字段，无需修改服务层
- `provider.go` / `provider.md` — 资源已注册，本次不新增资源
- vendor — SDK 已包含 `Category` 字段，无需升级依赖

### API 兼容性

- ✅ `GetTags` 响应 `Tags[].Category` 字段已存在于当前 vendored SDK（`Tag` 结构体，models.go:2446-2456）
- ✅ `CreateTag` 请求无 `Category` 入参，参数为 Computed，不影响创建行为
- ✅ 完全向后兼容：新增 Computed 字段不影响现有配置与 state

### 注意事项

- 参数名按需求指定为 `Category`（非 terraform 常规小写蛇形命名），schema key 以需求中的 SchemaName 为准
- SDK 注释标注 `Category` "此字段可能返回 null"，Read 中必须判空后再 Set
