---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_organ_members"
sidebar_current: "docs-tencentcloud-data_source-fwm_organ_members"
description: |-
  Provides a data source to query the member list of a FWM (Cloud Firewall Manager) organization.
---

# tencentcloud_fwm_organ_members

Provides a data source to query the member list of a FWM (Cloud Firewall Manager) organization.

## Example Usage

```hcl
data "tencentcloud_fwm_organ_members" "example" {
  filters {
    name          = "MemberId"
    values        = ["mem-3r5c8tzb47"]
    operator_type = 1
  }
  by    = "MemberCreateTime"
  order = "desc"
}

output "fwm_organ_member_count" {
  value = data.tencentcloud_fwm_organ_members.example.total_count
}

output "fwm_organ_member_ids" {
  value = [for m in data.tencentcloud_fwm_organ_members.example.members : m.member_id]
}
```

## Argument Reference

The following arguments are supported:

* `by` - (Optional, String) Field used to sort the returned members, for example `MemberCreateTime`.
* `filters` - (Optional, List) Search conditions used to filter the returned members, such as member ID, account name, role or management status. All conditions are evaluated together.
* `order` - (Optional, String) Sort direction: `asc` for ascending, `desc` for descending.

The `filters` object supports the following:

* `name` - (Optional, String) Name of the field to filter on, for example `MemberId`, `AppId` or `Nickname`.
* `operator_type` - (Optional, Int) Operator of this condition: 1 - exact match, 7 - in, 9 - fuzzy match.
* `values` - (Optional, Set) Values of the field to filter on.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `group_id` - Account group ID.
    * `group_name` - Account group name.
  * `account_group` - Account group the member belongs to.
  * `app_id` - AppId of the member account.
  * `cfw_capable` - Whether the member is capable of using cloud firewall.
  * `cfw_instance_id` - Cloud firewall billing instance ID. A non-empty value means cloud firewall has been purchased.
  * `cfw_managed` - Cloud firewall management status: 0 - not managed, 1 - managed.
  * `cfw_pay_status` - Cloud firewall billing status: 0 - not purchased, 2 - purchased, 3 - trial, 4 - expired.
  * `cfw_share_role_display` - Display name of the cloud firewall sharing role.
  * `cfw_share_role` - Cloud firewall sharing role: `sharer` - sharer, `user` - user, `none` - not set.
  * `cfw_sharer_app_id` - AppId of the cloud firewall sharer, set when the member acts as a user.
  * `is_cfw_post_pay` - Whether the member uses the pay-as-you-go cloud firewall edition: 1 - pay-as-you-go, 0 - not pay-as-you-go, -1 - unknown.
  * `join_type` - Way the account joined the organization.
  * `member_create_time` - Time when the member joined the organization.
  * `member_id` - Member ID.
  * `nickname` - Name of the member account.
  * `node_name` - Name of the organization node the member belongs to.
  * `policy_analysis_enabled` - Policy analysis permission: 0 - disabled, 1 - enabled.
  * `role_display` - Display name of the member role.
  * `role` - Role of the member: `admin` - administrator, `delegatedAdmin` - delegated administrator, `member` - ordinary member.
  * `sg_managed` - Private security group management switch: 0 - not managed, 1 - managed.
  * `sub_account_count` - Number of sub-accounts of the member.
  * `uin` - Uin of the member account.
* `members` - List of FWM organization members.
* `total_count` - Number of members that match the query conditions.


