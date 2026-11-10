Use this data source to query the TEO inference service deployment record list.

Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_01" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "isvc-2qtuhspy7cr6"
  sort_by    = "create-time"
  sort_order = "desc"
}
```
