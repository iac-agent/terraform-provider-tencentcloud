Use this data source to query the list of certificates bound to an Application Load Balancer listener.

Example Usage

```hcl
data "tencentcloud_alb_listener_certificates" "example" {
  certificate_type = "SVR"
  listener_id      = "lst-d9p3k7wa"
  load_balancer_id = "alb-f8q2xk9m"
}
```
