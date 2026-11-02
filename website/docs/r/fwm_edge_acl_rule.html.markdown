---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_edge_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_edge_acl_rule"
description: |-
  Provides a resource to manage an edge ACL rule of Cloud Firewall (FWM)
---

# tencentcloud_fwm_edge_acl_rule

Provides a resource to manage an edge ACL rule of Cloud Firewall (FWM)

## Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id          = "fwmrg-9z03vakjj8"
  direction         = 1
  source_content    = "0.0.0.0/0"
  source_type       = "ip"
  target_content    = "0.0.0.0/0"
  target_type       = "ip"
  protocol          = "TCP"
  port              = "80,443"
  rule_action       = "accept"
  scope             = "all"
  description       = "remark."
  belong_member_id  = "mem-xxxxxxxx"
  param_template_id = "pt-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int) Rule direction: 0-egress, 1-ingress. The direction is part of the resource id (`group_id#rule_id#direction`), so changing it recreates the rule.
* `group_id` - (Required, String) Edge ACL rule group ID. The rule group is the parent container and must already exist before creating the rule. The rule is always queried within this group, so changing it recreates the rule.
* `belong_member_id` - (Optional, String) Member account ID that the rule belongs to (used in multi-account scenarios). Required when source_type or target_type is instance or tag.
* `description` - (Optional, String) Rule description, length 0-256 characters.
* `param_template_id` - (Optional, String) Parameter template ID.
* `port` - (Optional, String) Port. Supports single port, port range and comma separated multi ports, such as 80, 1-65535, 80,443,8080.
* `protocol` - (Optional, String) Protocol type: TCP, UDP, ICMP, ANY.
* `rule_action` - (Optional, String) Rule action: accept-pass, drop-block, log-observe.
* `scope` - (Optional, String) Effective scope: serial, side, all.
* `source_content` - (Optional, String) Source address content. Format depends on source_type: IP/CIDR for ip, domain for domain, template ID for template, instance ID list (comma separated) for instance, tag key:value for tag.
* `source_type` - (Optional, String) Source address type: ip, domain, template, instance, tag. It can be modified in place; it only recreates the rule when it is changed alone while belong_member_id is configured, because the member account may no longer match the new address type.
* `target_content` - (Optional, String) Target address content. Format is the same as source_content.
* `target_type` - (Optional, String) Target address type: ip, domain, template, instance, tag. It can be modified in place; it only recreates the rule when it is changed alone while belong_member_id is configured, because the member account may no longer match the new address type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `belong_member` - Member account that the rule belongs to.
  * `app_id` - Member AppId.
  * `member_id` - Member ID, empty in the personal account scenario.
  * `nickname` - Member nickname.
  * `node_name` - Department the member belongs to.
  * `uin` - Member Uin.
* `city_code` - City ID, returned for region rules.
* `city_key` - Province/city short name, returned for region rules.
* `city_name` - Region name level 2, returned for region rules (used by regular matching).
* `cloud_code` - Cloud vendor code, returned for cloud vendor rules.
* `country_code` - Country ID, returned for region rules.
* `country_key` - Country short name, returned for region rules.
* `country_name` - Region name level 1, returned for region rules (used by regular matching).
* `create_time` - Rule creation time.
* `dns_parse_cnt` - Number of parsed domains.
* `hit_count` - Accumulated hit count of the rule over all delivered member accounts. It is the Cloud Firewall value without any time window conversion, and 0 when the rule never hit.
* `instance_name` - Instance name returned by the server.
* `invalid` - Whether the rule is invalid: 0-valid, 1-invalid.
* `is_cloud` - 0 means a normal rule, 1 means a cloud vendor rule.
* `is_region` - 0 means a normal rule, 1 means a region rule.
* `last_hit_time` - Latest hit time over all delivered member accounts, "-" when the rule never hit.
* `param_template_name` - Parameter template name returned by the server, only set when the rule uses a parameter template.
* `rule_id` - Rule ID, it is part of the resource id `group_id#rule_id#direction`.
* `sequence` - Rule execution sequence returned by the server (the read side alias of the create `OrderIndex`).
* `source_name` - Source address name returned by the server, only set when source_type is template, instance or tag.
* `target_name` - Target address name returned by the server, only set when target_type is template, instance or tag.
* `update_time` - Rule last update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEdgeAclRule](https://www.tencentcloud.com/search/?wd=CreateEdgeAclRule&type=document&lang=en)
* Read: [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
* Update: [ModifyEdgeAclRule](https://www.tencentcloud.com/search/?wd=ModifyEdgeAclRule&type=document&lang=en)
* Delete: [DeleteEdgeAclRule](https://www.tencentcloud.com/search/?wd=DeleteEdgeAclRule&type=document&lang=en)

## Import

FWM edge ACL rule can be imported using the `groupId#ruleId#direction`, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example fwmrg-9z03vakjj8#2fe0b492-5623-4646-b40e-7b0908307ca8#1
```

~> **NOTE:** The rule group ID and the direction are part of the resource ID, because `DescribeEdgeAclRules` can only query rules inside a given rule group. Changing `group_id` or `direction` recreates the rule, and the imported resource must be destroyed with a state that still holds both of them.

~> **NOTE:** `source_type` and `target_type` are applied in place through `ModifyEdgeAclRule`. When one of them is changed alone while `belong_member_id` is configured, the rule is recreated instead, because the configured member account may no longer match the new address type.

~> **NOTE:** `CreateEdgeAclRule` only returns a list of rule IDs, and `DescribeEdgeAclRules` returns the execution order as `sequence`. The create-side `OrderIndex` is not exposed as a configuration argument: a managed rule is always created with `OrderIndex` 1 and updated with the current `sequence`, so the rule priority is never silently reset.

