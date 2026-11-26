Provides a resource to manage a listener forwarding rule.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  listener_id      = "lst-12345678"
  load_balancer_id = "alb-12345678"
  dry_run          = false

  rules = [
    {
      priority  = 10
      rule_name = "tf-test-rule"

      conditions = [
        {
          type = "Host"
          host_config = ["example.com"]
        }
      ]

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
              enabled = true
              timeout = 1000
            }
          }
        }
      ]
    }
  ]
}
```

Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example rule-12345678
```
