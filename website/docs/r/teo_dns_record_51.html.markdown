---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_51"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_51"
description: |-
  Provides a resource to create a teo dns_record_51
---

# tencentcloud_teo_dns_record_51

Provides a resource to create a teo dns_record_51

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_51" "dns_record_51" {
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

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. valid values are A, AAAA, MX, CNAME, TXT, NS, CAA, SRV. different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Zone id.
* `location` - (Optional, String) DNS record resolution route, not specified as default, indicates the default resolution route, which is effective for all regions. the resolution of line configuration is only applicable when the Type (DNS record type) is A, AAAA, or CNAME. the analysis of line configuration is only applicable to standard and enterprise packages. please refer to the analysis of line and corresponding code enumeration for values.
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modify time.
* `record_id` - DNS record id.
* `status` - DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.


## Import

teo dns_record_51 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_51.dns_record_51 {zoneId}#{recordId}
```

