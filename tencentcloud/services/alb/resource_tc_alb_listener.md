Provides a resource to manage the forwarding rules of an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  listener_id      = "lst-d9p3k7wa"
  load_balancer_id = "alb-f8q2xk9m"

  rules = [{
    priority  = 10
    rule_name = "tf-test-rule"
    direction = "Request"

    actions = [{
      order = 1
      type  = "TargetGroup"
      target_group_config = {
        target_groups = [{
          target_group_id = "lbtg-a1b2c3d4"
          weight          = 10
        }]
      }
    }]

    conditions = [{
      type        = "Host"
      host_config = ["www.example.com"]
    }]

    tags = [{
      tag_key   = "env"
      tag_value = "test"
    }]
  }]
}
```

Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example rule-a1b2c3d4
```
