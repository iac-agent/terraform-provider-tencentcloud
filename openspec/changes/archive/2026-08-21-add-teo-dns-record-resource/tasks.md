## 1. Schema 定义

- [x] 1.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.go`，定义 `ResourceTencentCloudTeoDnsRecordV3()` 返回的 `schema.Resource`
- [x] 1.2 定义 Required 字段：`zone_id`（`ForceNew: true`）、`name`、`type`、`content`
- [x] 1.3 定义 Optional+Computed 字段：`location`、`ttl`、`weight`、`priority`（与 SDK 默认值语义一致）
- [x] 1.4 定义纯 Computed 字段：`record_id`、`status`、`created_on`、`modified_on`（不设 `Optional`）
- [x] 1.5 配置 `Importer`（`ImportStatePassthrough`，使用复合 ID）

## 2. CRUD 函数实现

- [x] 2.1 Create：构造 `CreateDnsRecordRequest`（`ZoneId`/`Name`/`Type`/`Content` + 可选 `Location`/`TTL`/`Weight`/`Priority`），在 `resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `CreateDnsRecordWithContext`；成功后检查 `Response`/`Response.RecordId` 非空（为空返回 `NonRetryableError`），设置 `d.SetId(zoneId#recordId)`（`tccommon.FILED_SP` 分隔）
- [x] 2.2 Read：拆分 `d.Id()` 得到 `zoneId`、`recordId`，复用 `service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)`；返回 nil 时先 `log.Printf` 打印资源 id 再 `d.SetId("")`；对 `DnsRecord` 各字段 nil-safe 调用 `d.Set`
- [x] 2.3 Update：以 `name`/`type`/`content`/`location`/`ttl`/`weight`/`priority` 为 `mutableArgs` 检测变更；有变更时构造 `ModifyDnsRecordsRequest`（`ZoneId` + 单元素 `DnsRecords`，元素仅设置 `RecordId` 与发生变化的可变字段，**不**设置 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn`），在 `resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `ModifyDnsRecordsWithContext`
- [x] 2.4 Delete：拆分 `d.Id()` 得到 `zoneId`、`recordId`，构造 `DeleteDnsRecordsRequest`（`ZoneId` + `RecordIds=[recordId]`），在 `resource.Retry(tccommon.WriteRetryTimeout)` 内调用 `DeleteDnsRecordsWithContext`
- [x] 2.5 确认 service 层直接复用已有 `TeoService.DescribeTeoDnsRecordById`，无需新增 helper

## 3. 资源注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_teo_dns_record_v3": teo.ResourceTencentCloudTeoDnsRecordV3()`
- [x] 3.2 在 `tencentcloud/provider.md` 中新增对应资源条目

## 4. 单元测试

- [x] 4.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v3_test.go`，使用 gomonkey mock `CreateDnsRecordWithContext`、`DescribeDnsRecords`、`ModifyDnsRecordsWithContext`、`DeleteDnsRecordsWithContext`
- [x] 4.2 覆盖 Create 成功分支与 `RecordId` 为 nil 的失败分支
- [x] 4.3 覆盖 Read 成功分支与空结果（`d.SetId("")`）分支
- [x] 4.4 覆盖 Update 存在可变字段变更与无可变字段变更两个分支
- [x] 4.5 覆盖 Delete 分支

## 5. 文档

- [x] 5.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.md`，包含一句话描述（带 TEO）、Example Usage、Import 示例（说明复合 ID `zone_id#record_id`）
- [x] 5.2 文档中不手写 `Argument Reference` / `Attribute Reference`（由工具自动生成）；website docs 由收尾阶段 `make doc` 生成

## 6. 正确性检查（验证）

- [x] 6.1 检查新增代码引用的 SDK 类型与函数签名与 `vendor` 下 `teov20220901` 完全一致（确保可编译）
- [x] 6.2 检查所有函数返回的 `error` 均已处理（必然不出错的用 `_ =` 承接），无未使用变量
- [x] 6.3 build/lint/test 等验证由其他流程执行；`gofmt`/`make doc`/changelog 由收尾阶段 `tfpacer-finalize` skill 执行
