Provides a resource to create a teo dns_record_51

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_51" "dns_record_51" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.makn.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
}
```

Import

teo dns_record_51 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_51.dns_record_51 {zoneId}#{recordId}
```
