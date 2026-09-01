## Context

The `tencentcloud_teo_l4_proxy_rule` resource wraps the TEO L4 proxy forwarding rule APIs: `CreateL4ProxyRules`, `DescribeL4ProxyRules`, `ModifyL4ProxyRules`, and `DeleteL4ProxyRules`. The current resource schema uses a nested `l4_proxy_rules` block (TypeList, MaxItems: 1) containing the forwarding rule fields.

The vendor SDK (`tencentcloud-sdk-go/tencentcloud/teo/v20220901`) already defines `BuId` and `L4ProxyRemoteAuth` (with `Switch`, `Address`, `ServerFaultyBehavior`) in the `L4ProxyRule` struct. These fields are returned by `DescribeL4ProxyRules` and `BuId` is accepted as optional input by `CreateL4ProxyRules` and `ModifyL4ProxyRules`.

Per the vendor SDK documentation:
- `RemoteAuth` is explicitly documented as ignored in Create/Modify: "RemoteAuth 在 CreateL4ProxyRules 或 ModifyL4ProxyRules 不可作为入参使用，如有传此参数，会忽略。在 DescribeL4ProxyRules 返回为空时，表示没有开启远程鉴权。"
- `BuId` has no such restriction and can be used as both input and output.

## Goals / Non-Goals

**Goals:**
- Add `BuId` as an Optional+Computed field to the `l4_proxy_rules` nested schema, readable from Describe and writable in Create/Modify
- Add `RemoteAuth` as a Computed-only nested block (TypeList, MaxItems: 1) with `Switch`, `Address`, `ServerFaultyBehavior` sub-fields, readable from Describe only
- Set `BuId` and `RemoteAuth` sub-fields in the Read function from the Describe response
- Include `BuId` in the Create and Update request building logic

**Non-Goals:**
- Expose `Offset`, `Limit`, `Filters` (Name/Values) — these are internal pagination/filtering parameters for the Describe API, used by the service layer, not exposed to Terraform users
- Expose `TotalCount` — this is a response-level field from `DescribeL4ProxyRulesResponse`, not part of the `L4ProxyRule` struct
- Make `RemoteAuth` writable — the API ignores it in Create/Modify
- Make `Status` writable — the API explicitly forbids filling it in Create/Modify; it remains Computed

## Decisions

1. **BuId: Optional+Computed** — The field is optional in API input and returned in Describe. Making it Computed ensures backward compatibility for existing state that doesn't have it.

2. **RemoteAuth: Computed-only** — Since the API explicitly ignores RemoteAuth in Create/Modify, making it writable would mislead users. Setting it as Computed-only ensures users can read remote auth configuration without expecting to set it through Terraform.

3. **RemoteAuth as TypeList nested block** — Following the existing pattern of `l4_proxy_rules` as TypeList with MaxItems: 1, the RemoteAuth block uses the same structure for consistency.

4. **No changes to service layer** — The `DescribeTeoL4ProxyRuleById` function already returns the full `L4ProxyRule` struct including `BuId` and `RemoteAuth`. No modifications needed.

5. **No changes to extension file** — The `_extension.go` file handles async state waiting and delete pre-processing. These operations are not affected by the new read-only fields.

## Risks / Trade-offs

- **Risk: RemoteAuth could become writable in a future API version** → Mitigation: Can be promoted to Optional+Computed later without breaking existing configurations
- **Risk: BuId may have validation constraints not documented in the SDK** → Mitigation: The field is a plain string in the SDK; any server-side validation errors will be returned by the API and surfaced to users