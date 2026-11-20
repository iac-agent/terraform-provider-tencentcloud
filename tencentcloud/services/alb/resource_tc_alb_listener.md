Provides a resource to create an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
      target_group_sticky_session {
        enabled = false
        timeout = 1000
      }
    }
  }

  gzip_enabled = true
  idle_timeout = 15
}
```

HTTPS listener with certificates:

```hcl
resource "tencentcloud_alb_listener" "https_example" {
  load_balancer_id   = "alb-12345678"
  listener_port      = 443
  listener_protocol  = "HTTPS"
  listener_name      = "tf-example-https"
  certificate_ids    = ["cert-12345678"]
  security_policy_id = "tls-12345678"
  http2enabled       = true

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  x_forwarded_for_config {
    x_forwarded_for_mode = "append"
  }
}
```

Import

ALB listener can be imported using the loadBalancerId#listenerId, e.g.

```
terraform import tencentcloud_alb_listener.example alb-12345678#lst-12345678
```
