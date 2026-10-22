# redis-backup-encrypt-config Specification

## Purpose
定义 `tencentcloud_redis_backup_encrypt` 资源的能力规格：通过云API `DescribeBackupEncryptConfig`/`ModifyBackupEncrypt` 管理腾讯云 Redis 实例的备份加密配置（开关、加密密钥类型、KMS 密钥与地域），涵盖 schema 定义、读取回填、更新提交、创建（配置初始化）、删除（no-op）、资源导入、Provider 注册及文档与单元测试要求。
## Requirements
### Requirement: 资源 tencentcloud_redis_backup_encrypt 的 Schema 定义

系统 SHALL 提供 `tencentcloud_redis_backup_encrypt` 资源（RESOURCE_KIND_CONFIG 类型），包含以下参数：

- `instance_id`（String，必填，ForceNew）：Redis 实例 ID。
- `backup_encrypt`（Bool，必填）：是否开启备份加密，`true` 开启、`false` 关闭。
- `encrypt_key_type`（Int，可选，Computed）：加密密钥类型，`0` 腾讯云托管（默认）、`1` 用户自定义 KMS 密钥。
- `kms_key_id`（String，可选，Computed）：KMS 密钥 ID，格式 `xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx`。
- `kms_region`（String，可选，Computed）：KMS 地域，格式 `ap-xxxx`。

资源 ID SHALL 为 `instance_id`（单实例单配置，无联合 ID）。

#### Scenario: 声明资源并指定全部参数
- **WHEN** 用户在 terraform 配置中声明 `tencentcloud_redis_backup_encrypt` 并指定 `instance_id`、`backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region`
- **THEN** schema 校验通过，资源被接受并进入 plan/apply 流程

#### Scenario: instance_id 变更触发重建
- **WHEN** 用户修改已有资源配置中的 `instance_id`
- **THEN** terraform 计划销毁旧资源并创建新资源（ForceNew 语义）

### Requirement: 读取备份加密配置

系统 SHALL 通过调用云API `DescribeBackupEncryptConfig`（request.InstanceId = 资源 ID）读取实例备份加密配置，并将 `response.Response.BackupEncrypt`、`response.Response.EncryptKeyType`、`response.Response.KmsKeyId`、`response.Response.KmsRegion` 分别回填到 `backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region` 字段。

回填字段前系统 SHALL 判断对应响应字段是否为 nil，为 nil 时跳过 set，避免空值写入 state。

当云API返回 `ResourceNotFound.InstanceNotExists` 错误时，系统 SHALL 先打印含资源 ID 的日志（`log.Printf("[CRUD] redis backup_encrypt id=%s", d.Id())`）再执行 `d.SetId("")` 并正常返回。

#### Scenario: 读取开启加密的配置
- **WHEN** 资源执行 Read，云API返回 `BackupEncrypt=true`、`EncryptKeyType=1`、`KmsKeyId=xxx`、`KmsRegion=ap-guangzhou`
- **THEN** state 中 `backup_encrypt=true`、`encrypt_key_type=1`、`kms_key_id=xxx`、`kms_region=ap-guangzhou`

#### Scenario: 实例不存在
- **WHEN** 资源执行 Read，云API返回错误码 `ResourceNotFound.InstanceNotExists`
- **THEN** 日志中打印资源 ID，state 中的资源被移除（ID 置空），Read 正常返回无错误

#### Scenario: 云API返回部分字段为空
- **WHEN** 资源执行 Read，云API返回 `KmsKeyId` 为 nil（如使用腾讯云托管密钥）
- **THEN** state 中 `kms_key_id` 不被覆盖为空值，其余非空字段正常回填

### Requirement: 更新备份加密配置

系统 SHALL 在任一可变字段（`backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region`）发生变化时，通过调用云API `ModifyBackupEncrypt` 提交更新，request 字段与 schema 字段一一对应（`InstanceId`、`BackupEncrypt`、`EncryptKeyType`、`KmsKeyId`、`KmsRegion`）。

系统 SHALL 使用 `tccommon.WriteRetryTimeout` 作为超时时间对 `ModifyBackupEncrypt` 调用添加 retry 处理，接口失败时使用 `tccommon.RetryError()` 包装错误；retry 块内仅调用云API接口。

系统 SHALL 在调用返回后检查 `response`/`response.Response` 是否为空，为空时返回 NonRetryableError。

更新提交成功后系统 SHALL 调用 Read 回读最新配置。

#### Scenario: 开启备份加密并使用腾讯云托管密钥
- **WHEN** 用户将 `backup_encrypt` 从 `false` 改为 `true`（`encrypt_key_type` 未指定或为 `0`）
- **THEN** 系统调用 `ModifyBackupEncrypt`（`BackupEncrypt=true`），成功后回读配置，state 中 `backup_encrypt=true`

#### Scenario: 切换为自定义 KMS 密钥
- **WHEN** 用户设置 `encrypt_key_type=1` 并指定 `kms_key_id` 与 `kms_region`
- **THEN** 系统调用 `ModifyBackupEncrypt` 时在 request 中携带 `EncryptKeyType=1`、`KmsKeyId`、`KmsRegion`，成功后 state 与云端配置一致

#### Scenario: 无字段变化时跳过更新
- **WHEN** 用户执行 plan/apply 且 `backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region` 均无变化
- **THEN** 系统不调用 `ModifyBackupEncrypt`，直接回读配置

#### Scenario: 云API调用失败
- **WHEN** `ModifyBackupEncrypt` 调用返回错误
- **THEN** 错误经 `tccommon.RetryError()` 包装进入重试，重试耗尽后错误返回给用户，state 不被破坏

### Requirement: 创建资源（配置初始化）

系统 SHALL 在 Create 时将资源 ID 设置为 `instance_id`（`d.SetId(instance_id)`），随后复用 Update 逻辑完成首次配置提交。

#### Scenario: 首次创建
- **WHEN** 用户 apply 一个新的 `tencentcloud_redis_backup_encrypt` 资源
- **THEN** 资源 ID 被设置为 `instance_id`，随后按更新逻辑调用 `ModifyBackupEncrypt` 完成配置提交，并回读配置

### Requirement: 删除资源为空操作

系统 SHALL 将 Delete 实现为空操作（no-op）：不调用任何云API，直接返回 nil。配置随实例存在，删除仅从 state 中移除。

#### Scenario: 销毁资源
- **WHEN** 用户执行 `terraform destroy` 销毁该资源
- **THEN** 系统不调用云API、不关闭备份加密，资源仅从 state 中移除

### Requirement: 资源导入

系统 SHALL 支持 `terraform import`（`schema.ImportStatePassthrough`），导入 ID 为 Redis 实例 ID。

#### Scenario: 导入已有配置
- **WHEN** 用户执行 `terraform import tencentcloud_redis_backup_encrypt.example crs-xxxx`
- **THEN** 资源被导入 state，随后 Read 回填云端当前配置

### Requirement: Provider 注册

系统 SHALL 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `tencentcloud_redis_backup_encrypt`，并在 `tencentcloud/provider.md` 资源清单中添加该资源名称。

#### Scenario: 资源可用
- **WHEN** 用户使用新版本 provider
- **THEN** `tencentcloud_redis_backup_encrypt` 资源类型可被 terraform 识别并使用

### Requirement: 文档与单元测试

系统 SHALL 提供 `resource_tc_redis_backup_encrypt_config.md` 资源文档（含 Example Usage 与 Import 示例，不手写 Argument Reference / Attribute Reference），并在 `resource_tc_redis_backup_encrypt_config_test.go` 中使用 gomonkey mock 云API 提供单元测试（覆盖 Read 正常/实例不存在、Update 变更提交等场景）。

#### Scenario: 文档完整
- **WHEN** 开发者查阅 `resource_tc_redis_backup_encrypt_config.md`
- **THEN** 文档包含一句话描述（带 Redis 云产品名）、Example Usage 与 Import 示例

#### Scenario: 单元测试通过
- **WHEN** 运行 `resource_tc_redis_backup_encrypt_config_test.go` 中的单元测试
- **THEN** mock 云API返回后资源 Read/Update 逻辑行为符合预期，测试全部通过

