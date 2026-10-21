## Why

TencentDB for Redis (CRS) supports encrypting backup files with either a Tencent Cloud managed key or a user-defined KMS key, but the provider currently has no way to manage this setting. Users who need backup encryption for compliance reasons must configure it manually in the console, which drifts from their Terraform-managed infrastructure. A dedicated config resource closes this gap so backup encryption becomes declarative and reviewable.

## What Changes

- Add a new `RESOURCE_KIND_CONFIG` resource `tencentcloud_redis_backup_encrypt` (Terraform Plugin Framework implementation) that manages the backup encryption configuration of an **existing** Redis instance.
  - Fields: `instance_id` (host instance ID, required, immutable), `backup_encrypt` (required, whether backup encryption is enabled), `encrypt_key_type` (optional, `0` = Tencent Cloud managed key, `1` = user-defined KMS key), `kms_key_id` (optional, KMS key ID used when `encrypt_key_type = 1`), `kms_region` (optional, KMS region used when `encrypt_key_type = 1`).
  - Create/Update both call `ModifyBackupEncrypt`; Read calls `DescribeBackupEncryptConfig`.
  - Destroy **disables backup encryption** (`ModifyBackupEncrypt` with `BackupEncrypt=false`), because this resource is an enable-switch style config: leaving encryption on after `terraform destroy` would silently keep the feature (and its KMS usage) alive. The host instance is never deleted.
  - State ID equals the host instance ID, importable with `terraform import ... <instance_id>`.
- Register the new resource in `tencentcloud/framework/registry.go` (`resourceFactories`) and in `tencentcloud/provider.md` under the `TencentDB for Redis(crs)` product node.
- Add the sibling documentation file `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md` consumed by `make doc`.
- Add unit tests (`resource_tc_redis_backup_encrypt_config_test.go`) using gomonkey-based cloud API mocks (no Terraform acceptance test suite).

No existing resource schema is modified, so there is no breaking change.

## Capabilities

### New Capabilities

- `redis-backup-encrypt-config`: Declarative management of the backup encryption configuration (encryption switch, key type, KMS key ID and KMS region) of an existing TencentDB for Redis instance, including read-back, update and destroy semantics.

### Modified Capabilities

<!-- None: no existing capability requirements are changing. -->

## Impact

- **New files**
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.go`
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config_test.go`
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md`
- **Modified files**
  - `tencentcloud/framework/registry.go` (one factory entry, one import)
  - `tencentcloud/provider.md` (one resource line under the Redis product node)
- **Cloud APIs used** (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`)
  - `DescribeBackupEncryptConfig` (read)
  - `ModifyBackupEncrypt` (create / update / destroy)
- **Dependencies**: none added; the redis SDK package is already vendored.
