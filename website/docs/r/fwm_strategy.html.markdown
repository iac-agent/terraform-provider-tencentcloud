---
subcategory: "Firewall Management(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_strategy"
sidebar_current: "docs-tencentcloud-resource-fwm_strategy"
description: |-
  # tencentcloud_fwm_strategy
---

# tencentcloud_fwm_strategy

# tencentcloud_fwm_strategy

Provides a resource to manage a FWM (Firewall Management) strategy.

## Example Usage

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product         = "clb"
  receive_account = ["100000000001", "100000000002"]
  pre_strategy {
    group_id = "rule-group-xxxx"
    sequence = 1
  }
  post_strategy {
    group_id = "rule-group-yyyy"
    sequence = 2
  }
  receive_group = ["account-group-xxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) Product type. This field is only used when creating and cannot be changed afterwards.
* `group_id` - (Optional, String) Rule group ID used when updating the strategy.
* `post_strategy` - (Optional, Object) Post-zone rule group. This field is only used when creating and cannot be changed afterwards.
* `pre_strategy` - (Optional, Object) Pre-zone rule group. This field is only used when creating and cannot be changed afterwards.
* `receive_account` - (Optional, List) List of accounts (uins) receiving the dispatch rules.
* `receive_group` - (Optional, List) List of account groups receiving the dispatch rules.
* `sequence` - (Optional, Int) Priority of the rule group used when updating the strategy.

The `post_strategy` object supports the following:

* `group_id` - (Optional, String) Rule group ID.
* `sequence` - (Optional, Int) Priority of the rule group.

The `pre_strategy` object supports the following:

* `group_id` - (Optional, String) Rule group ID.
* `sequence` - (Optional, Int) Priority of the rule group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `app_id` - AppId of the member.
    * `member_id` - ID of the member. Empty in personal scenarios.
    * `nickname` - Nickname of the member.
    * `node_name` - Department the member belongs to.
    * `uin` - Uin of the member.
  * `members` - Members of the account group. Only returned when `receiver_type` is `1`.
  * `nickname` - Name of the tenant.
  * `receiver_type` - Type of the receiver. `0` means account uin and `1` means account group.
  * `uin` - Uin of the tenant.
* `create_by` - Creator of the strategy.
* `create_name` - Name of the creator.
* `create_time` - Creation time of the strategy.
* `err_msg` - Reason why the dispatch failed.
* `error_type` - Type of the dispatch failure reason.
* `exec_area` - Execution area of the strategy.
* `group_name` - Name of the rule group.
* `id` - ID of the strategy, same as `strategy_id`.
* `receive_account_detail` - Details of the accounts receiving the dispatch rules.
* `rule_count` - Number of rules in the strategy.
* `rule_status` - Status of the strategy.
* `strategy_id` - ID of the strategy.
* `update_by` - Last updater of the strategy.
* `update_name` - Name of the last updater.
* `update_time` - Last updated time of the strategy.


## Import

FWM strategy can be imported using the strategy id, e.g.

```
terraform import tencentcloud_fwm_strategy.example strategy-xxxx
```

