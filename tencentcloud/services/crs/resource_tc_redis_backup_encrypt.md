Provides a resource to manage the backup encryption configuration of an existing TencentDB for Redis instance.

~> **NOTE:** Destroying this resource disables backup encryption on the cloud side (`backup_encrypt = false`). It never deletes or otherwise modifies the host Redis instance.

Example Usage

Enable backup encryption with the Tencent Cloud managed key

```hcl
resource "tencentcloud_redis_backup_encrypt" "managed" {
  instance_id    = tencentcloud_redis_instance.example.id
  backup_encrypt = true
}
```

Enable backup encryption with a user-defined KMS key

```hcl
resource "tencentcloud_redis_backup_encrypt" "custom" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = true
  encrypt_key_type = 1
  kms_key_id       = "8f6a1c2e-4b3d-4f5a-9c8e-2d1b0a9f8e7d"
  kms_region       = "ap-guangzhou"
}
```

Disable backup encryption

```hcl
resource "tencentcloud_redis_backup_encrypt" "disabled" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = false
  encrypt_key_type = 0
}
```

Import

Redis backup encryption configuration can be imported using the host instance ID, e.g.

```
$ terraform import tencentcloud_redis_backup_encrypt.example crs-cqdfdzvt
```
