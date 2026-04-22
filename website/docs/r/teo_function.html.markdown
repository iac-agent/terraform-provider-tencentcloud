---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function"
sidebar_current: "docs-tencentcloud-resource-teo_function"
description: |-
  Provides a resource to create a TEO (EdgeOne) function
---

# tencentcloud_teo_function

Provides a resource to create a TEO (EdgeOne) function

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

### With filters

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

  filters {
    name   = "name"
    values = ["test-function"]
  }

  filters {
    name   = "remark"
    values = ["test"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Function content, currently only supports JavaScript code, with a maximum size of 5MB.
* `name` - (Required, String) Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `filters` - (Optional, List) Filter conditions for querying functions. Supported filter names: `name` (fuzzy match by function name), `remark` (fuzzy match by function description).
* `remark` - (Optional, String) Function description, maximum support of 60 characters.

The `filters` object supports the following:

* `name` - (Required, String) Field name to filter on. Valid values: `name`, `remark`.
* `values` - (Required, List) Filter values for the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `domain` - The default domain name for the function.
* `function_id` - ID of the Function.
* `update_time` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.


## Import

TEO (EdgeOne) function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```

