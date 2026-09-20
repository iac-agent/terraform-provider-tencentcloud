Provides a resource to manage an EdgeOne (TEO) edge function.

Example Usage

```hcl
resource "tencentcloud_teo_function_v6" "teo_function_v6" {
    zone_id = "zone-2qtuhspy7cr6"
    name    = "test-function"
    remark  = "test"
    content = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
}
```

Import

EdgeOne (TEO) edge function can be imported using the composite id, e.g.

```
terraform import tencentcloud_teo_function_v6.teo_function_v6 zone_id#function_id
```
