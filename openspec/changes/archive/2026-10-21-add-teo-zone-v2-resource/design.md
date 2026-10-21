# Design: add-teo-zone-v2-resource

## Context

腾讯云 Terraform Provider 的 TEO（EdgeOne）服务目录 `tencentcloud/services/teo/` 已存在旧版站点资源 `tencentcloud_teo_zone`（`resource_tc_teo_zone.go`，基于同一套 `CreateZone`/`DescribeZones`/`ModifyZone`/`DeleteZone` 接口但 schema 陈旧：`type` 必填、缺少 `plan_id`/`allow_duplicates`/`jump_start`/`vanity_name_servers` 等新参数，且未透出 `resources`、`ns_detail`、`cname_detail`、`dns_pod_detail` 等新出参）。受"不能修改已有资源 schema（除新增 Optional 字段外）"的硬约束，本次以纯新增 `tencentcloud_teo_zone_v2` 资源的方式暴露新版站点 CRUD 全量能力。

SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（vendor 已存在，无需更新依赖）。客户端获取方式：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()`（返回 `*teov20220901.Client`；资源文件内 import 别名用 `teov20220901`，参考 `resource_tc_teo_edge_kv_namespace.go` 的写法）。

云API 四个接口均为**同步接口**（CreateZone 响应只含 `ZoneId`/`OwnershipVerification`，无 FlowId/TaskId/JobId），故无需异步轮询逻辑；Delete 为幂等删除（重复删除不报错），删除后无需等待。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_zone_v2` 资源（RESOURCE_KIND_GENERAL），完整实现站点 CRUD。
- Create 覆盖 `CreateZone` 全部入参（type/zone_name/area/plan_id/alias_zone_name/tags/allow_duplicates/jump_start），并将创建响应中的归属权验证信息写入 state。
- Read 覆盖 `DescribeZones` 返回的 Zone 全量字段（平铺到顶层，符合"禁止 xxx_set/xxx_list 包一层"的规则——本资源是单站点资源，非列表型 Describe）。
- Update 覆盖 `ModifyZone` 全部可变参数（type/vanity_name_servers/alias_zone_name/area/zone_name）。
- Delete 调用 `DeleteZone`。
- 单元测试用 gomonkey mock 云API，不使用 terraform 测试套件。
- 注册 provider.go / provider.md，编写 resource_tc_teo_zone_v2.md 文档。

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_zone`、`tencentcloud_teo_zone_setting` 资源。
- 不实现 `ModifyZoneStatus`（paused）、`ModifyZoneWorkMode`（work_mode_infos 的变更）——需求映射中 Update 仅含 `ModifyZone` 一个接口的参数，Update 入参不含 paused/work_mode_infos，故 schema 中这两个字段仅作为 Computed 只读出参（Describe 出参映射包含它们），不做变更逻辑。
- 不实现站点归属权验证（VerifyOwnership）、套餐购买（CreatePlan）、绑定套餐（BindZoneToPlan）等接口。
- 不生成 `_extension.go` 文件。

## Decisions

### D1: Schema 设计（参数映射核对结论）

资源 id：单字段 `ZoneId`（CreateZone 响应的 `zone_id`，形如 `zone-2noz78a8ev6k`），无联合 id，支持 Import（`schema.ImportStatePassthrough`）。RESOURCE_KIND_GENERAL 有 Import 部分。

入参（基于 vendor 中 CreateZoneRequest/ModifyZoneRequest 的实际字段核对，全部存在、类型一致）：

| SchemaName | Type | Required/Optional | ForceNew | 写入接口 | CRUD 可用性 |
|---|---|---|---|---|---|
| `type` | TypeString | Optional（有默认值 partial） | 否（ModifyZone 支持修改） | Create: `request.Type`；Update: `request.Type` | ✓ |
| `zone_name` | TypeString | Optional（noDomainAccess 时留空） | 否（ModifyZone 仅无域名接入切 CNAME 场景有效，作为可变参数） | Create: `request.ZoneName`；Update: `request.ZoneName` | ✓ |
| `area` | TypeString | Optional（默认 overseas） | 否 | Create: `request.Area`；Update: `request.Area` | ✓ |
| `plan_id` | TypeString | Optional | ForceNew（CreateZone 专属，ModifyZone 无此参数） | Create: `request.PlanId` | ✓（仅 Create） |
| `alias_zone_name` | TypeString | Optional | 否 | Create: `request.AliasZoneName`；Update: `request.AliasZoneName` | ✓ |
| `tags` | TypeList of {tag_key: Required, tag_value: Required} | Optional | 否 | Create: `request.Tags []*Tag`；Read 回填 | ✓（CreateZone 原生支持 Tags 入参；DescribeZones 出参含 Tags。注意：与旧版 teo_zone 的 TypeMap+svctag 方式不同，本次映射明确为 `request.Tags` 的 `TagKey/TagValue` 结构，直接用 CreateZone 的 Tags 参数实现） |
| `allow_duplicates` | TypeBool | Optional（默认 false） | ForceNew（CreateZone 专属） | Create: `request.AllowDuplicates` | ✓（仅 Create） |
| `jump_start` | TypeBool | Optional（默认 false） | ForceNew（CreateZone 专属） | Create: `request.JumpStart` | ✓（仅 Create） |
| `vanity_name_servers` | TypeList of {switch: Required(string on/off), servers: Optional(list string)} | Optional（Computed 同时设置，因为 Read 会回填） | 否 | Update: `request.VanityNameServers`；Read 从 Zone.VanityNameServers / NSDetail.VanityNameServers 回填 | ✓（仅 ModifyZone） |

出参（Computed，全部来自 DescribeZones 响应 Zone 结构；Create 后的 Read 会回填）：

顶层平铺：`zone_id`(string)、`status`、`cname_status`、`active_status`、`lock_status`、`paused`(bool)、`created_on`、`modified_on`（均为 TypeString 除 paused 外）。

嵌套结构（TypeList + Computed，Elem 为 Resource）：
- `resources`：`[{id, pay_mode(int), create_time, enable_time, expire_time, status, sv: [{key, value, pack, instance_id, protection_specs}], auto_renew_flag(int), plan_id, area, group, zone_number(int), type}]`
- `ns_detail`：`[{cname_speed_up, is_fake(int), ownership_verification, original_name_servers(list), name_servers(list), vanity_name_servers, vanity_name_servers_ips: [{name, i_pv4}]}]`（仅 type=full 返回）
- `cname_detail`：`[{is_fake(int), ownership_verification}]`（仅 type=partial 返回）
- `dns_pod_detail`：`[{is_fake(int)}]`（仅 type=dnsPodAccess 返回）
- `ownership_verification`：`[{dns_verification: [{subdomain, record_type, record_value}], file_verification: [{path, content}], ns_verification: [{name_servers(list)}]}]`（Create 响应与 Zone 出参均有此结构）
- `work_mode_infos`：`[{config_group_type, work_mode}]`
- `tags`（入参同时也是出参，Read 时回填）
- `name_servers`、`original_name_servers`（Zone 顶层的历史保留字段，TypeList of string，Computed）
- `vanity_name_servers_ips`（Zone 顶层，`[{name, i_pv4}]`，Computed）

字段命名注意：需求映射中 `IPv4` → SchemaName `i_pv4`（schema 不允许以数字结尾后接大写驼峰转换歧义，按映射表用 `i_pv4`）。

### D2: CRUD 流程（严格遵循 igtm_strategy / edge_kv_namespace 风格）

**Create**（`resourceTencentCloudTeoZoneV2Create`）：
1. `defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`。
2. 组装 request（`d.GetOk` 逐字段），tags 为 `[]*teov20220901.Tag{TagKey, TagValue}` 列表。
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `CreateZoneWithContext(ctx, request)`；失败 `return tccommon.RetryError(e)`；成功打 `[DEBUG]` 日志。
4. retry 外（错误处理后）：检查 `response == nil || response.Response == nil || response.Response.ZoneId == nil || *response.Response.ZoneId == ""`，为空则 `return fmt.Errorf(...)`（NonRetryableError 在 retry 内做 response nil 检查；ZoneId 检查在 retry 外）。检查前打印 `log.Printf("[DEBUG]%s create teo_zone_v2, logId: %s, id: %s", logId, logId, d.Id())`。
5. `d.SetId(*response.Response.ZoneId)`（同步接口，拿到 id 后立即 SetId，随后调 Read 回填字段）。
6. `return resourceTencentCloudTeoZoneV2Read(d, meta)`。

**Read**（`resourceTencentCloudTeoZoneV2Read`）：
1. `zoneId := d.Id()`。
2. 组装 `DescribeZonesRequest`：`Filters = []*teov20220901.AdvancedFilter{{Name: "zone-id", Values: []*string{&zoneId}}}`，`Limit = helper.Int64(100)`（云API注释标注的最大值）、`Offset = helper.Int64(0)`。
3. `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `DescribeZonesWithContext(ctx, request)`；失败 `tccommon.RetryError(e)`。
4. retry 外：`if response == nil || response.Response == nil || len(response.Response.Zones) == 0` → `log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())` 后 `d.SetId("")`、`return nil`（先打印后清空，保留现场）。
5. 取 `zone := response.Response.Zones[0]`，逐字段判 nil 后 `_ = d.Set(...)`；嵌套结构按 D1 组装 `[]map[string]interface{}` 后 Set（参考 igtm_strategy Read 的组装方式）。
6. tags 回填：从 `zone.Tags`（`[]*teov20220901.Tag`）组装 `[]map[string]interface{}{{"tag_key": ..., "tag_value": ...}}` 后 Set（不走 svctag.DescribeResourceTags，因本资源 tags 为 CreateZone 原生参数结构）。

**Update**（`resourceTencentCloudTeoZoneV2Update`）：
1. `zoneId := d.Id()`；`needChange` 检查 mutableArgs：`["type", "zone_name", "area", "alias_zone_name", "vanity_name_servers"]`（`d.HasChange` 任一命中）。
2. 命中则组装 `ModifyZoneRequest`（`request.ZoneId = &zoneId`；其余字段 `d.GetOk`；vanity_name_servers 组装 `teov20220901.VanityNameServers{Switch, Servers}`）。
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `ModifyZoneWithContext(ctx, request)`；retry 内检查 `result == nil || result.Response == nil` 返回 NonRetryableError。
4. 最后 `return resourceTencentCloudTeoZoneV2Read(d, meta)`。

**Delete**（`resourceTencentCloudTeoZoneV2Delete`）：
1. `zoneId := d.Id()`；`request.ZoneId = &zoneId`。
2. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `DeleteZoneWithContext(ctx, request)`；retry 内检查 response nil。
3. 成功后直接 `return nil`（同步删除，无轮询）。

注意：本资源四个接口均使用 `XxxWithContext` 版本并传 ctx（`tccommon.NewResourceLifeCycleHandleFuncContext` 创建）；旧版 service 层 `DescribeTeoZoneById` 使用不带 ctx 的 `DescribeZones`，本次 v2 不复用该 service 方法，直接在资源文件内实现 Read 调用（与 edge_kv_namespace 一致，避免引入 service 层重复方法）。

### D3: 单元测试（gomonkey，不用 terraform 测试套件）

参考 `resource_tc_teo_edge_kv_namespace_test.go` 的完整结构：
- `package teo_test`，定义 `mockMetaTeoZoneV2`（实现 `GetAPIV3Conn() *connectivity.TencentCloudClient`，`var _ tccommon.ProviderMeta = &mockMetaTeoZoneV2{}` 断言）。
- `patches.ApplyMethodReturn(mockMeta.client, "UseTeoV20220901Client", teoClient)`。
- `patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(...))` mock 各接口，返回 `teov20220901.NewCreateZoneResponse()` + `Response: &teov20220901.CreateZoneResponseParams{...}`。
- 用 `schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` 构造 d，断言 Create/Read/Update/Delete 的请求参数与 state 回填。
- 用例覆盖：Create 正常（含 tags/plan_id/allow_duplicates/jump_start 参数传递断言、ownership_verification 回填断言）、Create 响应 ZoneId 为空报错、Read 正常、Read NotFound 清空 id、Read 空 Response 清空 id、Update 修改 type/area/alias_zone_name/vanity_name_servers、Delete 正常。

### D4: Provider 注册与文档

- `provider.go`：在 `"tencentcloud_teo_zone"` 行附近追加 `"tencentcloud_teo_zone_v2": teo.ResourceTencentCloudTeoZoneV2(),`（保持资源注册 map 的对齐格式，参考第 2126 行附近的写法）。
- `provider.md`：在 "TencentCloud EdgeOne(TEO)" 的 Resource 列表（`tencentcloud_teo_zone` 行后）追加 `tencentcloud_teo_zone_v2`。
- `resource_tc_teo_zone_v2.md`：一句话描述 "Provides a resource to create a TEO zone"（带云产品名 TEO/EdgeOne）；Example Usage 含完整 HCL（tags 块、vanity_name_servers 块示例）；Import 部分说明使用 zone id 导入（单 id，无联合 id）；**不写** Argument Reference / Attribute Reference（工具自动生成）。

## Risks / Trade-offs

- [Risk] `tags` 参数语义与旧版 teo_zone（TypeMap + svctag 通用标签服务）不同，用户可能混淆 → [Mitigation] 在 schema Description 与 .md 文档中明确说明 tags 为 CreateZone 原生标签参数（tag_key/tag_value 结构），需先在标签控制台创建标签。
- [Risk] `zone_name` 不可变语义弱化（ModifyZone 仅特定场景支持修改 zone_name）→ [Mitigation] 遵循需求映射将 zone_name 作为可变参数传给 ModifyZone；云API会在不支持的场景自行报错， terraform 侧不做额外拦截。
- [Risk] Read 中 `ownership_verification` 等字段可能返回 null（vendor 注释标注"此字段可能返回 null"）→ [Mitigation] 所有 Set 前先判 nil；顶层 Computed list 字段在 nil 时 Set 空列表或跳过（按需，参考旧版对 ownership_verification 的 else 分支处理）。
- [Risk] 不同接入类型返回的详情结构不同（ns_detail 仅 full、cname_detail 仅 partial、dns_pod_detail 仅 dnsPodAccess）→ [Mitigation] 三者均定义为 Computed list，Read 时判 nil 跳过，不做类型强校验。
- [Trade-off] 不复用旧版 service 层 `DescribeTeoZoneById`（其使用不带 ctx 的调用与 limit=20），在 v2 资源内直接实现 Read，代码略有重复但获得 ctx 传递与 limit=100 的正确行为。

## Migration Plan

纯新增资源，无迁移。发布后用户可通过 `terraform import` 将已有站点导入 `tencentcloud_teo_zone_v2` 管理（import id 为 ZoneId）。回滚：从 provider 注册中移除该资源即可，无 state 损坏风险（旧资源不受影响）。

## Open Questions

（无。参数映射已与 vendor 云API定义逐项核对一致：CreateZoneRequest 9 个入参、ModifyZoneRequest 6 个入参、DeleteZoneRequest 1 个入参均存在；DescribeZonesRequest 的 Filters/Order/Direction 为查询入参，本资源 Read 仅使用 zone-id 过滤，Order/Direction 不暴露到 schema（资源 Read 无需排序）。）

---

## 附录 A: 云API request/response struct 关键字段（摘自 vendor，实施时无需回读大文件）

文件：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`

### CreateZoneRequest（行 5911-5949）
```go
type CreateZoneRequest struct {
	*tchttp.BaseRequest
	Type *string            // 站点接入类型: partial(CNAME)/full(NS)/noDomainAccess/dnsPodAccess/ai，不填默认 partial
	ZoneName *string        // CNAME/NS 接入传二级域名(example.com)；无域名接入保留为空
	Area *string            // Type=partial/full 时的加速区域: global/mainland/overseas，不填默认 overseas；noDomainAccess 时保留为空
	PlanId *string          // 待绑定的目标套餐 ID；不填则创建 init 状态未激活站点
	AliasZoneName *string   // 同名站点标识，数字/英文/. /-/_，200 字符以内
	Tags []*Tag             // 标签（需先在标签控制台创建）
	AllowDuplicates *bool   // 是否允许重复接入，不填默认 false
	JumpStart *bool         // 是否跳过站点现有 DNS 记录扫描，默认 false
}
```

### CreateZoneResponse（行 5977-5999）
```go
type CreateZoneResponseParams struct {
	ZoneId *string                                  // 站点 ID
	OwnershipVerification *OwnershipVerification    // 归属权验证信息（可能为 null；partial 时需 TXT/文件验证；full 时需切 DNS；noDomainAccess 时为空）
	RequestId *string
}
// CreateZoneResponse{ *tchttp.BaseResponse; Response *CreateZoneResponseParams }
```

### DescribeZonesRequest（行 15712-15737）
```go
type DescribeZonesRequest struct {
	*tchttp.BaseRequest
	Offset *int64                    // 默认 0
	Limit *int64                     // 默认 20，最大 100 ← Read 时取 100
	Filters []*AdvancedFilter        // 过滤条件（Values 上限 20）：zone-name/zone-type/zone-id/status/tag-key/tag-value/alias-zone-name；模糊查询仅支持 zone-name 或 alias-zone-name
	Order *string                    // type/area/create-time/zone-name/use-time/active-status
	Direction *string                // asc/desc，默认 desc
}
```

### DescribeZonesResponse（行 15762-15777）
```go
type DescribeZonesResponseParams struct {
	TotalCount *int64
	Zones []*Zone
	RequestId *string
}
```

### ModifyZoneRequest（行 23092-23118）
```go
type ModifyZoneRequest struct {
	*tchttp.BaseRequest
	ZoneId *string                       // 站点 ID（必填）
	Type *string                         // full/partial/dnsPodAccess，不填保持原有配置
	VanityNameServers *VanityNameServers // 自定义 NS，不填保持原有配置；无域名接入时不允许传
	AliasZoneName *string                // 同名站点标识
	Area *string                         // global/mainland/overseas；无域名接入时不允许传
	ZoneName *string                     // 仅无域名接入切 CNAME 场景有效
}
// ModifyZoneResponseParams{ RequestId *string } —— 无业务出参
```

### DeleteZoneRequest（行 8205-8210）
```go
type DeleteZoneRequest struct {
	*tchttp.BaseRequest
	ZoneId *string   // 站点 ID（必填）
}
// DeleteZoneResponseParams{ RequestId *string } —— 无业务出参
```

### Zone（行 26782-26883，DescribeZones 出参元素）
```go
type Zone struct {
	ZoneId *string
	ZoneName *string
	AliasZoneName *string
	Area *string
	Type *string            // full/partial/noDomainAccess/dnsPodAccess/pages/ai
	Tags []*Tag             // Tag{TagKey, TagValue}
	Resources []*Resource   // 计费资源列表
	NSDetail *NSDetail      // 仅 Type=full 时返回
	CNAMEDetail *CNAMEDetail // 仅 Type=partial 时返回
	DNSPodDetail *DNSPodDetail // 仅 Type=dnsPodAccess 时返回
	CreatedOn *string
	ModifiedOn *string
	Status *string          // active/pending/moved/deactivated/initializing
	CnameStatus *string     // finished/pending
	ActiveStatus *string    // active/inactive/paused
	LockStatus *string      // enable/disable/plan_migrate
	Paused *bool
	IsFake *int64
	CnameSpeedUp *string
	OwnershipVerification *OwnershipVerification  // 历史保留字段，可能 null
	OriginalNameServers []*string                 // 历史保留字段
	NameServers []*string                         // 历史保留字段
	VanityNameServers *VanityNameServers          // 历史保留字段，可能 null
	VanityNameServersIps []*VanityNameServersIps  // 历史保留字段，可能 null
	WorkModeInfos []*ConfigGroupWorkModeInfo
}
```

### 子结构体
```go
type Tag struct { TagKey *string; TagValue *string }                                // 行 26241
type AdvancedFilter struct { Name *string; Values []*string; Fuzzy *bool }          // 行 388
type VanityNameServers struct { Switch *string /* on/off */; Servers []*string }    // 行 26577
type VanityNameServersIps struct { Name *string; IPv4 *string }                     // 行 26587
type OwnershipVerification struct {                                                // 行 24327
	DnsVerification *DnsVerification
	FileVerification *FileVerification
	NsVerification *NsVerification
}
type DnsVerification struct { Subdomain, RecordType, RecordValue *string }          // 行 16102
type FileVerification struct { Path, Content *string }                              // 行 17108
type NsVerification struct { NameServers []*string }                                // 行 23722
type Resource struct {                                                              // 行 25055
	Id *string
	PayMode *int64
	CreateTime *string
	EnableTime *string
	ExpireTime *string
	Status *string
	Sv []*Sv
	AutoRenewFlag *int64
	PlanId *string
	Area *string
	Group *string        // 可能 null
	ZoneNumber *int64    // 可能 null
	Type *string         // 可能 null
}
type Sv struct { Key, Value, Pack, InstanceId, ProtectionSpecs *string }             // 行 26177（Pack/InstanceId/ProtectionSpecs 可能 null）
type NSDetail struct {                                                              // 行 23646
	CnameSpeedUp *string
	IsFake *int64
	OwnershipVerification *OwnershipVerification
	OriginalNameServers []*string
	NameServers []*string
	VanityNameServers *VanityNameServers
	VanityNameServersIps []*VanityNameServersIps
}
type CNAMEDetail struct { IsFake *int64; OwnershipVerification *OwnershipVerification } // 行 1437
type DNSPodDetail struct { IsFake *int64 }                                          // 行 6345
type ConfigGroupWorkModeInfo struct { ConfigGroupType, WorkMode *string }            // 行 2128
```

### client 方法（vendor client.go，均为同步方法，另有对应 XxxWithContext 版本）
```go
func (c *Client) CreateZone(request *CreateZoneRequest) (response *CreateZoneResponse, err error)     // 行 4331
func (c *Client) DeleteZone(request *DeleteZoneRequest) (response *DeleteZoneResponse, err error)     // 行 6222
func (c *Client) DescribeZones(request *DescribeZonesRequest) (response *DescribeZonesResponse, err error) // 行 11248
func (c *Client) ModifyZone(request *ModifyZoneRequest) (response *ModifyZoneResponse, err error)     // 行 16313
```

## 附录 B: 代码风格参照（tencentcloud_igtm_strategy / teo_edge_kv_namespace）

参照文件：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（RESOURCE_KIND_GENERAL 标准风格）、`tencentcloud/services/teo/resource_tc_teo_edge_kv_namespace.go`（teo v20220901 新式风格 + gomonkey 单测配套）。

资源骨架（igtm_strategy）：
```go
func ResourceTencentCloudTeoZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoZoneV2Create,
		Read:   resourceTencentCloudTeoZoneV2Read,
		Update: resourceTencentCloudTeoZoneV2Update,
		Delete: resourceTencentCloudTeoZoneV2Delete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{ ... },   // 嵌套 list 用 Elem: &schema.Resource{Schema: ...}
	}
}
```

CRUD 骨架（每个函数）：
```go
defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.create")()
defer tccommon.InconsistentCheck(d, meta)()

var (
	logId = tccommon.GetLogId(tccommon.ContextNil)
	ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	request = teov20220901.NewCreateZoneRequest()
)
// d.GetOk 逐字段组装 request
reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
	result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateZoneWithContext(ctx, request)
	if e != nil {
		return tccommon.RetryError(e)          // 云API失败统一用 tccommon.RetryError 包装
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	if result == nil || result.Response == nil {
		return resource.NonRetryableError(fmt.Errorf("Create teo_zone_v2 failed, Response is nil"))
	}
	response = result
	return nil
})
if reqErr != nil {
	log.Printf("[CRITAL]%s create teo_zone_v2 failed, reason:%+v", logId, reqErr)
	return reqErr
}
// 成功操作（SetId 等）放在 retry 块外
```

关键风格点（rules 沉淀）：
- 嵌套 list 的 Create/Update 组装：`for _, item := range v.([]interface{})` → `item.(map[string]interface{})` → SDK struct → append；Read 反向组装 `[]map[string]interface{}` → `_ = d.Set(name, list)`。
- Read 中 Response 字段为 nil 时不调用 setXX；资源不存在时先 `log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())` 再 `d.SetId("")`。
- 日志/错误信息统一使用资源名 `teo_zone_v2`（小写蛇形），如 `"[CRITAL]%s create teo_zone_v2 failed, reason:%+v"`。
- Update 的 mutableArgs 模式：`needChange := false; mutableArgs := []string{...}; for _, v := range mutableArgs { if d.HasChange(v) { needChange = true; break } }`。
- helper 工具：`helper.String(v.(string))`、`helper.Bool(v.(bool))`、`helper.Int64(n)`。
- import 分组：标准库 → terraform-plugin-sdk → tencentcloud-sdk-go（别名 `teov20220901`）→ tccommon/helper。

gomonkey 测试骨架（edge_kv_namespace_test.go）：
```go
package teo_test

type mockMetaTeoZoneV2 struct{ client *connectivity.TencentCloudClient }
func (m *mockMetaTeoZoneV2) GetAPIV3Conn() *connectivity.TencentCloudClient { return m.client }
var _ tccommon.ProviderMeta = &mockMetaTeoZoneV2{}

patches := gomonkey.NewPatches()
defer patches.Reset()
teoClient := &teov20220901.Client{}
patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)
patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
	resp := teov20220901.NewCreateZoneResponse()
	resp.Response = &teov20220901.CreateZoneResponseParams{ ZoneId: ptrString("zone-test123"), ... }
	return resp, nil
})
// schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...}) → res.Create(d, meta) → assert
```

provider 注册（provider.go 第 2126 行附近格式）：
```go
"tencentcloud_teo_zone_v2":    teo.ResourceTencentCloudTeoZoneV2(),
```
