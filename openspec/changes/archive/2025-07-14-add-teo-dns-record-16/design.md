## Context

TencentCloud EdgeOne (TEO) provides DNS record management for edge domains. The existing `tencentcloud_teo_dns_record` resource manages DNS records but includes `status` as a modifiable field via the `ModifyDnsRecordsStatus` API. In certain scenarios, managing status through Terraform causes unintended side effects (e.g., re-enabling records that were disabled through the console).

The new `tencentcloud_teo_dns_record_16` resource simplifies the design by treating `status` as a computed-only field, removing the `ModifyDnsRecordsStatus` API call from the update flow. This aligns with the principle that Terraform should manage infrastructure declarations rather than operational state toggles.

The resource follows the standard TEO resource patterns, using `UseTeoV20220901Client()` for API access and composite ID format `zone_id + FILED_SP + record_id`.

## Goals / Non-Goals

**Goals:**
- Provide a new TEO DNS record resource with full CRUD lifecycle management
- Simplify the resource by removing status management (computed-only)
- Support resource import using composite ID format
- Follow existing TEO resource code patterns and conventions
- Include comprehensive unit tests using gomonkey mock approach

**Non-Goals:**
- Modify or deprecate the existing `tencentcloud_teo_dns_record` resource
- Add data source for DNS records (out of scope for this change)
- Support batch operations (the resource manages a single DNS record)
- Manage DNS record status enable/disable operations

## Decisions

### 1. Resource ID Format: Composite ID (zone_id + FILED_SP + record_id)
**Rationale**: The `DescribeDnsRecords` API requires `ZoneId` as a required parameter to locate a specific DNS record. Using a composite ID ensures both the zone and record identifiers are available during Read, Update, and Delete operations. This matches the pattern used by the existing `tencentcloud_teo_dns_record` resource.

**Alternative considered**: Single record_id — rejected because all CRUD APIs require ZoneId, and without it the Read operation cannot function.

### 2. Status as Computed-Only Field
**Rationale**: Unlike the existing `tencentcloud_teo_dns_record` which uses `ModifyDnsRecordsStatus` to toggle record status, this resource treats `status` as read-only. This prevents Terraform from overriding operational status changes made through the console or other tools.

**Alternative considered**: Making status modifiable like the existing resource — rejected because this is the key differentiation of the new resource variant.

### 3. Read Implementation: Direct API Call vs TeoService Helper
**Rationale**: The new resource will call `DescribeDnsRecords` directly with filtering by record ID, rather than using the `TeoService.DescribeTeoDnsRecordById` helper. This gives more control over the filtering logic and avoids dependency on the service layer for a simple lookup. The filter will use `AdvancedFilter` with field "id" to locate the specific record.

**Alternative considered**: Using `TeoService.DescribeTeoDnsRecordById` — viable but adds unnecessary indirection; direct API call is cleaner for a new resource.

### 4. Update Implementation: ModifyDnsRecords with DnsRecord Struct
**Rationale**: The `ModifyDnsRecords` API accepts a list of `DnsRecord` objects. For a single-record resource, we construct one `DnsRecord` with `RecordId` set to identify the record and populate the changed mutable fields (name, type, content, location, ttl, weight, priority). Output-only fields (ZoneId, Status, CreatedOn, ModifiedOn) in the DnsRecord struct are not set.

### 5. Delete Implementation: DeleteDnsRecords with RecordIds
**Rationale**: The `DeleteDnsRecords` API accepts `ZoneId` and a list of `RecordIds`. For a single-record resource, we pass a single-element list containing the record ID extracted from the composite ID.

## Risks / Trade-offs

- **[Risk] Record not found after creation** → Mitigation: The Create function calls Read immediately after setting the ID to verify the resource exists. The Read function sets `d.SetId("")` if the record is not found, signaling Terraform to recreate it.

- **[Risk] DescribeDnsRecords returns multiple records** → Mitigation: Use `AdvancedFilter` with field "id" to filter by record ID, ensuring only the target record is returned. Additionally, check that exactly one record is returned and log a warning if multiple records match.

- **[Risk] ModifyDnsRecords is a batch API** → Mitigation: For single-record management, always construct a single-element DnsRecord list. The API documentation confirms this is valid.

- **[Trade-off] Computed-only status means users cannot enable/disable records via Terraform** → This is intentional and aligns with the design goal of simplifying the resource. Users who need status management should use the existing `tencentcloud_teo_dns_record` resource.
