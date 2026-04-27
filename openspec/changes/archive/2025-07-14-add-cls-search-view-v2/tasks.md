## 1. Schema 定义与 CRUD 函数实现

- [x] 1.1 创建 `tencentcloud/services/cls/resource_tc_cls_search_view_v2.go`，定义 ResourceTencentCloudClsSearchViewV2() 资源函数，包含完整 Schema 定义（logset_id, logset_region, view_name, view_type, topics, description, view_id_prefix, view_id）
- [x] 1.2 实现 resourceTencentCloudClsSearchViewV2Create 函数，调用 CreateSearchView API，映射所有入参，设置 resource ID 为返回的 ViewId
- [x] 1.3 实现 resourceTencentCloudClsSearchViewV2Read 函数，调用 DescribeSearchViews API，通过 Filter 按 viewId 过滤，从 SearchViewInfo 中填充 state
- [x] 1.4 实现 resourceTencentCloudClsSearchViewV2Update 函数，调用 ModifySearchView API，传入 view_id, view_name, view_type, topics, description
- [x] 1.5 实现 resourceTencentCloudClsSearchViewV2Delete 函数，调用 DeleteSearchView API，传入 view_id
- [x] 1.6 添加 Importer 支持（schema.ImportStatePassthrough）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 tencentcloud_cls_search_view_v2 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 tencentcloud_cls_search_view_v2 资源条目

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/cls/resource_tc_cls_search_view_v2_test.go`，使用 gomonkey mock 方式编写 Create/Read/Update/Delete 的单元测试
- [x] 3.2 使用 `go test -gcflags=all=-l` 运行单元测试并确保通过

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/cls/resource_tc_cls_search_view_v2.md`，包含一句话描述、Example Usage 和 Import 部分
