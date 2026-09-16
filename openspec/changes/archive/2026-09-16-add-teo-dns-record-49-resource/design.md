## Context

腾讯云 Terraform Provider（terraform-provider-tencentcloud）的 teo 服务目录 `tencentcloud/services/teo/` 下已存在资源 `tencentcloud_teo_dns_record`（`resource_tc_teo_dns_record.go`），其 Read 依赖服务层 `service_tencentcloud_teo.go` 中的 `DescribeTeoDnsRecordById(ctx, zoneId, recordId)`（内部调用 `DescribeDnsRecords`，用 `Filters` 的 `id` 条件过滤）。

本变更按需求映射规范新增一个独立的新资源 `tencentcloud_teo_dns_record_49`（RESOURCE_KIND_GENERAL），管理单个 EdgeOne DNS 记录。四个云 API 接口均已在 vendor 中验证存在：

- `CreateDnsRecord`：入参 `ZoneId`、`Name`、`Type`、`Content`（必填）、`Location`、`TTL`、`Weight`、`Priority`（可选）；出参 `RecordId`。
- `DescribeDnsRecords`：入参 `ZoneId`（必填）、`Filters []*AdvancedFilter`（`Name`/`Values`/`Fuzzy`）、`SortBy`、`SortOrder`、`Match`（可选）+ 分页 `Offset`/`Limit`（Limit 上限 1000）；出参 `TotalCount` 与 `DnsRecords []*DnsRecord`。
- `ModifyDnsRecords`：入参 `ZoneId`（必填）、`DnsRecords []*DnsRecord`（单次最多 100 条）；响应仅 RequestId。
- `DeleteDnsRecords`：入参 `ZoneId`、`RecordIds []*string`（上限 1000）；响应仅 RequestId。

`DnsRecord` 结构体包含 `ZoneId`、`RecordId`、`Name`、`Type`、`Location`、`Content`、`TTL`、`Weight`、`Priority`、`Status`、`CreatedOn`、`ModifiedOn`。注意云 API 注释明确说明：`ZoneId`、`Status`、`CreatedOn`、`ModifiedOn` 在 `ModifyDnsRecords` 中不可作为入参使用（如有传此参数会被忽略）。

## Goals / Non-Goals

**Goals:**

- 新增资源 `tencentcloud_teo_dns_record_49`，支持单个 DNS 记录的完整 CRUD 生命周期（创建、查询回填、更新、删除），并在 provider 中注册。
- 严格遵循 provider 的既有代码风格（参考 `resource_tc_teo_dns_record.go` 与 `resource_tc_igtm_strategy.go`）：`tccommon.LogElapsed` / `tccommon.InconsistentCheck` defer、`tccommon.NewResourceLifeCycleHandleFuncContext`、retry 包装、`helper.String/IntInt64` 等。
- 资源 ID 采用复合 ID `zoneId#recordId`（`tccommon.FILED_SP`），支持 import，且在 Read/Update/Delete 中均从 `d.Id()` 解析。
- 提供符合规范的 `.md` 文档与 gomonkey 单元测试。

**Non-Goals:**

- 不修改既有资源 `tencentcloud_teo_dns_record`（保持向后兼容，避免破坏现有 TF 配置和 state）。
- 不实现 DNS 记录的启停状态切换（`ModifyDnsRecordsStatus` 接口不在本次需求映射范围内，`status` 仅作为只读回填字段）。
- 不直接修改 `website/` 目录（文档由收尾阶段 `make doc` 生成）。
- 不在本阶段执行 `gofmt`、`make doc`、`.changelog` 文件创建（统一由收尾阶段 tfpacer-finalize skill 执行）。

## Decisions

1. **文件与函数命名**：新资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_49.go`，导出函数 `ResourceTencentCloudTeoDnsRecord49()`，内部 CRUD 函数 `resourceTencentCloudTeoDnsRecord49Create/Read/Update/Delete`。理由：遵循 `resource_tc_<Product>_<资源名>.go` 命名格式与既有风格（数字后缀直接拼在资源名后，与需求给定的资源名一致）。

2. **Schema 设计**（依据参数映射的必填/可选标注，以及 `DnsRecord` 出参结构）：
   - `zone_id`：`TypeString`，Required + ForceNew（Delete/Modify 均需 ZoneId，且记录不能跨 Zone 迁移）。
   - `name`、`type`、`content`：`TypeString`，Required（Create 必填）。
   - `location`：`TypeString`，Optional + Computed（云 API 不传时默认 "Default"）。
   - `ttl`、`weight`、`priority`：`TypeInt`，Optional + Computed（云 API 有默认值：TTL 默认 300、Weight 默认 -1、Priority 默认 0；使用 `GetOkExists` 读取以区分"未设置"）。
   - `status`、`created_on`、`modified_on`：`TypeString`，Computed（云 API 在 `ModifyDnsRecords` 中忽略这些字段，仅做查询回填；不提供可写入口，避免用户设置后产生无效 diff）。
   - `record_id`：`TypeString`，Computed（Create 出参回填，同时便于用户在输出中引用记录 ID）。
   - 资源支持 Import（`schema.ImportStatePassthrough`），import 时使用复合 ID `zoneId#recordId`。

3. **Read 实现**：复用服务层现有方法 `TeoService.DescribeTeoDnsRecordById(ctx, zoneId, recordId)`，该方法已实现 `Filters: id=recordId` 过滤 + `tccommon.ReadRetryTimeout` retry。不新增 Describe 的顶层 schema 参数（`filters`/`sort_by`/`sort_order`/`match` 属于 `DescribeDnsRecords` 列表查询语义，不属于单资源 Read 语义，故不纳入资源 schema；这符合"资源参数映射以 Create/Modify 入参为主体"的约束）。仅当 `respData == nil` 时先 `log.Printf("[CRUD] teo dns_record_49 id=%s", d.Id())` 再 `d.SetId("")`。回填字段前均判空。

4. **Update 实现**：可变字段为 `name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`（与既有 `tencentcloud_teo_dns_record` 的 mutableArgs 一致，且均为 `ModifyDnsRecords` 实际支持的入参字段）。当任一字段 `d.HasChange()` 时，构造单条 `teov20220901.DnsRecord{RecordId: recordId, ...}`（不设置 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn`，云 API 会忽略），调用 `ModifyDnsRecords`（WriteRetryTimeout retry）。`zone_id` 因 ForceNew 不会进入 Update 流程。

5. **Delete 实现**：从 `d.Id()` 解析出 `zoneId`、`recordId`，调用 `DeleteDnsRecords`，`RecordIds` 传单条记录，WriteRetryTimeout retry。

6. **Create 的空返回校验**：调用成功后校验 `response.Response == nil || response.Response.RecordId == nil || *response.Response.RecordId == ""`，命中则返回 `tccommon.NonRetryableError`，并在校验前打印 `logId` 与请求上下文，避免写入空 ID 造成 state 混乱。校验逻辑放在 retry 块外（成功路径）。

7. **分页参数取值**：`DescribeDnsRecords` 的 `Limit` 云 API 标注上限 1000。服务层现有 `DescribeTeoDnsRecordById` 未显式设置 Limit（默认 20），为保证按 id 精确过滤时一定能查到目标记录，本次在其内部设置 `request.Limit = helper.IntUint64(1000)`（等价 Int64），Offset 保持默认 0。此改动只影响新资源的查询路径（该服务层方法同时被旧资源复用，但按 id 精确过滤语义下提高 Limit 只会增强可靠性，不改变既有行为契约）。

8. **单元测试**：新增 `resource_tc_teo_dns_record_49_test.go`，使用 gomonkey patch `TeoService` 客户端相关方法 / `GetAPIV3Conn`，对 Create/Read/Update/Delete 业务逻辑做 mock 测试（不使用 terraform 测试套件），参考仓库内既有 gomonkey 风格测试文件。

9. **文档**：新增 `resource_tc_teo_dns_record_49.md`，格式为一句话描述（带云产品名 EdgeOne）+ Example Usage + Import（说明使用复合 ID `zoneId#recordId`）。不手写 Argument/Attribute Reference（由工具自动生成）。

## Risks / Trade-offs

- [服务层方法被旧资源复用] `DescribeTeoDnsRecordById` 同时被 `tencentcloud_teo_dns_record` 使用 → 本次只在该方法内补充 `Limit=1000`，不改变函数签名与返回语义；若评审认为必须完全隔离，可改为新增 `DescribeTeoDnsRecord49ById` 独立方法（备选方案，默认不采用，避免重复代码）。
- [`GetOkExists` 语义] `ttl`/`weight`/`priority` 使用 `GetOkExists` 以区分"用户未设置"与"设置为 0"（weight=0 表示不解析、priority=0 为默认值，均为合法值）→ 与既有 teo 资源保持一致的读取方式；风险是 SDK 内部对该 API 的行为约束，本仓库已大量使用该模式，风险可控。
- [ModifyDnsRecords 单次上限 100 条] 单资源更新只传 1 条，远低于上限 → 无需分批。
- [Create 后立即 Read 可能为空] 创建为同步接口，云 API 返回 RecordId 即表示记录已生成；Read 失败（记录不存在）时按资源被删除处理并打日志 → 遵循 provider 通用模式。
- [旧资源共存] 新资源名带 `_49` 后缀与旧 `tencentcloud_teo_dns_record` 并存，用户可能混淆 → 在文档描述中明确该资源基于新版接口映射规范实现；不删除旧资源（硬约束：向后兼容）。
