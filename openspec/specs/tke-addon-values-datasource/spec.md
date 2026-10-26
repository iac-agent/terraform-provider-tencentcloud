# tke-addon-values-datasource Specification

## Purpose
Provide the `tencentcloud_tke_addon_values` data source, which queries a TKE addon's parameter values through the TKE `DescribeAddonValues` API, exposing both the rendered parameters of an installed addon and the chart default parameters.

## Requirements

### Requirement: Data Source Schema Definition

The provider SHALL provide a data source named `tencentcloud_tke_addon_values` that queries the parameters of a TKE addon via the TKE API `DescribeAddonValues`.

The data source MUST define the following query inputs:

- `cluster_id` (String, Required): the ID of the TKE cluster.
- `addon_name` (String, Required): the name of the addon, as shown on the cluster's component management page.

The data source MUST define the following computed outputs:

- `values` (String, Computed): the rendered parameter JSON string of the addon. When the addon is installed, the cloud renders user-set parameters together with the chart default parameters; when the addon is not installed, the returned value is empty.
- `default_values` (String, Computed): the parameter JSON string supported by the addon, with values taken from the chart defaults.

Query inputs MUST NOT be marked Computed, and computed outputs MUST NOT be Required or Optional.

#### Scenario: Query an installed addon

- **WHEN** a user declares the data source with a valid `cluster_id` and the `addon_name` of an addon installed in that cluster
- **THEN** `values` contains the rendered parameter JSON string returned by the cloud
- **AND** `default_values` contains the default parameter JSON string returned by the cloud

#### Scenario: Query an addon that is not installed

- **WHEN** a user declares the data source with a valid `cluster_id` and an `addon_name` that is not installed in that cluster
- **THEN** the read completes successfully without an error
- **AND** `values` is empty (null/empty string)
- **AND** `default_values` still contains the default parameter JSON string

#### Scenario: Missing required input

- **WHEN** a user omits `cluster_id` or `addon_name` from the data source configuration
- **THEN** Terraform fails validation before any cloud API call, reporting the missing required argument

### Requirement: Read Behavior and Error Handling

The data source MUST implement only the `Read` operation of `datasource.DataSource` (together with `Metadata` and `Schema`) and MUST NOT modify any cloud-side state.

On `Read` the data source MUST call `DescribeAddonValues` exactly once per invocation, mapping `cluster_id` to `request.ClusterId` and `addon_name` to `request.AddonName`, then write `response.Response.Values` and `response.Response.DefaultValues` back to state.

The data source MUST surface the following conditions as diagnostics errors instead of silently returning empty values:

- the API call returns an error (including the `ResourceUnavailable` error code);
- the API response or its `Response` payload is nil.

An empty `values` string MUST NOT be treated as an error, because the cloud uses it to signal that the addon is not installed.

#### Scenario: API returns an error

- **WHEN** the `DescribeAddonValues` call fails, for example because the cluster does not exist or the caller lacks permission
- **THEN** the data source returns an error diagnostic that names the failing API and includes the underlying error message
- **AND** no partial values are written to the result

#### Scenario: API returns a nil response payload

- **WHEN** the cloud returns a response whose `Response` payload is nil
- **THEN** the data source returns an error diagnostic
- **AND** the read does not silently succeed with empty outputs

#### Scenario: Repeated reads are idempotent

- **WHEN** the same configuration is read multiple times without any cloud-side change
- **THEN** every read returns the same `values` and `default_values`
- **AND** no cloud resource is created, modified, or deleted

### Requirement: Provider Registration

The data source MUST be implemented with Terraform Plugin Framework in `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`, with the factory function `NewTkeAddonValuesDataSource()`.

The factory MUST be registered in the `dataSourceFactories` slice of `tencentcloud/framework/registry.go`, and the data source name MUST be listed under the TKE node in `tencentcloud/provider.md`.

The provider MUST NOT require changes to `go.mod` or `vendor/`, because `DescribeAddonValues` is already available in the vendored TKE SDK package `tke/v20180525`.

#### Scenario: Provider resolves the data source name

- **WHEN** the provider is initialized and `tencentcloud_tke_addon_values` is used in a configuration
- **THEN** the data source is resolved from the framework `dataSourceFactories` registry
- **AND** the configured provider client is injected so that the read can call the TKE API

#### Scenario: Documentation is generated from the sibling markdown file

- **WHEN** `make doc` runs during the finalize stage
- **THEN** the generated documentation page for `tencentcloud_tke_addon_values` contains the example usage from `data_source_tc_tke_addon_values.md`
- **AND** no file under `website/` is edited manually
