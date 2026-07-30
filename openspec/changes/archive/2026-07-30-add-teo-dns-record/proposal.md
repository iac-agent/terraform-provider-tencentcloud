## Why

需要为云产品 TEO (TencentCloud EdgeOne) 新增 `tencentcloud_teo_dns_record` 资源，以支持通过 Terraform 管理 TEO DNS 记录的完整生命周期（创建、读取、更新、删除）。该资源基于 TEO 的 CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords 四个云 API 接口实现。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record`，支持 TEO DNS 记录的 CRUD 操作
- 新增资源代码文件 `tencentcloud/services/teo/resource_tc_teo_dns_record.go`
- 新增资源测试文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_test.go`
- 新增资源文档文件 `tencentcloud/services/teo/resource_tc_teo_dns_record.md`
- 在 `tencentcloud/provider.go` 中注册新资源
- 在 `tencentcloud/provider.md` 中添加新资源文档条目

## Capabilities

### New Capabilities
- `teo-dns-record-resource`: TEO DNS 记录资源的完整 CRUD 实现，支持 zone_id、name、type、content、location、ttl、weight、priority 等参数，以及 status、created_on、modified_on 等只读计算字段

### Modified Capabilities

（无）

## Impact

- 新增资源文件：`tencentcloud/services/teo/resource_tc_teo_dns_record.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_test.go`
- 新增文档文件：`tencentcloud/services/teo/resource_tc_teo_dns_record.md`
- 修改 `tencentcloud/provider.go`：添加资源注册行
- 修改 `tencentcloud/provider.md`：添加资源文档条目
- 依赖云 API：`teo/v20220901` 包中的 CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords、ModifyDnsRecordsStatus 接口
- 复用已有的服务层方法 `TeoService.DescribeTeoDnsRecordById`