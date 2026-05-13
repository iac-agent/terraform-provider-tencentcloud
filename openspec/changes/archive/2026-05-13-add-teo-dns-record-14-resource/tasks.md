## 1. Schema 与 CRUD 函数实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_14.go`，定义 ResourceTencentCloudTeoDnsRecord14() 函数和 schema
- [x] 1.2 实现 resourceTencentCloudTeoDnsRecord14Create 函数，调用 CreateDnsRecord API，设置复合 ID (zoneId#recordId)
- [x] 1.3 实现 resourceTencentCloudTeoDnsRecord14Read 函数，调用 DescribeDnsRecords API，通过 Filters 按 id 过滤查询记录
- [x] 1.4 实现 resourceTencentCloudTeoDnsRecord14Update 函数，调用 ModifyDnsRecords API，传入 DnsRecord 对象更新可变字段
- [x] 1.5 实现 resourceTencentCloudTeoDnsRecord14Delete 函数，调用 DeleteDnsRecords API，传入 RecordIds 删除记录

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_dns_record_14` 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_14` 资源条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_14.md` 文档，包含一句话描述、Example Usage 和 Import 部分

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_14_test.go`，使用 gomonkey mock 云 API 编写单元测试
- [x] 4.2 运行 `go test -gcflags=all=-l` 确保单元测试通过
