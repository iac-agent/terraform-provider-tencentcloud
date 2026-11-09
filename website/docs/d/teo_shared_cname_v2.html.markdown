---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_shared_cname_v2"
sidebar_current: "docs-tencentcloud-data_source-teo_shared_cname_v2"
description: |-
  Use this data source to query the shared CNAME list of TEO (TencentCloud EdgeOne) within a zone.
---

# tencentcloud_teo_shared_cname_v2

Use this data source to query the shared CNAME list of TEO (TencentCloud EdgeOne) within a zone.

## Example Usage

```hcl
data "tencentcloud_teo_shared_cname_v2" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}

# Query with filters and sorting.
data "tencentcloud_teo_shared_cname_v2" "filtered" {
  zone_id = "zone-2qtuhspy7cr6"

  filters {
    name   = "shared-cname"
    values = ["prefix-1.2qtuhspy7cr6.share.dnse1.com"]
  }

  filters {
    name   = "description"
    values = ["demo"]
    fuzzy  = true
  }

  match     = "any"
  direction = "desc"
  order     = "create-time"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID to which the shared CNAMEs belong.
* `direction` - (Optional, String) Sorting direction of the list. Valid values: `asc`, `desc`. Default value: `asc`.
* `filters` - (Optional, List) Filter conditions of the shared CNAME list. The upper limit of Filters.Values is 20. Valid filter names: `shared-cname` (filter by shared CNAME), `type` (filter by shared CNAME type), `description` (filter by description).
* `match` - (Optional, String) Matching mode of multiple query conditions. Valid values: `all` (return shared CNAMEs matching all query conditions), `any` (return shared CNAMEs matching any one of the query conditions). Default value: `all`.
* `order` - (Optional, String) Sorting basis of the list. Valid values: `create-time` (creation time), `shared-cname` (shared CNAME). Default sorting is by the shared-cname attribute.

The `filters` object supports the following:

* `name` - (Required, String) Name of the field to filter by.
* `values` - (Required, Set) Values of the filtered field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `associated_domain` - Domain associated with the IP SSL configuration. Null when Status is `unbound`.
    * `instance` - Information of the referenced/referencing instance.
    * `status` - Association status of the IP SSL configuration. Valid values: `bound` (bound), `binding` (binding), `unbinding` (unbinding), `unbound` (not bound).
    * `type` - Instance type. Valid value: `acceleration-domain` (acceleration domain).
    * `zone_id` - Zone ID of the domain associated with the IP SSL configuration. Null when Status is `unbound`.
    * `zone_id` - Zone ID.
  * `acceleration_domains` - List of acceleration domains joined to the shared CNAME. Only the first 100 acceleration domains are returned when more than 100 domains are joined.
  * `bind_domain_count` - Number of acceleration domains bound to the shared CNAME.
  * `description` - Description of the shared CNAME.
  * `ipssl_config` - IP SSL configuration associated with the shared CNAME, returned only when type is `ip-ssl`. When Status is `unbound`, zone_id and associated_domain are null.
  * `shared_cname` - Name of the shared CNAME.
  * `type` - Type of the shared CNAME. Valid values: `custom` (custom shared CNAME created by user), `ip-ssl` (IP SSL shared CNAME), `zero-rating` (zero-rating shared CNAME), `preset` (preset resource shared CNAME).
* `shared_cname_info` - List of shared CNAMEs.


