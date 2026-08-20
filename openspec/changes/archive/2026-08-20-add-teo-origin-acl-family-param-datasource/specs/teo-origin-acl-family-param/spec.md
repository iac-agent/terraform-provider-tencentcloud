## ADDED Requirements

### Requirement: 数据源支持查询 OriginACLFamily 参数
Terraform 数据源 `tencentcloud_teo_origin_acl` SHALL 支持查询和返回 `OriginACLFamily` 参数，该参数表示源站防护回源ACL控制域信息。

数据源的 `origin_acl_info` 结构中 SHALL 包含一个名为 `origin_acl_family` 的字段，类型为 string，且为 Computed（只读）。

#### Scenario: 成功查询包含 OriginACLFamily 的源站ACL信息
- **WHEN** 用户调用 `tencentcloud_teo_origin_acl` 数据源查询指定 zone_id 的源站ACL信息
- **AND** 云API `DescribeOriginACL` 返回的 `OriginACLInfo.OriginACLFamily` 字段不为空
- **THEN** 数据源 SHALL 返回 `origin_acl_info.0.origin_acl_family` 字段，其值为云API返回的 `OriginACLFamily` 值

#### Scenario: 查询不包含 OriginACLFamily 的源站ACL信息
- **WHEN** 用户调用 `tencentcloud_teo_origin_acl` 数据源查询指定 zone_id 的源站ACL信息
- **AND** 云API `DescribeOriginACL` 返回的 `OriginACLInfo.OriginACLFamily` 字段为 nil 或空字符串
- **THEN** 数据源 SHALL 不包含 `origin_acl_info.0.origin_acl_family` 字段在输出中（或该字段为空字符串）

#### Scenario: 数据源 schema 定义正确
- **WHEN** Terraform 提供者初始化并加载 `tencentcloud_teo_origin_acl` 数据源
- **THEN** 数据源的 schema SHALL 包含 `origin_acl_info` 块
- **AND** `origin_acl_info` 块 SHALL 包含 `origin_acl_family` 字段
- **AND** `origin_acl_family` 字段的类型 SHALL 为 `schema.TypeString`
- **AND** `origin_acl_family` 字段的 `Computed` 属性 SHALL 为 `true`

### Requirement: 数据源 Read 方法正确处理 OriginACLFamily 字段
数据源的 Read 方法 SHALL 在调用云API `DescribeOriginACL` 后，正确处理响应中的 `OriginACLInfo.OriginACLFamily` 字段。

#### Scenario: Read 方法成功设置 OriginACLFamily 字段
- **WHEN** 数据源 Read 方法成功调用云API并获取到响应
- **AND** 响应中的 `respData.OriginACLInfo.OriginACLFamily` 不为 nil
- **THEN** Read 方法 SHALL 将 `OriginACLFamily` 的值设置到 `origin_acl_info` map 的 `origin_acl_family` 键中

#### Scenario: Read 方法处理 OriginACLFamily 为 nil 的情况
- **WHEN** 数据源 Read 方法成功调用云API并获取到响应
- **AND** 响应中的 `respData.OriginACLInfo.OriginACLFamily` 为 nil
- **THEN** Read 方法 SHALL 不设置 `origin_acl_family` 键到 `origin_acl_info` map 中
