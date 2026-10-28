Provides a resource to create a FWM VPC border ACL rule

Example Usage

```hcl
resource "tencentcloud_fwm_vpc_acl_rule" "example" {
  group_id    = "acg-nyalbmbi"
  fw_group_id = "cfwg-xxxxxxxx"
  rule_action = "accept"
  description = "tf-example"
  order_index = 1
}
```

Import

FWM VPC border ACL rule can be imported using the composite id `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule.example acg-nyalbmbi#rule-xxxxxxxx
```