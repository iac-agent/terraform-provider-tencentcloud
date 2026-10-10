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

The following computed attribute is exported when the function default domain is restricted in some regions due to compliance issues:

* `domain_compliance_restrictions` - (Optional) Region access restriction list of the function default domain due to compliance issues. Each element contains:
  * `reason` - The reason for the access restriction. Valid values: `ICP_RECORD_REQUIRED` (ICP filing not obtained), `GOVERNMENT_ORDER` (government order).
  * `region` - The country/region code of the restricted access area, using the ISO 3166 country/region code standard (e.g. `CN`).

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
