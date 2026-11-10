---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_1110_02"
sidebar_current: "docs-tencentcloud-resource-teo_zone_1110_02"
description: |-
  Provides a resource to create a zone in TencentCloud EdgeOne (TEO).
---

# tencentcloud_teo_zone_1110_02

Provides a resource to create a zone in TencentCloud EdgeOne (TEO).

## Example Usage

```hcl
resource "tencentcloud_teo_zone_1110_02" "example" {
  zone_name = "example.com"
  type      = "full"
  area      = "overseas"
  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}

resource "tencentcloud_teo_zone_1110_02" "example_with_plan" {
  zone_name = "example2.com"
  type      = "partial"
  area      = "overseas"
  plan_id   = "plan-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `alias_zone_name` - (Optional, String) Alias zone name for zones sharing the same name. Only digits, letters, `.`, `-` and `_` are allowed, up to 200 characters.
* `allow_duplicates` - (Optional, Bool) Whether to allow duplicate access. Valid values: true, false (default).
* `area` - (Optional, String) Acceleration area of layer-7 domains when Type is partial/full. Valid values: global, mainland, overseas (default). Leave empty for noDomainAccess access.
* `jump_start` - (Optional, Bool) Whether to skip scanning existing DNS records of the zone. Default value: false.
* `plan_id` - (Optional, String) ID of the plan to bind. If omitted, the zone is created in the `init` (inactive) state; bind a plan via BindZoneToPlan to activate it.
* `tags` - (Optional, List) Tags used for permission control and billing allocation of the zone.
* `type` - (Optional, String) Zone access type. Valid values: partial (CNAME access, default), full (NS access), noDomainAccess, dnsPodAccess (requires the domain hosted in DNSPod), ai (edge inference access).
* `vanity_name_servers` - (Optional, Object) Custom NS information to replace the system-assigned name servers. Not allowed for noDomainAccess zones.
* `zone_name` - (Optional, String) Zone name. For CNAME/NS access, use the second-level domain (example.com). Leave empty for noDomainAccess access.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `vanity_name_servers` object supports the following:

* `switch` - (Required, String) Custom NS switch. Valid values: on, off.
* `servers` - (Optional, Set) List of custom NS servers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

        * `content` - Content of the verification file.
        * `content` - Content of the verification file.
        * `name_servers` - DNS servers assigned by EdgeOne for NS access; switch the domain NameServer to these addresses.
        * `name_servers` - DNS servers assigned by EdgeOne for NS access; switch the domain NameServer to these addresses.
        * `path` - URL path of the verification file.
        * `path` - URL path of the verification file.
        * `record_type` - Record type.
        * `record_type` - Record type.
        * `record_value` - Record value.
        * `record_value` - Record value.
        * `subdomain` - Host record.
        * `subdomain` - Host record.
      * `content` - Content of the verification file.
      * `dns_verification` - Information required for DNS resolution verification.
      * `dns_verification` - Information required for DNS resolution verification.
      * `file_verification` - Information required for file verification.
      * `file_verification` - Information required for file verification.
      * `i_pv4` - Custom name server IPv4 address.
      * `instance_id` - Layer-4 proxy instance ID mapped by the inquiry parameter.
      * `key` - Inquiry parameter key.
      * `name_servers` - DNS servers assigned by EdgeOne for NS access; switch the domain NameServer to these addresses.
      * `name` - Custom name server name.
      * `ns_verification` - Information required for switching DNS servers in NS access mode.
      * `ns_verification` - Information required for switching DNS servers in NS access mode.
      * `pack` - Quota mapped by the inquiry parameter. Valid values: zone, custom-rule, rate-limiting-rule, l4-proxy-instance.
      * `path` - URL path of the verification file.
      * `protection_specs` - Protection level mapped by the inquiry parameter.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `servers` - List of custom NS servers.
      * `subdomain` - Host record.
      * `switch` - Custom NS switch. Valid values: on, off.
      * `value` - Inquiry parameter value.
    * `area` - Region. Valid values: mainland, overseas, global.
    * `auto_renew_flag` - Whether auto renewal is enabled. Valid values: 0 (default), 1 (auto renew), 2 (no auto renew).
    * `cname_speed_up` - Whether CNAME acceleration is enabled. Valid values: enabled, disabled.
    * `config_group_type` - Configuration group type. Valid values: l7_acceleration, edge_functions, web_security.
    * `content` - Content of the verification file.
    * `create_time` - Creation time.
    * `dns_verification` - Information required for DNS resolution verification.
    * `enable_time` - Effective time.
    * `expire_time` - Expiration time.
    * `file_verification` - Information required for file verification.
    * `group` - Resource type. Valid values: plan, pay-as-you-go, value-added.
    * `i_pv4` - Custom name server IPv4 address.
    * `id` - Resource ID.
    * `is_fake` - Whether a zone with the same name exists. Valid values: 0 (no), 1 (yes).
    * `is_fake` - Whether the zone is a fake zone. Valid values: 0 (no), 1 (yes).
    * `is_fake` - Whether the zone is a fake zone. Valid values: 0 (no), 1 (yes).
    * `name_servers` - DNS servers assigned by EdgeOne for NS access; switch the domain NameServer to these addresses.
    * `name_servers` - NS servers assigned by EdgeOne.
    * `name` - Custom name server name.
    * `ns_verification` - Information required for switching DNS servers in NS access mode.
    * `original_name_servers` - NS servers currently in use detected by EdgeOne.
    * `ownership_verification` - Ownership verification information.
    * `ownership_verification` - Ownership verification information.
    * `path` - URL path of the verification file.
    * `pay_mode` - Billing mode. Valid values: 0 (postpaid).
    * `plan_id` - Associated plan resource ID.
    * `record_type` - Record type.
    * `record_value` - Record value.
    * `servers` - List of custom NS servers.
    * `status` - Plan status. Valid values: normal, isolated, destroyed.
    * `subdomain` - Host record.
    * `sv` - Inquiry parameters.
    * `switch` - Custom NS switch. Valid values: on, off.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
    * `type` - Resource marker type.
    * `vanity_name_servers_ips` - IP address information of custom NS servers.
    * `vanity_name_servers` - Custom NS server information.
    * `work_mode` - Work mode. Valid values: immediate_effect, version_control.
    * `zone_number` - Number of zones bound to the resource.
  * `active_status` - Display status. Valid values: active, inactive, paused.
  * `alias_zone_name` - Alias zone name for zones sharing the same name.
  * `area` - Acceleration area of the zone. Valid values: global, mainland, overseas.
  * `cname_detail` - CNAME access zone details. Returned only when Type = partial.
  * `cname_speed_up` - Whether CNAME acceleration is enabled (legacy field, no longer maintained). Valid values: enabled, disabled.
  * `cname_status` - CNAME access status. Valid values: finished, pending.
  * `created_on` - Zone creation time.
  * `dns_pod_detail` - DNSPod hosted zone details. Returned only when Type = dnsPodAccess.
  * `dns_verification` - Information required for DNS resolution verification.
  * `file_verification` - Information required for file verification.
  * `is_fake` - Whether the zone is a fake zone (legacy field, no longer maintained). Valid values: 0 (no), 1 (yes).
  * `lock_status` - Lock status. Valid values: enable, disable, plan_migrate.
  * `modified_on` - Zone modification time.
  * `name_servers` - NS servers assigned by Tencent Cloud (legacy field, no longer maintained).
  * `ns_detail` - NS access zone details. Returned only when Type = full.
  * `ns_verification` - Information required for switching DNS servers in NS access mode.
  * `original_name_servers` - NS servers currently used by the zone (legacy field, no longer maintained).
  * `ownership_verification` - Ownership verification information (legacy field, no longer maintained).
  * `paused` - Whether the zone is paused.
  * `resources` - Billing resource list.
  * `status` - Zone status. Valid values: active, pending, moved, deactivated, initializing.
  * `tags` - Tags associated with the zone.
  * `type` - Zone access type. Valid values: full, partial, noDomainAccess, dnsPodAccess, pages, ai.
  * `vanity_name_servers_ips` - Custom NS IP information (legacy field, no longer maintained).
  * `vanity_name_servers` - Custom NS information (legacy field, no longer maintained).
  * `work_mode_infos` - Work mode of configuration groups for version management.
  * `zone_id` - Zone ID.
  * `zone_name` - Zone name.
* `id` - ID of the zone in the form `zone_id`.
* `ownership_verification` - Zone ownership verification information returned by CreateZone.
* `zones` - Zone detail list returned by DescribeZones. Always contains one element for this resource.

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

teo zone can be imported using the zone id, e.g.

```
terraform import tencentcloud_teo_zone_1110_02.example zone-2noz78a8ev6k
```

