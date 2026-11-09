---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain_1109_03"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain_1109_03"
description: |-
  Provides a resource to create a TEO alias domain.
---

# tencentcloud_teo_alias_domain_1109_03

Provides a resource to create a TEO alias domain.

~> **NOTE:** This feature is only supported by the Enterprise edition plan and is currently in beta testing. Please [contact us](https://cloud.tencent.com/online-service?from=connect-us) if you need to use it.

~> **NOTE:** `cert_type` and `cert_id` are not returned by the `DescribeAliasDomains` API, so Terraform cannot detect drift on these two fields. Removing either of them from the configuration keeps the previous cloud-side configuration; to remove a hosted certificate binding, set `cert_type` to `none` explicitly.

## Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_03" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

### With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_03" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String) Alias domain name. Changing it forces a new resource: ModifyAliasDomain uses AliasName to locate the alias domain, so it cannot be renamed in place.
* `target_name` - (Required, String) Target domain name. Always sent on update because ModifyAliasDomain requires it.
* `zone_id` - (Required, String) Site ID. Changing it forces a new resource: ModifyAliasDomain cannot relocate an alias domain to another site.
* `cert_id` - (Optional, Set) Certificate ID list, required when `cert_type` is `hosting`. When omitted on update the cloud side keeps the previous configuration. Not returned by DescribeAliasDomains, so drift cannot be detected on this field.
* `cert_type` - (Optional, String) Certificate configuration. Valid values: `none` (no certificate), `hosting` (SSL hosted certificate). Default value: `none`. When omitted on update, the cloud side keeps the previous configuration; pass `none` explicitly to remove the certificate binding. Not returned by DescribeAliasDomains, so drift cannot be detected on this field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_on` - Alias domain creation time.
* `forbid_mode` - Block mode. Valid values: `0` (not blocked), `11` (compliance blocked), `14` (not-registered blocked).
* `id` - Resource ID in the form `<zone_id>#<alias_name>`.
* `modified_on` - Alias domain modification time.
* `status` - Alias domain status. Valid values: `active` (effective), `pending` (deploying), `conflict` (reclaimed), `stop` (disabled).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateAliasDomain](https://www.tencentcloud.com/document/api/1145/51740)
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
* Read: [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
* Update:
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
  - [ModifyAliasDomain](https://www.tencentcloud.com/document/api/1145/51737)
* Delete:
  - [DeleteAliasDomain](https://www.tencentcloud.com/document/api/1145/51739)
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)

## Import

TEO alias domain can be imported using the composite id `<zone_id>#<alias_name>`, e.g.

```
terraform import tencentcloud_teo_alias_domain_1109_03.example zone-2qtuhspy7cr6#alias.demo.com
```

