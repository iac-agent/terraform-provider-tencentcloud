## Why

Users who manage TKE (TencentCloud Kubernetes Engine) addons with Terraform need to know the parameters an addon supports before writing `raw_values`. Today `tencentcloud_kubernetes_addon` requires a base64-encoded JSON `raw_values` string, but the provider exposes no way to discover the addon's default chart values or the currently rendered values of an installed addon. The only way to obtain them is the console or a manual `DescribeAddonValues` API call, which breaks the infrastructure-as-code workflow.

## What Changes

- Add a new Terraform data source `tencentcloud_tke_addon_values`.
- The data source calls the TKE API `DescribeAddonValues` (version `2018-05-25`) with the required inputs `cluster_id` and `addon_name`.
- It exposes two computed outputs: `values` (the rendered parameter JSON string; empty when the addon is not installed) and `default_values` (the default chart parameter JSON string).
- Implemented with Terraform Plugin Framework, registered in `tencentcloud/framework/registry.go` (`dataSourceFactories`) and documented in `tencentcloud/provider.md`.
- No existing resource, data source, or state is modified, so the change is fully backward compatible.

## Capabilities

### New Capabilities

- `tke-addon-values-datasource`: Querying the supported/default parameter values and the rendered parameter values of a TKE addon in a given cluster.

### Modified Capabilities

<!-- None: this change only introduces a new data source. -->

## Impact

- **Affected service**: TKE (`tencentcloud/services/tke/`)
- **New files**:
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.md`
- **Modified files**:
  - `tencentcloud/framework/registry.go` (add `tke.NewTkeAddonValuesDataSource` to `dataSourceFactories`; this is the first framework data source in the repository)
  - `tencentcloud/provider.md` (documentation registration under the TKE node)
- **New service-layer method**: `DescribeAddonValues` wrapper in `tencentcloud/services/tke/service_tencentcloud_tke.go`
- **API dependency**: `DescribeAddonValues` in
  `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525`
  (already present in `vendor/`, no `go.mod` change required)
- **Possible error code**: `ResourceUnavailable`
- **Docs**: `website/docs/d/tke_addon_values.html.markdown` is generated later by `make doc`
