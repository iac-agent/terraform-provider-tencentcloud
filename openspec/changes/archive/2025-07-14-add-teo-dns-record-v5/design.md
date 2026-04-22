## Context

The Terraform Provider for TencentCloud already has a `tencentcloud_teo_dns_record` resource that manages EdgeOne DNS records. However, a V5 version (`tencentcloud_teo_dns_record_v5`) is needed to follow the latest code generation patterns, referencing `tencentcloud_igtm_strategy` as the template for RESOURCE_KIND_GENERAL resources.

The existing `tencentcloud_teo_dns_record` resource uses the same cloud APIs:
- CreateDnsRecord / DescribeDnsRecords / ModifyDnsRecords / ModifyDnsRecordsStatus / DeleteDnsRecords
- All from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`

The V5 resource will use the same API interfaces but with updated code structure and patterns.

Current state:
- Existing `resource_tc_teo_dns_record.go` in `tencentcloud/services/teo/`
- Existing `DescribeTeoDnsRecordById` service method in `service_tencentcloud_teo.go`
- The existing resource uses composite ID format: `zoneId + FILED_SP + recordId`
- The DnsRecord struct fields: ZoneId, RecordId, Name, Type, Location, Content, TTL, Weight, Priority, Status, CreatedOn, ModifiedOn

## Goals / Non-Goals

**Goals:**
- Create a new `tencentcloud_teo_dns_record_v5` resource following the latest code patterns (referencing `tencentcloud_igtm_strategy`)
- Support full CRUD lifecycle for TEO DNS records
- Support import of existing DNS records
- Handle the `status` field update via ModifyDnsRecordsStatus API separately from other mutable fields via ModifyDnsRecords API
- Register the new resource in provider.go and provider.md
- Add unit tests using gomonkey mock approach
- Add .md documentation file

**Non-Goals:**
- Do not modify the existing `tencentcloud_teo_dns_record` resource
- Do not add data source for DNS records (out of scope)
- Do not change the cloud API client or SDK
- Do not modify the DescribeTeoDnsRecordById service method (reuse existing one)

## Decisions

### 1. Resource naming and file organization
- Resource name: `tencentcloud_teo_dns_record_v5`
- File: `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.go`
- Test: `tencentcloud/services/teo/resource_tc_teo_dns_record_v5_test.go`
- Doc: `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.md`
- Rationale: Follows existing naming conventions with `_v5` suffix to distinguish from existing resource

### 2. Schema design
- `zone_id`: TypeString, Required, ForceNew (zone cannot be changed after creation)
- `name`: TypeString, Required (DNS record name)
- `type`: TypeString, Required (DNS record type: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content`: TypeString, Required (DNS record content)
- `location`: TypeString, Optional+Computed (resolution route, default "Default")
- `ttl`: TypeInt, Optional+Computed (cache time, default 300)
- `weight`: TypeInt, Optional+Computed (weight, default -1)
- `priority`: TypeInt, Optional+Computed (MX priority, default 0)
- `status`: TypeString, Optional+Computed (enable/disable status)
- `record_id`: TypeString, Computed (DNS record ID, set from CreateDnsRecord response)
- `created_on`: TypeString, Computed (creation time)
- `modified_on`: TypeString, Computed (modification time)
- Rationale: Matches the cloud API DnsRecord struct fields, with computed fields for output-only values

### 3. Composite ID format
- ID format: `zoneId + FILED_SP + recordId`
- In read/update/delete: parse ID using `strings.Split(d.Id(), tccommon.FILED_SP)` and get zone_id/record_id from `d.Get()` as well
- Rationale: Follows existing pattern in `tencentcloud_teo_dns_record` resource

### 4. Update strategy - two APIs
- Mutable fields (name, type, content, location, ttl, weight, priority): Use `ModifyDnsRecords` API
- Status field: Use `ModifyDnsRecordsStatus` API separately
- Rationale: The cloud API separates record content modification from status modification into different APIs. This matches the existing `tencentcloud_teo_dns_record` implementation.

### 5. Service layer reuse
- Reuse the existing `DescribeTeoDnsRecordById` method in `service_tencentcloud_teo.go`
- This method already uses `DescribeDnsRecords` with AdvancedFilter to query by record_id
- Rationale: No need to duplicate service layer code; the existing method provides the required functionality

### 6. Test approach
- Use gomonkey mock for unit testing (not terraform test suite)
- Mock the cloud API calls to test business logic
- Rationale: Required by project rules for new resources

## Risks / Trade-offs

- **Risk: Dual maintenance** - Two DNS record resources exist (`tencentcloud_teo_dns_record` and `tencentcloud_teo_dns_record_v5`). Mitigation: The V5 version follows better patterns and can eventually replace the old one. Users can choose which to use.
- **Risk: API compatibility** - DescribeDnsRecords returns a list, and we filter by record_id. If the record doesn't exist, we get an empty list. Mitigation: The existing service method already handles this by returning nil.
- **Trade-off: Reusing service method** - By reusing `DescribeTeoDnsRecordById`, we maintain consistency but accept that any bugs in the shared method affect both resources.
