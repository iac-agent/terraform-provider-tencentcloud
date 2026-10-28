---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_vpc_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_vpc_acl_rule"
description: |-
  Provides a resource to create a FWM VPC border ACL rule
---

# tencentcloud_fwm_vpc_acl_rule

Provides a resource to create a FWM VPC border ACL rule

## Example Usage

```hcl
resource "tencentcloud_fwm_vpc_acl_rule" "example" {
  group_id    = "acg-nyalbmbi"
  fw_group_id = "cfwg-xxxxxxxx"
  rule_action = "accept"
  description = "tf-example"
  order_index = 1
}
```

## Argument Reference

The following arguments are supported:

* `fw_group_id` - (Required, String) Firewall instance ID that the rule applies to. `ALL` means global, `ccn-xxx` is a Cloud Connect Network instance and `cfwg-xxx` is a firewall group instance. Changing it forces a new rule.
* `group_id` - (Required, String) ID of the VPC border ACL rule group the rule belongs to. Changing it forces a new rule.
* `rule_action` - (Required, String) Action of the rule. Valid values: `accept`, `drop`, `log`. Changing it forces a new rule.
* `description` - (Optional, String) Description of the rule.
* `order_index` - (Optional, Int) Priority of the rule, starting from `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `app_id` - AppId of the member.
    * `member_id` - ID of the member. It is empty for personal accounts.
    * `nickname` - Nickname of the member.
    * `node_name` - Department the member belongs to.
    * `uin` - Uin of the member.
  * `belong_member_id` - ID of the member account the rule belongs to.
  * `belong_member` - Member account the rule belongs to.
  * `create_time` - Creation time of the rule.
  * `description` - Description of the rule.
  * `dest_content` - Destination address content of the rule.
  * `dest_name` - Name of the destination asset.
  * `dest_type` - Destination type of the rule. Valid values: `ip`, `url`, `template`, `instance`, `tag`.
  * `dns_parse_cnt` - Number of parsed domain names of the rule.
  * `edge_id` - ID of the border firewall. `ALL` means global and `CFWS-xxx` means the specified border.
  * `fw_group_id` - ID of the firewall instance.
  * `fw_group_name` - Name of the firewall group.
  * `hit_count` - Accumulated hit count of the rule across all delivered member accounts.
  * `invalid` - Whether the rule is invalid. `0` means valid, `1` means invalid.
  * `ip_version` - IP version of the rule. Valid values: `ipv4`, `ipv6`.
  * `last_hit_time` - Last hit time of the rule. It is `-` when the rule has never been hit.
  * `order_index` - Priority of the rule, starting from `1`.
  * `param_template_id` - ID of the port template.
  * `param_template_name` - Name of the port template.
  * `port` - Port of the rule. It is empty when the protocol is `ICMP`.
  * `protocol` - Protocol of the rule, e.g. `TCP`, `UDP`, `ICMP`, `ANY`, `HTTP`, `HTTPS`, `DNS`, `FTP`.
  * `rule_action` - Action of the rule. Valid values: `accept`, `drop`, `log`.
  * `rule_id` - ID of the rule.
  * `sequence` - Priority of the rule.
  * `source_content` - Source address content of the rule.
  * `source_name` - Name of the source asset.
  * `source_type` - Source type of the rule. Valid values: `ip`, `url`, `template`, `instance`, `tag`.
  * `update_time` - Last modification time of the rule.
* `all_total_count` - Total number of rules of the rule group without any filter.
* `id` - Resource identifier in the form `<group_id>#<rule_id>`.
* `rule_id` - ID of the VPC border ACL rule managed by this resource.
* `rule_ids` - IDs of the VPC border ACL rules created by this resource.
* `rules` - Full list of rules of the rule group as returned by the cloud API, including the rule managed by this resource.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateVpcAclRule](https://www.tencentcloud.com/search/?wd=CreateVpcAclRule&type=document&lang=en)
  - [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Read: [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Update: [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Delete: [DeleteVpcAclRule](https://www.tencentcloud.com/search/?wd=DeleteVpcAclRule&type=document&lang=en)

## Import

FWM VPC border ACL rule can be imported using the composite id `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule.example acg-nyalbmbi#rule-xxxxxxxx
```

