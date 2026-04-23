---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v2"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v2"
description: |-
  Provides a resource to create a TEO (EdgeOne) DNS record.
---

# tencentcloud_teo_dns_record_v2

Provides a resource to create a TEO (EdgeOne) DNS record.

## Example Usage

### A record

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

### CNAME record

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

### MX record

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

### TXT record

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

### AAAA record with custom resolution route

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

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in the corresponding content according to the type value.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `location` - (Optional, String) DNS record resolution route. If not specified, the default is Default.
* `priority` - (Optional, Int) MX record priority. Range: 0-50. Default: 0. Only takes effect when type is MX.
* `ttl` - (Optional, Int) Cache time in seconds. Range: 60-86400. Default: 300.
* `weight` - (Optional, Int) DNS record weight. Range: -1 to 100. Default: -1, which means no weight is set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modification time.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable, disable.


## Import

TEO DNS record v2 can be imported using the joint id "zone_id#record_id", e.g.

```
terraform import tencentcloud_teo_dns_record_v2.a_record zone-39quuimqg8r6#record-abc123
```

