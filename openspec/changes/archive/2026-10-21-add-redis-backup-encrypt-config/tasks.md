## 1. Service layer (tencentcloud/services/crs/service_tencentcloud_redis.go)

- [x] 1.1 Add `DescribeBackupEncryptConfig(ctx context.Context, instanceId string) (*redis.DescribeBackupEncryptConfigResponseParams, error)` to `RedisService`, calling `client.UseRedisClient().DescribeBackupEncryptConfigWithContext` and returning a wrapped error on failure
- [x] 1.2 Add `ModifyBackupEncrypt(ctx context.Context, instanceId string, backupEncrypt bool, encryptKeyType *int64, kmsKeyId, kmsRegion *string) error` to `RedisService`
- [x] 1.3 Confirm both helpers log request/response (`log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]")`) and reuse the existing `RedisInstanceNotFound` constant

## 2. Framework resource (tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.go)

- [x] 2.1 Create the file in `package crs` with the `resource.Resource`, `resource.ResourceWithConfigure` and `resource.ResourceWithImportState` interface assertions and the factory `NewCrsRedisBackupEncryptConfigResource()`
- [x] 2.2 Define `CrsRedisBackupEncryptConfigResourceModel` with `Id`, `InstanceId`, `BackupEncrypt`, `EncryptKeyType`, `KmsKeyId`, `KmsRegion` (`tfsdk` tags: `id`, `instance_id`, `backup_encrypt`, `encrypt_key_type`, `kms_key_id`, `kms_region`)
- [x] 2.3 Implement `Metadata` with `resp.TypeName = "tencentcloud_redis_backup_encrypt"`
- [x] 2.4 Implement `Schema`: `id` computed + `UseStateForUnknown`; `instance_id` required + `RequiresReplace`; `backup_encrypt` required bool; `encrypt_key_type` optional+computed int64 with `OneOf(0, 1)` validator; `kms_key_id` and `kms_region` optional+computed string; schema description states the destroy semantics (disables encryption, never deletes the instance)
- [x] 2.5 Implement `Configure` reading `*sharedmeta.ProviderMeta` and storing the shared client / `NewRedisService(meta.Client)`
- [x] 2.6 Implement the private `modifyBackupEncrypt(ctx, m *Model) error` helper: only configured (`!IsNull() && !IsUnknown()`) optional fields are set on `redis.NewModifyBackupEncryptRequest()`, the API call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(err)`
- [x] 2.7 Implement `Create`: call `modifyBackupEncrypt`, then (outside the retry block) set `plan.Id = plan.InstanceId`, `resp.State.Set`, then refresh state via the Read logic ensuring the id is never empty
- [x] 2.8 Implement `Read`: backfill `instance_id` from `id` when null; call `DescribeBackupEncryptConfig` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`; handle `ResourceNotFound.InstanceNotExists` by logging the id and `resp.State.RemoveResource(ctx)`; on nil `response.Response` log `[CRUD] redis_backup_encrypt id=...` and `AddError`; write fields back with `types.XxxPointerValue(...)`
- [x] 2.9 Implement `Update`: call `modifyBackupEncrypt` with the plan, then `resp.State.Set` and refresh state
- [x] 2.10 Implement `Delete`: call `ModifyBackupEncrypt` with `BackupEncrypt=false` (log the disable action), never touch the host instance
- [x] 2.11 Implement `ImportState` setting `id` from `req.ID` via `resp.State.SetAttribute(ctx, path.Root("id"), req.ID)`

## 3. Registration

- [x] 3.1 Add the `crs` package import to `tencentcloud/framework/registry.go` and append `crs.NewCrsRedisBackupEncryptConfigResource,` to `resourceFactories`
- [x] 3.2 Append `tencentcloud_redis_backup_encrypt` to the `Resource` section of the `TencentDB for Redis(crs)` node in `tencentcloud/provider.md`
- [x] 3.3 Do not modify `tencentcloud/provider.go` (Framework reference is registered only through `registry.go`)

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md` with a one-line English description starting with "Provides a resource to ...", an `Example Usage` HCL block (instance reference plus KMS key variant), a note that destroy disables backup encryption without deleting the instance, and an `Import` section using the host instance ID
- [x] 4.2 Do not add `Argument Reference` / `Attribute Reference` sections and do not touch `website/` (generated later by `make doc`)

## 5. Unit tests (tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config_test.go)

- [x] 5.1 Add tests using gomonkey to mock `UseRedisClient()` / `DescribeBackupEncryptConfig` / `ModifyBackupEncrypt`, without the Terraform acceptance test suite
- [x] 5.2 Cover `Create` request mapping (all five attributes forwarded) and `id` assignment from `instance_id`
- [x] 5.3 Cover `Read` success state mapping, `ResourceNotFound.InstanceNotExists` removal, and generic error propagation
- [x] 5.4 Cover `Update` (disable encryption and key-type change) and `Delete` calling the API with `BackupEncrypt=false`
- [x] 5.5 Cover `ImportState` and the `instance_id` backfill path
- [x] 5.6 Do not execute `go test`, `go build` or `go vet` in this phase

## 6. Finalize

- [x] 6.1 Verify all returned errors are checked and unused ones assigned to `_`
- [x] 6.2 Hand over to the finalize skill for `gofmt`, `make doc` and changelog generation
