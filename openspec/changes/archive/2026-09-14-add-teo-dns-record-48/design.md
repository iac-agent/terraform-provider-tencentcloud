## Context

TEO (TencentCloud EdgeOne) provides DNS record management through its cloud API. Currently, the Terraform provider does not have a resource to manage individual TEO DNS records, forcing users to manage DNS records through the console or API directly.

The TEO DNS record API supports:
- **CreateDnsRecord**: Creates a single DNS record, returns `RecordId`
- **DescribeDnsRecords**: Lists DNS records by `ZoneId` with optional filters, supports pagination (max Limit=1000)
- **ModifyDnsRecords**: Batch modifies DNS records by `ZoneId` + `DnsRecords` array (requires `RecordId` for each record to modify)
- **DeleteDnsRecords**: Deletes DNS records by `ZoneId` + `RecordIds` array

Key SDK structure: `DnsRecord` contains fields: `ZoneId`, `RecordId`, `Name`, `Type`, `Location`, `Content`, `TTL`, `Weight`, `Priority`, `Status`, `CreatedOn`, `ModifiedOn`. Note that `ZoneId`, `Status`, `CreatedOn`, and `ModifiedOn` are output-only fields in `ModifyDnsRecords` (ignored if passed as input).

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_teo_dns_record_48` resource with full CRUD lifecycle
- Support all input parameters: `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`
- Expose computed attributes: `record_id`, `status`, `created_on`, `modified_on`
- Use composite ID (`zone_id` + `record_id`) for resource identification
- Follow existing Terraform provider patterns (reference `tencentcloud_igtm_strategy` resource style)

**Non-Goals:**
- Batch DNS record operations (each resource manages a single record)
- DNS record status management (enable/disable via `ModifyDnsRecordsStatus` is out of scope)
- Custom timeout configurations beyond default retry settings

## Decisions

1. **Composite Resource ID**: Use `zone_id#record_id` as the Terraform resource ID (using `tccommon.FILED_SP` separator). This allows the Read function to reconstruct API requests from `d.Id()` without requiring additional state lookups.
   - Alternative: Use `record_id` alone — rejected because Read API requires `ZoneId`, making it impossible to query without also storing `zone_id`.

2. **Read Implementation**: Use `DescribeDnsRecords` with filter on `record-id` to read a single record. The API returns a list; match by `RecordId` to find the specific record.
   - Alternative: There is no single-record read API, so filtering is the only option.

3. **Update Implementation**: Use `ModifyDnsRecords` with a single-element `DnsRecords` array containing the `RecordId` and modified fields. Only include updatable fields (`Name`, `Type`, `Location`, `Content`, `TTL`, `Weight`, `Priority`) — exclude read-only fields (`ZoneId`, `Status`, `CreatedOn`, `ModifiedOn`) as documented by the SDK.
   - The API documentation states these read-only fields are ignored in ModifyDnsRecords, so omitting them is correct.

4. **Delete Implementation**: Use `DeleteDnsRecords` with the `ZoneId` and `RecordIds` array containing the single record ID.

5. **Computed-only fields**: `status`, `created_on`, `modified_on` are computed-only (not user-configurable). `record_id` is also computed (set by Create API response).

## Risks / Trade-offs

- [DescribeDnsRecords returns a list, not a single record] → Mitigation: Filter by `record-id` and match the first result. If not found, treat as resource deleted.
- [ModifyDnsRecords is a batch API but we use it for single record updates] → Mitigation: Always pass a single-element array, which is a valid usage pattern.
- [No single-record read API exists] → Mitigation: Use `DescribeDnsRecords` with `record-id` filter for efficient single-record lookup.
