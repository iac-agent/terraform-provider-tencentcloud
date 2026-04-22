---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v6"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v6"
description: |-
  Provides a resource to create a TEO DNS record.
---

# tencentcloud_teo_dns_record_v6

Provides a resource to create a TEO DNS record.

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_v6" "example" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "a.example.com"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  ttl      = 300
  weight   = -1
  priority = 0
  status   = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in the corresponding content according to the type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `location` - (Optional, String) DNS record resolution route. If not specified, the default is Default, which means the default resolution route and is effective in all regions. Resolution route configuration is only applicable when type is A, AAAA, or CNAME. Resolution route configuration is only applicable to Standard and Enterprise edition packages.
* `priority` - (Optional, Int) MX record priority, which takes effect only when type is MX. The smaller the value, the higher the priority. Users can specify a value range of 0-50. The default value is 0 if not specified.
* `status` - (Optional, String) DNS record resolution status. Valid values: enable (has taken effect); disable (has been disabled).
* `ttl` - (Optional, Int) Cache time in seconds. Users can specify a value range of 60-86400. Default value: 300.
* `weight` - (Optional, Int) DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, the default is -1, which means no weight is set. Weight configuration is only applicable when type is A, AAAA, or CNAME.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modification time.
* `record_id` - DNS record ID.


## Import

TEO DNS record can be imported using the joint id "zone_id#record_id", e.g.

```
terraform import tencentcloud_teo_dns_record_v6.example zone-39quuimqg8r6#record-abc123
```

