---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v4"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v4"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica_v4

Provides a resource to create a TEO edge function replica

~> **NOTE:** A replica carries an independent copy of the edge function code. Destroying this resource only removes the replica it manages; the function itself and the other replicas under the same function are left untouched.

~> **NOTE:** `replica_name` is part of the resource identity (the id is `<zone_id>#<function_id>#<replica_name>`), so changing any of `zone_id`, `function_id` or `replica_name` creates a new replica and removes the old one. Removing `remark` from the configuration keeps the current description on the cloud side; set it to an empty string to clear it.

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v4" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only JavaScript code is supported, up to 5 MB. Changes are applied in place.
* `function_id` - (Required, String) Function ID. The replica is created under this function and cannot be moved to another function in place.
* `replica_name` - (Required, String) Edge function replica name. It must be 1-50 characters long and may contain only lowercase letters, digits and `-`; `-` cannot be used alone or consecutively, and cannot be placed at the beginning or end. The name must be unique under the same function. It identifies the replica, so changing it creates a new replica and removes the old one.
* `zone_id` - (Required, String) Zone ID.
* `remark` - (Optional, String) Edge function replica description, up to 50 characters. Removing this attribute leaves the description on the cloud side unchanged.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource identifier in the form `<zone_id>#<function_id>#<replica_name>`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateFunctionReplica](https://www.tencentcloud.com/document/api/1145/80477)
  - [DescribeFunctionReplicas](https://www.tencentcloud.com/document/api/1145/80475)
* Read: [DescribeFunctionReplicas](https://www.tencentcloud.com/document/api/1145/80475)
* Update:
  - [DescribeFunctionReplicas](https://www.tencentcloud.com/document/api/1145/80475)
  - [ModifyFunctionReplica](https://www.tencentcloud.com/document/api/1145/80474)
* Delete: [DeleteFunctionReplica](https://www.tencentcloud.com/document/api/1145/80476)

## Import

TEO edge function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v4.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

