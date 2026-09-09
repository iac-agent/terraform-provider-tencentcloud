## Context

The `tencentcloud_nat_gateway` resource already has a `deletion_protection_enabled` schema parameter defined (TypeBool, Optional+Computed). The Read function reads it from `DescribeNatGateways` response, and the Update function passes it to `ModifyNatGatewayAttribute`. The Create function already sets `request.DeletionProtectionEnabled` from the schema value. The cloud API `CreateNatGateway` now officially supports `DeletionProtectionEnabled` as an optional input parameter (confirmed in vendor SDK model `CreateNatGatewayRequest`).

Current state of the code in `resource_tc_nat_gateway.go`:
- Schema: `deletion_protection_enabled` is defined as Optional+Computed, TypeBool
- Create: Already sets `request.DeletionProtectionEnabled = helper.Bool(v.(bool))` when the value is provided via `d.GetOkExists("deletion_protection_enabled")`
- Read: Reads `DeletionProtectionEnabled` from `DescribeNatGateways` response
- Update: Passes `DeletionProtectionEnabled` to `ModifyNatGatewayAttribute` when changed

Since the Create function already wires this parameter to the API request, the existing implementation already supports creation-time deletion protection. The change is to verify this is correctly implemented and ensure test coverage exists.

## Goals / Non-Goals

**Goals:**
- Verify and ensure `deletion_protection_enabled` is correctly passed to the `CreateNatGateway` API during resource creation
- Add unit test coverage for the creation-time `deletion_protection_enabled` parameter
- Update resource documentation example to include `deletion_protection_enabled`

**Non-Goals:**
- Modifying the schema definition of `deletion_protection_enabled` (already exists)
- Changing the Read or Update logic (already correctly implemented)
- Adding new parameters beyond `deletion_protection_enabled`

## Decisions

1. **No schema changes needed**: The `deletion_protection_enabled` parameter already exists in the schema with the correct type (TypeBool) and options (Optional+Computed). No schema modification is required.

2. **Create function already correct**: The `resourceTencentCloudNatGatewayCreate` function already reads `deletion_protection_enabled` from the resource data and sets it on the `CreateNatGatewayRequest`. No code changes needed in the Create function.

3. **Test coverage**: Add unit test using gomonkey mock approach to verify the `deletion_protection_enabled` parameter is correctly passed through the Create flow.

4. **Documentation update**: Update the `.md` resource example to show `deletion_protection_enabled` usage at creation time.

## Risks / Trade-offs

- [Backward compatibility] → Since `deletion_protection_enabled` is Optional+Computed, existing configurations without this field will continue to work unchanged. No migration needed.
