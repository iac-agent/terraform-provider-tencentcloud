Provides a resource to manage an Application Load Balancer (ALB) listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_name     = "tf-test-listener"
  listener_port     = 80
  listener_protocol = "HTTP"

  default_actions {
    type = "TargetGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  gzip_enabled   = true
  idle_timeout   = 15
  request_timeout = 60
}
```

Import

ALB listener can be imported using the composite id `load_balancer_id#listener_id`, e.g.

```
terraform import tencentcloud_alb_listener.example alb-12345678#lst-12345678
```
