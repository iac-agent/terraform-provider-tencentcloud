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

* `content` - (Required, String) DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. valid values are:
<li>A: points the domain name to an external IPv4 address, such as 8.8.8.8;</li>
<li>AAAA: points the domain name to an external IPv6 address;</li>
<li>MX: used for mail server. when there are multiple MX records, the lower the priority, the more preferred;</li>
<li>CNAME: points the domain name to another domain name, and then resolves the final IP address by that domain name;</li>
<li>TXT: identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);</li>
<li>NS: if you need to delegate the subdomain to another DNS provider for resolution, you need to add an NS record. NS records cannot be added to the root domain;</li>
<li>CAA: specifies the CA that can issue certificates for this site;</li>
<li>SRV: identifies that a server uses a certain service, commonly found in Microsoft system directory management.</li>
Different record types such as SRV and CAA have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [DNS Record Type Introduction](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Zone id.
* `location` - (Optional, String) DNS record resolution route. if not specified, the default is Default, which means the default resolution route and is effective in all regions.

- Resolution route configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME.
- Resolution route configuration is only applicable for standard and enterprise edition packages. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://cloud.tencent.com/document/product/1552/112542).
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (DNS record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME. note: for the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `dns_records` - DNS record list.
  * `content` - DNS record content. fill in the corresponding content according to the type value.
  * `created_on` - Creation time. note: CreatedOn is only used as an output parameter, and cannot be used as an input parameter in ModifyDnsRecords. if this parameter is passed, it will be ignored.
  * `location` - DNS record resolution route. if not specified, the default is Default, which means the default resolution route and is effective in all regions. resolution route configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://cloud.tencent.com/document/product/1552/112542).
  * `modified_on` - Modify time. note: ModifiedOn is only used as an output parameter, and cannot be used as an input parameter in ModifyDnsRecords. if this parameter is passed, it will be ignored.
  * `name` - DNS record name.
  * `priority` - MX record priority. value range 0-50. the smaller the value, the higher the priority.
  * `record_id` - DNS record id.
  * `status` - DNS record resolution status. valid values: enable: has taken effect; disable: has been disabled. note: Status is only used as an output parameter, and cannot be used as an input parameter in ModifyDnsRecords. if this parameter is passed, it will be ignored.
  * `ttl` - Cache time. value range 60-86400. the smaller the value, the faster the modification records will take effect in all regions. unit: seconds.
  * `type` - DNS record type. valid values are:
<li>A: points the domain name to an external IPv4 address, such as 8.8.8.8;</li>
<li>AAAA: points the domain name to an external IPv6 address;</li>
<li>MX: used for mail server. when there are multiple MX records, the lower the priority, the more preferred;</li>
<li>CNAME: points the domain name to another domain name, and then resolves the final IP address by that domain name;</li>
<li>TXT: identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);</li>
<li>NS: if you need to delegate the subdomain to another DNS provider for resolution, you need to add an NS record. NS records cannot be added to the root domain;</li>
<li>CAA: specifies the CA that can issue certificates for this site;</li>
<li>SRV: identifies that a server uses a certain service, commonly found in Microsoft system directory management.</li>.
  * `weight` - DNS record weight. value range -1 to 100. -1 means no weight is assigned, 0 means no resolution. weight configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME.
  * `zone_id` - Zone id. note: ZoneId is only used as an output parameter, and cannot be used as an input parameter in ModifyDnsRecords. if this parameter is passed, it will be ignored.
* `modified_on` - Modify time.
* `record_id` - DNS record id.
* `status` - DNS record resolution status. valid values:
<li>enable: has taken effect;</li>
<li>disable: has been disabled.</li>.


## Import

TEO dns record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_24.example zone-39quuimqg8r6#rec-abcdefghij
```

