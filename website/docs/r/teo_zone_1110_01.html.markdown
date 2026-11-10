---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_1110_01"
sidebar_current: "docs-tencentcloud-resource-teo_zone_1110_01"
description: |-
  Provides a resource to create a TEO zone.
---

# tencentcloud_teo_zone_1110_01

Provides a resource to create a TEO zone.

~> **NOTE:** A freshly created zone only serves traffic after a plan is bound to it (`plan_id`) and after the ownership verification returned by `ownership_verification` is finished. When `plan_id` is omitted the zone is created in the `init` state and it does not show up in the console until a plan is bound to it.

~> **NOTE:** `allow_duplicates` and `jump_start` are create-only switches that the TEO `DescribeZones` API never returns. Their configured values are kept as-is in state, so out-of-band changes of these two attributes cannot be detected. Changing either of them replaces the zone.

## Example Usage

```hcl
resource "tencentcloud_teo_zone_1110_01" "example" {
  type             = "partial"
  zone_name        = "example.com"
  area             = "overseas"
  plan_id          = "edgeone-2wdo315m2y4c"
  alias_zone_name  = "example-alias"
  allow_duplicates = false
  jump_start       = false

  tags {
    tag_key   = "env"
    tag_value = "test"
  }
}
```

### ### Zone without a domain

```hcl
resource "tencentcloud_teo_zone_1110_01" "no_domain" {
  type       = "noDomainAccess"
  plan_id    = "edgeone-2wdo315m2y4c"
  jump_start = true
}
```

## Argument Reference

The following arguments are supported:

* `alias_zone_name` - (Optional, String) Alias of the zone name, used to distinguish zones that share the same zone name. Up to 200 characters, only digits, letters, `.`, `-` and `_` are allowed. Leave it empty when this feature is not needed. Removing the attribute from the configuration is rejected, because the API does not document whether an omitted value keeps or clears the existing alias.
* `allow_duplicates` - (Optional, Bool) Whether duplicating the access of an existing zone is allowed. `true` allows duplicate access, `false` does not. Default value: `false`. This create-only switch is not returned by the cloud API, so the configured value is kept in state.
* `area` - (Optional, String) Acceleration area of the zone when `type` is `partial` or `full`. Valid values: `global`, `mainland`, `overseas`. Default value: `overseas`. It has to be empty when `type` is `noDomainAccess`.
* `jump_start` - (Optional, Bool) Whether the existing DNS record scan of the zone is skipped. Default value: `false`. This create-only switch is not returned by the cloud API, so the configured value is kept in state.
* `plan_id` - (Optional, String) ID of the plan to be bound to the zone. When a plan already exists in the account this parameter binds the zone to it directly. When it is omitted, a zone in the `init` state is created: such a zone is not activated and does not show up in the console until a plan is bound to it. Changing this value replaces the zone.
* `tags` - (Optional, List) Tags of the zone, used for permission and billing management. The tags have to be created in the tag console before they can be used here.
* `type` - (Optional, String) Access type of the zone. Valid values: `partial` (CNAME access), `full` (NS access), `noDomainAccess` (access without a domain), `dnsPodAccess` (DNSPod hosted access, the domain has to be hosted in DNSPod already) and `ai` (edge inference access). Default value: `partial`. `noDomainAccess` is a create-only value: an existing zone cannot be switched back to it. When it is used, both `zone_name` and `area` have to be empty.
* `zone_name` - (Optional, String) Zone name. Pass the second-level domain (for example `example.com`) for CNAME/NS access and leave it empty for access without a domain. The name can only be changed when the zone is switched from access without a domain to CNAME access, so changing it in any other case replaces the zone.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `content` - Content that has to be written into the file located at `path`.
    * `name_servers` - Name servers assigned to the zone. The name servers of the domain have to be switched to these addresses.
    * `path` - URL path of the verification file, which is fetched by the EdgeOne backend via `http://{Host}{Path}`.
    * `record_type` - Record type.
    * `record_value` - Record value.
    * `subdomain` - Host record.
  * `dns_verification` - Information needed to verify the ownership by adding a DNS record.
  * `file_verification` - Information needed to verify the ownership by uploading a file.
  * `ns_verification` - Information needed to verify the ownership by switching the name servers.
* `id` - ID of the resource. It is identical to `zone_id`.
* `ownership_verification` - Ownership verification information returned by CreateZone. The zone only serves traffic after the ownership verification is finished. It is a creation-time output: `terraform import` leaves it empty.
* `zone_id` - Zone ID, for example `zone-2noz78a8ev6k`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateZone](https://www.tencentcloud.com/document/api/1145/50487)
  - [DescribeZones](https://www.tencentcloud.com/document/api/1145/50481)
* Read: [DescribeZones](https://www.tencentcloud.com/document/api/1145/50481)
* Update:
  - [DescribeZones](https://www.tencentcloud.com/document/api/1145/50481)
  - [ModifyZone](https://www.tencentcloud.com/document/api/1145/50477)
* Delete: [DeleteZone](https://www.tencentcloud.com/document/api/1145/50485)

## Import

TEO zone can be imported using the zone ID, e.g.

```
terraform import tencentcloud_teo_zone_1110_01.example zone-2noz78a8ev6k
```

~> **NOTE:** `ownership_verification` is a creation-time output, so it stays empty on an imported zone.

