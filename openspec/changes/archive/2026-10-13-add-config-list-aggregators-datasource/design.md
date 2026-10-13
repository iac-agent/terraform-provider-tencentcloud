# Design: add-config-list-aggregators-datasource

## 云 API 摘录（vendor 按行区间）

包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`（vendor 路径 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802/`）

### ListAggregatorsRequest（models.go L3476-L3484）

```go
type ListAggregatorsRequest struct {
    *tchttp.BaseRequest
    // 每页显示数量
    Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
    // 起始
    Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`
}
```

- 请求仅有分页参数 `Limit`/`Offset`，无过滤参数。分页 Limit 取值参考同产品 `ListCompliancePacks` 使用 100。

### ListAggregatorsResponse / ResponseParams（models.go L3507-L3521）

```go
type ListAggregatorsResponseParams struct {
    // 总数
    Total *uint64 `json:"Total,omitnil,omitempty" name:"Total"`
    // 账号组列表
    Items []*Aggregator `json:"Items,omitnil,omitempty" name:"Items"`
    // 唯一请求 ID
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ListAggregatorsResponse struct {
    *tchttp.BaseResponse
    Response *ListAggregatorsResponseParams `json:"Response"`
}
```

### Aggregator（models.go L654-L682）

```go
type Aggregator struct {
    Name           *string `json:"Name,omitnil,omitempty"`           // 账号组名称
    Description    *string `json:"Description,omitnil,omitempty"`    // 账号组描述
    OwnerUin       *uint64 `json:"OwnerUin,omitnil,omitempty"`       // 创建者用户ID
    CreateTime     *string `json:"CreateTime,omitnil,omitempty"`     // 创建时间
    AccountCount   *uint64 `json:"AccountCount,omitnil,omitempty"`   // 账号组成员数量
    Type           *string `json:"Type,omitnil,omitempty"`           // RD:全局账号组
    AccountGroupId *string `json:"AccountGroupId,omitnil,omitempty"` // 账号组ID
    AggregatorStatus *uint64 `json:"AggregatorStatus,omitnil,omitempty"` // 账号组状态
    MemberName     *string `json:"MemberName,omitnil,omitempty"`     // 账号组名称（可能返回 null）
}
```

### Client 方法（client.go L2020 起）

- 构造: `configv20220802.NewListAggregatorsRequest()`
- 调用: `me.client.UseConfigV20220802Client().ListAggregatorsWithContext(ctx, request)`（`WithApiInfo("config", APIVersion, "ListAggregators")`）

## 代码风格参照

### 数据源骨架参照: `tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`

- `DataSourceTencentCloudIgtmInstanceList() *schema.Resource`：Schema 顶层含可选输入 + computed 列表（`Elem: &schema.Resource{Schema: ...}`）+ `result_output_file`（`Optional`, "Used to save results."）。
- Read 函数骨架：
  1. `defer tccommon.LogElapsed("data_source.tencentcloud_igtm_instance_list.read")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
  2. `logId = tccommon.GetLogId(nil)`; `ctx = tccommon.NewResourceLifeCycleHandleFuncContext(...)`
  3. `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中调用 service 层查询，失败 `tccommon.RetryError(e)`；retry 块外做 `d.Set(...)`
  4. `d.SetId(helper.BuildToken())`
  5. `result_output_file` 处理: `tccommon.WriteToFile(output.(string), d)`

### 服务层参照: `tencentcloud/services/config/service_tencentcloud_config.go` `DescribeConfigCompliancePacksByFilter`（L25-L96）

- `type ConfigService struct { client *connectivity.TencentCloudClient }`，构造函数 `NewConfigService`。
- 分页循环：初始 `Limit=100, Offset=0`，每次请求在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 API；retry 块内仅做 API 调用与校验（`result == nil || result.Response == nil` → `resource.NonRetryableError`；`ratelimit.Check(request.GetAction())`；`log.Printf("[DEBUG]%s api[%s] success, ...")`），retry 外累加 `*request.Offset += *request.Limit`，`len(items) >= total` 时退出循环。
- 新增方法签名建议：`func (me *ConfigService) DescribeConfigListAggregators(ctx context.Context) (total uint64, items []*configv20220802.Aggregator, errRet error)`

### 数据源为空返回处理（规则要求）

- 服务层或 Read retry 块内若 `result == nil || result.Response == nil` → `resource.NonRetryableError`，并 `log.Printf("[DATASOURCE] read empty, skip SetId")`。
- Read 中 retry 块外先 `d.Set("total", ...) / d.Set("items", ...)`，最后 `d.SetId(helper.BuildToken())`；不调用 `d.SetId("")`。

## 目标实现设计

新文件 `tencentcloud/services/config/data_source_tc_config_list_aggregators.go`，函数 `DataSourceTencentCloudConfigListAggregators()`。

Schema（顶层平铺，无额外 wrapper）：

| SchemaName | Type | 说明 |
|---|---|---|
| total | TypeInt, Computed | 账号组总数（Response.Total） |
| items | TypeList, Computed, Elem:&schema.Resource | 账号组列表（Response.Items），元素字段如下 |
| items.name | TypeString, Computed | 账号组名称 |
| items.description | TypeString, Computed | 账号组描述 |
| items.owner_uin | TypeInt, Computed | 创建者用户ID（uint64） |
| items.create_time | TypeString, Computed | 创建时间 |
| items.account_count | TypeInt, Computed | 账号组成员数量 |
| items.type | TypeString, Computed | RD:全局账号组 |
| items.account_group_id | TypeString, Computed | 账号组ID |
| items.aggregator_status | TypeInt, Computed | 账号组状态 |
| items.member_name | TypeString, Computed | 账号组名称（可能为 null） |
| result_output_file | TypeString, Optional | Used to save results. |

Read 逻辑：无输入参数 → 直接调用 service 层分页查询（service 内已含 retry，Read 外层无需再 retry）→ 遍历 `Items` 逐字段判 nil 后 set → `d.SetId(helper.BuildToken())` → `result_output_file` 输出。

## 注册与文档

- `tencentcloud/provider.go`：数据源 map 中按字母序位置添加 `"tencentcloud_config_list_aggregators": config.DataSourceTencentCloudConfigListAggregators(),`（参照 `"tencentcloud_igtm_instance_list"` 的注册格式）。
- `tencentcloud/provider.md`：同步添加一行 `- tencentcloud_config_list_aggregators: Use this data source to query config account group list.`（参照现有 config 数据源条目）。
- `data_source_tc_config_list_aggregators.md`：
  - 一句话描述带产品名（Config）："Use this data source to query TencentCloud Config (配置审计) account group (aggregator) list."
  - Example Usage（无输入参数，仅引用 total/items/result_output_file）
  - 无 Import 部分（DATASOURCE 类型不添加）
  - 不添加 Argument/Attribute Reference 两节

## 单元测试

- 文件 `data_source_tc_config_list_aggregators_test.go`，使用 gomonkey mock `ListAggregatorsWithContext`，验证 Read 后 `total`/`items` 字段与 `result_output_file` 输出；不使用 terraform 测试套件，不实际执行。

## 注意事项

- 该接口为纯查询接口（非异步），无需 `d.SetId()` 异步轮询逻辑。
- `uint64` 出参使用 `helper.Int/`schema.TypeInt 落地（参照现有 config 数据源对 `*uint64` 的处理）。
- 生成代码后禁止在本阶段执行 `go build`/`go vet`/`go test`；`gofmt`、`make doc`、changelog 由收尾阶段统一执行。
