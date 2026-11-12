Use this data source to query detailed information of a Config system preset rule.

Example Usage

Query a system preset rule by identifier

```hcl
data "tencentcloud_config_system_rule" "example" {
  identifier = "cam-user-mfa-enabled"
}
```
