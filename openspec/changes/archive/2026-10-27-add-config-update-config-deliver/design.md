# Design: `tencentcloud_config_update_config_deliver`

## Context

Tencent Cloud Config exposes a single per-account/per-region delivery setting controlled by two APIs:

| API | Request | Response |
| --- | --- | --- |
| `DescribeConfigDeliver` | no parameters | the whole setting, flattened |
| `UpdateConfigDeliver` | 6 optional fields, `Status` is the switch | only `RequestId` |

There is no `Create*` and no `Delete*` API. The setting always exists; it just starts with `Status = 0`
(disabled). That interface composition (`Modify` + `Describe` only) matches **RESOURCE_KIND_CONFIG**
unambiguously, so the `ResourceKind` given in the task and the API composition agree — no correction needed.

Baseline facts probed in this repository (they override anything the generic templates assume):

```text
ls tencentcloud/framework/fw        -> NO_FW_PACKAGE  (the fw helper package does NOT exist here)
ls tencentcloud/services/*/action_tc_*.go -> (empty)
framework/registry.go: resourceFactories = [ssm.NewSsmSecretVersionV2Resource]
framework/registry.go: dataSourceFactories = [] , actionFactories = []
```

Consequences: the config template's `fw.ResourceWithConfigure` / `fw.WithImportByID` embeds **cannot** be
used. `Configure` and `ImportState` must be hand-written, modelled on the only existing Framework resource
in the repo, `tencentcloud/services/ssm/resource_tc_ssm_secret_version_v2.go`.

Also probed: the vendor for `config/v20220802` already exists (it was added by the sibling change), so no
`go.mod` / `vendor/` work is needed. `tencentcloud_config_update_config_deliver` is not registered anywhere,
so there is no name collision.

> Note on a pre-existing sibling: `tencentcloud_config_deliver_config` (SDKv2, same two APIs) already
> exists and is registered in `tencentcloud/provider.go`. It is intentionally **left untouched** — this
> change adds a separate Framework resource with a different type name and a different state-ID / import
> model. Deprecating or merging the two is out of scope.

## Goals / Non-Goals

**Goals:**

- A Framework resource that reads the delivery setting and writes it back, with a schema that maps
  one-to-one onto the 6 `UpdateConfigDeliver` request fields plus the 7 `DescribeConfigDeliver` response
  fields.
- Correct Terraform lifecycle for a singleton config: adoption of pre-existing cloud state on first apply,
  no spurious diff, `terraform import` support, and a destruction that does not leak a still-enabled
  delivery pipeline.
- Make the downstream implementation phase self-sufficient: all SDK struct excerpts and Framework
  skeleton details are recorded here so it need not re-read `vendor/`.

**Non-Goals:**

- No change to the SDKv2 `tencentcloud_config_deliver_config` resource.
- No new data source / action; no `Timeouts` block (both APIs are synchronous single-call operations, see
  Decisions #7).
- No cross-region / multi-account iteration: one resource instance manages the setting of the account and
  region the provider is configured for.
- No provider-side validation of the `TargetArn` six-segment format or of `deliver_type` ↔ ARN-scheme
  consistency: the cloud API is the authority and its error is surfaced verbatim.

## Decisions

### 1. Resource kind, file name, TF name, factory name

| Item | Value |
| --- | --- |
| Kind | `RESOURCE_KIND_CONFIG` |
| TF type name | `tencentcloud_config_update_config_deliver` |
| Go file | `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.go` |
| Doc file | `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.md` |
| Test file | `tencentcloud/services/config/resource_tc_config_update_config_deliver_config_test.go` |
| Factory | `config.NewConfigUpdateConfigDeliverResource()` |
| Go types | `ConfigUpdateConfigDeliverResource`, `ConfigUpdateConfigDeliverResourceModel` |
| Registration | `tencentcloud/framework/registry.go` → `resourceFactories` |

The TF-name prefix is **not** derived from the directory: the Config product family in
`tencentcloud/provider.md` (lines 2735–2748) and `provider.go` uses the `tencentcloud_config_` prefix for
all 10 existing entries, so `config` is confirmed as the prefix by rule 1 of the naming table (same-family
sibling). The required name already carries the `update_config_deliver` part, and the `_config` suffix is
appended to the **file name** only (config convention), not to the TF type name — the type name is used
verbatim as requested by the task.

### 2. State ID: fixed sentinel, plus import

The setting has no cloud-side identifier. Three options were considered:

| Option | Verdict |
| --- | --- |
| `helper.BuildToken()` (random) | Rejected. A random ID changes on every re-create, and after `terraform state rm` + re-apply the ID changes for the same cloud object — noisy and non-idempotent. It is what the sibling SDKv2 resource does; not worth copying. |
| Composite of the configured fields | Rejected. `deliver_name` / `target_arn` are optional, so the ID would be empty on a `status`-only resource, and it would churn on every rename — which for a resource that is a singleton makes `terraform import` unusable. |
| **Fixed sentinel constant** | **Chosen.** `const configUpdateConfigDeliverId = "config_update_config_deliver"` |

The sentinel makes `terraform import tencentcloud_config_update_config_deliver.x config_update_config_deliver`
work naturally (import ID == the sentinel), keeps the ID stable across re-creates, and honestly documents
that this resource is a singleton rather than pretending it has a per-object key.

### 3. `create_time` is Computed-only; everything else in the update request is Optional+Computed

- `create_time` is `Computed` (cloud-generated, never writable).
- `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_content_type` are
  `Optional + Computed` (no `Default`): the cloud side holds a real value for all of them once configured,
  and letting the cloud backfill keeps the provider from hard-coding cloud defaults.
- `status` is **`Required`**: it is the switch that decides whether logs are delivered, so it must be an
  explicit, intentional declaration in configuration. Making it `Optional + Computed` would let a user
  write a resource block containing only `deliver_name` and silently inherit whatever switch state the
  account happened to be in.

All 7 attributes are `tfsdk` snake_case names; the mapping is one-to-one with the API JSON names:

| Terraform attribute | `UpdateConfigDeliver` request | `DescribeConfigDeliver` response | Schema |
| --- | --- | --- | --- |
| `status` | `Status *uint64` (0/1) | `Status *uint64` | `Int64`, Required |
| `deliver_name` | `DeliverName *string` | `DeliverName *string` | `String`, Optional+Computed |
| `target_arn` | `TargetArn *string` | `TargetArn *string` | `String`, Optional+Computed |
| `deliver_prefix` | `DeliverPrefix *string` | `DeliverPrefix *string` | `String`, Optional+Computed |
| `deliver_type` | `DeliverType *string` | `DeliverType *string` | `String`, Optional+Computed, `OneOf("COS","CLS")` |
| `deliver_content_type` | `DeliverContentType *uint64` | `DeliverContentType *uint64` | `Int64`, Optional+Computed, `Between(1,3)` |
| `create_time` | — (not in request) | `CreateTime *string` | `String`, Computed |
| `id` | — | — | `String`, Computed, `UseStateForUnknown` |

`deliver_type` / `deliver_content_type` get validators from the vendor comments
(`投递类型` COS/CLS; `1：配置变更 2：资源列表 3：全部`). **They must be Optional** (not Required) because
the request fields are optional — a `status`-only destroy path calls the API without them.

`deliver_type` is deliberately **not** `RequiresReplace`: `UpdateConfigDeliver` accepts a new `DeliverType`
together with the new `TargetArn`, so switching COS→CLS is one in-place update, not a replace. (Validate
during manual testing that the API accepts a combined type+arn change; if it does not, the failure surfaces
as a plain API error and the user can express the two-step change — no schema change needed.)

### 4. Delete = turn off the switch (资源语义类型：开关启用型)

This resource is an **enable switch**: it turns log delivery on. Destroying it while `Status` stays `1`
would silently keep shipping audit logs (and cost) for a resource Terraform no longer tracks — a textbook
silent leak. So `Delete` calls `UpdateConfigDeliver` with `Status = 0`.

Binding rule from the template (implementation note #3): *if the Modify API can express "off", `Delete`
MUST turn it off*. `UpdateConfigDeliver.Status` can (`0 关闭`), so `Delete` must call it. `Delete` must
**never** attempt to delete the host configuration/account state — it only flips the switch.

The wording goes into both the `Schema.Description` and the `.md` file:
*destroying this resource disables Config log delivery on the cloud side; it never deletes the Config
service or any COS bucket / CLS topic.*

Delete uses `UpdateConfigDeliverRequest{Status: 0}` only — it does not need to clear the other fields, and
clearing `target_arn` first would make the disable call fail validation on some backends.

### 5. Read: nil-safe, no `RemoveResourceState` on a singleton

`DescribeConfigDeliver` is the only source of truth and cannot return "not found" in a meaningful sense.
Therefore Read:

- errors out (diagnostic) if `response == nil` or `response.Response == nil` — an empty body here is an API
  anomaly, not a deleted resource; so **do not** call `resp.State.RemoveResource`;
- writes each field via `types.StringPointerValue(...)` / `types.Int64ValueFromUint(...)`, i.e. a `nil`
  pointer becomes `Null` and is not written as a wrong zero value;
- backfills the sentinel ID when the state ID is empty (the post-import path);
- logs with `log.Printf("[CRITICAL] describe config_update_config_deliver failed, id=%s ...")` where
  `id` is the resource ID, so a failure can be located (per the logging rule, the resource is referred to
  as `config_update_config_deliver`, never "this resource").

Because the setting is a singleton, Read never needs a filter: it takes the whole response.

### 6. Create = first set, Update = re-set; both share one private `modifyConfig`

`Create` and `Update` both call the same `UpdateConfigDeliver`; the shared private helper is the signature
of the config pattern:

```text
Create: req.Plan.Get -> modifyConfig(plan) -> State.Set(plan)   // id = sentinel
Update: req.Plan.Get -> modifyConfig(plan) -> State.Set(plan, preserving id from state)
```

`modifyConfig` marshals only fields that are **set** — `if !v.IsNull() && !v.IsUnknown()` — mirroring the
"only send what the user/cloud has" rule. `status` is always sent because it is `Required`.

`status` is additionally guarded: if it is null or unknown, `modifyConfig` returns an error without calling
the API. `types.Int64.ValueInt64()` returns `0` for both cases, so an unguarded path would issue
`Status = 0` and silently disable audit-log delivery — the one failure mode of this resource that is both
invisible and costly. Because `status` is `Required` the guard is defensive, but it is the cheapest place to
stop that class of bug.

### 6.1 `*uint64` conversion stays local to the resource

`DescribeConfigDeliver` returns `Status` / `DeliverContentType` as `*uint64`, and
`UpdateConfigDeliver` expects `*uint64`. Two helper packages exist, and neither is usable here:

| Helper | Usable from `tencentcloud/services/config`? | Why |
| --- | --- | --- |
| `tencentcloud/framework/internal/helper` (`Int64ValueFromUint`) | **No** | It lives under `tencentcloud/framework/internal/`, so Go's `internal` visibility rule limits imports to the `tencentcloud/framework/...` subtree. Importing it from `tencentcloud/services/config` is a compile error. |
| `tencentcloud/internal/helper` (`Uint64`) | Yes | Imported as `tccommonhelper` for the **request** side. It has no `*uint64` → `types.Int64` flattener. |

So the request side uses `tccommonhelper.Uint64(...)` and the response side keeps a local
`uint64ToInt64Value(*uint64) types.Int64` that maps `nil` → `types.Int64Null()`.

`Create` does not need a pre-check against existing cloud state: `UpdateConfigDeliver` is idempotent and
overwrites, so first-apply adopts the existing setting in place. That is the documented behaviour for
config-type resources and avoids a read-modify-write race.

Everything happens inside one `retryFramework`-style loop over `tccommon.WriteRetryTimeout`; on failure the
error is wrapped with the action name and the resource ID.

### 7. No `Timeouts` block

Both APIs are synchronous, single-round-trip operations that return after the setting is stored; there is
no `FlowId` / `TaskId` to poll. The framework-side retry helper
`tencentcloud/framework/internal/helper.RetryFramework(ctx, timeout, retryableFn, fn)` is used with
`tccommon.WriteRetryTimeout` (write) / `tccommon.ReadRetryTimeout` (read) to absorb transient errors. Adding
a `Timeouts` block would advertise tunability the resource cannot honour.

### 8. Client injection: hand-written `Configure` (no `fw` package in this baseline)

```text
struct { client ConfigService }        // no fw.ResourceWithConfigure embed

func (r *X) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
    if !ok {
        resp.Diagnostics.AddError("Unexpected ProviderData type",
            fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. This is a bug in the provider; please report it.", req.ProviderData))
        return
    }
    r.client = NewConfigService(meta.Client)
}
```

Compile-time assertions to keep:
`var _ resource.Resource = &ConfigUpdateConfigDeliverResource{}`,
`var _ resource.ResourceWithConfigure = ...`,
`var _ resource.ResourceWithImportState = ...`.

`ImportState` is hand-written and simply accepts the ID and lets Read backfill the fields:

```text
resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
```

(validated against the sentinel: if `req.ID != sentinel`, return a clear diagnostic telling the user the
only valid import ID).

### 9. Service layer: methods on the existing `ConfigService`

`tencentcloud/services/config/service_tencentcloud_config.go` already owns `NewConfigService(client)` and
`ConfigService{client *connectivity.TencentCloudClient}`. Append (do not restructure):

```text
func (me *ConfigService) DescribeConfigDeliver(ctx context.Context) (*configv20220802.DescribeConfigDeliverResponseParams, error)
func (me *ConfigService) UpdateConfigDeliver(ctx context.Context, request *configv20220802.UpdateConfigDeliverRequest) error
```

Both use `me.client.UseConfigV20220802Client()` and `tccommon.GetLogId(ctx)` for `log.Printf` tracing —
the same client accessor the existing SDKv2 code in that file uses. They return the raw SDK

types so the Framework resource keeps full control over `nil`-vs-zero handling.

## Vendor reference (excerpts, verbatim from `vendor/.../config/v20220802/models.go`)

`sdk` below = `configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"`.

### `DescribeConfigDeliver` — no request parameters (models.go L2283–L2311)

```go
// Predefined struct for user
type DescribeConfigDeliverRequestParams struct {
}

type DescribeConfigDeliverRequest struct {
	*tchttp.BaseRequest
}
```

### `DescribeConfigDeliverResponseParams` (models.go L2312–L2342)

```go
type DescribeConfigDeliverResponseParams struct {
	// 投递名称
	// 注意：此字段可能返回 null，表示取不到有效值。
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`

	// 资源六段式
	// 注意：此字段可能返回 null，表示取不到有效值。
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`

	// 投递状态 DeliverStatus：0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`

	// 创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// 日志前缀
	// 注意：此字段可能返回 null，表示取不到有效值。
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`

	// 投递类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`

	// 1：配置变更   2： 资源列表 3：全部
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`

	// 唯一请求 ID ...
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *DescribeConfigDeliverResponseParams `json:"Response"`
}
```

Note `Status` and `DeliverContentType` are `*uint64` — use
`helper.Int64ValueFromUint(p)` (`tencentcloud/framework/internal/helper/types.go`) to convert, and
`helper.Int64(p)` on the request side; do **not** cast manually (loses the nil distinction).

### `UpdateConfigDeliverRequest` / `RequestParams` (models.go L5766–L5813)

```go
type UpdateConfigDeliverRequestParams struct {
	// 0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`

	// 投递服务名称
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`

	// 资源六段式
	// COS：qcs::cos:$region:$account:prefix/$appid/$BucketName
	// CLS:
	// qcs::cls:$region:$account:cls/topicId
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`

	// clonfig_fix
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`

	// 投递类型
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`

	// 1：配置变更 2： 资源列表 3：全选
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`
}

type UpdateConfigDeliverRequest struct {
	*tchttp.BaseRequest
	// (same six fields, same tags)
}
```

### `UpdateConfigDeliverResponseParams` (models.go L5839–L5848)

```go
type UpdateConfigDeliverResponseParams struct {
	// 唯一请求 ID ...
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type UpdateConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *UpdateConfigDeliverResponseParams `json:"Response"`
}
```

### Client constructors (client.go L1324–L1373, L3464–L3511)

```go
func NewDescribeConfigDeliverRequest() *DescribeConfigDeliverRequest
func (c *Client) DescribeConfigDeliverWithContext(ctx context.Context, request *DescribeConfigDeliverRequest) (*DescribeConfigDeliverResponse, error)

func NewUpdateConfigDeliverRequest() *UpdateConfigDeliverRequest
func (c *Client) UpdateConfigDeliverWithContext(ctx context.Context, request *UpdateConfigDeliverRequest) (*UpdateConfigDeliverResponse, error)
```

Both non-`WithContext` variants also exist but the `WithContext` form is used everywhere in this repo.

## Framework template reference (RESOURCE_KIND_CONFIG, adapted to this baseline)

The template `references/3-resource-config.md` is followed for **structure**; its `fw.*` embeds are
replaced per Decision #8 (see `references/0-baseline-adaptation.md`, column B).

### Skeleton

```go
package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const configUpdateConfigDeliverId = "config_update_config_deliver"

var (
	_ resource.Resource                = &ConfigUpdateConfigDeliverResource{}
	_ resource.ResourceWithConfigure   = &ConfigUpdateConfigDeliverResource{}
	_ resource.ResourceWithImportState = &ConfigUpdateConfigDeliverResource{}
)

func NewConfigUpdateConfigDeliverResource() resource.Resource { ... }

type ConfigUpdateConfigDeliverResource struct { client ConfigService }

type ConfigUpdateConfigDeliverResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Status             types.Int64  `tfsdk:"status"`
	DeliverName        types.String `tfsdk:"deliver_name"`
	TargetArn          types.String `tfsdk:"target_arn"`
	DeliverPrefix      types.String `tfsdk:"deliver_prefix"`
	DeliverType        types.String `tfsdk:"deliver_type"`
	DeliverContentType types.Int64  `tfsdk:"deliver_content_type"`
	CreateTime         types.String `tfsdk:"create_time"`
}

func (r *X) Metadata(...)  { resp.TypeName = "tencentcloud_config_update_config_deliver" }
func (r *X) Schema(...)    { /* table in Decision #3 */ }
func (r *X) Configure(...) { /* Decision #8 */ }
func (r *X) Create(...)    { /* plan -> modifyConfig -> id = sentinel -> State.Set */ }
func (r *X) Read(...)      { /* Decision #5 */ }
func (r *X) Update(...)    { /* plan -> modifyConfig -> State.Set (id from state) */ }
func (r *X) Delete(...)    { /* Decision #4: UpdateConfigDeliver{Status: 0} */ }
func (r *X) ImportState(...)
func (r *X) modifyConfig(ctx context.Context, m *Model) error { /* shared by Create/Update */ }
```

### Template checklist carried over (all must hold in the implementation)

1. Create and Update share one private `modifyConfig` — the signature trait of the config pattern. ✅
2. State ID = the host object's ID (here: the singleton sentinel, per Decision #2). ✅
3. Delete semantic chosen from the template table and written into **both** `Schema.Description` **and**
   the `.md`: 开关启用型 → turn off. ✅
4. Host ID uses plan modifiers — N/A here (there is no permissive host attribute); instead `id` gets
   `UseStateForUnknown`. Config-style fields are all mutable in place, so **no** `RequiresReplace` in this
   resource; a change to any attribute is an in-place `UpdateConfigDeliver`.
5. Optional+Computed for cloud-backfilled fields; `IsNull()` / `IsUnknown()` checked before reading. ✅
6. Naming: file `resource_tc_config_update_config_deliver_config.go`, factory
   `NewConfigUpdateConfigDeliverResource`, registered in `resourceFactories`. ✅

## Risks / Trade-offs

- **[Sentinel ID is global per provider instance]** Two `tencentcloud_config_update_config_deliver` blocks in
  one configuration would produce two identical state IDs, and Terraform would fight over the same
  singleton. → Mitigated by documentation ("exactly one instance per provider/region") and by the fact that
  the cloud object is itself a singleton, so no correctness loss; users who need multiple regions must use
  provider aliases.
- **[Destroy disables delivery for the whole account]** A user who had manually enabled delivery before
  adopting Terraform, then runs `terraform destroy` on this resource, turns delivery off. → This is the
  intended and documented behaviour (avoiding a silent leak is worse than the alternative); the `.md` states
  it explicitly.
- **[`Status` is `Required` while the API field is optional]** The provider is stricter than the API. →
  Deliberate: a switch must be declared. The only cost is that a user cannot write a resource block that
  "only fills in the other fields"; that is the desired semantic.
- **[Field-level validation may reject values the API would accept]** `deliver_type` `OneOf("COS","CLS")`
  and `deliver_content_type` `Between(1,3)` come from the vendor comments; if the cloud later adds a third
  delivery type, the validator becomes stale. → Acceptable: the enum is documented in the SDK and a stale
  validator produces a clear plan-time error rather than a confusing API error; relaxing it later is
  backward compatible.
- **[No `VerifyConfig` / pre-flight]** A first apply overwrites whatever the account had configured. →
  Same as any config-type resource; documented in the `.md`.
- **[Framework resource has no `Timeouts`]** Users cannot raise timeouts during cloud-side slowness. →
  Accepted (Decision #7); retries already absorb transient errors within `tccommon` timeouts.

## Migration Plan

Purely additive; nothing to migrate.

- Ship the resource. Existing `tencentcloud_config_deliver_config` (SDKv2) users are unaffected — different
  type name, different state file entry.
- Rollback = remove the new resource blocks and `terraform state rm` the instances; the cloud setting
  keeps whatever value it last had (except that a `destroy` already turned delivery off, so re-enabling is
  a one-attribute change).
- Adoption for users who already have a delivery setting configured out-of-band:
  `terraform import tencentcloud_config_update_config_deliver.x config_update_config_deliver`, then write
  matching HCL. Import accepts only the sentinel ID.

## Open Questions

None blocking. Two items to verify during the implementation/acceptance phase rather than here:

- Whether `UpdateConfigDeliver` accepts `DeliverType` + `TargetArn` type change in one call (Decision #3).
- Whether the backend rejects `Status=1` when `TargetArn` is absent — if it does, the API error is surfaced
  verbatim and the user must supply `target_arn`; no schema change required.
