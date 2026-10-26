## 1. Service Layer

- [x] 1.1 Add `func (me *TkeService) DescribeAddonValues(ctx context.Context, request *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error)` to `tencentcloud/services/tke/service_tencentcloud_tke.go`; call `me.client.UseTkeClient().DescribeAddonValues(request)` once (no pagination, no nested retry) and log the request body plus error via `log.Printf("[CRITAL]%s api[%s] fail, ...")` using `tccommon.GetLogId(ctx)`

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values.go` with interface assertions for `datasource.DataSource` and `datasource.DataSourceWithConfigure`, the factory `NewTkeAddonValuesDataSource()`, struct `TkeAddonValuesDataSource{ client TkeService }`, and model `TkeAddonValuesDataSourceModel` (`cluster_id`, `addon_name`, `values`, `default_values`, all `types.String`)
- [x] 2.2 Implement `Metadata` with `resp.TypeName = "tencentcloud_tke_addon_values"`
- [x] 2.3 Implement `Schema`: `cluster_id` and `addon_name` as `schema.StringAttribute{Required: true}`, `values` and `default_values` as `schema.StringAttribute{Computed: true}`; schema description starts with "Use this data source to query TKE addon parameter values."
- [x] 2.4 Implement `Configure`: type-assert `req.ProviderData` to `*sharedmeta.ProviderMeta`, build `d.client = NewTkeService(meta.Client)`, and emit an error diagnostic for an unexpected provider data type
- [x] 2.5 Implement `Read`: `req.Config.Get` into the model, set `request.ClusterId = config.ClusterId.ValueStringPointer()` and `request.AddonName = config.AddonName.ValueStringPointer()`, call `d.client.DescribeAddonValues`, then check `err != nil` and `response == nil || response.Response == nil` with `resp.Diagnostics.AddError`
- [x] 2.6 Write outputs back with `types.StringPointerValue(response.Response.Values)` and `types.StringPointerValue(response.Response.DefaultValues)`, leaving an empty `values` as a normal (null) value rather than an error, then `resp.State.Set(ctx, &config)`

## 3. Registration

- [x] 3.1 Add the `tke` import to `tencentcloud/framework/registry.go` and append `tke.NewTkeAddonValuesDataSource` to `dataSourceFactories`
- [x] 3.2 Add `tencentcloud_tke_addon_values` under the TKE node in `tencentcloud/provider.md`

## 4. Documentation Sibling File

- [x] 4.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` with a one-sentence `Use this data source to query ...` description naming TKE, an Example Usage block using the existing cluster/addon resources and showing `jsonencode()` where the JSON strings are consumed, and no `Argument Reference` / `Attribute Reference` sections

## 5. Unit Test

- [x] 5.1 Create `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go` using gomonkey (no Terraform acceptance suite): patch the `DescribeAddonValues` service method or the TKE client call
- [x] 5.2 Cover cases: addon installed (both outputs populated), addon not installed (`Values` nil pointer → `values` null, `default_values` populated, no diagnostics error), API error, and `response.Response == nil`; assert diagnostics and model values

## 6. Verification (no build/test execution in this stage)

- [x] 6.1 Self-review the generated Go code against `design.md` for a correct construction of `tke.NewDescribeAddonValuesRequest`, complete error handling on every returned error, and no unused imports or variables
- [x] 6.2 Confirm `.md` sibling file and registry entry are consistent with the Terraform name `tencentcloud_tke_addon_values`
- [x] 6.3 Confirm `website/` and `.changelog/` were not touched (they are produced in the finalize stage)
