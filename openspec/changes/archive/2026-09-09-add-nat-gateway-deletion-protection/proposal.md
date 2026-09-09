## Why

The `tencentcloud_nat_gateway` resource already defines `deletion_protection_enabled` in its schema and supports updating it via `ModifyNatGatewayAttribute`, but the Create API (`CreateNatGateway`) now also accepts `DeletionProtectionEnabled` as an input parameter. Currently, the Create function sets this field on the request but the parameter was originally added without being explicitly wired to the Create API's input. This change ensures that `deletion_protection_enabled` is properly passed to `CreateNatGateway` during resource creation, allowing users to enable deletion protection at creation time rather than requiring a separate update.

## What Changes

- Wire the `deletion_protection_enabled` schema parameter to the `CreateNatGateway` API request's `DeletionProtectionEnabled` field during resource creation (the schema and update/read logic already exist)

## Capabilities

### New Capabilities

_None_

### Modified Capabilities

- `nat-gateway-deletion-protection`: Adding creation-time support for `deletion_protection_enabled` parameter via `CreateNatGateway` API

## Impact

- `tencentcloud/services/vpc/resource_tc_nat_gateway.go`: Modify `resourceTencentCloudNatGatewayCreate` to ensure `deletion_protection_enabled` is passed to the `CreateNatGateway` request
- `tencentcloud/services/vpc/resource_tc_nat_gateway_test.go`: Add unit test coverage for the creation-time `deletion_protection_enabled` parameter
- `tencentcloud/services/vpc/resource_tc_nat_gateway.md`: Update documentation example to show `deletion_protection_enabled` usage
