Provides a resource to create a FWM (Firewall Manager) rule group

Example Usage

```hcl
resource "tencentcloud_fwm_rule_group" "example" {
  group_name = "tf-example-rule-group"
  product    = "cfw_vpc_acl"
}
```

Import

FWM rule group can be imported using the GroupId, e.g.

```
terraform import tencentcloud_fwm_rule_group.example fwmrg-nyalbmbi
```
