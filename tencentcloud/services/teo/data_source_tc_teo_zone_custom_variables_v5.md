Use this data source to query the zone-level custom variables of TencentCloud EdgeOne (TEO), including variable definitions and variable operation rules. When no custom variables are configured, both the variable list and the operation rule list are empty.

Example Usage

Query zone-level custom variables by zone ID

```hcl
data "tencentcloud_teo_zone_custom_variables_v5" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

Query zone-level custom variables referenced by a zone resource

```hcl
data "tencentcloud_teo_zone_custom_variables_v5" "example" {
  zone_id = tencentcloud_teo_zone.example.id
}
```
