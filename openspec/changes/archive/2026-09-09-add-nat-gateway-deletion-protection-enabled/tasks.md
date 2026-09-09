## 1. Schema Definition

- [x] 1.1 Add `deletion_protection_enabled` field to the resource schema in `resource_tc_nat_gateway.go` as `Optional` + `Computed` bool type

## 2. Create Implementation

- [x] 2.1 Add `deletion_protection_enabled` to `CreateNatGatewayRequest` in the Create function, using `d.GetOkExists` to check if the value is set

## 3. Read Implementation

- [x] 3.1 Read `deletion_protection_enabled` from `NatGateway.DeletionProtectionEnabled` response field in the Read function, with nil check before setting to state

## 4. Update Implementation

- [x] 4.1 Add `deletion_protection_enabled` change detection using `d.HasChange` in the Update function, passing value to `ModifyNatGatewayAttributeRequest.DeletionProtectionEnabled`

## 5. Documentation

- [x] 5.1 Update `resource_tc_nat_gateway.md` to include the new `deletion_protection_enabled` parameter with example usage

## 6. Code Quality

- [x] 6.1 Run `gofmt` to format the modified code
- [x] 6.2 Verify the code compiles without errors