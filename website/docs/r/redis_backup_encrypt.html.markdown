---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_encrypt"
sidebar_current: "docs-tencentcloud-resource-redis_backup_encrypt"
description: |-
  Provides a resource to manage the backup encryption config of a TencentDB for Redis (redis) instance.
---

# tencentcloud_redis_backup_encrypt

Provides a resource to manage the backup encryption config of a TencentDB for Redis (redis) instance.

~> NOTE: The backup encryption config only takes effect for new backups. Historical backups generated before enabling encryption are not affected. Destroying this resource only removes it from state and does not disable backup encryption.

## Example Usage

### Enable backup encryption with the TencentDB managed key

```hcl
resource "tencentcloud_redis_instance" "example" {
  type_id            = 8
  availability_zone  = "ap-guangzhou-3"
  name               = "tf_example"
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = 1
  redis_replicas_num = 1
  port               = 6379
}

resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id    = tencentcloud_redis_instance.example.id
  backup_encrypt = true
}
```

### Enable backup encryption with a customer-defined KMS key

```hcl
resource "tencentcloud_kms_key" "example" {
  alias      = "tf_example_kms_key"
  is_enabled = true
}

resource "tencentcloud_redis_instance" "example" {
  type_id            = 8
  availability_zone  = "ap-guangzhou-3"
  name               = "tf_example"
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = 1
  redis_replicas_num = 1
  port               = 6379
}

resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = true
  encrypt_key_type = 1
  kms_key_id       = tencentcloud_kms_key.example.id
  kms_region       = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `backup_encrypt` - (Required, Bool) Indicates whether to enable backup encryption. `true`: enable, `false`: disable.
* `instance_id` - (Required, String, ForceNew) ID of the redis instance.
* `encrypt_key_type` - (Optional, Int) Encryption key type. Valid values: `0`: TencentDB managed key (default), `1`: customer-defined KMS key.
* `kms_key_id` - (Optional, String) KMS key ID, format: `xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx`.
* `kms_region` - (Optional, String) KMS region, format: `ap-xxxx`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Redis backup encryption config can be imported using the id of the redis instance, e.g.

```
$ terraform import tencentcloud_redis_backup_encrypt.example crs-cqnlbzze
```

