---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_kv_namespace_config"
sidebar_current: "docs-tencentcloud-resource-teo_edge_kv_namespace_config"
description: |-
  Provides a resource to manage TEO Edge KV namespace configuration
---

# tencentcloud_teo_edge_kv_namespace_config

Provides a resource to manage TEO Edge KV namespace configuration

## Example Usage

```hcl
resource "tencentcloud_teo_edge_kv_namespace_config" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  remark    = "This is an example namespace config"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Namespace name. Supports 1-50 characters, allowed characters are a-z, A-Z, 0-9, -, and - cannot be used alone or consecutively, cannot be placed at the beginning or end. The name must be unique within the same site.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `remark` - (Optional, String) Namespace description. Used to describe the purpose or business meaning of the namespace. Maximum 256 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `capacity_used` - KV storage used capacity in bytes.
* `capacity` - KV storage available capacity in bytes.
* `created_on` - Creation time in ISO 8601 format (UTC).
* `modified_on` - Last modification time in ISO 8601 format (UTC).


## Import

TEO Edge KV namespace config can be imported using the zone_id#namespace, e.g.

```
terraform import tencentcloud_teo_edge_kv_namespace_config.example zone-2o3h21ed2t68#example-namespace
```

