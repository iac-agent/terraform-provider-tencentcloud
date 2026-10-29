---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_rule_group"
sidebar_current: "docs-tencentcloud-resource-fwm_rule_group"
description: |-
  Provides a resource to create a FWM (Firewall Manager) rule group
---

# tencentcloud_fwm_rule_group

Provides a resource to create a FWM (Firewall Manager) rule group

## Example Usage

```hcl
resource "tencentcloud_fwm_rule_group" "example" {
  group_name = "tf-example-rule-group"
  product    = "cfw_vpc_acl"
}
```

## Argument Reference

The following arguments are supported:

* `by` - (Optional, String) Field used to sort the query result.
* `filters` - (Optional, List) Filters used to query the rule group list.
* `group_id` - (Optional, String) ID of the rule group used by the modify operation. It must be the ID of the rule group managed by this resource.
* `group_ids` - (Optional, List) IDs of the rule groups used by the delete operation. It must contain the ID of the rule group managed by this resource.
* `group_name` - (Optional, String) Name of the rule group. It cannot be empty and can be up to 50 characters long (Chinese characters count as one character).
* `order` - (Optional, String) Sorting order of the query result, e.g. `desc` or `asc`.
* `product` - (Optional, String) Product the rule group belongs to. Valid values: `enterprise_sg` (enterprise security group), `cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall). It is required and does not accept `all`.

The `filters` object supports the following:

* `name` - (Optional, String) Name of the filter field.
* `operator_type` - (Optional, Int) Operator type. Valid values: `1` (exact match), `7` (in), `9` (fuzzy match).
* `values` - (Optional, List) Values of the filter field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `app_id` - Member AppId.
      * `member_id` - Member ID, empty in personal scenarios.
      * `nickname` - Member nickname.
      * `node_name` - Department the member belongs to.
      * `uin` - Member Uin.
    * `members` - Member list, returned only when `receiver_type` is `1`.
    * `nickname` - Tenant name.
    * `receiver_type` - Receiver type. `0` means account Uin, `1` means account group.
    * `uin` - Tenant Uin.
  * `app_id` - Tenant AppId.
  * `create_by` - SubUin of the creator.
  * `create_name` - Name of the creator.
  * `create_time` - Creation time.
  * `direction` - Rule direction. `1` means inbound, `0` means outbound, `-1` means no direction.
  * `exec_areas` - Execution stages the rule group belongs to.
  * `group_id` - Rule group ID.
  * `group_name` - Rule group name.
  * `product` - Product type.
  * `receive_account` - Associated accounts the rule group is delivered to.
  * `rule_count` - Number of rules in the rule group.
  * `update_by` - SubUin of the last updater.
  * `update_name` - Name of the last updater.
  * `update_time` - Update time.
* `id` - ID of the FWM rule group resource.
* `rule_groups` - List of rule groups.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRuleGroup](https://www.tencentcloud.com/search/?wd=CreateRuleGroup&type=document&lang=en)
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Read: [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Update:
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
  - [ModifyRuleGroup](https://www.tencentcloud.com/search/?wd=ModifyRuleGroup&type=document&lang=en)
* Delete: [DeleteRuleGroup](https://www.tencentcloud.com/search/?wd=DeleteRuleGroup&type=document&lang=en)

## Import

FWM rule group can be imported using the GroupId, e.g.

```
terraform import tencentcloud_fwm_rule_group.example fwmrg-nyalbmbi
```

