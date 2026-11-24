Provides a resource to create forwarding rules for an ALB (Application Load Balancer) listener.

Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  load_balancer_id = "alb-xxxxxx"
  listener_id      = "lst-xxxxxx"

  rules {
    priority  = 1
    direction = "Request"
    rule_name = "tf-example-rule"

    conditions {
      type        = "Host"
      host_config = ["www.example.com"]
    }

    conditions {
      type        = "Path"
      path_config = ["/api"]
    }

    actions {
      order = 1
      type  = "InsertHeader"

      insert_header_config {
        key        = "X-Custom-Env"
        value      = "test"
        value_type = "UserDefined"
      }
    }

    actions {
      order = 2
      type  = "TargetGroup"

      target_group_config {
        target_groups {
          target_group_id = "lbtg-xxxxxx"
          weight          = 10
        }

        target_group_sticky_session {
          enabled = true
          timeout = 1000
        }
      }
    }

    tags {
      tag_key   = "createdBy"
      tag_value = "terraform"
    }
  }
}
```

Import

ALB forwarding rules can be imported using the composite id load_balancer_id#listener_id, e.g.

```
terraform import tencentcloud_alb_rules.example alb-xxxxxx#lst-xxxxxx
```
