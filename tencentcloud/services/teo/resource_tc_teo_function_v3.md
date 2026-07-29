Provides a resource to create a TEO edge function.

Example Usage

```hcl
resource "tencentcloud_teo_function_v3" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  name    = "test-function"
  content = <<-EOT
    addEventListener('fetch', e => {
      const response = new Response('Hello World!!');
      e.respondWith(response);
    });
  EOT
  remark  = "test remark"
}
```

Import

TEO edge function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_v3.example zone_id#function_id
```
