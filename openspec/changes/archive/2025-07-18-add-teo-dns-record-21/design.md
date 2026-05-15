## Context

The TEO (EdgeOne) service provides DNS record management for edge zones. An existing resource `tencentcloud_teo_dns_record` already implements this, but a new version `tencentcloud_teo_dns_record_21` is needed. The implementation follows the standard Terraform Provider pattern for TencentCloud:

- Service layer in `service_tencentcloud_teo.go` for API calls (the `DescribeTeoDnsRecordById` method already exists and can be reused)
- Resource layer in `resource_tc_teo_dns_record_21.go` for Terraform lifecycle
- Uses `tencentcloud-sdk-go` TEO client (`teov20220901`) for API interactions

**API Summary:**
- `CreateDnsRecord`: Creates a single DNS record with parameters: ZoneId, Name, Type, Content, Location, TTL, Weight, Priority. Returns RecordId.
- `DescribeDnsRecords`: Queries DNS records with pagination and filters (AdvancedFilter). Returns list of DnsRecord objects with fields: ZoneId, RecordId, Name, Type, Content, Location, TTL, Weight, Priority, Status, CreatedOn, ModifiedOn.
- `ModifyDnsRecords`: Batch modifies DNS records. Takes ZoneId and list of DnsRecord objects (RecordId required, plus mutable fields).
- `DeleteDnsRecords`: Batch deletes DNS records. Takes ZoneId and list of RecordIds.

**Constraints:**
- `CreateDnsRecord` only creates a single record per call (not batch)
- `ModifyDnsRecords` and `DeleteDnsRecords` are batch APIs but the Terraform resource manages a single record
- `DescribeDnsRecords` uses `AdvancedFilter` (with Fuzzy field), not the simple `Filter` struct
- Must follow existing TEO resource patterns for consistency
- Must maintain backward compatibility (new resource, no breaking changes)

## Goals / Non-Goals

**Goals:**
- Provide Terraform resource `tencentcloud_teo_dns_record_21` for TEO DNS record lifecycle management
- Support standard Terraform operations: create, read, update, delete
- Support import of existing DNS records
- Include retry logic for eventual consistency
- Provide comprehensive unit test coverage using mock (gomonkey)
- Generate proper documentation

**Non-Goals:**
- Batch DNS record management (each resource manages one record)
- Data source for querying DNS records (separate concern)
- Management of DNS record status (enable/disable via `ModifyDnsRecordsStatus`) - not part of CRUD lifecycle

## Decisions

### 1. Resource Naming: `tencentcloud_teo_dns_record_21`

**Decision:** Use `tencentcloud_teo_dns_record_21` as the resource name with `_21` suffix to distinguish from the existing `tencentcloud_teo_dns_record`.

**Rationale:**
- Required by the specification to create a new versioned resource
- File naming: `resource_tc_teo_dns_record_21.go`

### 2. Composite ID Strategy: zone_id + record_id

**Decision:** Use composite ID with `tccommon.FILED_SP` separator joining zone_id and record_id.

**Rationale:**
- Consistent with the existing `tencentcloud_teo_dns_record` resource pattern
- All CRUD APIs require ZoneId as a parameter alongside the record-specific data
- Record ID alone is not sufficient for API calls

**Alternative considered:** Single record_id - rejected because all API calls require ZoneId.

### 3. Schema Design

**Decision:**
- Required fields: `zone_id` (ForceNew), `name`, `type`, `content`
- Optional fields: `location`, `ttl`, `weight`, `priority`
- Computed fields: `record_id`, `status`, `created_on`, `modified_on`

**Rationale:**
- `zone_id` is ForceNew because a DNS record belongs to a specific zone
- `name`, `type`, `content` are required as they are essential DNS record properties
- `location`, `ttl`, `weight`, `priority` are optional with API defaults
- `record_id`, `status`, `created_on`, `modified_on` are read-only computed fields from the API response

### 4. Update Strategy

**Decision:** Use `ModifyDnsRecords` API for updates, sending a single DnsRecord object with the changed fields.

**Rationale:**
- `ModifyDnsRecords` accepts a list of DnsRecord objects, each with RecordId and mutable fields
- For a single-record resource, send a list with one entry
- Only call the update API when relevant fields have changed (detected via `d.HasChange()`)

### 5. Service Layer Reuse

**Decision:** Reuse the existing `DescribeTeoDnsRecordById` method from `service_tencentcloud_teo.go` for the read operation.

**Rationale:**
- The method already exists and handles filtering by record ID using `AdvancedFilter`
- Avoids code duplication
- Consistent with the existing resource pattern

### 6. Error Handling and Retry

**Decision:** Use `resource.Retry` with `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations.

**Rationale:**
- Standard pattern used throughout the provider
- Handles eventual consistency issues common in cloud APIs
- Uses `tccommon.RetryError()` to wrap API errors for retry framework

## Risks / Trade-offs

### Risk: API Rate Limits
**Impact:** Rapid create/update/delete cycles could hit TEO API rate limits.
**Mitigation:**
- Use `ratelimit.Check()` before each API call (if applicable)
- Implement proper retry logic with exponential backoff

### Risk: Existing Resource Coexistence
**Impact:** Both `tencentcloud_teo_dns_record` and `tencentcloud_teo_dns_record_21` manage the same cloud resource type.
**Mitigation:**
- Clear documentation that both resources manage the same underlying API
- No state migration needed - users choose which resource to use

### Trade-off: No Status Management
**Decision:** Not implementing `ModifyDnsRecordsStatus` in this resource.
**Rationale:**
- Status (enable/disable) is a separate operation from CRUD lifecycle
- Keeps the resource focused on DNS record content management
- Can be added incrementally if needed
