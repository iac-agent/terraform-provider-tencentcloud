---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_vpc_acl_rule2"
sidebar_current: "docs-tencentcloud-resource-fwm_vpc_acl_rule2"
description: |-
  Provides a resource to create a VPC border firewall ACL rule of Cloud Firewall (FWM)
---

# tencentcloud_fwm_vpc_acl_rule2

Provides a resource to create a VPC border firewall ACL rule of Cloud Firewall (FWM)

~> **NOTE:** One resource manages exactly one ACL rule inside the rule group given by `group_id`. The cloud API accepts a list of rules when creating and deleting, but this resource always sends a single rule and requires exactly one rule ID in the response. The rule group itself is **not** managed by this resource; destroying this resource only removes the referenced rule and leaves the rule group and all other rules in it untouched.

~> **NOTE:** `order_index` sets the priority of the rule inside the rule group. The cloud side may re-sequence the other rules of the group as a result. The effective priority is exposed by the computed `sequence` attribute, which is also used to read `order_index` back.

## Example Usage

```hcl
resource "tencentcloud_fwm_vpc_acl_rule2" "example" {
  group_id          = "fwmrg-9z03vakjj8"
  source_content    = "10.0.0.0/16"
  source_type       = "ip"
  dest_content      = "0.0.0.0/0"
  dest_type         = "ip"
  protocol          = "TCP"
  port              = "80,443"
  rule_action       = "accept"
  order_index       = 1
  edge_id           = "ALL"
  fw_group_id       = "ALL"
  description       = "allow web traffic"
  param_template_id = "pt-xxxx"
  belong_member_id  = "mem-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `dest_content` - (Required, String) Destination address content. Its format depends on dest_type and follows the same rules as source_content.
* `dest_type` - (Required, String) Destination address type: ip, url, template, instance, tag.
* `edge_id` - (Required, String) Border firewall ID: ALL means global, CFWS-xxx means a specific border.
* `fw_group_id` - (Required, String) Firewall instance ID that the rule applies to: ALL means globally effective, ccn-xxx a Cloud Connect Network instance, cfwg-xxx a firewall group instance.
* `group_id` - (Required, String) ID of the ACL rule group that owns the rule, such as `fwmrg-9z03vakjj8`. The rule group is the parent container and must already exist before creating the rule. Changing this forces a new resource.
* `order_index` - (Required, Int) Priority of the rule, starting from 1. The rule is placed at this position inside the rule group and the other rules of the group may be re-sequenced by the cloud side as a result. The effective sequence is exposed by the computed `sequence` attribute.
* `protocol` - (Required, String) Protocol, such as TCP, UDP, ICMP, ANY, HTTP, HTTPS, DNS, FTP.
* `rule_action` - (Required, String) Action of the rule: accept, drop, log.
* `source_content` - (Required, String) Source address content. Its format depends on source_type: CIDR for ip, domain name for url, template ID for template, instance ID list for instance, tag for tag.
* `source_type` - (Required, String) Source address type: ip, url, template, instance, tag.
* `belong_member_id` - (Optional, String) ID of the member account that the rule belongs to. Required when fw_group_id is a cfwg-xxx firewall group or when source_type/dest_type is instance or tag.
* `description` - (Optional, String) Description of the rule.
* `param_template_id` - (Optional, String) ID of the port template used by the rule.
* `port` - (Optional, String) Port or port range. Empty for the ICMP protocol.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `app_id` - AppId of the member account that the rule belongs to.
* `create_time` - Creation time of the rule.
* `dest_name` - Name of the destination asset, when the destination is an instance or a tag.
* `dns_parse_cnt` - Number of domains resolved by the rule.
* `fw_group_name` - Name of the firewall group that the rule applies to.
* `hit_count` - Accumulated hit count of the rule over all member accounts it was delivered to, as returned by the cloud side. Zero when the rule has never been hit.
* `id` - Resource identifier in the form `<group_id>#<rule_id>`.
* `invalid` - Whether the rule is invalid: 0 means valid, 1 means invalid.
* `ip_version` - IP version of the rule: ipv4 or ipv6.
* `last_hit_time` - Latest hit time of the rule over all member accounts it was delivered to. "-" when the rule has never been hit.
* `nickname` - Nickname of the member account that the rule belongs to.
* `node_name` - Department that the member account belongs to.
* `param_template_name` - Name of the port template used by the rule.
* `rule_id` - ID of the ACL rule, returned by the API on create and used to locate the rule afterwards.
* `sequence` - Effective priority of the rule inside the rule group, as returned by the cloud side.
* `source_name` - Name of the source asset, when the source is an instance or a tag.
* `uin` - Uin of the member account that the rule belongs to.
* `update_time` - Last modification time of the rule.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateVpcAclRule](https://www.tencentcloud.com/search/?wd=CreateVpcAclRule&type=document&lang=en)
  - [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Read: [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Update:
  - [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
  - [ModifyVpcAclRule](https://www.tencentcloud.com/search/?wd=ModifyVpcAclRule&type=document&lang=en)
* Delete: [DeleteVpcAclRule](https://www.tencentcloud.com/search/?wd=DeleteVpcAclRule&type=document&lang=en)

## Import

FWM VPC ACL rule can be imported using the `groupId#ruleId`, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule2.example fwmrg-9z03vakjj8#cfwrule-2fe0b492
```

