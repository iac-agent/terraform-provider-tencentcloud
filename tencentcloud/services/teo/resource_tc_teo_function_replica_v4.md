Provides a resource to create a TEO edge function replica (v4)

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

NOTE: `tencentcloud_teo_function_replica_v4` and the legacy `tencentcloud_teo_function_replica` resource manage the same cloud object. Do not declare both resources for the same `zone_id`/`function_id`/`replica_name` combination in the same configuration, otherwise they will conflict with each other.

Import

TEO function replica v4 can be imported using the composite id `zone_id#function_id#replica_name`, e.g.

```
terraform import tencentcloud_teo_function_replica_v4.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
