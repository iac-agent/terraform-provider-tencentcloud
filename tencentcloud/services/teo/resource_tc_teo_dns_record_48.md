Provides a resource to create a TEO (TencentCloud EdgeOne) dns_record

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_48" "example" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  name     = "a.example.cn"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
}
```

Import

TEO dns_record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_48.example {zoneId}#{recordId}
```
