## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘安全加速平台，提供 DNS 解析管理能力。当前 TEO 的 DNS 记录只能通过控制台或 API 手动管理，需要新增 Terraform 资源 `tencentcloud_teo_dns_record_14` 来支持基础设施即代码管理。

云 API 接口分析：
- `CreateDnsRecord`：创建 DNS 记录，返回 RecordId（同步接口）
- `DescribeDnsRecords`：查询 DNS 记录列表，支持过滤和分页（同步接口）
- `ModifyDnsRecords`：批量修改 DNS 记录，传入 DnsRecord 列表（同步接口）
- `DeleteDnsRecords`：批量删除 DNS 记录，传入 RecordId 列表（同步接口）

所有接口均为同步接口，无需轮询等待。

参考资源：`tencentcloud_igtm_strategy`（通用 CRUD 资源模式）

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_14` 资源的完整 CRUD 生命周期管理
- 支持 DNS 记录的核心属性：zone_id、name、type、content、location、ttl、weight、priority
- 支持 computed 属性：record_id、status、created_on、modified_on
- 支持资源导入（import）
- 在 provider.go 和 provider.md 中正确注册新资源
- 编写单元测试覆盖核心业务逻辑

**Non-Goals:**
- 不实现 DNS 记录的批量管理（每次操作单条记录）
- 不实现 DNS 记录的启停操作（Status 字段为只读出参）
- 不修改已有 TEO 资源的行为

## Decisions

### 1. 资源 ID 组成方式
**决定**: 使用 `zoneId#recordId` 作为复合 ID，使用 `tccommon.FILED_SP`（"#"）作为分隔符。

**理由**: `DescribeDnsRecords` 接口需要 `ZoneId` 作为必填参数来查询记录列表，而 `RecordId` 用于在列表中定位具体记录。两个参数在 Read/Update/Delete 操作中均需使用，因此采用复合 ID 的方式存储。

### 2. type 字段不可变更
**决定**: 将 `type` 字段设为 `ForceNew: true`，修改类型需要销毁重建资源。

**理由**: 虽然 `ModifyDnsRecords` 的 DnsRecord 结构体包含 Type 字段，但 DNS 记录类型的变更在语义上等同于创建新记录。将 type 设为 ForceNew 符合 Terraform 的最佳实践，避免因类型变更导致的解析异常。

### 3. zone_id 字段不可变更
**决定**: 将 `zone_id` 字段设为 `ForceNew: true`。

**理由**: DNS 记录属于特定站点（Zone），变更 zone_id 意味着记录需要在不同站点重新创建，所有 API 调用均需 ZoneId 参数。

### 4. Read 实现方式
**决定**: 调用 `DescribeDnsRecords` 接口，通过 `Filters` 按 `id` 过滤查询特定记录。

**理由**: 不存在按单个 RecordId 查询的接口，只能通过列表查询加过滤的方式获取特定记录。使用 `id` 过滤条件可以在 API 层面缩小查询范围，减少不必要的数据传输。

### 5. Update 实现方式
**决定**: 调用 `ModifyDnsRecords` 接口，将当前资源的所有可变字段打包为 DnsRecord 对象传入。

**理由**: `ModifyDnsRecords` 接口需要传入完整的 DnsRecord 对象（包含 RecordId 和所有可变字段），而非仅传入变更字段。每次 Update 时传入完整记录数据，确保服务端数据与 Terraform state 一致。

### 6. Delete 实现方式
**决定**: 调用 `DeleteDnsRecords` 接口，传入 RecordIds 列表删除单条记录。

**理由**: `DeleteDnsRecords` 接口支持批量删除，但对于 Terraform 资源而言每次只管理单条记录，因此只传入一个 RecordId。

### 7. Computed 字段处理
**决定**: 在 Read 方法中，只有当 Response 中的字段不为 nil 时才调用 d.Set()。

**理由**: 遵循项目规范，避免将 nil 值设置到 state 中导致不必要的 diff。

## Risks / Trade-offs

- **[DescribeDnsRecords 过滤可靠性]** → 使用 `id` 过滤条件精确匹配 recordId，并在代码中验证返回结果中是否包含目标记录，若未找到则标记资源已删除（`d.SetId("")`）
- **[ModifyDnsRecords 批量接口限制]** → 当前只管理单条记录，不会触发批量限制（最大 100 条），但需注意 DnsRecord 结构体中 ZoneId/Status/CreatedOn/ModifiedOn 为只读字段，不应传入
- **[DNS 记录名称解析]** → DNS 记录的 Name 字段对于 CJK 域名需要转换为 punycode，此转换由用户在 Terraform 配置中处理，资源本身不做自动转换
