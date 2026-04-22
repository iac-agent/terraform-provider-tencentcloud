---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function"
sidebar_current: "docs-tencentcloud-resource-teo_function"
description: |-
  Provides a resource to create a TEO (EdgeOne) edge function
---

# tencentcloud_teo_function

Provides a resource to create a TEO (EdgeOne) edge function

## Example Usage

```hcl
resource "tencentcloud_teo_function" "teo_function" {
  content = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
  name    = "aaa-zone-2qtuhspy7cr6-1310708577"
  remark  = "test"
  zone_id = "zone-2qtuhspy7cr6"
}
```

### Query functions by function IDs

```hcl
resource "tencentcloud_teo_function" "teo_function" {
  zone_id      = "zone-2qtuhspy7cr6"
  name         = "test-function"
  content      = "addEventListener('fetch', e => { e.respondWith(new Response('Hello')); });"
  function_ids = ["func-001", "func-002"]
}
```

### Query functions by filters

```hcl
resource "tencentcloud_teo_function" "teo_function" {
  zone_id = "zone-2qtuhspy7cr6"
  name    = "test-function"
  content = "addEventListener('fetch', e => { e.respondWith(new Response('Hello')); });"
  filters {
    name   = "name"
    values = ["test-func"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Function content, currently only supports JavaScript code, with a maximum size of 5MB.
* `name` - (Required, String) Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `filters` - (Optional, List) Filter conditions for querying functions. Support filtering by `name` (function name fuzzy match) and `remark` (function description fuzzy match).
* `function_ids` - (Optional, List: [`String`]) List of function IDs to filter by.
* `remark` - (Optional, String) Function description, maximum support of 60 characters.

The `filters` object supports the following:

* `name` - (Required, String) Filter name. Valid values: `name`, `remark`.
* `values` - (Optional, Set) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `domain` - The default domain name for the function.
* `function_id` - ID of the Function.
* `functions` - A list of functions matching the query conditions. Each element contains the following attributes:
  * `content` - Function content.
  * `create_time` - Creation time.
  * `domain` - Default domain name for the function.
  * `function_id` - ID of the function.
  * `name` - Function name.
  * `remark` - Function description.
  * `update_time` - Modification time.
  * `zone_id` - ID of the site.
* `update_time` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.


## Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```

