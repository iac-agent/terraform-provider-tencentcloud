Provides a resource to create a TEO inference API token

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token_v7" "example" {
  zone_id = "zone-27q0p0bali16"
  name    = "my-inference-token"
}
```

Import

The resource can be imported by using the `token_id`, e.g.

```sh
terraform import tencentcloud_teo_inference_api_token_v7.example <token_id>
```