Use this data source to query the parameter values of a TKE addon.

Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}

output "addon_default_values" {
  value = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}
```

Use the returned JSON to build `raw_values` of `tencentcloud_kubernetes_addon`:

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  addon_name = "ip-masq-agent"
}

resource "tencentcloud_kubernetes_addon" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  addon_name = "ip-masq-agent"
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
