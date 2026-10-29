Provides a resource to create a FWM (Firewall Manager) strategy.

A FWM strategy describes the rule groups that are executed before and after the firewall, and the accounts or account groups the rules are delivered to. A single create request may generate more than one strategy, which is reflected in the computed `strategy` attribute.

Example Usage

```hcl
resource "tencentcloud_fwm_account_group" "example" {
  group_name  = "tf-example"
  description = "remark."
}

resource "tencentcloud_fwm_strategy_v2" "example" {
  product         = "cfw"
  receive_account = ["100000000001"]
  receive_group   = [tencentcloud_fwm_account_group.example.group_id]

  pre_strategy {
    group_id = "grp-xxxxxxxx"
    sequence = 1
  }

  post_strategy {
    group_id = "grp-yyyyyyyy"
    sequence = 2
  }
}
```

Import

FWM strategy can be imported using the strategy id, e.g.

```
terraform import tencentcloud_fwm_strategy_v2.example sg-xxxxxxxx
```

-> **Note:** When a single create request generated more than one strategy, the resource id is a list of strategy ids joined by `#`. Such a resource can be imported with the full joined id, and it can only be destroyed or recreated afterwards, because `ModifyStrategy` operates on a single strategy.
