---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_nat_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_nat_acl_rule"
description: |-
  Provides a resource to manage a FWM (Firewall Manager) NAT boundary ACL rule.
---

# tencentcloud_fwm_nat_acl_rule

Provides a resource to manage a FWM (Firewall Manager) NAT boundary ACL rule.

## Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "acg-1a2b3c4d"
  source_content = "10.0.0.0/8"
  source_type    = "ip"
  target_content = "192.168.1.0/24"
  target_type    = "ip"
  protocol       = "TCP"
  rule_action    = "accept"
  order_index    = 1
  scope          = "ALL"
  direction      = 0
  port           = "80,443"
  description    = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) ID of the NAT boundary ACL rule group the rule belongs to.
* `order_index` - (Required, Int) Priority of the rule, starting from 1.
* `protocol` - (Required, String) Protocol, such as `TCP`, `UDP`, `ICMP`, `ANY`, `HTTP`, `HTTPS`, `DNS`, `FTP`.
* `rule_action` - (Required, String) Action of the rule, valid values: `accept`, `drop`, `log`.
* `scope` - (Required, String) Effective scope of the rule: `ALL` takes effect globally, `ap-xxx` takes effect in a region, `cfwnat-xxx` takes effect on a NAT firewall instance.
* `source_content` - (Required, String) Source address content.
* `source_type` - (Required, String) Source type, valid values: `ip`, `url`, `template`, `instance`, `tag`.
* `target_content` - (Required, String) Target address content.
* `target_type` - (Required, String) Target type, valid values: `ip`, `url`, `template`, `instance`, `tag`.
* `belong_member_id` - (Optional, String) ID of the member account the rule belongs to. Required when `scope` is `cfwnat-xxx` or `source_type`/`target_type` is `instance`/`tag`.
* `description` - (Optional, String) Description of the rule.
* `direction` - (Optional, Int) Direction of the rule: `0` outbound, `1` inbound.
* `param_template_id` - (Optional, String) ID of the port template.
* `port` - (Optional, String) Port of the rule. It is empty when the protocol is `ICMP`.
* `rule_id` - (Optional, String) ID of the rule. It is returned by the create API and is required when the rule is modified; it is the identifier used by read, update and delete.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `app_id` - AppId of the member account.
  * `member_id` - ID of the member account. It is empty in the personal account scenario.
  * `nickname` - Nickname of the member account.
  * `node_name` - Organization node the member account belongs to.
  * `uin` - Uin of the member account.
* `all_total_count` - Total number of rules in the rule group without filters.
* `belong_member` - Member account the rule belongs to.
* `city_code` - Code of the city.
* `city_key` - Key of the city.
* `city_name` - Name of the city.
* `cloud_code` - Code of the cloud provider.
* `country_code` - Code of the country.
* `country_key` - Key of the country.
* `country_name` - Name of the country.
* `create_time` - Creation time of the rule.
* `dns_parse_cnt` - Number of parsed domains.
* `fw_ins_id` - ID of the firewall instance.
* `hit_count` - Cumulative hit count of the rule across all delivered member accounts.
* `id` - Resource identifier in the form `<group_id>#<rule_id>`.
* `invalid` - Whether the rule is invalid: `0` valid, `1` invalid.
* `is_cloud` - Whether the rule is a cloud provider rule: `0` no, `1` yes.
* `is_region` - Whether the rule is a region rule: `0` no, `1` yes.
* `last_hit_time` - Latest hit time of the rule across all delivered member accounts. It is `-` when the rule has never been hit.
* `param_template_name` - Name of the port template.
* `rule_ids` - Rule ID list returned by the create API. The cloud API creates one rule per call, so it holds a single element.
* `scope_desc` - Description of the effective scope.
* `sequence` - Priority of the rule as returned by the cloud API.
* `source_name` - Name of the source asset.
* `target_name` - Name of the target asset.
* `update_time` - Last update time of the rule.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateNatAclRule](https://www.tencentcloud.com/search/?wd=CreateNatAclRule&type=document&lang=en)
  - [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
* Read: [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
* Update:
  - [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
  - [ModifyNatAclRule](https://www.tencentcloud.com/search/?wd=ModifyNatAclRule&type=document&lang=en)
* Delete: [DeleteNatAclRule](https://www.tencentcloud.com/search/?wd=DeleteNatAclRule&type=document&lang=en)

## Import

FWM NAT ACL rule can be imported using the composite id `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example acg-1a2b3c4d#aclr-9x8y7z6w
```

