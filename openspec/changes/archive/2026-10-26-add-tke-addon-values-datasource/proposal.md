## Why

TencentCloud TKE users currently have no Terraform-native way to read the parameter template and current parameter values of a cluster addon before installing or upgrading it. Without this data, HCL authors must hand-craft the `raw_values` JSON blob for `tencentcloud_kubernetes_addon` / `tencentcloud_kubernetes_addon_config` by copying values out of the console, which is error-prone and drifts from the chart defaults.

The cloud API `DescribeAddonValues` (TKE v20180525) already exposes exactly this information: it returns the chart's default parameter values and, when the addon is installed, the rendered parameter values currently in effect. Exposing it as a data source lets configurations derive addon parameters dynamically (e.g. read `default_values`, patch a few keys, and feed the result into `raw_values`), keeping addon configuration idempotent and reproducible.

## What Changes

- New data source `tencentcloud_tke_addon_values` (RESOURCE_KIND_DATASOURCE), implemented with Terraform Plugin Framework.
- Implements a single read path that calls TKE `DescribeAddonValues`.
  - Inputs (Required): `cluster_id`, `addon_name`
  - Outputs (Computed): `values`, `default_values`
- Registers the factory in `tencentcloud/framework/registry.go` -> `dataSourceFactories` and lists the reference under the TKE `Data Source` node of `tencentcloud/provider.md`.
- Ships sibling documentation `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` consumed by `make doc`.
- Adds a mock-based unit test (gomonkey) covering the read path; no SDKv2 acceptance-test suite is used because this is a new reference.

No existing schema, resource, or data source is modified, so there is no breaking change.

## Capabilities

### New Capabilities

- `tke-addon-values-datasource`: read-only query of a TKE cluster addon's parameter values via `DescribeAddonValues`, exposing the effective values and chart default values as computed attributes.

### Modified Capabilities

<!-- None: no existing capability requirements are changing. -->

## Impact

- **Affected service**: TKE (`tencentcloud/services/tke`, package `tke`).
- **New files**:
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.md`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go`
- **Modified files**:
  - `tencentcloud/framework/registry.go` (one factory line + package import if missing)
  - `tencentcloud/provider.md` (one data source entry under the TKE node)
- **Cloud API dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525`, method `DescribeAddonValues`, with documented error code `ResourceUnavailable`.
- **No `website/` edits**: generated later by `make doc`.
- **Runtime impact**: none on existing resources; mux startup only gains one additional data source type name (`tencentcloud_tke_addon_values`), which does not collide with any existing type name.
