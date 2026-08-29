## Context

TEO (TencentCloud EdgeOne) 产品的 `tencentcloud_teo_origin_acl` 资源用于管理源站访问控制。当前资源实现已经支持了基本的源站 ACL 配置功能，包括 zone_id、l7_hosts、l4_proxy_ids 和 origin_acl_family 等参数。

根据云 API 的接口定义，EnableOriginACL、ModifyOriginACL 和 DescribeOriginACL 接口已经相对稳定。本变更旨在进一步增强资源的配置能力，可能需要新增一些高级配置参数以支持更多的使用场景。

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_teo_origin_acl` 资源新增必要的配置参数，提升资源的完整性和灵活性
- 保持与云 API 接口的一致性，确保所有 API 参数都有对应的 Terraform 配置选项
- 保持向后兼容性，不影响现有 Terraform 配置和 state

**Non-Goals:**
- 不修改已有的资源 schema 结构（除非只新增 Optional 字段）
- 不删除或重命名现有参数
- 不进行大规模的架构重构

## Decisions

### 决策 1: 参数新增策略
- **决策**: 基于云 API 接口的实际可用参数，逐步新增 Terraform 资源配置参数
- **理由**: 确保 Terraform 资源能够完整覆盖云 API 的功能，提供更好的用户体验
- **替代方案**: 一次性新增所有参数 vs 按需新增。选择按需新增以避免维护过多可能不常用的参数

### 决策 2: 实现方式
- **决策**: 遵循现有的 Terraform 资源实现模式，保持代码风格一致
- **理由**: 确保代码可维护性，减少学习成本
- **具体做法**: 参考现有的 `tencentcloud_igtm_strategy` 等资源实现方式

## Risks / Trade-offs

- **[风险] 新增参数可能导致资源 schema 变更** → 缓解措施：只新增 Optional 字段，确保向后兼容
- **[风险] 云 API 接口变更可能影响实现** → 缓解措施：仔细核对 API 文档和 SDK 定义
- **[权衡] 新增参数增加代码复杂度** → 收益：提升资源功能完整性，改善用户体验

## Open Questions

1. 具体需要新增哪些参数？需要根据实际的业务需求和云 API 能力来确定
2. 是否需要更新相关的数据源以支持新增参数的读取？
3. 新增参数是否需要在文档中特别说明使用场景和限制？
