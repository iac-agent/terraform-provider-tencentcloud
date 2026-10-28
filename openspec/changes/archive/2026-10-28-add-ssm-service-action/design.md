## Context

### Repository baseline (probed, facts not assumptions)

| Fact | Value |
| --- | --- |
| Framework helper package | `tencentcloud/framework/fw` **exists** (`action_with_configure.go`, `data_source_with_configure.go`, `ephemeral_resource_with_configure.go`, `resource_with_configure.go`, `with_import_by_id.go`, `with_meta.go`, `with_no_update.go`) |
| Available helpers | `fw.ActionWithConfigure` (embeds `withMeta` → `Client() *connectivity.TencentCloudClient`), `fw.ResourceWithConfigure`, `fw.DataSourceWithConfigure`, `fw.EphemeralResourceWithConfigure`, `fw.WithImportByID`, `fw.WithNoUpdate` |
| Framework version | `github.com/hashicorp/terraform-plugin-framework v1.19.0` (native `action` package vendored at `vendor/github.com/hashicorp/terraform-plugin-framework/action/`) |
| timeouts module | `github.com/hashicorp/terraform-plugin-framework-timeouts v0.7.0` → `actiontimeouts "…-timeouts/action/timeouts"` |
| Registry | `tencentcloud/framework/registry.go`; `actionFactories` already holds 6 entries; `ssm` package is **already imported** there (used for resource + ephemeral) |
| SDKv2 provider registration | **not touched** for framework references. Note: `tencentcloud/framework/provider.go` has **no** `CreateContext` bridge, so a framework reference cannot reuse an SDKv2 `schema.Resource` CRUD function; its `invoke` logic must live in the framework file itself |
| Action precedents | `services/teo/action_tc_teo_confirm_origin_acl_update.go`, `services/bdrc/action_tc_bdrc_run_copy_pair_tasks.go`, `services/dlc/action_tc_dlc_initialize_tc_lake.go`, `services/mongodb/action_tc_mongodb_restore_db_instance.go`, `services/antiddos/action_tc_antiddos_unblock_resources.go`, `services/tco/action_tc_organization_provision_role_configuration.go` |
| SSM service layer | `services/ssm/service_tencentcloud_ssm.go` → `NewSsmService(client *connectivity.TencentCloudClient) SsmService`, client accessor `me.client.UseSsmClient()`; SDKv2 three-piece convention (`ratelimit.Check`, `[CRITAL]` defer, `[DEBUG]` success log) |
| Existing status read | `SsmService.DescribeSsmServiceStatusByFilter(ctx) (*ssm.GetServiceStatusResponseParams, errRet error)` already wraps `GetServiceStatus` — reused as the convergence read, no new read wrapper is needed |
| Retry budgets | `tccommon.WriteRetryTimeout` (default 5m, `PROVIDER_WRITE_RETRY_TIMEOUT`), `tccommon.ReadRetryTimeout` (default 3m) from `tencentcloud/common/common.go` |
| Existing related refs | SDKv2 data source `tencentcloud_ssm_service_status` (`data_source_tc_ssm_service_status.go`) already surfaces the same `GetServiceStatus` fields — do **not** modify it |

### Resource kind self-check (mandatory, done before reading templates)

Hard criteria applied to the supplied API set `{CreateService, GetServiceStatus}`:

- Create-type trigger API present → **not** `RESOURCE_KIND_DATASOURCE`.
- The create side is a **trigger/one-shot** call: `CreateService` activates an account-level service,
  it has no resource identity of its own in Terraform's sense (the returned `ResourceId` is `ssm_<AppId>`,
  an account-level handle, and `GetServiceStatus` is a *status probe* with an **empty request**, not a
  `DescribeXxxById` lookup). No Delete/Modify API exists (`GetServiceStatus` + `CreateService` only).
- Therefore `Create + Describe + Delete` (GENERAL) and `Modify + Describe` (CONFIG) do **not** apply;
  there is no attach/detach pair either.
- Matches `RESOURCE_KIND_OPERATION`: "trigger a one-shot far-side side effect, no lifecycle state".

**Conclusion: no contradiction — the declared `RESOURCE_KIND_OPERATION` is correct.** No type correction is
performed, so the pre-selected template `references/5-action.md` remains the applicable template
(`RESOURCE_KIND_RESOLVED: NONE`).

### API contract excerpted from vendor (implementation stage does not need to re-read `models.go`)

Package: `ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"`

`models.go` L449–513:

```go
type CreateServiceRequestParams struct {
	// Only PostPay is supported, case sensitive. Enum: PostPay (pay-as-you-go).
	PayMode *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`
	// Client-generated creation request identifier. Must reuse the same value when retrying
	// the same logical creation. 1–64 visible ASCII chars ('!'..'~'); no space / Chinese /
	// control chars. Suggested: ssm-<UUID>. Required. SSM forwards it to billing and does not
	// store or return it.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}

type CreateServiceRequest struct {
	*tchttp.BaseRequest
	PayMode     *string `json:"PayMode,..."`
	ClientToken *string `json:"ClientToken,..."`
}
// FromJsonString rejects unknown keys: only "PayMode" and "ClientToken" are accepted.

type CreateServiceResponseParams struct {
	// SSM service resource ID: currently "ssm_" + AppId, e.g. ssm_1250000000. May be empty.
	// NOTE: the resource may still be provisioning on return; it does not mean the service
	// is already usable.
	ResourceId *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`
	// Billing order ID used to reconcile the activation order. No fixed format promised.
	OrderId    *string `json:"OrderId,omitnil,omitempty" name:"OrderId"`
	RequestId  *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateServiceResponse struct {
	*tchttp.BaseResponse
	Response *CreateServiceResponseParams `json:"Response"`
}
```

`models.go` L1378–1450 (`GetServiceStatusRequest` has **no** fields):

```go
type GetServiceStatusRequest struct { *tchttp.BaseRequest }   // no parameters

type GetServiceStatusResponseParams struct {
	ServiceEnabled          *bool   // true: activated, false: not activated
	InvalidType             *int64  // 0 not purchased, 1 normal, 2 suspended (arrears), 3 resource released
	AccessKeyEscrowEnabled  *bool   // key secure-hosting available
	ExpireTime              *string
	QPSLimit                *int64
	SecretLimit             *int64
	PayModel                *string
	RenewFlag               *int64  // 0 manual, 1 auto, 2 no renewal
	ResourceId              *string
	TotalCount              *int64
	ResourceRegion          *int64
	RequestId               *string
}
```

Client signatures (`client.go`): `CreateService(request)`, `CreateServiceWithContext(ctx, request)`,
`GetServiceStatus(request)`, `GetServiceStatusWithContext(ctx, request)`
(L326/L339/L1335/L1345).

Error code (`errors.go` L81): `ssm.RESOURCEINUSE_SERVICEALREADYENABLED = "ResourceInUse.ServiceAlreadyEnabled"`
— this is the "already at target state" signal for idempotent re-invoke.

> Parameter-mapping caveat (recorded in the spec): every response-side name suggested by the change request
> is `Computed`-only in nature. A Framework **action has no output attributes at all**
> (`action.InvokeResponse` only carries `Diagnostics` + `SendProgress`), so `resource_id`, `order_id`,
> `service_enabled`, `invalid_type`, `access_key_escrow_enabled`, `expire_time`, `qps_limit`,
> `secret_limit`, `pay_model`, `renew_flag`, `resource_region` can only be **logged**, not landed in state.
> They are therefore implemented as observability logging + (for a subset) as preconditions in the
> convergence loop, not as schema attributes.

### Template key structure (`references/5-action.md`, applicable form in this baseline)

The template is directly applicable because `tencentcloud/framework/fw` exists and `action` precedents exist.

```go
package ssm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ action.ActionWithConfigure = &SsmService{}

func NewSsmServiceAction() action.Action { return &SsmService{} }

type SsmService struct{ fw.ActionWithConfigure }

type SsmServiceModel struct {
	PayMode     types.String `tfsdk:"pay_mode"`
	ClientToken types.String `tfsdk:"client_token"`
	Timeouts    actiontimeouts.Value `tfsdk:"timeouts"` // only when a wait budget is exposed
}

func (a *SsmService) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_service"
}

func (a *SsmService) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to activate … CreateService … one-time; no state is persisted; " +
			"repeated invokes re-trigger the remote call; destroy does not deactivate.",
		Attributes: map[string]schema.Attribute{ /* pay_mode Required, client_token Optional+Computed */ },
		Blocks:     map[string]schema.Block{ /* optional "timeouts" */ },
	}
}

func (a *SsmService) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	// 1. req.Config.Get → tri-state input validation (IsNull/IsUnknown/empty)
	// 2. a.Client() nil guard
	// 3. call CreateService (service layer, retry inside)
	// 4. poll GetServiceStatus until ServiceEnabled == true, budget = tccommon.WriteRetryTimeout
	// 5. all failures → resp.Diagnostics.AddError
}
```

Lifecycle skeleton for an operation-type action: **only** `Metadata` / `Schema` / `Configure` (promoted from
`fw.ActionWithConfigure`) / `Invoke`. No `Create` / `Read` / `Update` / `Delete`, no `Importer`,
no `ModifyPlan` (unless config validation at plan time is needed).

## Goals / Non-Goals

**Goals:**

- Activate SSM (Secrets Manager) for the invoked account/region from Terraform, idempotently enough to be
  re-runnable after a partial failure.
- Wait for real convergence: return success only after `GetServiceStatus.ServiceEnabled == true`,
  because `CreateService` explicitly documents "may still be provisioning".
- Keep the framework action stateless: no id, no output attributes, no destroy semantics.
- Reuse the existing SSM service layer and its SDKv2 conventions (three-piece logging, `ratelimit.Check`,
  `tccommon.RetryError`), so the action behaves like the rest of the SSM code.
- Ship registration, provider.md entry, example `.md`, gomonkey unit test, and let `make doc` generate
  `website/docs/a/ssm_service.html.markdown`.

**Non-Goals:**

- No deactivation / cancellation on destroy (the API surface has no such operation).
- No management of the billing order (`OrderId` is logged only) and no `pay_mode` other than `PostPay`.
- No new SSM secrets-manager resource; no change to `tencentcloud_ssm_service_status`.
- No SDK/vendor upgrade in this change (both APIs are already vendored).
- No `terraform import`; an action cannot be imported.

## Decisions

### Decision 1: Framework native Action, not an SDKv2 pseudo-resource

**Chosen:** `action.Action` implemented in `tencentcloud/services/ssm/action_tc_ssm_service.go`,
factory `NewSsmServiceAction()`, registered in `framework/registry.go` `actionFactories`.

**Alternatives:** (a) SDKv2 `resource_tc_ssm_service_operation.go` pseudo resource (the older
`*_operation.go` pattern, e.g. `tencentcloud_ssm_rotate_product_secret`);
(b) an SDKv2 resource with `CustomizeDiff` and a no-op delete.

**Rationale:** the operation has no identity, no read-back identity, and no delete; a pseudo resource would
store a meaningless id and imply destroy semantics that do not exist. Six native action precedents already
ship in this repo, and `fw.ActionWithConfigure` is available, so the native form is the least surprising.
Registration **must not** touch `tencentcloud/provider.go` (framework refs live only in `registry.go`).

> Naming caveat to settle first: the `ssm` package **already declares**
> `func NewSsmService(client *connectivity.TencentCloudClient) SsmService` in
> `service_tencentcloud_ssm.go`. Go has no overloading, so the framework factory **cannot** be named
> `NewSsmService()` in the same package — that would be a duplicate declaration and a hard compile error.
> The factory is therefore named **`NewSsmServiceAction()`** (documented deviation from the template's
> "no `Action` suffix" rule, forced by the collision). `TypeName` remains `tencentcloud_ssm_service` and the
> registry line is `ssm.NewSsmServiceAction,`. The alternative (renaming the existing service constructor)
> is rejected: it would touch every SSM file for a cosmetic reason.

### Decision 2: Schema shape

| Attribute | Type | Required/Optional | Source | Rationale |
| --- | --- | --- | --- | --- |
| `pay_mode` | String | **Required** (schema-level `RequiresReplace` not applicable — actions have no plan/replace) | `CreateServiceRequest.PayMode` | Only `PostPay` is accepted by the API; keep it explicit so the user sees the billing mode. Runtime validation in `Invoke` (see Decision 5) rather than a plan-time validator, matching sibling actions which validate in `Invoke`. |
| `client_token` | String | **Optional + Computed** | `CreateServiceRequest.ClientToken` | Not hard-coding a Default (provider rule: cloud/default values are not hard-coded as `Default`); `Optional+Computed` lets the action generate `ssm-<token>` when the user omits it, while still allowing an explicit replay token. Note: for an action the computed value is evaluated at invoke time and is **not** persisted, so omitting it yields a fresh token per invoke (documented). |

`client_token` generation detail (verified): `tencentcloud/internal/helper.BuildToken()` returns
`base64.StdEncoding.EncodeToString(16 random bytes)` = 24 characters from `A–Za–z0–9+/`, i.e. plain ASCII
without padding-free issues and within the 1–64 visible-ASCII window of the API. It must be prefixed with
`ssm-` (API suggestion) — `fmt.Sprintf("ssm-%s", helper.BuildToken())` still fits 1–64 chars and matches the
recommended `ssm-<UUID>` shape.
| `timeouts` | SingleNestedBlock | Optional | `actiontimeouts.BlockWithOpts` (`invoke` only) | Follows `mongodb_restore_db_instance` / `tco_provision_role_configuration`: the action waits for an async convergence, so the wait budget should be overridable. Default from `defaultSsmServiceInvokeTimeout = 10 * time.Minute`. |

No `resource_id` / `order_id` / `service_enabled` / … attributes: Framework actions cannot publish outputs
(`InvokeResponse` has no state/output field). They are logged at `[DEBUG]` level instead.

### Decision 3: Service layer stays SDKv2-style; a new `CreateSsmService` is added

`service_tencentcloud_ssm.go` gains:

```go
func (me *SsmService) CreateSsmService(ctx context.Context, payMode, clientToken *string) (resp params, errRet error)
```

with the three mandatory conventions copied from the neighbouring methods (verified in-file:
`ratelimit.Check` 21 occurrences, `[DEBUG]` 22, `[CRITAL]` 10):

```go
logId   = tccommon.GetLogId(ctx)
request = ssm.NewCreateServiceRequest()
defer func() { if errRet != nil { log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error()) } }()
ratelimit.Check(request.GetAction())
response, err := me.client.UseSsmClient().CreateService(request)
log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
```

The convergence read reuses `DescribeSsmServiceStatusByFilter(ctx)` **unchanged** (it already wraps
`GetServiceStatus` with the same conventions). No context adaptation is needed: framework actions hold
`*connectivity.TencentCloudClient` directly (via `fw.ActionWithConfigure.Client()`), so
`NewSsmService(a.Client())` works and `tccommon.GetLogId(ctx)` reads the framework request `ctx`
(same pattern as `NewSsmService(r.Client())` in `resource_tc_ssm_secret_version_v2.go`).

### Decision 4: Retry placement and error wrapping

- The `CreateService` call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, …)` returning
  `tccommon.RetryError(e)` on failure, per provider rules; **no** nested retry.
- **`ClientToken` must be generated exactly once, before entering the retry block**, so every retry inside
  one invoke reuses the same token (API contract: the same logical creation must reuse the same value
  when retried).
- Idempotent-error handling: if the retry ultimately returns `ResourceInUse.ServiceAlreadyEnabled`,
  it is **not** a failure — log it, `resp.SendProgress` a "service already enabled" message, and continue
  to the convergence check (a previous invoke, a manual console activation, or a retry after a lost response).
  Implemented as a package-level `func isSsmServiceAlreadyEnabled(err error) bool` so it is unit-testable
  (detect via `*tcerr.TencentCloudSDKError` code comparison and/or substring fallback).
- Anything else → `resp.Diagnostics.AddError("activating SSM service (pay_mode=%s)", err)`.

### Decision 5: Convergence loop (action-appropriate async rules)

```go
// package-level so unit tests can shorten it
var ssmServicePollInterval = 5 * time.Second
```

- Loop until `GetServiceStatus.ServiceEnabled == true` **or** budget (`timeouts { invoke = … }`, default
  10m) exhausted; sleep `ssmServicePollInterval` between attempts.
- Each attempt uses `DescribeSsmServiceStatusByFilter` (already retried internally at the SDK-call level;
  no extra nesting).
- Progress reporting: `if resp.SendProgress != nil { resp.SendProgress(action.InvokeProgressEvent{Message: …}) }`
  before/after the create and on each poll tick.
- Treat `ServiceEnabled == false` but `InvalidType == 2` (suspended due to arrears) or `3` (resource released)
  as a **terminal** failure with a clear message rather than burning the whole budget.
- On timeout: `AddError("SSM service did not converge", fmt.Sprintf("last observed: service_enabled=%v, invalid_type=%d, resource_id=%s", …))`
  — the last snapshot must be echoed, never a bare "timeout".
- `ResourceId` returned by `CreateService` and by `GetServiceStatus` are logged for reconciliation;
  a mismatch (when both non-empty) is a **warning** via `resp.Diagnostics.AddWarning`, not an error,
  because the pay mode / AppId mapping is cloud-owned.

### Decision 6: `pay_mode` / `client_token` validation location

Validated in `Invoke` (not via `ValidateConfig`), consistent with every action precedent
(`bdrc`, `teo`, `mongodb` all validate in `Invoke` and `AddError`):

- `pay_mode` null/unknown/empty → `Missing pay_mode`.
- `pay_mode != "PostPay"` → `Invalid pay_mode`, detail "only `PostPay` is supported by the CreateService API".
- `client_token` (when set) must match `^[\x21-\x7E]{1,64}$` → otherwise `Invalid client_token`.
  When null/unknown → generate `helper.BuildToken()` (already prefixed by the provider helper) once.

### Decision 7: Documentation & naming

- TF name `tencentcloud_ssm_service`; file `action_tc_ssm_service.go` (product = directory name `ssm`,
  `<name>` = `service`).
- Factory **`NewSsmServiceAction()`** — the template's "no `Action` suffix" preference cannot be honoured here,
  because the `ssm` package already declares
  `func NewSsmService(client *connectivity.TencentCloudClient) SsmService` (Go has no overloading, so
  `NewSsmService()` would be a duplicate declaration / compile error). Decision: keep the framework factory as
  `NewSsmServiceAction()` and keep `TypeName = "tencentcloud_ssm_service"`; the registry entry is
  `ssm.NewSsmServiceAction,`. Renaming the existing service constructor is rejected (touches every SSM file
  for cosmetics). See Decision 1 for the same note.
- `<stem>.md` example: one-line description mentioning SSM, `Example Usage` with an `action` block +
  `terraform apply -invoke=action.tencentcloud_ssm_service.example`, a NOTE for Terraform >= 1.14 and
  a NOTE that repeated invoke re-triggers the activation and destroy does not deactivate. **No** Import
  section, **no** Argument/Attribute Reference sections (generated by `make doc`).
- `website/docs/a/ssm_service.html.markdown` is generated by `make doc` only.
- `provider.md`: add under `Secrets Manager(SSM)`, after `Ephemeral Resource`, an `Action` node with
  `tencentcloud_ssm_service`.

### Decision 8: Unit test strategy (gomonkey, no acceptance suite)

`action_tc_ssm_service_test.go` in package `ssm_test`, modelled on
`action_tc_bdrc_run_copy_pair_tasks_test.go`:

- Build `action.InvokeRequest` by hand: `tftypes.NewValue(tftypes.Object{AttributeTypes: …}, …)` with
  `pay_mode`, `client_token`, and (when the timeouts block exists) a `timeouts` object value; obtain the
  schema from `a.Schema(...)`.
- Inject the client via `meta := &sharedmeta.ProviderMeta{Client: client}` + `a.Configure(...)`.
- `patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)`; then
  `patches.ApplyMethodFunc(ssmClient, "CreateService", …)` and
  `patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", …)`.
- Set `ssmServicePollInterval = time.Millisecond` (package var) inside the convergence tests to keep them fast.
- Cases: success (poll returns `ServiceEnabled=true` immediately) / create API error → diagnostic /
  `ResourceInUse.ServiceAlreadyEnabled` → success without failing / `pay_mode` missing or invalid and
  `client_token` malformed → diagnostic + no API call / provider client nil → diagnostic /
  never-converging poll with an exhausted budget → timeout diagnostic containing the last snapshot.
- Test names `TestSsmServiceAction_*`; run hint comment
  `// go test ./tencentcloud/services/ssm/ -run "TestSsmServiceAction" -v -count=1 -gcflags="all=-l"`.

## Risks / Trade-offs

- **Repeated invoke re-activates (billing)** → The action places a PostPay activation order each time it is
  invoked; a stable `client_token` replayed by the user deduplicates on the billing side. Mitigation:
  document it in the schema `Description`, the example `.md`, and treat `ResourceInUse.ServiceAlreadyEnabled`
  as success so re-running is safe.
- **Lost response after a successful order** (SDK timeout while the order was already placed) → the retry
  inside one invoke reuses the same `ClientToken`, and a subsequent `ServiceAlreadyEnabled` is accepted.
  Cross-invoke duplicates can only be avoided by the user pinning `client_token`; documented.
- **`ClientToken` is not persisted** (action has no state) → omitting it produces a new token per invoke.
  Mitigation: `Optional + Computed` semantics + explicit note in the docs that a pinned value should be
  provided to make retries replayable.
- **No output attributes** → users cannot read `resource_id` / `order_id` from state; they appear only in
  `TF_LOG=DEBUG` output. Mitigation: log at `[DEBUG]`, and point users to
  `data.tencentcloud_ssm_service_status` for post-activation introspection.
- **Convergence may be slow** (billing order → provisioning) → budget defaults to 10 minutes and is
  overridable via `timeouts { invoke = "…" }`; timeout error reports the last observed snapshot so the
  user can tell "not yet effective" from "activated but predicate wrong".
- **Account/region scope** → activation is per account+region; running the action in another region's
  provider alias activates that region. Documented; no code-level guard (the API is region-scoped by
  endpoint already).
- **Symbol collision** `NewSsmService` → resolved by naming the factory `NewSsmServiceAction()`
  (see Decision 1 / Decision 7); a wrong registry symbol would be caught immediately by the compiler.

## Migration Plan

- Purely additive: new files + one registry line + one `provider.md` line + service-layer method.
- No state migration; existing SSM resources and the `tencentcloud_ssm_service_status` data source are
  untouched, so no plan drift is possible.
- Rollback: remove the action file, registry line, `provider.md` line, example `.md`, unit test, and the
  `CreateSsmService` service method. Nothing else references them.
- Release-note type: `new-action`.

## Implementation checklist (the template's "Implementation Notes" carried over)

1. Only `Invoke` (+ `Metadata` / `Schema` / promoted `Configure`); no `Read` / `Update` / `Delete` /
   `ImportState`, no id.
2. Input validation (`IsNull` / `IsUnknown` / empty / enum / charset) at the top of `Invoke`, all through
   `resp.Diagnostics.AddError` (never silent return, never panic).
3. Never write state; never set output attributes (impossible for an action — do not try).
4. Poll budget from `tccommon.WriteRetryTimeout` (or the `timeouts.invoke` override), poll interval as a
   package-level `var`, progress via `resp.SendProgress` guarded by a nil check.
5. Idempotent error code (`ResourceInUse.ServiceAlreadyEnabled`) → success path with explicit log/progress.
6. Timeout diagnostics must include the last observed snapshot.
7. Register in `framework/registry.go` `actionFactories` only — do **not** touch `tencentcloud/provider.go`.
8. All produced artifacts (Go code, comments, `.md`, tests) in English.
