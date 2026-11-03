Provides a resource to create a FWM NAT boundary (NAT ACL) rule

Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "fwmrg-in3q0qkq0n"
  direction      = 0
  source_content = "0.0.0.0/0"
  source_type    = "ip"
  target_content = "0.0.0.0/0"
  target_type    = "ip"
  protocol       = "ANY"
  rule_action    = "accept"
  order_index    = 1
  scope          = "ALL"
  port           = "-1/-1"
  description    = "remark."
}
```

Import

FWM NAT boundary rule can be imported using the `group_id#rule_id#direction` composite id, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-in3q0qkq0n#7520b898-30ba-415d-ab27-ff4542071ea9#0
```

~> **NOTE:** `group_id` has no dedicated field in the query response, so it is restored from the import id. The
configured `order_index` is only verified through the rule priority (`sequence`) returned by the cloud API.
