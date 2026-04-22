Provides a resource to create a TEO (EdgeOne) DNS record v5

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v5" "teo_dns_record_v5" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.makn.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 5
  status   = "enable"
}
```

Import

teo dns record v5 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v5.teo_dns_record_v5 zone-39quuimqg8r6#record-abc123
```
