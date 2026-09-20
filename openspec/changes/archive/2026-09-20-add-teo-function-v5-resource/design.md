## Context

当前 Terraform Provider 已有 `tencentcloud_teo_function` 资源（`tencentcloud/services/teo/resource_tc_teo_function.go`），管理 EdgeOne 边缘函数的增删改查，但其 Read 未同步云 API `DescribeFunctions` 返回的 `DomainComplianceRestrictions`（默认域名因合规问题产生的地区访问限制列表）。本变更新增 `tencentcloud_teo_function_v5` 资源，复用同一套 CRUD 接口并在 Read 中补充合规限制字段，作为新版资源供用户使用。

代码风格严格参考 `tencentcloud_igtm_strategy`（`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`）与现有 `tencentcloud_teo_function`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_v5` 资源，完整实现 Create / Read / Update / Delete / Import。
- Read 中暴露 `domain_compliance_restrictions`（列表，含 `reason`、`region`）。
- Create 为异步接口：拿到 `function_id` 后先 `d.SetId()`，再轮询 `DescribeFunctions` 直到 `Domain` 返回。
- 复合 ID 使用 `zone_id#function_id`（`tccommon.FILED_SP`）。
- 注册到 `provider.go` / `provider.md`，补充单元测试与文档。

**Non-Goals:**
- 不修改现有 `tencentcloud_teo_function` 资源 schema 或行为。
- 不暴露 `function_ids` / `filters` / `Offset` / `Limit` 给资源 schema（这些属于 Describe 查询入参，资源 Read 通过 ID 查询，无需暴露）。
- 不实现数据源（datasource）。

## Decisions

### 1. 资源命名与文件命名
- 资源名：`tencentcloud_teo_function_v5`，函数名：`ResourceTencentCloudTeoFunctionV5`。
- 文件：`resource_tc_teo_function_v5.go`、`resource_tc_teo_function_v5_test.go`、`resource_tc_teo_function_v5.md`。
- 理由：与现有 `tencentcloud_teo_function` 区分，避免破坏向后兼容；后缀 `_v5` 对应云 API 资源版本语义。

### 2. Schema 字段设计
入参（与 CreateFunction 入参一一对应）：
- `zone_id`：TypeString，Required，ForceNew（站点 ID）。
- `name`：TypeString，Required（函数名称）。
- `content`：TypeString，Required（函数内容 JS 代码）。
- `remark`：TypeString，Optional（函数描述）。

计算字段（与 DescribeFunctions 出参 `Function` struct 对应）：
- `function_id`：TypeString，Computed（函数 ID）。
- `domain`：TypeString，Computed（默认域名）。
- `domain_compliance_restrictions`：TypeList，Computed，Elem 为 Resource，含 `reason`（TypeString，Computed）、`region`（TypeString，Computed）。
- `create_time`：TypeString，Computed。
- `update_time`：TypeString，Computed。

### 3. 复合 ID
- 格式：`zone_id#function_id`，使用 `tccommon.FILED_SP` 连接。
- Read / Update / Delete 中通过 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析，长度必须为 2。
- Import 支持联合 id（文档说明使用 `zone_id#function_id`）。

### 4. Create 异步轮询
- `CreateFunction` 同步返回 `FunctionId`，但实际生效需等待 `Domain` 字段下发。
- 拿到 `FunctionId` 后立即 `d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))`，再执行 `resource.StateChangeConf` 轮询。
- 轮询逻辑参考现有 `resourceTeoFunctionCreateStateRefreshFunc_0_0`：调用 `DescribeFunctions`，通过 go-template 判断 `Functions[0].Domain` 是否存在；存在则 `true`，否则 `false`。
- 注意：需新建独立的 `resourceTeoFunctionV5CreateStateRefreshFunc` 函数，避免与现有函数命名冲突。

### 5. Update 不可变参数
- `immutableArgs = []string{"zone_id", "name"}`：变更时返回 error（`zone_id` 已 ForceNew，但 `name` 需在 update 中拦截）。
- 可变参数：`remark`、`content`，通过 `ModifyFunction` 更新。
- 当 `needChange` 为 true 时构建 `ModifyFunctionRequest`，retry 调用。

### 6. Read 服务方法
- 在 `service_tencentcloud_teo.go` 新增 `DescribeTeoFunctionV5ById(ctx, zoneId, functionId) (*teo.Function, error)`。
- 实现与现有 `DescribeTeoFunctionById` 基本一致：`DescribeFunctions` 入参设 `ZoneId` + `FunctionIds`，取 `Functions[0]` 返回。
- Read 中按 nil 判断逐字段 `d.Set`，`domain_compliance_restrictions` 需遍历 `DomainComplianceRestrictions` 切片构造 `[]map[string]interface{}`。

### 7. Read 空响应处理
- 当 `respData == nil` 时：先 `log.Printf("[CRUD] teo_function_v5 id=%s", d.Id())` 保留现场，再 `d.SetId("")`。

### 8. 单元测试
- 使用 gomonkey mock 云 API 方法，不使用 TF ACC 测试套件。
- mock `CreateFunctionWithContext`、`DescribeFunctions`（含轮询调用）、`ModifyFunctionWithContext`、`DeleteFunctionWithContext`。
- 覆盖 Create（含轮询成功）、Read（含 domain_compliance_restrictions）、Update、Delete 流程。

### 9. Provider 注册
- `provider.go`：`"tencentcloud_teo_function_v5": teo.ResourceTencentCloudTeoFunctionV5(),`，按字母序插入到 `tencentcloud_teo_function_component_binding` 之后。
- `provider.md`：新增 `tencentcloud_teo_function_v5` 行。

## 云 API 关键字段定义（vendor 摘录）

### CreateFunctionRequest / Response
```go
// vendor/.../teo/v20220901/models.go:3386
type CreateFunctionRequest struct {
    ZoneId  *string // 站点 ID。
    Name    *string // 函数名称。
    Content *string // 函数内容（JS，最大 5MB）。
    Remark  *string // 函数描述（最大 60 字符）。
}
// models.go:3433
type CreateFunctionResponse struct {
    Response *CreateFunctionResponseParams
}
type CreateFunctionResponseParams struct {
    FunctionId *string // 函数 ID。
    RequestId  *string
}
```

### DescribeFunctionsRequest / Response
```go
// models.go:10484
type DescribeFunctionsRequest struct {
    ZoneId      *string
    FunctionIds []*string
    Filters     []*Filter
    Offset      *int64
    Limit       *int64 // 默认 20，最大 200
}
// models.go:10529
type DescribeFunctionsResponseParams struct {
    TotalCount *int64
    Functions  []*Function
    RequestId  *string
}
```

### Function struct（Describe 出参元素）
```go
// models.go:16932
type Function struct {
    FunctionId                 *string
    ZoneId                     *string
    Name                       *string
    Remark                     *string
    Content                    *string
    Domain                     *string
    DomainComplianceRestrictions []*ComplianceRestriction
    CreateTime                 *string
    UpdateTime                 *string
}
```

### ComplianceRestriction struct
```go
// models.go:1982
type ComplianceRestriction struct {
    Reason *string // 下发访问限制的原因。枚举：ICP_RECORD_REQUIRED、GOVERNMENT_ORDER。
    Region *string // 限制访问地区的具体国家/地区码（ISO 3166）。
}
```

### Filter struct
```go
// models.go:16822
type Filter struct {
    Name   *string
    Values []*string
}
```

### ModifyFunctionRequest
```go
// models.go:20106
type ModifyFunctionRequest struct {
    ZoneId     *string
    FunctionId *string
    Remark     *string // 不填写保持原有配置。
    Content    *string // 不填写保持原有配置。
}
```

### DeleteFunctionRequest
```go
// models.go:6892
type DeleteFunctionRequest struct {
    ZoneId     *string
    FunctionId *string
}
```

## 代码风格参照说明

### tencentcloud_igtm_strategy 关键结构
- `ResourceTencentCloudIgtmStrategy()` 返回 `*schema.Resource`，含 Create/Read/Update/Delete/Importer。
- Create：`defer LogElapsed/InconsistentCheck`，`logId`/`ctx` 初始化，构建 request，`resource.Retry(WriteRetryTimeout)`，返回值 nil 检查后 `d.SetId()` 再调 Read。
- Read：拆分复合 ID，调 service Describe，nil 检查后逐字段 `d.Set`（含嵌套列表字段遍历构造 map）。
- Update：拆分 ID，`immutableArgs` 拦截 + `mutableArgs` 检测 `HasChange`，retry 调用 Modify，再调 Read。
- Delete：拆分 ID，构建 request，retry 调用 Delete。

### tencentcloud_teo_function 异步轮询参照
- `resourceTeoFunctionCreateStateRefreshFunc_0_0`：使用 `text/template` 解析 `DescribeFunctions.Response`，判断 `Functions[0].Domain` 是否存在。
- `StateChangeConf{Delay: 10s, MinTimeout: 3s, Pending: ["false"], Target: ["true"], Timeout: 600s}`。
- `d.SetId()` 在轮询完成后执行（本资源调整为轮询前 SetId 以满足规范要求）。

## Risks / Trade-offs

- **[重复资源]** `tencentcloud_teo_function_v5` 与 `tencentcloud_teo_function` 管理同一云资源 → 通过命名后缀区分，文档中说明二者关系，不自动迁移。
- **[异步轮询超时]** Create 后 `Domain` 下发可能较慢 → 复用 600s 超时与 3s 最小间隔，超时后返回 error 由用户重试。
- **[DomainComplianceRestrictions 可能为空]** 云 API 该字段可能为 nil → Read 中先判断切片非空再 set，空时不 set（terraform 自动置空列表）。
- **[Name 不变约束]** 云 API `ModifyFunction` 不支持修改 Name → update 中 `name` 加入 immutableArgs 拦截。
