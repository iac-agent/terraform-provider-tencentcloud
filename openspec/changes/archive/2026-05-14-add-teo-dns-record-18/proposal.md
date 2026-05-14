## Why

TEO (EdgeOne) 产品需要新增 `tencentcloud_teo_dns_record_18` 资源，以支持通过 Terraform 管理 TEO DNS 记录的完整生命周期（创建、读取、更新、删除）。该资源基于最新的云API接口 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords` 实现，为用户提供标准化的 DNS 记录管理能力。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_18`，文件名为 `resource_tc_teo_dns_record_18.go`
- 实现完整的 CRUD 操作：Create（调用 CreateDnsRecord）、Read（调用 DescribeDnsRecords）、Update（调用 ModifyDnsRecords）、Delete（调用 DeleteDnsRecords）
- 资源 ID 由 `zone_id` 和 `record_id` 组合而成，使用 `FILED_SP` 分隔符
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源
- 新增对应的单元测试文件 `resource_tc_teo_dns_record_18_test.go`（使用 gomonkey mock）
- 新增资源文档 `resource_tc_teo_dns_record_18.md`

## Capabilities

### New Capabilities
- `teo-dns-record-18-resource`: TEO DNS 记录资源的完整 CRUD 管理，支持 zone_id、name、type、content、location、ttl、weight、priority 等字段的创建和更新，以及 record_id、status、created_on、modified_on 等只读字段的读取

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_18.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_18_test.go`
- 新增文件：`gendoc/resource/tencentcloud_teo_dns_record_18.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）
- 修改文件：`tencentcloud/provider.md`（添加资源文档条目）
- 依赖的云API包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- 复用现有的 TeoService 服务层中的 `DescribeTeoDnsRecordById` 方法
