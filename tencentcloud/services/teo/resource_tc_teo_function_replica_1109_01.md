Provides a resource to create a TEO edge function replica.

Example Usage

```hcl
resource "tencentcloud_teo_function_replica_1109_01" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

-> **Note:** `zone_id` is not echoed back by the DescribeFunctionReplicas API. The value from the state is kept on every refresh, so drift detection for this field is limited: an out-of-band change of the replica's parent zone is not detected until the replica itself disappears from the query result.

-> **Note:** Deleting this resource only removes the specified function replica itself; the host edge function is not affected. Deleting the host function does not delete the replica (cascade behavior is undocumented).

Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_1109_01.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
