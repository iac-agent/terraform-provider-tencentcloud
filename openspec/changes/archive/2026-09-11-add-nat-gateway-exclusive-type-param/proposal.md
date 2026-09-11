## Why

The `tencentcloud_nat_gateway` resource currently supports basic NAT gateway management but lacks support for the exclusive (dedicated) instance type feature. The TencentCloud VPC API already supports creating and managing exclusive NAT gateways via the `ExclusiveType` parameter, which allows users to specify dedicated instance specifications (`ExclusiveSmall`, `ExclusiveMedium1`, `ExclusiveLarge1`). Adding this parameter enables Terraform users to provision and manage exclusive NAT gateways infrastructure-as-code.

## What Changes

- Add a new `exclusive_type` optional parameter to `tencentcloud_nat_gateway` resource
- The parameter maps to `ExclusiveType` field in `CreateNatGatewayRequest`, `NatGateway` (response), and `ResetNatGatewayConnectionRequest`
- Supports three values: `ExclusiveSmall`, `ExclusiveMedium1`, `ExclusiveLarge1`
- Setting this parameter changes the NAT gateway to exclusive (dedicated) mode
- The parameter is `ForceNew` for the Create path, but can be updated via `ResetNatGatewayConnection` for in-place changes to exclusive type

## Capabilities

### New Capabilities
- `exclusive-nat-gateway`: Support for exclusive (dedicated) NAT gateway instance type specification, allowing users to specify dedicated instance tiers through the `exclusive_type` parameter on `tencentcloud_nat_gateway`

### Modified Capabilities

- (none)

## Impact

- **Code**: Modify `tencentcloud/services/vpc/resource_tc_nat_gateway.go` - add `exclusive_type` to schema, create, read, and update logic
- **Documentation**: Update resource documentation in the `.md` file and regenerate website docs via `make doc`
- **API**: Uses existing VPC SDK `ExclusiveType` field (already present in vendor SDK)