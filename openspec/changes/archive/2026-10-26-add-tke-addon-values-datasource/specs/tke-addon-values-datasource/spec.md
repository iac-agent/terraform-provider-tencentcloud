## ADDED Requirements

### Requirement: Query TKE addon parameter values

The provider SHALL expose a data source `tencentcloud_tke_addon_values` that reads addon parameter
information for a TKE cluster by calling the TKE v20180525 API `DescribeAddonValues`. The data source
MUST require `cluster_id` and `addon_name` as inputs and MUST expose `values` and `default_values` as
computed JSON-string outputs. The data source MUST NOT mutate any cloud-side state.

#### Scenario: Addon installed and both value sets returned

- **WHEN** a configuration declares `cluster_id = "cls-5yezvaxo"` and `addon_name = "cbs"` and the
  addon is installed on that cluster
- **THEN** `terraform plan`/`apply` succeeds
- **AND** `values` contains the JSON string returned in `response.Response.Values`
- **AND** `default_values` contains the JSON string returned in `response.Response.DefaultValues`
- **AND** `jsondecode(data.tencentcloud_tke_addon_values.example.values)` yields an object

#### Scenario: Addon not installed

- **WHEN** the requested `addon_name` is not installed on the cluster and the API returns an empty
  `Values` field
- **THEN** the read succeeds without an error
- **AND** `values` is null in the resulting data
- **AND** `default_values` is still populated from `response.Response.DefaultValues`

#### Scenario: Both value fields absent

- **WHEN** the API returns a response whose `Values` and `DefaultValues` are both nil
- **THEN** the read succeeds without an error
- **AND** both `values` and `default_values` are null (not empty strings)

### Requirement: Required inputs are validated before the API call

The data source SHALL reject empty `cluster_id` or `addon_name` values with a diagnostic instead of
issuing an API call.

#### Scenario: Empty cluster_id

- **WHEN** a configuration sets `cluster_id = ""`
- **THEN** the read fails with a diagnostic mentioning `cluster_id`
- **AND** no `DescribeAddonValues` request is sent

#### Scenario: Empty addon_name

- **WHEN** a configuration sets `addon_name = ""`
- **THEN** the read fails with a diagnostic mentioning `addon_name`
- **AND** no `DescribeAddonValues` request is sent

### Requirement: API failures surface as actionable diagnostics

The data source SHALL return a Terraform diagnostic, including the API name and the queried cluster
and addon, whenever the API call fails or returns an empty response envelope. It MUST NOT panic or
silently succeed.

#### Scenario: API returns an error

- **WHEN** `DescribeAddonValues` fails with `ResourceUnavailable` (for example a nonexistent cluster)
- **THEN** the read fails with a diagnostic whose detail names `DescribeAddonValues` together with
  the cluster ID and addon name
- **AND** no partial data is written

#### Scenario: Empty response envelope

- **WHEN** `DescribeAddonValues` returns a nil error but a nil `Response` envelope
- **THEN** the read fails with a diagnostic stating the response was empty

### Requirement: Framework registration and documentation

The data source SHALL be implemented with Terraform Plugin Framework, registered through the
framework factory manifest, and documented for `make doc` generation.

#### Scenario: Registered in the framework manifest

- **WHEN** the provider binary starts
- **THEN** `tencentcloud/framework/registry.go` contains the `tke.NewTkeAddonValuesDataSource` factory
  in `dataSourceFactories`
- **AND** `tencentcloud/provider.go`'s SDKv2 maps are unchanged
- **AND** the provider starts mux without a duplicate type-name error for
  `tencentcloud_tke_addon_values`

#### Scenario: Documentation entry present

- **WHEN** the change is complete
- **THEN** `tencentcloud/provider.md` lists `tencentcloud_tke_addon_values` under
  `Tencent Kubernetes Engine(TKE)` in the `Data Source` section
- **AND** `tencentcloud/services/tke/data_source_tc_tke_addon_values.md` exists with a one-line
  description, an `Example Usage` HCL block, and no `Argument Reference`/`Attribute Reference`
  sections
