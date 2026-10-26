---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tke_addon_values"
sidebar_current: "docs-tencentcloud-data_source-tke_addon_values"
description: |-
  Use this data source to query the parameter values of a TKE cluster addon.
---

# tencentcloud_tke_addon_values

Use this data source to query the parameter values of a TKE cluster addon.

The data source returns the parameter values currently rendered for the addon (with the values already set on it) and the default parameter values supported by the addon chart. Both are JSON strings, so `jsondecode` can be used to read individual parameters.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String) Name of the addon, as listed in the cluster's addon management page.
* `cluster_id` - (Required, String) ID of the TKE cluster the addon belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_values` - Default parameter values supported by the addon, as a JSON string.
* `values` - Parameter values of the addon rendered with the values already set, as a JSON string. Empty or null when the addon is not installed.


