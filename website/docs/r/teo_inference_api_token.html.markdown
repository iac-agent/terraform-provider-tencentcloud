---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token"
description: |-
  Provides a resource to create a TEO Inference API Token
---

# tencentcloud_teo_inference_api_token

Provides a resource to create a TEO Inference API Token

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-2o3h21ed2t68"
  name    = "example-token"

  offset = 0
  limit  = 20
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Inference API Token name, cannot exceed 30 characters.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `limit` - (Optional, Int) Pagination query limit. Default value: 20, maximum value: 100.
* `offset` - (Optional, Int) Pagination query offset. Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - Inference API Token content.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `token_id` - Inference API Token ID.
* `total_count` - Total number of Inference API Tokens.


## Import

TEO Inference API Token can be imported using the zone_id#token_id, e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-2o3h21ed2t68#token-xxxxx
```

