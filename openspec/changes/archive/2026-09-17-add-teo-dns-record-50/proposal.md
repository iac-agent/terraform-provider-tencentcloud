## Why

TencentCloud EdgeOne (TEO) 的 DNS 记录管理当前已有 `tencentcloud_teo_dns_record` 资源，但该资源未暴露 `DescribeDnsRecords` 接口提供的查询增强能力（`filters` 过滤、`sort_by`/`sort_order` 排序、`match` 匹配方式）。本次按云 API 最新参数映射新增 `tencentcloud_teo_dns_record_50` 资源，提供 TEO DNS 记录的完整 CRUD 生命周期管理，并覆盖 Describe 接口的全部入参能力。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_50`（RESOURCE_KIND_GENERAL），支持 TEO DNS 记录的完整 CRUD 生命周期管理
- Create：调用 `CreateDnsRecord` 接口创建 DNS 记录（zone_id/name/type/content 必填，location/ttl/weight/priority 可选）
- Read：调用 `DescribeDnsRecords` 接口按 record id 过滤查询记录详情，同时暴露 filters/sort_by/sort_order/match 查询参数
- Update：调用 `ModifyDnsRecords` 接口批量修改 DNS 记录（单条记录场景，封装为长度为 1 的 DnsRecords 列表）
- Delete：调用 `DeleteDnsRecords` 接口批量删除 DNS 记录（单条记录场景，RecordIds 封装为长度为 1 的列表）
- 资源 ID 使用 `zone_id` 与 `record_id` 的联合 ID（分隔符 `tccommon.FILED_SP`，即 `#`）
- 计算属性：`record_id`、`status`、`created_on`、`modified_on`

## Capabilities

### New Capabilities
- `teo-dns-record-50-resource`: 提供 `tencentcloud_teo_dns_record_50` 资源的完整 CRUD 生命周期管理，包括创建、读取、更新和删除 TEO DNS 记录

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_dns_record_50.go`
- 新增测试文件: `tencentcloud/services/teo/resource_tc_teo_dns_record_50_test.go`（gomonkey mock 单元测试）
- 新增文档文件: `tencentcloud/services/teo/resource_tc_teo_dns_record_50.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（在 TEO Resource 列表添加资源名）
- 依赖云 API SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中，无需改动）
