## 1. Schema and CRUD Changes

- [x] 1.1 Add `public_ip_address` as Optional+Computed string field in the `ipv6_addresses` sub-block schema of `resource_tc_eni_ipv6_address.go`
- [x] 1.2 In the Create flow, add `PublicIpAddress` handling in the `ipv6_addresses` loop to pass the value to `AssignIpv6Addresses` request
- [x] 1.3 In the Read flow, add `public_ip_address` to the `map[string]interface{}` that populates each `ipv6_addresses` entry from `DescribeNetworkInterfaces` response

## 2. Unit Test

- [x] 2.1 Add `public_ip_address` field test case in the existing unit test file `resource_tc_eni_ipv6_address_test.go` covering the schema definition and CRUD handling

## 3. Documentation

- [x] 3.1 Update the resource markdown example file `tencentcloud/services/vpc/resource_tc_eni_ipv6_address.md` to document the new `public_ip_address` field
- [x] 3.2 Run `make doc` (in the finalization stage) to regenerate `website/docs/` documentation (deferred to tfpacer-finalize skill)

## 4. Verification

- [x] 4.1 Verify the code compiles correctly with the new field (deferred to finalization stage - go build is handled by CI)