## Context

当前 `tencentcloud_teo_origin_acl` 数据源用于查询 TEO 站点的源站防护 ACL 配置信息。该数据源调用 `DescribeOriginACL` 接口获取配置，但在现有的实现中，虽然 `OriginACLInfo` 结构中包含 `OriginACLFamily` 字段（表示源站防护回源ACL控制域），但该字段未被添加到 Terraform 数据源的输出参数中。

通过检查代码发现：
- 云 API 的 `OriginACLInfo` 结构已包含 `OriginACLFamily *string` 字段
- 现有数据源代码中，schema 定义已包含 `origin_acl_family` 参数（第309-313行）
- 现有数据源代码中，读取逻辑已实现（第478-480行）
- 但可能缺少文档或需要验证完整性

## Goals / Non-Goals

**Goals:**
- 确保 `origin_acl_family` 参数在数据源中正确定义和输出
- 验证参数从云 API 响应到 Terraform 状态的完整数据流
- 确保文档正确生成

**Non-Goals:**
- 不修改其他数据源或资源
- 不修改云 API 调用逻辑（仅需确认读取正确）
- 不改变参数的类型或行为

## Decisions

### 决策 1: 参数定义位置

**决定**: 在 `origin_acl_info` 嵌套结构中添加 `origin_acl_family` 参数

**理由**: 该参数属于 `OriginACLInfo` 对象的一部分，放在嵌套结构中与现有参数组织方式一致，便于用户理解参数间的层级关系。

**替代方案**:
- 将参数提升为顶层参数：会破坏现有结构的一致性，且不符合云 API 返回数据的层次结构

### 决策 2: 参数类型

**决定**: 使用 `schema.TypeString` 和 `Computed: true`

**理由**:
- `OriginACLFamily` 在云 API 中是 string 类型
- 作为数据源输出参数，应使用 Computed 而非 Optional
- 用户不需要（也不能）在配置中设置此值

### 决策 3: 实现验证

**决定**: 验证现有代码完整性，必要时补充

**理由**: 通过代码检查，发现 `origin_acl_family` 参数似乎已经实现。需要验证：
1. Schema 定义是否正确
2. 读取逻辑是否完整
3. 是否有 nil 检查等防御性编程

## Risks / Trade-offs

**[风险] 参数已存在但未正确工作** → **缓解措施**: 验证现有实现，确保读取逻辑正确处理 nil 情况和类型转换

**[风险] 文档未自动生成** → **缓解措施**: 在收尾阶段执行 `make doc` 确保文档生成

**[风险] 向后兼容性** → **缓解措施**: 仅新增 Computed 参数，不影响现有配置和 state
