Use this data source to query the deployment records of a TEO inference service.

Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_03" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  service_id  = "svc-1a2b3c4d"
  sort_by     = "create-time"
  sort_order  = "desc"
}
```