---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_tokens"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_api_tokens"
description: |-
  Use this data source to query the inference API tokens of the specified TEO zone, including the token ID, name, content and creation time.
---

# tencentcloud_teo_inference_api_tokens

Use this data source to query the inference API tokens of the specified TEO zone, including the token ID, name, content and creation time.

## Example Usage

```hcl
data "tencentcloud_teo_inference_api_tokens" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `content` - Inference API Token content.
  * `create_time` - Creation time in ISO date format.
  * `name` - Inference API Token name.
  * `token_id` - Inference API Token ID.
* `tokens` - List of inference API tokens.


