---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_22"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_22"
description: |-
  Provides a resource to create a TEO (EdgeOne) DNS record.
---

# tencentcloud_teo_dns_record_22

Provides a resource to create a TEO (EdgeOne) DNS record.

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_22" "example" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.example.com"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in the corresponding content based on the Type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: <li>A: Points the domain to an external IPv4 address, e.g., 8.8.8.8;</li><li>AAAA: Points the domain to an external IPv6 address;</li><li>MX: Used for mail servers. Lower priority values are preferred when multiple MX records exist;</li><li>CNAME: Points the domain to another domain, which then resolves to the final IP address;</li><li>TXT: Identifies and describes the domain, commonly used for domain verification and SPF records (anti-spam);</li><li>NS: Required when delegating subdomain resolution to other DNS providers. NS records cannot be added to the root domain;</li><li>CAA: Specifies the CA that can issue certificates for this site;</li><li>SRV: Identifies a server using a specific service, commonly found in Microsoft system directory management.</li>
Different record types (e.g., SRV, CAA) have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [DNS Record Type Introduction](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Site ID. Cannot be null or empty string.
* `location` - (Optional, String) DNS record resolution line. Defaults to Default, which means the default resolution line that takes effect for all regions.

- Resolution line configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.
- Resolution line configuration only applies to Standard and Enterprise editions. For valid values, please refer to: [Resolution Line and Corresponding Code Enumeration](https://cloud.tencent.com/document/product/1552/112542).
* `priority` - (Optional, Int) MX record priority. This parameter only takes effect when Type (DNS record type) is MX. The smaller the value, the higher the priority. The user-specified value range is 0~50. Default is 0.
* `status` - (Optional, String) DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.
* `ttl` - (Optional, Int) Cache time. The user-specified value range is 60~86400. The smaller the value, the faster the record modification takes effect in each region. Default is 300, in seconds.
* `weight` - (Optional, Int) DNS record weight. The user-specified value range is -1~100. Setting to 0 means no resolution. Default is -1, which means no weight is set. Weight configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.<br>Note: Under the same subdomain, different DNS records with the same resolution line should either all have weights set or all have no weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `dns_records` - List of DNS records.
  * `content` - DNS record content. Fill in the corresponding content based on the Type value.
  * `created_on` - Creation time.<br>Note: CreatedOn is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.
  * `location` - DNS record resolution line. Defaults to Default, which means the default resolution line that takes effect for all regions.<br>Resolution line configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.<br>For valid values, please refer to: [Resolution Line and Corresponding Code Enumeration](https://cloud.tencent.com/document/product/1552/112542).
  * `modified_on` - Modification time.<br>Note: ModifiedOn is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.
  * `name` - DNS record name.
  * `priority` - MX record priority. Value range is 0~50. The smaller the value, the higher the priority.
  * `record_id` - DNS record ID.
  * `status` - DNS record resolution status. Valid values: <li>enable: has taken effect;</li><li>disable: has been disabled.</li>Note: Status is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.
  * `ttl` - Cache time. Value range is 60~86400. The smaller the value, the faster the record modification takes effect in each region, in seconds.
  * `type` - DNS record type. Valid values:
<li>A: Points the domain to an external IPv4 address, e.g., 8.8.8.8;</li>
<li>AAAA: Points the domain to an external IPv6 address;</li>
<li>MX: Used for mail servers. Lower priority values are preferred when multiple MX records exist;</li>
<li>CNAME: Points the domain to another domain, which then resolves to the final IP address;</li>
<li>TXT: Identifies and describes the domain, commonly used for domain verification and SPF records (anti-spam);</li>
<li>NS: Required when delegating subdomain resolution to other DNS providers. NS records cannot be added to the root domain;</li>
<li>CAA: Specifies the CA that can issue certificates for this site;</li>
<li>SRV: Identifies a server using a specific service, commonly found in Microsoft system directory management.</li>.
  * `weight` - DNS record weight. Value range is -1~100. -1 means no weight is assigned, 0 means no resolution. Weight configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.
  * `zone_id` - Site ID.<br>Note: ZoneId is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.
* `modified_on` - Modify time.
* `record_id` - DNS record ID.


## Import

TEO (EdgeOne) DNS record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_22.example {zoneId}#{recordId}
```

