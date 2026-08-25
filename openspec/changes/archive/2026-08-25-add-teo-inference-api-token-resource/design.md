## Context

腾讯云 TEO 提供"推理 API Token"用于访问 TEO 托管的推理服务。该能力对应的云 API（`teo/v20220901`）仅提供以下三个接口：

- `CreateInferenceAPIToken`：入参 `ZoneId`、`Name`；出参 `TokenId`、`Content`
- `DescribeInferenceAPITokens`：入参 `ZoneId`、`Offset`、`Limit`；出参 `TotalCount`、`Tokens[]`（每项含 `TokenId`、`Name`、`Content`、`CreateTime`）
- `DeleteInferenceAPIToken`：入参 `ZoneId`、`TokenId`

**关键事实**：不存在 `ModifyInferenceAPIToken` 接口。`ModifyInferenceService` 是针对"推理服务"（service）的接口，其字段（`ServiceId`、`ListenPort`、`Containers` 等）与"推理 API Token"无关，不能用于更新 Token。因此本资源只支持 CRD，不支持原地更新。

当前状态：
- provider 中无 `tencentcloud_teo_inference_api_token` 资源
- SDK 中已 vendored 上述三个接口的 Request/Response 结构体
- `DescribeInferenceAPITokens` 返回的是数组，无单条查询接口，需通过 `TokenId` 在列表中匹配

约束：
- 复合 ID：资源的唯一标识需要 `ZoneId` + `TokenId`（删除接口必须同时提供两者），采用 `tccommon.FIELD_SP`（`#`）分隔。
- 无独立更新接口：依据 provider 规范，仅 CRD 的资源需把 `Id()` 字段设为 ForceNew，并把其余顶层字段加入 `immutableArgs`，在 Update 方法中检测到非法变更时返回 error。
- Read 接口为列表查询：需要在 service 层封装一个按 `TokenId` 过滤的方法，遍历 `Tokens` 找到目标项。
- 错误处理遵循 provider 统一模式：`resource.Retry` + `tccommon.RetryError`，超时使用 `tccommon.ReadRetryTimeout`（读）/ `tccommon.WriteRetryTimeout`（写）。

## Goals / Non-Goals

**Goals:**
- 提供声明式创建、查询、删除 TEO 推理 API Token 的能力
- 通过 `immutableArgs` 机制正确处理"无更新接口"的约束（变更 `name` 触发重建，而非调用不存在的更新接口）
- 复合 ID 设计使 Delete 能还原出 `ZoneId` 与 `TokenId`
- 单测基于 gomonkey mock 云 API，覆盖 Create/Read/Delete 及不可变字段校验分支

**Non-Goals:**
- 不实现原地更新（云 API 不支持）
- 不引入 `ModifyInferenceService`（与 Token 无关，会破坏语义）
- 不新增数据源（`tencentcloud_teo_inference_api_tokens` 数据源不在本次范围）
- 不暴露 `Offset`/`Limit` 分页参数给用户（Read 内部固定 `Limit=100` 一次拉取并在结果中匹配）

## Decisions

### Decision 1: 资源类型为 CRD-only，不做原地更新

**选择**：仅实现 Create / Read / Delete；`name` 字段标记为 ForceNew，并在 Update 的 `immutableArgs` 数组中加入 `name`，命中即返回 error。

**备选**：强行调用 `ModifyInferenceService` 假装更新。

**理由**：
- 云 API 没有 `ModifyInferenceAPIToken`，`ModifyInferenceService` 的字段（`ServiceId`、`Containers` 等）与 Token 语义完全不匹配，强行调用会触发服务端错误或修改到无关资源
- 遵循 provider 既定规范："若一个资源只有 CRD 接口，则只将 Id() 字段设置成 ForceNew，并在资源 update 方法中将其余顶层字段加入 immutableArgs 数组"
- 用户如需修改 `name`，重建资源是唯一正确语义

### Decision 2: 复合 ID 使用 `ZoneId # TokenId`

**选择**：`d.SetId(fmt.Sprintf("%s%s%s", zoneId, tccommon.FILED_SP, tokenId))`，分隔符为 `tccommon.FILED_SP`（`#`）。

**理由**：
- Delete 接口同时需要 `ZoneId` 与 `TokenId`，单字段 ID 无法承载
- 与 provider 规范一致（复合 ID 使用 `tccommon.FILED_SP`）
- Read / Delete 方法中均通过 `d.Id()` 解析出两个字段，不依赖 schema 中其他字段是否被清空

### Decision 3: Read 通过 `DescribeInferenceAPITokens` 列表匹配

**选择**：在 service 层封装 `DescribeTeoInferenceApiTokenById(zoneId, tokenId)`，内部调用 `DescribeInferenceAPITokens`（`Offset=0`、`Limit=100`），遍历 `Tokens` 按 `TokenId` 精确匹配返回。

**理由**：
- 无单条查询接口，列表查询是唯一途径
- `Limit=100` 为云 API 注释标注的最大值，单站点 Token 数量预期不会超过该值
- 匹配失败时返回 `dataNotFound` 错误，由 Read 方法触发 `d.SetId("")`（先打印 `[CRUD]` 日志保留 id 现场）

### Decision 4: Create 返回值必须校验非空

**选择**：`CreateInferenceAPIToken` 调用后，校验 `response == nil` / `response.Response == nil` / `response.Response.TokenId == nil` / `*response.Response.TokenId == ""`，任一为空则返回 `NonRetryableError`。

**理由**：
- 遵循 provider 规范第 9 条：必须检查返回值是否为空，避免写入空 id 触发后续状态混乱
- 校验前打印 `logId` 与相关信息便于排障

### Decision 5: `token_id` 字段处理

**选择**：`token_id` 同时声明为 `Computed`（由 Create/Read 回填）并参与复合 ID。`zone_id`、`name` 声明为 `Required` + ForceNew。

**理由**：
- `token_id` 是云 API 生成的唯一标识，用户不应主动填写，故为 Computed
- 由于资源只有 CRD，`token_id` 一旦确定即不可变（重建才会生成新 id），ForceNew 语义天然满足

## Risks / Trade-offs

- **Risk**：单站点 Token 数量超过 100 时，`DescribeInferenceAPITokens` 一次拉取不全导致 Read 找不到目标 → **Mitigation**：以云 API 最大 `Limit=100` 拉取；当前产品阶段 Token 数量远小于该值，属可接受范围；如未来超出可改为分页累加。
- **Risk**：用户误改 `name` 触发重建导致现有 Token 被删除再创建 → **Mitigation**：这是 CRD-only 资源的正确语义，文档中明确说明 `name` 不可原地修改。
- **Trade-off**：不暴露分页参数给用户，Read 只取前 100 条并匹配 → 简化用户接口，符合 provider 数据源/资源分页内部处理约定。
- **Risk**：`content` 是敏感信息（Token 明文），明文写入 state 可能存在泄露风险 → **Mitigation**：该字段为云 API 出参，按现有资源惯例（如其他 token 类资源）以 Computed 存放；用户可结合 Terraform `sensitive` 策略进一步控制，本次不引入额外复杂度。
