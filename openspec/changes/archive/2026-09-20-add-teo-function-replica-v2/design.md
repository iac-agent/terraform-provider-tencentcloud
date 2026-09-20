## Context

腾讯云 Terraform Provider（Terraform Plugin SDK v2，Go 1.17+）需要为 EdgeOne（TEO）新增边缘函数副本资源 `tencentcloud_teo_function_replica_v2`。

- 云 API 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（vendor 已就绪，无需更新依赖）。
- 仓库已有同名 V1 资源 `resource_tc_teo_function_replica.go`，但其在 provider.go 中未注册（实际不可用）、无 md 文档、Read/Update/Delete 逻辑不规范。本变更新增 V2 资源，不动 V1。
- 资源类型：RESOURCE_KIND_GENERAL（完整 CRUD 生命周期管理）。
- 四个接口均为同步接口（响应中仅 RequestId，无 FlowId/TaskId/JobId），无需异步轮询。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_replica_v2` 资源，完整对接 Create/Describe/Modify/Delete 四个云 API。
- 严格参照 `tencentcloud_igtm_strategy` 的代码风格（schema 组织、CRUD 骨架、retry 包装、错误处理、日志规范）。
- 支持复合 ID `zone_id#function_id#replica_name` 与 `terraform import`。
- 补充 gomonkey mock 单元测试与资源 md 文档，并在 provider.go/provider.md 注册。

**Non-Goals:**
- 不修改/不废弃 V1 `tencentcloud_teo_function_replica` 资源。
- 不新增数据源（RESOURCE_KIND_DATASOURCE）。
- 不修改 vendor 依赖。
- 不在 website/ 目录手写文档（由收尾阶段 `make doc` 生成）。

## Decisions

### D1. 复合 ID 与字段可变性

ID 格式：`strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP)`（即 `#` 分隔）。
Read/Update/Delete 一律从 `d.Id()` split 解析出三元组定位资源（云 API 无单查接口，只能按 `replica-name` Filter 查询列表后匹配）。
`zone_id`/`function_id`/`replica_name` 三个字段 `ForceNew: true`（Create 专有定位字段，改名即重建）；`content`/`remark` 可更新。
Update 中维护 `immutableArgs := []string{"zone_id", "function_id", "replica_name"}`：任一发生 `d.HasChange` 即返回 error（而非静默重建）。
Delete 使用 `ReplicaNames []*string`（云 API 列表型入参），传单个 `replicaName`。

### D2. Schema 顶层平铺，查询辅助参数不进 schema

`DescribeFunctionReplicas` 的入参 `sort_by`/`sort_order`/`filters(name/values/fuzzy)` 仅是查询能力，不属于资源属性，**不**放入资源 schema（与 V1 一致，避免与"资源属性"语义混淆；资源定位靠 zone_id+function_id+replica-name Filter）。响应中的列表 `function_replicas` 也不作为 schema 字段——本资源是"单个副本"资源，Read 从列表中匹配 `replica_name` 后平铺 set 各字段；`created_on`/`modified_on` 作为 Computed 输出字段。

### D3. CRUD 骨架与 retry 规范（参照 igtm_strategy）

- Create：`resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调 `CreateFunctionReplicaWithContext`；response 为 nil 时返回 `resource.NonRetryableError`；retry 成功后（retry 块外）`d.SetId(...)` 并调 Read。
- Read：`resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调 `DescribeFunctionReplicasWithContext`；`request.Limit = helper.Int64(200)`（云 API 注释标注的最大值）；retry 块外判断 response/列表为空或未匹配到目标副本时，先 `log.Printf("[CRUD] ...")` 保留现场再 `d.SetId("")` 返回 nil。
- Update：仅当 mutable 字段（content/remark）`d.HasChange` 时构造 `ModifyFunctionReplicaRequest` 调用；随后调 Read。
- Delete：`resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调 `DeleteFunctionReplicaWithContext`。
- 每个 CRUD 函数：`defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.<op>")()` + `defer tccommon.InconsistentCheck(d, meta)()`；`ctx = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`；错误包装用 `tccommon.RetryError(e)`；调试日志 `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`。
- 客户端获取：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()`。

### D4. 测试与文档

- 单测：gomonkey mock（参照 `resource_tc_teo_function_replica_test.go` 的既有 mock 模式：`mockMeta` 实现 `tccommon.ProviderMeta`、`patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)`、`patches.ApplyMethodFunc` mock 各 WithContext 方法、`schema.TestResourceDataRaw` 构造 d），覆盖 Create/Read/ReadNotFound/Update/UpdateImmutableError/Delete。
- md：`Provides a resource to create a TEO function replica`（含云产品名 TEO）+ Example Usage + Import（说明使用 `zoneId#functionId#replicaName` 复合 ID）；不写 Argument/Attribute Reference。

## 云 API 参考（从 vendor 摘录，实施时无需回读 vendor）

文件：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`

### CreateFunctionReplicaRequest（行 3307-3324）
```go
type CreateFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId      *string `json:"ZoneId,omitnil,omitempty"`      // 站点 ID
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 副本名，1-50 字符，a-z/0-9/-，-不能单独/连续/首尾，同一 FunctionId 下唯一
	Content     *string `json:"Content,omitnil,omitempty"`     // 副本内容，JavaScript，最大 5MB
	Remark      *string `json:"Remark,omitnil,omitempty"`      // 副本描述，最大 50 字符
}
```
`CreateFunctionReplicaResponse` 仅含 `RequestId`（无资源 ID 返回，同步接口）。客户端方法：`CreateFunctionReplicaWithContext(ctx, request)`。

### DescribeFunctionReplicasRequest（行 10252-10275）
```go
type DescribeFunctionReplicasRequest struct {
	*tchttp.BaseRequest
	ZoneId      *string           `json:"ZoneId,omitnil,omitempty"`      // 站点 ID
	FunctionId  *string           `json:"FunctionId,omitnil,omitempty"`  // 函数 ID
	Offset      *int64            `json:"Offset,omitnil,omitempty"`      // 默认 0
	Limit       *int64            `json:"Limit,omitnil,omitempty"`       // 默认 20，最大 200
	SortBy      *string           `json:"SortBy,omitnil,omitempty"`      // created-on（默认）
	SortOrder   *string           `json:"SortOrder,omitnil,omitempty"`   // asc（默认）/desc
	Filters     []*AdvancedFilter `json:"Filters,omitnil,omitempty"`     // replica-name 过滤，支持模糊
}
```
### DescribeFunctionReplicasResponseParams（行 10303-10312）
```go
type DescribeFunctionReplicasResponseParams struct {
	TotalCount        *int64             `json:"TotalCount,omitnil,omitempty"`
	FunctionReplicas  []*FunctionReplica `json:"FunctionReplicas,omitnil,omitempty"`
	RequestId         *string            `json:"RequestId,omitnil,omitempty"`
}
```
### AdvancedFilter（行 374-383）
```go
type AdvancedFilter struct {
	Name   *string   `json:"Name,omitnil,omitempty"`   // 字段名，如 "replica-name"
	Values []*string `json:"Values,omitnil,omitempty"` // 过滤值，上限 20
	Fuzzy  *bool     `json:"Fuzzy,omitnil,omitempty"`  // 是否模糊查询
}
```
### FunctionReplica（行 16998-17016）
```go
type FunctionReplica struct {
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"`
	Content     *string `json:"Content,omitnil,omitempty"`
	Remark      *string `json:"Remark,omitnil,omitempty"`
	CreatedOn   *string `json:"CreatedOn,omitnil,omitempty"`  // 创建时间
	ModifiedOn  *string `json:"ModifiedOn,omitnil,omitempty"` // 更新时间
}
```
客户端方法：`DescribeFunctionReplicasWithContext(ctx, request)`。

### ModifyFunctionReplicaRequest（行 20027-20044）
```go
type ModifyFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId      *string `json:"ZoneId,omitnil,omitempty"`
	FunctionId  *string `json:"FunctionId,omitnil,omitempty"`
	ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 需要修改的副本名称（定位字段）
	Content     *string `json:"Content,omitnil,omitempty"`     // 可选更新
	Remark      *string `json:"Remark,omitnil,omitempty"`      // 可选更新
}
```
`ModifyFunctionReplicaResponse` 仅含 `RequestId`。客户端方法：`ModifyFunctionReplicaWithContext(ctx, request)`。

### DeleteFunctionReplicaRequest（行 6827-6838）
```go
type DeleteFunctionReplicaRequest struct {
	*tchttp.BaseRequest
	ZoneId       *string   `json:"ZoneId,omitnil,omitempty"`
	FunctionId   *string   `json:"FunctionId,omitnil,omitempty"`
	ReplicaNames []*string `json:"ReplicaNames,omitnil,omitempty"` // 支持列表批量，本资源传单个
}
```
`DeleteFunctionReplicaResponse` 仅含 `RequestId`。客户端方法：`DeleteFunctionReplicaWithContext(ctx, request)`。

## 参照代码结构说明

### 风格参照：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`

- `ResourceTencentCloudIgtmStrategy() *schema.Resource`：Create/Read/Update/Delete 四个函数指针 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}` + Schema map。
- Create 骨架：`defer LogElapsed/InconsistentCheck` → var 块（logId/ctx/request/局部变量）→ `d.GetOk` 逐字段填充 request → `resource.Retry(WriteRetryTimeout)`（内含 API 调用 + nil 检查 + response 赋值）→ retry 后取 ID 字段并 `d.SetId(strings.Join(..., tccommon.FILED_SP))` → `return resourceXxxRead(d, meta)`。
- Update 骨架：split ID → `needChange`/`mutableArgs` 模式（先扫描 HasChange）→ 需要变更才构造 request 并 `resource.Retry(WriteRetryTimeout)` → Read 收尾。
- Delete 骨架：split ID → 填充 request → `resource.Retry(WriteRetryTimeout)` → 返回 nil。
- Read 骨架：split ID → 调查询（服务层或 retry 直调）→ 未查到时 `log.Printf("[WARN]...")` + `d.SetId("")` + return nil → 逐字段 `if respData.X != nil { _ = d.Set(...) }`。

### 同域参照：`tencentcloud/services/teo/resource_tc_teo_function_replica.go`（V1，仅借鉴 API 调用方式与 Filter 用法）

- Read 中 Filter 构造：`request.Filters = []*teov20220901.AdvancedFilter{{Name: helper.String("replica-name"), Values: []*string{helper.String(replicaName)}}}`，`request.Limit = helper.Int64(200)`，随后遍历 `response.Response.FunctionReplicas` 匹配 `*replica.ReplicaName == replicaName`。
- 其 mock 单测文件 `resource_tc_teo_function_replica_test.go` 可直接作为 V2 单测模板（mockMeta/ApplyMethodReturn/ApplyMethodFunc/TestResourceDataRaw 模式已验证可编译）。

### provider 注册参照

- `tencentcloud/provider.go` 约 2154 行附近，格式：`"tencentcloud_teo_function_replica_v2":  teo.ResourceTencentCloudTeoFunctionReplicaV2(),`（按字典序插入 teo 资源块）。
- `tencentcloud/provider.md` 约 1624 行 `tencentcloud_teo_function_replica` 之后追加一行 `tencentcloud_teo_function_replica_v2`。

## Risks / Trade-offs

- [Read 依赖列表查询而非单查接口] → 用 `replica-name` Filter + 200 Limit 精确匹配；若同函数下副本数超过 200（罕见），匹配失败会被误判为 not found；后续如有需要可加 Offset 翻页，当前与 V1 行为保持一致，风险可接受。
- [Create API 不返回资源 ID] → 以请求参数三元组构造复合 ID；Create 成功后立即调 Read 校验副本确实存在（若未查到会清空 ID，符合 TF 语义）。
- [V1 资源已存在但未注册] → V2 使用独立资源名与独立文件，命名空间互不冲突；不处理 V1 的历史问题，避免破坏任何既有 state。
- [immutable 字段被用户在配置中修改] → Update 显式返回 error 提示重建，而非静默忽略或自动重建，行为可预期。

## Migration Plan

纯新增资源，无迁移。回滚 = revert 该变更（删除新增文件 + 移除 provider 注册两行）。

## Open Questions

无。
