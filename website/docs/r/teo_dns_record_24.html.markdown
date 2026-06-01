---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_24"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_24"
description: |-
  Provides a resource to create a TEO dns record
---

# tencentcloud_teo_dns_record_24

Provides a resource to create a TEO dns record

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_24" "example" {
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

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in the corresponding content according to the Type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values are:
	- A: Points the domain name to an external IPv4 address, such as 8.8.8.8;
	- AAAA: Points the domain name to an external IPv6 address;
	- MX: Used for email servers. When there are multiple MX records, the lower the priority value, the higher the priority;
	- CNAME: Points the domain name to another domain name, which then resolves to the final IP address;
	- TXT: Identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);
	- NS: If you need to delegate the subdomain to another DNS service provider for resolution, you need to add an NS record. The root domain cannot add NS records;
	- CAA: Specifies the CA that can issue certificates for this site;
	- SRV: Identifies a server using a service, commonly used in Microsoft's directory management.
Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [Introduction to DNS Record Types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Site ID.
* `location` - (Optional, String) DNS record resolution route. If not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.

- Resolution route configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.
- Resolution route configuration is only applicable to Standard Edition and Enterprise Edition packages. For valid values, please refer to: [Resolution Routes and Corresponding Code Enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).
* `priority` - (Optional, Int) MX record priority, which takes effect only when Type (DNS record type) is MX. The smaller the value, the higher the priority. Users can specify a value range of 0-50. The default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. Users can specify a value range of 60-86400. The smaller the value, the faster the modification records will take effect in all regions. Default value: 300. Unit: seconds.
* `weight` - (Optional, Int) DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, the default is -1, which means no weight is set. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME. Note: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `dns_records` - DNS record list.
  * `content` - DNS record content. Fill in the corresponding content according to the Type value.
  * `created_on` - Creation time. Note: CreatedOn is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.
  * `location` - DNS record resolution route. If not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions. Resolution route configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME. For valid values, please refer to: [Resolution Routes and Corresponding Code Enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).
  * `modified_on` - Modification time. Note: ModifiedOn is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.
  * `name` - DNS record name.
  * `priority` - MX record priority. Value range: 0-50. The smaller the value, the higher the priority.
  * `record_id` - DNS record ID.
  * `status` - DNS record resolution status. Valid values:
	- enable: has taken effect;
	- disable: has been disabled.
Note: Status is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.
  * `ttl` - Cache time. Value range: 60-86400. The smaller the value, the faster the modification records will take effect in all regions. Unit: seconds.
  * `type` - DNS record type. Valid values are:
	- A: Points the domain name to an external IPv4 address, such as 8.8.8.8;
	- AAAA: Points the domain name to an external IPv6 address;
	- MX: Used for email servers. When there are multiple MX records, the lower the priority value, the higher the priority;
	- CNAME: Points the domain name to another domain name, which then resolves to the final IP address;
	- TXT: Identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);
	- NS: If you need to delegate the subdomain to another DNS service provider for resolution, you need to add an NS record. The root domain cannot add NS records;
	- CAA: Specifies the CA that can issue certificates for this site;
	- SRV: Identifies a server using a service, commonly used in Microsoft's directory management.
  * `weight` - DNS record weight. Value range: -1 to 100. A value of -1 means no weight is assigned, and a value of 0 means no resolution. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.
  * `zone_id` - Site ID. Note: ZoneId is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.
* `modified_on` - Modify time.
* `record_id` - DNS record id.
* `status` - DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.


## Import

TEO dns record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_24.example zone-39quuimqg8r6#rec-abcdefghij
```

