Use this data source to query the security group quota limits of the current TencentCloud VPC account in the configured region.

~> **NOTE:** This data source only reads quota information, it does not create, modify or delete any security group or other cloud resource. The returned values are the quotas granted to the account, changing them requires a quota adjustment request on the TencentCloud console.

Example Usage

```hcl
data "tencentcloud_vpc_security_group_limits_1110_01" "example" {}

output "security_group_limit" {
  value = data.tencentcloud_vpc_security_group_limits_1110_01.example.security_group_limit_set[0].security_group_limit
}

output "security_group_policy_limit" {
  value = data.tencentcloud_vpc_security_group_limits_1110_01.example.security_group_limit_set[0].security_group_policy_limit
}
```
