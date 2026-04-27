---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_line"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway_line"
description: |-
  Provides a resource to create a TEO multi-path gateway line.
---

# tencentcloud_teo_multi_path_gateway_line

Provides a resource to create a TEO multi-path gateway line.

## Example Usage

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  gateway_id   = "gw-abcdefgh"
  line_type    = "custom"
  line_address = "1.2.3.4:8080"
}
```



```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "proxy_example" {
  zone_id      = "zone-2qtuhspy7cr6"
  gateway_id   = "gw-abcdefgh"
  line_type    = "proxy"
  line_address = "5.6.7.8:443"
  proxy_id     = "sid-xxxxxxxx"
  rule_id      = "rule-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Multi-path gateway ID.
* `line_address` - (Required, String) Line address, format: ip:port.
* `line_type` - (Required, String) Line type. Valid values: direct, proxy, custom.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `proxy_id` - (Optional, String) Layer 4 proxy instance ID, required when line_type is proxy.
* `rule_id` - (Optional, String) Forwarding rule ID, required when line_type is proxy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `line_id` - Line ID.


## Import

TEO multi-path gateway line can be imported using the joint id "zone_id#gateway_id#line_id", e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-2qtuhspy7cr6#gw-abcdefgh#line-2
```

