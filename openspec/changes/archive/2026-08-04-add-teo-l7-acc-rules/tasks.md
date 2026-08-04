## 1. Resource Schema and CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go` with resource schema definition (zone_id, rules, filters, rule_ids), reusing existing `TencentTeoL7RuleBranchBasicInfo` helper for the `branches` sub-schema
- [x] 1.2 Implement `ResourceTencentCloudTeoL7AccRulesCreate` function: call `CreateL7AccRules` API with `zone_id` and `rules`, validate response is not empty, set ID to `zone_id`, then call Read
- [x] 1.3 Implement `ResourceTencentCloudTeoL7AccRulesRead` function: call `DescribeL7AccRules` API with `zone_id` and optional `filters`, paginate with Limit=1000, map response to state; if response is empty, log and return NonRetryableError
- [x] 1.4 Implement `ResourceTencentCloudTeoL7AccRulesUpdate` function: compute diff between old and new `rules` lists (match by `rule_name`), call `ModifyL7AccRule` for changed rules, `CreateL7AccRules` for new rules, `DeleteL7AccRules` for removed rules
- [x] 1.5 Implement `ResourceTencentCloudTeoL7AccRulesDelete` function: call `DeleteL7AccRules` API with all `rule_ids` from state, then call `d.SetId("")`
- [x] 1.6 Implement `Import` support using `zone_id` as the import identifier

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_l7_acc_rules` in `tencentcloud/provider.go` with factory function `ResourceTencentCloudTeoL7AccRules`
- [x] 2.2 Add resource documentation entry in `tencentcloud/provider.md`

## 3. Unit Tests

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go` with unit tests using gomonkey to mock cloud API calls for Create, Read, Update, Delete, and Import scenarios

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md` with example usage and import instructions

## 5. Finalization

- [x] 5.1 Run `gofmt` on all changed Go files (handled by tfpacer-finalize skill)
- [x] 5.2 Run `make doc` to generate website documentation (handled by tfpacer-finalize skill)
- [x] 5.3 Create `.changelog` entry for the new resource (handled by tfpacer-finalize skill)