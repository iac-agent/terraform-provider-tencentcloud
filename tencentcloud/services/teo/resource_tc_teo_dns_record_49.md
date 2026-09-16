Provides a resource to create a TEO dns record

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_49" "example" {
  zone_id  = "zone-2o3h21ed2t68"
  name     = "www.example.com"
  type     = "A"
  content  = "1.2.3.4"
  location = "Default"
  ttl      = 300
  weight   = -1
}
```

Import

TEO dns record can be imported using the composite id with the format `zone_id#record_id`, e.g.

```
terraform import tencentcloud_teo_dns_record_49.example zone-2o3h21ed2t68#record-3qcqn0ejch9v
```
