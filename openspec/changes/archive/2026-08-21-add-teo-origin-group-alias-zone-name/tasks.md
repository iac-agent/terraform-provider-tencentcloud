## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_origin_group.go` 的 `ResourceTencentCloudTeoOriginGroup()` 中，为 `references` 嵌套块新增 `alias_zone_name` 字段（`Type: schema.TypeString`、`Computed: true`），Description 说明其为被引用实例的别名站点名称
- [x] 1.2 确认 `references` 块仍是 computed-only，`alias_zone_name` 不加入 Create/Update/Delete 的请求构建逻辑

## 2. Read 逻辑

- [x] 2.1 在 `resourceTencentCloudTeoOriginGroupRead()` 的 `references` 循环中，新增对 `references.AliasZoneName` 的 nil-check，并在非 nil 时设置 `referencesMap["alias_zone_name"]`
- [x] 2.2 确认 `service.DescribeTeoOriginGroupById` 已返回完整 `OriginGroupReference`（含 `AliasZoneName`），无需修改 service 层

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go` 中使用 gomonkey mock `DescribeOriginGroup`（或 `DescribeTeoOriginGroupById`），构造返回包含 `AliasZoneName` 的响应
- [x] 3.2 新增测试用例验证 Read 后 state 中 `references` 块条目正确包含 `alias_zone_name`
- [x] 3.3 新增测试用例验证当 `AliasZoneName` 为 nil 时，state 中不设置 `alias_zone_name`（nil-skip 行为）

## 4. 文档同步

- [x] 4.1 检查 `tencentcloud/services/teo/resource_tc_teo_origin_group.md`，确认 `references` 为 computed 块、示例无需新增 computed 字段；如描述有误则同步修正
- [x] 4.2 说明文档由 `make doc` 自动生成 `website/docs/`，禁止手动修改 website 目录文件

## 5. 验证

- [x] 5.1 运行 `gofmt` 格式化修改的 Go 文件
- [x] 5.2 确认新增字段为纯 computed、不破坏现有 TF 配置与 state（向后兼容检查）
