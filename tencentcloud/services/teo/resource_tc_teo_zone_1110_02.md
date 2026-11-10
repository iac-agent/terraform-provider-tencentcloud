---
subcategory: "TencentCloud EdgeOne (TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_1110_02"
description: |-
  Provides a resource to create a zone in TencentCloud EdgeOne (TEO).
---

# tencentcloud_teo_zone_1110_02

Provides a resource to create a zone in TencentCloud EdgeOne (TEO).

## Example Usage

```hcl
resource "tencentcloud_teo_zone_1110_02" "example" {
  zone_name = "example.com"
  type      = "full"
  area      = "overseas"
  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}

resource "tencentcloud_teo_zone_1110_02" "example_with_plan" {
  zone_name = "example2.com"
  type      = "partial"
  area      = "overseas"
  plan_id   = "plan-xxxxxx"
}

resource "tencentcloud_teo_zone_1110_02" "example_with_tags" {
  zone_name = "example3.com"
  type      = "partial"
  area      = "mainland"
  tags {
    tag_key   = "env"
    tag_value = jsonencode({ "key1" = "value1", "key2" = "value2" })
  }
}
```

## Import

teo zone can be imported using the zone id, e.g.

```
terraform import tencentcloud_teo_zone_1110_02.example zone-2noz78a8ev6k
```
