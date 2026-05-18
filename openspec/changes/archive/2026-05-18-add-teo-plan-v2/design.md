## Context

The existing `tencentcloud_teo_plan` resource manages TE0 (EdgeOne) plans but has incomplete API parameter coverage. It is missing the `auto_use_voucher` parameter from CreatePlan and the `deal_name` return value. The v2 resource provides a more complete and accurate mapping of the TE0 Plan cloud API while following the standard RESOURCE_KIND_GENERAL pattern.

Current state:
- `tencentcloud_teo_plan` exists with basic CRUD (CreatePlan, DescribePlans, UpgradePlan/RenewPlan/ModifyPlan for update, DestroyPlan for delete)
- Missing: `auto_use_voucher` parameter, `deal_name` computed field
- The existing resource handles `renew_flag` as a simple string in `prepaid_plan_param`, but ModifyPlan API uses a `RenewFlag` struct with a `Switch` field

## Goals / Non-Goals

**Goals:**
- Add a new `tencentcloud_teo_plan_v2` resource with complete API parameter coverage
- Include `auto_use_voucher` as an optional parameter for CreatePlan
- Include `deal_name` as a computed field from CreatePlan response
- Properly model `renew_flag` update via ModifyPlan API with the `RenewFlag` struct
- Support resource import via `plan_id`
- Follow the igtm_strategy resource pattern for code style
- Provide unit tests using gomonkey mock approach

**Non-Goals:**
- Modify the existing `tencentcloud_teo_plan` resource (backward compatibility)
- Support UpgradePlan or RenewPlan operations in v2 (only ModifyPlan for renew_flag update)
- Add additional computed fields beyond what the DescribePlans API returns

## Decisions

### Decision 1: New v2 resource instead of modifying existing resource
- **Choice**: Create `tencentcloud_teo_plan_v2` as a separate resource
- **Rationale**: Adding `auto_use_voucher` (which affects creation behavior) and changing the update logic model constitutes a meaningful schema change. Creating a v2 resource preserves backward compatibility for existing `tencentcloud_teo_plan` users.

### Decision 2: Update scope limited to ModifyPlan
- **Choice**: The v2 resource's update method only supports ModifyPlan (renew_flag changes). `plan_type` and `prepaid_plan_param.period` are ForceNew.
- **Rationale**: The user's API mapping only specifies ModifyPlan for update. UpgradePlan and RenewPlan are not included in the v2 specification. Fields that cannot be updated through ModifyPlan should be ForceNew.

### Decision 3: renew_flag modeled as top-level optional field
- **Choice**: Model `renew_flag` as a top-level optional string field (values: "on"/"off") in the schema, mapped to the `RenewFlag` struct's `Switch` field in the ModifyPlan request.
- **Rationale**: ModifyPlan API uses `RenewFlag{Switch: string}` struct. While CreatePlan's `PrepaidPlanParam` also has a `renew_flag` string field, these are semantically the same (auto-renewal on/off). The top-level field simplifies the update path and aligns with the ModifyPlan API's structure.

### Decision 4: Reuse existing DescribeTeoPlansById service method
- **Choice**: Reuse the existing `TeoService.DescribeTeoPlansById` method for Read operation
- **Rationale**: The existing service method already handles DescribePlans with plan-id filter, retry logic, and response parsing. No need to duplicate this logic.

### Decision 5: Computed fields from Plan struct
- **Choice**: Include the following computed fields from DescribePlans response: `area`, `status`, `pay_mode`, `enabled_time`, `expired_time`
- **Rationale**: These fields are returned by the Plan struct in DescribePlans and provide useful information about the plan state.

## Risks / Trade-offs

- [Risk] Existing `tencentcloud_teo_plan` users may be confused by the v2 resource → Mitigation: Clear documentation explaining the differences
- [Risk] `plan_type` and `prepaid_plan_param` changes will force resource recreation → Mitigation: This is expected behavior since these fields can only be changed through UpgradePlan/RenewPlan APIs which are not in scope
- [Risk] DestroyPlan is a destructive operation that cannot be undone → Mitigation: This is standard Terraform behavior; users must explicitly run `terraform destroy`
