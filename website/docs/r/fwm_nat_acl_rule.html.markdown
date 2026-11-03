---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_nat_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_nat_acl_rule"
description: |-
  Provides a resource to create a FWM NAT boundary (NAT ACL) rule
---

# tencentcloud_fwm_nat_acl_rule

Provides a resource to create a FWM NAT boundary (NAT ACL) rule

## Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "fwmrg-in3q0qkq0n"
  direction      = 0
  source_content = "0.0.0.0/0"
  source_type    = "ip"
  target_content = "0.0.0.0/0"
  target_type    = "ip"
  protocol       = "ANY"
  rule_action    = "accept"
  order_index    = 1
  scope          = "ALL"
  port           = "-1/-1"
  description    = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int, ForceNew) Rule direction: 0-outbound, 1-inbound. It is also part of the resource id, so changing it will recreate the rule.
* `group_id` - (Required, String, ForceNew) NAT ACL rule group ID. The group is the parent container and must already exist. Changing it will recreate the rule.
* `order_index` - (Required, Int) Rule priority, starting from 1.
* `protocol` - (Required, String) Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP...
* `rule_action` - (Required, String) Action: accept/drop/log.
* `scope` - (Required, String) Effective range: ALL-global, ap-xxx-region, cfwnat-xxx-NAT firewall instance.
* `source_content` - (Required, String) Source address content.
* `source_type` - (Required, String) Source type: ip/url/template/instance/tag.
* `target_content` - (Required, String) Target address content.
* `target_type` - (Required, String) Target type: ip/url/template/instance/tag.
* `belong_member_id` - (Optional, String) Member account ID the rule belongs to. Required when scope is cfwnat-xxx or source_type/target_type is instance/tag.
* `description` - (Optional, String) Rule description.
* `param_template_id` - (Optional, String) Port template ID.
* `port` - (Optional, String) Port, e.g. 80, 80,443, 80-90; leave empty when protocol is ICMP.
* `rule_id` - (Optional, String, ForceNew) Rule ID. It is filled by the server after creation and is part of the resource id `group_id#rule_id#direction`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `belong_member` - Member account the rule belongs to.
  * `app_id` - Member app ID.
  * `member_id` - Member ID (empty in personal scenario).
  * `nickname` - Member nickname.
  * `node_name` - Department the member belongs to.
  * `uin` - Member UIN.
* `city_code` - City code.
* `city_key` - City key.
* `city_name` - City name.
* `cloud_code` - Cloud vendor code.
* `country_code` - Country code.
* `country_key` - Country key.
* `country_name` - Country name.
* `create_time` - Creation time.
* `dns_parse_cnt` - Number of domains.
* `fw_ins_id` - Firewall instance ID.
* `hit_count` - Accumulated hit count of the rule under all delivered member accounts. It is 0 when never hit.
* `invalid` - Whether the rule is invalid: 0-valid, 1-invalid.
* `is_cloud` - Whether the rule is a cloud vendor rule: 0-no, 1-yes.
* `is_region` - Whether the rule is a region rule: 0-no, 1-yes.
* `last_hit_time` - Latest hit time of the rule under all delivered member accounts. It is "-" when never hit.
* `param_template_name` - Port template name.
* `scope_desc` - Effective range description.
* `sequence` - Rule priority returned by the server, the readable counterpart of `order_index`.
* `source_name` - Source asset name.
* `target_name` - Target asset name.
* `update_time` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNatAclRule](https://www.tencentcloud.com/search/?wd=CreateNatAclRule&type=document&lang=en)
* Read: [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
* Update: [ModifyNatAclRule](https://www.tencentcloud.com/search/?wd=ModifyNatAclRule&type=document&lang=en)
* Delete: [DeleteNatAclRule](https://www.tencentcloud.com/search/?wd=DeleteNatAclRule&type=document&lang=en)

## Import

FWM NAT boundary rule can be imported using the `group_id#rule_id#direction` composite id, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-in3q0qkq0n#7520b898-30ba-415d-ab27-ff4542071ea9#0
```

~> **NOTE:** `group_id` has no dedicated field in the query response, so it is restored from the import id. The
configured `order_index` is only verified through the rule priority (`sequence`) returned by the cloud API.

