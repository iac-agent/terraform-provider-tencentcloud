Provides a resource to create a TEO inference domain

Example Usage

```hcl
resource "tencentcloud_teo_inference_domain" "example" {
  zone_id     = "zone-3v0yhgt4m1z9"
  service_id  = "inf-vs06pb1cdz61"
  domain      = "test-ai.eo.cn"
  auth_switch = "On"
}
```

Import

teo inference domain can be imported using the id, e.g.

```hcl
terraform import tencentcloud_teo_inference_domain.example zone-3v0yhgt4m1z9#inf-vs06pb1cdz61#test-ai.eo.cn
```
