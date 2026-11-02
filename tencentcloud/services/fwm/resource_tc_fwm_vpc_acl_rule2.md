Provides a resource to create a FWM VPC border firewall ACL rule.

A rule group holds many rules and the cloud API creates and deletes rules in batches, but this resource always manages exactly one rule: both the create and the delete call carry a single element, and the rule is read back by its rule id.

~> **NOTE:** The VPC ACL query API has no read-back field for the rule group id and returns the priority as `sequence`. The resource keeps the configured `order_index` when the API does not return a sequence, so a drift on the priority can only be detected when the API reports it.

Example Usage

```hcl
resource "tencentcloud_fwm_vpc_acl_rule2" "example" {
  group_id          = "fwmrg-in3q0qkq0n"
  source_content    = "10.0.0.0/8"
  source_type       = "ip"
  dest_content      = "192.168.0.0/16"
  dest_type         = "ip"
  protocol          = "TCP"
  port              = "80,443"
  rule_action       = "accept"
  order_index       = 1
  edge_id           = "ALL"
  fw_group_id       = "ALL"
  belong_member_id  = "mem-xxxxxxxx"
  param_template_id = "ptb-xxxxxxxx"
  description       = "managed by terraform"
}
```

Import

FWM VPC ACL rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule2.example fwmrg-in3q0qkq0n#7520b898-30ba-415d-ab27-ff4542071ea9
```
