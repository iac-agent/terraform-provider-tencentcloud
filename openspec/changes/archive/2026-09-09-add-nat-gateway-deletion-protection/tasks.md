## 1. Verify Create Function Implementation

- [x] 1.1 Verify that `resourceTencentCloudNatGatewayCreate` in `tencentcloud/services/vpc/resource_tc_nat_gateway.go` correctly reads `deletion_protection_enabled` from schema and passes it to `CreateNatGatewayRequest.DeletionProtectionEnabled`
- [x] 1.2 Verify that `resourceTencentCloudNatGatewayRead` correctly reads `DeletionProtectionEnabled` from `DescribeNatGateways` response and sets it on the resource data
- [x] 1.3 Verify that `resourceTencentCloudNatGatewayUpdate` correctly handles `deletion_protection_enabled` changes via `ModifyNatGatewayAttribute`

## 2. Unit Test Coverage

- [x] 2.1 Add unit test in `tencentcloud/services/vpc/resource_tc_nat_gateway_test.go` using gomonkey mock to verify `deletion_protection_enabled` is correctly passed to `CreateNatGateway` API during creation
- [x] 2.2 Add unit test to verify `deletion_protection_enabled` is correctly read from `DescribeNatGateways` response
- [x] 2.3 Add unit test to verify `deletion_protection_enabled` update via `ModifyNatGatewayAttribute`

## 3. Documentation Update

- [x] 3.1 Update `tencentcloud/services/vpc/resource_tc_nat_gateway.md` to include `deletion_protection_enabled` in the example usage
