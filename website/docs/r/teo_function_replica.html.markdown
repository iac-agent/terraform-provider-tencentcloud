---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica

Provides a resource to create a TEO edge function replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qtuhspy7cr6"
  replica_name = "my-replica"
  content      = <<-EOT
    addEventListener('fetch', e => {
      const response = new Response('Hello Replica!!');
      e.respondWith(response);
    });
  EOT
  remark       = "test replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only JavaScript code is supported, with a maximum size of 5MB.
* `function_id` - (Required, String, ForceNew) ID of the edge function.
* `replica_name` - (Required, String, ForceNew) Edge function replica name. It can only contain lowercase letters, numbers, and hyphens. It must start and end with a letter or number, cannot have consecutive hyphens, and has a maximum length of 50 characters. The replica name must be unique under the same FunctionId.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `remark` - (Optional, String) Edge function replica description. Maximum support of 50 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time of the edge function replica.
* `modified_on` - Last modification time of the edge function replica.


## Import

TEO edge function replica can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_replica.example zone-2qtuhspy7cr6#ef-2qtuhspy7cr6#my-replica
```

