## Why

Terraform Provider for TencentCloud 的字段描述（Description）目前几乎全部为英文（约 38,312 行英文描述 vs 仅 1 行中文描述），而该 Provider 主要面向中国用户。将 schema 中的字段描述从英文转为中文，可以显著提升中���用户的配置体验，降低理解成本，使 Terraform 配置更易于编写和维护。

## What Changes

- 将 `tencentcloud/services/` 目录下所有 Go 源文件中 Terraform schema 定义的 `Description` 字段从英文翻译为简体中文
- 涉及约 2,183 个 Go 源文件、约 38,312 行 Description 描述
- 不修改 `tencentcloud/provider.go` 中的 provider 级别配置描述（这些面向全球用户）
- 不涉及任何代码逻辑、云 API 行为、schema 结构、字段类型或字段名称的变更
- 仅为描述性文本的翻译，**不涉及向后兼容性问题**（Description 仅用于文档生成和 terraform 提示，不影响 state 或配置解析）

## Capabilities

### New Capabilities
- `i18n-field-descriptions`: 将 Terraform schema 中字段的英文 Description 翻译为简体中文，使 terraform plan/apply 输出的字段提示和生成文档均为中文

### Modified Capabilities
<!-- 无现有 spec 需要修改，此变更仅涉及描述文本翻译，不影响任何 spec 级别的行为 -->

## Impact

- 受影响文件：`tencentcloud/services/` 目录下约 2,183 个 Go 源文件（资源文件 `resource_tc_*.go` 和数据源文件 `data_source_tc_*.go`）
- 受影响代码：Terraform schema 定义中的 `Description:` 字段值（纯文本字符串翻译）
- 不涉及任何 API 调用、业务逻辑、测试代码的变更
- 文档生成（`make doc`）将产出中文描述的 `website/docs/` 文档
- 不影响 provider 本身的功能行为和兼容性
