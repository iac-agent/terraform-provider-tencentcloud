Provides a resource to create a listener for an ALB (Application Load Balancer) instance.

Example Usage

Create an HTTP listener forwarding to a target group

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-xxxxxxxx"
  listener_protocol = "HTTP"
  listener_port     = 80
  listener_name     = "tf-example"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-xxxxxxxx"
        weight          = 10
      }
    }
  }
}
```

Create an HTTPS listener with server certificate and request/idle timeout

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id    = "alb-xxxxxxxx"
  listener_protocol   = "HTTPS"
  listener_port       = 443
  listener_name       = "tf-example"
  certificate_ids     = ["cert-xxxxxxxx"]
  security_policy_id  = "tls-xxxxxxxx"
  idle_timeout        = 30
  request_timeout     = 60
  gzip_enabled        = true
  http2enabled        = true

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-xxxxxxxx"
        weight          = 10
      }
      target_group_sticky_session {
        enabled = true
        timeout = 1000
      }
    }
  }
}
```

Import

ALB listener can be imported using the `load_balancer_id#listener_id` format, e.g.

```
terraform import tencentcloud_alb_listener.example alb-xxxxxxxx#lst-xxxxxxxx
```
