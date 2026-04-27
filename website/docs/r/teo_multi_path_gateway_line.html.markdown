---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_line"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway_line"
description: |-
  Provides a resource to create a TEO multi-path gateway line
---

# tencentcloud_teo_multi_path_gateway_line

Provides a resource to create a TEO multi-path gateway line

## Example Usage

### custom line

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "custom_line" {
  zone_id      = "zone-2xkazzl8yf6k"
  gateway_id   = "gw-3lchxitnb5pb"
  line_type    = "custom"
  line_address = "1.2.3.4:8080"
}
```

### proxy line

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "proxy_line" {
  zone_id      = "zone-2xkazzl8yf6k"
  gateway_id   = "gw-3lchxitnb5pb"
  line_type    = "proxy"
  line_address = "1.2.3.4:443"
  proxy_id     = "sid-2xzwkzljmm9b"
  rule_id      = "rule-2xzwkzljmm9b"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Multi-path gateway ID.
* `line_address` - (Required, String) Line address, format is ip:port.
* `line_type` - (Required, String) Line type. Valid values: direct, proxy, custom.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `proxy_id` - (Optional, String) L4 proxy instance ID. Required when line_type is proxy.
* `rule_id` - (Optional, String) Forwarding rule ID. Required when line_type is proxy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `line_id` - Line ID.


## Import

TEO multi-path gateway line can be imported using the zoneId#gatewayId#lineId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-2xkazzl8yf6k#gw-3lchxitnb5pb#line-2
```

