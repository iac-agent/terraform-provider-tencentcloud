Provides a resource to create a TEO DNS record

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_16" "example" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  name     = "a.example.com"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
}
```

Import

TEO DNS record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_16.example {zoneId}#{recordId}
```
