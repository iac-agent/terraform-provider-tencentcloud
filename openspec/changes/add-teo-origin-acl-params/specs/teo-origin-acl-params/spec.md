## ADDED Requirements

### Requirement: 新增源站ACL资源配置参数
系统 SHALL 为 `tencentcloud_teo_origin_acl` 资源新增必要的配置参数，以支持更多源站访问控制配置选项。

#### Scenario: 成功创建带有新增参数的源站ACL资源
- **WHEN** 用户创建 `tencentcloud_teo_origin_acl` 资源并指定新增的可选参数
- **THEN** 资源 SHALL 成功创建，并且所有指定的参数 SHALL 被正确应用到云 API 调用中

#### Scenario: 成功更新源站ACL资源配置
- **WHEN** 用户更新 `tencentcloud_teo_origin_acl` 资源的新增参数
- **THEN** 资源 SHALL 成功更新，并且更新后的参数 SHALL 在后续读取操作中正确返回

#### Scenario: 读取源站ACL资源配置
- **WHEN** 用户读取 `tencentcloud_teo_origin_acl` 资源的状态
- **THEN** 资源 SHALL 返回所有配置的参数值，包括新增的参数

### Requirement: 保持向后兼容性
系统 SHALL 确保新增参数不会影响现有 Terraform 配置和 state 的兼容性。

#### Scenario: 现有配置继续正常工作
- **WHEN** 用户使用不包含新增参数的现有 Terraform 配置
- **THEN** 配置 SHALL 继续正常工作，所有现有功能 SHALL 保持可用

#### Scenario: 新增参数为可选
- **WHEN** 用户创建或更新资源时未指定新增参数
- **THEN** 新增参数 SHALL 使用默认值或保持为空，不影响资源的基本功能

### Requirement: 更新相关数据源
系统 SHALL 更新相关的数据源以支持新增参数的读取和展示。

#### Scenario: 数据源正确返回新增参数
- **WHEN** 用户通过数据源查询源站ACL信息
- **THEN** 数据源 SHALL 正确返回所有支持的参数信息，包括新增的参数

### Requirement: 文档更新
系统 SHALL 更新相关文档以说明新增参数的用法和限制。

#### Scenario: 资源文档包含新增参数说明
- **WHEN** 用户查看 `tencentcloud_teo_origin_acl` 资源的文档
- **THEN** 文档 SHALL 包含新增参数的详细说明、使用场景和示例

#### Scenario: 数据源文档包含新增参数说明
- **WHEN** 用户查看相关数据源的文档
- **THEN** 文档 SHALL 包含新增参数的说明和读取方式
