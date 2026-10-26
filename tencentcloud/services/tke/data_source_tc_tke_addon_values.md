Use this data source to query the parameter values of a TKE cluster addon.

The data source returns the parameter values currently rendered for the addon (with the values already set on it) and the default parameter values supported by the addon chart. Both are JSON strings, so `jsondecode` can be used to read individual parameters.

Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "cbs"
}

output "addon_default_values" {
  value = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}

output "addon_values" {
  value = jsondecode(data.tencentcloud_tke_addon_values.example.values)
}
```
