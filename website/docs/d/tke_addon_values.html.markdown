---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tke_addon_values"
sidebar_current: "docs-tencentcloud-data_source-tke_addon_values"
description: |-
  Use this data source to query the parameter values of a TKE addon.
---

# tencentcloud_tke_addon_values

Use this data source to query the parameter values of a TKE addon.

## Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}

output "addon_default_values" {
  value = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}
```

### :

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  addon_name = "ip-masq-agent"
}

resource "tencentcloud_kubernetes_addon" "example" {
  cluster_id    = tencentcloud_kubernetes_cluster.example.id
  addon_name    = "ip-masq-agent"
  addon_version = "v1.0.0"

  raw_values = base64encode(jsonencode(merge(
    jsondecode(data.tencentcloud_tke_addon_values.example.default_values),
    {
      "nodeSelector" = {
        "node-role.kubernetes.io/agent" = ""
      }
    }
  )))
}
```

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String) Name of the addon, which can be obtained from the component management page of the cluster.
* `cluster_id` - (Required, String) ID of the TKE cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_values` - Parameter JSON string supported by the addon, with values taken from the chart defaults.
* `values` - Parameter JSON string of the addon. When the addon is installed, the user-set parameters are rendered together with the chart default parameters; when the addon is not installed, the value is empty.


