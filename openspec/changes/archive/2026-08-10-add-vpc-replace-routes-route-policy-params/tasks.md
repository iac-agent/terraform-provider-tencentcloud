## 1. Service 层扩展

- [x] 1.1 在 `tencentcloud/services/vpc/service_tencentcloud_vpc.go` 的 `DescribeRouteTables` 方法签名末尾新增 `needRouterInfo *bool`、`name string`、`values []*string` 三个参数，保持现有调用方传零值兼容。
- [x] 1.2 在 `DescribeRouteTables` 方法实现中：当 `needRouterInfo != nil` 时设置 `request.NeedRouterInfo`；当 `name != ""` 且 `len(values) > 0` 时构造 `Filter{Name: name, Values: values}` 追加到 `request.Filters`。
- [x] 1.3 调整方法返回值，将首次响应的 `TotalCount` 一并返回（或在方法内通过返回结构携带），供资源 Read 写入 `total_count`。若不便改返回值，则在 service 层新增一个返回 `totalCount` 的封装或扩展返回结构。

## 2. 资源 Schema 与 Read 逻辑

- [x] 2.1 在 `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config.go` 的 Schema 中新增 `need_router_info`（Optional bool）、`name`（Optional string）、`values`（Optional TypeList, Elem TypeString）、`total_count`（Computed int）四个字段。
- [x] 2.2 在 `resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigRead` 中读取 `need_router_info`、`name`、`values` 并透传给 `DescribeRouteTables` service 方法；当 `name` 非空且 `values` 非空时组装 filter。
- [x] 2.3 在 Read 方法中，当 `DescribeRouteTables` 返回空时，先 `log.Printf("[CRUD] ...")` 保留 id 现场，再 `d.SetId("")`；当返回非空时，`_ = d.Set("total_count", totalCount)` 写回出参，并保持现有 `route_table_id` 的 set 逻辑。
- [x] 2.4 确认 Create/Update/Delete 方法无需改动（新增参数仅影响 Read），Update 方法中 `route_table_id` 仍通过 `d.Id()` 获取。

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config_test.go` 中使用 gomonkey mock 云API（`DescribeRouteTables`、`ReplaceRoutesWithRoutePolicy`），补充覆盖 `need_router_info`、`name`、`values`、`total_count` 的单元测试用例，验证参数透传与 state 回填逻辑。
- [x] 3.2 补充 Read 返回空响应场景的测试用例，验证 `d.SetId("")` 行为及日志保留。

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config.md` 的 Example Usage，补充 `need_router_info`、`name`、`values` 参数示例（`total_count` 为 Computed，仅在出参说明，无需在 example 中显式设置）。
- [x] 4.2 确认 `tencentcloud/provider.go` 中资源已注册（无需新增注册），`website/docs/` 文档由收尾阶段 `make doc` 自动生成。
