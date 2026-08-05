## Context

腾讯云 EdgeOne (TEO) 推理服务提供了 API Token 管理功能，用户需要通过 API Token 来使用推理 API。当前 SDK 已支持 `CreateInferenceAPIToken`、`DescribeInferenceAPITokens` 和 `DeleteInferenceAPIToken` 三个 API，但 Terraform Provider 尚未提供对应的资源。

**Current State:**
- TEO 服务已有多个 Terraform 资源（zone、origin_group、l4_proxy 等），遵循统一的代码模式
- SDK (`tencentcloud-sdk-go/tencentcloud/teo/v20220901`) 已包含推理 API Token 相关接口，且已在 vendor 中
- 现有的 TEO 资源代码风格和模式可作为参考

**Constraints:**
- 该资源只有 CRD 接口（Create/Read/Delete），没有 Update/Modify 接口
- Read 接口（DescribeInferenceAPITokens）是列表接口，需要通过 TokenId 过滤找到目标资源
- Content 字段是敏感信息（Token 内容），需要在 Schema 中标记 `Sensitive: true`
- 必须遵循项目命名规范：`tencentcloud_teo_inference_api_token_v7`
- API 调用必须使用 `tccommon.ReadRetryTimeout` 超时，并使用 retry 逻辑

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_inference_api_token_v7` 资源，支持创建、读取和删除推理 API Token
- 提供完整的 Schema 定义，包含 `zone_id`、`name`、`token_id`、`content` 字段
- 支持 Import 操作（通过 TokenId 导入）
- 添加单元测试（使用 gomonkey mock）和完整文档
- 遵循现有 TEO 资源代码模式

**Non-Goals:**
- 不支持 Update 操作（API 不提供修改接口）
- 不创建数据源资源（本需求仅要求 RESOURCE_KIND_GENERAL 类型）
- 不修改现有 TEO 资源或服务层实现

## Decisions

### Decision 1: 资源 ID 设计
**Choice:** 使用 `TokenId` 作为资源 ID

**Rationale:**
- TokenId 是 API 返回的唯一标识符，在 Create 响应中返回
- 在 Describe 接口中可以通过遍历列表找到匹配的 TokenId
- 简单直接，无需复合 ID

**Alternatives Considered:**
- 使用 `ZoneId#TokenId` 复合 ID：增加复杂度，TokenId 本身已是全局唯一，无需复合

### Decision 2: Schema 设计
```go
{
  "zone_id": {
    Type:     schema.TypeString,
    Required: true,
    ForceNew: true,
  },
  "name": {
    Type:     schema.TypeString,
    Required: true,
    ForceNew: true,
  },
  "token_id": {
    Type:     schema.TypeString,
    Computed: true,
  },
  "content": {
    Type:      schema.TypeString,
    Computed:  true,
    Sensitive: true,
  },
}
```

**Rationale:**
- `zone_id` 和 `name` 是创建 API 的必填参数，标记为 Required 和 ForceNew（无 Update 接口）
- `token_id` 和 `content` 是 API 返回的计算字段，标记为 Computed
- `content` 是 Token 内容，属于敏感信息，标记为 Sensitive
- 所有用户可配置字段都是 ForceNew，因为 API 不支持更新

### Decision 3: Read 操作实现
**Approach:** 在 Read 方法中调用 `DescribeInferenceAPITokens`，遍历 `Tokens` 列表找到匹配的 TokenId

**Details:**
- 使用 `d.Id()` 获取 TokenId
- 从 `d.GetOk("zone_id")` 获取 ZoneId 作为请求参数
- 设置分页参数 `Limit: 100`（最大值），循环遍历直到找到匹配或遍历完所有结果
- 若找到匹配的 Token，设置 `zone_id`、`name`、`token_id`、`content` 到 state
- 若未找到，打印日志后 `d.SetId("")` 清空资源

### Decision 4: 错误处理策略
**Approach:**
- 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用
- 在 API 调用失败时使用 `tccommon.RetryError()` 包装错误
- 在 Create 方法中，检查返回值是否为空（Response 为 nil、TokenId 为 nil 或空字符串），若为空返回 `NonRetryableError`
- 在 Read 方法中，若 API 返回空列表（资源不存在），打印日志并清空 ID
- 在 Delete 方法中，即使 API 返回错误也继续（资源可能已被外部删除）

### Decision 5: Import 支持
**Format:** `terraform import tencentcloud_teo_inference_api_token_v7.foo <token_id>`

**Implementation:**
```go
Importer: &schema.ResourceImporter{
  State: schema.ImportStatePassthrough,
},
```

在 Read 函数中，使用 `d.Id()` 即 TokenId，加上从配置中读取的 `zone_id`，调用 Describe 接口验证资源存在。

### Decision 6: 单元测试
使用 gomonkey 对云 API 进行 mock，不依赖真实的云环境。测试覆盖：
- Create 成功场景
- Create 失败场景（API 返回空 TokenId）
- Read 成功场景
- Read 资源不存在场景
- Delete 成功场景

## Risks / Trade-offs

### Risk 1: Describe 接口是分页列表接口
**Issue:** 如果 Token 数量很多，可能需要多次分页查询才能找到目标 Token
**Mitigation:** 设置 `Limit: 100`（最大值），减少分页次数；在找到目标 Token 后立即返回

### Risk 2: Content 字段在 Read 后可能为空
**Issue:** 部分 API 在 Describe 响应中可能不返回 Content 字段
**Mitigation:** 在 Read 中检查 Content 是否为 nil，仅在非 nil 时设置；由于 Content 是 Computed 且 Sensitive，Terraform 不会因为其变化而触发 drift

### Risk 3: Token 在外部被删除
**Issue:** Token 可能通过控制台或其他方式被删除，导致 Terraform state 不一致
**Mitigation:** Read 操作中检查 Token 是否存在，不存在则清空 state，下次 plan 会检测到 drift 并提示重建

### Trade-off: 不支持 Update
**Reason:** 腾讯云 API 不提供修改推理 API Token 的接口
**Impact:** 用户需要删除并重建 Token 来修改任何属性
**Benefit:** 实现简单，与 API 能力对齐

## Migration Plan

### Deployment Steps:
1. 创建资源文件 `resource_tc_teo_inference_api_token_v7.go`
2. 创建单元测试文件 `resource_tc_teo_inference_api_token_v7_test.go`
3. 创建文档文件 `resource_tc_teo_inference_api_token_v7.md`
4. 在 `provider.go` 和 `provider.md` 中注册新资源
5. 运行 `make doc` 生成 website 文档

### Rollback:
- 如果发现问题，可以在 `provider.go` 中注释掉资源注册
- 不影响现有资源，因为是全新添加

## Open Questions

1. ~~是否需要支持 Update?~~ → 不需要，API 不提供修改接口
2. ~~资源 ID 使用什么格式?~~ → 使用 TokenId 作为 ID
3. ~~是否需要数据源?~~ → 本需求仅要求 RESOURCE_KIND_GENERAL 类型资源，不涉及数据源