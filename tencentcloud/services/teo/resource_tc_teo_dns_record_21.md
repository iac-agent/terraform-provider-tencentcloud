提供 TEO（EdgeOne）DNS 记录资源。

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_21" "example" {
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

TEO（EdgeOne）DNS 记录可以通过 ID 导入，例如：

```
terraform import tencentcloud_teo_dns_record_21.example zone-39quuimqg8r6#rec-abcdefghij
```