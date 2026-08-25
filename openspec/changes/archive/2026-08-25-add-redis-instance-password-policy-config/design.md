## Context

腾讯云 Redis（CRS）提供实例级密码复杂度策略，用于在创建或重置实例密码时强制校验密码的字母、数字、特殊字符最小数量与密码最小总长度。该能力对应两个云 API（包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`，已在 vendor 中验证）：

- `DescribeInstancePasswordPolicy`：入参 `InstanceId`，出参 `PasswordPolicy{Enabled, MinLetterCount, MinDigitCount, MinSpecialCount, MinLength}`。
- `ModifyInstancePasswordPolicy`：入参 `InstanceId`、`PasswordPolicy{Enabled(必填), MinLetterCount, MinDigitCount, MinSpecialCount, MinLength}`。

当前 Terraform Provider 已有大量 Redis 配置类资源（如 `tencentcloud_redis_maintenance_window`、`tencentcloud_redis_connection_config`、`tencentcloud_redis_read_only` 等），均位于 `tencentcloud/services/crs/` 目录，采用「配置依附于实例、以 instance_id 作为资源 ID」的惯例。本变更遵循同一惯例。

两个云 API 均为同步接口，无需轮询。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_redis_instance_password_policy` 资源（RESOURCE_KIND_CONFIG），封装 Redis 实例密码复杂度策略的读取与更新。
- 提供完整的 Create/Read/Update/Delete 行为，并以 `instance_id` 作为资源 ID。
- 严格遵循现有 `crs` 服务的代码风格（参考 `resource_tc_redis_maintenance_window.go` 与 `resource_tc_redis_connection_config.go`）。
- 通过 `tccommon.WriteRetryTimeout` / `tccommon.ReadRetryTimeout` 提供最终一致性重试。
- 保持向后兼容（纯新增）。

**Non-Goals:**
- 不负责实例的创建与生命周期管理（由 `tencentcloud_redis_instance` 资源负责）。
- 不实现删除即重置的策略——密码复杂度策略依附于实例，删除 Terraform 资源仅从 state 移除，不对实例做实际销毁（沿用 `maintenance_window` 的 delete 空实现惯例，避免误操作导致实例配置被破坏）。
- 不暴露 `password_policy` 这一层嵌套 schema 容器；按规则将列表/对象展开为顶层字段（本资源为单对象，直接平铺顶层字段）。

## Decisions

### Decision 1: 资源 ID 使用 instance_id（而非 UUID）

**选择**：资源 ID 直接使用 `instance_id`，Create 时 `d.SetId(instanceId)`。

**理由**：密码复杂度配置与实例一一对应（一个实例只有一份密码策略），且 Read/Update/Delete 接口均以 `InstanceId` 作为关键入参。这与同目录下 `tencentcloud_redis_maintenance_window`、`tencentcloud_redis_connection_config` 的实现完全一致，无需额外生成 UUID，也无需复合 ID 分隔符。

**备选**：使用 `helper.UUIDGenerator()` 生成 ID（如 `bh_auth_mode_config`）。不采用，因为该全局配置无实例维度，而 Redis 密码策略明确归属于某个实例。

### Decision 2: 字段平铺，不引入 `password_policy` 嵌套层

**选择**：将 `Enabled`、`MinLetterCount`、`MinDigitCount`、`MinSpecialCount`、`MinLength` 直接作为 schema 顶层字段（`enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length`），不包裹在 `password_policy` 块内。

**理由**：遵循「禁止创建该资源列表型/对象型嵌套 schema 这一层」的硬约束，使每个字段都可被 Terraform 单独 set/read，符合现有 Redis 配置类资源的扁平化风格。

**说明**：需求映射中提到的 `password_policy` 作为逻辑容器存在，但在 Terraform schema 中以平铺字段体现。

### Decision 3: Delete 为空实现（不重置配置）

**选择**：`resourceTencentCloudRedisInstancePasswordPolicyDelete` 不调用任何云 API，直接返回 nil。

**理由**：云 API 未提供「关闭/重置密码策略」的独立语义（`ModifyInstancePasswordPolicy` 仅修改，且 `Enabled=false` 也是合法配置但属于用户的主动行为）。为避免 `terraform destroy` 误把实例上的安全策略清空，沿用 `maintenance_window` 的空 delete 惯例——删除资源仅从 state 移除，实例配置保持不变。

### Decision 4: Create 复用 Update 逻辑

**选择**：`Create` 在 `d.SetId(instanceId)` 后直接调用 `Update`（其内调用 `ModifyInstancePasswordPolicy` 再回读），与 `maintenance_window` 一致。

**理由**：CONFIG 类资源没有真正的「创建」语义，配置随实例存在；Create 即「首次写入配置」，与 Update 行为完全相同，复用可避免重复代码。

### Decision 5: retry 与错误处理遵循通用规则

**选择**：
- Read 使用 `tccommon.ReadRetryTimeout`，Update/Create 使用 `tccommon.WriteRetryTimeout`。
- 云 API 调用失败时通过 `tccommon.RetryError(e)` 包装返回。
- `d.SetId("")` 前先 `log.Printf("[CRUD] ...")` 保留现场。
- Create 后检查返回值非空。
- 日志统一使用资源名 `redis_instance_password_policy`（蛇形）。

### Decision 6: enabled 设为 Required，其余 min_* 字段设为 Optional

**选择**：`enabled` 为 `Required`（对应云 API 中 `PasswordPolicy.Enabled` 为必填）；`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length` 为 `Optional`（对应云 API 中这些字段为可选）。

**理由**：与云 API 入参约束保持一致。

## Risks / Trade-offs

- **[风险] 删除资源不清空实例配置可能令用户困惑** → 缓解：在文档与代码注释中明确说明该资源删除仅移除 state，不影响实例上的实际策略；这是 CONFIG 类资源与实例绑定的固有特性。
- **[风险] `enabled=false` 与字段未设置难以区分** → 缓解：`enabled` 设为 Required，用户必须显式声明；min_* 字段缺失时不传给云 API（使用 `omitnil`，云 API 端使用默认值）。
- **[风险] Read 时云 API 返回空导致 state 被清空** → 缓解：在 retry 内返回 `NonRetryableError` 由上层处理，并对 `response == nil` / `response.Response == nil` 做判空保护后再 set 字段。
- **[权衡] 平铺字段而非嵌套 `password_policy` 块** → 收益：每个字段独立可读写、风格统一；代价：与云 API 的对象结构在 schema 层不再 1:1，但这是本仓库的既定规范。

## Migration Plan

- 全新增资源，无需数据迁移、无需 state 升级。
- 部署：合并代码并在 `provider.go` 注册后即可使用；用户可通过 `terraform import tencentcloud_redis_instance_password_policy.example <instance_id>` 导入既有配置。
- 回滚：移除 `provider.go` 注册项与资源文件即可，不影响既有资源与 state（已创建的该资源 state 在回滚后需手动 `terraform state rm`）。

## Open Questions

无。云 API 入参/出参与字段约束已在 vendor 中确认，均为同步接口，无需额外异步轮询设计。
