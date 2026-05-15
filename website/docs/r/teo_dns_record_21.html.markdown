---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_21"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_21"
description: |-
  Provides a resource to create a TEO dns record
---

# tencentcloud_teo_dns_record_21

Provides a resource to create a TEO dns record

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record_21" "example" {
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

* `content` - (Required, String) DNS record content. Fill in the corresponding content according to the Type value. If the domain is in Chinese, Korean, or Japanese, convert it to punycode before input.
* `name` - (Required, String) DNS record name. If the domain is in Chinese, Korean, or Japanese, convert it to punycode before input.
* `type` - (Required, String) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV. Different record types have different requirements for the host record name and record value format. For details, see: [DNS Record Types](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Site ID. This field is required.
* `location` - (Optional, String) DNS record resolution line. Default is Default, representing all regions. Resolution line configuration only applies when Type is A, AAAA, or CNAME. For standard and enterprise editions only. For valid values, see: [Resolution Lines and Codes](https://cloud.tencent.com/document/product/1552/112542).
* `priority` - (Optional, Int) MX record priority. Only effective when Type is MX. Lower values indicate higher priority. Valid range: 0-50. Default: 0.
* `ttl` - (Optional, Int) Cache time in seconds. Valid range: 60-86400. Smaller values mean faster propagation. Default: 300.
* `weight` - (Optional, Int) DNS record weight. Valid range: -1 to 100. 0 means no resolution, -1 means no weight. Weight configuration only applies when Type is A, AAAA, or CNAME. Note: DNS records under the same subdomain with the same resolution line should either all have weights set or all have weights unset.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `dns_records` - List of DNS records.
  * `content` - DNS record content. Fill in the corresponding content according to the Type value.
  * `created_on` - Creation time. Note: CreatedOn is output-only in ModifyDnsRecords and will be ignored if passed as input.
  * `location` - DNS record resolution line. Default is Default, representing all regions. Resolution line configuration only applies when Type is A, AAAA, or CNAME. For valid values, see: [Resolution Lines and Codes](https://cloud.tencent.com/document/product/1552/112542).
  * `modified_on` - Last modification time. Note: ModifiedOn is output-only in ModifyDnsRecords and will be ignored if passed as input.
  * `name` - DNS record name.
  * `priority` - MX record priority. Valid range: 0-50. Lower values indicate higher priority.
  * `record_id` - DNS record ID.
  * `status` - DNS record resolution status. Valid values: enable (has taken effect), disable (has been disabled). Note: Status is output-only in ModifyDnsRecords and will be ignored if passed as input.
  * `ttl` - Cache time in seconds. Valid range: 60-86400. Smaller values mean faster propagation.
  * `type` - DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
  * `weight` - DNS record weight. Valid range: -1 to 100. -1 means no weight, 0 means no resolution. Weight configuration only applies when Type is A, AAAA, or CNAME.
  * `zone_id` - Site ID. Note: ZoneId is output-only in ModifyDnsRecords and will be ignored if passed as input.
* `modified_on` - Last modification time.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable (has taken effect), disable (has been disabled).


## Import

TEO dns record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_21.example zone-39quuimqg8r6#rec-abcdefghij
```

