ALB listener certificates data source provides the list of certificates bound to an ALB listener.

Example Usage

```hcl
data "tencentcloud_alb_listener_certificates" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  certificate_type = "SVR"
}
```
