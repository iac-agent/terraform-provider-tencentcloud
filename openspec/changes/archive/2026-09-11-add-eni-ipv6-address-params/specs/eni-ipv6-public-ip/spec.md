## ADDED Requirements

### Requirement: The resource SHALL expose `public_ip_address` in the `ipv6_addresses` block

The `tencentcloud_eni_ipv6_address` resource SHALL support `public_ip_address` as an Optional+Computed string field within the `ipv6_addresses` sub-block. This field represents the public IP address associated with a ULA-type IPv6 address.

#### Scenario: User specifies public_ip_address in ipv6_addresses
- **WHEN** a user configures `public_ip_address` in the `ipv6_addresses` block of `tencentcloud_eni_ipv6_address`
- **THEN** the resource SHALL pass the value to `AssignIpv6Addresses` API as `Ipv6Addresses[].PublicIpAddress`
- **AND** the resource SHALL read and store the returned `PublicIpAddress` in the Terraform state

#### Scenario: public_ip_address is read back from the API
- **WHEN** the resource reads IPv6 addresses via `DescribeNetworkInterfaces`
- **THEN** the `public_ip_address` field SHALL be populated from `Response.NetworkInterfaceSet.Ipv6AddressSet[].PublicIpAddress` in the Terraform state

#### Scenario: public_ip_address is not set by user
- **WHEN** a user does NOT set `public_ip_address` in their configuration
- **THEN** the field SHALL be computed from the API response if available
- **AND** the field SHALL NOT be sent in the `AssignIpv6Addresses` request