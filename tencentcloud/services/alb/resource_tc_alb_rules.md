Provides a resource to manage a set of Application Load Balancer forwarding rules.

Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  listener_id      = "lst-12345678"
  load_balancer_id = "alb-12345678"

  rules = [
    {
      actions = [
        {
          order = 1
          type  = "TargetGroup"
          target_group_config = {
            target_groups = [
              {
                target_group_id = "lbtg-12345678"
                weight          = 10
              }
            ]
            target_group_sticky_session = {
              enabled = false
              timeout = 1000
            }
          }
        }
      ]
      conditions = [
        {
          type        = "Host"
          host_config = ["example.com"]
        }
      ]
      priority = 1
      rule_name = "tf-test-rule"
    }
  ]
}
```

Import

Alb rules can be imported using the id, e.g.

```
terraform import tencentcloud_alb_rules.example rule-12345678
```
