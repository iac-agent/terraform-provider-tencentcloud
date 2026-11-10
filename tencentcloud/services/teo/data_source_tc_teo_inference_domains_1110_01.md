# data_source_tc_teo_inference_domains_1110_01

Use this data source to query the list of TEO (EdgeOne) inference service domains.

## Example Usage

```hcl
data "tencentcloud_teo_inference_domains_1110_01" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-2rugn89bkla9"
  sort_by    = "CreateTime"
  sort_order = "Desc"
}
```
