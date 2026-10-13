---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_17"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_17"
description: |-
  Provides a resource to create a teo teo_dns_record_17
---

# tencentcloud_teo_dns_record_17

Provides a resource to create a teo teo_dns_record_17

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_17" "teo_dns_record_17" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
  # Optional pagination limit for the DescribeDnsRecords api call used by read,
  # range: 1-1000, defaults to 1000 (the cloud api max value) when not set.
  Limit = 1000
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. valid values are:
	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;
	- AAAA: points the domain name to an external ipv6 address;
	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;
	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;
	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);
	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;
	- CAA: specifies the ca that can issue certificates for this site;
	- SRV: identifies a server using a service, commonly used in microsoft's directory management.
* `zone_id` - (Required, String, ForceNew) Zone id.
* `Limit` - (Optional, Int) Pagination limit for the DescribeDnsRecords api call used by read. default value: 20, max value: 1000. if not set, 1000 will be used to make sure the query by record id filter is not affected by the api default pagination.
* `location` - (Optional, String) DNS record resolution route, not specified as default, indicates the default resolution route, which is effective for all regions.
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.
* `status` - (Optional, String) DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modify time.
* `record_id` - DNS record id.


## Import

teo teo_dns_record_17 can be imported using the composite id (joined with `#`, in the form `{zoneId}#{recordId}`), e.g.

```
terraform import tencentcloud_teo_dns_record_17.teo_dns_record_17 {zoneId}#{recordId}
```

