## Why

The existing `tencentcloud_teo_plan` resource is missing the `auto_use_voucher` parameter from the CreatePlan API and the `deal_name` return value. Additionally, the `renew_flag` field in the ModifyPlan API uses a `RenewFlag` struct (containing a `Switch` field) rather than a simple string, which requires proper modeling. A new v2 resource is needed to provide a more complete and accurate mapping of the TE0 Plan cloud API.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_plan_v2` (RESOURCE_KIND_GENERAL) for managing TE0 Plan lifecycle
- Support the following cloud API interfaces:
  - **CreatePlan**: Create a plan with `plan_type`, `auto_use_voucher` (new), and `prepaid_plan_param` parameters; return `plan_id` and `deal_name` (new computed field)
  - **DescribePlans**: Read plan information using `filters`, `order`, `direction`; return full plan details
  - **ModifyPlan**: Update plan configuration with `plan_id` and `renew_flag` (using `RenewFlag` struct with `Switch` field)
  - **DestroyPlan**: Delete plan with `plan_id`
- The v2 resource adds the missing `auto_use_voucher` field (whether to auto-use vouchers during plan creation)
- The v2 resource adds the `deal_name` computed field (order number returned from CreatePlan)
- The `renew_flag` field in the update flow properly models the SDK's `RenewFlag` struct

## Capabilities

### New Capabilities
- `teo-plan-v2-resource`: New TE0 Plan v2 resource with full CRUD support, including auto_use_voucher and deal_name fields

### Modified Capabilities

## Impact

- New resource file: `tencentcloud/services/teo/resource_tc_teo_plan_v2.go`
- New test file: `tencentcloud/services/teo/resource_tc_teo_plan_v2_test.go`
- New documentation file: `tencentcloud/services/teo/resource_tc_teo_plan_v2.md`
- Registration in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Reuses existing `TeoService.DescribeTeoPlansById` service method for Read operation
- Cloud APIs used: CreatePlan, DescribePlans, ModifyPlan, DestroyPlan (from teo v20220901 SDK)
