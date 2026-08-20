## Context

TEO (TencentCloud EdgeOne) 提供了源站防护功能，通过 `DescribeOriginACL` 接口可以查询源站ACL信息。该接口返回的 `OriginACLInfo` 结构中包含 `OriginACLFamily` 字段，用于指定源站防护回源ACL控制域。

当前 Terraform 数据源 `tencentcloud_teo_origin_acl` 已经实现了查询源站ACL信息的功能，但缺少对 `OriginACLFamily` 字段的映射和处理。

数据源文件位置：`tencentcloud/services/teo/data_source_tc_teo_origin_acl.go`

## Goals / Non-Goals

**Goals:**
- 在数据源 `tencentcloud_teo_origin_acl` 中新增 `OriginACLFamily` 参数
- 将云API返回的 `OriginACLInfo.OriginACLFamily` 字段正确映射到 Terraform 数据源的 schema 中
- 确保参数类型为 string，且为 Computed（只读）
- 保持向后兼容性，不影响现有 Terraform 配置

**Non-Goals:**
- 不涉及其他数据源或资源的修改
- 不涉及云API的调用逻辑修改（仅新增字段映射）
- 不涉及单元测试的修改（参数已在代码中实现，测试代码已存在）

## Decisions

### 1. 参数位置选择

**决策**: 将 `OriginACLFamily` 参数添加到 `origin_acl_info` 结构的顶层

**理由**:
- `OriginACLFamily` 是 `OriginACLInfo` 的直接字段，与 `Status`、`L7Hosts` 等字段同级
- 保持与云API响应结构的一致性
- 符合现有数据源的参数组织方式

**替代方案**:
- 创建独立的顶层参数：不符合数据结构层次，且与其他 `OriginACLInfo` 字段的组织方式不一致

### 2. 参数类型定义

**决策**: 使用 `schema.TypeString` 类型，设置为 `Computed: true`

**理由**:
- `OriginACLFamily` 在云API中为 string 类型
- 作为数据源参数，该字段为只读属性（由云API返回）
- `Computed: true` 表示该字段将由 provider 自动填充

### 3. 实现方式

**决策**: 在数据源的 `Read` 方法中添加对 `OriginACLFamily` 字段的处理

**理由**:
- 参考现有代码中其他字段的处理方式（如 `Status` 字段）
- 在设置 `origin_acl_info` map 时，检查 `respData.OriginACLInfo.OriginACLFamily` 是否为 nil，若不为 nil 则设置到 map 中
- 保持与现有代码风格一致

## Risks / Trade-offs

**[风险] 云API字段为空时的处理**
- **风险**: 如果云API返回的 `OriginACLFamily` 字段为 nil 或空字符串，可能导致数据源输出不一致
- **缓解措施**: 在 Read 方法中先检查字段是否为 nil，只有在不为 nil 时才设置到 map 中（已实现）

**[风险] 向后兼容性**
- **风险**: 新增字段理论上不影响向后兼容性，但需确保不会破坏现有 state
- **缓解措施**: 新增字段为 Computed 类型，不会影响现有 Terraform 配置的解析和状态管理

**[权衡] 参数必要性**
- **权衡**: 该参数已经在代码中存在（第309-313行和第478-480行），本次提案主要是规范化变更管理流程
- **说明**: 代码实现已完成，通过 openspec 流程确保变更的可追溯性和规范性
