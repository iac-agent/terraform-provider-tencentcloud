Provides a resource to create a teo teo_dns_record_15

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_15" "teo_dns_record_15" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
}
```

Import

teo teo_dns_record_15 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_15.teo_dns_record_15 {zoneId}#{recordId}
```
