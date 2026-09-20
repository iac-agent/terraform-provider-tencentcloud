---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v3"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v3"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica_v3

Provides a resource to create a TEO edge function replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v3" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
  replica_names = [
    "replica-example",
  ]

  sort_by    = "created-on"
  sort_order = "asc"

  filters {
    name   = "replica-name"
    values = ["replica-example"]
    fuzzy  = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only supports JavaScript code, maximum 5MB.
* `function_id` - (Required, String, ForceNew) Function ID.
* `replica_name` - (Required, String, ForceNew) Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9, -, and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique under the same FunctionId.
* `replica_names` - (Required, List: [`String`]) List of edge function replica names to be deleted.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `filters` - (Optional, List) Filter conditions for querying the function replica list. The maximum number of Values is 20. Supported filter key: replica-name (filter by replica name, fuzzy query supported).
* `remark` - (Optional, String) Edge function replica description. Maximum 50 characters.
* `sort_by` - (Optional, String) Sort by field for querying the function replica list. Valid value: created-on (creation time). Default: created-on.
* `sort_order` - (Optional, String) Sort order for querying the function replica list. Valid values: asc, desc. Default: asc.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, List) Filter values of the field.
* `fuzzy` - (Optional, Bool) Indicates whether fuzzy query is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time of the edge function replica.
* `function_replicas` - List of the edge function replicas returned by the DescribeFunctionReplicas API.
  * `content` - Edge function replica content (JavaScript).
  * `created_on` - Creation time of the edge function replica.
  * `function_id` - Function ID.
  * `modified_on` - Last modification time of the edge function replica.
  * `remark` - Edge function replica description.
  * `replica_name` - Edge function replica name.
* `modified_on` - Last modification time of the edge function replica.


## Import

TEO function replica v3 can be imported using the composite ID with the format `zone_id#function_id#replica_name`, e.g.

```
terraform import tencentcloud_teo_function_replica_v3.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

