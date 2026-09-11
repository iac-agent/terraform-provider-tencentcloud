## Context

The `tencentcloud_eni_ipv6_address` resource manages IPv6 addresses assigned to a network interface. The existing schema already includes `ipv6_addresses` as a TypeSet with sub-fields (`address`, `description`, `primary`, `address_id`, `is_wan_ip_blocked`, `state`), and `ipv6_address_count`.

The vendored SDK struct `Ipv6Address` already defines a `PublicIpAddress *string` field (ULA type public IP), but this field is not exposed in the Terraform resource schema or handled in the Create/Read flows. The corresponding cloud API `AssignIpv6Addresses` passes `Ipv6Address` objects which include this field, and `DescribeNetworkInterfaces` returns it in the response.

## Goals / Non-Goals

**Goals:**
- Add `public_ip_address` as an Optional+Computed field in the `ipv6_addresses` sub-block of the resource schema
- Pass `PublicIpAddress` to `AssignIpv6Addresses` request when set by the user
- Read and store `PublicIpAddress` from `DescribeNetworkInterfaces` response during Read
- Ensure backward compatibility — existing configurations and state files remain valid

**Non-Goals:**
- No new top-level resource parameters (only within the existing `ipv6_addresses` sub-block)
- No changes to `ipv6_address_count` behavior
- No new CRUD operations or API calls
- No state migration needed

## Decisions

1. **Field type: `TypeString`, Optional+Computed**
   - Rationale: The field is a string type (IP address). Optional allows users to set it; Computed allows the API to populate it on Read even if not specified in the config. This matches patterns used by other optional fields in the same block (e.g., `description`, `address_id`).

2. **Place existing `PublicIpAddress` field after `address_id` in schema ordering**
   - Rationale: The SDK struct lists `PublicIpAddress` after `State` and before `AddressType`. In the existing schema, fields follow a logical order: `address`, `description`, `primary`, `address_id`, `is_wan_ip_blocked`, `state`. Adding `public_ip_address` after `state` keeps related fields together.

3. **No ForceNew required**
   - Rationale: `PublicIpAddress` is informational metadata for ULA-type IPv6 addresses, not an identifier. Adding ForceNew would cause unnecessary resource recreation. However, looking at the existing pattern — all sub-fields in `ipv6_addresses` use `ForceNew: true` — we follow the same convention for consistency.

4. **Use `GetOk` + type assertion pattern for reading the field**
   - Rationale: This matches the existing code pattern used for all other fields in the `ipv6_addresses` loop.

## Risks / Trade-offs

- **[Low Risk] Field is nil-safe**: The SDK uses `omitempty` so an unset `PublicIpAddress` is simply omitted from the API request. The Read path checks for nil before setting via the current `map[string]interface{}` approach.
- **[Trade-off] ForceNew adds unnecessary churn**: The field is purely informational for ULA addresses — changing it doesn't require recreating the IPv6 address. But to stay consistent with the existing pattern (all sub-fields are ForceNew), we follow the same convention.
- **[No Risk] Backward compatibility**: Adding an Optional+Computed field to an existing TypeSet sub-block does not break existing configurations or state.