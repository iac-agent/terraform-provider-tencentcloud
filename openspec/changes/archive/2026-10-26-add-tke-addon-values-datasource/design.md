# Design: tencentcloud_tke_addon_values Data Source

## Context

The provider currently has no way to read TKE addon parameter metadata/value templates. The TKE
v20180525 SDK (already vendored in this repo) provides `DescribeAddonValues`, a pure read API that
takes a cluster ID plus an addon name and returns two JSON strings:

- `Values`: parameter list rendered with the parameters already set on the addon (effectively the
  current values); **empty when the addon is not installed**.
- `DefaultValues`: parameter list taken from the chart defaults.

Both are plain JSON strings (not base64) per the API documentation. There is no pagination and no
asynchronous behavior — a single synchronous call returns everything.

### Resource kind verification (mandatory self-check)

The requested kind is `RESOURCE_KIND_DATASOURCE`. Checking the hard criteria against the interface
set actually used by this change:

| Criterion | Observation |
| --- | --- |
| Only Describe/query interfaces | Yes — `DescribeAddonValues` is the only interface |
| Any Create / Modify / Delete / Attach / Detach / run-once trigger | None |

Conclusion: the interface composition is **consistent** with `RESOURCE_KIND_DATASOURCE`. No
correction is required, so the `references/2-data-source.md` template selected up front applies.

A secondary check on the requested schema path: the mapping names the resource
`tencentcloud_tke_addon_values`, which reads naturally, but the TKE section of this provider uses the
`tencentcloud_kubernetes_*` family for its addon references (`tencentcloud_kubernetes_addons`,
`tencentcloud_kubernetes_addon`, `tencentcloud_kubernetes_addon_config`). The Terraform type name is
kept exactly as requested — `tencentcloud_tke_addon_values` — because it was explicitly specified,
while the **Go file name** follows the mandatory `<type-prefix>_tc_<product>_<name>.go` convention
where `<product>` is `tke`:
`tencentcloud/services/tke/data_source_tc_tke_addon_values.go`.

### Constraints

- Terraform Plugin Framework (not SDKv2/helper-schema): the file must implement
  `datasource.DataSource` + `datasource.DataSourceWithConfigure`.
- Credentials/SDK client are owned by the SDKv2 provider; framework code only retrieves the shared
  client via `sharedmeta.ProviderMeta` in `Configure`.
- `vendor/` is authoritative — no new module may be added; the TKE v20180525 package is already
  vendored (see the excerpt below).
- No `website/` edits and no `.changelog/` files in this phase; docs are produced by `make doc`.

## Goals / Non-Goals

**Goals:**

- Expose `DescribeAddonValues` as a read-only data source with `cluster_id` + `addon_name` inputs and
  `values` + `default_values` computed outputs.
- Keep the implementation small and idiomatic: one `Read`, no retries/loops beyond what the shared
  client already performs.
- Follow the team empty-result convention for data sources (see Decisions D3).
- Provide everything the implementation phase needs inline in this document (SDK structs, template
  skeleton, file/registration checklist) so `vendor/models.go` never has to be re-read.

**Non-Goals:**

- No create/update/delete of addons; this data source never mutates cloud state.
- No decoded/parsed JSON output fields (e.g. a `values_json` object). The API returns opaque JSON
  strings whose schema differs per addon; returning them verbatim keeps the contract stable. Users
  can pipe them through `jsondecode()` in HCL if needed.
- No addon version / cluster region arguments: not present in the API.
- Not a `tencentcloud_kubernetes_addons` replacement; that data source lists installed addons, this
  one returns parameter templates for one addon.

## Decisions

### D1. File, type name, factory, registration

| Item | Value |
| --- | --- |
| Go file | `tencentcloud/services/tke/data_source_tc_tke_addon_values.go` (package `tke`) |
| TF type name | `tencentcloud_tke_addon_values` (must match the requested name) |
| Factory | `NewTkeAddonValuesDataSource() datasource.DataSource` |
| Interface assertions | `var _ datasource.DataSourceWithConfigure = &TkeAddonValuesDataSource{}` |
| Registry | one line in `tencentcloud/framework/registry.go` → `dataSourceFactories` |
| Doc index | one line under `Tencent Kubernetes Engine(TKE)` → `Data Source` in `tencentcloud/provider.md` |
| Sibling doc | `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` |
| Unit test | `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go` (gomonkey mock) |

`dataSourceFactories` is currently empty (`var dataSourceFactories = []func() datasource.DataSource{}`),
and `registry.go` only imports `services/ssm`. Registration therefore needs **two** edits: add the
`tke` import and add the factory line:

```go
import (
	...
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
)

// dataSourceFactories lists every framework DataSource factory.
var dataSourceFactories = []func() datasource.DataSource{
	tke.NewTkeAddonValuesDataSource,
}
```

Rationale: per `tencentcloud/framework/README.md`, framework references are registered **only** in
`registry.go`; `tencentcloud/provider.go` (SDKv2 map) must not be touched.

### D2. Client wiring — embed the framework configure helper and reuse `NewTkeService`

The framework-side `fw` helper package referenced by the template does not exist in this repository
snapshot (`tencentcloud/framework/fw/` is absent; there is no `DataSourceWithConfigure` symbol
anywhere under `tencentcloud/`). The working, already-shipped pattern is the zero-helper one used by
`tencentcloud/services/ssm/ephemeral_tc_ssm_secret_version.go`:

```go
type TkeAddonValuesDataSource struct {
	client TkeService
}

func (d *TkeAddonValuesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError("Unexpected ProviderData type", fmt.Sprintf(...))
		return
	}
	d.client = NewTkeService(meta.Client)
}
```

`NewTkeService(client *connectivity.TencentCloudClient) TkeService` already exists at
`tencentcloud/services/tke/service_tencentcloud_tke.go:150`, and `TkeService` holds
`client *connectivity.TencentCloudClient` with the accessor `me.client.UseTkeClient()` returning
`*tke.Client`. Reusing it keeps credential/UA/retry behavior identical to SDKv2 resources.

Because `TkeService` methods take the standard SDKv2 style (`ctx`, `logId`, deferred logging), the
read path goes through a **new service-layer method** rather than calling the SDK client directly
from the data source. This matches how every other TKE data source in the package is built.

Alternatives considered:

- Direct `d.client.client.UseTkeClient().DescribeAddonValuesWithContext(...)` inside `Read` — fewer
  lines, but bypasses the package convention of one service method per API and loses the
  `[CRITAL]` failure log.
- Introducing a `fw.DataSourceWithConfigure` helper package — rejected: it would be new shared
  infrastructure out of scope for this change, and the ephemeral implementation proves the
  `sharedmeta` path works today.

### D3. Empty-result semantics

The output shape is flat (two scalar strings), i.e. the template's "single object" case, but the API
contract makes one of the two fields legitimately empty: `Values` is empty when the addon is not
installed. The data source never has an "item not found" state — a nonexistent cluster or addon name
surfaces as an API error (`ResourceUnavailable` / `InvalidParameter`), not as an empty response.

Therefore:

- `err != nil` → `AddError` with the API name and inputs in the detail; return.
- `response.Response == nil` → `AddError` ("empty response body"); return. This is an interface
  failure, not a business state.
- `Response.Values == nil` / `Response.DefaultValues == nil` → **not an error**: write
  `types.StringNull()` for the missing one (a nil SDK pointer must not become `""`, otherwise
  `values` would misleadingly show an empty string for a non-installed addon and users could not
  distinguish it from the API returning `""`). No `SetId`/`RemoveResource` is involved because
  framework data sources have no persistent state.

No retry wrapper is added: `connectivity.TencentCloudClient` already applies the provider-wide retry
policy, and `DescribeAddonValues` is idempotent and read-only.

### D4. Schema shape (inputs Required, outputs Computed)

```go
resp.Schema = schema.Schema{
	Description: "Use this data source to query the parameter values of a TKE addon.",
	Attributes: map[string]schema.Attribute{
		"cluster_id": schema.StringAttribute{
			Required:    true,
			Description: "ID of the TKE cluster the addon belongs to.",
		},
		"addon_name": schema.StringAttribute{
			Required:    true,
			Description: "Name of the addon, as listed in the cluster's addon management page.",
		},
		"values": schema.StringAttribute{
			Computed:    true,
			Description: "Parameter values of the addon rendered with the values already set, as a JSON string. Empty or null when the addon is not installed.",
		},
		"default_values": schema.StringAttribute{
			Computed:    true,
			Description: "Default parameter values supported by the addon, as a JSON string.",
		},
	},
}
```

No `Sensitive` flags: these are chart parameter templates, not credentials. No `Default`, no
`Validators` (addon names are open-ended), no plan modifiers (data sources expose no plan).

### D5. Read implementation skeleton

```go
func (d *TkeAddonValuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TkeAddonValuesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Required attributes: guard against empty strings explicitly, since Framework only
	// enforces presence, not non-emptiness.
	clusterId := config.ClusterId.ValueString()
	addonName := config.AddonName.ValueString()
	if clusterId == "" {
		resp.Diagnostics.AddError("Missing cluster_id", "cluster_id is required and cannot be empty.")
		return
	}
	if addonName == "" {
		resp.Diagnostics.AddError("Missing addon_name", "addon_name is required and cannot be empty.")
		return
	}

	values, defaultValues, err := d.client.DescribeAddonValues(ctx, clusterId, addonName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TKE addon values",
			fmt.Sprintf("DescribeAddonValues failed for cluster %q addon %q: %v", clusterId, addonName, err),
		)
		return
	}

	// nil pointer -> Null (do NOT collapse to ""), see D3.
	config.Values = types.StringPointerValue(values)
	config.DefaultValues = types.StringPointerValue(defaultValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
```

Service-layer method (added to `service_tencentcloud_tke.go`, next to the other `Describe*`
helpers):

```go
func (me *TkeService) DescribeAddonValues(ctx context.Context, clusterId, addonName string) (values *string, defaultValues *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeAddonValuesRequest()
	request.ClusterId = &clusterId
	request.AddonName = &addonName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, errRet := me.client.UseTkeClient().DescribeAddonValuesWithContext(ctx, request)
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeAddonValues returned an empty response")
		return
	}

	values = response.Response.Values
	defaultValues = response.Response.DefaultValues
	return
}
```

Note: use `WithContext` so cancellation/credential context propagate; the deprecated `Client.DescribeAddonValues`
wrapper drops `ctx`.

### D6. Model struct

```go
type TkeAddonValuesDataSourceModel struct {
	ClusterId     types.String `tfsdk:"cluster_id"`
	AddonName     types.String `tfsdk:"addon_name"`
	Values        types.String `tfsdk:"values"`
	DefaultValues types.String `tfsdk:"default_values"`
}
```

### D7. Documentation file

`tencentcloud/services/tke/data_source_tc_tke_addon_values.md` — one-line description starting with
`Use this data source to query ...` plus `Example Usage`, and **no** `Argument Reference` /
`Attribute Reference` sections (auto-generated by the doc tool). Show `jsondecode()` usage because the
outputs are JSON strings:

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "cbs"
}

output "addon_default_values" {
  value = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}
```

### D8. Tests

Per project rules a **new** reference is unit-tested with gomonkey mocks only (no Terraform
acceptance suite). The test uses the exported model + factory and mocks
`(*tke.Client).DescribeAddonValuesWithContext`, so it needs a constructed `TkeService`. Two viable
seams, in order of preference:

1. Mock the SDK client method via `gomonkey.ApplyMethodReturn` on `*tke.Client` and build the
   `TkeService` with `NewTkeService(&connectivity.TencentCloudClient{})` — requires the client
   cache to hand back the mocked instance, so verify at implementation time; if the cache cannot be
   injected in a unit test, fall back to option 2.
2. Cover the pure mapping helper: extract the `*DescribeAddonValuesResponseParams -> model` mapping
   into a small unexported function (e.g. `fillTkeAddonValues(config *TkeAddonValuesDataSourceModel, in *tke.DescribeAddonValuesResponseParams)`)
   and table-test it directly (nil pointers → Null, empty strings → "", both set → values). This is
   deterministic and needs no network.

Both are listed; the implementation should ship option 2 at minimum and option 1 if the client seam
works. Tests must not call `go test` during this workflow.

## Reference: vendored SDK definitions (TKE v20180525)

`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525/models.go`

Request (`models.go:7302-7310`):

```go
type DescribeAddonValuesRequest struct {
	*tchttp.BaseRequest

	// 集群ID，请从容器服务控制台集群列表中获取（https://console.cloud.tencent.com/tke2/cluster）。
	ClusterId *string `json:"ClusterId,omitnil,omitempty" name:"ClusterId"`

	// addon名称，请从容器服务控制台进入集群详情页后，在组件管理页面中获取。
	AddonName *string `json:"AddonName,omitnil,omitempty" name:"AddonName"`
}
```

Response params (`models.go:7333-7342`):

```go
type DescribeAddonValuesResponseParams struct {
	// 参数列表，如果addon已安装，会使用已设置的参数和chart里的默认参数做渲染，是一个json格式的字符串，未安装addon时返回为空值。
	Values *string `json:"Values,omitnil,omitempty" name:"Values"`

	// addon支持的参数列表，值为chart的默认值，是一个json格式的字符串。
	DefaultValues *string `json:"DefaultValues,omitnil,omitempty" name:"DefaultValues"`

	// 唯一请求 ID ...
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeAddonValuesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAddonValuesResponseParams `json:"Response"`
}
```

Client (`client.go:5763-5792`) — constructor + `WithContext` signature:

```go
func NewDescribeAddonValuesRequest() (request *DescribeAddonValuesRequest)
func NewDescribeAddonValuesResponse() (response *DescribeAddonValuesResponse)
func (c *Client) DescribeAddonValuesWithContext(ctx context.Context, request *DescribeAddonValuesRequest) (response *DescribeAddonValuesResponse, err error)
```

Documented error code: `RESOURCEUNAVAILABLE = "ResourceUnavailable"` (plus generic
`InvalidParameter*` at the platform level). No pagination fields; response has no list to iterate.

### Field mapping (authoritative)

| Cloud API (JsonPath) | SDK field | TF attribute | Kind |
| --- | --- | --- | --- |
| `request.ClusterId` | `ClusterId *string` | `cluster_id` (string) | Required input |
| `request.AddonName` | `AddonName *string` | `addon_name` (string) | Required input |
| `response.Response.Values` | `Values *string` | `values` (string) | Computed output |
| `response.Response.DefaultValues` | `DefaultValues *string` | `default_values` (string) | Computed output |

`response.Response.RequestId` is intentionally not exposed (framework/provider logs already carry it).

## Reference: Framework data-source template essentials (`references/2-data-source.md`)

Applied structure (with the D2 deviation to `sharedmeta` wiring instead of the nonexistent `fw` package):

- Implements `datasource.DataSource` → `Metadata` (sets `resp.TypeName`), `Schema`, `Read`; plus
  `Configure` from `datasource.DataSourceWithConfigure` for client injection.
- Inputs `Required`/`Optional` only — never `Computed`. Outputs `Computed` only.
- `Read` starts with `req.Config.Get(ctx, &config)` and appends diagnostics; early-returns on
  `resp.Diagnostics.HasError()`.
- Never mutates cloud state.
- No plan/refresh distinction — `Read` runs on every plan, so no drift handling needed.
- Nested lists (not used here) would need `schema.ListNestedAttribute` + `types.ListValueFrom`;
  blocks are deprecated.
- Pointer/value conversion via `ValueStringPointer()` / `types.StringPointerValue()`
  (`references/6-field-types-and-computed.md` §1).
- Framework model fields use `tfsdk:"<snake_case>"` tags matching schema keys exactly.
- Naming/registration checklist (SKILL.md 命名对照表): file `data_source_tc_<product>_<name>.go`,
  factory `New<Product><NameCamel>DataSource()`, slice `dataSourceFactories`, plus
  `provider.md` + sibling `<stem>.md`.

## Risks / Trade-offs

- **[Opaque JSON outputs]** `values` / `default_values` content differs per addon and is a raw
  string; provider cannot validate or type it → Mitigation: document the JSON nature and demonstrate
  `jsondecode()`; keep verbatim pass-through so a cloud-side schema change never breaks the provider.
- **[`values` empty when addon not installed]** users may expect current values always →
  Mitigation: `Description` states it explicitly, and nil → `types.StringNull()` so the state
  distinguishes "not installed" (null) from "installed with empty parameters" (`""`).
- **[Type name `tencentcloud_tke_*` vs package convention `tencentcloud_kubernetes_*`]** mild
  inconsistency with sibling addon references → Accepted: the name is explicitly requested; file
  name still follows the product-based convention. If maintainers later prefer
  `tencentcloud_kubernetes_addon_values`, it is a rename-only change (no state migration needed for
  a data source, since data sources are not persisted).
- **[No retry wrapper]** transient `ResourceUnavailable` errors will fail the plan →
  Mitigation: the shared client already retries per provider policy; failing fast on
  `ResourceUnavailable` is correct because the cluster/addon genuinely may not exist.
- **[Unit test seam for client mocking]** `connectivity.TencentCloudClient` client cache may not be
  injectable in a pure unit test → Mitigation: D8 option 2 (pure mapping helper) is the guaranteed
  coverage path; option 1 is additive.
- **[Empty-response branch is defensive]** `response.Response == nil` is unlikely from the SDK but
  would nil-panic without the guard → Mitigation: explicit guard + `AddError`.

## Migration Plan

1. Add the service method `TkeService.DescribeAddonValues` to `service_tencentcloud_tke.go`.
2. Add `data_source_tc_tke_addon_values.go` (schema + model + Metadata/Configure/Read).
3. Register in `tencentcloud/framework/registry.go` (`tke` import + `dataSourceFactories` entry).
4. Add `tencentcloud/provider.md` entry under the TKE `Data Source` node.
5. Add sibling `.md` and the gomonkey unit test.
6. Finalization phase (out of scope here): `gofmt`, `make doc`, changelog entry.

Rollback: delete the new files and revert the two registration edits — no state or schema migration
is involved because the artifact is a brand-new, read-only data source.

## Open Questions

- None blocking. Open (deferred to implementation review): whether maintainers want the type renamed
  to `tencentcloud_kubernetes_addon_values` for family consistency (see Risks); this design proceeds
  with the requested `tencentcloud_tke_addon_values`.
