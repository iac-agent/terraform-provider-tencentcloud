Provides a resource to create a TEO (EdgeOne) inference API token.

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-3fkff38fyw8s"
  name    = "tf-example"
}
```

Import

TEO inference API token can be imported using the composite id zoneId#tokenId, e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-3fkff38fyw8s#token-abcd1234
```
