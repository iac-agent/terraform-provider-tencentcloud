Provides a resource to create a TEO edge function replica.

-> **Note:** `CreateFunctionReplica` does not return a resource-level id, so the resource id is the composite key `<zone_id>#<function_id>#<replica_name>`. Changing `zone_id`, `function_id` or `replica_name` forces resource recreation.

-> **Note:** `DescribeFunctionReplicas` does not return `ZoneId`, so drift of `zone_id` cannot be detected from the cloud side; the value configured in Terraform is kept in state.

-> **Note:** Deleting this resource only deletes the function replica itself via `DeleteFunctionReplica`; the host edge function (FunctionId) is not deleted. The effect of deleting a replica on live traffic is not documented by the API.

-> **Note:** Removing `remark` from the configuration is rejected: whether `ModifyFunctionReplica` treats an omitted `Remark` as "keep" or "clear" is not documented yet. Keep `remark` set, or destroy and recreate the replica to change it.

Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v6" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

Import

TEO function replica can be imported using the `zone_id#function_id#replica_name`, e.g.

```
terraform import tencentcloud_teo_function_replica_v6.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
