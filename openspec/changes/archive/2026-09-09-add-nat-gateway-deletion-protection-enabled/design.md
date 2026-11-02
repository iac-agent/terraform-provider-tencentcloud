## Context

当前 `tencentcloud_nat_gateway` 资源已支持 VPC 下 NAT 网关的完整生命周期管理，包括创建、读取、更新、删除操作。NAT 网关的删除保护功能（`DeletionProtectionEnabled`）在云 API 中已支持，但 Terraform 资源尚未暴露此参数。用户无法通过 Terraform 控制 NAT 网关的删除保护状态。

涉及云 API 接口：
- `CreateNatGateway`：支持 `DeletionProtectionEnabled` 入参
- `ModifyNatGatewayAttribute`：支持 `DeletionProtectionEnabled` 入参
- `DescribeNatGateways`：返回 `NatGateway.DeletionProtectionEnabled` 字段

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_nat_gateway` 资源新增 `deletion_protection_enabled` 可选参数
- 支持在创建 NAT 网关时设置删除保护
- 支持在更新 NAT 网关时修改删除保护
- 支持在读取时获取删除保护状态并写入 state
- 保持向后兼容，不破坏现有用户配置

**Non-Goals:**
- 不涉及 NAT 网关其他参数的修改
- 不涉及删除逻辑的变更（删除保护仅阻止删除操作，由云API控制，Terraform侧无需额外处理）
- 不涉及数据源 `tencentcloud_nat_gateways` 的修改

## Decisions

1. **参数类型：`Optional` + `Computed`**
   - 选择 `Optional` + `Computed` 而非 `Optional` + `Default(false)`，因为云 API 创建时若不传此参数，服务端会使用默认值。使用 `Computed` 可在读取时获取实际值，避免用户配置中未设置时产生 diff。
   
2. **在 Create 中传入**
   - 使用 `d.GetOkExists("deletion_protection_enabled")` 判断用户是否设置，设置后才传入请求，避免不必要的 API 参数传递。

3. **在 Update 中传入**
   - 使用 `d.HasChange("deletion_protection_enabled")` 检测变更，仅在变更时通过 `ModifyNatGatewayAttribute` 接口更新。

4. **在 Read 中读取**
   - 从 `NatGateway.DeletionProtectionEnabled` 字段读取，先判断非 nil 后再设置到 state。

## Risks / Trade-offs

- **[风险] 删除保护开启后， terraform destroy 会失败** → 用户需先通过 Terraform 将 `deletion_protection_enabled` 设置为 `false` 后再执行 destroy，或手动在控制台关闭删除保护。这是云 API 的预期行为，文档中需说明。
- **[风险] 兼容性** → 新增 `Optional` + `Computed` 参数不会破坏现有配置，已有 state 在升级后通过 read 操作自动填充该值。