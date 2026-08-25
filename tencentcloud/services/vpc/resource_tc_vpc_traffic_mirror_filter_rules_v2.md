Provides a resource to manage VPC traffic mirror filter rules.

Example Usage

```hcl
resource "tencentcloud_vpc_traffic_mirror_filter_rules_v2" "example" {
  traffic_mirror_id = "imgf-xxxxxx"

  ingress_filter_rules {
    src_net    = "10.0.0.0/24"
    dst_net    = "10.0.1.0/24"
    protocol   = "TCP"
    src_port   = "80"
    dst_port   = "8080"
    priority   = 1
    action     = "ACCEPT"
    description = "ingress rule"
  }

  egress_filter_rules {
    src_net    = "10.0.1.0/24"
    dst_net    = "10.0.0.0/24"
    protocol   = "TCP"
    src_port   = "8080"
    dst_port   = "80"
    priority   = 1
    action     = "ACCEPT"
    description = "egress rule"
  }
}
```

Import

VPC traffic mirror filter rules can be imported using the traffic mirror id, e.g.

```
terraform import tencentcloud_vpc_traffic_mirror_filter_rules_v2.example imgf-xxxxxx
```