# Tasks: add-config-list-aggregators-datasource

## 1. 服务层查询方法

- [x] 1.1 在 `tencentcloud/services/config/service_tencentcloud_config.go` 中新增 `DescribeConfigListAggregators(ctx context.Context) (total uint64, items []*configv20220802.Aggregator, errRet error)`：
  - 使用 `configv20220802.NewListAggregatorsRequest()`，初始 `Limit=100, Offset=0`
  - 分页循环（参照 `DescribeConfigCompliancePacksByFilter`，L25-L96）：每次请求在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `UseConfigV20220802Client().ListAggregatorsWithContext(ctx, request)`
  - retry 块内：`ratelimit.Check`、API 调用、`result == nil || result.Response == nil` 时返回 `resource.NonRetryableError`（附 `log.Printf("[DATASOURCE] read empty, skip SetId")` 类调试日志）、DEBUG 日志；retry 块外累加 Offset、累加 items、判断 `len(items) >= total` 退出

## 2. 数据源实现

- [x] 2.1 新增 `tencentcloud/services/config/data_source_tc_config_list_aggregators.go`：
  - `DataSourceTencentCloudConfigListAggregators()`，schema 严格按 design.md 的表格：`total`、`items`（Elem 含 9 个字段）、`result_output_file`
  - Read 函数：`tccommon.LogElapsed` / `InconsistentCheck` defer、`NewResourceLifeCycleHandleFuncContext`、调用 service 层（service 内已有 retry，外层不再 retry）
  - 遍历 `Items` 逐字段判 nil 后 set；`d.SetId(helper.BuildToken())`；处理 `result_output_file`
- [x] 2.2 新增 `data_source_tc_config_list_aggregators_test.go`：gomonkey mock `ListAggregatorsWithContext`，覆盖成功查询与响应为空两个用例；不使用 terraform 测试套件，不实际执行

## 3. 注册与文档

- [x] 3.1 `tencentcloud/provider.go` 数据源 map 中注册 `"tencentcloud_config_list_aggregators": config.DataSourceTencentCloudConfigListAggregators(),`（按现有排序位置插入）
- [x] 3.2 `tencentcloud/provider.md` 添加对应注册说明行
- [x] 3.3 新增 `tencentcloud/services/config/data_source_tc_config_list_aggregators.md`：
  - 一句话描述（带 Config 产品名，"Use this data source to query ..."）
  - Example Usage（引用 total/items/result_output_file）
  - 无 Import 部分；不添加 Argument/Attribute Reference

## 4. 收尾（由后续流程执行）

- [ ] 4.1 收尾阶段执行 `gofmt`、`make doc`、生成 `.changelog` 文件并推送（tfpacer-finalize skill）
