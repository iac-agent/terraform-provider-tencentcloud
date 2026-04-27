## Context

Terraform Provider for TencentCloud 需要新增 CLS（日志服务）查询视图（SearchView）的通用资源管理能力。CLS 云产品已提供完整的查询视图 CRUD API（CreateSearchView / DescribeSearchViews / ModifySearchView / DeleteSearchView），需要将其封装为 Terraform RESOURCE_KIND_GENERAL 资源，使用户能够通过 Terraform 声明式管理查询视图的完整生命周期。

当前 CLS 服务目录下已有多个资源（alarm、logset、topic、dashboard 等），新资源需遵循现有代码组织结构和编码模式。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_cls_search_view_v2` 资源的完整 CRUD（Create/Read/Update/Delete）
- 支持通过 DescribeSearchViews 按 ViewId 过滤读取资源状态
- 支持资源 Import
- 在 provider.go 中正确注册资源
- 编写使用 gomonkey mock 的单元测试
- 生成 .md 文档用于 make doc

**Non-Goals:**
- 不实现查询视图的数据源（data source），仅实现资源
- 不修改已有的 CLS 资源代码
- 不处理异步接口轮询（所有接口均为同步）

## Decisions

### 1. 资源 ID 使用 ViewId
- **决策**: 使用 CreateSearchView 返回的 ViewId 作为 Terraform 资源 ID
- **理由**: ViewId 是云 API 返回的唯一标识，且在 Delete/Modify/Describe 中均作为关键入参，无需复合 ID

### 2. Read 操作使用 DescribeSearchViews + Filter
- **决策**: 在 Read 中使用 DescribeSearchViews 接口，通过 Filters 按 viewId 过滤获取资源详情
- **理由**: 云 API 未提供 DescribeSearchView（单条）接口，DescribeSearchViews 支持按 viewId 过滤，可精确获取单条记录
- **替代方案**: 无其他可用接口

### 3. Topics 字段使用 TypeList + schema.Resource
- **决策**: Topics 字段使用 TypeList 嵌套 schema.Resource，包含 region、logset_id、topic_id 三个子字段
- **理由**: ViewSearchTopic 结构体包含三个属性，需要作为嵌套块在 Terraform 中配置

### 4. LogsetId 和 LogsetRegion 为 ForceNew
- **决策**: logset_id 和 logset_region 设置为 ForceNew：true，因为 ModifySearchView 接口不支持修改这两个字段
- **理由**: 创建后无法修改日志集归属，修改时需要销毁重建

### 5. view_id_prefix 为 Optional + ForceNew
- **决策**: view_id_prefix 设置为 Optional 且 ForceNew：true
- **理由**: 该参数仅在创建时使用，ModifySearchView 接口不包含此参数

### 6. view_id 为 Computed
- **决策**: view_id 设置为 Computed，从 CreateSearchView 响应中获取
- **理由**: view_id 由云 API 生成返回，格式为 ${ViewIdPrefix}-view

### 7. 编码风格参考 igtm_strategy
- **决策**: 代码结构和编码风格参考 tencentcloud_igtm_strategy 资源
- **理由**: 该资源是 RESOURCE_KIND_GENERAL 的标准实现范例

## Risks / Trade-offs

- [DescribeSearchViews 返回空] → 如果资源在 Terraform 之外被删除，Read 时 DescribeSearchViews 按 viewId 过滤后 Infos 为空，需要正确处理为资源已不存在（d.SetId("")）
- [Topics 嵌套块变更检测] → TypeList 嵌套 schema.Resource 的变更检测依赖 Terraform SDK 的默认行为，需确保 Update 时正确传入完整的 Topics 列表
