Provides a resource to manage the Tencent Cloud Config log delivery setting.

The delivery setting is a global singleton per account and region. It always exists and cannot be created or
deleted, only read and updated, so this resource adopts the pre-existing setting in place on first apply and
overwrites it with the declared values. Declare exactly one instance of this resource per provider
configuration; use provider aliases if you need to manage the setting of several regions.

~> **NOTE:** Destroying this resource disables Config log delivery by setting `status` to `0` on the cloud
side. This resource is an enable switch, so its deletion turns delivery off; it never deletes the Config
service, nor any COS bucket or CLS topic referenced by `target_arn`.

Example Usage

Deliver to a COS bucket

```hcl
resource "tencentcloud_config_update_config_deliver" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_content_type = 3
}
```

Deliver to a CLS topic

```hcl
resource "tencentcloud_config_update_config_deliver" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cls:ap-guangzhou:uin/100000005287:cls/8f3c0f9a-1111-2222-3333-444455556666"
  deliver_prefix       = "config"
  deliver_type         = "CLS"
  deliver_content_type = 1
}
```

Enable only the switch

```hcl
resource "tencentcloud_config_update_config_deliver" "example" {
  status = 0
}
```

Import

Config update config deliver can be imported using the fixed singleton ID `config_update_config_deliver`, e.g.

```
terraform import tencentcloud_config_update_config_deliver.example config_update_config_deliver
```
