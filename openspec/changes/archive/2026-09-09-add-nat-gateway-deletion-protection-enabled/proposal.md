## Why

NAT网关实例支持开启删除保护功能，开启后可以防止误删除操作。当前 `tencentcloud_nat_gateway` 资源未支持此参数，用户无法通过 Terraform 管理 NAT 网关的删除保护开关。

## What Changes

- 为 `tencentcloud_nat_gateway` 资源新增 `deletion_protection_enabled` 参数（可选，bool 类型）
- 创建 NAT 网关时支持传入 `deletion_protection_enabled` 参数
- 更新 NAT 网关时支持修改 `deletion_protection_enabled` 参数
- 读取 NAT 网关时返回 `deletion_protection_enabled` 参数值

## Capabilities

### New Capabilities
- `nat-gateway-deletion-protection`: 管理 NAT 网关的删除保护开关，支持创建时设置、更新时修改、读取时获取状态

### Modified Capabilities
- 无（仅修改现有资源参数，不涉及已有能力变更）

## Impact

- 修改文件：`tencentcloud/services/vpc/resource_tc_nat_gateway.go`
- 新增字段：`deletion_protection_enabled`，类型为 `schema.TypeBool`，`Optional` + `Computed`
- 涉及云API接口：
  - `CreateNatGateway`：入参新增 `DeletionProtectionEnabled`
  - `ModifyNatGatewayAttribute`：入参新增 `DeletionProtectionEnabled`
  - `DescribeNatGateways`：出参 `NatGateway.DeletionProtectionEnabled`
- 保持向后兼容，不影响现有配置和 state