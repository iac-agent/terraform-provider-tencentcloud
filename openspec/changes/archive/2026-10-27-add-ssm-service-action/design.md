# Design: tencentcloud_ssm_service (Terraform Plugin Framework Action)

## Context

### Resource kind self-check (done before reading the framework templates)

The provided interface set is:

| API | Direction | Role |
| --- | --- | --- |
| `CreateService` | request/response | subscribes the account to SSM — creates **no trackable object** |
| `GetServiceStatus` | response only | returns an account-level subscription snapshot |

Applying the hard criteria:

- Not "Describe only" → not `RESOURCE_KIND_DATASOURCE`.
- `CreateService` does **not** return an object with a lifecycle; there is no `DeleteService` /
- not `RESOURCE_KIND_GENERAL`.
- Not "Modify + Describe only" → not `RESOURCE_KIND_CONFIG`.
- No Attach/Detach pair → not `RESOURCE_KIND_ATTACHMENT`.
- `CreateService` is a one-shot trigger whose effect is account-wide, irreversible, and whose re-run
  returns `ResourceInUse.ServiceAlreadyEnabled` → **`RESOURCE_KIND_OPERATION`** ✔

**Conclusion: no correction required.** The announced `RESOURCE_KIND_OPERATION` matches the interface
composition, so the action template (`references/5-action.md`) is the right one to follow. This is
recorded here for the sake of auditability: using a resource would force Terraform to track a phantom
object whose "destroy" cannot be implemented.

### Repository baseline (probed, authoritative)

| Probe | Result | Consequence |
| --- | --- | --- |
| `ls tencentcloud/framework/fw` | **absent** | Templates' `fw.ActionWithConfigure` / `a.Client()` cannot be used. Hand-write `Configure` on the `action` types and store `SsmService` in the struct. |
| `ls tencentcloud/services/*/action_tc_*.go` | **no match** | This repository has no framework action precedent; `tencentcloud_ssm_service` is the first. |
| `tencentcloud/framework/registry.go` `actionFactories` | **empty** | One factory entry must be added; `import .../services/ssm` already exists. |
| `tencentcloud/framework/provider.go` | `resp.ActionData = meta` already wired (line 346), `Actions()` returns `frameworkActions()` (line 381–383) | No provider plumbing change needed. |
| `github.com/hashicorp/terraform-plugin-framework` | `v1.19.0`; `action`, `action/schema`, `action/validator` packages vendored | Native action support available. |
| `connectivity.TencentCloudClient` | `UseSsmClient() *ssm.Client` (client.go:1101) | Reuse the shared SDKv2 client — no new credential parsing. |

### Reference implementation for the hand-written `Configure`

`tencentcloud/services/ssm/resource_tc_ssm_secret_version_v2.go` (lines 104–123) is the in-repo pattern:

```go
func (r *SsmSecretVersionV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	r.client = NewSsmService(meta.Client)
}
```

`sharedmeta.ProviderMeta` is `struct{ Client *connectivity.TencentCloudClient }`; the framework provider
sets `resp.ActionData = meta` in its own `Configure`, so the same assertion works for actions with
`req.ProviderData.(*sharedmeta.ProviderMeta)`.

## Goals / Non-Goals

**Goals**

- One-shot subscription of the current account to SSM through `CreateService`, exposed as
  `tencentcloud_ssm_service`.
- Correct handling of the asynchronous nature of `CreateService`: poll `GetServiceStatus` until
  `ServiceEnabled == true` and surface a diagnostic if it does not converge.
- Make the "already enabled" case usable rather than fatal where the API contract allows it (see D4).
- Zero Terraform state; the action only reports results through diagnostics and invoke output.

**Non-Goals**

- No managed resource, no `Read`/`Update`/`Delete`, no importer — nothing is tracked in state.
- No rollback on destroy: SSM subscription has no de-provisioning API. This must be stated in the docs.
- No new service-layer API wrapper beyond a thin addition to `service_tencentcloud_ssm.go` (see D2).
- No changes to `tencentcloud/provider.go` (framework actions are registered only in
  `framework/registry.go`).

## Decisions

### D1 — Terraform-facing name and file layout

Per the naming table in `tfpacer-tf-templates/SKILL.md`, the TF prefix must come from existing siblings.
All eleven SSM references already use the `tencentcloud_ssm_` prefix (`provider.md` lines 1177–1196), and
the product directory name equals the prefix, so the derived name is unambiguous:

| Item | Value |
| --- | --- |
| TF type name | `tencentcloud_ssm_service` |
| Go file | `tencentcloud/services/ssm/action_tc_ssm_service.go` |
| Doc example | `tencentcloud/services/ssm/action_tc_ssm_service.md` |
| Test | `tencentcloud/services/ssm/action_tc_ssm_service_test.go` |
| Factory | `NewSsmService()` → `action.Action` (**no** `Action` suffix) |
| Go type | `SsmServiceAction` (not exported; returned as `action.Action`) |
| Registration | `tencentcloud/framework/registry.go` → `actionFactories` |
| Package | `package ssm` |

Note: the Go type name must not collide with the existing service struct `SsmService`
(`service_tencentcloud_ssm.go`), hence `SsmServiceAction`.

The factory name `NewSsmService()` does collide with the *existing* package-level constructor
`func NewSsmService(client *connectivity.TencentCloudClient) SsmService` in
`service_tencentcloud_ssm.go`, because Go has no function overloading. The factory name is fixed by
the action naming convention (`New<Product><NameCamel>()`, no `Action` suffix), so the pre-existing
service constructor — an internal helper that is never part of the public `tencentcloud` package
surface — was renamed to `NewSsmServiceClient(...)`. All in-repo call sites were updated:
`resource_tc_ssm_secret_version_v2.go`, `ephemeral_tc_ssm_secret_version.go`,
`resource_tc_ssm_secret_test.go` and `resource_tc_ssm_secret_version_test.go`.

### D2 — Client injection and service-layer access

`framework/fw` does not exist, so `Configure` is hand-written (adaptation table column B):

```go
type SsmServiceAction struct {
	client SsmService
}

var (
	_ action.Action              = &SsmServiceAction{}
	_ action.ActionWithConfigure = &SsmServiceAction{}
)

func (a *SsmServiceAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	a.client = NewSsmService(meta.Client)
}
```

Two new methods are added to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`,
following the file's existing conventions (`ratelimit.Check(request.GetAction())`,
`log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`,
`defer func() { if errRet != nil { log.Printf("[CRITAL]...") } }()`):

```go
func (me *SsmService) CreateSsmService(ctx context.Context, payMode, clientToken string) (service *ssm.CreateServiceResponseParams, errRet error)

func (me *SsmService) DescribeSsmServiceStatusByFilter(ctx context.Context) (ServiceStatus *ssm.GetServiceStatusResponseParams, errRet error)
```

`DescribeSsmServiceStatusByFilter` **already exists** (`service_tencentcloud_ssm.go` lines 638–666) and is
reused as-is — it returns `response.Response` with a nil check, which is exactly what the polling loop
needs.

### D3 — Schema

The action has no managed attributes; every field is a *config* input or an *invoke-time* observation.
Framework action attributes therefore only carry the two inputs; the response fields are reported through
`resp.SendProgress` / diagnostics, because Framework actions have **no state** to write observed values
into (`InvokeResponse` exposes only `Diagnostics` and `SendProgress` — verified in
`vendor/.../action/invoke.go`).

| Attribute | Type | Required/Optional | API mapping | Notes |
| --- | --- | --- | --- | --- |
| `pay_mode` | `schema.StringAttribute` | **Required** | `CreateServiceRequest.PayMode` | Only `PostPay` is accepted (case-sensitive); `ssmServicePayModeValidator` |
| `client_token` | `schema.StringAttribute` | **Required** | `CreateServiceRequest.ClientToken` | 1–64 visible ASCII chars, no spaces/Chinese/control chars; `ssmServiceClientTokenValidator` |

`github.com/hashicorp/terraform-plugin-framework-validators` (the package that ships
`stringvalidator.OneOf` / `stringvalidator.LengthBetween`) is **not** a dependency of this module and
is not vendored, so the two validators are hand-written against the framework's native
`validator.String` interface (`Describer` + `ValidateString`). They live in the same action file and
are also re-applied inside `Invoke` as a defensive check, because `Invoke` may be reached with a config
that skipped the plan-phase validation.

Both are required by the API contract (`ClientToken` must be supplied by the caller and be identical
across retries of the same logical creation; `PayMode` has no server-side default for a new
subscription). `client_token` is **not** marked `Sensitive` — it is an idempotency token, not a secret,
and the API docs explicitly ask callers to generate `ssm-<UUID>` values.

Because this is an action and not a resource, `resource_attribute.md` conventions for `Computed` do not
apply: there is no `Computed` at all. The nine `GetServiceStatus` response fields listed in the
requirement are **read-only observations**, reported to the operator through progress messages and, on
failure, error details:

```
service_enabled=<bool> invalid_type=<n> access_key_escrow_enabled=<bool> expire_time=<str>
qps_limit=<n> secret_limit=<n> pay_model=<str> renew_flag=<n> resource_region=<n> resource_id=<str>
```

This is the only framework-legal way to expose them: an action cannot persist attributes between runs.
The design therefore states explicitly that **the nine observation fields are not Terraform attributes**
and must not be declared in the schema (declaring them would create dead `Required`/`Optional`
attributes that could never be populated).

### D4 — Invoke lifecycle

```
Invoke
 ├─ req.Config.Get(ctx, &data)                    // read pay_mode, client_token
 ├─ validate: pay_mode non-null/unknown/non-empty, OneOf enforced by validators
 ├─ validate: client_token non-null/unknown/non-empty, 1..64
 ├─ a.client must be configured (defensive nil check → AddError)
 ├─ CreateServiceWithContext(ctx, req{PayMode, ClientToken})   [retry: tccommon.WriteRetryTimeout]
 │    ├─ on ResourceInUse.ServiceAlreadyEnabled → do NOT fail; log + fall through to the polling loop
 │    │    (the desired end state — a subscribed account — is already reached; the API is idempotent by
 │    │     contract and this is the only re-invoke-safe behaviour)
 │    └─ on other error → AddError("Error subscribing to SSM", "<CreateService failed: %v>")
 ├─ report resource_id / order_id through SendProgress
 └─ poll GetServiceStatusWithContext            [retry loop until ServiceEnabled == true or timeout]
      ├─ ServiceEnabled == true → SendProgress summary of the nine observation fields, success
      └─ timeout → AddError("SSM service subscription did not become enabled…", last snapshot)
```

Design notes:

1. **Retry, not `d.SetId`.** The SDKv2-oriented rule "set the id before polling" has no analogue here:
   actions have no state, so there is nothing to lose if polling fails. The equivalent safety property is
   that the subscription itself is account-level and idempotent — a later re-invoke is a no-op that
   falls into the `ServiceAlreadyEnabled` branch.
2. **`client_token` reuse across retries.** The token is read once from config and reused unchanged for
   every retry inside `resource.Retry`, so a network-level retry cannot create a second billing order
   (the API explicitly requires "同一次逻辑创建发生重试时，须使用相同的值").
2b. **`ResourceInUse.ServiceAlreadyEnabled` must be intercepted before `tccommon.RetryError`.**
   `tccommon.retryableErrorCode` contains the short code `ResourceInUse`, and
   `tccommon.IsExpectError` falls back to the short code of a dotted error code — therefore
   `ResourceInUse.ServiceAlreadyEnabled` would otherwise be classified as **retryable** and spin for the
   whole `WriteRetryTimeout`. The action checks `isSsmServiceAlreadyEnabled(err)` first and returns
   `nil` from the retry closure, so the "already subscribed" case exits immediately and falls through to
   the polling branch.
3. **Polling budget.** `tccommon.WriteRetryTimeout` for the create call; an explicit bounded loop for the
   status polling (a `for` with `time.Sleep(5s)` capped by `tccommon.WriteRetryTimeout` measured from the
   start of `Invoke`), erroring out with the last observed snapshot. `ServiceEnabled` may legitimately
   stay `false` for a while right after the order is placed, so a single status read is not sufficient.
4. **`resource.Retry` usage.** The project's generic rule requires retry wrapping with
   `tccommon.RetryError()`. The action reuses `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource`
   exactly like the surrounding SSM service code does, so the retry semantics (retryable error codes,
   timeouts) stay identical to the rest of the provider.

### D5 — Vendored API contract (extracted, no need to re-read the SDK later)

`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923/`

Client methods (`client.go`):
- `func (c *Client) CreateServiceWithContext(ctx context.Context, request *CreateServiceRequest) (response *CreateServiceResponse, err error)` — line 339
- `func (c *Client) GetServiceStatusWithContext(ctx context.Context, request *GetServiceStatusRequest) (response *GetServiceStatusResponse, err error)` — line 1345
- Constructors: `ssm.NewCreateServiceRequest()`, `ssm.NewGetServiceStatusRequest()`

`CreateServiceRequest` (`models.go` 457–465):

```go
type CreateServiceRequest struct {
	*tchttp.BaseRequest

	// Only PostPay is supported; case-sensitive.
	PayMode *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`

	// Client-generated creation token. Retries of the SAME logical creation must reuse the same value.
	// Format: 1-64 visible ASCII characters (! to ~). Required. No spaces, Chinese or control chars.
	// Recommended shape: ssm-<UUID>.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}
```

`CreateServiceResponseParams` (`models.go` 488–497):

```go
type CreateServiceResponseParams struct {
	// SSM service resource id; for postpaid phase 1 it is "ssm_" + AppId, e.g. ssm_1250000000.
	// The resource may still be provisioning when returned — this does NOT mean the service is usable.
	ResourceId *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`

	// Billing order id, used to reconcile this subscription order.
	OrderId *string `json:"OrderId,omitnil,omitempty" name:"OrderId"`

	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateServiceResponse struct {
	*tchttp.BaseResponse
	Response *CreateServiceResponseParams `json:"Response"`
}
```

`GetServiceStatusRequest` (`models.go` 1383–1386): **no fields** —
`type GetServiceStatusRequest struct { *tchttp.BaseRequest }`.

`GetServiceStatusResponseParams` (`models.go` 1408–1444):

```go
type GetServiceStatusResponseParams struct {
	ServiceEnabled         *bool   // true: subscribed, false: not subscribed
	InvalidType            *int64  // 0 not purchased, 1 normal, 2 suspended for arrears, 3 released
	AccessKeyEscrowEnabled *bool   // key escrow capability available
	ExpireTime             *string // expiration time
	QPSLimit               *int64  // QPS limit
	SecretLimit            *int64  // secret count limit
	PayModel               *string // billing mode
	RenewFlag              *int64  // 0 manual renewal, 1 auto renewal, 2 no renewal at expiry
	ResourceId             *string // resource id
	TotalCount             *int64  // number of hosted secrets  (present in SDK, NOT in requirement list)
	ResourceRegion         *int64  // region id for prepaid SSM resources
	RequestId              *string
}

type GetServiceStatusResponse struct {
	*tchttp.BaseResponse
	Response *GetServiceStatusResponseParams `json:"Response"`
}
```

Documented error code of interest (`errors.go` line 81 and `client.go` 320–326 comment block):

```go
RESOURCEINUSE_SERVICEALREADYENABLED = "ResourceInUse.ServiceAlreadyEnabled"
```

`CreateService` may otherwise return `FAILEDOPERATION`, `INTERNALERROR`, `INVALIDPARAMETERVALUE`,
`UNAUTHORIZEDOPERATION`.

### D6 — Framework action structure (template `references/5-action.md`, adapted to baseline B)

The template skeleton, with `fw.ActionWithConfigure` replaced by the hand-written `Configure` from D2 and
`a.Client()` replaced by `a.client` (an `SsmService`), is:

```go
package ssm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ action.Action              = &SsmServiceAction{}
	_ action.ActionWithConfigure = &SsmServiceAction{}
)

func NewSsmService() action.Action {
	return &SsmServiceAction{}
}

type SsmServiceAction struct {
	client SsmService
}

type SsmServiceActionModel struct {
	PayMode     types.String `tfsdk:"pay_mode"`
	ClientToken types.String `tfsdk:"client_token"`
}

func (a *SsmServiceAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_service"
}

func (a *SsmServiceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Subscribes the current account to Secrets Manager (SSM) by calling CreateService. " +
			"This is a one-time action: it keeps no Terraform state, every invoke performs a real API call, " +
			"and there is no way to unsubscribe (destroy does not revert it).",
		Attributes: map[string]schema.Attribute{
			"pay_mode": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{ /* stringvalidator.OneOf("PostPay") */ },
				Description: "Billing mode of the SSM subscription. Only `PostPay` is supported (case-sensitive).",
			},
			"client_token": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{ /* stringvalidator.LengthBetween(1, 64) */ },
				Description: "Client-generated idempotency token for the subscription request, 1 to 64 visible " +
					"ASCII characters. Reuse the same value when retrying the same logical creation.",
			},
		},
	}
}

func (a *SsmServiceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) { ... }
```

Template key points carried over verbatim (from `references/5-action.md` "实现要点"):

1. Only `Invoke` (+ `Metadata` / `Schema` / `Configure`); no `Read`/`Delete`/`Update`.
2. Input validation (`IsNull()` / `IsUnknown()` / empty value) is written at the top of `Invoke`.
3. All errors go through `resp.Diagnostics.AddError`; never return silently.
4. The Description must state "repeated invoke = repeated operation" and "destroy does not roll back".
5. Do not register a pseudo-resource to emulate the one-shot operation.
6. Naming: file `action_tc_ssm_service.go`, TF name `tencentcloud_ssm_service`, factory
   `NewSsmService()` (no `Action` suffix), registered in `framework/registry.go` `actionFactories`.
7. Action requires Terraform CLI ≥ 1.14; the repository is already mux-wired, so `main.go` is untouched.

Field-type rules from `references/6-field-types-and-computed.md` that apply here: no `Computed` /
`Default` / plan modifiers are used at all (actions have no state and no plan); every read of a config
value performs the `!IsNull() && !IsUnknown()` double check before calling `ValueString()`.

### D7 — Unit test strategy

`tencentcloud/services/ssm/action_tc_ssm_service_test.go` uses **gomonkey** to mock the cloud API — no
acceptance test suite, no real credentials:

- `TestSsmServiceAction_Metadata` / `TestSsmServiceAction_Schema` — assert `TypeName`, that the two
  attributes exist and are required, and that no extra attributes are declared.
- `TestSsmServiceAction_Invoke_InvalidConfig` — `pay_mode` null / empty and `client_token` empty must
  produce diagnostics without touching the mocked API.
- `TestSsmServiceAction_Invoke_CreateServiceError` — mocked `CreateService` returning a non-retryable
  error must yield an error diagnostic.
- `TestSsmServiceAction_Invoke_AlreadyEnabled` — mocked `CreateService` returning
  `ResourceInUse.ServiceAlreadyEnabled` must not fail, and the polling path must succeed when
  `GetServiceStatus` returns `ServiceEnabled=true`.
- `TestSsmServiceAction_Invoke_Success` — full happy path: `CreateService` returns `ResourceId`/`OrderId`,
  `GetServiceStatus` returns `ServiceEnabled=true`, no diagnostics.
- `TestSsmServiceAction_Invoke_PollingTimeout` — `ServiceEnabled` stays `false`; assert the timeout error
  diagnostic. Should be run with a shortened polling budget (package-level test hook variable) so the
  case does not sleep for the whole write timeout.

Building `action.InvokeRequest` / `action.InvokeResponse` in a unit test requires a `tfsdk.Config`; the
test constructs one from the action's own `Schema` with `req.Config.Raw = tftypes.NewValue(...)` so the
assertions exercise real framework decoding rather than a hand-rolled stub.

## Risks / Trade-offs

| Risk | Mitigation |
| --- | --- |
| Subscription is a paid, irreversible account-level change | Documented in the action Description and the `.md` example; the action is never implicitly executed (requires explicit `-invoke`) |
| `CreateService` is asynchronous; a false `ServiceEnabled` right after ordering | Bounded polling loop until `ServiceEnabled == true`, with the last snapshot included in the timeout diagnostic |
| Re-invoke on an already subscribed account returns `ResourceInUse.ServiceAlreadyEnabled` | Treated as success (end state already reached) and routed into the polling branch, making the action re-runnable without manual intervention |
| First framework action in this repository — no precedent to copy | `Configure` aligned byte-for-byte with the in-repo framework resource precedent; only `action`-specific types (`action.InvokeRequest`, `resp.SendProgress`) are new |
| Blocking poll inside `Invoke` holds the provider RPC | Poll interval 5 s, hard cap tied to `tccommon.WriteRetryTimeout`; progress events keep the CLI informed |
| The nine `GetServiceStatus` fields cannot be persisted by an action | Explicitly documented as invoke-time observations reported via progress/diagnostics; no phantom schema attributes |

## Migration Plan

Purely additive: one new Go file, one new `.md`, one new test, one registry line, one `provider.md`
section, one generated website page. Rollback = revert those files. No state migration, no schema change
to existing references, so no user-side action is required.

## Open Questions

None. The one ambiguity in the requirement text — listing `GetServiceStatus` response fields as if they
were resource attributes — is resolved in D3: they are reported as invoke-time observations because a
Framework action cannot hold state.
