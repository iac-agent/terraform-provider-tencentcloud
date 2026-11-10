## 1. Schema Definition

- [x] 1.1 Add `filters` optional parameter to `tencentcloud_teo_rule_engine` resource schema in `resource_tc_teo_rule_engine.go`: TypeList of objects with `name` (TypeString, Required) and `values` (TypeSet of TypeString, Required) sub-fields
- [x] 1.2 Add `rule_items` computed parameter to `tencentcloud_teo_rule_engine` resource schema in `resource_tc_teo_rule_engine.go`: TypeList of objects with `rule_id` (TypeString), `rule_name` (TypeString), `status` (TypeString), `rules` (TypeList, same nested structure as existing `rules`), `rule_priority` (TypeInt), and `tags` (TypeSet of TypeString) sub-fields

## 2. Service Layer

- [x] 2.1 Add `DescribeTeoRuleEngineByFilters` service method in `service_tencentcloud_teo.go` that accepts `zoneId` and custom `filters` parameters, calls `DescribeRules` API with the provided filters, and returns the full response including `ZoneId` and `RuleItems`
- [x] 2.2 Keep the existing `DescribeTeoRuleEngineById` method unchanged to maintain backward compatibility for the default read path

## 3. Read Function

- [x] 3.1 Modify `resourceTencentCloudTeoRuleEngineRead` in `resource_tc_teo_rule_engine.go` to read the `filters` parameter from schema and pass it to the service layer when present
- [x] 3.2 Add logic in the read function to flatten `RuleItems` from the `DescribeRules` response into the `rule_items` schema field
- [x] 3.3 Ensure backward compatibility: when `filters` is not specified, the read function continues to use the existing `DescribeTeoRuleEngineById` method with the hardcoded `rule-id` filter

## 4. Unit Tests

- [x] 4.1 Add unit tests for the `filters` parameter in `resource_tc_teo_rule_engine_test.go` using gomonkey mock to verify that custom filters are correctly passed to the `DescribeRules` API
- [x] 4.2 Add unit tests for the `rule_items` computed parameter to verify that the response is correctly flattened into the schema
- [x] 4.3 Add unit tests for backward compatibility: verify that when `filters` is not specified, the existing behavior is preserved

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_rule_engine.md` to document the new `filters` and `rule_items` parameters with example usage

## 6. Verification

- [x] 6.1 Run `go test` on the modified test files to ensure all unit tests pass
- [x] 6.2 Run `gofmt` to ensure code formatting is correct
