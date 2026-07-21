## Why

The existing `tencentcloud_teo_function_rule` resource does not expose the `TriggerType` parameter, which is available in the TEO cloud API (`CreateFunctionRule`, `ModifyFunctionRule`, `DescribeFunctionRules`). This parameter controls the function selection strategy (direct/weight/region) and is essential for users who need to configure weight-based or region-based function routing.

## What Changes

- Add a new optional `trigger_type` parameter (TypeString, Optional, Computed) to the `tencentcloud_teo_function_rule` resource schema
- Wire the parameter into the Create (CreateFunctionRule), Read (DescribeFunctionRules), and Update (ModifyFunctionRule) operations
- No changes needed for Delete (DeleteFunctionRules) as the parameter is not applicable

## Capabilities

### New Capabilities
- `teo-function-rule-trigger-type`: Add `trigger_type` parameter to the `tencentcloud_teo_function_rule` resource to support function selection strategy configuration (direct, weight, region)

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- Affected code: `tencentcloud/services/teo/resource_tc_teo_function_rule.go`
- Affected docs: `tencentcloud/services/teo/resource_tc_teo_function_rule.md`
- SDK dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already contains `TriggerType` field in `CreateFunctionRuleRequest`, `ModifyFunctionRuleRequest`, and `FunctionRule` response struct)
- Backward compatible: Yes, the new field is Optional and Computed, existing configurations remain valid