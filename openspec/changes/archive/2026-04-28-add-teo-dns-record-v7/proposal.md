## Why

腾讯云 EdgeOne (TEO) 产品的 DNS 记录管理能力需要新增 `tencentcloud_teo_dns_record_v7` 资源。现有的 `tencentcloud_teo_dns_record` 资源已存在，但需要创建 v7 版本以提供新的资源入口。新资源使用 TEO v20220901 SDK 的 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords` 和 `ModifyDnsRecordsStatus` 接口，实现 DNS 记录的完整 CRUD 生命周期管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_v7`，支持 DNS 记录的创建、读取、更新、删除和导入操作
- 支持的 DNS 记录字段：zone_id、name、type、content、location、ttl、weight、priority、status
- 支持计算字段：created_on、modified_on
- 资源 ID 使用复合 ID 格式：`{zoneId}{FILED_SP}{recordId}`
- 更新操作分为两部分：记录内容修改通过 `ModifyDnsRecords` 接口，状态切换通过 `ModifyDnsRecordsStatus` 接口
- 在 `provider.go` 中注册新资源
- 生成对应的 `.md` 文档示例文件

## Capabilities

### New Capabilities
- `teo-dns-record-v7-resource`: 新增 TEO DNS 记录 v7 资源，支持完整的 CRUD 操作，包含记录内容管理和状态管理

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v7.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v7_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v7.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）
- 修改文件：`tencentcloud/provider.md`（添加资源文档）
- 复用已有服务层方法：`DescribeTeoDnsRecordById`（在 `service_tencentcloud_teo.go` 中已存在）
- 依赖的云 API：CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords、ModifyDnsRecordsStatus
