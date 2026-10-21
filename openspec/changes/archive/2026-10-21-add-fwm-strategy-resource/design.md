# Design: add-fwm-strategy-resource

## Overview

New Terraform Plugin Framework resource `tencentcloud_fwm_strategy` (RESOURCE_KIND_GENERAL) managing TencentCloud FWM 下发规则 (strategy) with full CRUD + import. Implementation follows template `~/.codebuddy/skills/tfpacer-tf-templates/references/1-resource-general.md` and the common rules in `references/6-field-types-and-computed.md`. Reference implementation style: `tencentcloud/services/ssm/resource_tc_ssm_secret_version_v2.go` (Configure via `*sharedmeta.ProviderMeta`).

## Cloud API struct definitions (from vendor, verified)

Package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` (`vendor/.../fwm/v20250611/models.go`).

### CreateStrategyRequest (models.go L886-903)
```go
type CreateStrategyRequest struct {
	*tchttp.BaseRequest
	Product        *string          // 产品类型 (required by cloud)
	ReceiveAccount []*string        // 下发账号 (list of uins)
	PreStrategy    []*StrategyReq   // 前区规则组
	PostStrategy   []*StrategyReq   // 后区规则组
	ReceiveGroup   []*string        // 下发账号组
}
// CreateStrategyResponseParams: StrategyIds []*string  (L929-935)
```

### StrategyReq (models.go L4449-4455)
```go
type StrategyReq struct {
	GroupId  *string // 规则组Id
	Sequence *int64  // 优先级
}
```

### DescribeStrategyRequest / Response (models.go L2636-2674)
```go
type DescribeStrategyRequest struct {
	*tchttp.BaseRequest
	StrategyId *string
}
// DescribeStrategyResponseParams: Strategy *StrategyResp
```

### StrategyResp (models.go L4457-4505)
```go
type StrategyResp struct {
	StrategyId     *string            // 策略Id
	GroupId        *string            // 规则组Id
	GroupName      *string            // 规则组名称
	RuleCount      *int64             // 规则数
	RuleStatus     *int64             // 策略状态
	ReceiveAccount []*ReceiveAccount  // 下发账号 (objects in response)
	Sequence       *int64             // 优先级
	ErrMsg         *string            // 下发失败原因
	ErrorType      *string            // 下发失败原因类型
	UpdateTime     *string            // 更新时间
	CreateBy       *string            // 创建人
	UpdateBy       *string            // 更新人
	ExecArea       *string            // 执行区域
	CreateName     *string            // 创建人名称
	UpdateName     *string            // 更新人名称
	CreateTime     *string            // 创建时间
}
```

### ReceiveAccount (models.go L4068-4080) and MemberInfo (L3049-3064)
```go
type ReceiveAccount struct {
	Uin          *string       // 租户 uin
	Nickname     *string       // 租户名称
	ReceiverType *int64        // 0=账号uin，1=账号组
	Members      []*MemberInfo // 只有ReceiverType 是 1 时才返回
}
type MemberInfo struct {
	AppId    *string
	Uin      *string
	Nickname *string
	MemberId *string  // 个人场景为空
	NodeName *string  // 所属部门
}
```

### ModifyStrategyRequest (models.go L3464-3481)
```go
type ModifyStrategyRequest struct {
	*tchttp.BaseRequest
	StrategyId     *string
	ReceiveAccount []*string
	Sequence       *int64
	GroupId        *string
	ReceiveGroup   []*string
}
```
Note: Modify has NO `Product`, `PreStrategy`, `PostStrategy` — these are create-only.

### DeleteStrategyRequest (models.go L1331-1336)
```go
type DeleteStrategyRequest struct {
	*tchttp.BaseRequest
	StrategyIds []*string
}
```

## Schema design (attribute mapping decisions)

| TF attribute | Type | Flags | Notes |
| --- | --- | --- | --- |
| `id` | String | Computed + `UseStateForUnknown()` | = strategy_id |
| `product` | String | Required, `RequiresReplace()` | Create-only; not in Modify → immutable |
| `receive_account` | List(String) | Optional | Request-facing uin list (Create/Modify `[]*string`) |
| `receive_group` | List(String) | Optional | Create/Modify `[]*string` |
| `pre_strategy` | SingleNested {`group_id` String Optional, `sequence` Int64 Optional} | Optional, `RequiresReplace()` | Create-only (`StrategyReq`) |
| `post_strategy` | SingleNested {`group_id` String Optional, `sequence` Int64 Optional} | Optional, `RequiresReplace()` | Create-only (`StrategyReq`) |
| `group_id` | String | Optional + Computed | Top-level in Modify request AND StrategyResp |
| `sequence` | Int64 | Optional + Computed | Top-level in Modify request AND StrategyResp |
| `strategy_id` | String | Computed + `UseStateForUnknown()` | From StrategyResp |
| `group_name`, `err_msg`, `error_type`, `update_time`, `create_by`, `update_by`, `exec_area`, `create_name`, `update_name`, `create_time` | String | Computed | From StrategyResp |
| `rule_count`, `rule_status` | Int64 | Computed | From StrategyResp |
| `receive_account_detail` | ListNested {`uin` String, `nickname` String, `receiver_type` Int64, `members` ListNested {`app_id`, `uin`, `nickname`, `member_id`, `node_name` String}} | Computed | Deviation note: response `StrategyResp.ReceiveAccount` is `[]*ReceiveAccount` objects, conflicting with request-side `receive_account` (list of strings) which shares the mapped SchemaName. A separate computed block is used so both sides are representable. |

All attributes have `Description` in English. No `Default` anywhere (rule: cloud defaults flow back via Read).

## ID semantics

- `Create` returns `Response.StrategyIds []*string` — take the first non-empty element as the resource id; if `Response` is nil or `StrategyIds` empty, log the logId and return `NonRetryableError` (must not write an empty id).
- All four APIs are synchronous (no FlowId/TaskId), so no async polling is required; Create finishes with an immediate Read-back.
- Read: `DescribeStrategy(StrategyId=id)`; if `response == nil || response.Response == nil || response.Response.Strategy == nil` → `log.Printf` with the id, `AddWarning`, `resp.State.RemoveResource(ctx)`.

## Template key structure (resource-general, from 1-resource-general.md)

- Struct embeds `fw.ResourceWithConfigure` + `fw.WithImportByID` (import passthrough) — follow the SSM v2 reference: `Configure` asserts `req.ProviderData.(*sharedmeta.ProviderMeta)` and stores `meta.Client` (`*connectivity.TencentCloudClient`).
- Compile-time assertions for `ResourceWithConfigure` / `ResourceWithImportState` only.
- `Metadata` sets `resp.TypeName = "tencentcloud_fwm_strategy"`.
- Lifecycle skeleton: `Create` (Plan.Get → build `CreateStrategyRequest` with null/unknown guards → call API → validate id → SetId → read back → refreshModel → State.Set), `Read` (State.Get → read() → nil ⇒ RemoveResource+Warning → refreshModel → State.Set), `Update` (Plan.Get → `ModifyStrategyRequest` with StrategyId from plan.Id → call API → read back → State.Set), `Delete` (State.Get → `DeleteStrategyRequest{StrategyIds: [id]}` → call API).
- Private `read(ctx, id)` helper returns `(*fwm.StrategyResp, error)`; nil when not found.
- Private `refreshModel(m *model, item *fwm.StrategyResp)` copies cloud-owned fields via `types.XxxPointerValue`.
- Value rules (6-field-types-and-computed.md): every read checks `IsNull()` AND `IsUnknown()`; SDK pointer writes use `types.XxxPointerValue()`; nested single object via `SingleNestedAttribute` with model struct (`req.Plan.Get` auto-expands); string lists via `types.List` + `ElementsAs`.
- Errors only via `resp.Diagnostics.AddError(title, detail)` including the API name and resource identifier; never panic or swallow errors.

## Registration & client accessor

1. `tencentcloud/framework/registry.go`: add `fwm.NewFwmStrategyResource` to `resourceFactories` (import `.../tencentcloud/services/fwm`).
2. `tencentcloud/provider.md`: register `tencentcloud_fwm_strategy` under the `fwm` node (add the node if absent, following the `tencentcloud_igtm_strategy` entry style).
3. `tencentcloud/connectivity/client.go`: add
   ```go
   // UseFwmClient returns FWM client for service
   func (me *TencentCloudClient) UseFwmClient() *fwm.Client {
       return &fwm.Client{Client: me.BaseClient}
   }
   ```
   with import `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`.

## Files

| File | Action |
| --- | --- |
| `tencentcloud/services/fwm/resource_tc_fwm_strategy.go` | new (Framework resource) |
| `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go` | new (gomonkey mock unit tests; no TF test suite; NOT executed) |
| `tencentcloud/services/fwm/resource_tc_fwm_strategy.md` | new doc: "Provides a resource to ..." + Example Usage + Import section (no Argument/Attribute Reference) |
| `tencentcloud/framework/registry.go` | add factory line |
| `tencentcloud/provider.md` | add resource entry |
| `tencentcloud/connectivity/client.go` | add `UseFwmClient` |

All generated content (Go code, comments, .md docs, tests) is English only. No `go build`/`go vet`; only `gofmt`/`make doc`/changelog in the finalize stage.

## Correctness checks

- Every Create request field exists in `CreateStrategyRequest`; every Update field in `ModifyStrategyRequest`; Read uses `DescribeStrategy`; Delete uses `DeleteStrategy` with `StrategyIds = [id]`. Verified against vendor structs above.
- Response fields set only after nil checks (pointer-safe `types.XxxPointerValue` handles nils).
- Update takes values from Plan, not old State.
