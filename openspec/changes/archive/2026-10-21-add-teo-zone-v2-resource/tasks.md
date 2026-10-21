## 1. 资源实现（tencentcloud/services/teo/resource_tc_teo_zone_v2.go）

- [x] 1.1 定义 `ResourceTencentCloudTeoZoneV2()` 的 schema：可选入参 `type`/`zone_name`/`area`/`plan_id`（ForceNew）/`alias_zone_name`/`tags`（list of {tag_key Required, tag_value Required}）/`allow_duplicates`（ForceNew）/`jump_start`（ForceNew）/`vanity_name_servers`（Optional+Computed，{switch Required, servers Optional}）；Computed 出参 `zone_id`、`paused`、`status`、`cname_status`、`active_status`、`lock_status`、`created_on`、`modified_on`、`name_servers`、`original_name_servers`、`resources`、`ns_detail`、`cname_detail`、`dns_pod_detail`、`ownership_verification`、`vanity_name_servers_ips`、`work_mode_infos`（嵌套结构按 design.md D1 的字段清单展开，`IPv4` 对应 schema 名 `i_pv4`）；注册 Importer（ImportStatePassthrough）。不生成 `_extension.go`。
- [x] 1.2 实现 `resourceTencentCloudTeoZoneV2Create`：`d.GetOk` 组装 CreateZoneRequest 全部入参（Tags 组装为 `[]*teov20220901.Tag`）；`resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `UseTeoV20220901Client().CreateZoneWithContext(ctx, request)`，失败 `tccommon.RetryError(e)`，retry 内检查 response/Response 为 nil 返回 NonRetryableError；retry 外打印 logId 与 d.Id() 排障日志，检查 `ZoneId` 为 nil/空则报错，然后 `d.SetId(zoneId)`，最后调 Read 回填。
- [x] 1.3 实现 `resourceTencentCloudTeoZoneV2Read`：`zoneId := d.Id()`；组装 DescribeZonesRequest（Filters=zone-id、`Limit = helper.Int64(100)`、Offset=0）；`resource.Retry(tccommon.ReadRetryTimeout)` 内调用 `DescribeZonesWithContext`；retry 外检查空响应时先 `log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())` 再 `d.SetId("")` 返回 nil；取 Zones[0] 逐字段判 nil 后 `_ = d.Set(...)`，嵌套结构（resources/sv、ns_detail、cname_detail、dns_pod_detail、ownership_verification 的 dns/file/ns verification、vanity_name_servers、vanity_name_servers_ips、work_mode_infos、tags）按 igtm_strategy 风格组装 `[]map[string]interface{}`。
- [x] 1.4 实现 `resourceTencentCloudTeoZoneV2Update`：mutableArgs 为 `["type", "zone_name", "area", "alias_zone_name", "vanity_name_servers"]`；命中变更则组装 ModifyZoneRequest（`request.ZoneId = &zoneId` 从 d.Id() 取值，vanity_name_servers 组装为 `teov20220901.VanityNameServers`）；`resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `ModifyZoneWithContext`，retry 内检查 response nil 返回 NonRetryableError；完成后调 Read。
- [x] 1.5 实现 `resourceTencentCloudTeoZoneV2Delete`：`zoneId := d.Id()`；DeleteZoneRequest.ZoneId；`resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `DeleteZoneWithContext`，retry 内检查 response nil；同步接口无需轮询，成功后返回 nil。
- [x] 1.6 代码自查：所有函数返回的 error 均已处理（必定不出错的 `_ = d.Set(...)`）；日志/错误信息统一使用 `teo_zone_v2`；云API调用均以 tccommon.ReadRetryTimeout/WriteRetryTimeout 做超时 retry 且 retry 块内只有接口调用；不修改旧版 `resource_tc_teo_zone.go` 与其 `_extension.go`。

## 2. 单元测试（tencentcloud/services/teo/resource_tc_teo_zone_v2_test.go，gomonkey 方式）

- [x] 2.1 搭建测试骨架：`package teo_test`，`mockMetaTeoZoneV2` 实现 `tccommon.ProviderMeta`，`patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)`，`schema.TestResourceDataRaw` 构造 ResourceData（参考 resource_tc_teo_edge_kv_namespace_test.go）。
- [x] 2.2 编写 Create 用例：mock `CreateZoneWithContext` 断言入参映射（含 tags/plan_id/allow_duplicates/jump_start），mock `DescribeZonesWithContext` 返回带 ownership_verification/resources/ns_detail 的 Zone，断言 `d.Id()` 与 state 回填；另加 Create 响应 ZoneId 为空时报错的用例。
- [x] 2.3 编写 Read 用例：正常回填（含嵌套结构）、NotFound（Zones 为空）清空 id、Response 为 nil 清空 id。
- [x] 2.4 编写 Update 用例：mock `ModifyZoneWithContext` 断言 type/area/alias_zone_name/vanity_name_servers 入参映射与 ZoneId 来源。
- [x] 2.5 编写 Delete 用例：mock `DeleteZoneWithContext` 断言 ZoneId 入参。
- [x] 2.6 确认测试代码在当前环境可正确编译（不执行 go test，仅保证构建正确性；禁止使用 terraform 测试套件）。

## 3. Provider 注册与文档

- [x] 3.1 在 `tencentcloud/provider.go` 的资源注册 map 中追加 `"tencentcloud_teo_zone_v2": teo.ResourceTencentCloudTeoZoneV2(),`（在 `tencentcloud_teo_zone` 行附近，保持对齐格式）。
- [x] 3.2 在 `tencentcloud/provider.md` 的 TencentCloud EdgeOne(TEO) Resource 列表中追加 `tencentcloud_teo_zone_v2`（`tencentcloud_teo_zone` 行后）。
- [x] 3.3 新增 `tencentcloud/services/teo/resource_tc_teo_zone_v2.md`：一句话描述（"Provides a resource to create a TEO zone"，带上云产品名称）；Example Usage 完整 HCL（含 tags 块、vanity_name_servers 块）；Import 部分（说明使用 zone id 导入，单 id 无需联合 id）；不添加 Argument Reference / Attribute Reference 部分。

## 4. 验证（不执行 go build/vet/lint/test）

- [x] 4.1 人工核对参数映射：Create 参数在 CreateZoneRequest 中存在、Update 参数在 ModifyZoneRequest 中存在、Delete 参数在 DeleteZoneRequest 中存在、Read 出参字段在 Zone/子结构体中存在（对照 design.md 附录 A）。
- [x] 4.2 人工核对文件清单：新增 `resource_tc_teo_zone_v2.go`、`resource_tc_teo_zone_v2_test.go`、`resource_tc_teo_zone_v2.md`；修改 `provider.go`、`provider.md`；未触碰 website/ 目录（website 文档由收尾阶段 make doc 生成）；未新增 .changelog 文件（由 tfpacer-finalize 统一处理）。
