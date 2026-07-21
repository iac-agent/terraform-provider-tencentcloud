## 1. Schema Definition

- [x] 1.1 Add `trigger_type` field (TypeString, Optional, Computed) to the resource schema in `resource_tc_teo_function_rule.go`

## 2. Create Handler

- [x] 2.1 Add `trigger_type` parameter wiring in `resourceTencentCloudTeoFunctionRuleCreate`: read from `d.GetOk("trigger_type")` and set `request.TriggerType`

## 3. Read Handler

- [x] 3.1 Add `trigger_type` read-back in `resourceTencentCloudTeoFunctionRuleRead`: set `d.Set("trigger_type", respData.TriggerType)` when response data is not nil

## 4. Update Handler

- [x] 4.1 Add `trigger_type` to the `mutableArgs` array in `resourceTencentCloudTeoFunctionRuleUpdate`
- [x] 4.2 Add `trigger_type` parameter wiring in the update request body: read from `d.GetOk("trigger_type")` and set `request.TriggerType`

## 5. Unit Tests

- [x] 5.1 Add unit test cases for `trigger_type` parameter in `resource_tc_teo_function_rule_test.go` covering: create with trigger_type, read trigger_type, update trigger_type

## 6. Documentation

- [x] 6.1 Update `resource_tc_teo_function_rule.md` with `trigger_type` in the Example Usage section

## 7. Verification

- [x] 7.1 Run `gofmt` on modified Go files
- [x] 7.2 Run `make doc` to generate website documentation
- [x] 7.3 Verify the code compiles successfully