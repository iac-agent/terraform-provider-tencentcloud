## ADDED Requirements

### Requirement: exclusive_type parameter on tencentcloud_nat_gateway

The `tencentcloud_nat_gateway` resource SHALL support an optional `exclusive_type` parameter to specify the exclusive (dedicated) NAT gateway instance specification.

- The parameter SHALL be of type `string`
- The parameter SHALL be optional
- The parameter SHALL accept values: `ExclusiveSmall`, `ExclusiveMedium1`, `ExclusiveLarge1`
- The parameter SHALL be `Computed` to allow the system to read back the value
- The parameter SHALL be `ForceNew` for creation but also support in-place update via `ResetNatGatewayConnection`
- Validation SHALL only accept the three allowed values
- When not set, the NAT gateway SHALL be created as a standard (non-exclusive) gateway

#### Scenario: Create NAT gateway with exclusive_type
- **WHEN** a user creates `tencentcloud_nat_gateway` with `exclusive_type = "ExclusiveSmall"`
- **THEN** the CreateNatGateway API SHALL be called with `ExclusiveType = "ExclusiveSmall"`
- **THEN** the created NAT gateway SHALL be an exclusive (dedicated) type gateway

#### Scenario: Read exclusive_type from API response
- **WHEN** reading a NAT gateway that has `ExclusiveType` set
- **THEN** the `exclusive_type` field SHALL be populated with the value from the API response

#### Scenario: Update exclusive_type via ResetNatGatewayConnection
- **WHEN** a user updates `exclusive_type` on an existing NAT gateway
- **THEN** the ResetNatGatewayConnection API SHALL be called with the new `ExclusiveType` value

#### Scenario: Backward compatibility - no exclusive_type set
- **WHEN** an existing configuration does not set `exclusive_type`
- **THEN** the NAT gateway SHALL be created as a standard (non-exclusive) gateway
- **THEN** no `ExclusiveType` field SHALL be sent in the API request