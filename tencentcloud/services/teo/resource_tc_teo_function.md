Provides a resource to create a teo teo_function

Example Usage

```hcl
resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
```

The `domain_compliance_restrictions` computed attribute lists the regional access restrictions of the function default domain for compliance reasons, e.g. when the default domain is inaccessible in some regions due to ICP filing not obtained or government order:

```
domain_compliance_restrictions {
    reason = "ICP_RECORD_REQUIRED"
    region = "CN"
}
```

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
