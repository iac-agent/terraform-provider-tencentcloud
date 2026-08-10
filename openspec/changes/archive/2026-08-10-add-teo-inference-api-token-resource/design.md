## Context

TEO (Tencent EdgeOne) 推理 API Token 是管理 EdgeOne 边缘推理服务认证凭据的资源。云API 提供了三个接口：`CreateInferenceAPIToken`（创建）、`DescribeInferenceAPITokens`（查询列表）、`DeleteInferenceAPIToken`（删除）。云API 未提供更新接口，因此该资源仅支持 CRD（Create/Read/Delete）操作。

当前 Terraform Provider 已有 `tencentcloud/services/teo/` 目录，包含多个 TEO 资源，代码风格已成熟。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_inference_api_token_v9` 资源，支持创建、查询和删除推理 API Token
- 遵循现有 TEO 资源的代码风格和模式
- 所有字段因云 API 不支持更新而标记为 ForceNew

**Non-Goals:**
- 不支持更新推理 API Token（云 API 无更新接口）
- 不实现 Import 功能（Token 创建后 Content 仅返回一次，Import 无法恢复 Content）

## Decisions

### 1. 资源 ID 设计
- **决策**: 使用 `TokenId` 作为资源 ID（`d.SetId(tokenId)`）
- **理由**: `TokenId` 是云 API 返回的唯一标识，且 Delete 接口仅需 `ZoneId` + `TokenId`
- **替代方案**: 使用 `ZoneId#TokenId` 联合 ID。但 `TokenId` 本身已全局唯一，无需联合 ID。且联合 ID 会增加 Read 和 Delete 时解析的复杂度。

### 2. Schema 字段设计
- **决策**: 所有顶层字段标记为 `ForceNew: true`
- **理由**: 云 API 无 Update 接口，任何字段变更都需要重建资源
- **字段映射**:
  - `zone_id` (Required, ForceNew): → Create.ZoneId / Describe.ZoneId / Delete.ZoneId
  - `name` (Required, ForceNew): → Create.Name
  - `token_id` (Computed): → CreateResponse.TokenId
  - `content` (Computed, Sensitive): → CreateResponse.Content（创建时一次性返回，后续查询不会返回明文）

### 3. Read 实现策略
- **决策**: 调用 `DescribeInferenceAPITokens` 传入 `ZoneId`，设置 `Limit=100`（最大值），遍历返回的 `Tokens` 列表，按 `TokenId` 匹配
- **理由**: Describe 接口不支持按 TokenId 过滤，只能通过列表遍历匹配。这是 TEO 多个资源的标准做法。
- **注意**: Read 时不设置 `content` 字段（云 API 列表查询不返回 Content 明文），Content 仅在 Create 时返回一次

### 4. 不使用 Update 方法
- **决策**: 不定义 Update 方法，资源仅包含 Create/Read/Delete
- **理由**: 云 API 无 Update 接口，且所有字段均为 ForceNew，Terraform 检测到变更时会触发 destroy + create

### 5. 不使用 Import
- **决策**: 不支持 Import
- **理由**: Token 的 `content` 仅在创建时返回一次，Import 无法恢复 Content 字段，导致 state 不完整

## Risks / Trade-offs

- **[风险] Content 字段在 Create 后仅返回一次**: Read 操作无法恢复 Content 值。如果 state 丢失，Content 将永久丢失。→ 缓解：标记 Content 为 Sensitive，在 Read 中不覆盖 Content（保留 state 中已有值），仅在 d.SetId("") 时清空。
- **[风险] DescribeInferenceAPITokens 是分页接口**: 如果 Token 数量超过 100 个，当前分页查询可能找不到目标 Token。→ 缓解：设置 Limit=100（最大值），当前 EdgeOne 推理 API Token 数量通常远小于此上限。

## Open Questions

- 无。云 API 接口已在 vendor 中验证，参数映射清晰。