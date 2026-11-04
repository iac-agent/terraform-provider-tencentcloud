Use this data source to query the list of TEO inference service domains under a specified zone and inference service.

Example Usage

```hcl
data "tencentcloud_teo_inference_domains_v2" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-123456"
}
```

Sort by creation time in ascending order:

```hcl
data "tencentcloud_teo_inference_domains_v2" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-123456"
  sort_by    = "CreateTime"
  sort_order = "Asc"
}
```
