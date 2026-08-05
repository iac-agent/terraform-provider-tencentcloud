---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token"
description: |-
  Provides a resource to create a TEO inference API token
---

# tencentcloud_teo_inference_api_token

Provides a resource to create a TEO inference API token

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-xxxxxxxx"
  name    = "my-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Inference API Token name, with a length limit of no more than 30 characters.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - Inference API Token content.
* `create_time` - Creation time, in ISO date format.
* `token_id` - Inference API Token ID.


## Import

TEO inference API token can be imported using the zone ID and token ID, separated by "#", e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-xxxxxxxx#token-xxxxxxxx
```

