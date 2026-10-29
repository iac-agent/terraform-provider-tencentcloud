---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_security_group_rule"
description: |-
  Provides a resource to manage FWM enterprise security group rules.
---

# tencentcloud_fwm_security_group_rule

Provides a resource to manage FWM enterprise security group rules.

Each apply creates one rule per element of `rules`, and `terraform destroy` deletes every rule tracked by the resource.

~> **NOTE:** The resource ID is a compound ID in the form `<group_id>#<rule_id>#<rule_id>...`, which is also the format required by `terraform import`.

## Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id = "fwg-4x1y2z3a"

  rules {
    ip_version     = "v4"
    source_type    = "net"
    source_content = "10.0.0.0/8"
    dest_type      = "net"
    dest_content   = "192.168.0.0/16"
    protocol       = "TCP"
    port           = "80"
    rule_action    = "2"
    description    = "allow http from intranet"
    scope          = "1"
  }

  rules {
    ip_version     = "v4"
    source_type    = "net"
    source_content = "0.0.0.0/0"
    dest_type      = "net"
    dest_content   = "10.0.1.0/24"
    protocol       = "TCP"
    port           = "22"
    rule_action    = "1"
    description    = "block ssh from internet"
    scope          = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, String) ID of the enterprise security group rule group.
* `rules` - (Optional, List) List of security group rules to create.

The `rules` object supports the following:

* `belong_member_id` - (Optional, String) ID of the member account the rule belongs to. It is required when group_id is a cfwg-xxx group, or when source_type/dest_type is `instance` or `tag`.
* `description` - (Optional, String) Description of the rule.
* `dest_content` - (Optional, String) Destination address of the rule.
* `dest_type` - (Optional, String) Destination type of the rule.
* `ip_version` - (Optional, String) IP version of the rule.
* `order_index` - (Optional, Int) Priority of the rule.
* `port` - (Optional, String) Port range of the rule.
* `protocol_port_type` - (Optional, Int) Protocol port type of the rule.
* `protocol` - (Optional, String) Protocol of the rule.
* `rule_action` - (Optional, String) Action of the rule, for example `1` for block and `2` for allow.
* `rule_id` - (Optional, String) ID of the rule.
* `scope` - (Optional, String) Effective scope of the rule.
* `service_template_id` - (Optional, String) Protocol port template ID of the rule.
* `source_content` - (Optional, String) Source address of the rule.
* `source_type` - (Optional, String) Source type of the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `detail` - Description of the rule.
  * `dst_content` - Destination content of the rule.
  * `dst_port` - Destination port of the rule. Single port, multiple ports and port ranges are supported.
  * `dst_type` - Destination type of the rule.
  * `ip_version` - IP version of the rule.
  * `protocol_port_type` - Protocol port template parameter of the rule.
  * `protocol` - Protocol of the rule.
  * `region` - Region of the rule.
  * `rule_action` - Action of the rule, `1` for block and `2` for allow.
  * `rule_id` - ID of the rule.
  * `rule_source` - Source of the rule, `0` means it was added from the user console.
  * `scope` - Effective scope of the rule, bitwise, `1` for SG (enterprise security group) and `2` for LH (lightweight server).
  * `sequence` - Priority of the rule.
  * `service_template_id` - Protocol port template ID of the rule.
  * `src_content` - Source content of the rule.
  * `src_type` - Source type of the rule.
* `id` - Resource identifier in the form `<group_id>#<rule_id>#<rule_id>...`.
* `rule` - Rule details returned by the cloud API.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateSecurityGroupRule](https://www.tencentcloud.com/search/?wd=CreateSecurityGroupRule&type=document&lang=en)
  - [DescribeSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRule&type=document&lang=en)
* Read: [DescribeSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRule&type=document&lang=en)
* Update:
  - [DescribeSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRule&type=document&lang=en)
  - [ModifySecurityGroupRule](https://www.tencentcloud.com/search/?wd=ModifySecurityGroupRule&type=document&lang=en)
* Delete: [DeleteSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DeleteSecurityGroupRule&type=document&lang=en)

## Import

`tencentcloud_fwm_security_group_rule` can be imported using the compound ID `<group_id>#<rule_id>#<rule_id>...`, for example:

```
$ terraform import tencentcloud_fwm_security_group_rule.example fwg-4x1y2z3a#fwr-a1b2c3d4#fwr-e5f6g7h8
```

