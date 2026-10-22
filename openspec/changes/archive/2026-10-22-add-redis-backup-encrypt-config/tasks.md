# 实施任务清单：tencentcloud_redis_backup_encrypt

## 1. 资源实现（tencentcloud/services/crs/）

- [x] 1.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.go`：定义 `ResourceTencentCloudRedisBackupEncryptConfig()` 及 schema（`instance_id` String Required ForceNew；`backup_encrypt` Bool Required；`encrypt_key_type` Int Optional Computed；`kms_key_id` String Optional Computed；`kms_region` String Optional Computed），带 `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`
- [x] 1.2 实现 `resourceTencentCloudRedisBackupEncryptConfigCreate`：`d.SetId(instance_id)` 后复用 Update 逻辑
- [x] 1.3 实现 `resourceTencentCloudRedisBackupEncryptConfigRead`：在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeBackupEncryptConfig(request)`（request.InstanceId = d.Id()），错误用 `tccommon.RetryError(e)` 包装；`RedisInstanceNotFound` 错误码时先 `log.Printf("[CRUD] redis backup_encrypt id=%s", d.Id())` 再 `d.SetId("")` 返回 nil；回填 `backup_encrypt`/`encrypt_key_type`/`kms_key_id`/`kms_region` 前逐一判空
- [x] 1.4 实现 `resourceTencentCloudRedisBackupEncryptConfigUpdate`：用 `d.HasChange` 检查 `backup_encrypt`/`encrypt_key_type`/`kms_key_id`/`kms_region` 是否有变化，有变化时在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `ModifyBackupEncrypt`（request 字段与 schema 一一对应，Optional 字段用 `d.GetOkExists`/`d.GetOk` 取值），`response`/`response.Response` 为空时返回 `resource.NonRetryableError`；成功后调用 Read 回读
- [x] 1.5 实现 `resourceTencentCloudRedisBackupEncryptConfigDelete`：no-op，直接返回 nil
- [x] 1.6 检查代码规范：文件开头无注释；所有 err 均被处理（必定不出错的用 `_ =`）；日志统一使用资源名 `redis backup_encrypt`；函数返回值 error 均检查

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` ResourcesMap 中注册 `"tencentcloud_redis_backup_encrypt": crs.ResourceTencentCloudRedisBackupEncryptConfig()`（放在 redis 资源区块，参考 `tencentcloud_redis_backup_config` 的位置）
- [x] 2.2 在 `tencentcloud/provider.md` 的 Resource 清单中（redis 区块）添加 `tencentcloud_redis_backup_encrypt`

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md`：一句话描述（含 Redis 云产品名，"Provides a resource to ..." 格式）、Example Usage、Import 示例（说明导入 ID 为 Redis 实例 ID）；不添加 Argument Reference / Attribute Reference 章节

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config_test.go`（package crs_test，使用 gomonkey mock 云API，参照 `resource_tc_cfw_cluster_fw_bypass_config_test.go` 风格：mockMeta + `patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)` + `patches.ApplyMethodFunc`）
- [x] 4.2 编写 Read 正常用例：mock `DescribeBackupEncryptConfig` 返回 `BackupEncrypt=true`、`EncryptKeyType=1`、`KmsKeyId`、`KmsRegion`，断言字段回填正确且 ID 不变
- [x] 4.3 编写 Read 实例不存在用例：mock 返回 `ResourceNotFound.InstanceNotExists` 错误，断言 ID 被清空且无错误
- [x] 4.4 编写 Update 用例：mock `ModifyBackupEncrypt` 成功与 `DescribeBackupEncryptConfig` 回读，断言调用无错误、字段符合预期
- [x] 4.5 确认测试代码可正确编译（不执行 `go test`）

## 5. 收尾（由 tfpacer-finalize 统一执行，本阶段禁止执行）

- [ ] 5.1 `gofmt` 格式化变更的 Go 文件
- [ ] 5.2 `make doc` 生成 website/ 文档与 changelog（禁止手动修改 website/ 与 .changelog/）
