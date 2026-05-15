## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_23.go`，包含 schema 定义和 CRUD 函数（参考 resource_tc_teo_dns_record_21.go）
- [x] 1.2 实现 Create 函数：调用 CreateDnsRecord API，设置复合 ID（zone_id#record_id），包含 nil 检查和 NonRetryableError 处理
- [x] 1.3 实现 Read 函数：复用 TeoService.DescribeTeoDnsRecordById，设置所有 schema 字段，包含 nil 检查
- [x] 1.4 实现 Update 函数：检测 mutableArgs 变更，调用 ModifyDnsRecords API
- [x] 1.5 实现 Delete 函数：调用 DeleteDnsRecords API，使用 RecordIds 传入单个 record_id
- [x] 1.6 实现 Import：使用 schema.ImportStatePassthrough，支持 zone_id#record_id 格式

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中添加 `tencentcloud_teo_dns_record_23` 资源注册行
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_23` 文档条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_23.md`，包含描述、Example Usage 和 Import 部分

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_23_test.go`，使用 gomonkey mock 方式编写 CRUD 单元测试
- [x] 4.2 运行 `go test -gcflags=all=-l` 验证单元测试通过

## 5. 代码正确性验证

- [x] 5.1 检查 Create 函数中的参数是否在 CreateDnsRecord API 入参中存在
- [x] 5.2 检查 Update 函数中的参数是否在 ModifyDnsRecords API 入参中存在
- [x] 5.3 检查 Read 函数中的参数是否在 DescribeDnsRecords API 出参中存在
- [x] 5.4 检查 Delete 函数中的参数是否在 DeleteDnsRecords API 入参中存在
