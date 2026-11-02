Provides a resource to manage an edge ACL rule of Cloud Firewall (FWM)

Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id          = "fwmrg-9z03vakjj8"
  direction         = 1
  source_content    = "0.0.0.0/0"
  source_type       = "ip"
  target_content    = "0.0.0.0/0"
  target_type       = "ip"
  protocol          = "TCP"
  port              = "80,443"
  rule_action       = "accept"
  scope             = "all"
  description       = "remark."
  belong_member_id  = "mem-xxxxxxxx"
  param_template_id = "pt-xxxxxxxx"
}
```

Import

FWM edge ACL rule can be imported using the `groupId#ruleId#direction`, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example fwmrg-9z03vakjj8#2fe0b492-5623-4646-b40e-7b0908307ca8#1
```

~> **NOTE:** The rule group ID and the direction are part of the resource ID, because `DescribeEdgeAclRules` can only query rules inside a given rule group. Changing `group_id` or `direction` recreates the rule, and the imported resource must be destroyed with a state that still holds both of them.

~> **NOTE:** `source_type` and `target_type` are applied in place through `ModifyEdgeAclRule`. When one of them is changed alone while `belong_member_id` is configured, the rule is recreated instead, because the configured member account may no longer match the new address type.

~> **NOTE:** `CreateEdgeAclRule` only returns a list of rule IDs, and `DescribeEdgeAclRules` returns the execution order as `sequence`. The create-side `OrderIndex` is not exposed as a configuration argument: a managed rule is always created with `OrderIndex` 1 and updated with the current `sequence`, so the rule priority is never silently reset.
