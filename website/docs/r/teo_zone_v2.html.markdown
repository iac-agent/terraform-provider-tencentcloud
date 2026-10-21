---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_v2"
sidebar_current: "docs-tencentcloud-resource-teo_zone_v2"
description: |-
  Provides a resource to create a TEO (EdgeOne) zone
---

# tencentcloud_teo_zone_v2

Provides a resource to create a TEO (EdgeOne) zone

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_zone_v2" "zone" {
  type             = "partial"
  zone_name        = "tf-teo.com"
  area             = "global"
  plan_id          = "edgeone-2kfv1h391n6w"
  alias_zone_name  = "teo-test"
  allow_duplicates = false
  jump_start       = false

  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

### Use Custom Name Servers

```hcl
resource "tencentcloud_teo_zone_v2" "zone" {
  type            = "full"
  zone_name       = "tf-teo-vanity.com"
  area            = "global"
  plan_id         = "edgeone-2kfv1h391n6w"
  alias_zone_name = "teo-vanity-test"

  vanity_name_servers {
    switch  = "on"
    servers = ["ns1.example.com", "ns2.example.com"]
  }

  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias_zone_name` - (Optional, String) Alias zone name, which can contain up to 200 characters, including digits, letters, `.`, `- and `_`. Leave it empty if not needed.
* `allow_duplicates` - (Optional, Bool, ForceNew) Indicates whether to allow repeated access. Default value: `false`.
* `area` - (Optional, String) Acceleration area of the site when Type is partial/full. Valid values: `global`, `mainland`, `overseas`. Default value: `overseas`. Leave it empty when Type is noDomainAccess.
* `jump_start` - (Optional, Bool, ForceNew) Indicates whether to skip scanning the existing DNS records of the site. Default value: `false`.
* `plan_id` - (Optional, String, ForceNew) ID of the target plan to be bound. If this parameter is left empty, a site in init status (not activated) will be created, which will not be displayed in the console until it is bound to a plan via BindZoneToPlan.
* `tags` - (Optional, List) Tag list of the site. The tags need to be created in the tag console in advance.
* `type` - (Optional, String) Site access type. Valid values: `partial` (CNAME access), `full` (NS access), `noDomainAccess` (no-domain access), `dnsPodAccess` (DNSPod hosted access), `ai` (edge inference access). Default value: `partial`.
* `vanity_name_servers` - (Optional, List) Custom name servers of the site, which replace the default name servers assigned by the system. Only effective when the site is not in no-domain access mode.
* `zone_name` - (Optional, String) Site name. For CNAME/NS access, pass in a second-level domain (example.com) as the site name; for no-domain access, leave this value empty.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `vanity_name_servers` object supports the following:

* `switch` - (Required, String) Whether to enable custom name servers. Valid values: `on`, `off`.
* `servers` - (Optional, List) Custom name server list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `active_status` - Display status. Valid values: `active`, `inactive`, `paused`.
* `cname_detail` - Details of the CNAME access type site, which is returned only when Type is partial.
  * `is_fake` - Whether it is a fake site. Valid values: `0` (no), `1` (yes).
  * `ownership_verification` - Ownership verification info of the CNAME access type site.
    * `dns_verification` - Info required for DNS resolution verification.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `subdomain` - Host record.
    * `file_verification` - Info required for file verification.
      * `content` - Content of the verification file.
      * `path` - URL Path part of the verification file.
    * `ns_verification` - Info required to switch DNS servers for NS access.
      * `name_servers` - DNS server address assigned for NS access.
* `cname_status` - CNAME access status. Valid values: `finished`, `pending`.
* `created_on` - Site creation time.
* `dns_pod_detail` - Details of the DNSPod hosted access type site, which is returned only when Type is dnsPodAccess.
  * `is_fake` - Whether it is a fake site. Valid values: `0` (no), `1` (yes).
* `lock_status` - Lock status. Valid values: `enable`, `disable`, `plan_migrate`.
* `modified_on` - Site modification time.
* `name_servers` - List of NS servers assigned by Tencent Cloud.
* `ns_detail` - Details of the NS access type site, which is returned only when Type is full.
  * `cname_speed_up` - Whether to enable CNAME acceleration. Valid values: `enabled`, `disabled`.
  * `is_fake` - Whether a site with the same name exists. Valid values: `0` (no), `1` (yes).
  * `name_servers` - List of NS servers assigned by Tencent Cloud EdgeOne.
  * `original_name_servers` - List of NS servers currently used by the site.
  * `ownership_verification` - Ownership verification info of the NS access type site.
    * `dns_verification` - Info required for DNS resolution verification.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `subdomain` - Host record.
    * `file_verification` - Info required for file verification.
      * `content` - Content of the verification file.
      * `path` - URL Path part of the verification file.
    * `ns_verification` - Info required to switch DNS servers for NS access.
      * `name_servers` - DNS server address assigned for NS access.
  * `vanity_name_servers_ips` - IP address info of custom name servers of the user.
    * `i_pv4` - Custom name server IPv4 address.
    * `name` - Custom name server name.
  * `vanity_name_servers` - Custom name server domain info of the user.
    * `servers` - Custom name server list.
    * `switch` - Whether to enable custom name servers.
* `original_name_servers` - List of NS servers currently used by the site.
* `ownership_verification` - Site ownership verification info. After the site is created, ownership verification must be completed before the site can serve properly.
  * `dns_verification` - Info required for DNS resolution verification.
    * `record_type` - Record type.
    * `record_value` - Record value.
    * `subdomain` - Host record.
  * `file_verification` - Info required for file verification.
    * `content` - Content of the verification file.
    * `path` - URL Path part of the verification file.
  * `ns_verification` - Info required to switch DNS servers for NS access.
    * `name_servers` - DNS server address assigned for NS access.
* `paused` - Whether the site is closed.
* `resources` - Billing resource list of the site.
  * `area` - Region. Valid values: `mainland`, `overseas`, `global`.
  * `auto_renew_flag` - Auto renewal flag. Valid values: `0` (default), `1` (auto-renewal), `2` (no auto-renewal).
  * `create_time` - Creation time.
  * `enable_time` - Effective time.
  * `expire_time` - Expiration time.
  * `group` - Resource type.
  * `id` - Resource ID.
  * `pay_mode` - Billing mode. Valid values: `0` (pay-as-you-go).
  * `plan_id` - Plan-associated resource ID.
  * `status` - Plan status. Valid values: `normal`, `isolated`, `destroyed`.
  * `sv` - Inquiry parameters.
    * `instance_id` - Layer-4 proxy instance ID mapped by the inquiry parameter.
    * `key` - Inquiry parameter key.
    * `pack` - Quota mapped by the inquiry parameter.
    * `protection_specs` - Protection level of the inquiry parameter.
    * `value` - Inquiry parameter value.
  * `type` - Resource tag type.
  * `zone_number` - Number of sites bound to the current resource.
* `status` - Site status. Valid values: `active`, `pending`, `moved`, `deactivated`, `initializing`.
* `vanity_name_servers_ips` - IP address info of custom name servers of the user.
  * `i_pv4` - Custom name server IPv4 address.
  * `name` - Custom name server name.
* `work_mode_infos` - Version management config group working mode of the site.
  * `config_group_type` - Config group type. Valid values: `l7_acceleration`, `edge_functions`, `web_security`.
  * `work_mode` - Working mode. Valid values: `immediate_effect`, `version_control`.
* `zone_id` - Site ID.


## Import

TEO (EdgeOne) zone can be imported using the zone id, e.g.
```
terraform import tencentcloud_teo_zone_v2.zone zone-2noz78a8ev6k
```

