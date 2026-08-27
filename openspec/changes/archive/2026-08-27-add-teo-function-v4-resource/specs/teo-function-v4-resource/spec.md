## ADDED Requirements

### Requirement: 资源注册
provider SHALL 暴露资源类型 `tencentcloud_teo_function_v4`，用于管理单个腾讯云 TEO 边缘函数（Edge Function）的完整生命周期。该资源 MUST 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册，并在 `tencentcloud/provider.md` 中追加对应注册说明。

#### Scenario: 资源类型可被发现
- **WHEN** 用户在配置中声明 `resource "tencentcloud_teo_function_v4" "<name>"` 并执行 `terraform plan`
- **THEN** Terraform 能正确解析该资源类型，不报 "unknown resource" 错误，并展示计划创建。

#### Scenario: provider 编译通过
- **WHEN** 执行 `go build ./tencentcloud/...`
- **THEN** 构建成功，无与该新资源相关的编译错误。

### Requirement: schema 与 CreateFunction 入参对齐
资源 schema SHALL 暴露 `CreateFunction` 接受的所有入参作为顶层属性，且不重命名、不合并：
- `zone_id`（string，Required，**ForceNew**）：站点 ID。
- `name`（string，Required，**ForceNew**）：函数名称，只能包含小写字母、数字、连字符，以数字或字母开头和结尾，最大 30 个字符。
- `content`（string，Required）：函数内容，当前仅支持 JavaScript 代码，最大 5MB。
- `remark`（string，Optional）：函数描述，最大 60 个字符。

资源 SHALL 额外暴露来自 `DescribeFunctions` 响应的只读 computed 字段：
- `function_id`（string，Computed）：函数 ID。
- `domain`（string，Computed）：函数默认域名。
- `create_time`（string，Computed）：创建时间（UTC，ISO 8601）。
- `update_time`（string，Computed）：修改时间（UTC，ISO 8601）。

#### Scenario: 必填 SDK 字段均存在
- **WHEN** 开发者检视资源 schema
- **THEN** `CreateFunctionRequest` 的所有入参（`ZoneId`、`Name`、`Content`、`Remark`）都以语义等价的类型出现在 schema 中。

#### Scenario: 不存在多余的 schema 字段
- **WHEN** 开发者检视资源 schema
- **THEN** 除上述列出的输入与 computed 字段外，不存在任何派生标志或合成字段。

### Requirement: 复合资源 ID
资源 ID SHALL 为 `zone_id` 与 `function_id` 以 `tccommon.FILED_SP` 拼接的复合 ID（形如 `zone-xxx#func-yyy`）。资源 SHALL 支持使用该联合 ID 进行 `terraform import`。

#### Scenario: Create 设置复合 ID
- **WHEN** `CreateFunction` 返回非空的 `FunctionId`
- **THEN** 资源在 Create 流程末尾调用 `d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))`，随后调用 Read 回填。

#### Scenario: 通过联合 ID 导入
- **WHEN** 用户执行 `terraform import tencentcloud_teo_function_v4.example zone-xxx#func-yyy`
- **THEN** 资源从复合 ID 解析出 `zone_id` 与 `function_id`，调用 `DescribeFunctions` 回填 state，无需手动拆分。

#### Scenario: 复合 ID 损坏时报错
- **WHEN** `d.Id()` 无法按 `tccommon.FILED_SP` 拆分为两段
- **THEN** Read/Update/Delete 返回明确的 `id is broken` 错误，而非继续解引用空切片。

### Requirement: Create 调用 CreateFunction
Create SHALL 构造 `CreateFunctionRequest`（`ZoneId`、`Name`、`Content` 必填；`Remark` 在用户配置时填入），并在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `CreateFunctionWithContext`；失败使用 `tccommon.RetryError(e)` 包装。

#### Scenario: 创建成功
- **WHEN** `CreateFunctionWithContext` 成功且返回非空 `FunctionId`
- **THEN** 资源设置复合 ID 并调用 Read，返回无错误。

#### Scenario: 空 FunctionId 防御
- **WHEN** `CreateFunctionWithContext` 返回 `nil` Response、`nil` FunctionId 或空字符串 FunctionId
- **THEN** 资源打印 `logId` 与 `d.Id()` 后返回 `tccommon.NonRetryableError`，而非设置空 id。

#### Scenario: SDK 瞬时错误重试
- **WHEN** `CreateFunctionWithContext` 返回可重试的腾讯云 SDK 错误
- **THEN** 该错误通过 `tccommon.RetryError(e)` 在 retry 预算内重试，直至成功或预算耗尽。

### Requirement: Read 调用 DescribeFunctions 并拆分 name
Read SHALL 调用新增的服务层辅助 `TeoService.DescribeTeoFunctionV4ById(ctx, zoneId, functionId)`，该辅助：
- 构造 `DescribeFunctionsRequest`，设置 `ZoneId = zoneId`、`FunctionIds = []*string{functionId}`。
- 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中调用 `DescribeFunctions`。
- 当 `response.Response == nil || len(response.Response.Functions) < 1` 时返回 `(nil, nil)`。
- 严格相等校验 `*function.FunctionId == functionId` 后取第一项。

当辅助返回 `(nil, nil)` 时，Read SHALL 先 `log.Printf("[CRUD] teo_function_v4 id=%s", d.Id())` 保留现场，再 `d.SetId("")`，返回无错误。

Read SHALL 对每个响应字段先判 nil 再 `_ = d.Set(...)`；并对返回的拼接 `name` 调用 `parseTeoFunctionV4OriginalName(name, zoneId)` 拆分后仅 set 原始 name。

#### Scenario: 资源存在
- **WHEN** 辅助找到匹配的 `Function`
- **THEN** Read 填充所有输入字段与 computed 字段（`function_id`、`zone_id`、`name`（已拆分）、`remark`、`content`、`domain`、`create_time`、`update_time`）。

#### Scenario: 资源已被外部删除
- **WHEN** 辅助返回 `(nil, nil)`
- **THEN** Read 打印 `[CRUD]` 日志后调用 `d.SetId("")`，返回无错误。

#### Scenario: 拼接 name 拆分
- **WHEN** `DescribeFunctions` 返回的 `Name` 为 `my-func-zone-2qtuhspy7cr6-1310708577`，`zone_id` 为 `zone-2qtuhspy7cr6`
- **THEN** 设置到 state 的 `name` 值 SHALL 为 `my-func`。

#### Scenario: 拼接 name 不含后缀时保留原值
- **WHEN** `DescribeFunctions` 返回的 `Name` 不包含 `-zoneId` 后缀
- **THEN** 设置到 state 的 `name` 值 SHALL 保留原始返回值，不做截断。

### Requirement: name 拆分辅助函数
资源 SHALL 提供 `parseTeoFunctionV4OriginalName(name, zoneId string) string` 辅助函数，算法与既有 `ParseTeoFunctionOriginalName` 一致：以 `-` + `zoneId` 作为后缀，在 `name` 中定位最后一次出现的位置作为分割点，返回分割点之前的子串；找不到后缀时返回原值。

#### Scenario: 提取原始 name
- **WHEN** 调用 `parseTeoFunctionV4OriginalName("my-zone-func-zone-2qtuhspy7cr6-1310708577", "zone-2qtuhspy7cr6")`
- **THEN** 返回 `my-zone-func`。

#### Scenario: 不含后缀
- **WHEN** 调用 `parseTeoFunctionV4OriginalName("myfunc", "zone-2qtuhspy7cr6")`
- **THEN** 返回 `myfunc`。

#### Scenario: 空字符串
- **WHEN** 调用 `parseTeoFunctionV4OriginalName("", "zone-2qtuhspy7cr6")`
- **THEN** 返回空字符串 `""`。

### Requirement: Update 调用 ModifyFunction 并校验不可变字段
Update SHALL：
- 校验 `immutableArgs`（包含 `name`）：若任一不可变字段发生变更，返回 `argument %s cannot be changed` 错误。
- 当 `mutableArgs`（`remark`、`content`）任一发生变更时，构造 `ModifyFunctionRequest`（`ZoneId`、`FunctionId` 从复合 ID 解析；`Remark`、`Content` 在用户配置时填入），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `ModifyFunctionWithContext`，失败使用 `tccommon.RetryError(e)` 包装。
- retry 块内只执行接口调用，不设置 id 等成功操作。
- 完成后调用 Read 回填。

#### Scenario: 仅可变字段变更
- **WHEN** 仅 `remark` 与 `content` 发生变更
- **THEN** 资源调用 `ModifyFunction` 传入新的 `Remark`、`Content`，随后 Read 回填，返回无错误。

#### Scenario: 不可变字段变更被拒绝
- **WHEN** `name` 发生变更
- **THEN** 资源返回 `argument name cannot be changed` 错误，不调用 `ModifyFunction`。

#### Scenario: 无字段变更跳过调用
- **WHEN** 没有任何可变字段变更
- **THEN** 资源不调用 `ModifyFunction`，直接返回（或仅调用 Read）。

### Requirement: Delete 调用 DeleteFunction
Delete SHALL 从复合 ID 解析 `zoneId`、`functionId`，构造 `DeleteFunctionRequest`（`ZoneId`、`FunctionId`），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `DeleteFunctionWithContext`，失败使用 `tccommon.RetryError(e)` 包装。

#### Scenario: 删除成功
- **WHEN** `DeleteFunctionWithContext` 成功
- **THEN** 资源返回无错误，Terraform 标记资源已销毁。

#### Scenario: SDK 瞬时错误重试
- **WHEN** `DeleteFunctionWithContext` 返回可重试错误
- **THEN** 该错误在 retry 预算内重试，直至成功或预算耗尽。

### Requirement: retry 覆盖
所有 SDK 调用（`CreateFunctionWithContext`、`DescribeFunctions`、`ModifyFunctionWithContext`、`DeleteFunctionWithContext`）SHALL 在 `resource.Retry` 块内执行。retry 预算：写操作使用 `tccommon.WriteRetryTimeout`，读操作使用 `tccommon.ReadRetryTimeout`。失败统一通过 `tccommon.RetryError(e)` 包装。retry 块内不得再次嵌套 retry，也不得执行设置 id 等成功操作。

#### Scenario: 瞬时错误重试
- **WHEN** 任一 SDK 调用返回可重试错误
- **THEN** 在对应 retry 预算内重试，直至成功或预算耗尽。

#### Scenario: 非可重试错误立即抛出
- **WHEN** 任一 SDK 调用返回非可重试错误
- **THEN** `tccommon.RetryError(e)` 将其作为非可重试错误返回，retry 立即终止。

### Requirement: 日志约定
每个 CRUD 函数 SHALL：
- 顶部 `defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v4.<op>")()`。
- 顶部 `defer tccommon.InconsistentCheck(d, meta)()`。
- 每次 SDK 调用成功后打印 `[DEBUG]` 行（含 action、request body、response body）。
- 每个 retry 块失败打印 `[CRITAL]%s ... failed, reason:%+v` 行。
- Read 中资源被外部删除时打印 `[CRUD] teo_function_v4 id=%s` 行后再 `d.SetId("")`。

日志与错误描述统一使用资源名 `teo_function_v4`（小写蛇形），禁止使用"该资源""目标资源"等模糊措辞。

#### Scenario: 标准日志输出
- **WHEN** 任一 CRUD 操作执行
- **THEN** 通过 `tccommon.LogElapsed` 记录耗时、`tccommon.InconsistentCheck` 检查一致性，并在 SDK 调用处输出 `[DEBUG]` 行。

### Requirement: 服务层查询辅助函数
在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中 SHALL 新增 `DescribeTeoFunctionV4ById(ctx, zoneId, functionId string) (*teo.Function, error)`，实现与 `DescribeTeoFunctionById` 等价的行为（`FunctionIds` 精确查询、`ReadRetryTimeout` retry、空返回 `(nil, nil)`、严格相等校验）。

#### Scenario: 辅助函数精确查询
- **WHEN** 调用 `DescribeTeoFunctionV4ById(ctx, "zone-xxx", "func-yyy")`
- **THEN** 该函数构造 `FunctionIds=["func-yyy"]`、`ZoneId="zone-xxx"` 的 `DescribeFunctionsRequest`，在 retry 内调用并返回匹配的 `*Function`。

#### Scenario: 未找到返回空
- **WHEN** `DescribeFunctions` 返回空 `Functions` 列表
- **THEN** 辅助函数返回 `(nil, nil)`，不返回 error。

### Requirement: 文档与测试
本次变更 SHALL 包含：
- 文档 `tencentcloud/services/teo/resource_tc_teo_function_v4.md`：一句话描述（含云产品名称 TEO），`Example Usage`（HCL 示例展示 `zone_id`、`name`、`content`、`remark`），`Import` 部分说明需使用联合 ID `zone_id#function_id`。不手写 `Argument Reference` / `Attribute Reference`（由工具自动生成）。
- 单元测试 `tencentcloud/services/teo/resource_tc_teo_function_v4_test.go`：使用 gomonkey mock 云 API 客户端方法（不使用 terraform 测试套件），覆盖 Create/Read/Update/Delete 业务逻辑分支及 `parseTeoFunctionV4OriginalName` 辅助函数；禁止通过 `go test` 执行，但代码须在当前环境下可正确构建执行。

#### Scenario: 文档存在
- **WHEN** 变更合并
- **THEN** `resource_tc_teo_function_v4.md` 存在，包含 HCL 示例与使用联合 ID 的 `import` 示例。

#### Scenario: 测试文件存在
- **WHEN** 变更合并
- **THEN** `resource_tc_teo_function_v4_test.go` 存在，使用 gomonkey mock，而非 `resource.TestCase` 验收测试套件。

### Requirement: SDK 约束
实现 SHALL NOT 修改 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/` 下任何文件。若缺少所需接口，实现者 MUST 停止并请求 SDK 升级，而非自行修补 SDK 源码。

#### Scenario: vendor SDK 已足够
- **WHEN** 实现开始
- **THEN** 已在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` 中确认 `CreateFunction`、`DescribeFunctions`、`ModifyFunction`、`DeleteFunction` 四个接口及其入参/出参结构体存在。
