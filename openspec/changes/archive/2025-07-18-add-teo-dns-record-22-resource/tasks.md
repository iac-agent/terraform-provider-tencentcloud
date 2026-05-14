## 1. 资源 Schema 与 CRUD 实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_22.go`，定义 `ResourceTencentCloudTeoDnsRecord22` 函数，包含完整的 Schema 定义（zone_id、name、type、content、location、ttl、weight、priority、status、created_on、modified_on、record_id）
- [x] 1.2 实现 `resourceTencentCloudTeoDnsRecord22Create` 函数，调用 CreateDnsRecord API，设置联合 ID（zone_id#record_id）
- [x] 1.3 实现 `resourceTencentCloudTeoDnsRecord22Read` 函数，调用 DescribeDnsRecords API，通过 RecordId 过滤获取单条记录，设置所有 schema 字段
- [x] 1.4 实现 `resourceTencentCloudTeoDnsRecord22Update` 函数，分别调用 ModifyDnsRecords（修改 name/type/content/location/ttl/weight/priority）和 ModifyDnsRecordsStatus（修改 status）
- [x] 1.5 实现 `resourceTencentCloudTeoDnsRecord22Delete` 函数，调用 DeleteDnsRecords API

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DescribeTeoDnsRecord22ById` 方法，封装 DescribeDnsRecords 调用并通过 RecordId 过滤获取单条记录

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_dns_record_22` 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_22` 资源条目

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_22.md`，包含一句话描述、Example Usage 和 Import 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_22_test.go`，使用 gomonkey mock 方式编写 Create、Read、Update、Delete 的单元测试
- [x] 5.2 运行 `go test -gcflags=all=-l` 确保所有单元测试通过
