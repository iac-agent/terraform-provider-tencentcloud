## Context

CLS (Cloud Log Service) 当前不支持通过 Terraform 管理查询视图（Search View）资源。查询视图用于组织和管理日志查询，用户需要通过 Terraform 来自动化创建、修改和删除查询视图。

当前仓库中 CLS 服务已有多个资源实现（如 `tencentcloud_cls_alarm`、`tencentcloud_cls_alarm_notice` 等），新增资源需要遵循现有的代码模式和规范。

云 API 接口情况：
- `CreateSearchView`：创建查询视图，入参包括 LogsetId、LogsetRegion、ViewName、ViewType、Topics、ViewIdPrefix、Description，出参返回 ViewId
- `DescribeSearchViews`：查询视图列表，入参包括 Filters（支持 viewId、viewName、logsetId 过滤）、Offset、Limit，出参返回 Infos（SearchViewInfo 列表）和 Total
- `ModifySearchView`：修改查询视图，入参包括 ViewId、ViewName、ViewType、Topics、Description
- `DeleteSearchView`：删除查询视图，入参包括 ViewId

关键数据结构：
- `SearchViewInfo`：包含 ViewId、ViewName、ViewType、LogsetId、LogsetRegion、Topics（ViewSearchTopic 列表）、Description、CreateTime、UpdateTime
- `ViewSearchTopic`：包含 Region、LogsetId、TopicId
- `Filter`：包含 Key、Values

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_cls_search_view_v2` 资源，支持查询视图的完整 CRUD 生命周期管理
- 在 provider.go 和 provider.md 中注册新资源
- 编写单元测试验证业务逻辑
- 生成资源文档 .md 文件

**Non-Goals:**
- 不修改已有资源的 schema 或行为
- 不添加查询视图的数据源（data source）
- 不处理异步接口轮询（所有接口均为同步接口）

## Decisions

### 1. 资源 ID 设计
**决策**：使用 `view_id` 作为资源唯一标识符（由 CreateSearchView 返回的 ViewId）

**理由**：ViewId 是云 API 返回的唯一标识，且在 Delete、Modify、Describe 接口中均作为关键参数使用。无需复合 ID。

### 2. Read 接口查询策略
**决策**：调用 `DescribeSearchViews` 接口，使用 Filters 按 viewId 过滤来读取单个资源详情

**理由**：CLS 没有提供单条记录查询接口，DescribeSearchViews 支持 viewId 过滤，可以精确定位到目标资源。

### 3. Topics 字段 Schema 设计
**决策**：将 Topics 定义为 TypeList + TypeMap 的嵌套结构，每个 Topic 包含 region、logset_id、topic_id 三个字符串字段

**理由**：ViewSearchTopic 结构体包含 Region、LogsetId、TopicId 三个字段，均为字符串类型，使用 TypeMap 可以简化用户配置。

### 4. Create 和 Update 参数差异处理
**决策**：
- Create 接口参数：logset_id、logset_region、view_name、view_type、topics、view_id_prefix、description
- Update 接口参数：view_id、view_name、view_type、topics、description
- logset_id 和 logset_region 仅在 Create 时设置，不可更新（ModifySearchView 不支持修改这两个字段）
- view_id_prefix 仅在 Create 时设置，不可更新（ModifySearchView 不支持修改此字段）

**理由**：严格对齐云 API 接口参数，确保 Create 入参与 CreateSearchView 接口一致，Update 入参与 ModifySearchView 接口一致。

### 5. Computed 字段
**决策**：create_time 和 update_time 设为 Computed 字段，从 DescribeSearchViews 返回的 SearchViewInfo 中读取

**理由**：这两个字段由服务端生成和管理，用户不应设置。

### 6. 参考资源
**决策**：代码风格严格参考 `tencentcloud_igtm_strategy` 资源

**理由**：需求明确要求参考该资源的代码样式。

## Risks / Trade-offs

- **DescribeSearchViews 分页查询** → 使用 Limit=100（最大值）单次查询，若结果不足则无需分页。正常情况下一个 viewId 只匹配一条记录，不存在分页问题。
- **DescribeSearchViews 接口使用 Filters 而非直接 ID 查询** → 需要遍历返回结果匹配 viewId，增加了少量代码复杂度，但这是唯一可用的读取方式。
- **ModifySearchView 不支持修改 logset_id 和 logset_region** → 这两个字段在 Terraform schema 中设置 ForceNew，修改时触发资源重建。
