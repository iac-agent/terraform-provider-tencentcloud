Provides a resource to manage a forwarding rule of an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  rule_name        = "tf-test-rule"
  priority         = 10
  direction        = "Request"

  actions {
    order = 1

    type = "TargetGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }

      target_group_sticky_session {
        enabled = true
        timeout = 1000
      }
    }
  }

  conditions {
    type        = "Host"
    host_config = ["www.example.com"]
  }

  tags {
    tag_key   = "env"
    tag_value = "test"
  }
}
```

Import

ALB rules can be imported using the composite id load_balancer_id#listener_id#rule_id, e.g.

```
terraform import tencentcloud_alb_rules.example alb-12345678#lst-12345678#rule-12345678
```
