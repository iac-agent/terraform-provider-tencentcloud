## Why

目前 TEO 产品的 DNS 记录管理已有 `tencentcloud_teo_dns_record` 资源，但需要新增 `tencentcloud_teo_dns_record_10` 资源以支持 DNS 记录的完整 CRUD 生命周期管理，包括创建、读取、更新和删除操作。该新资源将基于最新的云 API 接口规范实现，确保参数映射与 API 接口完全一致。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_10`，用于管理 TEO DNS 记录的完整生命周期
- 实现资源 Create 方法，调用 `CreateDnsRecord` 接口，支持 zone_id、name、type、content、location、ttl、weight、priority 参数
- 实现资源 Read 方法，调用 `DescribeDnsRecords` 接口，通过 zone_id 和 record_id 查询单条 DNS 记录
- 实现资源 Update 方法，调用 `ModifyDnsRecords` 接口，支持修改 name、type、content、location、ttl、weight、priority 参数
- 实现资源 Delete 方法，调用 `DeleteDnsRecords` 接口，通过 zone_id 和 record_id 删除 DNS 记录
- 使用复合 ID（zone_id + FILED_SP + record_id）作为资源标识
- 支持 Terraform Import 功能
- 在 provider.go 和 provider.md 中注册新资源
- 编写单元测试和资源文档

## Capabilities

### New Capabilities
- `teo-dns-record-10-resource`: TEO DNS 记录资源的完整 CRUD 生命周期管理，包括 CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords 四个云 API 接口的封装

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_10.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_10_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_10.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）
- 修改文件：`tencentcloud/provider.md`（添加资源文档条目）
- 修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（添加 DescribeTeoDnsRecord10ById 服务方法）
- 依赖云 API 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
