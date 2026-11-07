Provides a resource to create a TEO edge function replica

~> **NOTE:** A replica carries an independent copy of the edge function code. Destroying this resource only removes the replica it manages; the function itself and the other replicas under the same function are left untouched.

~> **NOTE:** `replica_name` is part of the resource identity (the id is `<zone_id>#<function_id>#<replica_name>`), so changing any of `zone_id`, `function_id` or `replica_name` creates a new replica and removes the old one. Removing `remark` from the configuration keeps the current description on the cloud side; set it to an empty string to clear it.

Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v4" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

Import

TEO edge function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v4.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
