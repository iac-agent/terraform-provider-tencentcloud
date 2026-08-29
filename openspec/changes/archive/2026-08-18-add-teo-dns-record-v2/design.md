## Context

当前 TEO 服务已存在 `tencentcloud_teo_dns_record`（v1）资源，封装了 `CreateDnsRecord` / `DescribeDnsRecords` / `ModifyDnsRecords` / `ModifyDnsRecordsStatus` / `DeleteDnsRecords`，支持记录内容与启停状态的原地更新。本次新增 `tencentcloud_teo_dns_record_v2`（v2），聚焦"创建后不可原地变更"的场景，仅映射以下云 API：

| Terraform 操作 | 云 API |
|----------------|--------|
| Create | `CreateDnsRecord` |
| Read | `DescribeDnsRecords`（按 `id` 过滤） |
| Delete | `DeleteDnsRecords` |

约束：
- 复用已 vendored 的 `teov20220901` SDK，不新增依赖。
- 向后兼容：不改动现有 v1 资源。
- 资源 ID 沿用项目复合 ID 规范：`{zoneId}#{recordId}`（`tccommon.FILED_SP` 分隔）。
- 云 API 会为 `location` / `ttl` / `weight` / `priority` 提供默认值，因此这些字段使用 `Optional + Computed`。

## Goals / Non-Goals

**Goals:**
- 新增一个只读 + 创建 + 删除（不可原地更新）的 TEO DNS 记录资源。
- 记录字段变更时通过 Update 返回错误，引导用户显式重建。
- 提供 gomonkey 单测覆盖 Create / Read / Delete 与 Update 不可变校验。
- 提供 `.md` 示例文档，支持 `make doc` 生成 website 文档。

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_dns_record` 资源。
- 不引入 `ModifyDnsRecords` / `ModifyDnsRecordsStatus` 更新能力。
- 不暴露 `filters` / `sort_by` / `sort_order` / `match` / `offset` / `limit` 等查询参数给用户（这些仅作为 Read 内部实现）。
- 不新增 `status` / `created_on` / `modified_on` 等未在本次映射中的输出字段。

## Decisions

### Decision 1: 资源为不可变（CRD-only），Update 通过 immutableArgs 返回错误

**选择**：schema 中仅 `zone_id` 设置 `ForceNew`；`name` / `type` / `content` / `location` / `ttl` / `weight` / `priority` 不设置 `ForceNew`，而是在 Update 函数中通过 `immutableArgs` 校验，发现变化即 `return fmt.Errorf("argument `%s` cannot be changed", v)`。

**备选**：所有字段都设置 `ForceNew`，让 Terraform 自动 destroy + create。

**理由**：
- 与项目内 CRD-only 资源（如 `cfs_user_quota`）的实现模式一致。
- 对 DNS 记录这类资源，自动 destroy + create 可能造成解析中断，显式报错更安全，让用户有意识地 `terraform taint` / 重建。

### Decision 2: 资源 ID 使用 `{zoneId}#{recordId}`

**选择**：`d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))`。

**理由**：
- 与现有 v1 资源及项目规范一致，import 时也使用该复合 ID。

### Decision 3: Read 复用 `TeoService.DescribeTeoDnsRecordById`

**选择**：复用 `service_tencentcloud_teo.go` 中已有的 `DescribeTeoDnsRecordById`（内部通过 `DescribeDnsRecords` + `Filters=[{Name:"id", Values:[recordId]}]` 查询）。

**理由**：避免重复实现，保持 v1/v2 查询逻辑一致。

### Decision 4: 字段类型与 Optional/Computed

`location` / `ttl` / `weight` / `priority` 使用 `Optional + Computed`，因为云 API 存在默认值（`Default` / `300` / `-1` / `0`），Read 需回填。

### Decision 5: Create 响应为空校验

Create 调用 `CreateDnsRecord` 后，校验 `response == nil || response.Response == nil || response.Response.RecordId == nil`，为空则返回 `resource.NonRetryableError`，避免写入空 id。

## Risks / Trade-offs

- [Risk] 用户修改不可变字段时 apply 直接失败 → Mitigation：在文档与 schema Description 中明确"不可变"，引导用户重建。
- [Risk] Read 复用现有 `DescribeTeoDnsRecordById` 未显式设置 `Limit`，若未来 `id` 过滤语义变化可能返回多条 → Mitigation：该 helper 已按 `id` 精确过滤，云 API 语义稳定；如需要再单独补充分页。
- [Trade-off] 不可变设计牺牲了原地更新便利性 → 与本次"仅 CRD"需求一致，是刻意取舍。

## Migration Plan

- 新增资源，无状态迁移。
- 文档：新增 `resource_tc_teo_dns_record_v2.md`，`make doc` 生成 website 文档。
- 回滚：删除资源注册与文件即可，不影响现有资源。

## Open Questions

- 无
