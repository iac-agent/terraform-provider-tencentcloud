Provides a resource to create a teo dns_record_50

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_50" "dns_record_50" {
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

```hcl
resource "tencentcloud_teo_dns_record_50" "dns_record_50" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.makn.cn"
  type     = "A"
  content  = "1.2.3.5"

  filters {
    name   = "name"
    values = ["a.makn.cn"]
    fuzzy  = true
  }

  sort_by    = "created-on"
  sort_order = "desc"
  match      = "all"
}
```

Import

teo dns_record_50 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_50.dns_record_50 {zoneId}#{recordId}
```
