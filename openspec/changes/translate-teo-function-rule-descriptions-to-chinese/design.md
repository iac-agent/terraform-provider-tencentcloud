## Context

`tencentcloud_teo_function_rule` 是 TEO（边缘安全加速平台）产品下的一个 Terraform 资源，用于管理边缘函数触发规则。该资源已完整实现 CRUD 逻辑，当前所有 Schema 字段的 `Description` 均为英文。本次变更仅将字段描述翻译为中文，不改动任何业务逻辑。

## Goals / Non-Goals

**Goals:**
- 将 `resource_tc_teo_function_rule.go` 中所有 Schema 字段的 `Description` 从英文翻译为中文
- 保持翻译后的描述准确、简洁，符合中文技术文档习惯

**Non-Goals:**
- 不修改任何业务逻辑代码（CRUD 函数）
- 不修改任何云 API 接口调用
- 不修改资源 Schema 结构（字段名、类型、Required/Optional/Computed 等属性）
- 不修改测试文件
- 不修改 `.md` 文档文件

## Decisions

1. **仅修改 Description 字符串**：因为这是纯文档性质的变更，不涉及任何功能行为，所以只需修改 `Description` 字段值即可。不需要修改 Schema 结构、Provider 注册代码等。

2. **使用中文技术术语翻译**：参考腾讯云官方文档中使用的术语（如 "站点"、"功能"、"边缘函数"、"操作符" 等），确保翻译与产品文档一致。

3. **保留枚举值原文**：`operator` 和 `target` 字段的 Description 中包含枚举值名称（如 `equals`、`filename`），这些是 API 层面的标识符，保持不变，仅翻译其说明文字。

## Risks / Trade-offs

- **风险**：无。此变更仅修改 Description 字符串，不影响任何功能行为、API 调用或状态管理。
- **向后兼容性**：完全兼容。Terraform 的 Description 字段仅用于文档展示，不影响配置解析、状态管理或 API 调用。