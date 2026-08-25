## 1. 资源代码实现

- [x] 1.1 创建资源文件 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.go`，声明 `package crs` 与导入（`context`、`log`、`resource.Retry`、`schema`、redis sdk `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`、`tccommon`、`helper`），代码风格参考 `resource_tc_redis_maintenance_window.go` 与 `resource_tc_redis_connection_config.go`，文件开头不加注释。
- [x] 1.2 实现 `ResourceTencentCloudRedisInstancePasswordPolicy()`，定义扁平化 schema：`instance_id`(Required, TypeString)、`enabled`(Required, TypeBool)、`min_letter_count`(Optional, TypeInt)、`min_digit_count`(Optional, TypeInt)、`min_special_count`(Optional, TypeInt)、`min_length`(Optional, TypeInt)；声明 Create/Read/Update/Delete 四个回调与 `Importer`（`schema.ImportStatePassthrough`）。
- [x] 1.3 实现 `resourceTencentCloudRedisInstancePasswordPolicyCreate`：`defer tccommon.LogElapsed` / `defer tccommon.InconsistentCheck`；从 schema 取 `instance_id` 赋给局部变量；`d.SetId(instanceId)`；直接 `return resourceTencentCloudRedisInstancePasswordPolicyUpdate(d, meta)`。
- [x] 1.4 实现 `resourceTencentCloudRedisInstancePasswordPolicyRead`：以 `d.Id()` 作为 `InstanceId`，构建 `redis.NewDescribeInstancePasswordPolicyRequest`，使用 `tccommon.ReadRetryTimeout` + `tccommon.RetryError` 调用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeInstancePasswordPolicy(request)`；对 `response==nil`/`response.Response==nil`/`PasswordPolicy==nil` 判空（先 `log.Printf("[CRUD] redis_instance_password_policy id=%s", d.Id())` 再 `d.SetId("")`）；仅当各字段非 nil 时 `_ = d.Set(...)` 回填 `instance_id`、`enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length`。
- [x] 1.5 实现 `resourceTencentCloudRedisInstancePasswordPolicyUpdate`：构建 `redis.NewModifyInstancePasswordPolicyRequest`，`InstanceId` 取自 `d.Id()`；构建 `redis.PasswordPolicy{}` 并填充 `Enabled`（必填，始终发送）、`MinLetterCount`/`MinDigitCount`/`MinSpecialCount`/`MinLength`（仅当对应字段 `d.GetOk` 为真时设置）；使用 `tccommon.WriteRetryTimeout` + `tccommon.RetryError` 调用 `ModifyInstancePasswordPolicy`；成功后 `return resourceTencentCloudRedisInstancePasswordPolicyRead(d, meta)`；retry 块内只调用接口，不执行 set 操作。
- [x] 1.6 实现 `resourceTencentCloudRedisInstancePasswordPolicyDelete`：仅 `defer tccommon.LogElapsed` / `defer tccommon.InconsistentCheck`，直接 `return nil`，不调用任何云 API（配置依附于实例，删除仅移除 state）。

## 2. Provider 注册与文档

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中，于 redis 资源组（约 1684–1712 行）按字母顺序插入 `"tencentcloud_redis_instance_password_policy": crs.ResourceTencentCloudRedisInstancePasswordPolicy(),`。
- [x] 2.2 创建文档文件 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.md`：一句话描述带上云产品名 Redis；Example Usage（含依赖实例的示例）；Import 段说明使用 `instance_id` 导入；不添加 Argument/Attribute Reference（由工具生成）。

## 3. 单元测试

- [x] 3.1 创建测试文件 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件），覆盖 Create/Read/Update/Delete 业务逻辑（含字段回填、空响应处理、enabled 必填、可选字段按需发送）；检查所有 error 返回，无错路径用 `_ =` 处理。

## 4. 验证（收尾阶段执行）

- [ ] 4.1 执行 `gofmt` 格式化新增/修改的 Go 文件（由 tfpacer-finalize 统一执行）。
- [ ] 4.2 通过 `make doc` 生成 website/docs 文档（由 tfpacer-finalize 统一执行，禁止手动编写 website/ 文件）。
