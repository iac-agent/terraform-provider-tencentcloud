Provides a resource to create a TEO inference API token

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-xxxxxxxx"
  name    = "my-token"
}
```

Import

TEO inference API token can be imported using the zone ID and token ID, separated by "#", e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-xxxxxxxx#token-xxxxxxxx
```