---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v2"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v2"
description: |-
  Provides a resource to create a teo teo_dns_record_v2
---

# tencentcloud_teo_dns_record_v2

Provides a resource to create a teo teo_dns_record_v2

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v2" "teo_dns_record_v2" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.makn.cn"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 5
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. fill in the corresponding content according to the type value.
* `name` - (Required, String) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type.
* `zone_id` - (Required, String, ForceNew) Zone id.
* `location` - (Optional, String) DNS record resolution route, not specified as default, indicates the default resolution route, which is effective for all regions.
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `record_id` - DNS record id.


## Import

teo teo_dns_record_v2 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v2.teo_dns_record_v2 {zoneId}#{recordId}
```

