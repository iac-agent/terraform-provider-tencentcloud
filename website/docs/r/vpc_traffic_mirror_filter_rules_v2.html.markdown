---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_traffic_mirror_filter_rules_v2"
sidebar_current: "docs-tencentcloud-resource-vpc_traffic_mirror_filter_rules_v2"
description: |-
  Provides a resource to manage VPC traffic mirror filter rules.
---

# tencentcloud_vpc_traffic_mirror_filter_rules_v2

Provides a resource to manage VPC traffic mirror filter rules.

## Example Usage

```hcl
resource "tencentcloud_vpc_traffic_mirror_filter_rules_v2" "example" {
  traffic_mirror_id = "imgf-xxxxxx"

  ingress_filter_rules {
    src_net     = "10.0.0.0/24"
    dst_net     = "10.0.1.0/24"
    protocol    = "TCP"
    src_port    = "80"
    dst_port    = "8080"
    priority    = 1
    action      = "ACCEPT"
    description = "ingress rule"
  }

  egress_filter_rules {
    src_net     = "10.0.1.0/24"
    dst_net     = "10.0.0.0/24"
    protocol    = "TCP"
    src_port    = "8080"
    dst_port    = "80"
    priority    = 1
    action      = "ACCEPT"
    description = "egress rule"
  }
}
```

## Argument Reference

The following arguments are supported:

* `traffic_mirror_id` - (Required, String, ForceNew) Traffic mirror instance ID.
* `egress_filter_rules` - (Optional, List) Egress filter rules.
* `ingress_filter_rules` - (Optional, List) Ingress filter rules.

The `egress_filter_rules` object supports the following:

* `action` - (Optional, String) Traffic mirror filter rule policy, support types: `ACCEPT`, `DROP`.
* `description` - (Optional, String) Traffic mirror filter rule description.
* `dst_net` - (Optional, String) Destination network segment of filter rule.
* `dst_port` - (Optional, String) Destination port of filter rule, default 1-65535.
* `priority` - (Optional, Int) Traffic mirror filter rule priority.
* `protocol` - (Optional, String) Protocol of filter rule.
* `src_net` - (Optional, String) Source network segment of filter rule.
* `src_port` - (Optional, String) Source port of filter rule, default 1-65535.

The `ingress_filter_rules` object supports the following:

* `action` - (Optional, String) Traffic mirror filter rule policy, support types: `ACCEPT`, `DROP`.
* `description` - (Optional, String) Traffic mirror filter rule description.
* `dst_net` - (Optional, String) Destination network segment of filter rule.
* `dst_port` - (Optional, String) Destination port of filter rule, default 1-65535.
* `priority` - (Optional, Int) Traffic mirror filter rule priority.
* `protocol` - (Optional, String) Protocol of filter rule.
* `src_net` - (Optional, String) Source network segment of filter rule.
* `src_port` - (Optional, String) Source port of filter rule, default 1-65535.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPC traffic mirror filter rules can be imported using the traffic mirror id, e.g.

```
terraform import tencentcloud_vpc_traffic_mirror_filter_rules_v2.example imgf-xxxxxx
```

