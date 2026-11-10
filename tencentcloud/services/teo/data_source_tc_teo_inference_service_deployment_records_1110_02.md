Use this data source to query the deployment history records of a TEO inference service.

Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_02" "deployment_records" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "svc-2qtuhspyeq8r6"
  sort_by    = "create-time"
  sort_order = "desc"
}
```
