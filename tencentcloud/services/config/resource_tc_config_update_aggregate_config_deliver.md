Provides a resource to manage the Config account group delivery setting.

Example Usage

```hcl
resource "tencentcloud_config_update_aggregate_config_deliver" "example" {
  account_group_id     = "ca-xxxxxxxx"
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_uin          = 0
  deliver_content_type = 3
}
```

Import

Config update aggregate config deliver can be imported using the account group id, e.g.

```
$ terraform import tencentcloud_config_update_aggregate_config_deliver.example ca-xxxxxxxx
```
