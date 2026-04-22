## Why

Terraform Provider for TencentCloud currently has `tencentcloud_teo_dns_record` resource for managing TEO DNS records, but it does not support the `record_id` field as an explicit schema field, nor does it expose `status`, `created_on`, `modified_on` as computed fields in a v6-compatible schema structure. A new `tencentcloud_teo_dns_record_v6` resource is needed to provide a clean, well-structured resource definition that aligns with the latest TEO DNS API capabilities and follows the v6 resource naming convention for improved schema design.

## What Changes

- Add new Terraform resource `tencentcloud_teo_dns_record_v6` (RESOURCE_KIND_GENERAL) for managing TEO DNS records
- The resource will support full CRUD operations using the following cloud APIs:
  - **Create**: `CreateDnsRecord` - creates a DNS record with zone_id, name, type, content, location, ttl, weight, priority
  - **Read**: `DescribeDnsRecords` - queries DNS records by zone_id and filters, returning record details including record_id, status, created_on, modified_on
  - **Update**: `ModifyDnsRecords` - modifies DNS record fields (name, type, content, location, ttl, weight, priority) and status via `ModifyDnsRecordsStatus`
  - **Delete**: `DeleteDnsRecords` - deletes DNS records by zone_id and record_ids
- The resource will use composite ID: `zone_id` + `record_id` joined by `FILED_SP`
- Register the new resource in `provider.go` and `provider.md`
- Add unit tests using gomonkey mock approach
- Add resource documentation `.md` file

## Capabilities

### New Capabilities
- `teo-dns-record-v6`: Terraform resource for managing TEO DNS records with full CRUD support, including record creation, reading, updating (fields and status), and deletion. Supports fields: zone_id, name, type, content, location, ttl, weight, priority, record_id, status, created_on, modified_on.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **New files**:
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v6.go` - Resource definition and CRUD implementation
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v6_test.go` - Unit tests with gomonkey mocks
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v6.md` - Resource documentation
- **Modified files**:
  - `tencentcloud/provider.go` - Register new resource
  - `tencentcloud/provider.md` - Add resource documentation reference
  - `tencentcloud/services/teo/service_tencentcloud_teo.go` - Add DescribeTeoDnsRecordV6ById service method
- **Cloud APIs**: Uses `teo/v20220901` SDK package - CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, ModifyDnsRecordsStatus, DeleteDnsRecords
- **Dependencies**: No new external dependencies needed; uses existing SDK and helper packages
