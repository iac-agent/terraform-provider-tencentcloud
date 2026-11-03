---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v6"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v6"
description: |-
  Provides a resource to create a TEO edge function replica.
---

# tencentcloud_teo_function_replica_v6

Provides a resource to create a TEO edge function replica.

-> **Note:** `CreateFunctionReplica` does not return a resource-level id, so the resource id is the composite key `<zone_id>#<function_id>#<replica_name>`. Changing `zone_id`, `function_id` or `replica_name` forces resource recreation.

-> **Note:** `DescribeFunctionReplicas` does not return `ZoneId`, so drift of `zone_id` cannot be detected from the cloud side; the value configured in Terraform is kept in state.

-> **Note:** Deleting this resource only deletes the function replica itself via `DeleteFunctionReplica`; the host edge function (FunctionId) is not deleted. The effect of deleting a replica on live traffic is not documented by the API.

-> **Note:** Removing `remark` from the configuration is rejected: whether `ModifyFunctionReplica` treats an omitted `Remark` as "keep" or "clear" is not documented yet. Keep `remark` set, or destroy and recreate the replica to change it.

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v6" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only JavaScript code is supported, with a maximum size of 5MB.
* `function_id` - (Required, String) Function ID.
* `replica_name` - (Required, String) Edge function replica name. Length should be 1-50 characters, allowed characters are a-z, 0-9 and -, and - cannot be used alone or consecutively, nor placed at the beginning or the end. Replica names must be unique under the same FunctionId.
* `zone_id` - (Required, String) Site ID.
* `remark` - (Optional, String) Edge function replica description. Maximum 50 characters.

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
* Delete:
  - [DeleteFunctionReplica](https://www.tencentcloud.com/document/api/1145/80476)
  - [DescribeFunctionReplicas](https://www.tencentcloud.com/document/api/1145/80475)

## Import

TEO function replica can be imported using the `zone_id#function_id#replica_name`, e.g.

```
terraform import tencentcloud_teo_function_replica_v6.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

