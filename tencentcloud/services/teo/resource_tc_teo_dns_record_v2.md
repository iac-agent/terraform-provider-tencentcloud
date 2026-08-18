Provides a resource to create a teo teo_dns_record_v2

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v2" "teo_dns_record_v2" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.makn.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 5
}
```

Import

teo teo_dns_record_v2 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v2.teo_dns_record_v2 {zoneId}#{recordId}
```
