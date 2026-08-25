## Why

云数据库 Redis（CRS）目前缺少通过 Terraform 管理实例密码复杂度策略的能力。腾讯云 Redis 已提供 `DescribeInstancePasswordPolicy` 与 `ModifyInstancePasswordPolicy` 两个云 API 用于查询和修改存量实例的密码复杂度配置（包含开关、字母/数字/特殊字符最小数量、密码最小总长度），但 Terraform Provider 尚未封装该能力。用户若需统一管控实例密码安全基线，必须手动在控制台操作，无法纳入基础设施即代码流程。

新增 `tencentcloud_redis_instance_password_policy` 资源（RESOURCE_KIND_CONFIG）可填补这一空白，使 Redis 实例的密码复杂度配置可以被声明式地管理和复用。

## What Changes

- 新增资源 `tencentcloud_redis_instance_password_policy`（配置类资源，RESOURCE_KIND_CONFIG），文件名为 `resource_tc_redis_instance_password_policy_config.go`。
- 资源仅提供 **RU** 行为（配置的读取与更新），复用现有实例：
  - **Create**：调用 `ModifyInstancePasswordPolicy` 写入配置，使用实例 ID 作为资源 ID（配置依附于实例存在）。
  - **Read**：调用 `DescribeInstancePasswordPolicy` 读取当前密码复杂度配置并回填 state。
  - **Update**：检测到字段变更时调用 `ModifyInstancePasswordPolicy` 更新配置，随后回读校验。
  - **Delete**：配置依附于实例，删除资源不做实际销毁操作（保留实例，仅从 Terraform state 移除）。
- 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册新资源。
- 新增文档文件 `resource_tc_redis_instance_password_policy_config.md`（由 `make doc` 生成）。
- 新增单元测试文件 `resource_tc_redis_instance_password_policy_config_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

## Capabilities

### New Capabilities
- `redis-instance-password-policy`: 管理 Redis 实例的密码复杂度配置策略，支持读取（DescribeInstancePasswordPolicy）与更新（ModifyInstancePasswordPolicy）实例级密码复杂度规则（开关、字母/数字/特殊字符最小数量、密码最小总长度）。

### Modified Capabilities
- 无。本次变更全部为新增资源，不修改既有 capability 的需求。

## Impact

- **新增代码**：
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.go`（资源实现）
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config_test.go`（单元测试）
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.md`（文档，由 make doc 生成）
- **修改代码**：
  - `tencentcloud/provider.go`：在 redis 资源组中注册 `tencentcloud_redis_instance_password_policy`。
- **依赖的云 API**（已在 vendor 中验证存在，包名 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`）：
  - `DescribeInstancePasswordPolicy`：入参 `InstanceId`；出参 `PasswordPolicy{Enabled, MinLetterCount, MinDigitCount, MinSpecialCount, MinLength}`。
  - `ModifyInstancePasswordPolicy`：入参 `InstanceId`、`PasswordPolicy{Enabled(必填), MinLetterCount, MinDigitCount, MinSpecialCount, MinLength}`。
- **兼容性**：纯新增资源，不改变既有资源 schema 与 state，完全向后兼容。
- **文档**：由收尾阶段 `make doc` 自动生成，不在本变更阶段直接新增 website/ 文件。
