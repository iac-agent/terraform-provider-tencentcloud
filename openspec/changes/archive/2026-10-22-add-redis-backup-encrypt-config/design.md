# Design: tencentcloud_redis_backup_encrypt（Redis 备份加密配置资源）

## Context

腾讯云 Redis 实例提供备份加密能力（开启/关闭备份加密、选择加密密钥类型、绑定自定义 KMS 密钥），用于保障备份数据的安全合规。当前 Terraform Provider（`tencentcloud/services/crs/` 为 redis 产品的 service 目录）尚未提供管理该配置的资源。

本资源类型为 **RESOURCE_KIND_CONFIG**：单实例单配置项，配置随实例存在而存在，无独立的创建/删除接口。CRUD 退化为：
- Create = `d.SetId(instance_id)` 后复用 Update 逻辑；
- Read = 调用 `DescribeBackupEncryptConfig`；
- Update = 调用 `ModifyBackupEncrypt`；
- Delete = no-op（直接返回 nil，仅从 state 中移除）。

资源代码文件名按 CONFIG 类型命名规范为 `resource_tc_redis_backup_encrypt_config.go`，注册的资源名称为 `tencentcloud_redis_backup_encrypt`。

## Goals / Non-Goals

**Goals:**

- 新增 `tencentcloud_redis_backup_encrypt` 资源，支持读取与更新 Redis 实例备份加密配置。
- 支持 `terraform import`（ID 即 `instance_id`）。
- 代码风格与既有 crs 资源（`resource_tc_redis_backup_config.go`、`resource_tc_redis_connection_config.go`）保持一致。
- 单元测试使用 gomonkey mock 云 API（不使用 terraform 验收测试套件）。

**Non-Goals:**

- 不管理 KMS 密钥本身的创建/轮换（密钥由 `tencentcloud_kms_key` 等资源管理）。
- 不修改既有资源（如 `tencentcloud_redis_instance`、`tencentcloud_redis_backup_config`）的 schema。
- 不生成 `_extension.go` 文件。

## 云 API 定义（vendor 摘录）

SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`
（vendor 路径：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412/`）

### DescribeBackupEncryptConfig（查询实例备份加密配置）

`models.go` 行 2341-2404：

```go
type DescribeBackupEncryptConfigRequest struct {
	*tchttp.BaseRequest
	// 实例 ID，请登录Redis控制台在实例列表复制实例 ID。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

type DescribeBackupEncryptConfigResponseParams struct {
	// 备份是否加密。枚举值：true：加密；false：不加密
	BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
	// 加密密钥类型。枚举值：0：腾讯云托管（默认）；1：用户自定义 KMS 密钥
	EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
	// KMS的密钥ID。参数格式：xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	// KMS地域信息。参数格式：ap-xxxx
	KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
	// 唯一请求 ID
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeBackupEncryptConfigResponse struct {
	*tchttp.BaseResponse
	Response *DescribeBackupEncryptConfigResponseParams `json:"Response"`
}
```

client 方法（client.go 行 1797 起）：
```go
func (c *Client) DescribeBackupEncryptConfigWithContext(ctx context.Context, request *DescribeBackupEncryptConfigRequest) (response *DescribeBackupEncryptConfigResponse, err error)
func (c *Client) DescribeBackupEncryptConfig(request *DescribeBackupEncryptConfigRequest) (response *DescribeBackupEncryptConfigResponse, err error)
```

### ModifyBackupEncrypt（更改备份加密，开启/关闭加密、修改密钥）

`models.go` 行 8698-8777：

```go
type ModifyBackupEncryptRequest struct {
	*tchttp.BaseRequest
	// 实例 ID，请登录Redis控制台在实例列表复制实例 ID。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
	// 是否开启备份加密。默认值：false
	BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
	// 加密密钥类型。枚举值：0：腾讯云托管（默认）；1：用户自定义 KMS 密钥。默认值：0
	EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
	// KMS的密钥ID。参数格式：xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	// 地域信息。参数格式：ap-xxxx。取值参考：KMS支持的地域列表
	KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
}

type ModifyBackupEncryptResponse struct {
	*tchttp.BaseResponse
	Response *ModifyBackupEncryptResponseParams `json:"Response"` // 仅含 RequestId
}
```

client 方法（client.go 行 5609 起）：
```go
func (c *Client) ModifyBackupEncryptWithContext(ctx context.Context, request *ModifyBackupEncryptRequest) (response *ModifyBackupEncryptResponse, err error)
func (c *Client) ModifyBackupEncrypt(request *ModifyBackupEncryptRequest) (response *ModifyBackupEncryptResponse, err error)
```

两个接口均为同步接口（ModifyBackupEncrypt 响应不含 FlowId/TaskId/JobId，无需异步轮询）。

## 参数映射（Schema 设计）

| SchemaName | Type | Required/Optional | 云API JsonPath（Read 出参 / Modify 入参） |
|---|---|---|---|
| `instance_id` | String | Required, ForceNew | `request.InstanceId` / `response.Response.InstanceId`（Read 中用 d.Id()） |
| `backup_encrypt` | Bool | Required | `response.Response.BackupEncrypt` / `request.BackupEncrypt` |
| `encrypt_key_type` | Int | Optional, Computed | `response.Response.EncryptKeyType` / `request.EncryptKeyType` |
| `kms_key_id` | String | Optional, Computed | `response.Response.KmsKeyId` / `request.KmsKeyId` |
| `kms_region` | String | Optional, Computed | `response.Response.KmsRegion` / `request.KmsRegion` |

说明：
- `instance_id` 为 ForceNew：更换实例视同新建配置。
- `backup_encrypt` 必填（对应 ModifyBackupEncrypt 的必填入参）。
- `encrypt_key_type`/`kms_key_id`/`kms_region` 为 Optional+Computed：云API返回值有值时回填 state，用户未配置时以云端现状为准（Computed 行为）。
- 资源 ID：`instance_id`（单实例单配置，无需联合 ID、无需 FILED_SP 分隔符）。

## Decisions

### D1. Create 复用 Update（CONFIG 资源标准骨架）

参照 `resource_tc_redis_backup_config.go`（行 59-65）与 `resource_tc_redis_connection_config.go`（行 74-88）：

```go
func resourceTencentCloudRedisBackupEncryptConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	d.SetId(instanceId)
	return resourceTencentCloudRedisBackupEncryptConfigUpdate(d, meta)
}
```

### D2. Read：retry 内检查空响应，NotFound 处理

- 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `DescribeBackupEncryptConfig`。
- `response == nil || response.Response == nil` 时返回 `tccommon.RetryError(e)` 包装的错误 / `NonRetryableError`（按 Read 空响应规范：CONFIG 资源 Read 中，若云API错误码为 `ResourceNotFound.InstanceNotExists`（`crs` 包已有常量 `RedisInstanceNotFound = "ResourceNotFound.InstanceNotExists"`，定义于 `extension_redis.go` 行 74），先 `log.Printf("[CRUD] redis backup_encrypt id=%s", d.Id())` 再 `d.SetId("")` 并返回 nil。
- 设置字段前逐一判空：`BackupEncrypt`/`EncryptKeyType`/`KmsKeyId`/`KmsRegion` 为 nil 时跳过 `d.Set(...)`。
- retry 块内仅调用云API；`d.Set` 等成功操作放在 retry 外/错误处理后。由于 Read 需要用 response 内容 set 字段，参照 `resource_tc_redis_backup_config.go` 的写法：在 retry 内拿到数据后，在 retry 成功路径上完成 set（或先取出数据，retry 结束后 set——两种均可，保持与 crs 现有代码风格一致即可；注意 set 前判空）。

### D3. Update：直接全量提交，成功后回读

- 用 `d.HasChange(...)` 判断任一字段变化后调用 `ModifyBackupEncrypt`（参照 igtm_strategy 的 `needChange` + `mutableArgs` 模式，字段列表：`backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region`）。
- `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 `ModifyBackupEncrypt` 调用，错误用 `tccommon.RetryError(e)` 包装。
- `response == nil || response.Response == nil` 返回 `resource.NonRetryableError`。
- Update 成功后调用 Read 回读最新配置（保证 Computed 字段与云端一致）。
- ModifyBackupEncrypt 为同步接口，无需异步轮询、无需在 update 中额外等待任务。

### D4. Delete：no-op

配置随实例存在，删除仅移除 state：

```go
func resourceTencentCloudRedisBackupEncryptConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	return nil
}
```

> 不选择"Delete 时调用 ModifyBackupEncrypt 关闭加密"：关闭加密属于破坏性操作，用户 destroy 时意外关闭备份加密风险大于收益；与 `resource_tc_redis_connection_config.go` 的 no-op Delete 保持一致。

### D5. 服务层实现位置

不新建 service 文件，直接在资源文件中通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient()` 调用（参照 `resource_tc_redis_connection_config.go` 行 174-175 的直调风格，以及 igtm_strategy 在 CRUD 内直调 client 的风格）。如需封装，可复用 `RedisService`（`service_tencentcloud_redis.go` 行 26 的 struct），但本资源字段简单，直调即可，避免冗余抽象。

### D6. Import 支持

```go
Importer: &schema.ResourceImporter{
	State: schema.ImportStatePassthrough,
},
```

ID 即 `instance_id`。.md 文档 Import 示例说明使用实例 ID。

### D7. 日志与命名规范

- 资源名（小写蛇形）统一使用 `redis backup_encrypt`（如 `log.Printf("[CRUD] redis backup_encrypt id=%s", d.Id())`），日志检索友好。
- 不在资源 go 文件开头添加注释。
- `defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.<op>")()` + `defer tccommon.InconsistentCheck(d, meta)()`。

## 代码风格参照（关键结构说明）

### 参照 1：tencentcloud_igtm_strategy（用户指定主参照）

文件：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`

- Schema 组织：顶层 map[string]*schema.Schema，Required/Optional/Computed/ForceNew 明确，Description 齐全。
- CRUD 骨架：Create（构造 request → retry 调用 → 校验响应非空 → SetId → Read 回读）；Read（拆 ID → service 查询 → 判空处理 → 逐字段判空 set）；Update（`needChange` + `mutableArgs` + `d.HasChange` → 构造 request → retry 调用 → 响应校验 → Read 回读）；Delete（retry 调用删除接口）。
- retry 模式：`resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { result, e := client.Xxx(ctx, request); if e != nil { return tccommon.RetryError(e) } else { log.Printf("[DEBUG]...") }; if result == nil || result.Response == nil { return resource.NonRetryableError(...) }; response = result; return nil })`。
- 错误日志：`log.Printf("[CRITAL]%s xxx failed, reason:%+v", logId, reqErr)`。
- 注意：本资源 Delete 为 no-op、Create 为 SetId+Update，与 igtm_strategy 的 Create/Delete（真实创建/删除）不同，这一点参照 crs 的 CONFIG 资源。

### 参照 2：crs 现有 CONFIG 资源（同产品目录风格）

文件：`tencentcloud/services/crs/resource_tc_redis_backup_config.go`、`tencentcloud/services/crs/resource_tc_redis_connection_config.go`

- Create：`d.SetId(<instance_id>)` 后直接 `return resourceXxxUpdate(d, meta)`。
- Delete：no-op 直接 `return nil`（connection_config）。
- `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。
- Read retry 内处理 `RedisInstanceNotFound` 错误码（`sdkErr.Code == RedisInstanceNotFound` → `d.SetId("")`）。
- 直调 client：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient()`。

### 参照 3：单测风格（gomonkey mock，不跑验收测试）

文件：`tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config_test.go`（gomonkey 风格范本）

```go
type mockMetaXxx struct { client *connectivity.TencentCloudClient }
func (m *mockMetaXxx) GetAPIV3Conn() *connectivity.TencentCloudClient { return m.client }
var _ tccommon.ProviderMeta = &mockMetaXxx{}

patches := gomonkey.NewPatches()
defer patches.Reset()
redisClient := &redis.Client{}
patches.ApplyMethodReturn(newMockMeta().client, "UseRedisClient", redisClient)
patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfig", func(request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) { ... })
d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})
d.SetId("crs-xxxx")
err := res.Read(d, meta)
```

## Risks / Trade-offs

- [关闭加密为破坏性操作] → Delete 采用 no-op，不在 destroy 时自动关闭备份加密（D4），文档中明确说明。
- [加密开启后历史备份不受影响] → 文档说明该配置只对新备份生效，属云API行为，provider 不做额外处理。
- [`kms_key_id`/`kms_region` 为 Optional+Computed，用户关闭加密后云端可能返回空] → Read set 前判空跳过，避免将空值写入 state。
- [`backup_encrypt=false` 时传 `encrypt_key_type=1` 等参数可能报错] → Update 全量透传用户配置，云API校验错误直接透出给用户，不做本地拦截。
- [Read 中空响应与 NotFound 判定] → `ResourceNotFound.InstanceNotExists` 视为资源不存在（清空 ID）；其余错误进入 retry 重试。

## Migration Plan

纯新增资源，无迁移。部署顺序：合入代码 → `make doc` 生成文档 → 发布新版本。回滚：移除 provider 注册即可，无 state 迁移成本。

## Open Questions

无。云API字段与需求映射已在 vendor 中核实（两个接口的入参/出参与 JsonPath 完全对应）。
