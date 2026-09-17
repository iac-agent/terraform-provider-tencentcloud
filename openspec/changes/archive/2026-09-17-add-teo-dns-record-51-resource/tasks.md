## 1. Service layer

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecord51ById(ctx context.Context, zoneId, recordId string) (ret *teov20220901.DnsRecord, errRet error)`：构造 `DescribeDnsRecordsRequest`，设置 `ZoneId`、`Filters = []*teov20220901.AdvancedFilter{{Name: "id", Values: [recordId]}}`、`Limit = helper.IntInt64(1000)`（云 API 注释标注的最大值）、`Offset = 0`；在 `resource.Retry(tccommon.ReadRetryTimeout)` 内调用 `me.client.UseTeoClient().DescribeDnsRecords(request)`（含 `ratelimit.Check`），错误用 `tccommon.RetryError(e)` 包装；命中 `len(response.Response.DnsRecords) > 0` 时返回首个元素；保留既有 `[CRITAL]`/`[DEBUG]` 日志风格。不改动既有 `DescribeTeoDnsRecordById`（向后兼容）。

## 2. Resource implementation

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_51.go`，定义 `ResourceTencentCloudTeoDnsRecord51() *schema.Resource`：Create/Read/Update/Delete 四个回调 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`；schema 字段按序：`zone_id`(Required,ForceNew,TypeString)、`name`(Required,TypeString)、`type`(Required,TypeString，Description 列出 A/AAAA/MX/CNAME/TXT/NS/CAA/SRV 取值)、`content`(Required,TypeString)、`location`(Optional,Computed,TypeString)、`ttl`(Optional,Computed,TypeInt，范围 60~86400 默认 300)、`weight`(Optional,Computed,TypeInt，范围 -1~100 默认 -1)、`priority`(Optional,Computed,TypeInt，MX 专用 0~50 默认 0)、`record_id`(Computed,TypeString)、`status`(Computed,TypeString enable/disable)、`created_on`(Computed,TypeString)、`modified_on`(Computed,TypeString)。代码风格严格参考 `resource_tc_igtm_strategy.go`；文件开头不添加注释。
- [x] 2.2 实现 `resourceTencentCloudTeoDnsRecord51Create`：`defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_51.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`；`ctx := tccommon.NewResourceLifeCycleHandleFuncContext(...)`；从 schema 读取 zone_id/name/type/content（GetOk）与 location（GetOk）/ttl/weight/priority（GetOkExists）填充 `teov20220901.NewCreateDnsRecordRequest()`；在 `resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `CreateDnsRecordWithContext`，失败用 `tccommon.RetryError(e)` 包装，成功打印 `[DEBUG]` 日志（logId、action、请求/响应体）；块内校验 `result == nil || result.Response == nil || result.Response.RecordId == nil || *result.Response.RecordId == ""` 任一命中返回 `resource.NonRetryableError`（不写入空 id）；retry 成功后在块外执行 `d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))`（设置 id 等成功操作放在 retry 块外），再调用 Read。
- [x] 2.3 实现 `resourceTencentCloudTeoDnsRecord51Read`：解析 `d.Id()`（`strings.Split(d.Id(), tccommon.FILED_SP)`，len != 2 返回 `fmt.Errorf("id is broken,%s", d.Id())`）；调用 `service.DescribeTeoDnsRecord51ById(ctx, zoneId, recordId)`；若 `respData == nil`，先 `log.Printf("[CRUD] teo_dns_record_51 id=%s", d.Id())` 保留现场再 `d.SetId("")` 返回 nil；否则逐字段判 nil 后 `d.Set`：zone_id/name/type/location/content/ttl/weight/priority/status/record_id/created_on/modified_on。
- [x] 2.4 实现 `resourceTencentCloudTeoDnsRecord51Update`：解析复合 id 得到 zoneId/recordId；`mutableArgs := []string{"name", "type", "content", "location", "ttl", "weight", "priority"}`，任一 `d.HasChange` 才构造 `teov20220901.NewModifyDnsRecordsRequest()`：顶层 `request.ZoneId = helper.String(zoneId)`；`DnsRecords` 单元素 `&teov20220901.DnsRecord{RecordId: helper.String(recordId), ...可变字段}`，**不设置** DnsRecord 元素的 ZoneId/Status/CreatedOn/ModifiedOn（云 API 注释标注这些字段在 ModifyDnsRecords 中仅做出参、传入被忽略）；在 `resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `ModifyDnsRecordsWithContext`，`tccommon.RetryError(e)` 包装 + `[DEBUG]` 日志 + nil Response 防御；无字段变化则跳过 API 调用；最后调用 Read 回刷。
- [x] 2.5 实现 `resourceTencentCloudTeoDnsRecord51Delete`：解析复合 id；构造 `teov20220901.NewDeleteDnsRecordsRequest()`，设置 `ZoneId` 与 `RecordIds = helper.Strings([]string{recordId})`；`resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `DeleteDnsRecordsWithContext`，`tccommon.RetryError(e)` 包装 + `[DEBUG]` 日志 + `result == nil || result.Response == nil` 防御返回 NonRetryableError；成功后返回 nil。
- [x] 2.6 检查所有函数返回的 error 均被处理（`_ = d.Set(...)` 风格），资源命名统一使用 `teo_dns_record_51`（小写蛇形）出现在日志/错误信息中，不使用"该资源"等模糊措辞；确认无 `_extension.go` 文件生成。

## 3. Provider registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中新增 `"tencentcloud_teo_dns_record_51": teo.ResourceTencentCloudTeoDnsRecord51()`，放置于既有 `tencentcloud_teo_dns_record` 注册项附近保持 teo 命名空间连续。
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO Resource 分类下按字母序新增 `tencentcloud_teo_dns_record_51` 行（参考既有 `tencentcloud_teo_dns_record` 的格式）。

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_51.md`：一句话描述 "Provides a resource to create a teo dns_record_51"（含云产品名称）；`Example Usage` HCL 示例覆盖 zone_id/name/type/content 及 location/ttl/weight/priority 可选字段；`Import` 部分说明使用联合 id（`terraform import tencentcloud_teo_dns_record_51.xxx {zoneId}#{recordId}`）；不添加 `Argument Reference`/`Attribute Reference` 部分。

## 5. Unit tests (gomonkey)

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_51_test.go`（package teo_test），参考 `resource_tc_teo_bind_security_template_test.go` 的 mockMeta/gomonkey 写法：mock `UseTeoV20220901Client` 返回空 client，再以 `patches.ApplyMethodFunc` 分别 mock `CreateDnsRecordWithContext`、`DescribeDnsRecords`、`ModifyDnsRecordsWithContext`、`DeleteDnsRecordsWithContext`。
- [x] 5.2 编写用例：Create 成功路径（mock 返回 RecordId，断言 `d.Id()` 为 `zoneId#recordId`、state 字段正确回填）。
- [x] 5.3 编写用例：Create 空 RecordId（RecordId 为 nil 或空串）返回错误且不写入 id（NonRetryableError 路径）。
- [x] 5.4 编写用例：Read 命中（DescribeDnsRecords 返回含目标记录的列表，断言各字段 set）与未命中（返回空列表 → SetId("")）路径。
- [x] 5.5 编写用例：Update 触发 Modify（修改 content 后断言 ModifyDnsRecords 请求体单元素包含 RecordId 与新值、且不含 ZoneId/Status 等出参字段）与无变化不触发 API 的路径。
- [x] 5.6 编写用例：Delete 成功路径（断言 DeleteDnsRecords 请求 ZoneId/RecordIds 正确）。
- [x] 5.7 编写用例：复合 id 解析错误分支（id 缺段时 Read/Update/Delete 返回 "id is broken" 错误）。保证全部测试代码可编译（不执行 go test）。

## 6. Verification (handled by later phases, not in this openspec phase)

- [x] 6.1 代码正确性检查：核对 CRUD 入参与 vendor 云 API 入参一一对应（Create↔CreateDnsRecordRequest、Update↔ModifyDnsRecordsRequest 的可变字段、Delete↔DeleteDnsRecordsRequest、Read↔DescribeDnsRecordsRequest），确认编译无误（由后续验证流程执行 go build）。
- [ ] 6.2 gofmt 格式化、`make doc` 生成 website 文档、`.changelog` 文件创建：统一由收尾阶段 tfpacer-finalize skill 执行（本阶段禁止执行）。
