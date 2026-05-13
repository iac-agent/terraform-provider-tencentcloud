## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecord15ById` 方法，使用 DescribeDnsRecords 接口通过 AdvancedFilter 按 id 过滤查询单条 DNS 记录，支持 ReadRetryTimeout 重试

## 2. 资源文件实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_15.go`，定义 ResourceTencentCloudTeoDnsRecord15 函数，包含完整 schema 定义（zone_id、name、type、content、location、ttl、weight、priority、record_id、status、created_on、modified_on）和 Importer 支持
- [x] 2.2 实现 resourceTencentCloudTeoDnsRecord15Create 函数：调用 CreateDnsRecord API，验证返回的 RecordId 非空，设置复合 ID（zone_id + FILED_SP + record_id），调用 Read 填充字段
- [x] 2.3 实现 resourceTencentCloudTeoDnsRecord15Read 函数：解析复合 ID，调用服务层 DescribeTeoDnsRecord15ById，检查 nil 字段后 d.Set() 各字段，资源不存在时设置 d.SetId("")
- [x] 2.4 实现 resourceTencentCloudTeoDnsRecord15Update 函数：检测可变字段变更（name、type、content、location、ttl、weight、priority），调用 ModifyDnsRecords API 更新记录
- [x] 2.5 实现 resourceTencentCloudTeoDnsRecord15Delete 函数：解析复合 ID，调用 DeleteDnsRecords API 删除记录

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 Resources map 中添加 `tencentcloud_teo_dns_record_15` 资源注册
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO 部分添加 `tencentcloud_teo_dns_record_15` 资源条目

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_15.md`，包含一句话描述、Example Usage 和 Import 示例

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_15_test.go`，使用 gomonkey mock 云 API，编写 Create/Read/Update/Delete 操作的单元测试用例
- [x] 5.2 使用 `go test -gcflags=all=-l` 运行单元测试并确认通过

## 6. 代码正确性验证

- [x] 6.1 检查 Create 函数中所有 CreateDnsRecord 入参与云 API 接口参数一致
- [x] 6.2 检查 Read 函数中所有 DescribeDnsRecords 出参与云 API 接口返回字段一致
- [x] 6.3 检查 Update 函数中 ModifyDnsRecords 入参与云 API 接口参数一致
- [x] 6.4 检查 Delete 函数中 DeleteDnsRecords 入参与云 API 接口参数一致
