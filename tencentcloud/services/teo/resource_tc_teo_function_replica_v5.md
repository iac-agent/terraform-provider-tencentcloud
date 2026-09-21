Provides a resource to create a TEO edge function replica

Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v5" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
  sort_by      = "created-on"
  sort_order   = "desc"

  filters {
    name   = "replica-name"
    values = ["replica-example"]
    fuzzy  = false
  }
}
```

If you want to delete multiple replicas by name when destroying the resource, you can configure `replica_names`:

```hcl
resource "tencentcloud_teo_function_replica_v5" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  replica_names = ["replica-example", "replica-example-2"]
}
```

Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v5.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
