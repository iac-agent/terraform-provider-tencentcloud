---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_encrypt"
sidebar_current: "docs-tencentcloud-resource-redis_backup_encrypt"
description: |-
  Provides a resource to manage the backup encryption configuration of an existing TencentDB for Redis instance.
---

# tencentcloud_redis_backup_encrypt

Provides a resource to manage the backup encryption configuration of an existing TencentDB for Redis instance.

~> **NOTE:** Destroying this resource disables backup encryption on the cloud side (`backup_encrypt = false`). It never deletes or otherwise modifies the host Redis instance.

## Example Usage

### Enable backup encryption with the Tencent Cloud managed key

```hcl
resource "tencentcloud_redis_backup_encrypt" "managed" {
  instance_id    = tencentcloud_redis_instance.example.id
  backup_encrypt = true
}
```

### Enable backup encryption with a user-defined KMS key

```hcl
resource "tencentcloud_redis_backup_encrypt" "custom" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = true
  encrypt_key_type = 1
  kms_key_id       = "8f6a1c2e-4b3d-4f5a-9c8e-2d1b0a9f8e7d"
  kms_region       = "ap-guangzhou"
}
```

### Disable backup encryption

```hcl
resource "tencentcloud_redis_backup_encrypt" "disabled" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = false
  encrypt_key_type = 0
}
```

## Argument Reference

The following arguments are supported:

* `backup_encrypt` - (Required, Bool) Whether to enable backup encryption. `true`: backups are encrypted; `false`: backups are not encrypted. Destroying this resource sets it to `false`.
* `instance_id` - (Required, String) ID of the Redis instance whose backup encryption is managed. Changing it forces a new resource.
* `encrypt_key_type` - (Optional, Int) Encryption key type. `0`: Tencent Cloud managed key (cloud default); `1`: user-defined KMS key. Only meaningful when `backup_encrypt` is `true`.
* `kms_key_id` - (Optional, String) ID of the KMS key, format `xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx`. Only meaningful when `encrypt_key_type` is `1`.
* `kms_region` - (Optional, String) Region of the KMS key, format `ap-xxxx`. Only meaningful when `encrypt_key_type` is `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the backup encryption configuration, same as the host instance ID.


## Import

Redis backup encryption configuration can be imported using the host instance ID, e.g.

```
$ terraform import tencentcloud_redis_backup_encrypt.example crs-cqdfdzvt
```

