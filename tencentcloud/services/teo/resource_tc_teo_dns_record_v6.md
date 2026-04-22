Provides a resource to create a TEO DNS record.

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v6" "example" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.example.com"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
  status   = "enable"
}
```

Import

TEO DNS record can be imported using the joint id "zone_id#record_id", e.g.

```
terraform import tencentcloud_teo_dns_record_v6.example zone-39quuimqg8r6#record-abc123
```
