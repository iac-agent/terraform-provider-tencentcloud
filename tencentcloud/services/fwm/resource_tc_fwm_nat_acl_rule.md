Provides a resource to manage a FWM (Firewall Manager) NAT boundary ACL rule.

Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "acg-1a2b3c4d"
  source_content = "10.0.0.0/8"
  source_type    = "ip"
  target_content = "192.168.1.0/24"
  target_type    = "ip"
  protocol       = "TCP"
  rule_action    = "accept"
  order_index    = 1
  scope          = "ALL"
  direction      = 0
  port           = "80,443"
  description    = "tf-example"
}
```

Import

FWM NAT ACL rule can be imported using the composite id `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example acg-1a2b3c4d#aclr-9x8y7z6w
```
