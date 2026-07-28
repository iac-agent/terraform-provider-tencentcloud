## Why

TEO (Tencent EdgeOne) DNS records management requires a Terraform resource to support the complete lifecycle (CRUD) of individual DNS records. While previous versions (21, 22, 23) exist, version 24 is needed to align with the latest TEO v20220901 SDK API signatures and provide a properly structured general-purpose resource.

## What Changes

- Add new `tencentcloud_teo_dns_record_24` resource (`RESOURCE_KIND_GENERAL`) for managing TEO DNS records
- Support Create: `CreateDnsRecord` API for creating individual DNS records
- Support Read: `DescribeDnsRecords` API for querying DNS record details by zone and record ID
- Support Update: `ModifyDnsRecords` API for batch modifying DNS record properties
- Support Delete: `DeleteDnsRecords` API for batch deleting DNS records
- Support Import: via composite ID format `zone_id#record_id`

## Capabilities

### New Capabilities
- `teo-dns-record-24-resource`: Full lifecycle management (CRUD + import) of a single TEO DNS record, including fields for zone_id, name, type, content, location, ttl, weight, priority, and computed fields record_id, status, created_on, modified_on.

### Modified Capabilities
<!-- None - this is a new resource with no impact on existing specifications -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.go` - Resource implementation (already created)
- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.md` - Documentation template (already created)
- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_24_test.go` - Unit tests (already created)
- **Modified file**: `tencentcloud/provider.go` - Resource registration (already registered at line 2090)
- **Dependencies**: Uses existing `tencentcloud-sdk-go/teo/v20220901` package and existing `DescribeTeoDnsRecordById` service method
