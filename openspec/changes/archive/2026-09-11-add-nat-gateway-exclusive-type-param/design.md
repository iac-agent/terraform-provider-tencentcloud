## Context

The `tencentcloud_nat_gateway` resource currently supports standard NAT gateway management (creation, reading, updating, deletion). The TencentCloud VPC API supports an "exclusive" (dedicated) NAT gateway mode via the `ExclusiveType` field, which allows users to specify a dedicated instance specification. This capability is already available in the vendored SDK (`tencentcloud-sdk-go/tencentcloud/vpc/v20170312`) on the following API operations:

- **CreateNatGatewayRequest.ExclusiveType** - specified at creation time
- **NatGateway.ExclusiveType** - returned in DescribeNatGateways response
- **ResetNatGatewayConnectionRequest.ExclusiveType** - used for modifying the specification

The current resource code does not expose this field to Terraform users, limiting their ability to manage exclusive NAT gateways via IaC.

## Goals / Non-Goals

**Goals:**
- Add a new optional `exclusive_type` parameter to the `tencentcloud_nat_gateway` resource schema
- Support the parameter in Create (CreateNatGateway), Read (DescribeNatGateways), and Update (ResetNatGatewayConnection) operations
- Maintain backward compatibility - existing configurations continue to work unchanged
- Follow existing code patterns used by the resource

**Non-Goals:**
- Adding other new parameters not related to exclusive type
- Changing the behavior of existing parameters
- Implementing new data sources or resources

## Decisions

1. **Schema Type: TypeString, Optional, Computed, ForceNew**
   - Rationale: `ExclusiveType` is an optional string field in the API. It's set at creation time (`ForceNew` ensures recreation if the exclusive type fundamentally changes). However, the `ResetNatGatewayConnection` API also accepts `ExclusiveType` for in-place updates, so we will support updates via this API to allow changing the exclusive type without recreation. `Computed` allows the system to detect the value from the read response.
   - Note: Looking at the existing resource code, `ResetNatGatewayConnection` is already used for updating `max_concurrent`. We follow the same pattern.

2. **Field naming: `exclusive_type`** (snake_case, following Terraform convention)
   - Rationale: Consistent with existing Terraform parameter naming (e.g., `nat_product_version`).

3. **Update via ResetNatGatewayConnection**
   - The `ResetNatGatewayConnection` API already supports `ExclusiveType`. We will add logic to handle `exclusive_type` changes in the update function using the same API call pattern as `max_concurrent`.

4. **Validation: ValidateStringLengthInRange with allowed values**
   - The exclusive type values are: `ExclusiveSmall`, `ExclusiveMedium1`, `ExclusiveLarge1`
   - We will use a validate function to ensure only valid values are accepted.

## Risks / Trade-offs

- [Risk] The `ResetNatGatewayConnection` API's `ExclusiveType` field may have specific constraints when changing the exclusive type (e.g., can only upgrade, not downgrade). → Mitigation: The API will return an error which will be propagated to the user.
- [Risk] The `ExclusiveType` parameter might not be supported in all regions. → Mitigation: This is a standard API parameter; if the region does not support it, the API will return an appropriate error.
- [Trade-off] Setting `ForceNew` on the field is conservative but may cause unnecessary recreation for users who only want to update the exclusive type in-place. We mitigate this by implementing the update path via `ResetNatGatewayConnection` to allow in-place changes when supported by the API.