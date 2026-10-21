# tencentcloud_fwm_strategy

Provides a resource to manage a FWM (Firewall Management) strategy.

Example Usage

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

Import

FWM strategy can be imported using the strategy id, e.g.

```
terraform import tencentcloud_fwm_strategy.example strategy-xxxx
```
