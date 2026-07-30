## Context

当前 TEO DNS 记录资源（`tencentcloud_teo_dns_record`）已经存在于代码库中，具备完整的 CRUD 实现。本变更提案旨在形式化地记录该资源的设计与实现，作为 openspec 变更的一部分。

已有资源实现使用以下云 API 接口：
- `CreateDnsRecord`：创建 DNS 记录
- `DescribeDnsRecords`：查询 DNS 记录（通过 `TeoService.DescribeTeoDnsRecordById` 封装）
- `ModifyDnsRecords`：修改 DNS 记录（含 `ModifyDnsRecordsStatus` 用于启用/禁用状态变更）
- `DeleteDnsRecords`：删除 DNS 记录

资源采用复合 ID（`zone_id#record_id`）作为唯一标识，支持 import 操作。

## Goals / Non-Goals

**Goals:**
- 提供 TEO DNS 记录的完整 CRUD Terraform 资源
- 支持 DNS 记录的启用/禁用状态管理
- 支持通过复合 ID 导入已有 DNS 记录
- 代码风格与现有 TEO 资源（如 `tencentcloud_teo_dns_record_23`）保持一致

**Non-Goals:**
- 不提供批量操作（批量创建/修改/删除多条记录）
- 不修改已有 versioned 资源（`_21`, `_22`, `_23`, `_24`）的行为

## Decisions

### 1. 复合 ID 设计
使用 `zone_id#record_id` 格式作为资源 ID。`zone_id` 设置为 `ForceNew`，确保 zone 变更时资源重建。

**理由**：TEO DNS 记录的唯一标识需要同时包含 zone_id 和 record_id，因为 record_id 在不同 zone 之间可能不唯一。使用 `tccommon.FILED_SP`（`#`）作为分隔符与现有 TEO 资源一致。

### 2. Schema 字段设计
- 必填字段：`zone_id`（ForceNew）、`name`、`type`、`content`
- 可选字段：`location`、`ttl`、`weight`、`priority`、`status`
- 计算字段：`created_on`、`modified_on`

**理由**：`zone_id` 设为 ForceNew 是因为 DNS 记录绑定到特定 zone，不应跨 zone 迁移。`location`、`ttl`、`weight`、`priority` 为可选，云 API 有默认值。`status` 为可选+计算，支持用户主动设置或从 API 读取。

### 3. Update 实现分离
将 DNS 记录内容更新（`ModifyDnsRecords`）和状态更新（`ModifyDnsRecordsStatus`）分开处理，各自独立检测变更。

**理由**：云 API 将修改记录内容和修改记录状态分为两个独立接口，分离处理避免不必要的 API 调用。

### 4. 服务层复用
Read 操作复用 `TeoService.DescribeTeoDnsRecordById` 方法，该方法通过 `DescribeDnsRecords` 接口按 record_id 过滤查询。

**理由**：`DescribeTeoDnsRecordById` 已由其他 versioned 资源实现并验证，复用避免重复代码。

### 5. 错误处理
- 所有 API 调用使用 `retry` + `tccommon.RetryError` 包装
- 超时使用 `tccommon.ReadRetryTimeout`（读）和 `tccommon.WriteRetryTimeout`（写）
- Read 时若记录不存在，设置 `d.SetId("")` 并记录日志

## Risks / Trade-offs

- **风险**：`DescribeTeoDnsRecordById` 返回多条记录时只取第一条 → 通过 `id` 过滤确保唯一性，正常情况不会返回多条
- **风险**：`ModifyDnsRecords` 接口是批量接口（接受 `DnsRecords` 数组），但资源只操作单条记录 → 传入只含一个元素的数组，符合接口设计
- **风险**：`status` 字段更新依赖于 `ModifyDnsRecordsStatus` 接口，该接口传入 `RecordsToEnable` 或 `RecordsToDisable` → 实现中根据 status 值选择对应参数