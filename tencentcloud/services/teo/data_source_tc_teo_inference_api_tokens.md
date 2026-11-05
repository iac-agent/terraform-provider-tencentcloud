Use this data source to query the inference API tokens of the specified TEO zone, including the token ID, name, content and creation time.

Example Usage

```hcl
data "tencentcloud_teo_inference_api_tokens" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```
