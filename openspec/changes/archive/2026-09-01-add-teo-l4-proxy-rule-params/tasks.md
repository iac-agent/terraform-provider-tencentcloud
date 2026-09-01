## 1. Schema Changes

- [x] 1.1 Add `bu_id` field (Optional, Computed, TypeString) to the `l4_proxy_rules` nested schema in `resource_tc_teo_l4_proxy_rule.go`
- [x] 1.2 Add `remote_auth` nested block (Computed, TypeList, MaxItems: 1) to the `l4_proxy_rules` nested schema with sub-fields: `switch` (TypeString), `address` (TypeString), `server_faulty_behavior` (TypeString)

## 2. Create Function Changes

- [x] 2.1 Add `BuId` to the `L4ProxyRule` struct construction in `resourceTencentCloudTeoL4ProxyRuleCreate` from `l4_proxy_rules` block

## 3. Read Function Changes

- [x] 3.1 Add `BuId` reading from `l4ProxyRule.BuId` and setting to `l4ProxyRuleMap["bu_id"]` in `resourceTencentCloudTeoL4ProxyRuleRead`
- [x] 3.2 Add `RemoteAuth` reading from `l4ProxyRule.RemoteAuth` and setting to `l4ProxyRuleMap["remote_auth"]` with sub-fields `switch`, `address`, `server_faulty_behavior` in `resourceTencentCloudTeoL4ProxyRuleRead`

## 4. Update Function Changes

- [x] 4.1 Add `BuId` to the `L4ProxyRule` struct construction in `resourceTencentCloudTeoL4ProxyRuleUpdate` from `l4_proxy_rules` block

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l4_proxy_rule.md` with `bu_id` and `remote_auth` example usage

## 6. Unit Tests

- [x] 6.1 Update `resource_tc_teo_l4_proxy_rule_test.go` with test cases covering `bu_id` and `remote_auth` fields