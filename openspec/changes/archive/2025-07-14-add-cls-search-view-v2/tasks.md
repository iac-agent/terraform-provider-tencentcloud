## Tasks

- [x] 1. 创建资源文件 `tencentcloud/services/cls/resource_tc_cls_search_view_v2.go`，实现完整的 CRUD 逻辑（Create、Read、Update、Delete），参考 `tencentcloud_igtm_strategy` 资源代码风格
  - Schema 定义：view_id(Computed)、logset_id(Required,ForceNew)、logset_region(Required,ForceNew)、view_name(Required)、view_type(Required)、topics(Required,TypeList)、view_id_prefix(Optional,ForceNew)、description(Optional)、create_time(Computed)、update_time(Computed)
  - Create: 调用 CreateSearchView，存储返回的 ViewId 作为资源 ID
  - Read: 调用 DescribeSearchViews，使用 Filter 按 viewId 过滤，映射 SearchViewInfo 到 state
  - Update: 调用 ModifySearchView，传入 ViewId 及可更新字段
  - Delete: 调用 DeleteSearchView，传入 ViewId
  - API 调用均需添加 retry 处理，使用 tccommon.ReadRetryTimeout
  - Topics 嵌套结构包含 region、logset_id、topic_id 三个字段

- [x] 2. 创建单元测试文件 `tencentcloud/services/cls/resource_tc_cls_search_view_v2_test.go`，使用 gomonkey mock 云 API 接口，覆盖 CRUD 操作

- [x] 3. 创建资源文档文件 `tencentcloud/services/cls/resource_tc_cls_search_view_v2.md`，包含一句话描述、Example Usage 和 Import 部分

- [x] 4. 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cls_search_view_v2` 资源

- [x] 5. 在 `tencentcloud/provider.md` 中添加 `tencentcloud_cls_search_view_v2` 资源文档条目

- [x] 6. 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
