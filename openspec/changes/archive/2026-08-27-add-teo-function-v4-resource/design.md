## Context

腾讯云 EdgeFunctions（边缘函数，TEO）允许用户在边缘节点执行 JavaScript 代码。当前 vendor 中 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 已经暴露了完整的边缘函数 CRUD 同步接口：

- `CreateFunctionWithContext` → 返回 `{ FunctionId }`（同步接口，无 TaskId）。
- `DescribeFunctionsWithContext` → 分页列表，返回 `Functions []*Function`，`Function` 结构体包含 `FunctionId`、`ZoneId`、`Name`、`Remark`、`Content`、`Domain`、`CreateTime`、`UpdateTime`。
- `ModifyFunctionWithContext` → 入参 `ZoneId`、`FunctionId`、`Remark`、`Content`（同步接口，不返回 TaskId）。
- `DeleteFunctionWithContext` → 入参 `ZoneId`、`FunctionId`（同步接口）。

provider 已经存在 `tencentcloud_teo_function` 资源（文件 `resource_tc_teo_function.go`），但它使用了 `io.ReadAll` + `text/template` 的旧式异步轮询 `resourceTeoFunctionCreateStateRefreshFunc_0_0` 来等待 `Domain` 字段就绪。本次需求是新增一个**独立**的 `tencentcloud_teo_function_v4` 资源（`RESOURCE_KIND_GENERAL`），按统一的代码风格（参考 `tencentcloud_igtm_strategy`）实现完整生命周期，不修改既有 `tencentcloud_teo_function`。

代码生成硬约束明确要求：
- 非数据源（DATASOURCE）资源代码风格严格参考 `tencentcloud_igtm_strategy`。
- 在调用云 API 接口时以 `tccommon.ReadRetryTimeout` 作为超时时间添加 retry 处理；调用失败使用 `tccommon.RetryError()` 包装错误。
- 复合 ID 使用 `tccommon.FILED_SP` 分隔；导入示例需说明使用联合 ID。
- Read/Create 的空返回需先打印日志再 `d.SetId("")`。
- 资源仅有 CRUD 接口时，将 Id 字段设为 ForceNew，其余顶层字段加入 `immutableArgs`，发现变更则报错。
- 服务层查询接口展开列表，不引入 `xxx_set`/`xxx_list` 嵌套层。

现有 `TeoService.DescribeTeoFunctionById(ctx, zoneId, functionId)` 已使用 `FunctionIds` 精确查询 `DescribeFunctions`，并返回 `*teo.Function`。本次可直接复用其调用模式新增一个 `DescribeTeoFunctionV4ById` 辅助函数（避免与既有资源耦合，便于后续单独演进）。

## Goals / Non-Goals

**Goals:**
- 提供一个独立、可单独演进的 `tencentcloud_teo_function_v4` 资源，覆盖 TEO 边缘函数的创建 / 查询 / 修改 / 删除 / 导入。
- schema 字段与 `CreateFunction` 入参严格对齐：`zone_id`（ForceNew）、`name`（ForceNew）、`content`、`remark`；并暴露 `DescribeFunctions` 返回的 computed 字段 `function_id`、`domain`、`create_time`、`update_time`。
- 所有 SDK 调用包裹在 `resource.Retry(...)` 中（写操作 `tccommon.WriteRetryTimeout`、读操作 `tccommon.ReadRetryTimeout`），失败统一使用 `tccommon.RetryError(e)` 包装。
- 复合 ID `zone_id#function_id`，支持 `terraform import` 使用联合 ID。
- 复用既有 `ParseTeoFunctionOriginalName` 的 name 拆分逻辑，避免 `DescribeFunctions` 返回的拼接 name 触发 plan 差异。
- 代码风格与 `tencentcloud_igtm_strategy` 保持一致（单文件资源布局、retry、nil 防御）。
- 新增 gomonkey mock 单元测试，覆盖业务逻辑（不使用 terraform 测试套件、不执行 `go test`）。

**Non-Goals:**
- 不修改既有 `tencentcloud_teo_function` 资源的任何代码或 state。
- 不实现边缘函数副本（Replica）、规则（Rule）、运行时环境（RuntimeEnvironment）、组件绑定（ComponentBinding）等其他相关资源（已有独立资源覆盖）。
- 不实现数据源 `data_source_tc_teo_function_v4`（本次仅资源）。
- 不引入异步任务轮询（四个目标接口均为同步接口，无 `TaskId`）。
- 不引入 tag 管理（边缘函数 API 不支持 tag 参数）。
- 不修改 vendor 下任何 SDK 文件。

## Decisions

### D1. 新增独立资源而非修改既有资源
**Why**: 需求明确要求新增 `tencentcloud_teo_function_v4`，且硬约束要求保持向后兼容、不破坏既有 TF 配置与 state。新建独立文件 `resource_tc_teo_function_v4.go` 与函数符号（`ResourceTencentCloudTeoFunctionV4` 等）可彻底隔离影响面。
**Alternative**: 复用既有 `tencentcloud_teo_function` 文件并新增符号。Rejected — 会引入同一文件内两个生命周期实现，维护成本高，且违背"新增资源"语义。

### D2. 复合 ID = `zone_id#function_id`
**Why**: `CreateFunction` 返回 `FunctionId`，但 `DescribeFunctions`/`ModifyFunction`/`DeleteFunction` 均要求同时传入 `ZoneId` 与 `FunctionId`。复合 ID 让 Read/Update/Delete 可直接从 `d.Id()` 解析出两个必要参数，无需额外 state 字段串联。
**Alternative**: 仅用 `FunctionId` 作为 ID，从 `zone_id` schema 字段取 ZoneId。Rejected — 导入场景下 `zone_id` 为空会导致 Update/Delete 失败；复合 ID 是 provider 内统一约定（`tencentcloud_teo_function` 即采用此约定）。

### D3. ForceNew 与 immutableArgs 设计
- `function_id`：computed，资源标识，不可由用户配置。
- `zone_id`：`Required + ForceNew`（所有四个接口都需要它，修改站点等同于重建）。
- `name`：`Required + ForceNew`。`ModifyFunction` 不接受 `name` 参数，因此 name 变更只能重建；按硬约束纳入 `immutableArgs` 校验（若发生变更返回 error）。
- `content`、`remark`：可更新字段，纳入 `mutableArgs`，触发 `ModifyFunction`。
**Why**: 与云 API 能力严格一致，避免向 `ModifyFunction` 传入其不支持的参数。
**Alternative**: 放开 `name` 为可更新。Rejected — SDK `ModifyFunctionRequest` 无 `Name` 字段，强行调用会被服务端拒绝或忽略，造成 state drift。

### D4. 复用 `DescribeFunctions` + `FunctionIds` 精确查询
新增 `TeoService.DescribeTeoFunctionV4ById(ctx, zoneId, functionId) (*teo.Function, error)`：
- 构造 `DescribeFunctionsRequest`，`ZoneId = zoneId`，`FunctionIds = []*string{functionId}`。
- 包裹 `resource.Retry(tccommon.ReadRetryTimeout, ...)`。
- 当 `response.Response == nil || len(response.Response.Functions) < 1` 时返回 `(nil, nil)`。
- 严格相等校验 `*function.FunctionId == functionId` 后取第一项。
**Why**: 与既有 `DescribeTeoFunctionById` 完全一致的模式，保证两个资源互不影响、可独立演进。
**Alternative**: 直接复用既有 `DescribeTeoFunctionById`。Rejected — 需求要求独立资源，服务层辅助也保持独立命名，避免后续任一资源演进时相互牵连；但实现完全等价。

### D5. name 拆分（ParseTeoFunctionOriginalName 复用）
`DescribeFunctions` 返回的 `Name` 是拼接形式（`<原始name>-<zoneId后缀>`，如 `my-func-zone-2qtuhspy7cr6-1310708577`）。若直接 set 到 state，会导致每次 plan 都因 name 字段不一致而误报变更。
**Decision**: 在新资源文件内新增等价的 `parseTeoFunctionV4OriginalName(name, zoneId string) string` 辅助函数（逻辑与既有 `ParseTeoFunctionOriginalName` 完全一致：以 `-` + zoneId 为后缀定位分割点），并在 Read 中对返回 name 做拆分后再 set。
**Alternative**: 直接复用既有导出函数 `ParseTeoFunctionOriginalName`。Rejected — 该函数绑死在 `tencentcloud_teo_function` 的语义与单测上，为保持资源彻底独立，在新资源内自含同名私有辅助函数；二者算法等价，行为一致。

### D6. retry 拓扑（同步接口，无任务轮询）
所有四个 SDK 调用均在 `resource.Retry` 内执行：
- Create / Modify / Delete：`tccommon.WriteRetryTimeout`。
- Read（服务层）：`tccommon.ReadRetryTimeout`。
- 失败统一 `tccommon.RetryError(e)`。
无 `TaskId`、无异步轮询、无 `StateChangeConf`（区别于既有 `tencentcloud_teo_function` 的 `resourceTeoFunctionCreateStateRefreshFunc_0_0`）。
**Why**: 这四个接口均为同步接口，调用成功即代表操作生效，无需额外的就绪轮询。按代码生成硬约束，retry 块内只执行接口调用，设置 id 等成功操作放到 retry 块外。

### D7. 空 id / 空返回防御
- Create：调用完成后检查 `response == nil || response.Response == nil || response.Response.FunctionId == nil || *response.Response.FunctionId == ""`，任一为空则返回 `tccommon.NonRetryableError`；检查前打印 `logId` 与 `d.Id()` 便于排障。
- Read：服务层返回 `(nil, nil)` 时，先 `log.Printf("[CRUD] teo_function_v4 id=%s", d.Id())` 保留现场，再 `d.SetId("")`。
- Read 中每个 `_ = d.Set(k, v)` 前先判断响应字段是否为 nil，nil 则跳过 set。
**Why**: 避免写入空 id 触发后续状态混乱；符合代码生成硬约束第 8、9 条。

### D8. 文件布局（单文件资源）
按用户既有的严格反馈，整个资源集中在一个文件 `resource_tc_teo_function_v4.go`：package + imports → `ResourceTencentCloudTeoFunctionV4()` schema → CRUD 函数 → `parseTeoFunctionV4OriginalName` 辅助。服务层查询函数置于既有 `service_tencentcloud_teo.go`。
**Why**: 与 `tencentcloud_igtm_strategy` 单文件风格一致，避免 `_crud.go`/`_helpers.go` 拆分。

### D9. 单元测试使用 gomonkey mock
按代码生成硬约束，新增资源的单元测试不使用 terraform 测试套件，而是用 gomonkey mock 云 API 客户端方法（`CreateFunctionWithContext`、`DescribeFunctions`、`ModifyFunctionWithContext`、`DeleteFunctionWithContext`），覆盖 Create/Read/Update/Delete 的业务逻辑分支（含成功路径、空返回防御、immutableArgs 校验、name 拆分辅助函数）。
**Why**: 业务逻辑单测可在无云账号环境下验证代码正确性，符合"不执行集成测试"的约束。

## Risks / Trade-offs

- **[Risk]** `DescribeFunctions` 返回的拼接 `name` 格式可能随 API 版本变化，导致 `parseTeoFunctionV4OriginalName` 无法正确拆分。→ **Mitigation**: 算法对找不到 `-zoneId` 后缀的情况保留原值不做截断（与既有 `ParseTeoFunctionOriginalName` 一致），最坏情况退化为不拆分，不会丢数据。
- **[Risk]** `CreateFunction` 创建后边缘函数的 `Domain` 字段可能需要短暂时间才生效，既有 `tencentcloud_teo_function` 通过轮询 `Domain` 就绪来规避。本次新资源按需求采用同步接口、不轮询 `Domain`，Create 成功后立即走 Read 回填（若 `Domain` 暂未就绪则为空字符串，下次 plan 会自动校正）。→ **Mitigation**: 这是同步语义下的合理行为；`Domain` 为 computed 字段，不参与 diff 决策，不影响资源可用性；如确需等待，可在后续迭代中追加 `StateChangeConf`（本次范围外）。
- **[Trade-off]** 同时存在 `tencentcloud_teo_function` 与 `tencentcloud_teo_function_v4` 两个语义相近的资源。→ 用户明确要求新增独立资源以单独演进，二者互不干扰；文档中会清晰说明二者均管理边缘函数，用户可按需选择。
- **[Risk]** 单元测试依赖 gomonkey 对 SDK client 方法打桩，若 SDK 方法签名变更需要同步更新测试。→ **Mitigation**: mock 目标明确限定为 `UseTeoV20220901Client()` 返回的 `*teo.Client` 的四个 `...WithContext` 方法，签名稳定。

## Migration Plan

纯新增，无需迁移：
1. 落地新资源文件、服务层辅助、provider 注册、文档、单元测试。
2. 发布后用户通过 `resource "tencentcloud_teo_function_v4" "x" { ... }` 采用。
3. 既有 `tencentcloud_teo_function` 资源完全不受影响。

回滚：仅 revert 新增文件与 `provider.go` / `provider.md` 的注册行，无任何 state 变更需处理。

## Open Questions

- 无需用户介入的待决项。四个目标接口均为同步接口且已在 vendor 中确认存在，所有决策（复合 ID、ForceNew 集合、name 拆分、retry 拓扑、测试方式）均已由需求与硬约束确定。
