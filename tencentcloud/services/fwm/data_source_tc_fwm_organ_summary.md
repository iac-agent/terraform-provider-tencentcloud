Use this data source to query the organisation overview of FWM (Firewall Manager).

Example Usage

```hcl
data "tencentcloud_fwm_organ_summary" "example" {}

output "fwm_group_name" {
  value = data.tencentcloud_fwm_organ_summary.example.summary.group_name
}

output "fwm_joined_member_count" {
  value = data.tencentcloud_fwm_organ_summary.example.summary.joined_member_count
}

output "fwm_departments" {
  value = data.tencentcloud_fwm_organ_summary.example.summary.departments
}
```
