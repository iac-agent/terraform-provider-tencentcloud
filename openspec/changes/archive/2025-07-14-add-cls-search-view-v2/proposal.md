## Why

CLS (Cloud Log Service) 当前缺少对查询视图（Search View）资源的 Terraform 管理支持。查询视图是 CLS 中用于组织和管理日志查询的重要功能，用户需要通过 Terraform 来自动化创建、修改和删除查询视图，实现基础设施即代码的管理方式。

## What Changes

- 新增 Terraform 资源 `tencentcloud_cls_search_view_v2`，支持查询视图的完整 CRUD 生命周期管理
  - Create: 调用 `CreateSearchView` 接口创建查询视图，支持配置日志集、视图名称、视图类型、主题列表、描述信息及视图ID前缀
  - Read: 调用 `DescribeSearchViews` 接口，通过 viewId 过滤读取查询视图详情
  - Update: 调用 `ModifySearchView` 接口修改查询视图的名称、类型、主题列表和描述信息
  - Delete: 调用 `DeleteSearchView` 接口删除查询视图
- 在 `provider.go` 和 `provider.md` 中注册新资源

## Capabilities

### New Capabilities
- `cls-search-view-v2`: 新增 CLS 查询视图资源的 Terraform CRUD 管理，包含创建、读取、更新、删除操作及相关的 schema 定义、单元测试和文档

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/cls/resource_tc_cls_search_view_v2.go`、`tencentcloud/services/cls/resource_tc_cls_search_view_v2_test.go`、`tencentcloud/services/cls/resource_tc_cls_search_view_v2.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（文档更新）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016` 中的 `CreateSearchView`、`DescribeSearchViews`、`ModifySearchView`、`DeleteSearchView` 接口
