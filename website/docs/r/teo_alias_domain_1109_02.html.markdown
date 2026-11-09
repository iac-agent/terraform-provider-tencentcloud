---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain_1109_02"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain_1109_02"
description: |-
  Provides a resource to create a TEO alias domain. An alias domain points to a target domain within the same site and shares its acceleration configuration.
---

# tencentcloud_teo_alias_domain_1109_02

Provides a resource to create a TEO alias domain. An alias domain points to a target domain within the same site and shares its acceleration configuration.

~> **NOTE:** This resource is implemented with terraform-plugin-framework. `cert_type` and `cert_id` are not returned by the DescribeAliasDomains API, so their values cannot be drift-detected after apply: they are kept from the configuration, and out-of-band changes on the cloud side are not detected.

## Example Usage

### Basic

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
}
```

### With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

### Switch the certificate configuration back to none explicitly

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String) Alias domain name.
* `target_name` - (Required, String) Target domain name that the alias domain points to.
* `zone_id` - (Required, String) Site ID.
* `cert_id` - (Optional, Set) Certificate ID list, required when `cert_type` is `hosting`. DescribeAliasDomains does not return this field, so its value is kept from the configuration and drift detection for it is not possible.
* `cert_type` - (Optional, String) Certificate configuration. Valid values: `none` (no configuration), `hosting` (SSL hosted certificate). Default value: `none`. On update, omitting this attribute keeps the previous configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alias_domains` - Alias domain details. Contains exactly one element describing this alias domain.
* `id` - Resource identifier in the form `<zone_id>#<alias_name>`.

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

TEO alias domains can be imported using the composite id `<zone_id>#<alias_name>`, e.g.

```
$ terraform import tencentcloud_teo_alias_domain_1109_02.example zone-3fkff38fyw8s#alias.demo.com
```

When destroying this resource, only the alias domain itself is deleted; the target domain and any hosted certificates are left untouched. Removing `cert_type` or `cert_id` from the configuration keeps the previous certificate configuration on the cloud side (the documented update semantics); set `cert_type = "none"` explicitly to switch to no certificate.

