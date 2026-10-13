## Context

TencentCloud EdgeOne (TEO) provides DNS management capabilities for zones. An existing `tencentcloud_teo_dns_record` resource already exists in the Terraform provider, but a new version `tencentcloud_teo_dns_record_17` is needed as a separate resource variant. The underlying cloud APIs (CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords) are all synchronous and already vendored in the SDK.

The existing `tencentcloud_teo_dns_record` resource follows the standard RESOURCE_KIND_GENERAL pattern with CRUD operations, composite ID (zone_id#record_id), and uses the TeoService's `DescribeTeoDnsRecordById` for reading individual records. The new resource `tencentcloud_teo_dns_record_17` will follow the same pattern.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_teo_dns_record_17` as a new RESOURCE_KIND_GENERAL resource
- Support full CRUD lifecycle: Create, Read, Update, Delete
- Support resource import using composite ID (zone_id#record_id)
- Support all DNS record fields: zone_id, name, type, content, location, ttl, weight, priority
- Include computed fields: status, created_on, modified_on
- Update handling for status changes via ModifyDnsRecordsStatus API
- Add unit tests using gomonkey mock approach
- Add resource documentation markdown

**Non-Goals:**
- Modifying the existing `tencentcloud_teo_dns_record` resource
- Adding datasource support (not requested)
- Supporting batch operations (each resource manages a single DNS record)

## Decisions

### 1. Resource File Naming
**Decision**: Name the file `resource_tc_teo_dns_record_17.go` following the convention `resource_tc_<Product>_<name>.go`.

**Rationale**: Consistent with existing naming conventions in the TEO service directory.

### 2. Composite Resource ID
**Decision**: Use `zone_id + FILED_SP + record_id` as the composite ID, same as the existing dns_record resource.

**Rationale**: The CreateDnsRecord API requires zone_id, and the record_id is returned from the create response. Both are needed for Read/Update/Delete operations. The `#` separator (tccommon.FILED_SP) is the standard pattern for composite IDs.

### 3. Read Implementation
**Decision**: Reuse the existing `TeoService.DescribeTeoDnsRecordById()` method from the service layer.

**Rationale**: This method already implements the correct pattern: calling DescribeDnsRecords with a filter on record ID, handling pagination, and returning a single DnsRecord struct. No duplication needed.

### 4. Update Implementation
**Decision**: Split update into two operations:
1. For mutable fields (name, type, content, location, ttl, weight, priority): use ModifyDnsRecords API
2. For status field: use ModifyDnsRecordsStatus API

**Rationale**: This follows the exact same pattern as the existing `tencentcloud_teo_dns_record` resource. The ModifyDnsRecords API handles content changes while status changes go through a separate API.

### 5. Schema Fields
**Decision**: Use the same schema as the existing dns_record resource, with:
- Required: zone_id (ForceNew), name, type, content
- Optional+Computed: location, ttl, weight, priority, status
- Computed: created_on, modified_on

**Rationale**: Directly reflects the cloud API parameters. The zone_id is ForceNew since changing the zone means a different resource. status is Optional+Computed to allow both reading the current status and explicitly setting it.

### 6. Test Strategy
**Decision**: Use gomonkey mock approach for unit tests (not Terraform acceptance tests).

**Rationale**: Per the project requirements, new resources must use mock-based unit tests with gomonkey rather than Terraform test suites.

## Risks / Trade-offs

- **[API Consistency]** The DescribeDnsRecords API returns a list and uses filtering to find a specific record. If the record ID filter doesn't work correctly, Read might fail → Mitigation: This is already handled by the existing `DescribeTeoDnsRecordById` service method which has been proven to work.
- **[Naming Confusion]** Having both `tencentcloud_teo_dns_record` and `tencentcloud_teo_dns_record_17` might confuse users → Mitigation: The `_17` suffix differentiates the resources; this is a deliberate versioning decision.
- **[Backward Compatibility]** The new resource must not modify any existing resource schemas or behavior → Mitigation: All changes are additive (new files, new provider registration entry).
