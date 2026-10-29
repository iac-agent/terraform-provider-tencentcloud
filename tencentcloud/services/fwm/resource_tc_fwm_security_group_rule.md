Provides a resource to create a FWM (Firewall Manager) enterprise security group rule

Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id = "cfwg-nyalbmbi"

  rules {
    ip_version         = "ipv4"
    source_content     = "10.0.0.0/16"
    source_type        = "net"
    dest_content       = "10.0.1.0/24"
    dest_type          = "net"
    protocol           = "TCP"
    port               = "80"
    rule_action        = "2"
    description        = "allow http"
    order_index        = 1
    scope              = "1"
    protocol_port_type = 0
  }
}
```

Import

FWM security group rule can be imported using the GroupId and the RuleId separated by `#`, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example cfwg-nyalbmbi#rule-nyalbmbi
```
