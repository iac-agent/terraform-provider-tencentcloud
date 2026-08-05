## Context

TEO 推理服务提供 `CreateInferenceAPIToken`、`DescribeInferenceAPITokens`、`DeleteInferenceAPIToken` 三个云 API 用于管理推理 API Token。当前 Terraform Provider 尚未封装这些接口为 Terraform 资源。本设计新增 `tencentcloud_teo_inference_api_token` 资源（RESOURCE_KIND_GENERAL），覆盖 CRD 生命周期。

- 云 API 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（vendor 中已存在）
- 参考资源：`tencentcloud_teo_content_identifier`（CRUD 模式）、`tencentcloud_teo_ownership_verify`（CRD 模式）
- 资源路径：`tencentcloud/services/teo/`

## Goals / Non-Goals

**Goals:**
- 支持通过 Terraform 创建 TEO 推理 API Token（`CreateInferenceAPIToken`）
- 支持通过 Terraform 读取 TEO 推理 API Token（`DescribeInferenceAPITokens`）
- 支持通过 Terraform 删除 TEO 推理 API Token（`DeleteInferenceAPIToken`）
- Token 内容（`content`）仅在创建时返回，需标记为 `Sensitive`

**Non-Goals:**
- 不支持修改/更新 Token（云 API 未提供 Update/Modify 接口）
- 不做 Token 的导入（import）支持（除非有实际需求）
- 不创建 datasource 资源

## Decisions

### 1. 资源 ID 使用 TokenId

`CreateInferenceAPIToken` 返回 `TokenId`（String 类型），将此值作为 Terraform 资源的 ID。
- 理由：`TokenId` 是唯一标识，`DeleteInferenceAPIToken` 和 `DescribeInferenceAPITokens` 均依赖它定位资源。

### 2. Schema 设计：所有用户输入字段标记 ForceNew

由于云 API 没有 Update 接口，资源 Schema 中所有 Required/Optional 字段均标记 `ForceNew: true`。资源不注册 Update 函数。
- `zone_id`：Required, ForceNew — 站点 ID
- `name`：Required, ForceNew — Token 名称（限制 ≤30 字符）
- `token_id`：Computed — 创建后由云 API 返回
- `content`：Computed, Sensitive — Token 内容（仅在创建时返回一次）
- `create_time`：Computed — 从 Describe 响应中读取

### 3. Read 实现：通过 DescribeInferenceAPITokens 列表查找

`DescribeInferenceAPITokens` 是列表接口，仅接受 `ZoneId` 作为过滤条件（不含 TokenId 筛选）。Read 时需要：
1. 调用 `DescribeInferenceAPITokens`，传入 `ZoneId`（从 state 读取），设置分页 Limit=100（最大值）
2. 遍历返回的 `Tokens` 列表，匹配 `TokenId == d.Id()`
3. 若匹配成功，将对应字段写入 state
4. 若未找到匹配项，打印警告日志后 `d.SetId("")`

### 4. 调用云 API 的客户端方式

使用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()` 获取 teo 客户端（与现有 teo 资源一致）。

### 5. Read 接口的空结果处理

若 `DescribeInferenceAPITokens` 返回的 `Response` 为 nil 或 `Tokens` 为空（或未找到匹配），按照规范：
- 先打印 `log.Printf("[CRUD] teo_inference_api_token id=%s", d.Id())` 保留现场
- 再 `d.SetId("")`

### 6. Service 层辅助函数

在 `service_tencentcloud_teo.go` 中增加 `DescribeTeoInferenceAPITokenById` 方法，封装 `DescribeInferenceAPITokens` 的调用和 `TokenId` 匹配逻辑，保持资源代码简洁。

## Risks / Trade-offs

- **[风险] Token Content 丢失**：`content` 仅在 `CreateInferenceAPIToken` 响应中返回一次。若 state 丢失（如 drift），Read 时虽然 `DescribeInferenceAPITokens` 的 `InferenceAPIToken` 结构也包含 Content 字段，但云 API 是否返回需要验证。
  → **缓解**：若 Read 时 `Content` 为 nil 则不覆盖已有值，避免空覆盖。
  
- **[风险] DescribeInferenceAPITokens 分页**：若用户在同一 Zone 下创建超过 100 个 Token，单次分页查询可能遗漏目标。但云 API 文档标注最大值 100，且每个站点最多创建 100 个 Token（根据 client.go 注释），单次查询即可覆盖全部。
  → **缓解**：Limit 固定为 100，满足当前上限。

- **[权衡] 不支持 Import**：由于 Read 需要 `zone_id` 才能调用 `DescribeInferenceAPITokens`，而 Import 时只有 `token_id`，无法获取 `zone_id`。若未来需要 support import，需考虑使用联合 ID（`zone_id#token_id`）的方案。
  → 当前不实现 Import，保持简洁。
