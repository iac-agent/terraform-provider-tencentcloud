## Context

腾讯云 TEO 边缘函数副本（Edge Function Replica）是挂在某个边缘函数（Function）下的代码副本，每个副本在同一个 `FunctionId` 下通过 `replica_name` 唯一标识——云 API **没有**为副本分配独立的资源 ID，`DeleteFunctionReplica` 按名称列表删除。因此 Terraform 资源必须采用复合 ID：`zone_id#function_id#replica_name`（分隔符 `tccommon.FILED_SP`，即 `#`）。

当前状态：
- 旧资源 `tencentcloud_teo_function_replica`（`tencentcloud/services/teo/resource_tc_teo_function_replica.go`）仅支持 `zone_id` / `function_id` / `replica_name` / `content` / `remark` 五个参数，无 `sort_by` / `sort_order` / `filters` / `replica_names` / `created_on` / `modified_on`。
- 本次需求为新增独立资源 `tencentcloud_teo_function_replica_v5`（RESOURCE_KIND_GENERAL），不动旧资源（旧 spec `teo-function-replica` 不受影响）。
- 四个云 API 全部位于 vendor 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（client 方法均已存在，models.go 中 struct 定义完整），无需变更 vendor。
- 代码风格参照 `tencentcloud_igtm_strategy`（`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`）。

约束：
- 同步接口：`CreateFunctionReplica` / `ModifyFunctionReplica` / `DeleteFunctionReplica` 的响应仅含 `RequestId`，无 FlowId/TaskId/JobId，属于**同步**接口，无需异步轮询等待（Create 成功后直接 SetId 并调用 Read 回填即可）。
- Read 接口 `DescribeFunctionReplicas` 为列表查询（分页 `Offset` / `Limit`，`Limit` 注释标注"默认值：20，最大值：200"），资源 Read 必须以 `Limit=200` 查询并按 `replica_name` 精确定位目标副本。
- PROVIDER 规范：Read 中 set 前需判 nil；查询接口 `Limit` 取注释标注的最大值 200；错误用 `tccommon.RetryError()` 包装；`FILED_SP` 复合 ID；Retry 块内只做 API 调用，SetId 等成功操作放在 Retry 块外。
- 使用 mock（gomonkey）方式编写单元测试，不使用 terraform 测试套件。
- 不生成 `_extension.go` 文件；go 文件开头不加注释；日志统一使用资源蛇形名 `teo function replica v5` / `tencentcloud_teo_function_replica_v5`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_replica_v5` 资源，完整对齐四个云 API 的参数能力（创建/查询/更新/删除侧参数全部暴露）
- Read 回填 `created_on` / `modified_on` Computed 属性
- 支持复合 ID 的 `terraform import`
- provider.go / provider.md 注册与文档、gomonkey 单元测试齐备
- design.md 沉淀 vendor struct 关键字段定义与参考代码结构，实施阶段无需回读 vendor 大文件

**Non-Goals:**
- 不修改旧资源 `tencentcloud_teo_function_replica` 及其测试、文档
- 不为 Describe 的 `TotalCount`、`Offset` 暴露 schema 参数（非资源参数映射范围）
- 不引入 `function_replicas` 嵌套列表层（遵循"列表平铺"规范，按列表元素字段结构平铺到顶层）
- 不手工修改 `website/` 目录（由收尾阶段 `make doc` 生成）
- 不生成 `.changelog/` 文件（由收尾阶段 tfpacer-finalize 统一处理）

## Vendor 云 API 关键定义（实施阶段直接引用，无需回读 vendor）

以下摘录自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`。

### CreateFunctionReplicaRequest（models.go:3307）

```go
type CreateFunctionReplicaRequest struct {
	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// 函数 ID。
	FunctionId *string `json:"FunctionId,omitnil,omitempty" name:"FunctionId"`
	// 边缘函数副本名称。限制可输入 1-50 个字符，允许的字符为a-z、0-9、-，且-不能单独注册或连续使用，不能放在开头或结尾。同一 FunctionId 下副本名称需唯一。
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty" name:"ReplicaName"`
	// 边缘函数副本内容，当前仅支持 JavaScript 代码，最大支持 5MB。
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
	// 边缘函数副本描述。最大支持 50 个字符。
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}
```

`CreateFunctionReplicaResponse`（models.go:3355）仅含 `Response.RequestId *string`（同步接口，无任务 ID、无资源 ID 出参）。SDK 注释（client.go:1531 起）：创建边缘函数副本。

### DescribeFunctionReplicasRequest（models.go:10323）

```go
type DescribeFunctionReplicasRequest struct {
	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// 函数 ID。
	FunctionId *string `json:"FunctionId,omitnil,omitempty" name:"FunctionId"`
	// 分页查询偏移量。默认值：0。
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`
	// 分页查询限制数目。默认值：20，最大值：200。
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
	// 排序依据，取值有：created-on：创建时间。 默认根据 created-on 属性排序。
	SortBy *string `json:"SortBy,omitnil,omitempty" name:"SortBy"`
	// 列表排序方式，取值有：asc：升序排列；desc：降序排列。 默认值为 asc。
	SortOrder *string `json:"SortOrder,omitnil,omitempty" name:"SortOrder"`
	// 过滤条件，Filters.Values 的上限为 20。该参数不填写时，返回函数 ID 下全部函数副本。
	// 详细的过滤条件如下：replica-name：按照函数副本名称进行过滤，支持模糊查询。
	Filters []*AdvancedFilter `json:"Filters,omitnil,omitempty" name:"Filters"`
}
```

`DescribeFunctionReplicasResponse`（models.go:10385）：
```go
type DescribeFunctionReplicasResponseParams struct {
	// 边缘函数副本总数。
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	// 边缘函数副本列表。
	FunctionReplicas []*FunctionReplica `json:"FunctionReplicas,omitnil,omitempty" name:"FunctionReplicas"`
	RequestId *string
}
```

`FunctionReplica`（models.go:17069）：
```go
type FunctionReplica struct {
	// 函数 ID。
	FunctionId *string
	// 边缘函数副本名称。
	ReplicaName *string
	// 边缘函数副本内容。格式为 JavaScript 代码。
	Content *string
	// 边缘函数副本描述。
	Remark *string
	// 边缘函数副本创建时间。
	CreatedOn *string
	// 边缘函数副本更新时间。
	ModifiedOn *string
}
```

`AdvancedFilter`（models.go:374）：
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

### ModifyFunctionReplicaRequest（models.go:20084）

```go
type ModifyFunctionReplicaRequest struct {
	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// 函数 ID。
	FunctionId *string `json:"FunctionId,omitnil,omitempty" name:"FunctionId"`
	// 需要修改的边缘函数副本名称。
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty" name:"ReplicaName"`
	// 边缘函数副本内容，当前仅支持 JavaScript 代码，最大支持 5MB。
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
	// 边缘函数副本描述。最大支持 50 个字符。
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}
```

`ModifyFunctionReplicaResponse`（models.go:20132）仅含 `Response.RequestId`（同步接口）。

### DeleteFunctionReplicaRequest（models.go:6817）

```go
type DeleteFunctionReplicaRequest struct {
	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// 函数 ID。
	FunctionId *string `json:"FunctionId,omitnil,omitempty" name:"FunctionId"`
	// 需要删除的函数的副本名称。支持以列表的形式传入。
	ReplicaNames []*string `json:"ReplicaNames,omitnil,omitempty" name:"ReplicaNames"`
}
```

`DeleteFunctionReplicaResponse`（models.go:6857）仅含 `Response.RequestId`（同步接口）。

### Client 方法（client.go，均存在）

- `CreateFunctionReplicaWithContext(ctx, request) (*CreateFunctionReplicaResponse, error)`（client.go:1567）
- `DescribeFunctionReplicasWithContext(ctx, request) (*DescribeFunctionReplicasResponse, error)`（client.go:7595）
- `ModifyFunctionReplicaWithContext(ctx, request) (*ModifyFunctionReplicaResponse, error)`（client.go:13047）
- `DeleteFunctionReplicaWithContext(ctx, request) (*DeleteFunctionReplicaResponse, error)`（client.go:4881）

客户端获取方式：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()`。

## 参考代码结构说明

### tencentcloud_igtm_strategy（风格基准，tencentcloud/services/igtm/resource_tc_igtm_strategy.go）

- **Schema 组织**：`ResourceTencentCloudIgtmStrategy()` 返回 `&schema.Resource{}`，`Create/Read/Update/Delete` 四个函数指针 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}` + `Schema: map[string]*schema.Schema{}`；每个字段带 `Type` / `Required|Optional|Computed` / `ForceNew`（身份字段）/ `Description`（英文、句号结尾）
- **CRUD 骨架**（每个函数统一模式）：
  1. `defer tccommon.LogElapsed("resource.tencentcloud_<name>.<op>")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
  2. `logId := tccommon.GetLogId(tccommon.ContextNil)`；`ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`
  3. `request = <sdk>.NewXxxRequest()`，用 `d.GetOk()` + `helper.String()` / `helper.Int64()` 填充入参
  4. `reqErr := resource.Retry(tccommon.WriteRetryTimeout | tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`：API 错误返回 `tccommon.RetryError(e)`；成功打 `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`；`result == nil || result.Response == nil` 返回 `resource.NonRetryableError(...)`
  5. Retry 块外：`if reqErr != nil { log.Printf("[CRITAL]%s ... reason:%+v", logId, reqErr); return reqErr }`
  6. Create 末尾 `d.SetId(strings.Join([]string{...}, tccommon.FILED_SP))` 后 `return resourceXxxRead(d, meta)`；Update 末尾同样回 Read；Delete 成功直接 `return nil`
- **Update 模式**：拆 ID → `needChange := false; mutableArgs := []string{...}` + `d.HasChange(v)` 判断 → 构造 request → WriteRetry 调用 → 回 Read
- **导入路径**：`tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"` 与 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"`（SDK import 别名 `teov20220901`）

### 旧资源 resource_tc_teo_function_replica.go（同 API 的现有实现，字段语义与 Read 定位逻辑直接可复用）

- Read：`idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`（len==3 校验，错误信息 `"id is broken, id is %s"`）→ `request.Filters = []*teov20220901.AdvancedFilter{{Name: helper.String("replica-name"), Values: []*string{helper.String(replicaName)}}}` + `request.Limit = helper.Int64(200)` → ReadRetry → `response == nil || response.Response == nil` 与未匹配到 `targetReplica` 时打印 `[WARN]...not found` 后 `d.SetId("")`（本变更需按规范加 `log.Printf("[CRUD] teo_function_replica_v5 id=%s", d.Id())` 现场日志）→ 逐字段判 nil 后 `_ = d.Set(...)`
- Delete：`request.ReplicaNames = []*string{helper.String(replicaName)}`（单元素列表）
- 新资源函数命名：`ResourceTencentCloudTeoFunctionReplicaV5` / `resourceTencentCloudTeoFunctionReplicaV5Create/Read/Update/Delete`

### AdvancedFilter 的 schema 构造先例（data_source_tc_teo_zones.go:407）

```go
if v, ok := d.GetOk("filters"); ok {
	filtersSet := v.([]interface{})
	tmpSet := make([]*teov20220901.AdvancedFilter, 0, len(filtersSet))
	for _, item := range filtersSet {
		filtersMap := item.(map[string]interface{})
		advancedFilter := teov20220901.AdvancedFilter{}
		if v, ok := filtersMap["name"].(string); ok && v != "" {
			advancedFilter.Name = helper.String(v)
		}
		if v, ok := filtersMap["values"]; ok {
			valuesSet := v.(*schema.Set).List()
			for i := range valuesSet {
				advancedFilter.Values = append(advancedFilter.Values, helper.String(valuesSet[i].(string)))
			}
		}
		if v, ok := filtersMap["fuzzy"].(bool); ok {
			advancedFilter.Fuzzy = helper.Bool(v)
		}
		tmpSet = append(tmpSet, &advancedFilter)
	}
}
```

对应 filters schema（TypeList + Elem Resource；`values` 用 TypeSet of String）：
```go
"filters": {
	Type:        schema.TypeList,
	Optional:    true,
	Description: "...",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":  {Type: schema.TypeString, Required: true, Description: "Field to be filtered."},
			"values": {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Value of the filtered field."},
			"fuzzy": {Type: schema.TypeBool, Optional: true, Description: "Whether to enable fuzzy query."},
		},
	},
},
```

### provider 注册（参考 tencentcloud_igtm_strategy）

`provider.go` ResourcesMap（teov 区段，位于 `tencentcloud_teo_function_replica` 附近，约 2154 行）：
```go
"tencentcloud_teo_function_replica_v5": teo.ResourceTencentCloudTeoFunctionReplicaV5(),
```
（对齐列宽注释风格与相邻条目一致。）

`provider.md` 资源清单（`tencentcloud_teo_function_replica` 下一行追加）：
```
tencentcloud_teo_function_replica_v5
```

### 单元测试（gomonkey，参考 resource_tc_teo_function_replica_test.go 风格）

- 包 `teo_test`；`mockMeta` 结构体实现 `tccommon.ProviderMeta`（`GetAPIV3Conn() *connectivity.TencentCloudClient`），`newMockMeta()` 返回含空 `&connectivity.TencentCloudClient{}` 的实例
- `patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)`（`teoClient := &teov20220901.Client{}`）
- `patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {...})` mock 各接口
- `d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` + `d.SetId(...)`，调用 `res.Create/Read/Update/Delete(d, meta)` 后 `assert.NoError` 与 `assert.Equal`
- 注意：Create 会经由 Read 回填，mock 需同时覆盖 `DescribeFunctionReplicasWithContext`；测试函数命名 `TestTeoFunctionReplicaV5_Create` 等，并附 `// go test ...` 运行注释

## Decisions

### Decision 1: 资源命名与文件命名采用 `_v5` 后缀（独立新资源）

**选择**：新资源 `tencentcloud_teo_function_replica_v5`，文件 `resource_tc_teo_function_replica_v5.go`，导出函数 `ResourceTencentCloudTeoFunctionReplicaV5`。

**备选**：直接扩展现有 `tencentcloud_teo_function_replica` 的 schema。

**理由**：
- 需求明确为"新增 RESOURCE_KIND_GENERAL 资源 tencentcloud_teo_function_replica_v5"，且 provider 已有版本化资源先例（`tencentcloud_teo_l7_acc_rule_v2`、`tencentcloud_emr_cluster_v2`、`tencentcloud_cynosdb_cluster_v2` 等），`_v2` 后缀模式成熟，`_v5` 同理
- 修改现有资源 schema 会触碰旧 spec `teo-function-replica`（需 MODIFIED delta），并可能对存量 state 引入 plan drift；新增资源是纯加法，风险最低

### Decision 2: 复合 ID 采用三段 `zone_id#function_id#replica_name`

**选择**：`d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))`，Read/Update/Delete 均从 `d.Id()` 拆出三段作为请求身份参数。

**理由**：
- 云 API 无独立资源 ID，副本在 `(zone_id, function_id, replica_name)` 三元组下唯一（SDK 注释：同一 FunctionId 下副本名称需唯一）
- 旧资源 `tencentcloud_teo_function_replica` 已采用同构三段 ID，行为一致便于用户从旧资源迁移到 `_v5`；import 示例需说明使用该联合 id
- 符合 PROVIDER 规范"多字段联合 id 使用 tccommon.FILED_SP 分隔，read/update/delete 从 d.Id() 获取 id 作为 request 一部分"

### Decision 3: 身份字段 ForceNew，仅 `content` / `remark` 可更新

**选择**：`zone_id` / `function_id` / `replica_name` 为 `Required + ForceNew`；`mutableArgs = ["content", "remark"]` 触发 `ModifyFunctionReplica`。

**理由**：
- `ModifyFunctionReplicaRequest` 只有 `ZoneId` / `FunctionId` / `ReplicaName`（定位身份）+ `Content` / `Remark`（可变属性），身份字段不可变（变更即销毁重建）
- Create 接口的五个入参与 Update 可变字段完全覆盖需求映射，无遗漏

### Decision 4: `sort_by` / `sort_order` / `filters` 作为查询侧参数透传，Read 定位仍用内部 replica-name 过滤

**选择**：schema 暴露 `sort_by`、`sort_order`（Optional, TypeString）与 `filters`（Optional, TypeList, AdvancedFilter 结构：name Required / values Required TypeSet / fuzzy Optional TypeBool）；Read 时将用户 filters 追加到 request.Filters，并**额外追加**内部过滤条件 `{Name: "replica-name", Values: [replica_name]}` 精确定位本资源，最后仍在内存中按 `ReplicaName == replicaName` 精确匹配兜底。

**备选**：仅暴露基础字段（旧资源做法），或将 filters 作为 Update 触发条件。

**理由**：
- 需求映射明确要求这四个查询参数进入资源 schema（`request.SortBy` / `request.SortOrder` / `request.Filters.*`）
- 排序/过滤参数只影响 Describe 的结果集与顺序，不影响资源实体，作为 Optional 参数在 Read 中透传即可，不进入 `mutableArgs`（不会触发 Modify 调用；云 API 也没有对应可写语义）
- 追加内部 `replica-name` 过滤可减少返回数据量并保证定位唯一性；内存精确匹配兜底防止模糊过滤误匹配（用户配置 `fuzzy=true` 时服务端可能返回多个副本）

### Decision 5: `replica_names` 作为 Delete 侧参数，回退复合 ID 的 replica_name

**选择**：schema 暴露 `replica_names`（Optional, TypeList of String）；Delete 时若用户配置了 `replica_names`（非空），以其构造 `request.ReplicaNames`；否则回退 `[]*string{helper.String(replicaName)}`（复合 ID 第三段）。

**备选**：不暴露 `replica_names`，Delete 永远只删复合 ID 中的单个副本（旧资源做法）。

**理由**：
- 需求映射明确要求 `request.ReplicaNames → replica_names`（必填标注）
- 云 API 语义支持按名称列表批量删除；暴露该参数对齐 API 能力，同时保留单资源删除的回退路径保证资源语义完整（一个 TF 资源对应一个副本，默认删除自身）
- 该参数不参与 Read 回填与 Update（Describe 响应无此字段、Modify 无此入参）

### Decision 6: `created_on` / `modified_on` 为 Computed 属性

**选择**：`created_on` / `modified_on`（Computed, TypeString），Read 时从匹配到的 `FunctionReplica` 判 nil 后 `_ = d.Set(...)`。

**理由**：
- 需求映射将 `response.Response.FunctionReplicas.CreatedOn` / `ModifiedOn` 映射到资源参数，二者仅出现在 Describe 出参，属于只读计算属性
- 符合 PROVIDER 规范"set 前判 nil"

### Decision 7: 列表响应平铺，不引入 `function_replicas` 包装层

**选择**：schema 顶层直接包含列表元素的字段（`function_id` / `replica_name` / `content` / `remark` / `created_on` / `modified_on`），不创建 `function_replicas` / `function_replica_set` 嵌套列表层。

**理由**：
- PROVIDER 规范明确：资源参数 schema 禁止"列表型数据"包装层，应把列表元素字段平铺到顶层（按"列表中第一项"的字段结构定义）
- 通用资源本身即描述单个副本，平铺语义自然

### Decision 8: Retry 与空响应处理遵循 PROVIDER 规范

**选择**：
- Create/Update/Delete 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，Read 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)`；API 错误一律 `tccommon.RetryError(e)` 包装
- Create Retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`；拿到响应后（块外）SetId——本接口响应无资源 ID 出参，以请求参数三段构造复合 ID；SetId 前打印 `logId` 便于排障
- Read 中 response 为空或未匹配到副本时：先 `log.Printf("[CRUD] teo_function_replica_v5 id=%s", d.Id())` 保留现场，再 `d.SetId("")`
- `DescribeFunctionReplicas` 的 `Limit` 固定为 `helper.Int64(200)`（API 注释标注的最大值）

**理由**：与 PROVIDER 规范及仓库现有模式（igtm_strategy、旧 function_replica）完全一致；同步接口无 FlowId，无需异步轮询，故 SetId 紧跟 Create 成功（先 SetId 后 Read 回填，满足"拿到 id 后先 SetId"的规范精神）。

## Risks / Trade-offs

- **Risk**：用户在 `filters` 中配置了与 `replica-name` 冲突的过滤条件（如模糊查询返回多副本）导致 Read 误匹配 → **Mitigation**：Read 在内存中按 `ReplicaName` 精确相等匹配兜底，匹配不到按"not found"处理（SetId("")）
- **Risk**：`replica_names` 允许用户删除列表中不包含自身 `replica_name` 的副本，导致资源实际存在但 TF 认为已删除 → **Mitigation**：文档 Description 中明确说明默认行为（不配置时仅删除资源自身对应的副本）；该字段为可选，常规用法留空
- **Risk**：`Limit=200` 之外还有更多副本（分页未取全）导致目标副本不在首页 → **Mitigation**：函数副本单函数数量远小于 200，且 Read 叠加了 `replica-name` 精确过滤条件；若极端场景仍未取全，Read 按 not found 处理并可通过 refresh 重试（与旧资源行为一致）
- **Trade-off**：新增 `_v5` 资源与旧 `tencentcloud_teo_function_replica` 功能重叠 → 需求方明确要求新资源名；旧资源保持不动，用户可自行选择，后续由产品决策是否废弃旧资源
- **Trade-off**：查询侧参数（sort_by/sort_order/filters）只影响 Read 查询，不影响云上资源状态 → 明确不进入 mutableArgs，避免无意义的 Modify 调用

## Migration Plan

- 纯新增资源，无迁移需求；现有 `tencentcloud_teo_function_replica` 用户不受影响
- 用户从旧资源迁移到 `_v5`：`terraform state mv` 或 import（`terraform import tencentcloud_teo_function_replica_v5.example zone-id#function-id#replica-name`，ID 格式与旧资源一致）
- 回滚：删除 provider.go 注册条目与三个新文件即可，无 state 兼容负担

## Open Questions

- 无（所有接口与字段均已在 vendor 中验证存在，映射关系与需求描述完全一致）
