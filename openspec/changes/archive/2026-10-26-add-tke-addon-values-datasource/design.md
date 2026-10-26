## Context

TKE (TencentCloud Kubernetes Engine) addons are managed by the existing resource `tencentcloud_kubernetes_addon`, which takes a `raw_values` argument: a base64-encoded JSON string of addon parameters. To fill that argument correctly, a user must know two things which are only obtainable from the cloud:

1. Which parameters the addon chart supports and what their default values are (`DefaultValues`).
2. What the parameter JSON renders to for an already-installed addon (`Values`, which mixes user-set parameters with chart defaults).

The TKE API `DescribeAddonValues` (v20180525) returns exactly these two JSON strings. The data source below exposes them so they can be consumed directly (`base64encode(data.tencentcloud_tke_addon_values.this.default_values)` feeds `raw_values` without a manual lookup).

### Resource kind verification

| Check | Result |
| --- | --- |
| APIs used | `DescribeAddonValues` only (one Describe/query API) |
| Create / Update / Delete APIs | None |
| Attach / Detach APIs | None |
| One-shot trigger APIs (`RunXxx` / `ConfirmXxx`) | None |
| Hard rule "Describe only → RESOURCE_KIND_DATASOURCE" | Matches |

Conclusion: the declared `RESOURCE_KIND_DATASOURCE` is consistent with the API composition. No correction is required.

### Current provider state

- `tencentcloud/framework/registry.go` already exists and has an empty `dataSourceFactories` slice — **this will be the first Framework data source in the repository**.
- The closest existing Framework analogue is the query-style ephemeral resource `tencentcloud/services/ssm/ephemeral_tc_ssm_secret_version.go`: it embeds a plain `client` field, receives `*sharedmeta.ProviderMeta` in `Configure`, and calls a single Describe API. The data source follows that shape, but implements `datasource.DataSource` (not `EphemeralResource`).
- `sharedmeta.ProviderMeta{Client *connectivity.TencentCloudClient}` is the `req.ProviderData` type that Framework references receive in `Configure`.
- TKE SDK package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525` (already vendored; no `go.mod` change).

## Goals / Non-Goals

**Goals:**

- Expose `DescribeAddonValues` as `data.tencentcloud_tke_addon_values` with required inputs `cluster_id`, `addon_name` and computed outputs `values`, `default_values`.
- Implement with Terraform Plugin Framework, registered in `tencentcloud/framework/registry.go`.
- Keep the implementation minimal: one Schema + one Read; no state mutation, no retry wrappers beyond what a read-only data source needs.
- Provide a unit test that mocks the TKE client (gomonkey) — no acceptance-test suite.

**Non-Goals:**

- Do not change `tencentcloud_kubernetes_addon` (no schema change, no new attribute).
- Do not JSON-decode/normalize `values` or `default_values` — they are returned verbatim as JSON strings, matching the cloud response and the `raw_values` round-trip (the existing `jsonencode()` + `base64encode()` pattern in `.md` examples covers user needs).
- Do not add pagination / filtering inputs; the API takes none.
- No asynchronous wait: `DescribeAddonValues` is synchronous.

## Decisions

### D1. Resource kind and template

Use **`RESOURCE_KIND_DATASOURCE` → `references/2-data-source.md`** (Framework data source template) plus the shared addendum `references/6-field-types-and-computed.md`.

Alternative rejected: SDKv2 `data_source_tc_*.go` following the neighbouring TKE files. Rejected because the task mandates Framework for new references and `registry.go` is the registration point for Framework references.

### D2. Naming

Per the tfpacer naming table (product `tke` → `Tke`):

| Item | Value |
| --- | --- |
| Terraform name | `tencentcloud_tke_addon_values` |
| Go file | `tencentcloud/services/tke/data_source_tc_tke_addon_values.go` |
| Factory | `NewTkeAddonValuesDataSource()` |
| Data source struct | `TkeAddonValuesDataSource` |
| Model struct | `TkeAddonValuesDataSourceModel` |
| Registration | `dataSourceFactories` in `tencentcloud/framework/registry.go` |
| Docs sibling file | `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` |

### D3. Client injection

Embed `fw.DataSourceWithConfigure` per the template **or** follow the existing repository precedent (`ssm` ephemeral resource) with an explicit `client TkeService` field set in `Configure`.

Decision: mirror the existing, compiling repository pattern — explicit field + `Configure` that type-asserts `req.ProviderData` to `*sharedmeta.ProviderMeta` and builds `NewTkeService(meta.Client)`. Rationale: `fw.DataSourceWithConfigure` does not exist yet in this repository (grep for `DataSourceWithConfigure` returns nothing), while `sharedmeta.ProviderMeta` + `New<Svc>Service` is proven to compile. The template's `d.Client()` usage must therefore be replaced by the service field.

```go
var _ datasource.DataSource = &TkeAddonValuesDataSource{}
var _ datasource.DataSourceWithConfigure = &TkeAddonValuesDataSource{}

type TkeAddonValuesDataSource struct {
    client TkeService
}

func (d *TkeAddonValuesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
    if !ok {
        resp.Diagnostics.AddError(...)
        return
    }
    d.client = NewTkeService(meta.Client)
}
```

### D4. Schema

| Attribute | Type | Mode | Mapping |
| --- | --- | --- | --- |
| `cluster_id` | `schema.StringAttribute` | `Required` | `request.ClusterId` |
| `addon_name` | `schema.StringAttribute` | `Required` | `request.AddonName` |
| `values` | `schema.StringAttribute` | `Computed` | `response.Response.Values` |
| `default_values` | `schema.StringAttribute` | `Computed` | `response.Response.DefaultValues` |

No `id`, no `result_output_file` (Framework data sources do not need them; the SDKv2 `result_output_file` convention is not part of the template). Inputs are strictly `Required` (no `Computed` on inputs, per template point 1).

Model:

```go
type TkeAddonValuesDataSourceModel struct {
    ClusterId     types.String `tfsdk:"cluster_id"`
    AddonName     types.String `tfsdk:"addon_name"`
    Values        types.String `tfsdk:"values"`
    DefaultValues types.String `tfsdk:"default_values"`
}
```

### D5. Read flow and empty-result handling

`DescribeAddonValues` returns a **single object with flat outputs** (not a list). Per `6-field-types-and-computed.md` §4, single-object queries may error when nothing is found; but here "not found" is expressed by an **empty `Values`** (API doc: "未安装addon时返回为空值"), not by a nil response. Therefore:

| Condition | Handling |
| --- | --- |
| `err != nil` | `AddError("Error querying TKE addon values", fmt.Sprintf("DescribeAddonValues failed: %v", err))` |
| `response == nil \|\| response.Response == nil` | `AddError(...)` — API anomaly |
| `Values` nil or `""` (addon not installed) | **Not an error**: write `types.StringNull()`-equivalent (leave `types.StringValue("")` / `types.StringPointerValue(nil)`) and return normally; `default_values` is still populated |
| `DefaultValues` nil | write null via `types.StringPointerValue(nil)` |

Write-back uses `types.StringPointerValue(item.Values)` / `types.StringPointerValue(item.DefaultValues)` so a nil pointer becomes Null with no zero-value ambiguity.

Pseudo-code skeleton:

```go
func (d *TkeAddonValuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var config TkeAddonValuesDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() { return }

    request := tke.NewDescribeAddonValuesRequest()
    request.ClusterId = config.ClusterId.ValueStringPointer()
    request.AddonName = config.AddonName.ValueStringPointer()

    response, err := d.client.DescribeAddonValues(ctx, request)
    if err != nil { resp.Diagnostics.AddError(...); return }
    if response == nil || response.Response == nil { resp.Diagnostics.AddError(...); return }

    config.Values = types.StringPointerValue(response.Response.Values)
    config.DefaultValues = types.StringPointerValue(response.Response.DefaultValues)

    resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
```

### D6. Service layer

Add to `tencentcloud/services/tke/service_tencentcloud_tke.go` (package-private helper following existing TKE style, using `tccommon.GetLogId` + `log.Printf` on error):

```go
func (me *TkeService) DescribeAddonValues(ctx context.Context, request *tke.DescribeAddonValuesRequest) (response *tke.DescribeAddonValuesResponse, errRet error)
```

Single call, no pagination (API has no limit/offset). Do **not** add `resource.Retry` in the service method: the data source is read-only and the client already performs SDK-level retries; keeping the service call single-shot avoids nested retries.

### D7. Vendor API reference (verbatim excerpts)

`vendor/.../tke/v20180525/models.go` lines 7293–7358:

```go
// Predefined struct for user
type DescribeAddonValuesRequestParams struct {
    // 集群ID，请从容器服务控制台集群列表中获取（https://console.cloud.tencent.com/tke2/cluster）。
    ClusterId *string `json:"ClusterId,omitnil,omitempty" name:"ClusterId"`

    // addon名称，请从容器服务控制台进入集群详情页后，在组件管理页面中获取。
    AddonName *string `json:"AddonName,omitnil,omitempty" name:"AddonName"`
}

type DescribeAddonValuesRequest struct {
    *tchttp.BaseRequest

    ClusterId *string `json:"ClusterId,omitnil,omitempty" name:"ClusterId"`
    AddonName *string `json:"AddonName,omitnil,omitempty" name:"AddonName"`
}

// Predefined struct for user
type DescribeAddonValuesResponseParams struct {
    // 参数列表，如果addon已安装，会使用已设置的参数和chart里的默认参数做渲染，是一个json格式的字符串，未安装addon时返回为空值。
    Values *string `json:"Values,omitnil,omitempty" name:"Values"`

    // addon支持的参数列表，值为chart的默认值，是一个json格式的字符串。
    DefaultValues *string `json:"DefaultValues,omitnil,omitempty" name:"DefaultValues"`

    // 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeAddonValuesResponse struct {
    *tchttp.BaseResponse
    Response *DescribeAddonValuesResponseParams `json:"Response"`
}
```

`vendor/.../tke/v20180525/client.go` lines 5763–5792: signature and error code.

```go
// DescribeAddonValues
// 获取一个addon的参数
//
// 可能返回的错误码:
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
func (c *Client) DescribeAddonValues(request *DescribeAddonValuesRequest) (response *DescribeAddonValuesResponse, err error)

func (c *Client) DescribeAddonValuesWithContext(ctx context.Context, request *DescribeAddonValuesRequest) (response *DescribeAddonValuesResponse, err error)
```

Constructors: `tke.NewDescribeAddonValuesRequest()`, `tke.NewDescribeAddonValuesResponse()`.

Both request fields are `*string` and optional at the SDK level, but the cloud requires them — hence `Required` in the schema.

### D8. Framework data source template key structure (to follow at implementation time)

From `tfpacer-tf-templates/references/2-data-source.md`:

- Interface: `datasource.DataSource` — only `Metadata` / `Schema` / `Read`.
- `Metadata` sets `resp.TypeName = "tencentcloud_tke_addon_values"`.
- `Schema` uses `schema.Schema{Description, Attributes}` from `terraform-plugin-framework/datasource/schema`.
- Inputs: only `Required`/`Optional` (never `Computed`). Outputs: only `Computed`.
- `Read`: `req.Config.Get` → build request → API call → `resp.State.Set`.
- Sibling `<stem>.md` doc + registration in `dataSourceFactories` + `tencentcloud/provider.md`.
- Addendum (`6-field-types-and-computed.md`): `types.StringPointerValue()` for pointer write-back; check `IsNull()`/`IsUnknown()` before reading values; data-source outputs are all `Computed`; single-object query may error when absent (: here "absent" is an empty `values` string, which is a normal state, not an error).

### D9. Unit test strategy

Per project rules, a new data source gets a **mock-based (gomonkey) unit test**, not a Terraform acceptance suite:

- `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go`
- Monkey-patch `(*TkeService).DescribeAddonValues` (or the client call) to return canned responses.
- Cases: (a) addon installed — both `values` / `default_values` populated; (b) addon not installed — `Values` nil pointer → `values` Null, `default_values` populated, no error; (c) API error → diagnostics contain error; (d) `response.Response == nil` → diagnostics contain error.
- Assertions via `NewTkeAddonValuesDataSource().Schema(...)` and direct `Read` invocation with a constructed `datasource.ReadRequest`.

## Risks / Trade-offs

- [First Framework data source in the repo — no in-tree precedent] → Follow the proven `ssm` ephemeral resource structure (`Configure` + `sharedmeta.ProviderMeta`); validate against that file while implementing.
- [`fw.DataSourceWithConfigure` referenced by the template does not exist in this repo] → Use the explicit `client TkeService` field pattern instead (D3), documented here so implementation does not chase a non-existent helper.
- [`vendor/` shows uncommitted SDK diffs in the working tree] → Do not touch `go.mod`/`vendor/`; confirm `DescribeAddonValues` exists in the vendored copy (verified above, models.go:7294-7358).
- [`values` can be an empty string for a valid, non-installed addon] → Explicitly treated as a normal state, not an error (D5), keeping semantic parity with the cloud API doc.
- [JSON payload size] → `values` / `default_values` are plain `StringAttribute`; large JSON is valid, no length constraint is imposed. No `Sensitive` marking: the values are configuration defaults, not secrets.

## Migration Plan

1. Add the service-layer helper.
2. Add the data source file + `.md` sibling doc.
3. Register `tke.NewTkeAddonValuesDataSource` in `dataSourceFactories` and add the resource name in `tencentcloud/provider.md`.
4. Add the gomonkey unit test.
5. Docs generation (`make doc`) happens in the finalize stage; `website/` must not be edited manually.

Rollback: revert the commit; the change is additive only, so no state migration is involved.
