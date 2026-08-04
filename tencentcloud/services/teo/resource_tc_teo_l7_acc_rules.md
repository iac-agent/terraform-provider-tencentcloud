Provides a resource to manage the complete set of TEO L7 acceleration rules under a zone.

~> **NOTE:** This resource manages all L7 acceleration rules under a zone as a single declarative block. For managing individual rules independently, use `tencentcloud_teo_l7_acc_rule_v2` instead.

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rules" "example" {
  zone_id = "zone-3fkff38fyw8s"

  rules {
    rule_name   = "Web Acceleration"
    status      = "enable"
    description = ["description"]
    branches {
      condition = "$${http.request.host} in ['www.example.com']"
      actions {
        name = "Cache"
        cache_parameters {
          custom_time {
            cache_time           = 2592000
            ignore_cache_control = "off"
            switch               = "on"
          }
        }
      }
    }
  }
}
```

Import

TEO l7 acc rules can be imported using the zone id, e.g.

```
terraform import tencentcloud_teo_l7_acc_rules.example zone-3fkff38fyw8s
```