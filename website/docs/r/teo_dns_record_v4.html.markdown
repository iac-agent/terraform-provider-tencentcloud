---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v4"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v4"
description: |-
  Provides a resource to create a TEO DNS record
---

# tencentcloud_teo_dns_record_v4

Provides a resource to create a TEO DNS record

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v4" "example" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in the corresponding content according to the type value.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `location` - (Optional, String) DNS record resolution route. If not specified, the default is Default, which means the default resolution route and is effective in all regions.
* `priority` - (Optional, Int) MX record priority, which takes effect only when type is MX. Value range: 0-50. Default: 0.
* `status` - (Optional, String) DNS record resolution status. Valid values: enable (active), disable (disabled).
* `ttl` - (Optional, Int) Cache time in seconds. Value range: 60-86400. Default: 300.
* `weight` - (Optional, Int) DNS record weight. Value range: -1 to 100. -1 means no weight is set. 0 means no resolution.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modification time.
* `record_id` - DNS record ID.


## Import

TEO DNS record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v4.example zone-39quuimqg8r6#record-abc123
```

