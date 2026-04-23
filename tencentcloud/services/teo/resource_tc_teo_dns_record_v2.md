Provides a resource to create a TEO (EdgeOne) DNS record.

Example Usage

A record

```hcl
resource "tencentcloud_teo_dns_record_v2" "a_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.example.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
}
```

CNAME record

```hcl
resource "tencentcloud_teo_dns_record_v2" "cname_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "cdn.example.cn"
  type     = "CNAME"
  content  = "cdn.example.cn.cdn.dnsv1.com"
  location = "Default"
  ttl      = 600
}
```

MX record

```hcl
resource "tencentcloud_teo_dns_record_v2" "mx_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "mail.example.cn"
  type     = "MX"
  content  = "mailserver.example.cn"
  priority = 10
  location = "Default"
  ttl      = 300
}
```

TXT record

```hcl
resource "tencentcloud_teo_dns_record_v2" "txt_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "_verification.example.cn"
  type     = "TXT"
  content  = "v=spf1 include:spf.mail.example.cn ~all"
  location = "Default"
  ttl      = 300
}
```

AAAA record with custom resolution route

```hcl
resource "tencentcloud_teo_dns_record_v2" "aaaa_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "ipv6.example.cn"
  type     = "AAAA"
  content  = "2402:4e00:1900:1700:0:9556:1bca:1fb"
  location = "CM"
  ttl      = 300
  weight   = 50
}
```

Import

TEO DNS record v2 can be imported using the joint id "zone_id#record_id", e.g.

```
terraform import tencentcloud_teo_dns_record_v2.a_record zone-39quuimqg8r6#record-abc123
```
