## Context

当前 Terraform Provider for TencentCloud 的 `tencentcloud/services/` 目录下约 2,183 个 Go 源文件中的 Terraform schema `Description` 字段几乎全部为英文（约 38,312 行），仅 1 行为中文描述。该 Provider 面向的主要用户群体为中文用户，英文描述增加了用户理解字段含义的成本。

Terraform 的 schema `Description` 字段仅用于：
1. 在 `terraform plan`/`apply` 输出中提供字段说明
2. 通过 `make doc` 生成 `website/docs/` 下的文档

Description 不参与 Terraform state 管理、配置解析或 API 调用，因此翻译操作是纯文本替换，不涉及任何功能变更。

## Goals / Non-Goals

**Goals:**
- 将 `tencentcloud/services/` 目录下所有 Go 源文件中 Terraform schema 的 `Description` 字段从英文翻译为简体中文
- 保持翻译的专业性和准确性，使用腾讯云产品文档中的标准术语
- 翻译后代码可正常编译，`make doc` 可正常生成中文文档

**Non-Goals:**
- 不翻译 `tencentcloud/provider.go` 中的 provider 级别配置描述（面向全球用户）
- 不翻译代码注释、日志输出、错误信息
- 不修改任何代码逻辑、schema 结构、字段类型或字段名称
- 不修改 `tencentcloud/` 下非 `services/` 目录的文件
- 不涉及测试代码（`*_test.go`）的修改

## Decisions

### 决策 1：翻译范围

**选择**：仅翻译 `tencentcloud/services/` 目录下 `*_tc_*.go` 文件中的 schema `Description` 字段。

**理由**：
- `provider.go` 的配置描述面向全球用户，应保持英文
- 测试文件中的描述不影响用户可见的文档和提示
- `services/` 目录下的资源/数据源是用户直接交互的入口

**替代方案**：翻译所有 Description 包括 provider.go → 被排除，因为 provider 配置需要面向国际化用户。

### 决策 2：翻译策略

**选择**：基于文件批量翻译，按服务和文件逐一处理。

**理由**：
- 约 38,312 行描述，量级大，需要分批处理
- 按文件/服务拆分便于 review 和回滚
- 确保每个文件翻译后仍可编译

**替代方案**：一次性全局替换 → 不可行，因为需要确保翻译质量和上下文准确性。

### 决策 3：翻译质量标准

**选择**：使用腾讯云产品文档中的标准中文术语，保持风格一致。

**理由**：
- 与腾讯云控制台和文档术语保持一致，降低用户认知负担
- 例如："instance id" → "实例 ID"，"VPC" → "私有网络"，"subnet" → "子网"

### 决策 4：不修改 schema 结构

**选择**：仅修改 `Description:` 后的字符串内容，不修改任何其他代码。

**理由**：
- 确保不引入任何向后兼容性问题
- Terraform schema 的 Description 是纯文档属性，不影响功能
- 降低风险，变更仅涉及字符串字面量

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 翻译不准确或不一致 | 参考腾讯云官方文档术语表；按服务逐文件 review |
| 翻译后文件编码问题 | 仅使用 UTF-8 编码的中文字符，Go 原生支持 |
| 翻译量大，可能遗漏 | 使用 grep 检查各服务目录下是否还有英文 Description |
| `make doc` 生成文档格式异常 | 翻译完成后执行 `make doc` 验证文档生成正常 |
| 翻译后字段描述过长影响显示 | 遵循与原文相近的长度，避免过度翻译 |

## Open Questions

- 是否需要建立术语对照表（glossary）以确保翻译一致性？建议在实施过程中积累。
- 翻译后的文档是否需要同时保留英文版本？当前设计为仅中文。