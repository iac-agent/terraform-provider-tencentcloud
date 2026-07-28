## Context

The `tencentcloud_teo_dns_record_24` resource provides Terraform management of individual TEO (Tencent EdgeOne) DNS records. This follows the established pattern from previous versions (21, 22, 23) of the same resource, using the TEO v20220901 SDK. The resource code has already been implemented following the v23 pattern exactly.

The TEO DNS record APIs are:
- **CreateDnsRecord**: Creates a single DNS record and returns its `RecordId`
- **DescribeDnsRecords**: Lists DNS records with filtering, pagination, and sorting
- **ModifyDnsRecords**: Batch modifies DNS records (up to 100 at a time)
- **DeleteDnsRecords**: Batch deletes DNS records (up to 1000 at a time)

All APIs are synchronous (no async polling needed).

## Goals / Non-Goals

**Goals:**
- Provide full CRUD lifecycle (Create, Read, Update, Delete) for a single TEO DNS record
- Support Terraform import using composite ID `zone_id#record_id`
- Follow the same code pattern as `tencentcloud_teo_dns_record_23`

**Non-Goals:**
- Batch management of multiple DNS records in a single resource instance
- Support for DNS record status modification (separate `ModifyDnsRecordsStatus` API)

## Decisions

### 1. Resource ID: Composite ZoneId + RecordId

**Decision**: Use `zone_id` and `record_id` joined by `tccommon.FILED_SP` as the Terraform resource ID.

**Rationale**: The `DescribeDnsRecords` API requires `ZoneId` to query records. Storing both in the ID enables direct lookup in Read/Update/Delete without additional state storage.

### 2. Schema: Required + Optional + Computed fields

**Decision**: 
- `zone_id`: Required + ForceNew (zone cannot change for existing record)
- `name`, `type`, `content`: Required (core DNS record properties)
- `location`, `ttl`, `weight`, `priority`: Optional + Computed (have defaults from API)
- `record_id`, `status`, `created_on`, `modified_on`: Computed only (API output)
- `dns_records`: Computed List (mirrors DescribeDnsRecords response)

**Rationale**: Aligns with the API contract. `zone_id` is ForceNew because changing the zone would require deleting and recreating the record.

### 3. Read via Service Layer

**Decision**: Use the existing `DescribeTeoDnsRecordById` service method that filters by `id` filter on `DescribeDnsRecords`.

**Rationale**: Reuses existing service infrastructure. The method returns `*DnsRecord` or nil, and the Read function handles nil by calling `d.SetId("")`.

### 4. Update via ModifyDnsRecords with single DnsRecord

**Decision**: Wrap the single record update in a `ModifyDnsRecords` request with a single-element `DnsRecords` slice.

**Rationale**: The API only supports batch modification. Sending a single record is the standard approach for individual resource management.

### 5. Delete via DeleteDnsRecords with single RecordId

**Decision**: Wrap the single record deletion in a `DeleteDnsRecords` request with a single-element `RecordIds` slice.

**Rationale**: Same as update - the API only supports batch deletion.

## Risks / Trade-offs

- **[API Consistency]**: The ModifyDnsRecords API ignores certain fields (`ZoneId`, `Status`, `CreatedOn`, `ModifiedOn`) when used as input. The implementation correctly omits these from the DnsRecord struct in update requests.
- **[Error Handling]**: All API calls use `resource.Retry` with `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations, providing resilience against transient failures.
- **[State Drift]**: If a DNS record is deleted outside of Terraform, the Read function will detect the absence and call `d.SetId("")`, allowing Terraform to detect the drift and plan recreation.