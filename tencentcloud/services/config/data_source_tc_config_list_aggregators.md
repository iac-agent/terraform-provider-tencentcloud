Use this data source to query the list of TencentCloud Config (配置审计) account groups (aggregators).

Example Usage

```hcl
data "tencentcloud_config_list_aggregators" "example" {}

output "total" {
  value = data.tencentcloud_config_list_aggregators.example.total
}

output "items" {
  value = data.tencentcloud_config_list_aggregators.example.items
}
```
