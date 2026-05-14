## 1. Schema 与 CRUD 函数实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_18.go` 文件，定义 `ResourceTencentCloudTeoDnsRecord18()` 函数，包含完整的 Schema 定义（zone_id Required+ForceNew, name/type/content Required, location/ttl/weight/priority Optional+Computed, record_id/status/created_on/modified_on Computed）和 Importer 支持
- [x] 1.2 实现 `resourceTencentCloudTeoDnsRecord18Create()` 函数，调用 CreateDnsRecord API，设置复合 ID（zone_id#record_id），处理空响应返回 NonRetryableError
- [x] 1.3 实现 `resourceTencentCloudTeoDnsRecord18Read()` 函数，解析复合 ID，调用 TeoService.DescribeTeoDnsRecordById()，设置所有非 nil 字段到 state，处理资源不存在情况
- [x] 1.4 实现 `resourceTencentCloudTeoDnsRecord18Update()` 函数，检测 name/type/content/location/ttl/weight/priority 字段变化，调用 ModifyDnsRecords API 构造 DnsRecord 子结构体
- [x] 1.5 实现 `resourceTencentCloudTeoDnsRecord18Delete()` 函数，调用 DeleteDnsRecords API，传入 ZoneId 和 RecordIds

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_dns_record_18` 资源，添加 factory 函数 `ResourceTencentCloudTeoDnsRecord18`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_18` 资源条目

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_18_test.go` 文件，使用 gomonkey mock 方式编写 Create、Read、Update、Delete 操作的单元测试
- [x] 3.2 运行 `go test -gcflags=all=-l` 验证单元测试通过

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_18.md` 文件，包含一句话描述（提及 TEO）、Example Usage、Import 部分（说明使用复合 ID 格式）
