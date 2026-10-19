## Why

The TencentCloud ENI IPv6 address resource (`tencentcloud_eni_ipv6_address`) currently supports managing IPv6 addresses on a network interface, but the `Ipv6Address` struct in the SDK already includes a `PublicIpAddress` field that is not exposed to Terraform users. Adding this field enables users to specify and read the public IP address associated with ULA (Unique Local Address) type IPv6 addresses, completing the parity between the Terraform resource and the underlying cloud API.

## What Changes

- Add `public_ip_address` as a new Optional+Computed field in the `ipv6_addresses` sub-block of `tencentcloud_eni_ipv6_address`
- In the `AssignIpv6Addresses` (Create) flow, pass the `PublicIpAddress` value to the API request when provided
- In the `DescribeNetworkInterfaces` (Read) flow, read and store the `PublicIpAddress` response value
- The change is fully backward-compatible since this is an Optional field with no existing dependency

## Capabilities

### New Capabilities

- `eni-ipv6-public-ip`: Expose the `PublicIpAddress` field of `Ipv6Address` in the `tencentcloud_eni_ipv6_address` resource, allowing users to specify and read the public IPv6 address bound to a ULA-type IPv6 address.

### Modified Capabilities

- *(none)*

## Impact

- **File**: `tencentcloud/services/vpc/resource_tc_eni_ipv6_address.go` — add `public_ip_address` to the `ipv6_addresses` sub-block schema, and handle it in Create and Read flows
- **SDK**: Already supported — the vendored SDK struct `Ipv6Address` already contains `PublicIpAddress *string` (field number 26656)
- **API**: `AssignIpv6Addresses` already passes `Ipv6Address` objects which include `PublicIpAddress`; `DescribeNetworkInterfaces` already returns it in the response
- **Docs**: Update the resource markdown example file to reflect the new field
- **Backward Compatibility**: Fully compatible — new Optional field, no state migration needed