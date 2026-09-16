## Why

EdgeOne（teo）现有 Terraform 资源 `tencentcloud_teo_dns_record` 基于旧版 DNS 记录管理接口实现。本次需求要求按新的接口与参数映射规范，新增一个独立的通用资源 `tencentcloud_teo_dns_record_49`，用于对单个 EdgeOne DNS 记录（站点 Zone 下的解析记录）进行完整的增删改查生命周期管理，为用户提供符合最新云 API 规范的资源管理能力。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_49`（RESOURCE_KIND_GENERAL，通用资源，文件名 `resource_tc_teo_dns_record_49.go`），实现单个 DNS 记录的完整 CRUD 生命周期：
  - **Create**：调用 `CreateDnsRecord` 创建 DNS 记录，返回 `RecordId` 作为资源 ID 的一部分，复合 ID 采用 `zoneId#recordId`（`tccommon.FILED_SP` 分隔）。
  - **Read**：调用 `DescribeDnsRecords`（通过 `Filters` 的 `id` 过滤条件，`Limit` 取云 API 标注的最大值 1000）查询该 DNS 记录并回填状态字段（`status`、`created_on`、`modified_on` 等）。
  - **Update**：调用 `ModifyDnsRecords`（批量接口，仅传入本资源对应的单条记录）修改可变字段（`name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`）。
  - **Delete**：调用 `DeleteDnsRecords`（批量接口，仅传入本资源对应的单条 `RecordIds`）删除该 DNS 记录。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册 `tencentcloud_teo_dns_record_49` 资源。
- 新增资源文档 `resource_tc_teo_dns_record_49.md`（一句话描述 + Example Usage + Import 说明，Argument/Attribute Reference 由工具自动生成）。
- 新增单元测试文件 `resource_tc_teo_dns_record_49_test.go`，使用 gomonkey mock 云 API 客户端进行业务逻辑测试（不使用 terraform 测试套件）。

不修改已有资源 `tencentcloud_teo_dns_record` 的任何 schema 或实现，完全向后兼容。

## Capabilities

### New Capabilities
- `teo-dns-record-49-resource`: 通过 `CreateDnsRecord` / `DescribeDnsRecords` / `ModifyDnsRecords` / `DeleteDnsRecords` 四个云 API 管理单个 EdgeOne DNS 记录（`tencentcloud_teo_dns_record_49`）的完整生命周期，包括复合 ID 设计、可变字段更新、状态/时间等只读字段回填。

### Modified Capabilities

（无 —— 本变更仅新增资源，不改变任何既有能力的需求。）

## Impact

- **代码**：
  - 新增 `tencentcloud/services/teo/resource_tc_teo_dns_record_49.go`（资源 schema + CRUD 实现，复用 `service_tencentcloud_teo.go` 中已有的 `DescribeTeoDnsRecordById` 服务层查询能力，如不满足则新增服务层方法）。
  - 新增 `tencentcloud/services/teo/resource_tc_teo_dns_record_49_test.go`（gomonkey mock 单元测试）。
  - 修改 `tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（资源清单）。
- **云 API 依赖**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 中以下接口（已存在于 vendor 中，已验证）：
  - `CreateDnsRecord`（入参 ZoneId/Name/Type/Content/Location/TTL/Weight/Priority，出参 RecordId）。
  - `DescribeDnsRecords`（入参 ZoneId/Filters/SortBy/SortOrder/Match，出参 DnsRecords 列表及 TotalCount）。
  - `ModifyDnsRecords`（入参 ZoneId/DnsRecords，其中 DnsRecord 的 ZoneId/Status/CreatedOn/ModifiedOn 会被服务端忽略，仅做出参使用）。
  - `DeleteDnsRecords`（入参 ZoneId/RecordIds）。
- **行为约束**：
  - `zone_id` 为 ForceNew（创建后不可变）；其余顶层参数为可变参数。
  - 状态字段 `status`、`created_on`、`modified_on` 在 `ModifyDnsRecords` 中被云 API 忽略，故只作为 Computed 只读字段，不做更新传参。
  - 所有云 API 调用均需使用 `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` 做 retry 包装，接口失败使用 `tccommon.RetryError()` 包装。
- **文档**：新增 `resource_tc_teo_dns_record_49.md`，最终通过 `make doc` 生成 `website/` 下文档（不在本变更中直接修改 `website/` 目录）。
