Provides a resource to manage TEO Edge KV namespace configuration

Example Usage

```hcl
resource "tencentcloud_teo_edge_kv_namespace_config" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  remark    = "This is an example namespace config"
}
```

Import

TEO Edge KV namespace config can be imported using the zone_id#namespace, e.g.

```
terraform import tencentcloud_teo_edge_kv_namespace_config.example zone-2o3h21ed2t68#example-namespace
```