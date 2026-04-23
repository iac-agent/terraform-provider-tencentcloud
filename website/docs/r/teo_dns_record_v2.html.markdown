---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v2"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v2"
description: |-
  Provides a resource to create a TEO DNS record v2.
---

# tencentcloud_teo_dns_record_v2

Provides a resource to create a TEO DNS record v2.

## Example Usage

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
terraform import tencentcloud_teo_dns_record_v2.example zone-39quuimqg8r6#record-abc123
```

