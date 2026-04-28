## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecordV8ById` 方法，使用 `DescribeDnsRecords` API 并通过 id 过滤器查询 DNS 记录

## 2. Resource Schema & CRUD

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v8.go`，定义 Schema（zone_id/Required+ForceNew, name/Required, type/Required, content/Required, location/Computed+Optional, ttl/Computed+Optional+TypeInt, weight/Computed+Optional+TypeInt, priority/Computed+Optional+TypeInt, status/Computed+Optional, created_on/Computed, modified_on/Computed），支持 Import
- [x] 2.2 实现 Create 函数 `resourceTencentCloudTeoDnsRecordV8Create`，调用 `CreateDnsRecord` API，设置复合 ID（zoneId#recordId），处理返回值为空的情况
- [x] 2.3 实现 Read 函数 `resourceTencentCloudTeoDnsRecordV8Read`，通过 `DescribeTeoDnsRecordV8ById` 查询记录，解析复合 ID，nil 检查后设置各字段
- [x] 2.4 实现 Update 函数 `resourceTencentCloudTeoDnsRecordV8Update`，分两部分：(1) 使用 `ModifyDnsRecords` 修改 name/type/content/location/ttl/weight/priority 字段；(2) 使用 `ModifyDnsRecordsStatus` 修改 status 字段
- [x] 2.5 实现 Delete 函数 `resourceTencentCloudTeoDnsRecordV8Delete`，调用 `DeleteDnsRecords` API

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_dns_record_v8` 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_dns_record_v8` 资源条目

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v8.md` 文档，包含一句话描述（提及 TEO）、Example Usage 和 Import 说明

## 5. Unit Tests

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v8_test.go`，使用 gomonkey mock 方式编写 Create/Read/Update/Delete 的单元测试

## 6. Verification

- [x] 6.1 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
