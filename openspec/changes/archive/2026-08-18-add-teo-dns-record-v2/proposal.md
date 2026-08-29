## Why

用户需要一种更简单、声明式的 TEO DNS 记录管理方式。现有 `tencentcloud_teo_dns_record` 资源通过 `ModifyDnsRecords` / `ModifyDnsRecordsStatus` 支持原地更新，字段与行为较复杂；而在部分运维场景下，DNS 记录一旦创建就不应被原地修改（避免解析中断或状态漂移）。因此新增一个仅提供 Create / Read / Delete 的 `tencentcloud_teo_dns_record_v2` 资源，将记录视为不可变实体，任何字段变更都要求显式重建。

## What Changes

- 新增 TEO 资源 `tencentcloud_teo_dns_record_v2`（RESOURCE_KIND_GENERAL）。
- 资源仅支持 Create / Read / Delete：
  - Create 调用 `CreateDnsRecord`，以返回的 `RecordId` 作为资源标识的一部分。
  - Read 调用 `DescribeDnsRecords`（按 `id` 过滤），回写 `zone_id`、`name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`、`record_id`。
  - Delete 调用 `DeleteDnsRecords`，按 `RecordId` 删除。
- 资源不提供原地更新：`zone_id`（资源身份字段）设置 `ForceNew`，其余顶层字段在 Update 中通过 `immutableArgs` 校验，变更时返回错误。
- 新增 schema 字段：`zone_id`、`name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`、`record_id`。
- 资源 ID 采用复合 ID：`{zoneId}#{recordId}`。
- 在 `tencentcloud/provider.go` / `tencentcloud/provider.md` 中注册资源。
- 新增单元测试 `resource_tc_teo_dns_record_v2_test.go`（使用 gomonkey mock 云 API）。
- 新增资源示例文档 `resource_tc_teo_dns_record_v2.md`。

## Capabilities

### New Capabilities
- `teo-dns-record-v2-resource`: 新增不可变的 TEO DNS 记录资源 `tencentcloud_teo_dns_record_v2` 的 schema 定义、Create/Read/Delete 生命周期、复合 ID 与文档。

### Modified Capabilities
<!-- 无新增修改 capability，本次仅新增 capability -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.go`（新增）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v2_test.go`（新增）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.md`（新增）
  - `tencentcloud/provider.go`、`tencentcloud/provider.md`（注册资源）
- 依赖：复用已 vendored 的 `tencentcloud-sdk-go/tencentcloud/teo/v20220901` 中 `CreateDnsRecord`、`DescribeDnsRecords`、`DeleteDnsRecords`，无需变更 vendor。
- 向后兼容：新增资源，不影响现有 `tencentcloud_teo_dns_record` 资源及任何已有 state。
- 文档：通过 `make doc` 自动生成 `website/docs/` 文档。
