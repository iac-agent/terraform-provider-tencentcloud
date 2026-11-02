---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_edge_acl_rule_group"
sidebar_current: "docs-tencentcloud-resource-fwm_edge_acl_rule_group"
description: |-
  Provides a resource to create an internet border rule group of Firewall Manager (FWM)
---

# tencentcloud_fwm_edge_acl_rule_group

Provides a resource to create an internet border rule group of Firewall Manager (FWM)

The rule group is created together with the rules listed in it, and destroying the
resource deletes that rule group only: no other rule group of the account is touched,
because the underlying batch delete API is always called with the single id managed by
this resource. Only `group_name` can be updated in place; changing `rules` or any rule
field replaces the rule group.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `group_name` - (Optional, String) Name of the rule group, 1 to 50 characters. It can be updated in place.
* `product` - (Optional, String) Product the rule group belongs to. The internet border rule group is fixed to `cfw_edge_acl`; this attribute mirrors the value returned by the cloud.
* `rules` - (Optional, List) Rules of the rule group. Their order is significant: order_index starts at 1 and increases strictly, the smaller the number the higher the priority. Any change replaces the rule group, because the update API cannot modify rules.

The `rules` object supports the following:

* `belong_member_id` - (Optional, String) ID of the member account the rule belongs to, used in multi account scenarios. Required when source_type or target_type is instance or tag.
* `description` - (Optional, String) Description of the rule, 0 to 256 characters.
* `direction` - (Optional, Int) Direction of the rule: 0 egress, 1 ingress.
* `order_index` - (Optional, Int) Execution order of the rule, starting at 1 and strictly increasing. The smaller the number the higher the priority.
* `param_template_id` - (Optional, String) ID of the parameter template used by the rule.
* `port` - (Optional, String) Port of the rule. Supports a single port, a port range and comma separated ports, such as 80, 1-65535 or 80,443,8080.
* `protocol` - (Optional, String) Protocol of the rule: TCP, UDP, ICMP, ANY.
* `rule_action` - (Optional, String) Action of the rule: accept to allow, drop to block, log to observe.
* `scope` - (Optional, String) Effective scope of the rule: serial, side or all.
* `source_content` - (Optional, String) Source content. Format depends on source_type: IP/CIDR for ip, a domain for domain, a template id for template, a comma separated instance id list for instance, and key:value for tag.
* `source_type` - (Optional, String) Source address type: ip, domain, template, instance, tag.
* `target_content` - (Optional, String) Target content. Format is the same as source_content.
* `target_type` - (Optional, String) Target address type: ip, domain, template, instance, tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_id` - ID of the rule group assigned by the cloud.
* `id` - ID of the rule group, equal to group_id.
* `rule_count` - Number of rules in the rule group.

The `rules` object exports the following:

* `rule_id` - ID of the rule assigned by the cloud.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateEdgeAclRuleGroup](https://www.tencentcloud.com/search/?wd=CreateEdgeAclRuleGroup&type=document&lang=en)
  - [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Read:
  - [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Update:
  - [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
  - [ModifyRuleGroup](https://www.tencentcloud.com/search/?wd=ModifyRuleGroup&type=document&lang=en)
* Delete: [DeleteRuleGroup](https://www.tencentcloud.com/search/?wd=DeleteRuleGroup&type=document&lang=en)

## Import

FWM internet border rule group can be imported using the group id, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule_group.example fwmrg-9z03vakjj8
```

