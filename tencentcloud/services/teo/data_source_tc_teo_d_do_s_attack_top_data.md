Use this data source to query detailed information of TEO DDoS attack top data

Example Usage

```hcl
data "tencentcloud_teo_d_do_s_attack_top_data" "example" {
  start_time  = "2024-01-01T00:00:00Z"
  end_time    = "2024-01-02T00:00:00Z"
  metric_name = "ddos_attackFlux_protocol"
  zone_ids    = ["zone-2qtuhspy7cr6"]
  area        = "overseas"
}
```