Provides a resource to create a TEO DNS record.

Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "teo_dns_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
```

Import

TEO DNS record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record.teo_dns_record {zoneId}#{recordId}
```
