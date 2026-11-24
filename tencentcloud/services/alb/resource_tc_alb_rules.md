Provides a resource to create a forwarding rule of an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"

  rules {
    priority  = 1
    direction = "Request"
    rule_name = "tf-rule-example"

    conditions {
      type        = "Host"
      host_config = ["www.example.com"]
    }

    conditions {
      type        = "Path"
      path_config = ["/api"]
    }

    actions {
      order = 10
      type  = "InsertHeader"
      insert_header_config {
        key        = "x-tf-header"
        value      = "terraform"
        value_type = "UserDefined"
      }
    }

    actions {
      order = 20
      type  = "TargetGroup"
      target_group_config {
        target_groups {
          target_group_id = "lbtg-12345678"
          weight          = 10
        }
        target_group_sticky_session {
          enabled = false
        }
      }
    }

    tags {
      tag_key   = "env"
      tag_value = "test"
    }
  }
}
```

Import

ALB forwarding rule can be imported using the composite id load_balancer_id#listener_id#rule_id, e.g.

```
terraform import tencentcloud_alb_rules.example alb-12345678#lst-12345678#rule-12345678
```
