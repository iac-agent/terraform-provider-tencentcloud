## 1. 资源 Schema 与 CRUD 实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.go` 文件，定义资源 schema（zone_id、name、type、content、location、ttl、weight、priority、record_id、status、created_on、modified_on），参考 `tencentcloud_igtm_strategy` 代码风格
- [x] 1.2 实现 Create 函数：调用 `CreateDnsRecord` API，设置复合 ID（zone_id#record_id），使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 重试，调用 Read 回填
- [x] 1.3 实现 Read 函数：从 `d.Get()` 获取 zone_id 和 record_id，调用 `DescribeDnsRecords` API 使用 AdvancedFilter 按 id 过滤查询，回填所有 schema 字段
- [x] 1.4 实现 Update 函数：检测 name、type、content、location、ttl、weight、priority 变更，调用 `ModifyDnsRecords` API，只传 RecordId 和变更字段
- [x] 1.5 实现 Delete 函数：调用 `DeleteDnsRecords` API，传入 ZoneId 和 RecordIds
- [x] 1.6 支持 Import：在 schema 中配置 `Importer: schema.ImportStatePassthrough`

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecordV2ById` 方法，使用 AdvancedFilter 按 record id 查询 DNS 记录

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_dns_record_v2` 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_v2` 资源条目

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.md` 文件，包含资源描述、Example Usage、Import 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2_test.go` 文件，使用 gomonkey mock 云 API，测试 Create、Read、Update、Delete 业务逻辑
- [x] 5.2 运行 `go test` 验证单元测试通过
