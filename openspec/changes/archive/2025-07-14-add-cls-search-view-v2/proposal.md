## Why

CLS（日志服务）新增了查询视图（SearchView）的云API接口（CreateSearchView / DescribeSearchViews / ModifySearchView / DeleteSearchView），需要为 Terraform Provider 新增对应的通用资源 `tencentcloud_cls_search_view_v2`，使用户能够通过 Terraform 管理 CLS 查询视图的完整生命周期。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_cls_search_view_v2`，支持 CLS 查询视图的创建、读取、更新、删除操作
- 资源文件: `tencentcloud/services/cls/resource_tc_cls_search_view_v2.go`
- 测试文件: `tencentcloud/services/cls/resource_tc_cls_search_view_v2_test.go`
- 文档文件: `tencentcloud/services/cls/resource_tc_cls_search_view_v2.md`
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源

## Capabilities

### New Capabilities
- `cls-search-view-v2`: CLS 查询视图资源的 CRUD 管理，包括日志集关联、视图类型配置、主题列表管理等功能

### Modified Capabilities
（无修改的已有能力）

## Impact

- 新增资源代码文件，不影响已有资源
- 需要在 `tencentcloud/provider.go` 中注册新资源
- 依赖云 API 接口: `cls` 产品的 `CreateSearchView`、`DescribeSearchViews`、`ModifySearchView`、`DeleteSearchView`
- SDK 包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`
