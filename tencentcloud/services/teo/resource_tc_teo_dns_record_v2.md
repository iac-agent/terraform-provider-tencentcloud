Provides a resource to create a TEO DNS record v2.

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v2" "example" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.example.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
}
```

Import

TEO DNS record v2 can be imported using the joint id "zone_id#record_id", e.g.

```
terraform import tencentcloud_teo_dns_record_v2.example zone-39quuimqg8r6#record-abc123
```
