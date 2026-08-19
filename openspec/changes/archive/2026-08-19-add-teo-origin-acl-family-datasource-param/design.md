## Context

TEO (TencentCloud EdgeOne) 的 `tencentcloud_teo_origin_acl` 数据源用于查询站点的回源 ACL 配置。当前数据源已经实现了对 `DescribeOriginACL` API 的调用，并且云 API 的返回结构体 `OriginACLInfo` 中已包含 `OriginACLFamily` 字段，该字段用于标识回源 ACL 的控制域（如全局、中国大陆、全球不包括中国大陆等）。

经过代码检查，发现：
1. 数据源 schema 中已定义 `origin_acl_family` 参数（第 309-313 行）
2. Read 方法中已实现对该字段的读取（第 478-480 行）
3. 云 API vendor 包中 `OriginACLInfo` 结构体已包含 `OriginACLFamily *string` 字段

本设计文档旨在确认当前实现的正确性，并确保该功能完整可用。

## Goals / Non-Goals

**Goals:**
- 确保 `tencentcloud_teo_origin_acl` 数据源正确暴露 `origin_acl_family` 参数
- 验证数据源 Read 方法正确读取并设置 `OriginACLFamily` 字段
- 确保文档正确生成并包含该参数说明
- 验证单元测试覆盖（如需要）

**Non-Goals:**
- 不修改资源 `tencentcloud_teo_origin_acl`（本变更仅针对数据源）
- 不修改云 API 调用逻辑（API 已支持该字段）
- 不修改 `OriginACLFamily` 字段的类型或语义

## Decisions

### 决策 1：确认现有实现的有效性

**决定**: 经代码验证，数据源已正确实现 `origin_acl_family` 参数的 schema 定义和 Read 逻辑。

**理由**:
- Schema 定义符合 Terraform 数据源规范（Computed, TypeString）
- Read 方法中已包含 nil 检查（`if respData.OriginACLInfo.OriginACLFamily != nil`）
- 使用 `d.Set()` 正确设置字段值

**替代方案**: 无需替代方案，现有实现已满足需求。

### 决策 2：文档生成策略

**决定**: 使用 `make doc` 命令自动生成文档，而非手动编写 `.md` 文件。

**理由**:
- 根据禁止事项第 7 条，禁止直接新增/修改 `website/` 目录下的文件
- 文档应由 `make doc` 命令统一生成
- 这确保了文档与代码的一致性

**替代方案**: 无（必须遵守禁止事项）。

### 决策 3：单元测试策略

**决定**: 对于数据源，使用 gomonkey 进行 mock 测试，而非 Terraform 测试套件。

**理由**:
- 根据 go代码生成要求第 1 条，对于 RESOURCE_KIND_DATASOURCE 类型，应使用 mock 方法进行单元测试
- 避免依赖真实的云 API 调用
- 提高测试执行速度

**替代方案**: 使用 Terraform 验收测试（需要 TF_ACC=1 和真实凭证），但不符合当前代码生成规范。

## Risks / Trade-offs

### Risk 1: 字段值为空的情况
**风险**: 如果云 API 返回的 `OriginACLFamily` 为 nil 或空字符串，数据源可能不会设置该字段。

**缓解措施**:
- 代码中已包含 nil 检查（`if respData.OriginACLInfo.OriginACLFamily != nil`）
- Terraform 对于 Computed 字段，如果未设置会返回空字符串，这是可接受的行为
- 已在 Read 方法中正确处理

### Risk 2: 文档生成失败
**风险**: `make doc` 可能因为代码问题而生成失败。

**缓解措施**:
- 确保 Go 代码可编译（不执行 `go build`，但代码语法必须正确）
- 遵循现有的数据源代码模式
- 在收尾阶段由 tfpacer-finalize skill 执行 `make doc`

### Risk 3: 与现有配置的兼容性
**风险**: 新增字段可能影响现有 Terraform 配置。

**缓解措施**:
- `origin_acl_family` 是 Computed 字段，不会改变现有配置的语义
- 完全向后兼容，不影响现有 state
- 现有配置无需修改即可继续工作

## Implementation Notes

由于代码已经实现，本变更的重点是：
1. 验证现有实现的正确性
2. 确保文档正确生成
3. 补充必要的单元测试（如缺失）

代码位置：
- 数据源文件: `tencentcloud/services/teo/data_source_tc_teo_origin_acl.go`
- 云 API 包: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`
