## 1. 资源实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 的 `DescribeTeoDnsRecordById` 方法内为 `DescribeDnsRecords` 请求补充 `Limit` 为云 API 标注最大值 1000（不改变方法签名与返回语义）
- [x] 1.2 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_49.go`，定义 `ResourceTencentCloudTeoDnsRecord49()` schema：`zone_id`（Required+ForceNew）、`name`/`type`/`content`（Required）、`location`/`ttl`/`weight`/`priority`（Optional+Computed）、`status`/`record_id`/`created_on`/`modified_on`（Computed），并声明 Importer（`ImportStatePassthrough`）
- [x] 1.3 实现 `resourceTencentCloudTeoDnsRecord49Create`：构建 `CreateDnsRecord` 请求（必填 + 可选字段，`GetOkExists` 读取 int 字段），WriteRetryTimeout retry 内调用 `CreateDnsRecordWithContext`，成功后在 retry 块外校验 `Response`/`RecordId` 非空（为空返回 NonRetryableError，校验前打印 logId），然后 `d.SetId(zoneId + FILED_SP + recordId)` 并调用 Read
- [x] 1.4 实现 `resourceTencentCloudTeoDnsRecord49Read`：解析复合 ID（两段，否则报 "id is broken"），调用服务层 `DescribeTeoDnsRecordById` 查询单条记录；查不到时先 `log.Printf("[CRUD] teo dns_record_49 id=%s", d.Id())` 再 `d.SetId("")`；查到后逐字段判空回填（含 `status`/`record_id`/`created_on`/`modified_on`）
- [x] 1.5 实现 `resourceTencentCloudTeoDnsRecord49Update`：解析复合 ID，检测 `name`/`type`/`content`/`location`/`ttl`/`weight`/`priority` 是否有变更；有变更则构建单条 `DnsRecord{RecordId, ...}`（不传 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn`）调用 `ModifyDnsRecordsWithContext`（WriteRetryTimeout retry），完成后调用 Read
- [x] 1.6 实现 `resourceTencentCloudTeoDnsRecord49Delete`：解析复合 ID，调用 `DeleteDnsRecordsWithContext`，`RecordIds` 传单条记录（WriteRetryTimeout retry）

## 2. Provider 注册与文档

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `"tencentcloud_teo_dns_record_49": teo.ResourceTencentCloudTeoDnsRecord49()`
- [x] 2.2 在 `tencentcloud/provider.md` 中补充 `tencentcloud_teo_dns_record_49` 资源条目
- [x] 2.3 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_49.md`：一句话描述（含 EdgeOne 云产品名）+ Example Usage + Import 部分（说明使用复合 ID `zone_id#record_id`），不手写 Argument/Attribute Reference

## 3. 单元测试

- [x] 3.1 新建 `tencentcloud/services/teo/resource_tc_teo_dns_record_49_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件），覆盖 Create（成功/空 RecordId/参数填充）、Read（回填/不存在/ID 损坏）、Update（可变字段触发 ModifyDnsRecords/无变更不调用）、Delete（单条 RecordIds）业务逻辑，保证代码可正确构建

## 4. 验证（收尾阶段执行）

- [x] 4.1 通过 openspec apply 阶段确认所有任务完成后，由收尾阶段 tfpacer-finalize skill 统一执行 `gofmt`、`make doc`（生成 website/docs 文档）及 `.changelog` 文件创建，本阶段不执行这些命令
