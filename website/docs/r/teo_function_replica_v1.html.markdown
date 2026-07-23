---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v1"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v1"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica_v1

Provides a resource to create a TEO edge function replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v1" "replica" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "func-abcdefghij"
  replica_name = "replica-test"
  content      = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
  remark       = "test replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Content of the edge function replica. Currently only supports JavaScript code, with a maximum size of 5MB.
* `function_id` - (Required, String, ForceNew) ID of the edge function.
* `remark` - (Required, String) Description of the edge function replica. Maximum support of 50 characters.
* `replica_name` - (Required, String, ForceNew) Name of the edge function replica. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, with a maximum length of 50 characters. The replica name must be unique under the same FunctionId.
* `zone_id` - (Required, String, ForceNew) ID of the site.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the edge function replica.
* `update_time` - Update time of the edge function replica.


