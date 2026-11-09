---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain_1109_04"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain_1109_04"
description: |-
  Provides a resource to create a TEO alias domain.
---

# tencentcloud_teo_alias_domain_1109_04

Provides a resource to create a TEO alias domain.

~> **NOTE:** An alias domain shares the plan of its target domain. Only the alias domain mapping itself is managed here: deleting this resource removes the alias domain only, and does not affect the target domain or the zone.

-> **NOTE:** `cert_type` and `cert_id` are write-only on the cloud side (the query API does not return them), so Terraform cannot detect drift on these two attributes; their configured values are kept in state.

## Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_04" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
}
```

### With an SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_04" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String) Alias domain name.
* `target_name` - (Required, String) Target domain name.
* `zone_id` - (Required, String) Site ID.
* `cert_id` - (Optional, Set) Certificate IDs. Required when `cert_type` is `hosting`. The field is not returned by the query API, so the configured value is kept in state (limited drift detection).
* `cert_type` - (Optional, String) Certificate configuration. Valid values: `none` (no certificate), `hosting` (SSL hosted certificate). The cloud side defaults to `none` when omitted. The field is not returned by the query API, so the configured value is kept in state (limited drift detection).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource identifier in the form `<zone_id>#<alias_name>`.
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

TEO alias domains can be imported using the `zone_id#alias_name` composite ID, e.g.
```
$ terraform import tencentcloud_teo_alias_domain_1109_04.example zone-3fkff38fyw8s#alias.demo.com
```

