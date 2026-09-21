---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v5"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v5"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica_v5

Provides a resource to create a TEO edge function replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v5" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
  sort_by      = "created-on"
  sort_order   = "desc"

  filters {
    name   = "replica-name"
    values = ["replica-example"]
    fuzzy  = false
  }
}
```

### :

```hcl
resource "tencentcloud_teo_function_replica_v5" "example" {
  zone_id       = "zone-2qtuhspy7cr6"
  function_id   = "ef-2qlxy8s7o96e"
  replica_name  = "replica-example"
  content       = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  replica_names = ["replica-example", "replica-example-2"]
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only supports JavaScript code, maximum 5MB.
* `function_id` - (Required, String, ForceNew) Function ID.
* `replica_name` - (Required, String, ForceNew) Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9, -, and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique under the same FunctionId.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `filters` - (Optional, List) Filter conditions of the replica list. The maximum value of Filters.Values is 20. If this parameter is not filled in, all function replicas under the function ID will be returned. Valid values: `replica-name`: filter by function replica name, which supports fuzzy query.
* `remark` - (Optional, String) Edge function replica description. Maximum 50 characters.
* `replica_names` - (Optional, List: [`String`]) Names of the function replicas to be deleted when destroying the resource. If not configured, only the replica corresponding to this resource (the `replica_name` in the resource ID) will be deleted by default.
* `sort_by` - (Optional, String) Sort basis of the replica list. Valid values: `created-on`: sort by creation time. Default sorted by the `created-on` attribute.
* `sort_order` - (Optional, String) Sort order of the replica list. Valid values: `asc`: sort in ascending order; `desc`: sort in descending order. Default value: `asc`.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered. Valid value: `replica-name`.
* `values` - (Required, Set) Value of the filtered field. The maximum number of values is 20.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time of the edge function replica.
* `modified_on` - Last modified time of the edge function replica.


## Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v5.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

