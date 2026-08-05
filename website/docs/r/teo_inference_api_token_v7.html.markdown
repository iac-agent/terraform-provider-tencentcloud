---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token_v7"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token_v7"
description: |-
  Provides a resource to create a TEO inference API token
---

# tencentcloud_teo_inference_api_token_v7

Provides a resource to create a TEO inference API token

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token_v7" "example" {
  zone_id = "zone-27q0p0bali16"
  name    = "my-inference-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the inference API token, with a maximum length of 30 characters.
* `zone_id` - (Required, String, ForceNew) The zone ID where the inference API token belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - The content of the inference API token.
* `token_id` - The unique ID of the inference API token.


## Import

The resource can be imported by using the `token_id`, e.g.

```sh
terraform import tencentcloud_teo_inference_api_token_v7.example <token_id>
```

