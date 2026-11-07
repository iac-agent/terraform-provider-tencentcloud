Provides a resource to create a TEO edge function replica.

The resource manages exactly one replica of one TEO edge function. Deleting the resource only removes
that single replica; the parent function and the sibling replicas of the same function are not affected.

~> **NOTE:** The `content` attribute is stored in the Terraform state file. Protect the state file
accordingly when the function code is sensitive.

~> **NOTE:** `zone_id`, `function_id` and `replica_name` are the identity of this resource and cannot be
changed in place. Modifying any of them destroys the current replica and creates a new one.

Example Usage

```hcl
resource "tencentcloud_teo_function_replica_v3" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica_v3.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
