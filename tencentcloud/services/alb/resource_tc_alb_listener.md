Provides a resource to manage an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-test-listener"

  default_actions {
    type = "ForwardGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-abcdefgh"
        weight          = 10
      }
    }
  }
}
```

Import

ALB listener can be imported using the composite id `load_balancer_id#listener_id`, e.g.

```
terraform import tencentcloud_alb_listener.example alb-12345678#lst-12345678
```
