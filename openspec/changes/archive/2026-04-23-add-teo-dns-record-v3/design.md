## Context

Terraform Provider for TencentCloud currently includes a `tencentcloud_teo_dns_record` resource that manages TEO DNS records. However, a new `tencentcloud_teo_dns_record_v3` resource is needed as a separate resource to provide a clean schema with full CRUD support using the latest cloud API. The existing `tencentcloud_teo_dns_record` resource supports `status` modification via `ModifyDnsRecordsStatus` API, while the new `v3` resource focuses purely on the core DNS record CRUD operations via `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, and `DeleteDnsRecords`.

All four APIs are synchronous (no async polling needed). The cloud SDK package is `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`.

## Goals / Non-Goals

**Goals:**
- Add a new `tencentcloud_teo_dns_record_v3` resource with full CRUD lifecycle management
- Support all fields from the `CreateDnsRecord` API: `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`
- Read computed fields from `DescribeDnsRecords`: `record_id`, `status`, `created_on`, `modified_on`
- Use composite ID format `zone_id#record_id` for resource identification
- Follow the existing TEO resource patterns (same as `tencentcloud_teo_dns_record`)
- Support import via `terraform import`
- Add unit tests using gomonkey mock approach
- Add `.md` documentation file

**Non-Goals:**
- Do NOT modify the existing `tencentcloud_teo_dns_record` resource
- Do NOT support `ModifyDnsRecordsStatus` API in the v3 resource (status is read-only/computed)
- Do NOT add a data source for DNS records (out of scope)

## Decisions

### Decision 1: Resource naming `tencentcloud_teo_dns_record_v3`
- **Choice**: Use `v3` suffix to distinguish from the existing `tencentcloud_teo_dns_record`
- **Rationale**: The new resource provides a clean separation. The `v3` naming follows the codebase convention (e.g., `tencentcloud_teo_l7_acc_rule_v2`)
- **Alternative**: Could modify the existing resource, but that risks breaking existing Terraform configurations

### Decision 2: Composite ID with `zone_id#record_id`
- **Choice**: Use `zone_id` and `record_id` joined by `tccommon.FILED_SP` (`#`)
- **Rationale**: Follows the exact same pattern as `tencentcloud_teo_dns_record`. Both `zone_id` and `record_id` are needed for all CRUD operations
- **Alternative**: Single ID with just `record_id` — rejected because the cloud API requires `zone_id` for all operations

### Decision 3: `status` field is Computed-only (read-only)
- **Choice**: Make `status` a Computed-only field in the v3 resource schema
- **Rationale**: The v3 resource does not call `ModifyDnsRecordsStatus`. The `status` is returned by `DescribeDnsRecords` as output and should be readable but not settable through this resource. This simplifies the resource and avoids the dual-update-path complexity of the original resource
- **Alternative**: Include status modification — rejected to keep the resource focused on core DNS record fields

### Decision 4: `zone_id` is Required + ForceNew
- **Choice**: `zone_id` is Required and ForceNew
- **Rationale**: `zone_id` identifies the TEO zone and is required for all API operations. Changing the zone_id means a completely different resource, so it must trigger recreation
- **Alternative**: Mutable zone_id — rejected because the cloud API doesn't support moving records between zones

### Decision 5: Mutable fields for Update
- **Choice**: `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority` are all mutable
- **Rationale**: The `ModifyDnsRecords` API accepts all these fields in its `DnsRecord` input struct, so they can all be updated
- **Alternative**: Making some fields ForceNew — rejected because the API supports modifying them

### Decision 6: Service layer method `DescribeTeoDnsRecordV3ById`
- **Choice**: Add a new `DescribeTeoDnsRecordV3ById` method to `TeoService`
- **Rationale**: Uses the `DescribeDnsRecords` API with `AdvancedFilter` on `id` to find a specific record by `record_id`. A separate method avoids interfering with the existing `DescribeTeoDnsRecordById`
- **Alternative**: Reuse existing method — rejected to avoid side effects on the existing resource

## Risks / Trade-offs

- **[Risk] Duplicate resource with existing `tencentcloud_teo_dns_record`** → Mitigation: The v3 resource has a simpler schema (no status modification), making it easier to maintain. Users can choose either resource.
- **[Risk] API rate limiting on DescribeDnsRecords** → Mitigation: Use `ratelimit.Check()` before API calls as per existing patterns, and use `ReadRetryTimeout` for read operations.
- **[Risk] `priority` field only applies to MX type records** → Mitigation: The cloud API validates this server-side. The Terraform schema description should document this constraint.
- **[Trade-off] Not supporting status modification in v3** → This simplifies the resource but means users who need to toggle DNS record status must use the original resource or manage it separately.
