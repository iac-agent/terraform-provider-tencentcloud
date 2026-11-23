Provides a resource to create a listener of an application load balancer (ALB).

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-xxxxxxxx"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example-listener"
  idle_timeout      = 30
  request_timeout   = 60

  default_actions = [
    {
      type = "Forward"
      target_group_config = {
        target_groups = [
          {
            target_group_id = "lbtg-xxxxxxxx"
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

  x_forwarded_for_config = {
    x_forwarded_for_mode = "append"
  }

  tags = [
    {
      tag_key   = "createdBy"
      tag_value = "terraform"
    }
  ]
}
```

HTTPS listener with certificates:

```hcl
resource "tencentcloud_alb_listener" "https" {
  load_balancer_id   = "alb-xxxxxxxx"
  listener_port      = 443
  listener_protocol  = "HTTPS"
  listener_name      = "tf-example-https-listener"
  certificate_ids    = ["cert-xxxxxxxx"]
  security_policy_id = "tls-xxxxxxxx"
  http2enabled       = true

  default_actions = [
    {
      type = "Forward"
      target_group_config = {
        target_groups = [
          {
            target_group_id = "lbtg-xxxxxxxx"
          }
        ]
      }
    }
  ]
}
```

Import

ALB listeners can be imported using the composite id, e.g. `load_balancer_id#listener_id`.

```
terraform import tencentcloud_alb_listener.example alb-xxxxxxxx#lst-xxxxxxxx
```

-> `load_balancer_id`, `listener_port` and `listener_protocol` are immutable: changing any of them forces resource replacement. The `client_token` attribute is only used for request idempotency and is never echoed back by the cloud, so it must not be changed between applies.
