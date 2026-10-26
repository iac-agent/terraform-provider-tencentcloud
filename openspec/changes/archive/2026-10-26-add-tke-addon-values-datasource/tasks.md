## 1. Service Layer

- [x] 1.1 In `tencentcloud/services/tke/service_tencentcloud_tke.go`, add
      `func (me *TkeService) DescribeAddonValues(ctx context.Context, clusterId, addonName string) (values *string, defaultValues *string, errRet error)`
      using `tke.NewDescribeAddonValuesRequest()`, `request.ClusterId`, `request.AddonName`, and
      `me.client.UseTkeClient().DescribeAddonValuesWithContext(ctx, request)` (skeleton in design.md D5)
- [x] 1.2 Add the `[CRITAL]` deferred failure log (`tccommon.GetLogId(ctx)`, `request.GetAction()`,
      `request.ToJsonString()`) matching the other `Describe*` helpers in the file
- [x] 1.3 Guard `response == nil || response.Response == nil` and return a non-nil `errRet`
      (`"DescribeAddonValues returned an empty response"`)
- [x] 1.4 Return `response.Response.Values` / `response.Response.DefaultValues` as raw `*string`
      (no nil collapsing, no JSON parsing)

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values.go` with `package tke`
      and imports: `context`, `fmt`, `datasource`, `datasource/schema`, `types`,
      `internal/sharedmeta`
- [x] 2.2 Define the factory `NewTkeAddonValuesDataSource() datasource.DataSource` returning
      `&TkeAddonValuesDataSource{}`, plus the interface assertion
      `var _ datasource.DataSourceWithConfigure = &TkeAddonValuesDataSource{}`
- [x] 2.3 Define `TkeAddonValuesDataSource` with the `client TkeService` field (design.md D2 — no
      `fw.DataSourceWithConfigure` helper exists in this repo; follow
      `services/ssm/ephemeral_tc_ssm_secret_version.go`)
- [x] 2.4 Define `TkeAddonValuesDataSourceModel` with `tfsdk` tags `cluster_id`, `addon_name`,
      `values`, `default_values` (design.md D6)
- [x] 2.5 Implement `Metadata` setting `resp.TypeName = "tencentcloud_tke_addon_values"`
- [x] 2.6 Implement `Schema` with `cluster_id`/`addon_name` as `Required` string attributes and
      `values`/`default_values` as `Computed` string attributes, all with English descriptions
      (design.md D4); no `Sensitive`, `Default`, `Validators` or plan modifiers
- [x] 2.7 Implement `Configure` type-asserting `req.ProviderData` to `*sharedmeta.ProviderMeta`,
      nil-guarding the input, and setting `d.client = NewTkeService(meta.Client)`; emit
      `AddError("Unexpected ProviderData type", ...)` on a failed assertion
- [x] 2.8 Implement `Read`: `req.Config.Get`, diagnostic early-return, explicit empty-string checks
      for `cluster_id` and `addon_name`, call `d.client.DescribeAddonValues`, wrap API errors with
      `AddError` including the API name plus cluster/addon, map results with
      `types.StringPointerValue` (nil → Null), then `resp.State.Set`
- [x] 2.9 Extract the response-to-model mapping into a small unexported helper so it can be unit
      tested without a live client (design.md D8 option 2)

## 3. Registration and Documentation Index

- [x] 3.1 In `tencentcloud/framework/registry.go`, add the
      `services/tke` import and the single entry `tke.NewTkeAddonValuesDataSource` inside
      `dataSourceFactories`
- [x] 3.2 In `tencentcloud/provider.md`, add `tencentcloud_tke_addon_values` under
      `Tencent Kubernetes Engine(TKE)` -> `Data Source`
- [x] 3.3 Do NOT modify `tencentcloud/provider.go` or anything under `website/`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` with a one-line
      description beginning `Use this data source to query ...` that names TKE
- [x] 4.2 Add an `Example Usage` HCL block using `cluster_id` + `addon_name`, plus a
      `jsondecode(...default_values)` output example (design.md D7)
- [x] 4.3 Ensure the `.md` contains no `Argument Reference` / `Attribute Reference` / `Import`
      sections (auto-generated later by `make doc`)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go` using gomonkey
      mocks only (no Terraform acceptance suite, no `TF_ACC`)
- [x] 5.2 Cover the mapping helper: both fields set, `Values` nil, `DefaultValues` nil, both nil
- [x] 5.3 Optionally cover `Read` end-to-end by mocking
      `(*tke.Client).DescribeAddonValuesWithContext` if the shared client can be injected in a unit
      test; skip cleanly otherwise
- [x] 5.4 Ensure all returned `error`/`diag.Diagnostics` values are checked or explicitly discarded
      so the package still compiles

## 6. Verification

- [x] 6.1 Re-read design.md and confirm every schema attribute, tag and registration matches it
- [x] 6.2 Confirm each cloud API parameter used exists in the vendored SDK
      (`ClusterId`, `AddonName` on the request; `Values`, `DefaultValues` on the response params)
- [x] 6.3 Confirm no `Go build`/`go vet`/`golint`/`gofmt`/`make doc`/changelog work is done in this
      phase (handled by the finalize step)
