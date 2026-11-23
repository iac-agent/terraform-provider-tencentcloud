Provides a resource to create an ALB (Application Load Balancer) listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id   = "alb-xxxxxxxx"
  listener_port      = 80
  listener_protocol  = "HTTP"
  listener_name      = "tf-example"
  idle_timeout       = 30
  request_timeout    = 60
  gzip_enabled       = true

  default_actions = [
    {
      type = "Forward"
      target_group_config = {
        target_groups = [
          {
            target_group_id = "lbtg-xxxxxxxx"
            weight          = 10
          },
        ]
        target_group_sticky_session = {
          enabled = false
          timeout = 1000
        }
      }
    },
  ]

  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

ALB listener can be imported using the composite id, e.g. `load_balancer_id#listener_id`.

```
terraform import tencentcloud_alb_listener.example alb-xxxxxxxx#lst-xxxxxxxx
```

-> `load_balancer_id`, `listener_port` and `listener_protocol` are immutable: changing any of them forces resource replacement. The `certificate_ids`, `ca_certificate_ids` and `default_actions` fields are not echoed back by the DescribeListeners API, so their values are preserved from the Terraform configuration on read.
