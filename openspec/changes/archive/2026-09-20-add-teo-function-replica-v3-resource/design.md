## Context

TencentCloud EdgeOne (TEO) 边缘函数副本通过 `zone_id + function_id + replica_name` 三元组唯一标识（Create 接口不返回独立 ID）。云 API（vendor 中 `teo/v20220901`）已提供完整 CRUD 接口，且均为同步接口（无 FlowId/TaskId/JobId 等异步任务标识，无需轮询）。

Provider 中已存在同类资源 `tencentcloud_teo_function_replica`（`tencentcloud/services/teo/resource_tc_teo_function_replica.go`，302 行），本次新增版本化资源 `tencentcloud_teo_function_replica_v3`（文件 `resource_tc_teo_function_replica_v3.go`），不修改既有资源。

代码风格严格参考 `tencentcloud_igtm_strategy`（`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`）。

### Vendor 云 API struct 摘录（实施时无需回读 vendor 大文件）

以下摘录自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`：

**CreateFunctionReplicaRequest**（models.go L3307-3324）:
```go
type CreateFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId      *string `json:"ZoneId,omitnil,omitempty"`      // 站点 ID。（必填）
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID。（必填）
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 副本名称。1-50 字符，允许 a-z、0-9、-，-不能单独/连续使用，不能放开头结尾。同一 FunctionId 下唯一。（必填）
	Content     *string `json:"Content,omitnil,omitempty"`     // 副本内容，仅支持 JavaScript，最大 5MB。（必填）
	Remark      *string `json:"Remark,omitnil,omitempty"`      // 副本描述，最大 50 字符。（可选）
}
```
**CreateFunctionReplicaResponseParams**（L3350-3353）: 仅含 `RequestId *string`，无业务返回字段 → 无独立 ID，用联合 ID。

**DescribeFunctionReplicasRequest**（L10252-10275）:
```go
type DescribeFunctionReplicasRequest struct {
	*tchttp.BaseRequest
	ZoneId     *string `json:"ZoneId,omitnil,omitempty"`     // 站点 ID。
	FunctionId *string `json:"FunctionId,omitnil,omitempty"` // 函数 ID。
	Offset     *int64  `json:"Offset,omitnil,omitempty"`     // 分页偏移量。默认 0。
	Limit      *int64  `json:"Limit,omitnil,omitempty"`      // 分页限制。默认 20，最大 200。
	SortBy     *string `json:"SortBy,omitnil,omitempty"`     // 排序依据，取值：created-on（创建时间）。默认 created-on。
	SortOrder  *string `json:"SortOrder,omitnil,omitempty"`  // 排序方式：asc / desc。默认 asc。
	Filters    []*AdvancedFilter `json:"Filters,omitnil,omitempty"` // 过滤条件，Values 上限 20。支持 replica-name（按副本名称过滤，支持模糊查询）。
}
```
**DescribeFunctionReplicasResponseParams**（L10303-10312）:
```go
type DescribeFunctionReplicasResponseParams struct {
	TotalCount       *int64             `json:"TotalCount,omitnil,omitempty"`       // 副本总数。
	FunctionReplicas []*FunctionReplica `json:"FunctionReplicas,omitnil,omitempty"` // 副本列表。
	RequestId        *string            `json:"RequestId,omitnil,omitempty"`
}
```
**FunctionReplica**（L16998-17016）:
```go
type FunctionReplica struct {
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID。
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 副本名称。
	Content     *string `json:"Content,omitnil,omitempty"`     // 副本内容（JavaScript）。
	Remark      *string `json:"Remark,omitnil,omitempty"`      // 副本描述。
	CreatedOn   *string `json:"CreatedOn,omitnil,omitempty"`   // 创建时间。
	ModifiedOn  *string `json:"ModifiedOn,omitnil,omitempty"`  // 更新时间。
}
```
**AdvancedFilter**（L374-383）:
```go
type AdvancedFilter struct {
	Name   *string  `json:"Name,omitnil,omitempty"`   // 需要过滤的字段。
	Values []*string `json:"Values,omitnil,omitempty"` // 字段的过滤值。
	Fuzzy  *bool    `json:"Fuzzy,omitnil,omitempty"`  // 是否启用模糊查询。
}
```
**ModifyFunctionReplicaRequest**（L20027-20044）:
```go
type ModifyFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId      *string `json:"ZoneId,omitnil,omitempty"`      // 站点 ID。
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID。
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 需要修改的副本名称（定位标识，不可改名）。
	Content     *string `json:"Content,omitnil,omitempty"`     // 副本内容。（可选）
	Remark      *string `json:"Remark,omitnil,omitempty"`      // 副本描述。（可选）
}
```
**ModifyFunctionReplicaResponseParams**（L20070-20073）: 仅含 `RequestId`。

**DeleteFunctionReplicaRequest**（L6827-6838）:
```go
type DeleteFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId       *string   `json:"ZoneId,omitnil,omitempty"`       // 站点 ID。
	FunctionId   *string   `json:"FunctionId,omitnil,omitempty"`   // 函数 ID。
	ReplicaNames []*string `json:"ReplicaNames,omitnil,omitempty"` // 需要删除的副本名称列表。
}
```
**DeleteFunctionReplicaResponseParams**（L6862-6865）: 仅含 `RequestId`。

**Client 方法**（client.go）: `CreateFunctionReplicaWithContext`（L1567）、`DescribeFunctionReplicasWithContext`（L7481）、`ModifyFunctionReplicaWithContext`（L12933）、`DeleteFunctionReplicaWithContext`（L4881）。Client 获取方式：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()`。

### 代码风格参照：igtm_strategy 关键结构

`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（680 行）：
- 资源函数命名：`ResourceTencentCloudIgtmStrategy()`；CRUD 函数命名：`resourceTencentCloudIgtmStrategyCreate/Read/Update/Delete`（小驼峰资源名）
- 骨架：`defer tccommon.LogElapsed("resource.tencentcloud_xx.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`；`logId := tccommon.GetLogId(tccommon.ContextNil)`；`ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`
- request 构建：`igtmv20231024.NewCreateStrategyRequest()`；字符串指针用 `helper.String(v.(string))`
- retry 骨架：`resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { result, e := ...XxxWithContext(ctx, request); if e != nil { return tccommon.RetryError(e) } else { log.Printf("[DEBUG]%s api[%s] success, ...") }; if result == nil || result.Response == nil { return resource.NonRetryableError(fmt.Errorf(...)) }; return nil })`；retry 外处理错误并 `return reqErr`，然后 SetId / 调 Read
- 联合 ID：`d.SetId(strings.Join([]string{a, b, c}, tccommon.FILED_SP))`；读取时 `idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`，段数不符返回 `fmt.Errorf("id is broken, id is %s", d.Id())`
- Update：`needChange := false; mutableArgs := []string{...}; for _, v := range mutableArgs { if d.HasChange(v) { needChange = true; break } }`，needChange 时构建 request 调 Modify，最后 `return XxxRead(d, meta)`
- Read 中设置字段前判断 nil：`if respData.Xxx != nil { _ = d.Set("xxx", respData.Xxx) }`
- Importer：`Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`

**既有同类实现参照**：`tencentcloud/services/teo/resource_tc_teo_function_replica.go`（v1，302 行）——本次 v3 与其 API 相同，schema 在其基础上增加查询/删除相关字段（见 Decisions）。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_function_replica_v3` 资源完整 CRUD（RESOURCE_KIND_GENERAL）
- schema 覆盖接口映射中列出的所有参数：zone_id、function_id、replica_name、content、remark、sort_by、sort_order、filters（含 name/values/fuzzy）、replica_names、created_on、modified_on、function_replicas（出参平铺）
- 使用 `tccommon.FILED_SP`（"#"）联合 ID：`zone_id#function_id#replica_name`，支持 Import
- CRUD 均使用 retry（Read 用 `tccommon.ReadRetryTimeout`，C/U/D 用 `tccommon.WriteRetryTimeout`），错误用 `tccommon.RetryError()` 包装
- 提供基于 gomonkey 的单元测试（mock 云 API，不使用 TF 验收测试套件）
- 在 provider.go / provider.md 注册资源；生成资源 .md 文档

**Non-Goals:**
- 不修改既有 `tencentcloud_teo_function_replica`（v1）资源
- 不实现 data source
- 不处理异步轮询（四个接口均同步，响应无任务 ID）
- 不生成 `_extension.go` 文件（无必要）
- 不在实施阶段执行 `gofmt`/`make doc`/`.changelog` 写入（由收尾阶段 tfpacer-finalize skill 统一处理）
- 不新增 `website/` 目录文件（由 `make doc` 生成）

## Decisions

### D1. Schema 设计（字段来源与必填性）

按接口映射（入参必填性以映射标注为准）：

| Schema 字段 | 类型 | Required/Optional/Computed | ForceNew | 来源接口字段 |
|---|---|---|---|---|
| `zone_id` | String | Required | 是 | Create/Describe/Modify/Delete `ZoneId` |
| `function_id` | String | Required | 是 | Create/Describe/Modify/Delete `FunctionId` |
| `replica_name` | String | Required | 是 | Create/Describe/Modify `ReplicaName` |
| `content` | String | Required | 否 | Create `Content`（必填）；Modify `Content`（可选）→ 以 Create 为准设 Required |
| `remark` | String | Optional | 否 | Create/Modify `Remark` |
| `sort_by` | String | Optional | 否 | Describe `SortBy`（仅查询用，Create/Modify/Delete 不含） |
| `sort_order` | String | Optional | 否 | Describe `SortOrder` |
| `filters` | TypeList → Resource{name(String,Required), values(TypeList→String,Required), fuzzy(Bool,Optional)} | Optional | 否 | Describe `Filters []*AdvancedFilter` |
| `replica_names` | TypeList → String | Required | 否 | Delete `ReplicaNames []*string`（映射标注必填） |
| `created_on` | String | Computed | — | `FunctionReplicas.CreatedOn` |
| `modified_on` | String | Computed | — | `FunctionReplicas.ModifiedOn` |
| `function_replicas` | TypeList → Resource{function_id, replica_name, content, remark, created_on, modified_on}（Computed） | Computed | — | `Response.FunctionReplicas` |

要点：
- **不创建"列表包装层"**（规则 13）：`function_replicas` 直接作为顶层 Computed 列表，元素字段（function_id/replica_name/content/remark/created_on/modified_on）平铺其中，不出现 `xxx_set` 再嵌套一层的结构。
- `filters.values` 用 TypeList→String（与 tccommon 常规做法一致，参照 `l7_acc_rule_priority_operation.rule_ids`）。
- `replica_names` 虽在 Delete 接口为列表入参，但资源粒度为单个副本；删除时取该列表传入 `request.ReplicaNames`（schema 定义 TypeList→String, Required，读取 `d.Get("replica_names").([]interface{})` 逐项 append `helper.String(item.(string))`）。
- `content` 在 Modify 中为可选，但 Create 必填 → schema 设 Required；Update 时 `d.HasChange("content")` 后按 `d.GetOk` 设置。

### D2. 资源 ID 与 CRUD 参数来源

- ID：`zone_id#function_id#replica_name`（3 段，`tccommon.FILED_SP` 分隔）。Create 成功后 `d.SetId(...)` 再调 Read。
- Read/Update/Delete 均从 `d.Id()` 拆段获取 zone_id/function_id/replica_name 作为请求参数（规则 6）。`len(idSplit) != 3` 时返回 `fmt.Errorf("id is broken, id is %s", d.Id())`。
- Import：`schema.ImportStatePassthrough`；.md 文档 Import 部分说明使用联合 ID `zone_id#function_id#replica_name`。

### D3. Read 实现

- 调 `DescribeFunctionReplicasWithContext`，入参：ZoneId、FunctionId、`Limit = helper.Int64(200)`（云 API 注释最大值，规则 5）、Filters 内部固定使用 `Name="replica-name"` + `Values=[replica_name]` 精确定位（参照 v1 实现）；若用户配置了 `sort_by`/`sort_order`/`filters`，则将用户配置映射到请求对应字段（sort_by→SortBy、sort_order→SortOrder、filters→Filters，注意此时需合并/保留 replica-name 过滤以定位单个副本——实现上以内部 replica-name 过滤为准，用户 filters 作为附加过滤条件传入）。
- retry 块内仅调用 API（规则 4），成功后将 response 赋值给外部变量，retry 块外处理。
- 返回空（`response == nil || response.Response == nil` 或匹配不到目标 replica）：**先 `log.Printf("[CRUD] xxx id=%s", d.Id())` 保留现场**（带 `teo_function_replica_v3` 资源名），再 `d.SetId("")`（规则 8）。
- 设置字段前判 nil（`if targetReplica.Content != nil { _ = d.Set("content", ...) }`）；同时 set zone_id/function_id/replica_name（来自 id 拆段）、created_on/modified_on（Computed）。
- `function_replicas` Computed 列表：将 `response.Response.FunctionReplicas` 映射为 `[]map[string]interface{}`（每个元素包含 function_id/replica_name/content/remark/created_on/modified_on，均判 nil 后放入 map），`_ = d.Set("function_replicas", list)`。

### D4. Update 实现

- `mutableArgs := []string{"content", "remark", "sort_by", "sort_order", "filters"}` 中任一 HasChange 时执行 Modify（注意：sort_by/sort_order/filters 为查询参数，Modify 接口无对应字段，不传入 Modify request；仅 content/remark 进入 Modify 请求体）。为避免无意义的 Modify 调用，实际判定变更的字段限定为 `{"content", "remark"}`（与云 API Modify 入参一致），sort/filters 变化不触发云 API 调用（它们不影响服务端状态）。
- request：ZoneId、FunctionId、ReplicaName（id 拆段）+ `d.GetOk("content")` / `d.GetOk("remark")`。
- retry 用 WriteRetryTimeout；成功后 `return Read(d, meta)`。

### D5. Delete 实现

- request：ZoneId、FunctionId（id 拆段）、`ReplicaNames` 取 `d.Get("replica_names")` 列表（`[]*string`）。若用户未配置 `replica_names`（理论上 Required 不会为空，但防御性处理），回退为 `[replica_name]`（id 拆段），保证删除目标正确。
- retry 用 WriteRetryTimeout；成功后直接 return nil（不回读，删除后资源不存在）。

### D6. Create 返回值校验（规则 9）

- retry 块内检查 `result == nil || result.Response == nil` → `resource.NonRetryableError`。
- Create 接口无独立 ID 返回，打印 `log.Printf("[DEBUG]%s create teo function replica v3, id: %s", logId, ...)` 后以联合 ID SetId。本接口为同步接口，无 FlowId/TaskId，**不需要**异步轮询，也无需提前 SetId 的异步场景处理。

### D7. Provider 注册与文档

- `tencentcloud/provider.go`：在 TEO 资源 map 中（`tencentcloud_teo_function_replica` 之后）添加 `"tencentcloud_teo_function_replica_v3": teo.ResourceTencentCloudTeoFunctionReplicaV3(),`（保持列对齐风格）。
- `tencentcloud/provider.md`：在 TEO 部分资源清单中 `tencentcloud_teo_function_replica` 之后添加 `tencentcloud_teo_function_replica_v3`。
- 资源函数命名：`ResourceTencentCloudTeoFunctionReplicaV3`；CRUD 函数命名：`resourceTencentCloudTeoFunctionReplicaV3Create/Read/Update/Delete`；日志资源名统一使用 `teo_function_replica_v3`（规则 12）。
- `resource_tc_teo_function_replica_v3.md`：一句话描述（带 TEO 产品名，"Provides a resource to create a TEO edge function replica"）+ Example Usage + Import（说明联合 ID 格式），不添加 Argument/Attribute Reference。

### D8. 单元测试（gomonkey，规则 go代码生成要求 1）

- 文件 `resource_tc_teo_function_replica_v3_test.go`，包名 `teo_test`。
- 参照既有 `resource_tc_teo_function_replica_test.go` 模式：`mockMeta` 实现 `tccommon.ProviderMeta`（GetAPIV3Conn 返回 `&connectivity.TencentCloudClient{}`）、`patches.ApplyMethodReturn(mockClient, "UseTeoV20220901Client", teoClient)`、`patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(...)...)`、`schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})`。
- 用例覆盖：Create（校验请求字段+ID 拼接）、Read（存在/不存在两种）、Update（Modify 请求字段校验）、Delete（ReplicaNames 传参校验）。
- 每个测试函数头注释 `// go test ./tencentcloud/services/teo/ -run "TestXxx" -v -count=1 -gcflags="all=-l"`。
- 禁止实际执行测试命令，仅保证代码可编译。

## Risks / Trade-offs

- [Risk] DescribeFunctionReplicas 的 replica-name 过滤支持模糊查询 → Mitigation：Read 中遍历 `FunctionReplicas`，精确匹配 `*replica.ReplicaName == replicaName` 后再取值（参照 v1 实现）。
- [Risk] Create 无返回 ID，若 Create 成功但 Read 失败可能出现状态偏差 → Mitigation：Create 成功即 SetId（联合 ID 由输入参数组成，不依赖返回值），Read 失败时保留 id 便于排障/销毁。
- [Risk] `replica_names` schema 为列表，用户可传多项，但资源粒度为单副本 → Mitigation：Delete 直接透传列表（与 API 语义一致）；.md 示例中演示单元素用法。
- [Trade-off] `sort_by`/`sort_order`/`filters` 暴露到资源 schema 但不影响服务端状态，Update 时这些字段变更不触发 Modify 调用 → 符合"资源参数覆盖接口参数"的映射要求，同时避免无效云 API 调用。
- [Trade-off] v3 与 v1 资源功能高度重叠 → 纯新增保持向后兼容，用户可自行选择；不迁移、不废弃 v1。

## Migration Plan

纯新增资源，无迁移需求。回滚策略：删除新增文件与 provider 注册行即可。

## Open Questions

无（云 API 四接口均已在 vendor 中确认，字段映射明确）。
