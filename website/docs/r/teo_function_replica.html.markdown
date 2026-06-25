---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica"
description: |-
  Provides a resource to create a TEO function_replica
---

# tencentcloud_teo_function_replica

Provides a resource to create a TEO function_replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica" "function_replica" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "func-xxxxxx"
  replica_name = "my-replica"
  content      = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World from replica!!');
          e.respondWith(response);
        });
    EOT
  remark       = "test replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Function replica content, currently only supports JavaScript code, with a maximum size of 5MB.
* `function_id` - (Required, String, ForceNew) ID of the Function.
* `replica_name` - (Required, String, ForceNew) The name of the function replica. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 50 characters.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `remark` - (Optional, String) Function replica description, maximum support of 50 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `modified_on` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.


## Import

teo function_replica can be imported using the composite id `zone_id#function_id#replica_name`, e.g.

```
terraform import tencentcloud_teo_function_replica.function_replica zone-2qtuhspy7cr6#func-xxxxxx#my-replica
```

