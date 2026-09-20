## Context

腾讯云 EdgeOne (TEO) 边缘函数（Function）允许用户在边缘节点运行 JavaScript 代码。云 API 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 提供了完整的 CRUD 接口。当前 provider 已存在由 iacg 代码生成器生成的 `tencentcloud_teo_function` 资源，但本次需求要求新增一个手写实现的 `tencentcloud_teo_function_v6` 资源，严格遵循项目代码规范，并完整覆盖云 API 的字段。

本设计文档沉淀后续实施所需的全部参考信息（云 API struct 关键字段、代码风格参照说明），实施阶段只需读取本文档即可获得所需信息，无需再回读 vendor 大文件。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_v6` 资源，完整覆盖 `CreateFunction`、`DescribeFunctions`、`ModifyFunction`、`DeleteFunction` 四个云 API 接口。
- 严格遵循项目代码规范：retry 处理、错误检查、复合 ID、异步创建轮询、日志打印等。
- 提供 gomonkey mock 单元测试（新增资源不使用 terraform 测试套件）。
- 提供 `.md` 文档与 provider 注册。

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_function` 资源。
- 不实现 `DescribeFunctions` 中 `Filters` 作为用户可配置的 schema 入参（Function 资源为通用资源，通过 `function_id` 查询，Filters 为只读查询字段，不在 schema 中暴露）。
- 不暴露分页参数 `offset`/`limit` 给用户（资源 Read 通过 FunctionIds 精确查询单个函数，无需分页）。

## Decisions

### 1. 资源命名与文件组织
- 资源名：`tencentcloud_teo_function_v6`
- 文件：`tencentcloud/services/teo/resource_tc_teo_function_v6.go`
- 服务层方法：`DescribeTeoFunctionV6ById`（在 `service_tencentcloud_teo.go` 中新增，避免与已有 `DescribeTeoFunctionById` 命名冲突）
- 函数命名前缀：`resourceTencentCloudTeoFunctionV6`（避免与已有 `resourceTencentCloudTeoFunction` 冲突）

### 2. Schema 字段定义

根据云 API 入参与出参映射，schema 字段如下：

**入参字段（用户可配置）：**
| SchemaName | 类型 | 必填/可选 | ForceNew | 说明 |
|---|---|---|---|---|
| `zone_id` | TypeString | Required | 是 | 站点 ID（CreateFunction 必填，ModifyFunction/DeleteFunction 也需要，且不可变更） |
| `name` | TypeString | Required | 是 | 函数名称（CreateFunction 必填；ModifyFunction 不支持修改 name，设为 ForceNew） |
| `content` | TypeString | Required | 否 | 函数内容（JS 代码；CreateFunction 必填，ModifyFunction 可选） |
| `remark` | TypeString | Optional | 否 | 函数描述（CreateFunction 可选，ModifyFunction 可选） |

**出参字段（Computed，只读）：**
| SchemaName | 类型 | 说明 |
|---|---|---|
| `function_id` | TypeString | 函数 ID（CreateFunction 返回） |
| `domain` | TypeString | 函数默认域名 |
| `domain_compliance_restrictions` | TypeList | 域名合规限制列表 |
| `domain_compliance_restrictions.reason` | TypeString | 限制原因 |
| `domain_compliance_restrictions.region` | TypeString | 限制地区 |
| `create_time` | TypeString | 创建时间 |
| `update_time` | TypeString | 更新时间 |

> 说明：`name` 字段设为 ForceNew，因为 `ModifyFunction` 接口不支持修改函数名称（接口入参无 Name 字段），修改 name 需要重建资源。

### 3. 复合 ID
- 使用 `zoneId#functionId` 作为复合 ID（分隔符 `tccommon.FILED_SP` = `#`）。
- 在 Read/Update/Delete 中通过 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析出 `zoneId` 与 `functionId`。
- 支持 import（RESOURCE_KIND_GENERAL 类型资源），导入时使用复合 ID。

### 4. 异步创建轮询
`CreateFunction` 是异步接口：创建后函数需部署到边缘节点，部署完成后 `DescribeFunctions` 返回的 `Function` 中 `Domain` 字段才会有值。参照现有 `resource_tc_teo_function.go` 的 `resourceTeoFunctionCreateStateRefreshFunc_0_0` 实现：
- 在拿到 `functionId` 后、轮询之前，**先 `d.SetId()`**（遵循规范：避免轮询失败时 tfstate 中无 id 导致无法 destroy）。
- 使用 `resource.StateChangeConf` 轮询 `DescribeFunctions`，检查 `Functions[0].Domain` 是否非空，判断部署是否完成。
- 轮询参数：Delay 10s，MinTimeout 3s，Timeout 600s。

### 5. CRUD 逻辑要点
- **Create**: 构建 `CreateFunctionRequest`（ZoneId、Name、Content、Remark）→ retry 调用 → 检查返回 FunctionId 是否为空（空则返回 NonRetryableError）→ `d.SetId()` → 异步轮询 → 调用 Read。
- **Read**: 解析复合 ID → 调用 `DescribeTeoFunctionV6ById` → 若返回为空先打印 `log.Printf("[CRUD] ...")` 再 `d.SetId("")` → 非空则逐字段 set（set 前判断 nil）。
- **Update**: 解析复合 ID → 检查 immutableArgs（`name`、`zone_id` 已为 ForceNew，但仍需检查 `name`）→ 若 remark/content 有变更则构建 `ModifyFunctionRequest`（ZoneId、FunctionId、Remark、Content）→ retry 调用 → 调用 Read。
- **Delete**: 解析复合 ID → 构建 `DeleteFunctionRequest`（ZoneId、FunctionId）→ retry 调用。

### 6. 服务层方法
新增 `DescribeTeoFunctionV6ById(ctx, zoneId, functionId) (*teo.Function, error)`，逻辑与已有 `DescribeTeoFunctionById` 一致：
- 构建 `DescribeFunctionsRequest`，设置 `ZoneId` 与 `FunctionIds`（单元素列表）。
- 使用 `tccommon.ReadRetryTimeout` retry 调用 `DescribeFunctions`。
- 若 `response.Response == nil || len(response.Response.Functions) < 1` 返回 nil。

### 7. 文档与注册
- `.md` 文件：一句话描述带上云产品名称（EdgeOne/TEO），包含 Example Usage 与 Import 部分（说明使用复合 id）。
- `provider.go`：在 resources map 中注册 `"tencentcloud_teo_function_v6": teo.ResourceTencentCloudTeoFunctionV6()`。
- `provider.md`：在资源列表中新增 `tencentcloud_teo_function_v6`。

## Risks / Trade-offs

- **[与现有资源并存]** → 新增 `tencentcloud_teo_function_v6` 与已有 `tencentcloud_teo_function` 功能重叠，但二者独立存在、互不影响，用户可自行选择。这是需求明确要求的，不属于重复资源冲突。
- **[异步轮询超时]** → CreateFunction 部署可能耗时，轮询超时设为 600s，超时后 terraform 报错但 id 已写入 state，用户可 destroy 重试。
- **[name 字段 ForceNew]** → ModifyFunction 不支持改 name，必须 ForceNew，符合云 API 能力。

## 云 API Struct 关键字段定义（vendor 摘录）

以下为后续实施所需的云 API request/response struct 关键字段，摘自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`。

### CreateFunctionRequest (models.go:3386)
```go
type CreateFunctionRequest struct {
    *tchttp.BaseRequest
    ZoneId  *string `json:"ZoneId,omitnil,omitempty"`   // 站点 ID。必填
    Name    *string `json:"Name,omitnil,omitempty"`     // 函数名称。必填
    Content *string `json:"Content,omitnil,omitempty"`  // 函数内容(JS)。必填
    Remark  *string `json:"Remark,omitnil,omitempty"`   // 函数描述。可选
}
```

### CreateFunctionResponse (models.go:3433)
```go
type CreateFunctionResponse struct {
    *tchttp.BaseResponse
    Response *CreateFunctionResponseParams `json:"Response"`
}
type CreateFunctionResponseParams struct {
    FunctionId *string `json:"FunctionId,omitnil,omitempty"`  // 函数 ID
    RequestId  *string `json:"RequestId,omitnil,omitempty"`
}
```

### DescribeFunctionsRequest (models.go:10484)
```go
type DescribeFunctionsRequest struct {
    *tchttp.BaseRequest
    ZoneId      *string   `json:"ZoneId,omitnil,omitempty"`      // 站点 ID
    FunctionIds []*string `json:"FunctionIds,omitnil,omitempty"` // 按函数 ID 过滤
    Filters     []*Filter `json:"Filters,omitnil,omitempty"`     // 过滤条件
    Offset      *int64    `json:"Offset,omitnil,omitempty"`      // 分页偏移
    Limit       *int64    `json:"Limit,omitnil,omitempty"`       // 分页限制，最大 200
}
```

### DescribeFunctionsResponse (models.go:10540)
```go
type DescribeFunctionsResponse struct {
    *tchttp.BaseResponse
    Response *DescribeFunctionsResponseParams `json:"Response"`
}
type DescribeFunctionsResponseParams struct {
    TotalCount *int64       `json:"TotalCount,omitnil,omitempty"`
    Functions  []*Function  `json:"Functions,omitnil,omitempty"`  // 函数列表
    RequestId  *string      `json:"RequestId,omitnil,omitempty"`
}
```

### Function (models.go:16932)
```go
type Function struct {
    FunctionId                 *string                 `json:"FunctionId,omitnil,omitempty"`                 // 函数 ID
    ZoneId                     *string                 `json:"ZoneId,omitnil,omitempty"`                     // 站点 ID
    Name                       *string                 `json:"Name,omitnil,omitempty"`                       // 函数名字
    Remark                     *string                 `json:"Remark,omitnil,omitempty"`                     // 函数描述
    Content                    *string                 `json:"Content,omitnil,omitempty"`                    // 函数内容
    Domain                     *string                 `json:"Domain,omitnil,omitempty"`                     // 函数默认域名
    DomainComplianceRestrictions []*ComplianceRestriction `json:"DomainComplianceRestrictions,omitnil,omitempty"` // 域名合规限制列表
    CreateTime                 *string                 `json:"CreateTime,omitnil,omitempty"`                 // 创建时间
    UpdateTime                 *string                 `json:"UpdateTime,omitnil,omitempty"`                 // 修改时间
}
```

### ComplianceRestriction (models.go:1982)
```go
type ComplianceRestriction struct {
    Reason *string `json:"Reason,omitnil,omitempty"`  // 限制原因
    Region *string `json:"Region,omitnil,omitempty"`  // 限制地区
}
```

### Filter (models.go:16822)
```go
type Filter struct {
    Name   *string   `json:"Name,omitnil,omitempty"`   // 过滤字段
    Values []*string `json:"Values,omitnil,omitempty"` // 过滤值
}
```

### ModifyFunctionRequest (models.go:20106)
```go
type ModifyFunctionRequest struct {
    *tchttp.BaseRequest
    ZoneId     *string `json:"ZoneId,omitnil,omitempty"`     // 站点 ID
    FunctionId *string `json:"FunctionId,omitnil,omitempty"` // 函数 ID
    Remark     *string `json:"Remark,omitnil,omitempty"`     // 函数描述（不填保持原配置）
    Content    *string `json:"Content,omitnil,omitempty"`    // 函数内容（不填保持原配置）
}
```
> 注意：ModifyFunctionRequest 无 Name 字段，故 `name` 必须 ForceNew。

### ModifyFunctionResponse (models.go:20150)
```go
type ModifyFunctionResponse struct {
    *tchttp.BaseResponse
    Response *ModifyFunctionResponseParams `json:"Response"`
}
type ModifyFunctionResponseParams struct {
    RequestId *string `json:"RequestId,omitnil,omitempty"`
}
```

### DeleteFunctionRequest (models.go:6892)
```go
type DeleteFunctionRequest struct {
    *tchttp.BaseRequest
    ZoneId     *string `json:"ZoneId,omitnil,omitempty"`     // 站点 ID
    FunctionId *string `json:"FunctionId,omitnil,omitempty"` // 函数 ID
}
```

### DeleteFunctionResponse (models.go:6928)
```go
type DeleteFunctionResponse struct {
    *tchttp.BaseResponse
    Response *DeleteFunctionResponseParams `json:"Response"`
}
type DeleteFunctionResponseParams struct {
    RequestId *string `json:"RequestId,omitnil,omitempty"`
}
```

## 代码风格参照文件说明

### 参照资源：tencentcloud_igtm_strategy (`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`)
该资源是项目代码规范的标准参照，关键结构说明：
- **Schema 组织**：`ResourceTencentCloudIgtmStrategy()` 返回 `*schema.Resource`，包含 Create/Read/Update/Delete 与 Importer（`schema.ImportStatePassthrough`）。Schema 为 `map[string]*schema.Schema`，Required 字段含 ForceNew，Computed 字段标注 Computed。
- **Create 骨架**：`defer LogElapsed/InconsistentCheck` → 获取 logId/ctx → 构建 request → 逐字段从 `d.GetOk` 填充 → `resource.Retry(WriteRetryTimeout, ...)` 调用，retry 内检查 `result == nil || result.Response == nil` 返回 NonRetryableError → 检查返回 id 字段非空 → `d.SetId(strings.Join(..., FILED_SP))` → 调用 Read。
- **Read 骨架**：解析复合 ID → 调用 service 层 Describe → 若空则打印 WARN 日志并 `d.SetId("")` → 逐字段判断 nil 后 `d.Set`。
- **Update 骨架**：解析复合 ID → 检查 immutableArgs → 检查 needChange → 构建 ModifyRequest → retry 调用 → 调用 Read。
- **Delete 骨架**：解析复合 ID → 构建 DeleteRequest → retry 调用。

### 参照现有资源：resource_tc_teo_function.go (`tencentcloud/services/teo/resource_tc_teo_function.go`)
该文件展示了同云 API 的异步创建轮询实现：
- **异步轮询**：`resourceTeoFunctionCreateStateRefreshFunc_0_0` 使用 `resource.StateChangeConf`，通过 go-template 判断 `Functions[0].Domain` 是否非空。
- **轮询时机**：在 `d.SetId()` 之前轮询（现有实现）；本资源按规范调整为先 `d.SetId()` 再轮询。
- **复合 ID**：`strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP)`。

### 服务层参照：DescribeTeoFunctionById (`tencentcloud/services/teo/service_tencentcloud_teo.go:1367`)
- 构建 `DescribeFunctionsRequest`，设置 `ZoneId` 与 `FunctionIds`。
- `ratelimit.Check(request.GetAction())`。
- `resource.Retry(ReadRetryTimeout, ...)` 调用 `DescribeFunctions`。
- 若 `response.Response == nil || len(response.Response.Functions) < 1` 返回 nil。
- 返回 `response.Response.Functions[0]`。

### Mock 测试参照：resource_tc_teo_function_rule_test.go
- 包名 `teo_test`。
- 使用 `github.com/agiledragon/gomonkey/v2` mock 云 API 方法。
- 使用 `github.com/stretchr/testify/assert` 断言。