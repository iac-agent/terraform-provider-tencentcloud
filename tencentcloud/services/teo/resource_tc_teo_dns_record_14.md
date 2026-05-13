Provides a resource to create a teo teo_dns_record_14

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_14" "teo_dns_record_14" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  name     = "a.makn.cn"
  location = "Default"
  priority = 5
  ttl      = 300
  weight   = -1
}
```

Import

teo teo_dns_record_14 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_14.teo_dns_record_14 {zoneId}#{recordId}
```
