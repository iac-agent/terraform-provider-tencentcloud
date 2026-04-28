## Context

腾讯云 EdgeOne (TEO) 产品需要新增 `tencentcloud_teo_dns_record_v7` 资源。已有的 `tencentcloud_teo_dns_record` 资源实现了相同的 API 逻辑，新资源是独立的新版本入口。

当前状态：
- 已有 `tencentcloud_teo_dns_record` 资源（`resource_tc_teo_dns_record.go`）
- 已有服务层方法 `DescribeTeoDnsRecordById`（在 `service_tencentcloud_teo.go` 中）
- 云 API SDK 位于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`

涉及的云 API 接口：
- **CreateDnsRecord**：创建 DNS 记录，返回 RecordId
- **DescribeDnsRecords**：查询 DNS 记录列表，支持按 id 过滤（AdvancedFilter）
- **ModifyDnsRecords**：批量修改 DNS 记录内容（name, type, content, location, ttl, weight, priority）
- **ModifyDnsRecordsStatus**：切换 DNS 记录状态（enable/disable）
- **DeleteDnsRecords**：批量删除 DNS 记录

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_v7` 资源，实现完整 CRUD 操作
- 支持所有 DNS 记录类型（A, AAAA, MX, CNAME, TXT, NS, CAA, SRV）
- 支持记录内容修改（通过 ModifyDnsRecords）和状态切换（通过 ModifyDnsRecordsStatus）
- 使用复合 ID 格式 `{zoneId}{FILED_SP}{recordId}` 标识资源
- 支持 Import 导入
- 代码风格严格参考 `resource_tc_igtm_strategy.go`

**Non-Goals:**
- 不修改已有的 `tencentcloud_teo_dns_record` 资源
- 不新增服务层方法（复用已有的 `DescribeTeoDnsRecordById`）
- 不支持批量创建/修改（Terraform 资源管理单条记录）

## Decisions

### 1. 复合 ID 格式

**决定**：使用 `{zoneId}{FILED_SP}{recordId}` 作为资源 ID。

**原因**：
- CreateDnsRecord 返回的 RecordId 需要配合 ZoneId 才能唯一定位记录
- DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords 都需要 ZoneId 参数
- 与已有 `tencentcloud_teo_dns_record` 资源保持一致的 ID 格式

### 2. Update 操作拆分为两步

**决定**：Update 函数中，记录内容字段变更调用 `ModifyDnsRecords`，状态字段变更调用 `ModifyDnsRecordsStatus`。

**原因**：
- 云 API 将记录内容修改和状态切换设计为两个独立接口
- `ModifyDnsRecords` 用于修改 name、type、content、location、ttl、weight、priority
- `ModifyDnsRecordsStatus` 用于启用/禁用记录
- 两个接口互不影响，需要分别调用

### 3. 服务层方法复用

**决定**：复用已有的 `DescribeTeoDnsRecordById` 服务层方法进行 Read 操作。

**原因**：
- 该方法已存在于 `service_tencentcloud_teo.go` 中，封装了 `DescribeDnsRecords` 的调用和重试逻辑
- 使用 AdvancedFilter 按 id 过滤，精确查询单条记录
- 避免代码重复

### 4. Schema 设计

**决定**：
- `zone_id`：Required, ForceNew（ZoneId 在创建后不可变更）
- `name`、`type`、`content`：Required（DNS 记录核心字段）
- `location`、`ttl`、`weight`、`priority`：Optional, Computed（有默认值）
- `status`：Optional, Computed（默认 enable）
- `created_on`、`modified_on`：Computed（只读字段）

**原因**：
- 与 CreateDnsRecord 请求参数对齐
- Computed 属性确保 Read 操作能正确回填云 API 返回的默认值

### 5. 创建返回值校验

**决定**：CreateDnsRecord 调用成功后，检查 response.Response.RecordId 是否为空，若为空则返回 NonRetryableError。

**原因**：遵循代码生成要求，确保创建接口返回有效 ID。

## Risks / Trade-offs

- [Risk] DescribeDnsRecords 使用模糊过滤（id 字段支持模糊匹配），可能返回多条记录 → 只取第一条，并在服务层方法中已做此处理
- [Risk] ModifyDnsRecords 和 ModifyDnsRecordsStatus 不是原子操作，如果第一步成功第二步失败，状态可能不一致 → Terraform 下次 apply 时会重新检测并修正
- [Risk] 新资源名称带有 v7 后缀，可能让用户困惑 → 通过文档说明这是新版本资源
