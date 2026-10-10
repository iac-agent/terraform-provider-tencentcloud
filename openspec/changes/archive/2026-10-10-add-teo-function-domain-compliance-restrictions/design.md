## Context

`tencentcloud_teo_function` 是 RESOURCE_KIND_GENERAL 资源，封装 TEO 边缘函数的完整生命周期：Create 调 `CreateFunction`，Read 通过 `TeoService.DescribeTeoFunctionById`（内部调 `DescribeFunctions`，按 `ZoneId`+`FunctionIds` 过滤并取 `Functions[0]`）回填 state，Update 调 `ModifyFunction`（可变参数仅 `remark`/`content`），Delete 调 `DeleteFunction`。资源 ID 为 `zoneId#functionId`（`tccommon.FILED_SP` 分隔）。

vendored SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）中 `Function` 结构体已新增 `DomainComplianceRestrictions` 出参，本设计将其中每个元素（`ComplianceRestriction`）的 `Reason`、`Region` 两个字段映射到资源 schema 顶层新增的 `domain_compliance_restrictions` 嵌套 Computed 列表中。

云 API 出参摘录（后续实施只需阅读本节，无需回读 vendor 大文件）：

### vendor SDK 关键定义（models.go，已核对）

`Function` struct（models.go:17226-17253，节选与本次相关字段）：

```go
type Function struct {
	// <p>函数 ID。</p>
	FunctionId *string `json:"FunctionId,omitnil,omitempty" name:"FunctionId"`
	// <p>站点 ID。</p>
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// <p>函数名字。</p>
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// <p>函数描述。</p>
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
	// <p>函数内容。</p>
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
	// <p>函数默认域名。</p>
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`
	// <p>边缘函数默认域名因合规问题产生的地区访问限制列表。</p>
	DomainComplianceRestrictions []*ComplianceRestriction `json:"DomainComplianceRestrictions,omitnil,omitempty" name:"DomainComplianceRestrictions"`
	// <p>创建时间。时间为世界标准时间（UTC）， 遵循 ISO 8601 标准的日期和时间格式。</p>
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`
	// <p>修改时间。时间为世界标准时间（UTC）， 遵循 ISO 8601 标准的日期和时间格式。</p>
	UpdateTime *string `json:"UpdateTime,omitnil,omitempty" name:"UpdateTime"`
}
```

`ComplianceRestriction` struct（models.go:1999-2005，完整定义）：

```go
type ComplianceRestriction struct {
	// <p>下发访问限制的原因。</p><p>枚举值：</p><ul><li>ICP_RECORD_REQUIRED： 未备案；</li><li>GOVERNMENT_ORDER： 政府指令。</li></ul>
	Reason *string `json:"Reason,omitnil,omitempty" name:"Reason"`

	// <p>限制访问地区的具体国家/地区码，使用"ISO 3166 国家/地区代码标准"。</p><p>参数格式：查看链接：https://www.iso.org/iso-3166-country-codes.html。</p>
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}
```

`DescribeFunctionsResponseParams`（models.go:10668-10677）：

```go
type DescribeFunctionsResponseParams struct {
	// 符合查询条件的函数总数。
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	// 符合查询条件的所有函数信息。
	Functions []*Function `json:"Functions,omitnil,omitempty" name:"Functions"`
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

注意：`CreateFunction`、`ModifyFunction` 请求结构中**均无** `DomainComplianceRestrictions` 相关入参（该字段为纯出参、由云端合规系统生成，用户不可配置），因此本次变更只涉及 Schema + Read 回填。

### 现有 schema 与 Read 关联结构

`resource_tc_teo_function.go` 的 Schema map（节选与新增参数相邻的字段）：

```go
"domain": {
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The default domain name for the function.",
},
// ← 新增 "domain_compliance_restrictions"（TypeList + Elem: &schema.Resource{...}）插入于此处
"create_time": {
	Type:        schema.TypeString,
	Computed:    true,
	Description: "Creation time. ...",
},
```

Read 方法当前回填模式（`resourceTencentCloudTeoFunctionRead`，节选）：

```go
respData, err := service.DescribeTeoFunctionById(ctx, zoneId, functionId)
if err != nil {
	return err
}
if respData == nil {
	d.SetId("")
	log.Printf("[WARN]%s resource `teo_function` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
	return nil
}
...
if respData.Domain != nil {
	_ = d.Set("domain", respData.Domain)
}
// ← 新增 domain_compliance_restrictions 回填逻辑插入于此处（UpdateTime 回填之前）
```

服务层 `DescribeTeoFunctionById`（`service_tencentcloud_teo.go:1367-1403`）已按 `resource.Retry(tccommon.ReadRetryTimeout, ...)` + `tccommon.RetryError(e)` 模式封装，`response.Response.Functions[0]` 即 `respData`，无需改动。

## Goals / Non-Goals

**Goals:**
- 用户可通过 `terraform state` / output 读取 `tencentcloud_teo_function.domain_compliance_restrictions`，获知函数默认域名在哪些地区（`region`）因何种原因（`reason`）被限制访问
- 向后兼容：纯新增 Computed 字段，已有配置升级后无 plan drift
- Read 回填遵循"nil 判断后再 set"的既有规范

**Non-Goals:**
- 不为 `DomainComplianceRestrictions` 提供任何写入能力（Create/Update 均不涉及，云 API 无对应入参）
- 不修改现有任何字段的 schema、Create/Update/Delete 逻辑、资源 ID 结构
- 不新增独立数据源（如 `tencentcloud_teo_functions` 列表数据源）
- 不修改 `provider.go` / `provider.md`（资源早已注册）
- 不在本阶段生成 website 文档与 changelog（由收尾阶段 `make doc` / tfpacer-finalize 统一处理）

## Decisions

### Decision 1: 参数命名与结构——单个嵌套 Computed 列表 `domain_compliance_restrictions`

**选择**：在资源 schema 顶层新增 `domain_compliance_restrictions`（`TypeList` + `Computed: true` + `Elem: &schema.Resource{ Schema: { reason, region } }`），把 `Reason`/`Region` 作为子字段平铺在元素结构内。

**备选**：
- a) 在顶层新增两个独立标量 `reason`/`region`：不可行——`DomainComplianceRestrictions` 是列表，同一函数可能在多个地区被限制且原因不同，两个顶层标量无法表达多元素。
- b) 再包一层 `xxx_set`/`xxx_list` 容器：被项目规范明令禁止（禁止"该资源列表型数据"这一层 schema）。

**理由**：JsonPath `response.Functions.DomainComplianceRestrictions.Reason` 对应 SchemaName `Reason`，即列表元素内的子字段；provider 中 Computed 嵌套列表的标准写法即 `TypeList + Elem: &schema.Resource{...}`（参考 `resource_tc_teo_identify_zone_operation.go` 的 `ascription` 块）。

### Decision 2: Read 回填采用"逐元素构造 map + 整体 Set"模式

**选择**：在 Read 中当 `respData.DomainComplianceRestrictions != nil && len(...) > 0` 时，遍历构造 `[]map[string]interface{}`（每个元素仅在对应子字段非 nil 时写入 map），再 `_ = d.Set("domain_compliance_restrictions", list)`。

**备选**：直接 `d.Set("domain_compliance_restrictions", respData.DomainComplianceRestrictions)`（依赖 SDK 结构体到 map 的隐式转换）。

**理由**：
- 与 `resource_tc_teo_zone.go` Read 中 `work_mode_infos` 的回填模式完全一致（项目既有惯例）
- 显式转换避免 SDK 结构体直接 set 进 state 导致的序列化问题；`d.Set` 的 error 用 `_ =` 忽略与其他字段回填保持一致
- 空列表时不调用 Set，state 保持空 `[]`，避免 SDK 返回 `null` 时写脏 state

### Decision 3: 不加入 Update/immutableArgs/mutableArgs

**选择**：新字段为纯 Computed，不出现在 `mutableArgs` 或 `immutableArgs` 中，Update 函数零改动。

**理由**：云 API `CreateFunction`/`ModifyFunction` 均无该入参（纯出参）；Computed 字段 Terraform 不允许用户配置，天然不会触发 `d.HasChange` 路径。

### Decision 4: 单测采用 gomonkey mock（不使用 TF 测试套件）

**选择**：在 `resource_tc_teo_function_test.go` 补充用例，gomonkey mock `TeoService.DescribeTeoFunctionById`（或底层 `DescribeFunctions`）返回带/不带 `DomainComplianceRestrictions` 的 `Function`，断言 state 中回填结果。

**理由**：项目规范要求"修改现有terraform资源的需求使用测试套件补充测试用例"，但同时允许对纯业务逻辑使用 mock；现有 `resource_tc_teo_function_test.go` 已经存在，本次按其既有 mock 风格补充，确保不引入真实云调用、可离线构建。

## Risks / Trade-offs

- [SDK 返回 `DomainComplianceRestrictions` 为 null（正常场景，无限制）] → 回填逻辑加 `!= nil && len > 0` 判断，null 时跳过 Set，state 为空列表，无报错
- [云端合规状态变化导致 plan 显示差异] → Computed 字段不参与 plan diff（terraform 对 computed-only 字段不产生变更提示），风险可控
- [函数默认域名被限制时 Create 的 StateRefresh 模板] → 现有 `resourceTeoFunctionCreateStateRefreshFunc_0_0` 仅判断 `.Domain` 非空，与 `DomainComplianceRestrictions` 无关，不受影响
- [子字段命名与 TF 惯例冲突] → `reason`/`region` 已是 snake_case，无需转换；与 JsonPath 映射（`Reason`→`Reason`）语义一致

## Migration Plan

无迁移需求：纯新增 Computed 字段，`terraform plan` 对旧配置零差异。回滚只需 revert 本变更涉及的三个文件即可，state 中多出的字段在降级后会被自动忽略。
