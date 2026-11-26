Use this data source to query the list of ALB listeners.

Example Usage

```hcl
data "tencentcloud_alb_listeners" "example" {
  load_balancer_id = "alb-f8q2xk9m"
  listener_ids     = ["lst-a1b2c3d4"]

  filters = [
    {
      name   = "Protocol"
      values = ["HTTP"]
    },
  ]
}
```
