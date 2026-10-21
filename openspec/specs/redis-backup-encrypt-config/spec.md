# redis-backup-encrypt-config Specification

## Purpose
TBD - created by archiving change add-redis-backup-encrypt-config. Update Purpose after archive.
## Requirements
### Requirement: Manage Redis backup encryption configuration declaratively

The provider SHALL expose a Terraform Plugin Framework resource named
`tencentcloud_redis_backup_encrypt` that manages the backup encryption configuration of an existing
TencentDB for Redis (CRS) instance. The resource SHALL declare the attributes `instance_id`
(required, immutable), `backup_encrypt` (required), `encrypt_key_type` (optional, computed),
`kms_key_id` (optional, computed) and `kms_region` (optional, computed), and a computed `id` equal
to the host instance ID.

#### Scenario: Create enables backup encryption

- **WHEN** a user applies a `tencentcloud_redis_backup_encrypt` resource with
  `instance_id = "crs-xxxx"` and `backup_encrypt = true`
- **THEN** the provider calls `ModifyBackupEncrypt` with `InstanceId`, `BackupEncrypt=true` and any
  configured `EncryptKeyType` / `KmsKeyId` / `KmsRegion`, records `id = "crs-xxxx"` in state, and
  refreshes the state through `DescribeBackupEncryptConfig`

#### Scenario: Create with user-defined KMS key

- **WHEN** the configuration sets `backup_encrypt = true`, `encrypt_key_type = 1`,
  `kms_key_id = "<kms-key-id>"` and `kms_region = "ap-guangzhou"`
- **THEN** the provider forwards all four values in the `ModifyBackupEncrypt` request and state
  mirrors the values returned by `DescribeBackupEncryptConfig`

#### Scenario: Instance ID change forces replacement

- **WHEN** a user changes `instance_id` from `crs-aaaa` to `crs-bbbb`
- **THEN** Terraform plans a destroy of the old config resource and a create against `crs-bbbb`
  instead of an in-place update

### Requirement: Read reflects the cloud-side backup encryption state

The `Read` callback SHALL call `DescribeBackupEncryptConfig` with `tccommon.ReadRetryTimeout`
retry and SHALL write `backup_encrypt`, `encrypt_key_type`, `kms_key_id` and `kms_region` back to
state from the response. When a response field pointer is nil the corresponding state value SHALL
be set to null rather than a zero value.

#### Scenario: Drift detected after out-of-band change

- **WHEN** backup encryption is disabled manually in the console while the resource state still has
  `backup_encrypt = true`
- **THEN** the next refresh/plan updates state to `backup_encrypt = false` and proposes an in-place
  update back to `true`

#### Scenario: Host instance no longer exists

- **WHEN** `DescribeBackupEncryptConfig` returns the error code `ResourceNotFound.InstanceNotExists`
- **THEN** the provider logs the resource id and removes the resource from state without failing

#### Scenario: Empty response payload

- **WHEN** `DescribeBackupEncryptConfig` returns a nil `Response` payload
- **THEN** the provider reports a diagnostic error including the resource id instead of clearing the
  state id silently

#### Scenario: Import by host instance ID

- **WHEN** a user runs `terraform import tencentcloud_redis_backup_encrypt.example crs-xxxx`
- **THEN** the provider stores `crs-xxxx` as the resource id, backfills `instance_id` from it during
  `Read`, and populates all configuration fields from the cloud

### Requirement: Update applies configuration changes through the modify API

The `Update` callback SHALL reuse the same internal helper as `Create` to call
`ModifyBackupEncrypt`, and SHALL persist the planned values to state after a successful call.

#### Scenario: Disable encryption through update

- **WHEN** a user changes `backup_encrypt` from `true` to `false`
- **THEN** the provider calls `ModifyBackupEncrypt` with `BackupEncrypt=false` and stores the new
  value in state

#### Scenario: Change the KMS key type

- **WHEN** a user changes `encrypt_key_type` from `0` to `1` and supplies `kms_key_id` and
  `kms_region`
- **THEN** the provider forwards the new key type, key id and region in `ModifyBackupEncrypt` and
  the resulting state matches the cloud response

### Requirement: Destroy disables backup encryption without deleting the instance

The `Delete` callback SHALL call `ModifyBackupEncrypt` with `BackupEncrypt=false`. It SHALL never
delete or otherwise modify the host Redis instance lifecycle. This behaviour SHALL be documented in
the resource schema descriptions and the resource Markdown documentation.

#### Scenario: Destroy turns encryption off

- **WHEN** a user runs `terraform destroy` on the resource while `backup_encrypt = true`
- **THEN** the provider calls `ModifyBackupEncrypt` with `BackupEncrypt=false` and the Redis
  instance still exists afterwards

#### Scenario: Destroy when encryption already off

- **WHEN** the resource state has `backup_encrypt = false` and the resource is destroyed
- **THEN** the provider still issues the disable call (idempotent on the cloud side) and removes the
  resource from state

### Requirement: Resource registration and documentation

The resource SHALL be registered in `tencentcloud/framework/registry.go` `resourceFactories` via
`crs.NewCrsRedisBackupEncryptConfigResource` and listed in `tencentcloud/provider.md` under the
`TencentDB for Redis(crs)` product node's `Resource` section. A sibling Markdown file
`tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md` SHALL provide the one-line
description, an HCL `Example Usage` and an `Import` section.

#### Scenario: Mux startup without type-name collision

- **WHEN** the provider starts and merges the SDKv2 and Framework servers
- **THEN** `tencentcloud_redis_backup_encrypt` is served by the Framework provider only, with no
  duplicate type-name error

#### Scenario: Documentation generation

- **WHEN** `make doc` runs
- **THEN** the generated documentation for `tencentcloud_redis_backup_encrypt` includes the example
  usage and import instructions, and the destroy semantics note

### Requirement: Unit test coverage via cloud API mocks

The resource SHALL ship a unit test file that mocks the Redis cloud API with gomonkey and does not
use the Terraform acceptance test suite.

#### Scenario: Create request mapping verified

- **WHEN** the unit test invokes `Create` with a plan containing all five attributes
- **THEN** the mocked `ModifyBackupEncrypt` receives the expected `InstanceId`, `BackupEncrypt`,
  `EncryptKeyType`, `KmsKeyId` and `KmsRegion` values and the resulting state id equals the
  instance id

#### Scenario: Read error propagation verified

- **WHEN** the mocked `DescribeBackupEncryptConfig` returns a non-not-found error
- **THEN** the `Read` callback reports a diagnostic error and leaves the state id unchanged

