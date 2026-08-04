Provides a resource to create an inference API token for TEO (EdgeOne) zones.

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-12345678"
  name    = "my-token"
}
```