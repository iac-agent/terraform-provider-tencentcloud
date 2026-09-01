## Why

The `tencentcloud_teo_l4_proxy_rule` resource is missing support for `BuId` (Business Unit ID) and `RemoteAuth` (remote authentication) fields in the `L4ProxyRules` block. These fields are returned by the `DescribeL4ProxyRules` API and are needed for users to read the complete state of L4 proxy forwarding rules. Additionally, `BuId` is accepted by the `CreateL4ProxyRules` and `ModifyL4ProxyRules` APIs as an optional input, enabling users to associate rules with business units.

## What Changes

- Add `BuId` (Optional, Computed) field to the `l4_proxy_rules` sub-schema — writable in Create/Modify and readable from Describe
- Add `RemoteAuth` block (Computed only) with `Switch`, `Address`, `ServerFaultyBehavior` sub-fields to the `l4_proxy_rules` sub-schema — readable from Describe response only (API ignores RemoteAuth in Create/Modify)

**Note on excluded parameters:** The following parameters from the API specification were evaluated and excluded due to vendor API constraints:
- `Status` — already present in the schema as Computed; API explicitly forbids filling it in Create/Modify
- `RuleId` — already handled internally in the code (set from the ID split in Modify, not filled in Create per API spec)
- `Offset`, `Limit`, `Filters` (Name/Values) — these are internal Describe API parameters used by the service layer, not exposed to Terraform users
- `TotalCount` — this is a Describe API response-level field, not part of the `L4ProxyRule` struct, and not applicable to a single-resource read
- `RuleIds` — already used internally in the Delete flow
- `L4ProxyRuleIds` — already used internally from the Create response

## Capabilities

### New Capabilities
- `teo-l4-proxy-rule-params`: Add `BuId` (Optional) and `RemoteAuth` (Computed) fields to the `l4_proxy_rules` sub-schema of `tencentcloud_teo_l4_proxy_rule`

### Modified Capabilities
<!-- No existing capability specs are modified -->

## Impact

- **Affected code**: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule.go` (schema and CRUD methods), `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule.md` (documentation)
- **Affected dependencies**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (SDK models already contain `BuId` and `L4ProxyRemoteAuth` fields)
- **Backward compatibility**: Fully backward compatible — all new fields are Optional or Computed; existing configurations and state are unaffected