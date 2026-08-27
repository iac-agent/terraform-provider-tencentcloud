Provides a resource to create a TEO Inference API Token

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-2o3h21ed2t68"
  name    = "example-token"

  offset = 0
  limit  = 20
}
```

Import

TEO Inference API Token can be imported using the zone_id#token_id, e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-2o3h21ed2t68#token-xxxxx
```