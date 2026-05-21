## Context

The TencentCloud EdgeOne (TEO) service provides DNS record management through its v20220901 API. Currently, the Terraform provider for TencentCloud does not include a resource for managing TEO DNS records. Users need to manage DNS records for their EdgeOne zones via the Terraform provider to enable full infrastructure-as-code workflows.

The TEO DNS record APIs available in the vendor directory are:
- `CreateDnsRecord`: Creates a single DNS record, returns `RecordId`
- `DescribeDnsRecords`: Lists DNS records with filtering, sorting, and pagination support
- `ModifyDnsRecords`: Batch modifies DNS records (accepts a list of `DnsRecord` objects)
- `DeleteDnsRecords`: Batch deletes DNS records by `RecordIds`

The `DnsRecord` struct contains: `ZoneId`, `RecordId`, `Name`, `Type`, `Location`, `Content`, `TTL`, `Weight`, `Priority`, `Status`, `CreatedOn`, `ModifiedOn`.

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_teo_dns_record_24` as a RESOURCE_KIND_GENERAL Terraform resource
- Support full CRUD lifecycle: create, read, update, delete a single DNS record
- Support all DNS record types: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- Support optional fields: `location`, `ttl`, `weight`, `priority`
- Use `zone_id` + `record_id` as the composite resource ID
- Register the resource in provider.go and provider.md
- Add unit tests using gomonkey mocks (not terraform acceptance tests)
- Add documentation .md file

**Non-Goals:**
- Batch management of multiple DNS records in a single resource
- Managing DNS record status (enable/disable) via `ModifyDnsRecordsStatus`
- Datasource for listing DNS records (separate resource)

## Decisions

### Decision 1: Resource ID format
Use `zone_id#record_id` as the composite ID (using `tccommon.FILED_SP` separator). This is consistent with other TEO resources that use composite IDs.

**Rationale**: The `RecordId` alone is not globally unique without the `ZoneId` context. Using a composite ID ensures unambiguous identification and supports import.

### Decision 2: Update via ModifyDnsRecords
The update operation uses `ModifyDnsRecords` which accepts a list of `DnsRecord` objects. For a single record update, we pass a list with one element containing the `RecordId` and updated fields.

**Rationale**: There is no single-record modify API; `ModifyDnsRecords` is the only update interface available.

### Decision 3: Read via DescribeDnsRecords with id filter
The read operation uses `DescribeDnsRecords` with a filter on `id` (the RecordId) to retrieve the specific record.

**Rationale**: There is no single-record describe API. Using the `id` filter in `DescribeDnsRecords` is the correct approach.

### Decision 4: Computed fields
`record_id`, `status`, `created_on`, `modified_on` are Computed-only fields set from the API response. They are not user-configurable.

### Decision 5: ForceNew fields
`zone_id` is ForceNew since changing the zone requires creating a new record. `name` and `type` are also ForceNew as DNS record name and type changes require recreation.

**Rationale**: The TEO API does not support changing `ZoneId`, `Name`, or `Type` on an existing record.

## Risks / Trade-offs

- [Risk] `ModifyDnsRecords` is a batch API; if the API returns partial success, error handling may be complex → Mitigation: treat any non-nil error as a full failure and surface it to the user.
- [Risk] `DescribeDnsRecords` pagination: if there are many records, the target record may not appear in the first page → Mitigation: use the `id` filter to directly target the specific record, which should return at most one result.
- [Risk] `Weight` field has special semantics (-1 means no weight, 0 means no resolution) → Mitigation: document this in the schema description and .md file.

## Migration Plan

This is a new resource with no existing state. No migration is needed. Users can import existing DNS records using `terraform import tencentcloud_teo_dns_record_24.<name> <zone_id>#<record_id>`.

## Open Questions

None.
