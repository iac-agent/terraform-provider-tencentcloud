Provides a resource to create a TEO dns record

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_24" "example" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.example.com"
  type     = "A"
  content  = "1.2.3.4"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
}
```

Import

TEO dns record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_24.example zone-39quuimqg8r6#rec-abcdefghij
```
