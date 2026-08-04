---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token"
description: |-
  Provides a resource to create an inference API token for TEO (EdgeOne) zones.
---

# tencentcloud_teo_inference_api_token

Provides a resource to create an inference API token for TEO (EdgeOne) zones.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-12345678"
  name    = "my-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Inference API Token name, length limit does not exceed 30 characters.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - Inference API Token content.
* `token_id` - Inference API Token ID.


