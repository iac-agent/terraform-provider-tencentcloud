## Why

TEO（EdgeOne）DNS 记录需要以 Terraform 资源形式进行声明式管理。虽然已有 `tencentcloud_teo_dns_record` 资源封装了同一组云 API，但本次需新增一个独立的 `tencentcloud_teo_dns_record_v3` 资源，将 DNS 记录的增删改查（`CreateDnsRecord` / `DescribeDnsRecords` / `ModifyDnsRecords` / `DeleteDnsRecords`）以更清晰的生命周期语义暴露，并将 `status`、`created_on`、`modified_on` 等云端只读字段收敛为纯 Computed 输出，避免用户通过 Terraform 修改状态造成与云端只读语义不一致。

## What Changes

- 新增 RESOURCE_KIND_GENERAL 资源 `tencentcloud_teo_dns_record_v3`，提供完整 CRUD：
  - Create 调用 `CreateDnsRecord` 创建单条 DNS 记录；
  - Read 调用 `DescribeDnsRecords`（通过 `Filters` 中 `id` 过滤条件）查询单条记录；
  - Update 调用 `ModifyDnsRecords`（批量接口，单条记录场景仅传一个元素）修改可变更字段；
  - Delete 调用 `DeleteDnsRecords`（批量接口，单条记录场景仅传一个 `RecordId`）删除记录。
- 定义资源 schema，字段如下：
  - Required：`zone_id`（ForceNew）、`name`、`type`、`content`；
  - Optional + Computed：`location`、`ttl`、`weight`、`priority`；
  - Computed：`record_id`、`status`、`created_on`、`modified_on`。
- 资源 ID 采用 `zone_id#record_id` 复合 ID（`tccommon.FILED_SP` 分隔），支持 `terraform import`。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册资源。
- 新增资源示例文档 `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.md`（最终由 `make doc` 生成 website docs）。

非破坏性：本次为纯新增资源，不修改既有 `tencentcloud_teo_dns_record` 资源及其 schema/state。

## Capabilities

### New Capabilities

- `teo-dns-record-v3-resource`: 提供 `tencentcloud_teo_dns_record_v3` 资源的 CRUD 生命周期、schema 定义、复合 ID 与 import 语义。

### Modified Capabilities

<!-- 无：本次为纯新增 capability，不修改既有 capability 的 requirement -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.go`（新增，schema 与 CRUD 实现）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v3_test.go`（新增，单测）
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.md`（新增，示例文档）
  - `tencentcloud/provider.go` / `tencentcloud/provider.md`（资源注册）
- 依赖：复用已 vendored 的 `tencentcloud-sdk-go/tencentcloud/teo/v20220901`（`CreateDnsRecord` / `DescribeDnsRecords` / `ModifyDnsRecords` / `DeleteDnsRecords`），无需变更 vendor。
- 向后兼容：纯新增资源，不影响既有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由 `make doc` 自动生成）。
