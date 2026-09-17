---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_50"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_50"
description: |-
  Provides a resource to create a teo dns_record_50
---

# tencentcloud_teo_dns_record_50

Provides a resource to create a teo dns_record_50

## Example Usage

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
  zone_id = "zone-39quuimqg8r6"
  name    = "a.makn.cn"
  type    = "A"
  content = "1.2.3.5"

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

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content, which should match the value of `type`.
* `name` - (Required, String) DNS record name. For domains in Chinese, Japanese or Korean, convert it to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: `A`, `AAAA`, `MX`, `CNAME`, `TXT`, `NS`, `CAA`, `SRV`.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `filters` - (Optional, List) Query filter conditions of `DescribeDnsRecords` used on each read. The precise `id` filter of the resource itself is always sent first, and user filters are appended after it.
* `location` - (Optional, String) Resolution route, defaults to `Default`. Only applicable when `type` is `A`, `AAAA` or `CNAME`.
* `match` - (Optional, String) Match mode of `DescribeDnsRecords`. Valid values: `all`, `any`. Defaults to `all`.
* `priority` - (Optional, Int) MX record priority, range 0-50. A smaller value indicates a higher priority. Only applicable when `type` is `MX`.
* `sort_by` - (Optional, String) Sort key of `DescribeDnsRecords`. Valid values: `content`, `created-on`, `name`, `ttl`, `type`. Defaults to a combined sort by `type` and `name`.
* `sort_order` - (Optional, String) Sort order of `DescribeDnsRecords`. Valid values: `asc`, `desc`. Defaults to `asc`.
* `ttl` - (Optional, Int) Cache time in seconds, range 60-86400, defaults to 300.
* `weight` - (Optional, Int) Record weight, range -1-100. `-1` means no weight, `0` means no resolution. Only applicable when `type` is `A`, `AAAA` or `CNAME`.

The `filters` object supports the following:

* `name` - (Required, String) Filter field. Valid values: `id`, `name`, `content`, `type`, `ttl`.
* `values` - (Required, List) Filter values, up to 20.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy matching.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Last modification time.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: `enable`, `disable`.


## Import

teo dns_record_50 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_50.dns_record_50 {zoneId}#{recordId}
```

