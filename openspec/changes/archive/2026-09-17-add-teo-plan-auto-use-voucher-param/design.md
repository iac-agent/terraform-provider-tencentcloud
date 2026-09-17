## Context

`tencentcloud_teo_plan` is a RESOURCE_KIND_GENERAL resource that wraps the TEO plan (套餐) lifecycle:

- Create → `CreatePlan` (returns `PlanId`)
- Read → `DescribePlans` (via service layer `DescribeTeoPlansById`, filter `plan-id`)
- Update → `UpgradePlan` (plan_type change) / `RenewPlan` (period change) / `ModifyPlan` (renew_flag change)
- Delete → `DestroyPlan`

Current state:

- Resource file: `tencentcloud/services/teo/resource_tc_teo_plan.go`
- Schema fields: `plan_type` (Required), `prepaid_plan_param` (Optional list: `period`, `renew_flag`), plus Computed fields (`plan_id`, `area`, `status`, `pay_mode`, `enabled_time`, `expired_time`)
- No `resource_tc_teo_plan_test.go` and no `resource_tc_teo_plan.md` exist yet
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (vendored)

**API behavior analysis (from vendored SDK):**

| API | AutoUseVoucher in Request | AutoUseVoucher in Response |
|-----|---------------------------|----------------------------|
| `CreatePlan` | Yes — `AutoUseVoucher *string`, values `true`/`false`, default `false`; only effective when PlanType is `personal`/`basic`/`standard` | N/A |
| `DescribePlans` | N/A | No — `Plan` struct has no `AutoUseVoucher` field |
| `UpgradePlan` | Yes — but out of scope for this change (only `CreatePlan` parameter is requested) | N/A |
| `RenewPlan` | Yes — but out of scope for this change | N/A |
| `ModifyPlan` | No | N/A |
| `DestroyPlan` | No | N/A |

Note: `UpgradePlan` and `RenewPlan` also accept `AutoUseVoucher`, but this change is scoped strictly to the `CreatePlan` input parameter (`request.AutoUseVoucher` → `auto_use_voucher`). Extending the parameter to upgrade/renew paths would change existing update behavior and is out of scope.

**Key constraint:** `AutoUseVoucher` is a create-time (and order-related) setting that is not persisted in the `Plan` entity returned by `DescribePlans`, so it cannot be refreshed in Read.

## Goals / Non-Goals

**Goals:**

- Add `auto_use_voucher` (Optional, TypeString, values `true`/`false`) to the `tencentcloud_teo_plan` resource schema, mapped to `request.AutoUseVoucher` of `CreatePlan`
- Pass the parameter to the `CreatePlan` API request only when explicitly set in the configuration (default `false` handled server-side)
- Reject post-creation changes of `auto_use_voucher` with a clear error in Update (immutable args pattern), because the value cannot be modified or refreshed through any plan API used by Read
- Add unit tests (`resource_tc_teo_plan_test.go`) using gomonkey mocks for the resource CRUD functions, covering the new parameter
- Create `resource_tc_teo_plan.md` documentation (one-line description mentioning TEO, Example Usage, Import section) since it does not exist

**Non-Goals:**

- Exposing `AutoUseVoucher` on `UpgradePlan` / `RenewPlan` update flows (out of scope; the task adds one Create parameter only)
- Marking the parameter as ForceNew (create-only params use immutable args errors instead, per provider convention)
- Modifying the datasource `tencentcloud_teo_plans`
- Adding Computed behavior to the parameter (API does not return it)

## Decisions

### Decision 1: Schema field naming — `auto_use_voucher` (snake_case)

**选择**: `auto_use_voucher` (snake_case).

**备选**: `AutoUseVoucher` (PascalCase, matching the requested SchemaName literally).

**理由**:
- Every schema field in the provider (including all teo resources, e.g. `plan_type`, `prepaid_plan_param`) uses snake_case; Terraform HCL convention is snake_case
- The SchemaName mapping `AutoUseVoucher` describes the target parameter identity (the `AutoUseVoucher` field of the resource); converting to provider naming convention keeps the codebase consistent
- The API JSON name remains `AutoUseVoucher` in the request, which is what the SDK marshals

### Decision 2: Type — `TypeString` with values `true`/`false`, no `ValidateFunc`

**选择**: `TypeString`, valid values `true` / `false`, description documents the semantics; no `ValidateFunc`.

**理由**:
- The SDK field is `*string`, and the API documentation defines values as strings `true`/`false`
- Existing provider resources handle similar boolean-ish API strings as TypeString (e.g. `renew_flag` with `on`/`off` uses ValidateFunc; but for plain `true`/`false` strings, letting the API validate keeps parity with many existing fields). `ValidateAllowedStringValue([]string{"true", "false"})` may be added; keeping parity with the teo `plan_type` field style, a `ValidateFunc` with allowed values is acceptable and improves UX. → Final: use `tccommon.ValidateAllowedStringValue([]string{"true", "false"})` for early validation, consistent with `plan_type` in the same resource.

### Decision 3: Create-time only; immutable args check in Update

**选择**: Do not set the field in Read. In Update, add `auto_use_voucher` to an `immutableArgs` array; if `d.HasChange("auto_use_voucher")`, return error `argument 'auto_use_voucher' cannot be changed`.

**备选**: ForceNew: true.

**理由**:
- `DescribePlans` does not return `AutoUseVoucher`, so Read cannot refresh it; if it were mutable, Terraform would see perpetual drift between config and state
- ForceNew would silently destroy and recreate a paid plan — dangerous and unnecessary; a clear immutable error is safer (same pattern as `auto_voucher` in `resource_tc_elasticsearch_logstash.go`, `resource_tc_sqlserver_general_cloud_ro_instance.go`, etc.)
- The parameter only influences order payment at creation time; changing it later has no cloud-side meaning through the plan resource's update APIs

### Decision 4: Pass to request only when explicitly configured

**选择**: In Create, use `if v, ok := d.GetOk("auto_use_voucher"); ok { request.AutoUseVoucher = helper.String(v.(string)) }`.

**理由**:
- Unset → omit from request → API defaults to `false`, preserving current behavior exactly (backward compatibility)
- Consistent with how `plan_type` is currently wired in the same function

### Decision 5: Unit tests with gomonkey mocks (new test file)

**选择**: Create `resource_tc_teo_plan_test.go` with unit tests that mock the SDK client methods (`CreatePlanWithContext`, `DescribePlans`, `UpgradePlanWithContext`, `RenewPlanWithContext`, `ModifyPlanWithContext`, `DestroyPlanWithContext`) via gomonkey, testing:

- Create with `auto_use_voucher = "true"` → request contains `AutoUseVoucher`
- Create without `auto_use_voucher` → request omits `AutoUseVoucher`
- Read populating computed fields
- Update changing `auto_use_voucher` → returns immutable error
- Delete flow

**理由**:
- Per repo rules, for new params/resources tests must use gomonkey mocks against cloud APIs (no terraform test suites, no live API calls)
- The resource currently has no test file; the new parameter's behavior (request wiring + immutability) is testable purely with mocks

### Decision 6: Documentation file creation

**选择**: Create `resource_tc_teo_plan.md` following gendoc/README.md conventions: one-line description ("Provides a resource to create a TEO plan."), Example Usage (with `auto_use_voucher`), Import section (plan id). Do not add `Argument Reference` / `Attribute Reference` (auto-generated).

**理由**:
- Repo rules require every resource to have a `.md` example file that feeds `make doc`
- The file is missing today; adding the parameter is the right moment to add it so the generated website docs include the new field

## Risks / Trade-offs

- **[Risk] Import drift**: after `terraform import`, `auto_use_voucher` is not in state (Read never sets it); a user adding it to config post-import will hit the immutable error instead of an in-place update → **Mitigation**: document in the `.md` that the parameter is create-time only and cannot be changed after creation (including import scenarios)
- **[Risk] `d.GetOk` treats `""` as unset**: acceptable, empty string is not a valid value anyway and the ValidateFunc rejects it at plan time when explicitly set
- **[Trade-off]**: Not exposing `AutoUseVoucher` on `RenewPlan`/`UpgradePlan` means voucher usage cannot be controlled on renew/upgrade via Terraform → out of scope for this single-parameter change; can be a follow-up if requested
- **[Trade-off]**: Boolean-like string (`"true"`/`"false"`) instead of TypeBool → matches SDK `*string` field and avoids conversion ambiguity; documented in the field description

## Migration Plan

- Pure additive schema change (new Optional field); no state migration needed
- Existing configurations: unaffected — the field is omitted from requests when unset
- Rollback: remove the schema entry, the Create wiring, and the immutable args entry; state values (if any) become unmanaged without breaking `terraform plan`

## Open Questions

- None
