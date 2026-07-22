## Why

当前 `tencentcloud_teo_function_rule` 资源的所有 Schema 字段描述（Description）均为英文，不利于中文用户理解和使用。为提升中文用户的体验，需要将所有字段描述翻译为中文。

## What Changes

- 将 `tencentcloud_teo_function_rule` 资源 Go 源文件中所有 Schema 字段的 `Description` 从英文翻译为中文，包括顶层字段和嵌套字段（`function_rule_conditions`、`rule_conditions` 等）。
- 不涉及任何业务逻辑修改，不涉及任何云 API 接口调用修改，仅修改 Description 字符串。

## Capabilities

### New Capabilities
<!-- 本次变更不引入新能力 -->
无

### Modified Capabilities
<!-- 本次变更不修改现有能力的需求 -->
无

## Impact

- 受影响文件：`tencentcloud/services/teo/resource_tc_teo_function_rule.go`
- 影响范围：仅 Schema 字段的 `Description` 属性值，为纯文档性质的变更
- 无 API 变更、无依赖变更、无破坏性变更