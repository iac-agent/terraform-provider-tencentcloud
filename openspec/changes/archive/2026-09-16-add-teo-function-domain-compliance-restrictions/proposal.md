## Why

The cloud API `DescribeFunctions` returns domain compliance restriction information for each edge function through the `Function.DomainComplianceRestrictions` field (JsonPath: `response.Functions.DomainComplianceRestrictions`), where each `ComplianceRestriction` entry carries a `Reason` (JsonPath: `response.Functions.DomainComplianceRestrictions.Reason`) and a `Region` (JsonPath: `response.Functions.DomainComplianceRestrictions.Region`). This information tells users in which regions the TEO edge function default domain is inaccessible (e.g. due to ICP record not obtained or government order). The `tencentcloud_teo_function` resource currently does not expose this information, so users cannot know the regional access restrictions of the function default domain through Terraform, which may cause confusion when the default domain is unreachable from certain regions.

## What Changes

- Add 2 new computed attributes to the `tencentcloud_teo_function` resource, exposed through a new `domain_compliance_restrictions` computed list block (each element is a flattened `ComplianceRestriction`):
  - `reason`: the reason the access restriction was issued (maps to `ComplianceRestriction.Reason`, enum values: `ICP_RECORD_REQUIRED` (ICP filing not obtained), `GOVERNMENT_ORDER` (government order))
  - `region`: the country/region code where access is restricted, in ISO 3166 format (maps to `ComplianceRestriction.Region`)
- Update the resource Read method to flatten the `DomainComplianceRestrictions` list from the `DescribeFunctions` response into the new computed block, following the nil-check pattern used for the existing `references`-style computed blocks.
- This is a read-only computed field sourced from the existing `DescribeFunctions` API response; no changes to Create/Update/Delete operations and no new API calls.

## Capabilities

### New Capabilities
- `teo-function-domain-compliance-restrictions`: Expose the `DomainComplianceRestrictions` entries (with `reason` and `region` attributes) returned by the `DescribeFunctions` API as a computed `domain_compliance_restrictions` block of the `tencentcloud_teo_function` resource.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_function.go`: add the `domain_compliance_restrictions` computed list block to the schema (with `reason` and `region` sub-fields) and the flattening logic in the Read method.
- `tencentcloud/services/teo/resource_tc_teo_function_test.go`: add mock (gomonkey) unit tests for the new computed fields.
- `tencentcloud/services/teo/resource_tc_teo_function.md`: update the example documentation.
- `tencentcloud/provider.go` / `tencentcloud/provider.md`: no registration change needed (the resource is already registered).
- APIs: no new API calls; the existing `DescribeTeoFunctionById` service method (wrapping `DescribeFunctions`) already returns `DomainComplianceRestrictions` in the vendored SDK `teov20220901.Function` struct.
- Dependencies: none — the vendored SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` already contains the `ComplianceRestriction` struct with `Reason` and `Region` fields.
- Backward compatibility: fully backward compatible; only a new Computed block is added, existing TF configurations and state are unaffected.
