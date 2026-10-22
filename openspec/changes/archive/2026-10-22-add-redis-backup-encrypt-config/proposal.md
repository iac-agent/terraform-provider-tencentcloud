## Why

Redis（云数据库 Redis）实例的备份加密能力（开启/关闭备份加密、切换加密密钥类型、绑定自定义 KMS 密钥）目前无法通过 Terraform 管理。用户需要在基础设施即代码的场景下查询与修改实例备份加密配置，以保证备份数据的安全合规性。

## What Changes

- 新增 Terraform 资源 `tencentcloud_redis_backup_encrypt`（RESOURCE_KIND_CONFIG 类型，单实例单配置项）：
  - Read：调用 `DescribeBackupEncryptConfig` 读取实例备份加密配置。
  - Update：调用 `ModifyBackupEncrypt` 更新备份加密配置。
  - Create：`d.SetId(instance_id)` 后直接复用 Update 逻辑（配置随实例存在，无独立创建接口）。
  - Delete：空操作（no-op），配置随实例存在，删除仅从 state 中移除。
  - 支持 Import（ID 为 `instance_id`）。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册该资源。
- 新增资源文档 `resource_tc_redis_backup_encrypt_config.md` 与单元测试 `resource_tc_redis_backup_encrypt_config_test.go`。

## Capabilities

### New Capabilities

- `redis-backup-encrypt-config`: 管理腾讯云 Redis 实例备份加密配置（读取与更新备份加密开关、加密密钥类型、KMS 密钥 ID 与 KMS 地域）。

### Modified Capabilities

<!-- 无既有 capability 的需求变更 -->

## Impact

- 新增代码文件：`tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.go`、`resource_tc_redis_backup_encrypt_config_test.go`、`resource_tc_redis_backup_encrypt_config.md`。
- 修改文件：`tencentcloud/provider.go`（注册资源）、`tencentcloud/provider.md`（资源清单）。
- 依赖云 API（SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`，已在 vendor 中）：
  - `DescribeBackupEncryptConfig`（查询实例备份加密配置）。
  - `ModifyBackupEncrypt`（更改备份加密，开启/关闭加密、修改密钥）。
- 向后兼容：纯新增资源，不影响任何现有资源的 schema 与 state。
