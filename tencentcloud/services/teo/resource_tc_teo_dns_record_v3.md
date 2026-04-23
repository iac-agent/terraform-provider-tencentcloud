Provides a resource to create a TEO (TencentCloud EdgeOne) DNS record v3

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v3" "a_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.example.com"
  ttl      = 300
  weight   = -1
}
```

```hcl
resource "tencentcloud_teo_dns_record_v3" "mx_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "MX"
  content  = "mail.example.com"
  name     = "mx.example.com"
  priority = 5
  ttl      = 300
}
```

Import

TEO DNS record v3 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v3.a_record {zoneId}#{recordId}
```
