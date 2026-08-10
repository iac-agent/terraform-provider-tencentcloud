---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token_v9"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token_v9"
description: |-
  Provides a resource to create a TEO inference API token.
---

# tencentcloud_teo_inference_api_token_v9

Provides a resource to create a TEO inference API token.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token_v9" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  name    = "my-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the inference API token, limited to 30 characters.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - Inference API Token content. Only returned once during creation, subsequent queries will not return this value.
* `token_id` - Inference API Token ID.


