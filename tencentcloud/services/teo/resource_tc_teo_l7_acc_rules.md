Provides a resource to manage TEO L7 acceleration rules.

~> **NOTE:** This feature only supports the sites in the plans of the Standard Edition and the Enterprise Edition.

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rules" "example" {
  zone_id = "zone-36bjhygh1bxe"

  rules {
    rule_name   = "Web Acceleration"
    status      = "enable"
    description = ["rule description 1"]
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
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

  rules {
    rule_name   = "API Acceleration"
    status      = "enable"
    description = ["rule description 2"]
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      actions {
        name = "Cache"
        cache_parameters {
          no_cache {
            switch = "on"
          }
        }
      }
    }
  }
}
```

Import

TEO L7 acceleration rules can be imported using the `zone_id`, e.g.

```
terraform import tencentcloud_teo_l7_acc_rules.example zone-36bjhygh1bxe
```