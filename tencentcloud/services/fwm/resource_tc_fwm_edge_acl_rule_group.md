Provides a resource to create an internet border rule group of Firewall Manager (FWM)

The rule group is created together with the rules listed in it, and destroying the
resource deletes that rule group only: no other rule group of the account is touched,
because the underlying batch delete API is always called with the single id managed by
this resource. Only `group_name` can be updated in place; changing `rules` or any rule
field replaces the rule group.

Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule_group" "example" {
  group_name = "tf-example-edge-acl-group"

  rules {
    order_index    = 1
    direction      = 1
    source_content = "0.0.0.0/0"
    source_type    = "ip"
    target_content = "1.1.1.1"
    target_type    = "ip"
    protocol       = "TCP"
    port           = "80,443"
    rule_action    = "accept"
    scope          = "all"
    description    = "allow web traffic"
  }

  rules {
    order_index    = 2
    direction      = 1
    source_content = "0.0.0.0/0"
    source_type    = "ip"
    target_content = "1.1.1.1"
    target_type    = "ip"
    protocol       = "ANY"
    rule_action    = "drop"
    scope          = "all"
    description    = "drop the remaining traffic"
  }
}
```

Import

FWM internet border rule group can be imported using the group id, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule_group.example fwmrg-9z03vakjj8
```