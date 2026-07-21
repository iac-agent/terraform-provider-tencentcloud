## Context

The `tencentcloud_teo_function_rule` resource manages TEO edge function trigger rules. The TEO cloud API (`CreateFunctionRule`, `ModifyFunctionRule`, `DescribeFunctionRules`) supports a `TriggerType` parameter that controls the function selection strategy. Currently, the Terraform resource does not expose this parameter, limiting users to the default `direct` behavior.

The SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) already includes `TriggerType` in:
- `CreateFunctionRuleRequest` (line 3421)
- `ModifyFunctionRuleRequest` (line 19778)
- `FunctionRule` response struct (line 16656)

No new SDK dependency or version bump is needed.

## Goals / Non-Goals

**Goals:**
- Add `trigger_type` (TypeString, Optional, Computed) to the `tencentcloud_teo_function_rule` resource schema
- Wire the parameter into Create, Read, and Update operations
- Keep backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Add `RegionMappingSelections` or `WeightedSelections` parameters (these are related but separate parameters for region/weight trigger types)
- Change the resource ID format or Import behavior
- Modify the Delete operation

## Decisions

1. **Parameter type: `TypeString` with `Optional: true, Computed: true`**
   - The field is optional in the cloud API (defaults to `direct` when not specified)
   - `Computed: true` ensures the value is properly read back from the API response during Read operations, even when not explicitly set in the configuration

2. **Wire into Create, Read, and Update only**
   - `CreateFunctionRule` accepts `TriggerType` — set it in the Create handler
   - `ModifyFunctionRule` accepts `TriggerType` — set it in the Update handler (alongside existing `function_rule_conditions` and `remark`)
   - `FunctionRule` response includes `TriggerType` — read it back in the Read handler
   - `DeleteFunctionRules` does not need `TriggerType` — no changes needed

3. **Add to `mutableArgs` in the Update function**
   - The `TriggerType` is mutable on the cloud API side, so it should be added to the `mutableArgs` array in `resourceTencentCloudTeoFunctionRuleUpdate`

## Risks / Trade-offs

- [Risk] If a user sets `trigger_type` to `weight` or `region` without also configuring the corresponding `RegionMappingSelections` or `WeightedSelections`, the API call may fail → Mitigation: This is standard Terraform behavior — the API error will be surfaced to the user. The linked parameters are separate future additions.
- [Trade-off] Not adding `RegionMappingSelections` and `WeightedSelections` in this change → These are independent parameters that can be added in follow-up changes without blocking the `trigger_type` addition.