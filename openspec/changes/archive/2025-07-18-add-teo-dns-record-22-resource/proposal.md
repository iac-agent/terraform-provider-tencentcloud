## Why

TEO (EdgeOne) DNS 记录资源需要新增 `tencentcloud_teo_dns_record_22` 资源，以提供更完整的 DNS 记录管理能力。当前已有的 `tencentcloud_teo_dns_record` 资源已在使用，但需要新增一个独立命名的资源来满足特定场景下的 DNS 记录管理需求。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_teo_dns_record_22`，支持 DNS 记录的完整 CRUD 操作
- Create: 调用 `CreateDnsRecord` 接口创建 DNS 记录，入参包括 ZoneId、Name、Type、Content、Location、TTL、Weight、Priority，出参 RecordId
- Read: 调用 `DescribeDnsRecords` 接口查询 DNS 记录，通过 ZoneId 和 RecordId 过滤获取单条记录
- Update: 调用 `ModifyDnsRecords` 接口修改 DNS 记录，支持修改 Name、Type、Content、Location、TTL、Weight、Priority 字段；调用 `ModifyDnsRecordsStatus` 接口修改 DNS 记录状态
- Delete: 调用 `DeleteDnsRecords` 接口删除 DNS 记录
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 新增资源文档 `resource_tc_teo_dns_record_22.md`

## Capabilities

### New Capabilities
- `teo-dns-record-22`: TEO DNS 记录资源的 CRUD 管理，包括创建、读取、更新、删除 DNS 记录，支持 Name、Type、Content、Location、TTL、Weight、Priority 等参数

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_22.go`、`tencentcloud/services/teo/resource_tc_teo_dns_record_22_test.go`、`tencentcloud/services/teo/resource_tc_teo_dns_record_22.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（新增资源文档条目）
- 依赖的云 API 接口：CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、ModifyDnsRecordsStatus、DeleteDnsRecords（teo v20220901 包）
