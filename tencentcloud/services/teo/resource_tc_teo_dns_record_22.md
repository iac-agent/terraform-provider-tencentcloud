提供 TEO（EdgeOne）DNS 记录资源。

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_22" "example" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.example.com"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
```

Import

TEO（EdgeOne）DNS 记录可以通过 ID 导入，例如：

```
terraform import tencentcloud_teo_dns_record_22.example {zoneId}#{recordId}
```