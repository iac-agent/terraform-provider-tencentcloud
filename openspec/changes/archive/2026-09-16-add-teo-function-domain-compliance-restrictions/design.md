## Context

The `tencentcloud_teo_function` resource manages TEO edge functions. Its Read method calls the service-layer `DescribeTeoFunctionById(ctx, zoneId, functionId)` (in `tencentcloud/services/teo/service_tencentcloud_teo.go`), which wraps the `DescribeFunctions` API of the vendored SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` and returns the first `teo.Function` entry of `Response.Functions`.

The SDK struct `Function` contains a field that is not yet surfaced in the Terraform schema:

- `DomainComplianceRestrictions []*ComplianceRestriction` — 边缘函数默认域名因合规问题产生的地区访问限制列表 (the list of regional access restrictions applied to the function default domain for compliance reasons).

Each `ComplianceRestriction` element contains exactly two fields:
- `Reason *string` — 下发访问限制的原因 (enum: `ICP_RECORD_REQUIRED` 未备案; `GOVERNMENT_ORDER` 政府指令)
- `Region *string` — 限制访问地区的具体国家/地区码，使用 ISO 3166 国家/地区代码标准 (e.g. `CN`)

Current state of the resource schema (`tencentcloud/services/teo/resource_tc_teo_function.go`):
- `zone_id` (Required, ForceNew), `name` (Required), `remark` (Optional), `content` (Required)
- Computed fields: `function_id`, `domain`, `create_time`, `update_time`
- No `domain_compliance_restrictions` block exists today.

The task requires mapping JsonPath `response.Functions.DomainComplianceRestrictions.Reason` → schema field `reason` and `response.Functions.DomainComplianceRestrictions.Region` → schema field `region` (both snake_case per Terraform convention).

## Goals / Non-Goals

**Goals:**
- Expose `DomainComplianceRestrictions` as a computed `domain_compliance_restrictions` list block on `tencentcloud_teo_function`, with each element containing the flattened `reason` and `region` string attributes.
- Populate the block in `resourceTencentCloudTeoFunctionRead` from the `DescribeFunctions` response returned by `DescribeTeoFunctionById`, following the same nil-guard flattening pattern as the `references` block in `resource_tc_teo_origin_group.go`.
- Maintain full backward compatibility (only a new Computed block is added; no existing schema field is altered or removed).
- Cover the new fields with gomonkey mock unit tests (no Terraform acceptance test suite for new cases, per repo convention for modified resources using mock tests where they already exist — the existing test file already contains plain Go tests such as `TestParseTeoFunctionOriginalName`).

**Non-Goals:**
- No changes to Create/Update/Delete operations (`DomainComplianceRestrictions` is a read-only API output; `CreateFunction`/`ModifyFunction`/`DeleteFunction` do not accept these parameters).
- No new top-level scalar parameters `reason`/`region` (the API returns a list of restriction entries; flattening the list into top-level scalars would lose information when multiple restrictions exist).
- No changes to the service layer `DescribeTeoFunctionById` (it already returns the full `*teo.Function` struct including `DomainComplianceRestrictions`).
- No changes to provider registration (resource already registered in `tencentcloud/provider.go`).

## Decisions

### Decision 1: Expose as a computed list block `domain_compliance_restrictions` with flattened sub-fields

**选择**: Add a top-level `Computed` `TypeList` block named `domain_compliance_restrictions`, with `Elem: &schema.Resource{...}` containing `reason` (Computed, TypeString) and `region` (Computed, TypeString). In the Read method, iterate over `respData.DomainComplianceRestrictions` and append a `map[string]interface{}` per entry (nil-guarding each field), then `d.Set("domain_compliance_restrictions", list)` — identical to the `references` block pattern in `resource_tc_teo_origin_group.go`.

**备选**:
1. Two top-level flat `TypeList` of `TypeString` attributes `reasons` / `regions` (like `references` in `resource_tc_teo_customize_error_page.go`, where the nested struct had only one field).
2. Top-level scalar `reason`/`region` fields.

**理由**:
- Each `ComplianceRestriction` carries a paired (Reason, Region) tuple; the region code only makes sense together with its reason. Two parallel flat lists would force users to correlate by index, and plain scalars would lose data when the API returns multiple restrictions.
- The nested computed block is the established pattern in this codebase for multi-field repeated API output (e.g. `references` in `resource_tc_teo_origin_group.go`, `ascription` in `resource_tc_teo_identify_zone_operation.go`).
- The task mapping (`Reason` → `reason`, `Region` → `region`) is naturally satisfied as sub-field names inside the block.

### Decision 2: Nil-guard the whole list and each sub-field

**选择**: Only set the attribute when the source pointer is non-nil:
- Skip the whole `d.Set("domain_compliance_restrictions", ...)` when `respData.DomainComplianceRestrictions == nil` (or len == 0), leaving state empty.
- Inside the loop, only add `reasonMap["reason"]` / `reasonMap["region"]` keys when the corresponding pointer is non-nil, matching the repo requirement "在调用setXX()设置字段前，请先判断Response中的字段是否为nil".

**理由**: The SDK marks these fields as possibly null (the list is only populated when the default domain is restricted). Avoids writing empty strings into state and matches existing conventions.

### Decision 3: No immutableArgs / mutableArgs changes

**选择**: Do not touch the Update method at all. The new block is Computed-only, so `d.HasChange` on it is not a trigger for `ModifyFunction`, and it is not added to `immutableArgs`.

**理由**: `ModifyFunction` request has no parameters for compliance restrictions — they are cloud-managed read-only output. Adding them to `immutableArgs` would cause spurious update failures.

### Decision 4: Unit tests via gomonkey mocks

**选择**: Add tests in `resource_tc_teo_function_test.go` using gomonkey to mock `DescribeFunctionsWithContext` on the teo client (following `resource_tc_teo_function_replica_test.go` / `resource_tc_teo_function_component_binding_test.go` patterns): one case where `DomainComplianceRestrictions` contains entries (assert both `reason` and `region` values in state), and one case where the list is nil/empty (assert the state is empty and no error occurs). Note the Read path goes through `DescribeTeoFunctionById` which calls `DescribeFunctions` (non-context variant), so the mock targets `DescribeFunctions` accordingly.

**理由**: Repo rule — for modified existing resources, unit tests must be runnable without real cloud credentials; gomonkey is the established approach for teo function resources.

### Decision 5: Documentation update only in the resource .md

**选择**: Update `tencentcloud/services/teo/resource_tc_teo_function.md` only if needed for example clarity (the Argument/Attribute Reference sections are auto-generated by `make doc`, which is executed by the tfpacer-finalize skill, not in this change).

**理由**: Repo rule forbids manual edits to `website/` and forbids running `make doc`/`gofmt` outside the finalize stage.

## Risks / Trade-offs

- **[Risk] API may return null or empty `DomainComplianceRestrictions` for functions without restrictions** → Mitigation: nil/length guard before `d.Set`; the state simply stays empty. Existing resources see no plan diff because Computed blocks are only refreshed.
- **[Risk] Multiple restriction entries could be returned** → Mitigation: the TypeList block preserves all entries in order.
- **[Risk] Adding a computed block affects existing state** → Mitigation: Computed-only additions are backward compatible; no state migration needed. Terraform populates the new attribute on the next refresh.
- **[Trade-off]** Users cannot filter or configure the restrictions (read-only by nature) — acceptable, as this information is cloud-managed compliance output.

## Migration Plan

- Pure additive change: new Computed block only. No migration required; existing TF configurations and state keep working.
- Rollback: revert the schema block and Read flattening code; state values for the block are simply dropped on the next refresh.

## Open Questions

- None.
