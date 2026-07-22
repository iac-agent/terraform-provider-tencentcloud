## ADDED Requirements

### Requirement: Terraform schema 字段描述为简体中文
`tencentcloud/services/` 目录下所有资源（resource）和数据源（datasource）Go 源文件中的 Terraform schema `Description` 字段 SHALL 使用简体中文编写。

#### Scenario: 资源文件中的 Description 为中文
- **WHEN** 用户查看任意资源文件（如 `resource_tc_vpc.go`）的 Terraform schema 定义
- **THEN** 所有 `Description:` 字段的值均为简体中文

#### Scenario: 数据源文件中的 Description 为中文
- **WHEN** 用户查看任意数据源文件（如 `data_source_tc_instances.go`）的 Terraform schema 定义
- **THEN** 所有 `Description:` 字段的值均为简体中文

#### Scenario: 生成的文档为中文
- **WHEN** 执行 `make doc` 命令生成文档
- **THEN** `website/docs/` 目录下生成的 Markdown 文档中 `Argument Reference` 和 `Attribute Reference` 的字段描述均为简体中文

### Requirement: 翻译不改变代码行为
Description 的翻译 SHALL 仅修改字符串字面量的内容，不修改任何代码逻辑、schema 结构、字段类型、字段名称或 API 调用。

#### Scenario: 翻译后代码可编译
- **WHEN** 翻译完成后
- **THEN** 代码可正常编译，无语法错误或类型错误

#### Scenario: 翻译后 schema 结构不变
- **WHEN** 翻译完成后，对比翻译前后的 schema 定义
- **THEN** 除 Description 字符串内容外，所有字段的 `Type`、`Required`、`Optional`、`Computed`、`ForceNew`、`Default` 等属性保持不变

### Requirement: Provider 级别配置描述保持不变
`tencentcloud/provider.go` 中的 provider 级别配置描述 SHALL 保持英文不变。

#### Scenario: provider.go 的 Description 仍为英文
- **WHEN** 翻译完成后
- **THEN** `tencentcloud/provider.go` 中所有 `Description:` 字段的值仍为英文

### Requirement: 翻译术语一致性
翻译 SHALL 使用腾讯云官方文档中的标准中文术语，确保各服务、各资源之间的术语一致性。

#### Scenario: 相同概念使用相同术语
- **WHEN** 多个资源或数据源中存在相同语义的字段（如 `instance_id`）
- **THEN** 这些字段的 Description 翻译使用一致的术语（如均翻译为"实例 ID"）