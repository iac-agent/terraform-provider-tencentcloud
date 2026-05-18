## 1. Schema and CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_v2.go` with resource schema definition including plan_type (Required, ForceNew), auto_use_voucher (Optional, ForceNew), prepaid_plan_param (Optional, ForceNew, MaxItems:1 with period and renew_flag), renew_flag (Optional, top-level), and computed fields (plan_id, deal_name, area, status, pay_mode, enabled_time, expired_time)
- [x] 1.2 Implement `ResourceTencentCloudTeoPlanV2Create` function: build CreatePlanRequest from schema, call CreatePlanWithContext with retry, validate response (check PlanId and DealName are not nil), set d.SetId(planId)
- [x] 1.3 Implement `ResourceTencentCloudTeoPlanV2Read` function: call TeoService.DescribeTeoPlansById, handle nil response (d.SetId("")), set all computed and output fields with nil checks
- [x] 1.4 Implement `ResourceTencentCloudTeoPlanV2Update` function: check d.HasChange("renew_flag"), build ModifyPlanRequest with RenewFlag struct (Switch field), call ModifyPlanWithContext with retry
- [x] 1.5 Implement `ResourceTencentCloudTeoPlanV2Delete` function: build DestroyPlanRequest with PlanId, call DestroyPlanWithContext with retry
- [x] 1.6 Add Importer support with schema.ImportStatePassthrough

## 2. Provider Registration

- [x] 2.1 Add `"tencentcloud_teo_plan_v2": teo.ResourceTencentCloudTeoPlanV2()` to ResourcesMap in `tencentcloud/provider.go`
- [x] 2.2 Add `tencentcloud_teo_plan_v2` entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_v2.md` with one-line description (including "TEO" product name), Example Usage section with jsonencode() for JSON strings, and Import section

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_v2_test.go` with gomonkey mock-based unit tests for Create, Read, Update, Delete operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` and ensure all tests pass

## 5. Code Correctness Verification

- [x] 5.1 Verify all CreatePlan input parameters exist in the CreatePlan SDK request struct
- [x] 5.2 Verify all ModifyPlan input parameters exist in the ModifyPlan SDK request struct
- [x] 5.3 Verify DestroyPlan request struct accepts PlanId parameter
- [x] 5.4 Verify DescribePlans response Plan struct contains all computed fields used in Read
