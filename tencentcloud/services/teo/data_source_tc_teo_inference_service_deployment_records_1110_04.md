# tencentcloud_teo_inference_service_deployment_records_1110_04

Use this data source to query the deployment history records of a TEO (TencentCloud EdgeOne) inference service.

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_04" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "service-abc123"
  sort_by    = "create-time"
  sort_order = "desc"
  offset     = 0
  limit      = 20
}

output "deployment_record_ids" {
  value = [for record in data.tencentcloud_teo_inference_service_deployment_records_1110_04.example.record_set : record.record_id]
}
```
