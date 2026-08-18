## 1. Schema 定义与资源骨架

- [x] 1.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.go`，定义 `ResourceTencentCloudTeoDnsRecordV2()`，Schema 包含 `zone_id`(Required,ForceNew)、`name`(Required)、`type`(Required)、`content`(Required)、`location`(Optional,Computed)、`ttl`(Optional,Computed)、`weight`(Optional,Computed)、`priority`(Optional,Computed)、`record_id`(Computed)
- [x] 1.2 为资源注册 `Importer`（`schema.ImportStatePassthrough`），import 格式为 `{zoneId}#{recordId}`

## 2. CRUD 函数实现

- [x] 2.1 实现 `resourceTencentCloudTeoDnsRecordV2Create`：构造 `CreateDnsRecordRequest`，填充 `ZoneId`/`Name`/`Type`/`Content`/`Location`/`TTL`/`Weight`/`Priority`，使用 `resource.Retry(tccommon.WriteRetryTimeout)` 调用 `CreateDnsRecordWithContext`
- [x] 2.2 Create 校验响应非空与 `RecordId` 非空，随后 `d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))` 并调用 Read 回写状态
- [x] 2.3 实现 `resourceTencentCloudTeoDnsRecordV2Read`：解析复合 ID，复用 `TeoService.DescribeTeoDnsRecordById(ctx, zoneId, recordId)`（无需新增 service 方法）；返回 nil 时打印日志并 `d.SetId("")`；否则逐字段判空后 `d.Set`
- [x] 2.4 实现 `resourceTencentCloudTeoDnsRecordV2Update`：`immutableArgs := []string{"name","type","content","location","ttl","weight","priority"}`，任一字段 `d.HasChange` 则返回错误；否则调用 Read
- [x] 2.5 实现 `resourceTencentCloudTeoDnsRecordV2Delete`：构造 `DeleteDnsRecordsRequest`，填充 `ZoneId` 与 `RecordIds=[recordId]`，使用 `resource.Retry(tccommon.WriteRetryTimeout)` 调用 `DeleteDnsRecordsWithContext`

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_teo_dns_record_v2": ResourceTencentCloudTeoDnsRecordV2()`
- [x] 3.2 在 `tencentcloud/provider.md` 中同步添加资源条目

## 4. 单元测试

- [x] 4.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2_test.go`，使用 gomonkey mock `CreateDnsRecord`、`DescribeDnsRecords`、`DeleteDnsRecords`
- [x] 4.2 覆盖 Create 成功、Create 空响应、Read 成功、Read 未找到、Delete 成功、Update 不可变字段变更报错、Update 无变更等分支

## 5. 文档

- [x] 5.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.md`，包含一句话描述（带 TEO）、Example Usage、Import 示例（`{zoneId}#{recordId}`）
- [x] 5.2 执行 `make doc` 生成 `website/docs/` 下的文档（禁止手改 website 目录）

## 6. 验证

- [x] 6.1 运行 `go build ./...` 确保编译通过
- [x] 6.2 运行 `go vet ./tencentcloud/services/teo/...`
- [x] 6.3 运行 `go test ./tencentcloud/services/teo/ -run TestAccTeoDnsRecordV2 -v -count=1 -gcflags="all=-l"` 确认单测通过
