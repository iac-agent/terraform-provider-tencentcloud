## Context

腾讯云 TEO（EdgeOne）通过 `teo/v20220901` SDK 提供 DNS 记录管理能力。本次变更为 master 分支新增通用资源 `tencentcloud_teo_dns_record_17`（RESOURCE_KIND_GENERAL，全 CRUD 生命周期），核心需求是新增 `DescribeDnsRecords` 接口的入参 `Limit`（分页查询限制数目）到 Terraform 资源 schema 中（SchemaName: `Limit`，Optional）。

现有 master 中同域资源 `tencentcloud_teo_dns_record`（`resource_tc_teo_dns_record.go`）不含 `Limit` 参数，其服务层方法 `DescribeTeoDnsRecordById` 调用 `DescribeDnsRecords` 时未设置 `Limit`（受 API 默认值 20 限制）。本次新增资源不复用该方法，而是新增独立的服务层方法 `DescribeTeoDnsRecord17ById` 以支持 `Limit` 传递。

### 涉及的文件

| 文件 | 作用 | 操作 |
|------|------|------|
| `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go` | 资源 schema 定义及 CRUD 逻辑 | 新增 |
| `tencentcloud/services/teo/resource_tc_teo_dns_record_17_test.go` | gomonkey mock 单元测试 | 新增 |
| `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md` | 资源文档（Example Usage + Import） | 新增 |
| `tencentcloud/services/teo/service_tencentcloud_teo.go` | 服务层，新增 `DescribeTeoDnsRecord17ById`（封装 DescribeDnsRecords 调用，支持 Limit） | 修改 |
| `tencentcloud/provider.go` | ResourcesMap 注册 | 修改 |
| `tencentcloud/provider.md` | TEO 分类下资源名列表 | 修改 |

### 云 API 定义摘录（vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go）

**DescribeDnsRecordsRequest（行 10026-10049）——本次新增参数 `Limit` 所在接口：**
```go
type DescribeDnsRecordsRequest struct {
	*tchttp.BaseRequest

	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// 分页查询偏移量，默认为 0。
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// 分页查询限制数目，默认值：20，上限：1000。
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// 过滤条件，Filters.Values 的上限为20。详细的过滤条件如下：
	// <li>id： 按照 DNS 记录 ID 进行过滤，支持模糊查询；</li>
	// <li>name：按照 DNS 记录名称进行过滤，支持模糊查询；</li>
	// <li>content：按照 DNS 记录内容进行过滤，支持模糊查询；</li>
	// <li>type：按照 DNS 记录类型进行过滤，不支持模糊查询。可选项：A / AAAA / CNAME / TXT / NS / CAA / SRV / MX</li>
	// <li>ttl：按照解析生效时间进行过滤，不支持模糊查询。</li>
	Filters []*AdvancedFilter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// 排序依据，取值有：content / created-on / name / ttl / type。默认根据 type, name 属性组合排序。
	SortBy *string `json:"SortBy,omitnil,omitempty" name:"SortBy"`

	// 列表排序方式，取值有：asc（升序）/ desc（降序）。默认值为 asc。
	SortOrder *string `json:"SortOrder,omitnil,omitempty" name:"SortOrder"`

	// 匹配方式，取值有：all（返回匹配所有查询条件的记录）/ any（返回匹配任意一个查询条件的记录）。默认值为 all。
	Match *string `json:"Match,omitnil,omitempty" name:"Match"`
}
```

> **关键点**：`Limit` 类型为 `*int64`；云 API 注释标注默认值 20、**上限 1000**。`Offset` 默认 0，无需暴露。按需求仅新增 `Limit` 一个参数，`Offset`/`SortBy`/`SortOrder`/`Match` 均不暴露。

**DescribeDnsRecordsResponse（行 10077-10091）：**
```go
type DescribeDnsRecordsResponseParams struct {
	// DNS 记录总数。
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// DNS 记录列表。
	DnsRecords []*DnsRecord `json:"DnsRecords,omitnil,omitempty" name:"DnsRecords"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得该次请求的 RequestId）。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDnsRecordsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDnsRecordsResponseParams `json:"Response"`
}
```

**CreateDnsRecordRequest（行 3155-3185）及 Response（行 3214-3225）：**
```go
type CreateDnsRecordRequest struct {
	*tchttp.BaseRequest

	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// DNS 记录名，如果是中文、韩文、日文域名，需要转换为 punycode 后输入。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// DNS 记录类型：A / AAAA / MX / CNAME / TXT / NS / CAA / SRV。
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`
	// DNS 记录内容，根据 Type 值填入与之相对应的内容。
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
	// DNS 记录解析线路，不指定默认为 Default。仅 Type 为 A、AAAA、CNAME 时适用。
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`
	// 缓存时间，60~86400，默认 300，单位：秒。
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
	// DNS 记录权重，-1~100，设置为 0 表示不解析，默认 -1（不设置）。仅 Type 为 A、AAAA、CNAME 时适用。
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`
	// MX 记录优先级，仅 Type 为 MX 时生效，0~50，默认 0。
	Priority *int64 `json:"Priority,omitnil,omitempty" name:"Priority"`
}

type CreateDnsRecordResponseParams struct {
	// DNS 记录 ID。
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`
	// 唯一请求 ID。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

> **注意**：`CreateDnsRecordRequest` 中**没有** `Limit` 参数——`Limit` 仅属于 `DescribeDnsRecords`（查询接口），故只在 Read 路径使用。

**DnsRecord struct（行 16056-16100）——Read 出参 / Modify 入参元素结构：**
```go
type DnsRecord struct {
	// 站点 ID。注意：ZoneId 仅做出参使用，在 ModifyDnsRecords 不可作为入参使用，如有传此参数，会忽略。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// DNS 记录 ID。
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`
	// DNS 记录名。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// DNS 记录类型：A / AAAA / MX / CNAME / TXT / NS / CAA / SRV。
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`
	// DNS 记录解析线路，不指定默认为 Default。
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`
	// DNS 记录内容。
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
	// 缓存时间，60~86400，单位：秒。
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
	// DNS 记录权重，-1~100，-1 表示不分配权重，0 表示不解析。
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`
	// MX 记录优先级，0~50。
	Priority *int64 `json:"Priority,omitnil,omitempty" name:"Priority"`
	// DNS 记录解析状态：enable（已生效）/ disable（已停用）。注意：Status 仅做出参使用，在 ModifyDnsRecords 中传入会被忽略。
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
	// 创建时间。注意：CreatedOn 仅做出参使用，在 ModifyDnsRecords 中传入会被忽略。
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`
	// 修改时间。注意：ModifiedOn 仅做出参使用，在 ModifyDnsRecords 中传入会被忽略。
	ModifiedOn *string `json:"ModifiedOn,omitnil,omitempty" name:"ModifiedOn"`
}
```

**ModifyDnsRecordsRequest（行 20106-20114）：**
```go
type ModifyDnsRecordsRequest struct {
	*tchttp.BaseRequest

	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// DNS 记录修改数据列表，一次最多修改100条。
	DnsRecords []*DnsRecord `json:"DnsRecords,omitnil,omitempty" name:"DnsRecords"`
}
```

**ModifyDnsRecordsStatusRequest（行 20170-20181）——status 字段更新接口：**
```go
type ModifyDnsRecordsStatusRequest struct {
	*tchttp.BaseRequest

	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// 待启用的 DNS 记录 ID 列表，上限：200。注意：同个 DNS 记录 ID 不能同时存在于 RecordsToEnable 和 RecordsToDisable。
	RecordsToEnable []*string `json:"RecordsToEnable,omitnil,omitempty" name:"RecordsToEnable"`

	// 待停用的 DNS 记录 ID 列表，上限：200。
	RecordsToDisable []*string `json:"RecordsToDisable,omitnil,omitempty" name:"RecordsToDisable"`
}
```

**DeleteDnsRecordsRequest（行 6784-6792）：**
```go
type DeleteDnsRecordsRequest struct {
	*tchttp.BaseRequest

	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// 待删除的 DNS 记录 ID 列表，上限：1000。
	RecordIds []*string `json:"RecordIds,omitnil,omitempty" name:"RecordIds"`
}
```

**AdvancedFilter struct（行 388-397）——Read 按记录 ID 过滤用：**
```go
type AdvancedFilter struct {
	// 需要过滤的字段。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 字段的过滤值。
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
	// 是否启用模糊查询。
	Fuzzy *bool `json:"Fuzzy,omitnil,omitempty" name:"Fuzzy"`
}
```

### 目标资源现有 schema 关联说明

master 中同域资源 `tencentcloud_teo_dns_record` 的 schema（`resource_tc_teo_dns_record.go` 行 26-111）字段为：`zone_id`（Required/ForceNew）、`name`、`type`、`content`（Required）、`location`、`ttl`、`weight`、`priority`、`status`（Optional+Computed）、`record_id`、`created_on`、`modified_on`（Computed）——**不含 `Limit` 参数**。其服务层方法 `DescribeTeoDnsRecordById`（`service_tencentcloud_teo.go` 行 1796-1831）调用 `DescribeDnsRecords` 时仅设置 `ZoneId` 和 `Filters`（`id` 精确过滤），未设置 `Limit`，受 API 默认分页值 20 限制。

本次新增的 `tencentcloud_teo_dns_record_17` 与上述字段结构保持一致（参考 `teo-dns-record-51` 同系列先例），在此基础上**新增 `Limit`（Optional，TypeInt）**，对应 `DescribeDnsRecordsRequest.Limit`（JsonPath: `request.Limit` → SchemaName: `Limit`）。`Limit` 属于查询行为参数，不参与资源 ID、不写入 state（不设 Computed、不在 Read 中回写），仅在 Read 路径读取并传给服务层。

### 同类先例参考

- **limit 参数化先例**：commit b2bdb09ed（`add-teo-l4-proxy-pagination-params`）为 `tencentcloud_teo_l4_proxy` 新增 `limit`（Optional，TypeInt），Read 中 `d.GetOk("limit")` 后经 `helper.IntUint64()` 传给服务层 `DescribeTeoL4ProxyById`，服务层将签名扩展为接受 limit 参数。
- **dns_record 系列新资源先例**：`tencentcloud_teo_dns_record_51`（commit fb35d02405），服务层方法 `DescribeTeoDnsRecord51ById` 固定设置 `request.Offset = helper.IntInt64(0)`、`request.Limit = helper.IntInt64(1000)`（云 API 注释标注的最大值）。本次差异点：`Limit` 由用户可配置（未配置时仍默认 1000）。

## Goals / Non-Goals

**Goals:**
- 新增资源 `tencentcloud_teo_dns_record_17` 的完整 CRUD 生命周期（含 import）
- **新增参数 `Limit`（Optional，TypeInt）**：控制 Read 路径 `DescribeDnsRecords` 的分页查询限制数目；未设置时默认 1000（云 API 上限）
- 保持与既有 `tencentcloud_teo_dns_record` 资源并存，零破坏性

**Non-Goals:**
- 不暴露 `Offset`/`SortBy`/`SortOrder`/`Match` 等其他 `DescribeDnsRecords` 入参（本次需求仅新增 `Limit` 一个参数）
- 不修改既有资源 `tencentcloud_teo_dns_record` 的任何代码与 schema
- 不新增 `_extension.go` 文件
- 不手改 `website/` 目录（由收尾阶段 `make doc` 生成）

## Decisions

### Decision 1: 新增参数 `Limit` 仅用于 Read 路径

**选择**: Schema 中新增 `Limit`（Optional，TypeInt，无 Computed），仅 `resourceTencentCloudTeoDnsRecord17Read` 读取并传递给服务层；Create/Update/Delete 路径不使用（`CreateDnsRecord`/`ModifyDnsRecords`/`ModifyDnsRecordsStatus`/`DeleteDnsRecords` 请求中无此参数）。

**理由**: `Limit` 是 `DescribeDnsRecords`（查询接口）的入参，只在 Read 路径有意义。云 API 各 CRUD 接口入参不一致是常态，已核对 `CreateDnsRecordRequest` 等 struct 中均无 `Limit` 字段。

### Decision 2: `Limit` 未设置时服务层默认 1000（云 API 上限）

**选择**: 服务层方法 `DescribeTeoDnsRecord17ById(ctx, zoneId, recordId string, limit int64)` 接收 limit 值；资源 Read 方法中 `d.GetOk("Limit")` 读取用户配置，未配置时传 1000。

**理由**: 云 API 注释标注 `Limit` 默认值 20、上限 1000。按记录 ID（`Filters: id`）过滤定位单条记录时，若该 zone 下记录总数超过 20，默认分页可能截断结果导致误判"资源不存在"。遵循 repo 规范"若查询接口中有分页字段，则给定值应该是云API接口注释中标注的最大值"，默认 1000。

### Decision 3: `Limit` 字段命名保持需求指定的 PascalCase

**选择**: Schema 字段名使用 `Limit`（与需求映射 `request.Limit → Limit` 一致），而非惯用的 snake_case `limit`。

**理由**: 需求明确指定 SchemaName 为 `Limit`。Terraform Plugin SDK v2 允许任意字符串作为字段名。注意：由于 `Limit` 非 Computed 且 Read 中不回写，不会出现 plan 阶段 diff 问题。

### Decision 4: 服务层新增独立方法 `DescribeTeoDnsRecord17ById`，不改既有方法签名

**选择**: 在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecord17ById(ctx context.Context, zoneId, recordId string, limit int64)`，设置 `request.Offset = helper.IntInt64(0)` 与 `request.Limit = helper.IntInt64(limit)`，按 `id` AdvancedFilter 过滤，retry 块内取 `DnsRecords[0]`。

**理由**: 不修改既有 `DescribeTeoDnsRecordById` 签名可避免影响既有资源 `tencentcloud_teo_dns_record` 的调用方；与 `DescribeTeoDnsRecord51ById` 的同系列资源隔离模式一致。

### Decision 5: Read 空结果处理遵循 repo 规范

**选择**: 服务层返回 nil 时，Read 方法中**先打印** `log.Printf("[CRUD] teo_dns_record_17 id=%s", d.Id())` 保留现场，**再** `d.SetId("")`，并打印 WARN 日志后 return nil。

**理由**: 遵循 repo 规范"先打印日志保留现场再 SetId("")，严禁直接 SetId("") 丢失 id 信息"。

### Decision 6: Create 响应校验 + status 更新路径

**选择**:
- Create 中 retry 成功后检查 `response.Response` 与 `RecordId` 是否为空/空串，若为空打印 logId 后返回 `NonRetryableError`；随后 `d.SetId(zoneId#recordId)` 再调用 Read。
- Update 中分两段：`name/type/content/location/ttl/weight/priority` 变更走 `ModifyDnsRecords`（单元素 `DnsRecords` 数组，仅含 RecordId + 可变字段，排除 ZoneId/Status/CreatedOn/ModifiedOn 出参字段）；`status` 变更走 `ModifyDnsRecordsStatus`（enable → RecordsToEnable，disable → RecordsToDisable）。

**理由**: `ModifyDnsRecords` 的 `DnsRecord` 中 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` 为出参字段（传入会被忽略），不可作为入参；`Limit` 亦不出现在 Modify 入参中。

### Decision 7: 单元测试使用 gomonkey mock

**选择**: 新增资源属于新增 terraform 资源，`resource_tc_teo_dns_record_17_test.go` 使用 gomonkey 对 `teov20220901.Client.CreateDnsRecord`/`DescribeDnsRecords`/`ModifyDnsRecords`/`ModifyDnsRecordsStatus`/`DeleteDnsRecords` 打桩，覆盖 CRUD 业务逻辑及 `Limit` 参数传递（验证用户设置 500 时 request.Limit=500、未设置时默认 1000）。

**理由**: 遵循 repo 规范"新增 terraform 资源的测试使用 mock（gomonkey）方法对云API进行 mock 处理，不使用 terraform 测试套件"。

## Risks / Trade-offs

- [Risk] `Limit` 为查询行为参数，用户在 `terraform plan` 中修改它只会影响下一次 refresh 的查询行为，不会触发云上资源变更。
  → Mitigation: `Limit` 为 Optional 且不在 mutableArgs 检查列表中，Update 方法中忽略其变更（不触发任何云 API 调用）。
- [Risk] PascalCase 字段名 `Limit` 与 repo 中惯用 snake_case 不一致，可能在文档生成时排序略有差异。
  → Mitigation: 需求明确指定 SchemaName 为 `Limit`，保持与需求映射一致；文档中正常列出。
- [Risk] `DescribeDnsRecords` 返回列表而非单条记录。
  → Mitigation: 按 `id` AdvancedFilter 精确过滤并取 `DnsRecords[0]`；空列表视为资源已删除（按 Decision 5 处理）。
- [Risk] 既有未合并分支 `feat/add-teo-dns-record-17` 中存在同名但不含 `Limit` 参数的资源实现。
  → Mitigation: 本提案基于 master 从零生成，以 `Limit` 参数为核心差异点，不依赖该分支代码；注册名与文件名一致不会冲突（master 中不存在）。

## Migration Plan / Compatibility

- 纯新增资源 + 新增 Optional 参数，完全向后兼容，无迁移需求。
- 不影响既有 `tencentcloud_teo_dns_record` 资源及其 state。
