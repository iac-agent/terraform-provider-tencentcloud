---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica_v3"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica_v3"
description: |-
  Provides a resource to create a TEO edge function replica.
---

# tencentcloud_teo_function_replica_v3

Provides a resource to create a TEO edge function replica.

The resource manages exactly one replica of one TEO edge function. Deleting the resource only removes
that single replica; the parent function and the sibling replicas of the same function are not affected.

~> **NOTE:** The `content` attribute is stored in the Terraform state file. Protect the state file
accordingly when the function code is sensitive.

~> **NOTE:** `zone_id`, `function_id` and `replica_name` are the identity of this resource and cannot be
changed in place. Modifying any of them destroys the current replica and creates a new one.

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v3" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Content of the edge function replica. Only JavaScript code is supported, and the maximum size is 5MB. Note that the content is stored in the Terraform state file, so the state file must be protected accordingly.
* `function_id` - (Required, String) ID of the edge function.
* `replica_name` - (Required, String) Name of the edge function replica. 1-50 characters, only `a-z`, `0-9` and `-` are allowed; `-` cannot be used alone, repeated, or placed at the beginning or the end. The name must be unique within the same function.
* `zone_id` - (Required, String) ID of the zone to which the function belongs.
* `remark` - (Optional, String) Description of the edge function replica. Maximum 50 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_on` - Creation time of the edge function replica.
* `id` - Resource identifier in the form `<zone_id>#<function_id>#<replica_name>`.
* `modified_on` - Last modification time of the edge function replica.

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

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v3.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

