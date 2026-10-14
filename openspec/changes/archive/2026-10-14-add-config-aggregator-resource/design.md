# Design: add-config-aggregator-resource

## 1. 云 API 关键结构定义（摘自 vendor config/v20220802/models.go）

文件: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802/models.go`

### AggregatorAccount (L684-691)

```go
type AggregatorAccount struct {
	MemberUin *uint64 `json:"MemberUin,omitnil,omitempty" name:"MemberUin"`
	MemberName *string `json:"MemberName,omitnil,omitempty" name:"MemberName"` // 可能为 null
}
```

### CreateAggregator (L1227-1291)

```go
type CreateAggregatorRequestParams struct {
	Name *string                 // 账号组名称
	Description *string         // 账号组描述
	Type *string                // 枚举: RD 全局账号组 / CUSTOM 自定义账号组
	AggregatorAccounts []*AggregatorAccount // 账号组成员信息列表，最多100个
}
type CreateAggregatorResponseParams struct {
	AccountGroupId *string // 账号组Id
	RequestId *string
}
```
注意：Create 请求**无 OwnerUin**。

### DescribeAggregator (L2130-2193)

```go
type DescribeAggregatorRequestParams struct {
	AccountGroupId *string // 账号组ID
	OwnerUin *uint64        // 账号组创建者ID
}
type DescribeAggregatorResponseParams struct {
	Name *string
	Description *string
	Type *string
	AggregatorAccounts []*AggregatorAccount // 可能返回 null
	AggregatorStatus *uint64               // 创建状态
	RequestId *string
}
```

### UpdateAggregator (L5432-5500)

```go
type UpdateAggregatorRequestParams struct {
	Name *string
	Description *string
	AccountGroupId *string
	OwnerUin *uint64
	AggregatorAccounts []*AggregatorAccount
}
type UpdateAggregatorResponseParams struct { RequestId *string } // 无业务出参
```

### DeleteAggregators (L1512-1559)

```go
type DeleteAggregatorsRequestParams struct {
	AccountGroupId *string
	OwnerUin *uint64
}
type DeleteAggregatorsResponseParams struct { RequestId *string }
```

客户端调用方式（参照 resource_tc_config_rule.go L251）：
```go
meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CreateAggregatorWithContext(ctx, request)
```

## 2. 代码风格参照：tencentcloud/services/igtm/resource_tc_igtm_strategy.go

结构要点：
- 函数命名: `ResourceTencentCloudConfigAggregator()` 返回 `*schema.Resource`，Create/Read/Update/Delete/Importer(`schema.ImportStatePassthrough`)。
- 每个 CRUD 函数开头:
  ```go
  defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.create")()
  defer tccommon.InconsistentCheck(d, meta)()
  var (
      logId = tccommon.GetLogId(tccommon.ContextNil)
      ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
      ...
  )
  ```
- Create: `resource.Retry(tccommon.WriteRetryTimeout, ...)`，retry 内只调 API；`RetryError(e)` 包装错误；成功后 `log.Printf("[DEBUG]%s api[%s] success...")`；retry 外校验 response/id 非空否则报错，然后 `d.SetId()`，最后 `return resourceTencentCloudConfigAggregatorRead(d, meta)`。
- Read: 从 `d.Id()` 用 `strings.Split(d.Id(), tccommon.FILED_SP)` 拆出 owner_uin 与 account_group_id；`d.Set()` 前判断字段非 nil；嵌套 TypeList 转换为 `[]map[string]interface{}`；资源不存在时先 `log.Printf("[CRUD] config_aggregator id=%s", d.Id())` 再 `d.SetId("")`。
- Update: `mutableArgs` 检查 `d.HasChange`（本资源除 computed 外均可变：name/description/aggregator_accounts/owner_uin），构造 `UpdateAggregatorRequest`，`request.AccountGroupId` 与 `request.OwnerUin` 从 id 拆分获取，retry 调 `UpdateAggregatorWithContext`。
- Delete: `DeleteAggregatorsRequest`，同上取 id 字段，retry 调用。
- helper 包: `helper.String` / `helper.IntUint64` / `helper.Uint64ToStr` / `helper.StrToUint64Point`。

## 3. 资源设计

### Schema

```go
Schema: map[string]*schema.Schema{
    "name":            {Type: TypeString, Required: true, Description: "账号组名称"},
    "description":     {Type: TypeString, Required: true, Description: "账号组描述"},
    "type":            {Type: TypeString, Required: true, Description: "账号组类型: RD(全局)/CUSTOM(自定义)"},
    "owner_uin":       {Type: TypeInt, Required: true, Description: "账号组创建者ID"},
    "aggregator_accounts": {Type: TypeList, Optional: true, Description: "账号组成员信息列表，最多100个",
        Elem: &schema.Resource{Schema: map[string]*schema.Schema{
            "member_uin":  {Type: TypeInt, Required: true, Description: "成员ID"},
            "member_name": {Type: TypeString, Required: true, Description: "成员名称"},
        }},
    },
    // computed
    "account_group_id": {Type: TypeString, Computed: true, Description: "账号组Id"},
    "aggregator_status": {Type: TypeInt, Computed: true, Description: "创建状态"},
}
```

### ID 格式

`<owner_uin><tccommon.FILED_SP><account_group_id>`，Create 成功后设置；Read/Update/Delete 从 id 拆分；支持 `terraform import`（联合 id）。

### 注意点

- `OwnerUin`/`MemberUin` 云 API 为 `*uint64`，schema 用 `TypeInt`，写入用 `helper.IntUint64(v)`；从 id 拆出时用 `helper.StrToUint64Point`。
- Describe 响应的 `AggregatorAccounts` 可能为 null，`d.Set` 前判空。
- Create 校验：`response == nil || response.Response == nil || response.Response.AccountGroupId == nil || *response.Response.AccountGroupId == ""` → `resource.NonRetryableError`；打印 logId 便于排障。
- 所有 4 个接口均为同步接口，无需轮询，Create 无需在 SetId 前做异步等待。
- 不生成 `_extension.go`；go 文件开头不加注释；日志中统一使用 `config_aggregator` 指代资源。

## 4. 文件清单

- `tencentcloud/services/config/resource_tc_config_aggregator.go`（新资源）
- `tencentcloud/services/config/resource_tc_config_aggregator_test.go`（gomonkey mock 单测，不用 tf 测试套件）
- `tencentcloud/services/config/resource_tc_config_aggregator.md`（一句话描述 + Example Usage + Import 联合 id 示例）
- `tencentcloud/provider.go` / `tencentcloud/provider.md` 注册 `"tencentcloud_config_aggregator": config.ResourceTencentCloudConfigAggregator()`（参照 `tencentcloud_igtm_strategy` 注册方式）
- 收尾阶段由 tfpacer-finalize 执行 gofmt / make doc / changelog。
