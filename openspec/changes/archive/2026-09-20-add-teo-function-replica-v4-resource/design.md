## Context

TencentCloud EdgeOne (TEO) 边缘函数副本（Function Replica）允许为边缘函数创建副本用于版本管理与灰度测试。客户端请求匹配触发规则或默认域名时，可通过请求头 `EO-Function-Replica-Name:[副本名称]` 访问特定副本（每个函数默认支持创建两个副本）。

当前状态：
- 已存在 v1 资源 `tencentcloud_teo_function_replica`（`tencentcloud/services/teo/resource_tc_teo_function_replica.go`），覆盖同一组云 API，但未暴露 `created_on`/`modified_on` Computed 字段
- 供应商内已有版本化资源命名先例：`tencentcloud_teo_l7_acc_rule_v2`（在 v1 基础上新增 `_v2` 后缀、独立文件、独立注册），本次 `tencentcloud_teo_function_replica_v4` 遵循同样的惯例
- 云 API 均位于 vendored SDK 包 `teov20220901`（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`），四个接口全部为**同步接口**（响应中无 TaskId/FlowId/JobId），无需异步轮询

### 云 API 接口与 struct 定义（从 vendor 按行区间摘录，实施阶段无需回读 vendor）

**1. CreateFunctionReplica（创建边缘函数副本）**
- client 方法: `func (c *Client) CreateFunctionReplicaWithContext(ctx context.Context, request *CreateFunctionReplicaRequest) (response *CreateFunctionReplicaResponse, err error)`
- `CreateFunctionReplicaRequest`（models.go L3307）:
  ```go
  type CreateFunctionReplicaRequest struct {
      *tchttp.BaseRequest
      ZoneId      *string `json:"ZoneId,omitnil,omitempty"`      // 站点 ID。（必填）
      FunctionId  *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID。（必填）
      ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 副本名称，1-50 字符，a-z/0-9/-，-不能单独/连续使用或首尾；同一 FunctionId 下唯一。（必填）
      Content     *string `json:"Content,omitnil,omitempty"`     // 副本内容，JavaScript 代码，最大 5MB。（必填）
      Remark      *string `json:"Remark,omitnil,omitempty"`      // 副本描述，最大 50 字符。（可选）
  }
  ```
- `CreateFunctionReplicaResponse`（L3355）: 仅含 `Response *CreateFunctionReplicaResponseParams`，ResponseParams 内仅 `RequestId *string`，**无返回 ID 字段** → 资源 ID 必须由入参组合

**2. DescribeFunctionReplicas（查询边缘函数副本列表）**
- client 方法: `DescribeFunctionReplicasWithContext`
- `DescribeFunctionReplicasRequest`（models.go L10252）:
  ```go
  type DescribeFunctionReplicasRequest struct {
      *tchttp.BaseRequest
      ZoneId     *string `json:"ZoneId,omitnil,omitempty"`     // 站点 ID。
      FunctionId *string `json:"FunctionId,omitnil,omitempty"` // 函数 ID。
      Offset     *int64  `json:"Offset,omitnil,omitempty"`     // 分页偏移，默认 0。
      Limit      *int64  `json:"Limit,omitnil,omitempty"`      // 分页限制，默认 20，最大 200。
      SortBy     *string `json:"SortBy,omitnil,omitempty"`     // 排序依据，取值 created-on（默认）。
      SortOrder  *string `json:"SortOrder,omitnil,omitempty"`  // asc/desc，默认 asc。
      Filters    []*AdvancedFilter `json:"Filters,omitnil,omitempty"` // 过滤条件：replica-name（支持模糊查询），Filters.Values 上限 20。
  }
  ```
- `AdvancedFilter`（models.go L374）:
  ```go
  type AdvancedFilter struct {
      Name   *string  `json:"Name,omitnil,omitempty"`   // 过滤字段名。
      Values []*string `json:"Values,omitnil,omitempty"` // 过滤值。
      Fuzzy  *bool    `json:"Fuzzy,omitnil,omitempty"`  // 是否模糊查询。
  }
  ```
- `DescribeFunctionReplicasResponse`（L10314）:
  ```go
  type DescribeFunctionReplicasResponseParams struct {
      TotalCount       *int64             `json:"TotalCount,omitnil,omitempty"`       // 副本总数。
      FunctionReplicas []*FunctionReplica `json:"FunctionReplicas,omitnil,omitempty"` // 副本列表。
      RequestId        *string            `json:"RequestId,omitnil,omitempty"`
  }
  ```
- `FunctionReplica`（models.go L16998）:
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

**3. ModifyFunctionReplica（编辑边缘函数副本）**
- client 方法: `ModifyFunctionReplicaWithContext`
- `ModifyFunctionReplicaRequest`（models.go L20027）:
  ```go
  type ModifyFunctionReplicaRequest struct {
      *tchttp.BaseRequest
      ZoneId      *string `json:"ZoneId,omitnil,omitempty"`
      FunctionId  *string `json:"FunctionId,omitnil,omitempty"`
      ReplicaName *string `json:"ReplicaName,omitnil,omitempty"` // 需要修改的副本名称（定位标识，不能改名）。
      Content     *string `json:"Content,omitnil,omitempty"`     // 可选。
      Remark      *string `json:"Remark,omitnil,omitempty"`      // 可选。
  }
  ```
- `ModifyFunctionReplicaResponse`: 仅含 RequestId，无业务字段

**4. DeleteFunctionReplica（删除边缘函数副本）**
- client 方法: `DeleteFunctionReplicaWithContext`
- `DeleteFunctionReplicaRequest`（models.go L6827）:
  ```go
  type DeleteFunctionReplicaRequest struct {
      *tchttp.BaseRequest
      ZoneId       *string   `json:"ZoneId,omitnil,omitempty"`
      FunctionId   *string   `json:"FunctionId,omitnil,omitempty"`
      ReplicaNames []*string `json:"ReplicaNames,omitnil,omitempty"` // 需要删除的副本名称列表（批量接口，TF 资源粒度传单元素）。
  }
  ```
- `DeleteFunctionReplicaResponse`: 仅含 RequestId，无业务字段

### 代码风格参照（tencentcloud_igtm_strategy / tencentcloud_teo_function_replica 骨架说明）

资源文件结构（单文件 `resource_tc_teo_function_replica_v4.go`，顺序）：
1. `ResourceTencentCloudTeoFunctionReplicaV4() *schema.Resource` — 注册 Create/Read/Update/Delete + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}` + Schema map
2. `resourceTencentCloudTeoFunctionReplicaV4Create` — `defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v4.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`；`logId := tccommon.GetLogId(tccommon.ContextNil)`；`ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`；用 `d.GetOk()` 逐字段填充 request；`resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，错误用 `tccommon.RetryError(e)` 包装；retry 块内校验 `result == nil || result.Response == nil` 返回 `resource.NonRetryableError`；retry 块外 `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))`；末尾调 Read
3. `resourceTencentCloudTeoFunctionReplicaV4Read` — 拆分 `d.Id()`（3 段，`len(idSplit) != 3` 报错）；填 ZoneId/FunctionId + `Filters = []*teov20220901.AdvancedFilter{{Name: "replica-name", Values: []*string{helper.String(replicaName)}}}` + `Limit = helper.Int64(200)`；`resource.Retry(tccommon.ReadRetryTimeout, ...)`；遍历结果精确匹配 `*replica.ReplicaName == replicaName`；未找到时**先** `log.Printf("[CRUD] ...")` 保留现场**再** `d.SetId("")`；set 字段前判 nil
4. `resourceTencentCloudTeoFunctionReplicaV4Update` — 拆 ID；`needChange` 检查 `mutableArgs := []string{"content", "remark"}` 的 `d.HasChange`；变化时调 ModifyFunctionReplica；末尾调 Read
5. `resourceTencentCloudTeoFunctionReplicaV4Delete` — 拆 ID；`request.ReplicaNames = []*string{helper.String(replicaName)}`；`resource.Retry(tccommon.WriteRetryTimeout, ...)`

关键 import 模式：
```go
import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

单测骨架（gomonkey mock，参照 `resource_tc_teo_function_replica_test.go`）：
- `mockMeta`（实现 `tccommon.ProviderMeta` 的 `GetAPIV3Conn()`）、`newMockMeta()` 返回 `&mockMeta{client: &connectivity.TencentCloudClient{}}`
- `patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)` 挂 mock client
- `patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *...) (*..., error) {...})` 逐接口 mock，闭包内 assert 请求字段
- `schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` 构造 `d`，`d.SetId(...)` 后调 `res.Create/Read/Update/Delete(d, meta)` 并断言结果

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_replica_v4` 资源，完整覆盖四个云 API 的参数映射（含 v1 缺失的 `created_on`/`modified_on` Computed 字段）
- 遵循 provider 现有代码风格（`tencentcloud_igtm_strategy` / 版本化资源 `tencentcloud_teo_l7_acc_rule_v2` 模式）
- 支持复合 ID import：`zone_id#function_id#replica_name`
- 提供基于 gomonkey 的单元测试（不用 TF 验收测试套件）
- 在 provider.go/provider.md 注册，.md 示例文件供 `make doc` 生成 website 文档

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_function_replica`（v1）资源的 schema 与行为（保持向后兼容）
- 不实现对应的 data source
- 不处理异步轮询（四个接口均为同步接口，无 TaskId/FlowId）
- 不将 DescribeFunctionReplicas 的查询类参数（`sort_by`/`sort_order`/`filters`/`offset`/`limit`）暴露为资源 schema 字段——Read 内部固定使用 `Filters(replica-name)` + `Limit=200` 定位目标副本
- 不支持批量删除（Delete 接口虽为批量语义，TF 资源粒度为单个副本，仅传单元素列表）

## Decisions

### D1. 新增 `_v4` 版本化资源而非修改 v1
**决策**: 新建独立资源 `tencentcloud_teo_function_replica_v4`，独立文件 `resource_tc_teo_function_replica_v4.go`，v1 完全不动。

**备选**: 直接在 v1 资源上新增 `created_on`/`modified_on` Computed 字段。

**理由**:
- 需求明确要求新增 `tencentcloud_teo_function_replica_v4`（版本化命名）
- provider 已有 `tencentcloud_teo_l7_acc_rule_v2`、`tencentcloud_ccn_attachment_v2` 等版本化资源先例
- 纯新增对存量 state 零影响，回滚只需 revert 新文件与注册行

### D2. 资源 ID 设计：三段联合 ID
**决策**: `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))`，即 `zone_id#function_id#replica_name`（`FILED_SP = "#"`）。

**理由**:
- CreateFunctionReplica 响应仅含 RequestId，无服务端生成 ID；副本由 `(zone_id, function_id, replica_name)` 三元组唯一标识（同一 FunctionId 下副本名唯一）
- Delete/Modify 均需要这三项入参，Read 需要前两项 + replica-name 过滤；三段 ID 使各 CRUD 方法可直接从 `d.Id()` 拆出全部所需参数（import 后亦然）
- 与 v1 资源及 `tencentcloud_teo_alias_domain` 等联合 ID 资源保持一致

### D3. Schema 字段设计
**决策**:

| Schema 字段 | 类型 | 属性 | 说明 |
|---|---|---|---|
| `zone_id` | TypeString | Required, **ForceNew** | 站点 ID |
| `function_id` | TypeString | Required, **ForceNew** | 函数 ID |
| `replica_name` | TypeString | Required, **ForceNew** | 副本名称（定位标识） |
| `content` | TypeString | Required | 副本内容（JavaScript 代码） |
| `remark` | TypeString | Optional | 副本描述 |
| `created_on` | TypeString | **Computed** | 创建时间（来自 DescribeFunctionReplicas） |
| `modified_on` | TypeString | **Computed** | 更新时间（来自 DescribeFunctionReplicas） |

**理由**:
- `zone_id`/`function_id` 是副本归属标识，`replica_name` 是 Modify/Delete 的定位键（API 不支持改名），三者改变均应触发重建 → ForceNew
- `content`/`remark` 均可经 ModifyFunctionReplica 原地更新 → 不设 ForceNew
- `created_on`/`modified_on` 仅存在于 Describe 响应（Create/Modify 请求无此字段）→ Computed-only，这是 v4 相对 v1 的核心增强
- DescribeFunctionReplicas 的 `SortBy`/`SortOrder`/`Filters` 是查询辅助参数，不描述资源本身属性，不暴露到 schema（v1 亦如此）

### D4. Read 实现方式：列表接口 + replica-name 过滤 + 精确匹配
**决策**: 构造 `DescribeFunctionReplicasRequest{ZoneId, FunctionId, Filters: [{Name: "replica-name", Values: [replicaName]}], Limit: 200}`，遍历 `FunctionReplicas` 精确匹配 `*replica.ReplicaName == replicaName`。

**理由**:
- 无单数 `DescribeFunctionReplica` 接口，只能用列表接口定位
- `Limit=200` 是云 API 注释标注的最大值（分页字段取最大值符合项目规范）；单函数副本上限极小（默认 2 个），一页必然覆盖
- Filters 的 replica-name 支持模糊查询 → 必须再精确匹配兜底
- 未找到（response/Response 为空、列表为空、无匹配项）时：**先** `log.Printf("[CRUD] read teo function replica v4 id=%s", d.Id())` 保留现场，**再** `d.SetId("")` 返回 nil（资源型资源的标准行为，区别于 DATASOURCE 资源的 NonRetryableError 约定）

### D5. Delete 实现方式：单元素列表
**决策**: `request.ReplicaNames = []*string{helper.String(replicaName)}`。

**理由**: DeleteFunctionReplica 为批量删除接口（接受名称列表），但 TF 资源粒度为单个副本；与 v1 及 `tencentcloud_teo_alias_domain`（`AliasNames = []*string{...}`）模式一致。不把 `replica_names` 暴露为 schema 列表字段——资源语义是管理单个副本，ID 为三段联合 ID。

### D6. Update 实现方式：needChange + mutableArgs 模式
**决策**: `mutableArgs := []string{"content", "remark"}`，任一 `d.HasChange` 才调用 ModifyFunctionReplica（两个非 ForceNew 字段都传入请求），末尾调 Read。

**理由**:
- Modify 接口 Content/Remark 均为可选入参，无变化时跳过调用可避免无意义写操作
- 与 v1 资源及 igtm_strategy 的 `needChange` 模式保持一致
- Modify 请求始终填充 ZoneId/FunctionId/ReplicaName 三个定位字段（从 `d.Id()` 拆出）

### D7. retry 拓扑
**决策**: 所有 SDK 调用（`*WithContext` 变体）均包裹 `resource.Retry`：
- Create/Modify/Delete → `tccommon.WriteRetryTimeout`
- Read → `tccommon.ReadRetryTimeout`
- 错误统一 `tccommon.RetryError(e)` 包装；retry 块内校验响应非空（`result == nil || result.Response == nil` → `resource.NonRetryableError`）
- `d.SetId()` 等成功操作一律放在 retry 块外、reqErr 判空之后

**理由**: 符合项目“调用云API接口需以 ReadRetryTimeout/WriteRetryTimeout 作为超时时间添加 retry”的硬性规范；四个接口均为同步接口，无嵌套轮询（无 retry-in-retry）。

### D8. 单元测试：gomonkey mock（不用 TF 验收测试套件）
**决策**: 新建 `resource_tc_teo_function_replica_v4_test.go`（package teo_test），使用 gomonkey：
- mock `UseTeoV20220901Client` 返回裸 `&teov20220901.Client{}`
- mock 四个 `*WithContext` 方法，闭包内 assert 请求字段并返回固定响应
- 用例覆盖：Create（全参数）、Create（无 remark）、Read（找到并回读 created_on/modified_on）、Read（未找到 → SetId("")）、Update（content+remark 变更）、Delete（单元素 ReplicaNames）

**理由**: 按项目规范，新增 terraform 资源的单测使用 mock 方式而非 TF 测试套件；v1 资源的测试文件即此模式，可直接参照（注意 helper 函数命名需避免与同包 v1 测试文件冲突，如 `mockMetaFunctionReplicaV4`、`ptrStringFunctionReplicaV4`）。

### D9. 注册与文档
**决策**:
- `tencentcloud/provider.go` 的 teo 资源 map 中新增 `"tencentcloud_teo_function_replica_v4": teo.ResourceTencentCloudTeoFunctionReplicaV4()`（紧邻 v1 注册行）
- `tencentcloud/provider.md` 的 TEO Resources 列表新增 `tencentcloud_teo_function_replica_v4`（紧邻 v1 条目）
- 新建 `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.md`：一句话描述（带 TEO 产品名）+ Example Usage + Import（说明三段联合 ID 格式）；不手写 Argument/Attribute Reference（`make doc` 自动生成；website/ 目录禁止手改）

## Risks / Trade-offs

- [Risk] DescribeFunctionReplicas 的 replica-name 过滤支持模糊查询，可能返回多个结果 → **Mitigation**: Read 中遍历精确匹配 `*replica.ReplicaName == replicaName`，未匹配则按“资源已删除”处理（先 log 后 SetId("")）
- [Risk] Create 接口无返回 ID，若 Create 成功但后续 Read 失败，tfstate 中可能残留无对应资源的 ID → **Mitigation**: Create 内 `d.SetId()` 在 API 成功后立即执行（retry 外），末尾 Read 失败时用户可 `terraform destroy` 清理；Read 未找到会正常清空 ID
- [Risk] `content` 为大字段（最大 5MB），存储于 state 可能膨胀 → **Mitigation**: 与 v1/`tencentcloud_teo_function` 行为一致，属于产品语义，接受该 trade-off；不设 Sensitive（非敏感内容）
- [Trade-off] `replica_name` ForceNew 意味着改名需销毁重建 → 由 API 限制决定（Modify 接口用 replica_name 定位副本，不支持改名）
- [Trade-off] v4 与 v1 资源管理同一云对象，同时声明两者会产生外部冲突（同 zone/function/replica_name） → 文档 .md 中以 NOTE 提示用户二选一使用

## Migration Plan

纯新增变更，无 state 迁移：
1. 合入新文件 `resource_tc_teo_function_replica_v4.go` / `_test.go` / `_v4.md` 与 provider.go/provider.md 注册行
2. `make doc` 生成 `website/docs/r/teo_function_replica_v4.html.markdown`（由收尾阶段 tfpacer-finalize skill 统一执行）
3. 用户按需在配置中新增 `resource "tencentcloud_teo_function_replica_v4" "xxx" { ... }`

回滚策略：revert 新增文件与两处注册行即可，无 state 变更需撤销。

## Open Questions

- 无。四个云 API 已在 vendor 中核实（含 WithContext 变体与完整入参/出参 struct），均为同步接口，所有 schema 字段与 API 参数一一对应。
