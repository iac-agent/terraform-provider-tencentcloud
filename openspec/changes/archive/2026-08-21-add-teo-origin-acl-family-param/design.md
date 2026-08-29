## Context

当前 TEO 数据源 `tencentcloud_teo_origin_acl` 用于查询边缘安全加速平台的源站ACL信息。该数据源调用云API `DescribeOriginACL` 接口获取数据。

经过代码审查发现，`origin_acl_family` 参数实际上已经在现有代码中实现：
- Schema 定义：在 `data_source_tc_teo_origin_acl.go` 第 309-313 行已定义
- 读取逻辑：在第 478-480 行已实现从 API 响应中读取并设置该参数

云API `DescribeOriginACL` 的响应结构 `OriginACLInfo` 中已包含 `OriginACLFamily` 字段（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）。

本变更提案的目的是正式记录该参数的存在，确保代码、文档和规范的一致性。

## Goals / Non-Goals

**Goals:**
- 正式记录 `origin_acl_family` 参数在 `tencentcloud_teo_origin_acl` 数据源中的实现
- 确保该参数在文档中得到正确展示（通过 `make doc` 自动生成）
- 验证参数的实现符合 Terraform 提供程序的标准模式

**Non-Goals:**
- 不修改现有参数的实现逻辑（因为已经实现）
- 不新增其他参数或功能
- 不改变数据源的现有行为

## Decisions

### 决策1：确认参数已实现，无需修改代码

**决策**: 经过代码审查，`origin_acl_family` 参数已在数据源中完整实现，无需修改 Go 代码。

**理由**:
- Schema 已正确定义（`TypeString`, `Computed`）
- Read 方法中已正确设置该参数的值
- 代码实现符合 Terraform Plugin SDK v2 的标准模式

**替代方案**:
- 如果参数未实现，需要在 schema 中新增字段并在 Read 方法中设置值

### 决策2：通过 openspec 提案规范记录该参数

**决策**: 创建 openspec 变更提案，正式记录该参数的新增。

**理由**:
- 提供变更的可追溯性
- 确保文档自动生成时包含该参数
- 符合项目的变更管理流程

## Risks / Trade-offs

**风险1：参数实际上未正确实现**
- **可能性**: 低（已通过代码审查确认）
- **缓解措施**: 在实施阶段（openspec apply）验证参数的功能
- **验证方法**: 检查 `make doc` 生成的文档是否包含该参数

**风险2：云API字段变更**
- **可能性**: 低（该字段已在云API稳定版本中）
- **缓解措施**: 遵循现有的云API调用模式，如果云API变更，按照项目标准流程更新

**Trade-offs**:
- 由于参数已存在，此提案主要是文档性和规范性工作，实际代码变更很小
- 这确保了变更的可追溯性和完整性
