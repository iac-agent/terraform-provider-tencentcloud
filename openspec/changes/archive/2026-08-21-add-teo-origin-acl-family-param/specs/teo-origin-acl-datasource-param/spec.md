## ADDED Requirements

### Requirement: 数据源必须返回 origin_acl_family 参数

Terraform 数据源 `tencentcloud_teo_origin_acl` SHALL 返回 `origin_acl_family` 参数，该参数表示源站防护回源ACL控制域。

该参数 MUST 满足以下要求：
- 参数名称: `origin_acl_family`
- 参数类型: `string`
- 参数属性: `Computed` (只读)
- 数据来源: 云API `DescribeOriginACL` 的响应 `Response.OriginACLInfo.OriginACLFamily`

#### Scenario: 成功读取源站ACL信息并包含 origin_acl_family

- **WHEN** 用户配置 `tencentcloud_teo_origin_acl` 数据源并指定有效的 `zone_id`
- **THEN** 数据源 SHALL 成功返回源站ACL信息，且返回结果中 MUST 包含 `origin_acl_info` 列表
- **AND** `origin_acl_info` 列表中的每个元素 MUST 包含 `origin_acl_family` 字段
- **AND** `origin_acl_family` 字段的值 SHALL 来自云API响应的 `OriginACLInfo.OriginACLFamily`

#### Scenario: 云API返回空值处理

- **WHEN** 云API `DescribeOriginACL` 的响应中 `OriginACLInfo.OriginACLFamily` 为 `nil` 或空值
- **THEN** 数据源 SHALL 不在 `origin_acl_info` 中设置 `origin_acl_family` 字段
- **AND** 数据源 SHALL 正常返回其他字段的值

#### Scenario: 参数在Terraform状态中的持久化

- **WHEN** 用户成功应用包含 `tencentcloud_teo_origin_acl` 数据源的配置
- **THEN** `origin_acl_family` 参数的值 SHALL 被正确保存在Terraform state 中
- **AND** 后续执行 `terraform refresh` 或 `terraform plan` SHALL 正确读取该参数的值

## ADDED Requirements

### Requirement: 数据源文档必须包含 origin_acl_family 参数说明

通过 `make doc` 生成的 `tencentcloud_teo_origin_acl` 数据源文档 SHALL 包含 `origin_acl_family` 参数的说明。

#### Scenario: 文档自动生成包含新参数

- **WHEN** 执行 `make doc` 命令生成文档
- **THEN** 生成的 `website/docs/d/tencentcloud_teo_origin_acl.html.markdown` 文件 MUST 包含 `origin_acl_family` 参数的说明
- **AND** 参数说明 SHALL 包含参数类型、是否必填、描述信息
