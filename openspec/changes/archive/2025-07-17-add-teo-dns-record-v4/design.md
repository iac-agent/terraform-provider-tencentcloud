## Context

The TencentCloud Terraform Provider already has `tencentcloud_teo_dns_record` resource that manages TEO DNS records. This resource uses `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `ModifyDnsRecordsStatus`, and `DeleteDnsRecords` APIs. The new `tencentcloud_teo_dns_record_v4` follows the v4 naming convention and provides a clean resource definition that mirrors the existing `tencentcloud_teo_dns_record` implementation pattern, ensuring full CRUD lifecycle support.

The resource manages a single DNS record within a TEO zone. The composite ID format is `zone_id#record_id`, using `tccommon.FILED_SP` as separator.

**Cloud API details (verified from vendor SDK):**
- `CreateDnsRecordRequest`: ZoneId, Name, Type, Content, Location, TTL, Weight, Priority → Response: RecordId
- `DescribeDnsRecordsRequest`: ZoneId, Offset, Limit, Filters([]*AdvancedFilter), SortBy, SortOrder, Match → Response: TotalCount, DnsRecords([]*DnsRecord)
- `ModifyDnsRecordsRequest`: ZoneId, DnsRecords([]*DnsRecord) → Response: RequestId
- `ModifyDnsRecordsStatusRequest`: ZoneId, RecordsToEnable, RecordsToDisable → Response: RequestId
- `DeleteDnsRecordsRequest`: ZoneId, RecordIds([]*string) → Response: RequestId

**DnsRecord struct fields**: ZoneId, RecordId, Name, Type, Location, Content, TTL, Weight, Priority, Status, CreatedOn, ModifiedOn
- ZoneId, Status, CreatedOn, ModifiedOn are output-only in ModifyDnsRecords (ignored as input)

**AdvancedFilter struct fields**: Name, Values([]*string), Fuzzy(*bool) — used in DescribeDnsRecords for filtering

## Goals / Non-Goals

**Goals:**
- Implement a new `tencentcloud_teo_dns_record_v4` resource with full CRUD lifecycle
- Follow the existing `tencentcloud_teo_dns_record` resource pattern for consistency
- Support all Create API input parameters: zone_id, name, type, content, location, ttl, weight, priority
- Support computed output fields: record_id, status, created_on, modified_on
- Support update of mutable fields (name, type, content, location, ttl, weight, priority) via ModifyDnsRecords
- Support update of status field via ModifyDnsRecordsStatus API
- Support import via composite ID (zone_id#record_id)
- Register the resource in provider.go and provider.md
- Generate .md example documentation file

**Non-Goals:**
- Modifying or deprecating the existing `tencentcloud_teo_dns_record` resource
- Supporting batch creation/modification of multiple DNS records in a single resource
- Adding data source for DNS records listing

## Decisions

### 1. Resource naming: `tencentcloud_teo_dns_record_v4`
**Rationale**: Follows the v4 suffix naming convention for TEO DNS record resources, consistent with other v4/v2 resources in the codebase (e.g., `tencentcloud_teo_l7_acc_rule_v2`).

### 2. Composite ID format: `zone_id#record_id`
**Rationale**: ZoneId is ForceNew and required for all API calls. RecordId is returned by Create API. Using `tccommon.FILED_SP` separator is the standard pattern across the provider. Both components are needed for Read, Update, and Delete operations.

### 3. Status field update via separate ModifyDnsRecordsStatus API
**Rationale**: The existing `tencentcloud_teo_dns_record` resource already follows this pattern — the ModifyDnsRecords API handles DNS record content changes, while ModifyDnsRecordsStatus API handles enable/disable status changes. This separation is required by the cloud API design.

### 4. Read operation uses DescribeDnsRecords with filter by RecordId
**Rationale**: There is no single-record Describe API. The DescribeDnsRecords API returns a list filtered by AdvancedFilter. We filter by `id` field to find the specific record.

### 5. Schema design follows existing tencentcloud_teo_dns_record pattern
**Rationale**: Consistency with existing resources. The v4 resource mirrors the schema of the existing resource with the same field names, types, and computed/optional settings.

### 6. Unit tests using gomonkey mocks
**Rationale**: Per project requirements, new resources must use gomonkey-based unit tests with mock cloud API responses, not Terraform acceptance test suites.

## Risks / Trade-offs

- **DescribeDnsRecords pagination**: The Read operation uses DescribeDnsRecords with filter by record ID. If the filter doesn't work as expected, the read could return incorrect results → Mitigation: Use AdvancedFilter with Name="id" and the record's ID value, which is explicitly supported per API documentation.
- **Duplicate resource with tencentcloud_teo_dns_record**: Having two resources for the same cloud API could cause user confusion → Mitigation: Clear documentation distinguishing v4 as the recommended resource. The v4 naming convention signals it as a newer version.
- **ModifyDnsRecordsStatus not included in original interface list**: The ModifyDnsRecordsStatus API is needed for status updates but wasn't in the provided interface mapping → Mitigation: This is already a proven pattern from the existing resource, and is required for full lifecycle support.
