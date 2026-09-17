## Context

TencentCloud EdgeOne (TEO) 提供站点 DNS 记录管理能力。provider 中已存在 `tencentcloud_teo_dns_record` 资源（`tencentcloud/services/teo/resource_tc_teo_dns_record.go`），本次按云 API 最新参数映射新增 `tencentcloud_teo_dns_record_50` 资源（RESOURCE_KIND_GENERAL），实现 DNS 记录完整 CRUD 生命周期，并暴露 `DescribeDnsRecords` 的查询增强参数（filters / sort_by / sort_order / match）。

云 API SDK 已在 vendor 中可用：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（包别名惯例：`teov20220901`）。四个接口均已在 `vendor/.../teo/v20220901/client.go` 中验证存在：
- `CreateDnsRecord` / `CreateDnsRecordWithContext`（client.go:1355/1365）
- `DescribeDnsRecords` / `DescribeDnsRecordsWithContext`（client.go:7244/7253）
- `ModifyDnsRecords` / `ModifyDnsRecordsWithContext`（client.go:12618/12629）
- `DeleteDnsRecords` / `DeleteDnsRecordsWithContext`（client.go:4695/4705）

所有接口均为同步接口（云 API 文档未标注异步，返回值无任务 ID/状态轮询字段），无需 Read 轮询等待生效。

**重要约束**：不修改现有 `tencentcloud_teo_dns_record` 资源；新增资源是完全独立的实现，遵循 `tencentcloud_igtm_strategy` 的代码组织风格（CRUD 全部在 resource 文件内直接实现，使用 `resource.Retry` + client 直调模式，而非旧 teo 资源经 service 层的模式）。

### vendor 云 API struct 关键字段定义（models.go，实施时直接参照，无需回读）

**CreateDnsRecordRequest**（models.go:3138）：
```go
type CreateDnsRecordRequest struct {
    *tchttp.BaseRequest
    ZoneId   *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`   // 站点 ID
    Name     *string `json:"Name,omitnil,omitempty" name:"Name"`       // DNS 记录名（中文等需 punycode）
    Type     *string `json:"Type,omitnil,omitempty" name:"Type"`       // A/AAAA/MX/CNAME/TXT/NS/CAA/SRV
    Content  *string `json:"Content,omitnil,omitempty" name:"Content"` // 记录内容
    Location *string `json:"Location,omitnil,omitempty" name:"Location"` // 解析线路，默认 Default
    TTL      *int64  `json:"TTL,omitnil,omitempty" name:"TTL"`         // 60~86400，默认 300
    Weight   *int64  `json:"Weight,omitnil,omitempty" name:"Weight"`   // -1~100，默认 -1
    Priority *int64  `json:"Priority,omitnil,omitempty" name:"Priority"` // MX 优先级 0~50，默认 0
}
```

**CreateDnsRecordResponse**（models.go:3205）：
```go
type CreateDnsRecordResponse struct {
    *tchttp.BaseResponse
    Response *CreateDnsRecordResponseParams `json:"Response"`
}
// CreateDnsRecordResponseParams 字段: RecordId *string, RequestId *string
```

**DescribeDnsRecordsRequest**（models.go:9895）：
```go
type DescribeDnsRecordsRequest struct {
    *tchttp.BaseRequest
    ZoneId    *string           `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`    // 站点 ID
    Offset    *int64            `json:"Offset,omitnil,omitempty" name:"Offset"`    // 默认 0
    Limit     *int64            `json:"Limit,omitnil,omitempty" name:"Limit"`      // 默认 20，上限 1000
    Filters   []*AdvancedFilter `json:"Filters,omitnil,omitempty" name:"Filters"`  // 过滤条件
    SortBy    *string           `json:"SortBy,omitnil,omitempty" name:"SortBy"`    // content/created-on/name/ttl/type
    SortOrder *string           `json:"SortOrder,omitnil,omitempty" name:"SortOrder"` // asc/desc，默认 asc
    Match     *string           `json:"Match,omitnil,omitempty" name:"Match"`       // all/any，默认 all
}
```
Filters 支持的 Name 取值：`id`（记录 ID，模糊）、`name`（记录名，模糊）、`content`（内容，模糊）、`type`（类型，不支持模糊）、`ttl`（不支持模糊）。Filters.Values 上限 20。

**AdvancedFilter**（models.go:374）：
```go
type AdvancedFilter struct {
    Name   *string  `json:"Name,omitnil,omitempty" name:"Name"`
    Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
    Fuzzy  *bool    `json:"Fuzzy,omitnil,omitempty" name:"Fuzzy"`
}
```

**DescribeDnsRecordsResponse**（models.go:9957）：
```go
type DescribeDnsRecordsResponse struct {
    *tchttp.BaseResponse
    Response *DescribeDnsRecordsResponseParams `json:"Response"`
}
// DescribeDnsRecordsResponseParams 字段:
//   TotalCount *int64
//   DnsRecords []*DnsRecord
//   RequestId  *string
```

**DnsRecord**（models.go:15770）：
```go
type DnsRecord struct {
    ZoneId     *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`     // 仅出参，ModifyDnsRecords 传入会被忽略
    RecordId   *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`
    Name       *string `json:"Name,omitnil,omitempty" name:"Name"`
    Type       *string `json:"Type,omitnil,omitempty" name:"Type"`
    Location   *string `json:"Location,omitnil,omitempty" name:"Location"`
    Content    *string `json:"Content,omitnil,omitempty" name:"Content"`
    TTL        *int64  `json:"TTL,omitnil,omitempty" name:"TTL"`
    Weight     *int64  `json:"Weight,omitnil,omitempty" name:"Weight"`
    Priority   *int64  `json:"Priority,omitnil,omitempty" name:"Priority"`
    Status     *string `json:"Status,omitnil,omitempty" name:"Status"`     // 仅出参，Modify 传入会被忽略
    CreatedOn  *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`  // 仅出参
    ModifiedOn *string `json:"ModifiedOn,omitnil,omitempty" name:"ModifiedOn"` // 仅出参
}
```

**ModifyDnsRecordsRequest**（models.go:19746）：
```go
type ModifyDnsRecordsRequest struct {
    *tchttp.BaseRequest
    ZoneId     *string       `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`     // 站点 ID
    DnsRecords []*DnsRecord  `json:"DnsRecords,omitnil,omitempty" name:"DnsRecords"` // 一次最多修改 100 条
}
```
**ModifyDnsRecordsResponse**：仅 `RequestId`（无业务数据返回）。
注意：DnsRecord 中 `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` 在 ModifyDnsRecords 中作为入参会被忽略（见 DnsRecord 注释）。

**DeleteDnsRecordsRequest**（models.go:6702）：
```go
type DeleteDnsRecordsRequest struct {
    *tchttp.BaseRequest
    ZoneId    *string   `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`      // 站点 ID
    RecordIds []*string `json:"RecordIds,omitnil,omitempty" name:"RecordIds"` // 待删除记录 ID 列表，上限 1000
}
```
**DeleteDnsRecordsResponse**：仅 `RequestId`。

### 代码风格参照（tencentcloud_igtm_strategy 模式，实施时直接照此骨架）

参考文件：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`、`tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`（后者是 teo 包下同风格的 client 直调实现，含 gomonkey 测试样例，更贴近本次实现）。

**资源定义骨架**：
```go
func ResourceTencentCloudTeoDnsRecord50() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudTeoDnsRecord50Create,
        Read:   resourceTencentCloudTeoDnsRecord50Read,
        Update: resourceTencentCloudTeoDnsRecord50Update,
        Delete: resourceTencentCloudTeoDnsRecord50Delete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{ ... },
    }
}
```

**CRUD 函数骨架**（每个函数固定模式）：
```go
func resourceTencentCloudTeoDnsRecord50Create(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.create")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId    = tccommon.GetLogId(tccommon.ContextNil)
        ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request  = teov20220901.NewCreateDnsRecordRequest()
        response = teov20220901.NewCreateDnsRecordResponse()
        zoneId   string
        recordId string
    )

    if v, ok := d.GetOk("zone_id"); ok {
        zoneId = v.(string)
        request.ZoneId = helper.String(zoneId)
    }
    // ... 其余参数

    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        } else {
            log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
        }
        if result == nil || result.Response == nil || result.Response.RecordId == nil {
            return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s create teo dns_record_50 failed, response is nil, request body [%s]", logId, request.ToJsonString()))
        }
        response = result
        return nil
    })
    if err != nil {
        log.Printf("[CRITAL]%s create teo dns_record_50 failed, reason:%+v", logId, err)
        return err
    }

    recordId = *response.Response.RecordId
    d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))
    return resourceTencentCloudTeoDnsRecord50Read(d, meta)
}
```

**import 头部**（固定）：
```go
import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

**client 获取方式**：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()`（已验证存在于 `tencentcloud/connectivity/client.go:1974`）。

**helper 函数**：`helper.String()`、`helper.Strings()`、`helper.IntInt64()`（int → *int64）、`helper.Bool()`（如需 Fuzzy）。

**联合 ID 解析模式**（Read/Update/Delete 中）：
```go
idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
if len(idSplit) != 2 {
    return fmt.Errorf("id is broken,%s", d.Id())
}
zoneId := idSplit[0]
recordId := idSplit[1]
```

**gomonkey 单测骨架**（参照 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`，package teo_test）：
```go
type mockMetaForDnsRecord50 struct {
    client *connectivity.TencentCloudClient
}
func (m *mockMetaForDnsRecord50) GetAPIV3Conn() *connectivity.TencentCloudClient { return m.client }
var _ tccommon.ProviderMeta = &mockMetaForDnsRecord50{}

patches := gomonkey.NewPatches()
defer patches.Reset()
teoClient := &teov20220901.Client{}
patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)
patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) { ... })
patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) { ... })
d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})
err := res.Create(d, meta)
```

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_50` 资源的完整 CRUD 生命周期（Create/Read/Update/Delete）
- 支持 `zone_id` + `record_id` 联合 ID（`tccommon.FILED_SP` 分隔），支持 import
- 暴露 DescribeDnsRecords 的查询参数：`filters`（AdvancedFilter 列表：name/values/fuzzy）、`sort_by`、`sort_order`、`match`
- Read 用 `filters` 中 `id` 过滤器精确查询单条记录（非模糊）
- Update 用 ModifyDnsRecords（单条记录，DnsRecords 列表长度为 1，仅设置 RecordId + 可变字段）
- Delete 用 DeleteDnsRecords（单条记录，RecordIds 列表长度为 1）
- 遵循现有代码模式：retry（ReadRetryTimeout/WriteRetryTimeout）、tccommon.RetryError、LogElapsed、InconsistentCheck
- 提供单元测试（gomonkey mock 云 API，不使用 TF_ACC 验收套件）

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_dns_record` 资源（保持向后兼容）
- 不实现数据源 `tencentcloud_teo_dns_records`（本次仅实现 RESOURCE_KIND_GENERAL 资源）
- 不实现 `ModifyDnsRecordsStatus`（启停 DNS 记录状态不在本次接口映射范围内）
- 不暴露 Offset/Limit 分页参数给用户（Read 按记录 ID 精确过滤，单条结果；filters.Values 上限 20、Describe Limit 内部固定取值足够）
- 不实现批量管理（ModifyDnsRecords/DeleteDnsRecords 虽为批量接口，本资源仅管理单条记录）

## Decisions

### 1. 资源命名与文件命名
**决策**：资源名 `tencentcloud_teo_dns_record_50`，文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_50.go`，入口函数 `ResourceTencentCloudTeoDnsRecord50`，CRUD 函数 `resourceTencentCloudTeoDnsRecord50{Create,Read,Update,Delete}`。

**理由**：需求明确指定资源名为 `tencentcloud_teo_dns_record_50`（数字后缀与现有 `tencentcloud_teo_dns_record` 并存，如同参数映射工具有序号区分的场景）。Go 函数名按 CamelCase 惯例将 `50` 写为 `50`（如 `DnsRecord50`）。

### 2. 资源 ID 设计
**决策**：`d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))`，即 `zoneId#recordId`。

**理由**：Read/Update/Delete 均需 `zone_id` + 记录标识。`record_id` 是 `DeleteDnsRecords.RecordIds` 与 `DnsRecords.RecordId` 的必要入参；`zone_id` 是所有接口的必填入参。import 示例需说明联合 ID 格式 `terraform import tencentcloud_teo_dns_record_50.xxx zone-xxx#record-xxx`。

### 3. Schema 字段设计
**决策**：

| 字段 | 类型 | 必填性 | ForceNew | 说明 |
|------|------|--------|----------|------|
| `zone_id` | String | Required | 是 | 站点 ID，创建后不可改（Modify 虽有 ZoneId 入参但它是请求顶层定位参数，改 zone 等价于换资源） |
| `name` | String | Required | 否 | DNS 记录名，可修改 |
| `type` | String | Required | 否 | 记录类型，可修改 |
| `content` | String | Required | 否 | 记录内容，可修改 |
| `location` | String | Optional+Computed | 否 | 解析线路，默认 Default |
| `ttl` | Int | Optional+Computed | 否 | 60~86400，默认 300 |
| `weight` | Int | Optional+Computed | 否 | -1~100，默认 -1 |
| `priority` | Int | Optional+Computed | 否 | MX 优先级 0~50，默认 0 |
| `filters` | TypeList | Optional | 否 | 查询过滤条件，elem 为 Resource{name(Required String), values(Required List[String]), fuzzy(Optional Bool)} |
| `sort_by` | String | Optional | 否 | content/created-on/name/ttl/type |
| `sort_order` | String | Optional | 否 | asc/desc |
| `match` | String | Optional | 否 | all/any |
| `record_id` | String | Computed | - | 记录 ID |
| `status` | String | Computed | - | enable/disable（ModifyDnsRecords 入参中 Status 被忽略，只读） |
| `created_on` | String | Computed | - | 创建时间（ModifyDnsRecords 入参中被忽略，只读） |
| `modified_on` | String | Computed | - | 修改时间（ModifyDnsRecords 入参中被忽略，只读） |

**理由**：参数映射指定 zone_id/name/type/content 必填；location/ttl/weight/priority 可选且云 API 有默认值 → Optional+Computed（与现有 teo_dns_record 一致，避免未传时 perpetual diff）。`filters` 的 name/values 按映射为必填、fuzzy 可选。`sort_by`/`sort_order`/`match` 按映射为可选（对应 DescribeDnsRecords 查询入参）。

**关键取舍——查询参数在通用资源中的定位**：`filters`/`sort_by`/`sort_order`/`match` 是 Describe 接口的查询入参。它们属于"用户指定的查询行为参数"，不影响资源实体属性；Read 的核心逻辑始终是按 `filters: [{name: "id", values: [recordId], fuzzy: false}]` 查询该资源本身。为满足参数映射要求，将这些字段纳入 schema（Optional），用户设置的值会在每次 Read 调用时随请求发送（作为附加过滤/排序条件），但在 Read 回填时不从响应中 set（响应无对应字段，避免覆盖用户配置）。资源实体字段（name/type/content/...）仍从响应的 DnsRecords[0] 回填。

### 4. Read 实现方式
**决策**：直接在 resource 文件内调用 `DescribeDnsRecords`（igtm_strategy 风格的 client 直调，不经 service 层新增方法），用 retry 包裹：
- 必带过滤：`Filters = [{Name: "id", Values: [recordId]}]`（非模糊，精确匹配记录 ID）
- 若用户在 schema 中配置了额外的 `filters`/`sort_by`/`sort_order`/`match`，则合并进请求（用户 filters 追加在 id 过滤器之后）
- `Limit` 不暴露给用户，请求中固定设为云 API 注释标注的最大值 1000（单条记录查询足够）
- 响应处理：`len(response.Response.DnsRecords) == 0` → 先 `log.Printf("[CRUD] read teo dns_record_50 id=%s", d.Id())` 保留现场，再 `d.SetId("")` 返回 nil（资源已不存在）
- 命中记录后逐字段 set（先判 nil 再 set）

**理由**：DescribeDnsRecords 无单条 Get 接口，只能列表过滤查询。固定 `id` 过滤器保证定位唯一性。规则要求"查询接口中有分页字段，则给定值应该是云API接口注释中标注的最大值"（Limit 上限 1000）。

**替代方案**：在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoDnsRecord50ById` 服务层方法（旧 teo 资源风格）——被否决，因为任务要求严格参照 igtm_strategy 风格，且同包 teo 下新资源（multi_path_gateway 等）也已采用 client 直调模式。

### 5. Update 实现方式
**决策**：`mutableArgs := []string{"name", "type", "content", "location", "ttl", "weight", "priority", "filters", "sort_by", "sort_order", "match"}` 检查 HasChange；仅当可变字段变化时调用 `ModifyDnsRecords`：
```go
request := teov20220901.NewModifyDnsRecordsRequest()
request.ZoneId = helper.String(zoneId)
dnsRecord := &teov20220901.DnsRecord{RecordId: helper.String(recordId)}
// 逐字段从 d.GetOk / d.GetOkExists 填充 name/type/content/location/ttl/weight/priority
request.DnsRecords = []*teov20220901.DnsRecord{dnsRecord}
```
不设置 DnsRecord.ZoneId/Status/CreatedOn/ModifiedOn（vendor 注释明确这四个字段在 ModifyDnsRecords 入参中被忽略）。更新成功后调用 Read 刷新状态。

**理由**：ModifyDnsRecords 是批量接口，单资源场景封装为长度 1 的列表。查询参数（filters/sort_by/...）不发给 Modify（它不接受这些参数），但纳入 HasChange 判断以触发 Read 刷新无意义、纯粹避免冗余 Modify 调用——实际实现中查询参数不参与 needChange 判断也可（它们不影响服务端资源状态）；**最终决策：mutableArgs 仅包含实体可变字段 `name/type/content/location/ttl/weight/priority`**，查询参数变化不触发 Modify 调用（Modify 接口不接受查询参数，无服务端效果，避免无意义写操作）。

### 6. Delete 实现方式
**决策**：
```go
request := teov20220901.NewDeleteDnsRecordsRequest()
request.ZoneId = helper.String(zoneId)
request.RecordIds = helper.Strings([]string{recordId})
```
retry 包裹 `DeleteDnsRecordsWithContext`（WriteRetryTimeout）。删除成功即返回，无需 Read 轮询（同步接口）。

**理由**：DeleteDnsRecords 是批量接口（上限 1000），单资源场景传长度 1 的列表。

### 7. Create 返回值校验
**决策**：retry 块内校验 `result == nil || result.Response == nil || result.Response.RecordId == nil || *result.Response.RecordId == ""` → 返回 `resource.NonRetryableError`（带 logId 与请求体，便于排障）。

**理由**：规则要求 Create 后必须检查返回值为空的全部形式，避免写入空 ID 造成状态混乱。

### 8. 测试方式
**决策**：使用 gomonkey mock 云 API 进行单元测试（`resource_tc_teo_dns_record_50_test.go`，package `teo_test`），覆盖：
- Create 成功（校验请求参数与联合 ID 写入）
- Create 返回空 RecordId → 报错
- Read 成功（校验 Describe 请求带 id 过滤器、字段回填）
- Read 记录不存在 → SetId("")
- Update 可变字段变化 → ModifyDnsRecords 调用（校验 DnsRecords 长度 1 且含 RecordId）
- Delete 成功（校验 RecordIds）
- 各 API 报错路径

**理由**：项目规则要求新增 terraform 资源使用 gomonkey mock 单测，不用 TF 验收套件。测试中需注意唯一 mock 辅助类型命名（避免与同包其他 *_test.go 的 mock 类型/辅助函数重名，如 `mockMetaForDnsRecord50`、`ptrStrDR50`/`ptrInt64DR50`），且运行需 `-gcflags=all=-l` 禁用内联。

### 9. 文档与注册
**决策**：
- `tencentcloud/provider.go`：在 ResourcesMap 中添加 `"tencentcloud_teo_dns_record_50": teo.ResourceTencentCloudTeoDnsRecord50(),`（放在 `tencentcloud_teo_dns_record` 相邻位置，保持 gofmt 对齐）
- `tencentcloud/provider.md`：在 TEO Resource 列表 `tencentcloud_teo_dns_record` 之后添加 `tencentcloud_teo_dns_record_50`
- `resource_tc_teo_dns_record_50.md`：一句话描述（"Provides a resource to create a teo dns_record_50"）+ Example Usage + Import 部分（说明联合 ID `zoneId#recordId`），不手写 Argument/Attribute Reference（工具自动生成）

**理由**：遵循项目注册与文档规范；RESOURCE_KIND_GENERAL 资源含 Import 部分。

## Risks / Trade-offs

- [Risk] 用户误将 `filters` 配置为会过滤掉自身记录的条件（如 name 不匹配）→ Read 查不到记录会 `d.SetId("")`，Terraform 会尝试重建。**缓解**：Read 的 id 精确过滤器始终置于 filters 首位，且 id 与用户附加 filters 在默认 `match=all` 下取交集，文档示例中说明 filters 用于附加过滤。
- [Risk] DescribeDnsRecords 返回空可能因云 API 短暂波动 → 资源被误判删除。**缓解**：Read 的 retry 块内对 API 错误重试；空结果按通用资源规则打印 `[CRUD]` 日志后 SetId("")（本资源非 DATASOURCE，不适用 NonRetryableError 强失败策略）。
- [Risk] `weight` 取值 -1 与 `d.GetOkExists` 语义：GetOkExists 对 0 值视为"未设置"，但 weight 的合法值含 0（表示不解析）。**缓解**：与现有 teo_dns_record 资源一致使用 GetOkExists（0 权重场景罕见，且 Computed 回填会覆盖），在文档中说明 weight=0 的含义。
- [Risk] ModifyDnsRecords 单条封装与批量接口的语义差异 → 若未来需要批量管理需要新资源。**缓解**：本资源定位为单条记录管理，与现有 teo_dns_record 一致。
- [Risk] 与现有 `tencentcloud_teo_dns_record` 功能重叠造成用户困惑。**缓解**：两者独立共存，向后兼容性不受影响；新资源额外支持查询参数。

## Migration Plan

纯新增资源，无迁移。回滚方式：移除 provider.go/provider.md 中的注册项与新增文件即可。
