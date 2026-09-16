# teo-function-domain-compliance-restrictions Specification

## Purpose
TBD - created by archiving change add-teo-function-domain-compliance-restrictions. Update Purpose after archive.
## Requirements
### Requirement: domain_compliance_restrictions computed block

The `tencentcloud_teo_function` resource SHALL expose a `domain_compliance_restrictions` computed list block at the top level of the schema. Each element of the block SHALL be a flattened `ComplianceRestriction` from the `DescribeFunctions` API response (JsonPath: `response.Functions.DomainComplianceRestrictions`) and SHALL contain the following computed string attributes:

- `reason`: the reason the regional access restriction was issued, mapped from `ComplianceRestriction.Reason` (JsonPath: `response.Functions.DomainComplianceRestrictions.Reason`). Valid values: `ICP_RECORD_REQUIRED` (ICP filing not obtained), `GOVERNMENT_ORDER` (government order).
- `region`: the country/region code where access to the function default domain is restricted, in ISO 3166 format, mapped from `ComplianceRestriction.Region` (JsonPath: `response.Functions.DomainComplianceRestrictions.Region`).

The block SHALL be computed-only and SHALL NOT be user-configurable.

#### Scenario: Read populates domain_compliance_restrictions when API returns entries

- **WHEN** the resource Read method calls `DescribeFunctions` and the returned `Function` contains `DomainComplianceRestrictions` with one or more non-nil entries
- **THEN** each entry SHALL be flattened into a `domain_compliance_restrictions` block element containing `reason` and `region` attributes matching the API response values, in the same order as returned by the API

#### Scenario: Read populates multiple restriction entries

- **WHEN** the API returns a `Function` with `DomainComplianceRestrictions` containing multiple entries (e.g. `[{Reason: "ICP_RECORD_REQUIRED", Region: "CN"}, {Reason: "GOVERNMENT_ORDER", Region: "XX"}]`)
- **THEN** the `domain_compliance_restrictions` block SHALL contain one element per API entry, preserving each paired `reason`/`region` value

#### Scenario: Read handles nil DomainComplianceRestrictions

- **WHEN** the API returns a `Function` whose `DomainComplianceRestrictions` field is nil or empty
- **THEN** the Read method SHALL skip setting the `domain_compliance_restrictions` attribute (no error SHALL be returned) and the resource state SHALL keep the attribute empty

#### Scenario: Read handles nil sub-fields inside an entry

- **WHEN** the API returns a `ComplianceRestriction` entry whose `Reason` or `Region` pointer is nil
- **THEN** the corresponding `reason` or `region` key SHALL be skipped for that block element (nil fields are not written as empty strings)

#### Scenario: domain_compliance_restrictions is computed and not user-settable

- **WHEN** a user creates or updates a `tencentcloud_teo_function` resource
- **THEN** the `domain_compliance_restrictions` block SHALL be computed-only and cannot be set by the user in the Terraform configuration, and Create/Update/Delete operations SHALL NOT send any compliance-restriction-related parameters to the cloud API

#### Scenario: Backward compatibility with existing state

- **WHEN** an existing `tencentcloud_teo_function` resource state (created before this change) is refreshed
- **THEN** all existing fields (`zone_id`, `function_id`, `name`, `remark`, `content`, `domain`, `create_time`, `update_time`) SHALL continue to work unchanged, and the new `domain_compliance_restrictions` block SHALL be populated on the next read without requiring any user configuration change or state migration

