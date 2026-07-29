## Context

The TEO product already has multiple versioned DNS record resources (`tencentcloud_teo_dns_record` through `tencentcloud_teo_dns_record_24`). This change adds version `_25` which is a simplified variant focusing on core CRUD parameters without `status` field management and without the nested `dns_records` list in the top-level schema. The existing `DescribeTeoDnsRecordById` service function in `service_tencentcloud_teo.go` is reused for reading individual records.

The upstream vendor SDK at `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` already provides all four required API methods:
- `CreateDnsRecord` / `CreateDnsRecordWithContext`
- `DescribeDnsRecords` / `DescribeDnsRecordsWithContext`
- `ModifyDnsRecords` / `ModifyDnsRecordsWithContext`
- `DeleteDnsRecords` / `DeleteDnsRecordsWithContext`

## Goals / Non-Goals

**Goals:**
- Create a new `tencentcloud_teo_dns_record_25` resource with full CRUD support
- Support core parameters: `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`, `record_id`
- The resource ID is a composite of `zone_id#record_id` using `tccommon.FILED_SP`
- Follow the same code style as `resource_tc_teo_dns_record_24.go`

**Non-Goals:**
- No `status` field management (no `ModifyDnsRecordsStatus` call)
- No `dns_records` nested list in the top-level schema
- No `filters`, `sort_by`, `sort_order`, `match` input parameters (these are DescribeDnsRecords inputs, not resource-level attributes)
- No modification to existing resources

## Decisions

1. **Resource schema**: Only include the parameters specified in the API mapping:
   - Required: `zone_id` (ForceNew), `name`, `type`, `content`
   - Optional/Computed: `location`, `ttl`, `weight`, `priority`
   - Computed: `record_id`

2. **Create flow**: Call `CreateDnsRecord`, validate response is non-nil, extract `RecordId`, set composite ID `zone_id#record_id`, then call Read.

3. **Read flow**: Parse composite ID, call `DescribeTeoDnsRecordById`, populate all schema fields. If not found, call `d.SetId("")`.

4. **Update flow**: Detect changes on mutable fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`), call `ModifyDnsRecords` with a single `DnsRecord` entry containing the `RecordId` and updated fields.

5. **Delete flow**: Parse composite ID, call `DeleteDnsRecords` with the single `RecordId`.

6. **Import support**: Use `schema.ImportStatePassthrough` with composite ID format `zone_id#record_id`.

7. **No `status` field**: Unlike the base resource, this variant does not expose or manage the `status` field.

## Risks / Trade-offs

- [Risk] `ModifyDnsRecords` is a batch API but only used for a single record → Mitigation: This is the same pattern used by all existing DNS record resources (e.g., `_24`), proven to work correctly.
- [Risk] Parameters like `ZoneId`, `Status`, `CreatedOn`, `ModifiedOn` in the `DnsRecord` struct are marked as "output only" for `ModifyDnsRecords` → Mitigation: The update request only sets `RecordId` and the mutable fields, avoiding these output-only fields.