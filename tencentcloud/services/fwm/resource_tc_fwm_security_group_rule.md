Provides a resource to manage FWM enterprise security group rules.

Each apply creates one rule per element of `rules`, and `terraform destroy` deletes every rule tracked by the resource.

~> **NOTE:** The resource ID is a compound ID in the form `<group_id>#<rule_id>#<rule_id>...`, which is also the format required by `terraform import`.

Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id = "fwg-4x1y2z3a"

  rules {
    ip_version   = "v4"
    source_type  = "net"
    source_content = "10.0.0.0/8"
    dest_type    = "net"
    dest_content = "192.168.0.0/16"
    protocol     = "TCP"
    port         = "80"
    rule_action  = "2"
    description  = "allow http from intranet"
    scope        = "1"
  }

  rules {
    ip_version   = "v4"
    source_type  = "net"
    source_content = "0.0.0.0/0"
    dest_type    = "net"
    dest_content = "10.0.1.0/24"
    protocol     = "TCP"
    port         = "22"
    rule_action  = "1"
    description  = "block ssh from internet"
    scope        = "1"
  }
}
```

Import

`tencentcloud_fwm_security_group_rule` can be imported using the compound ID `<group_id>#<rule_id>#<rule_id>...`, for example:

```
$ terraform import tencentcloud_fwm_security_group_rule.example fwg-4x1y2z3a#fwr-a1b2c3d4#fwr-e5f6g7h8
```
