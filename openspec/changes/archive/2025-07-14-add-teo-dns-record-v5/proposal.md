## Why

There is an existing `tencentcloud_teo_dns_record` resource, but it needs a V5 version (`tencentcloud_teo_dns_record_v5`) to follow the latest resource code generation patterns and provide improved implementation quality. The new V5 resource will use the same cloud API interfaces (CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords, ModifyDnsRecordsStatus) but with updated code structure referencing `tencentcloud_igtm_strategy` as the pattern template.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_dns_record_v5` of type RESOURCE_KIND_GENERAL for the TEO (EdgeOne) cloud product
- The resource manages DNS records with full CRUD lifecycle:
  - **Create**: Call `CreateDnsRecord` API with zone_id, name, type, content, location, ttl, weight, priority
  - **Read**: Call `DescribeDnsRecords` API with zone_id and filter by record_id
  - **Update**: Call `ModifyDnsRecords` API for mutable fields (name, type, content, location, ttl, weight, priority) and `ModifyDnsRecordsStatus` API for status field
  - **Delete**: Call `DeleteDnsRecords` API with zone_id and record_ids
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add service layer method in `service_tencentcloud_teo.go` for reading DNS record by ID
- Create unit tests using mock (gomonkey) approach
- Create `.md` documentation file for the resource

## Capabilities

### New Capabilities
- `teo-dns-record-v5`: New TEO DNS record V5 resource supporting full CRUD operations for managing EdgeOne DNS records

### Modified Capabilities

## Impact

- **New files**: `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v5_test.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.md`
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation index), `tencentcloud/services/teo/service_tencentcloud_teo.go` (service layer method)
- **Cloud APIs**: CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, ModifyDnsRecordsStatus, DeleteDnsRecords from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- **Dependencies**: No new external dependencies required
