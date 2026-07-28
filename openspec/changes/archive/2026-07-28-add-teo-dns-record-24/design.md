## Context

The TEO (Tencent EdgeOne) Terraform provider has a family of versioned DNS record resources (`tencentcloud_teo_dns_record`, `tencentcloud_teo_dns_record_21`, `tencentcloud_teo_dns_record_22`, `tencentcloud_teo_dns_record_23`). This change adds `tencentcloud_teo_dns_record_24` following the same established pattern for managing TEO DNS records via the v20220901 API.

The resource code, documentation template, and provider registration already exist in the codebase.

## Goals / Non-Goals

**Goals:**
- Provide full CRUD lifecycle management for a single TEO DNS record
- Support import using the compound ID format (`zone_id#record_id`)
- Expose computed read-only fields (`record_id`, `status`, `created_on`, `modified_on`)
- Follow the same code patterns as existing TEO DNS record resources (21, 23)

**Non-Goals:**
- Batch DNS record management (single record operations only)
- Support for the `dns_records` nested block in Create/Update/Delete operations (it is a computed output only)
- Modify the existing TEO DNS record resources

## Decisions

### 1. Resource ID Format
**Decision**: Use compound ID `zone_id#record_id` separated by `tccommon.FILED_SP`.

**Rationale**: This matches the pattern used by all existing TEO DNS record resources. The `DescribeDnsRecords` API requires both `zone_id` and `record_id` (via filters) to fetch a single record, so both values must be preserved in the Terraform state ID.

### 2. Schema Design
**Decision**: Top-level schema with `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority` as user inputs; `record_id`, `status`, `created_on`, `modified_on` as computed outputs; and `dns_records` as a computed nested block with all record fields.

**Rationale**: 
- The user-facing input fields (`name`, `type`, `content`, etc.) map directly to `CreateDnsRecord` API parameters
- The computed-only fields come from the Read API response
- The `dns_records` nested block reflects the `DescribeDnsRecords` API response structure for consistency

### 3. Read Pattern (DescribeTeoDnsRecordById)
**Decision**: Use the existing `TeoService.DescribeTeoDnsRecordById` service method instead of calling the API directly.

**Rationale**: This service method already handles the `DescribeDnsRecords` API call with proper filter construction (filtering by `record-id`) and pagination logic. Reusing it avoids code duplication across dns_record resources.

**Alternatives considered**: Direct API call - rejected because the service method already encapsulates the filter/pagination logic needed.

### 4. Update Strategy
**Decision**: Use `ModifyDnsRecords` API with a single `DnsRecord` entry, only calling the API when mutable fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) have changed.

**Rationale**: `ModifyDnsRecords` is a batch API, but since we manage a single record, we pass a single-item slice. Only triggering the API when changes exist avoids unnecessary API calls.

### 5. Delete Strategy
**Decision**: Use `DeleteDnsRecords` API with a single record ID string in `RecordIds`.

**Rationale**: `DeleteDnsRecords` is a batch API, but for single-record management, we pass a single record ID.

## Risks / Trade-offs

- **Risk**: The `ModifyDnsRecords` and `DeleteDnsRecords` APIs are batch operations, but the Terraform resource manages only one record at a time. If the same record is modified/deleted via other means (console, API, or another Terraform resource instance), there could be conflicts.
  - **Mitigation**: This is inherent to the API design and consistent with all existing TEO DNS record resources. Terraform's state management ensures it only tracks the record it created.

- **Risk**: The `dns_records` field in schema contains the full `DnsRecord` structure as a computed list, which may be redundant with the top-level fields.
  - **Mitigation**: This matches the existing pattern in resources like `tencentcloud_teo_dns_record_23` and provides a complete view of the API response for advanced users.
