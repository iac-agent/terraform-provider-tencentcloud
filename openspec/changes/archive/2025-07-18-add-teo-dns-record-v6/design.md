## Context

The Terraform Provider for TencentCloud already includes a `tencentcloud_teo_dns_record` resource that manages TEO DNS records using the `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `ModifyDnsRecordsStatus`, and `DeleteDnsRecords` APIs from the `teo/v20220901` SDK package.

A new `tencentcloud_teo_dns_record_v6` resource is being added following the v6 naming convention. This new resource provides the same CRUD functionality as the existing `tencentcloud_teo_dns_record` but with a clean, well-structured schema that aligns with the latest patterns used in other v6 resources.

The TEO DNS record API supports:
- **CreateDnsRecord**: Creates a DNS record with ZoneId, Name, Type, Content, Location, TTL, Weight, Priority; returns RecordId
- **DescribeDnsRecords**: Queries DNS records by ZoneId with optional Filters (AdvancedFilter: Name, Values, Fuzzy), SortBy, SortOrder, Match; paginated with Offset/Limit; returns TotalCount and DnsRecords list
- **ModifyDnsRecords**: Modifies DNS records by ZoneId and DnsRecords list (DnsRecord struct with RecordId, Name, Type, Content, Location, TTL, Weight, Priority; note: ZoneId, Status, CreatedOn, ModifiedOn are output-only and ignored in ModifyDnsRecords)
- **ModifyDnsRecordsStatus**: Enables/disables DNS records by ZoneId and RecordsToEnable/RecordsToDisable
- **DeleteDnsRecords**: Deletes DNS records by ZoneId and RecordIds list

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_teo_dns_record_v6` as a new RESOURCE_KIND_GENERAL resource
- Support full CRUD lifecycle: Create, Read, Update, Delete
- Use composite ID: `zone_id` + `record_id` joined by `tccommon.FILED_SP`
- Follow the same pattern as `tencentcloud_teo_dns_record` and `tencentcloud_igtm_strategy`
- Support `status` field update via `ModifyDnsRecordsStatus` API
- Register the resource in `provider.go` and `provider.md`
- Add service layer method `DescribeTeoDnsRecordV6ById` in `service_tencentcloud_teo.go`
- Add unit tests using gomonkey mock approach
- Add resource documentation `.md` file for `make doc` generation

**Non-Goals:**
- Do not modify the existing `tencentcloud_teo_dns_record` resource
- Do not add data source for DNS records (out of scope)
- Do not support batch operations (single record management only)

## Decisions

### 1. Composite ID Design
**Decision**: Use `zone_id` + `record_id` joined by `tccommon.FILED_SP` as the resource ID.
**Rationale**: This follows the existing pattern in `tencentcloud_teo_dns_record` and other TEO resources. The zone_id is required for all API calls, and record_id uniquely identifies the record within a zone.

### 2. Schema Fields
**Decision**: Schema will include:
- Required: `zone_id` (ForceNew), `name`, `type`, `content`
- Optional+Computed: `location`, `ttl`, `weight`, `priority`
- Optional+Computed: `status` (updated via ModifyDnsRecordsStatus)
- Computed: `record_id`, `created_on`, `modified_on`

**Rationale**: Mirrors the existing `tencentcloud_teo_dns_record` schema. The `record_id` is computed (returned from Create API response). The `zone_id` is ForceNew because changing the zone requires recreating the record.

### 3. Update Strategy
**Decision**: Split update into two operations:
1. For mutable fields (name, type, content, location, ttl, weight, priority): use `ModifyDnsRecords`
2. For status field: use `ModifyDnsRecordsStatus`

**Rationale**: This matches the existing `tencentcloud_teo_dns_record` pattern where status changes are handled by a separate API.

### 4. Read Strategy
**Decision**: Add a service layer method `DescribeTeoDnsRecordV6ById` that calls `DescribeDnsRecords` with a filter on record ID, then finds the matching record.
**Rationale**: There is no single-record Describe API; `DescribeDnsRecords` returns a list. We filter by record ID and return the first match. Pagination is handled by setting Limit=1000 (API max).

### 5. Test Strategy
**Decision**: Use gomonkey mock approach for unit tests, not Terraform acceptance tests.
**Rationale**: Per project requirements, new resources should use gomonkey mocks for business logic unit testing rather than TF_ACC acceptance tests.

### 6. DnsRecord Struct in ModifyDnsRecords
**Decision**: When building the ModifyDnsRecords request, only include mutable input fields (Name, Type, Content, Location, TTL, Weight, Priority, RecordId). Exclude output-only fields (ZoneId, Status, CreatedOn, ModifiedOn) as documented in the SDK.
**Rationale**: The DnsRecord struct documentation explicitly states that ZoneId, Status, CreatedOn, ModifiedOn are output-only and will be ignored in ModifyDnsRecords. Including them would be unnecessary.

## Risks / Trade-offs

- **[Risk] DescribeDnsRecords uses list API, not single-record API** → Mitigation: Filter by record ID using AdvancedFilter to narrow results. If record not found, treat as resource deleted (set d.SetId("")).
- **[Risk] Existing `tencentcloud_teo_dns_record` resource coexistence** → Mitigation: The new resource has a distinct name (`tencentcloud_teo_dns_record_v6`) and will not conflict with the existing one. Users can choose which to use.
- **[Risk] ModifyDnsRecords is a batch API** → Mitigation: Only send a single record in the batch for this resource, following the existing `tencentcloud_teo_dns_record` pattern.
