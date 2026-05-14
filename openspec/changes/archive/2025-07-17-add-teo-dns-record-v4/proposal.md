## Why

Terraform Provider for TencentCloud currently has `tencentcloud_teo_dns_record` resource for managing TEO DNS records, but it lacks support for the `status` field as an updatable parameter via `ModifyDnsRecordsStatus` API and doesn't fully align with the latest API capabilities. A new `tencentcloud_teo_dns_record_v4` resource is needed to provide a clean, comprehensive resource definition that fully covers the TEO DNS Record CRUD lifecycle using the latest cloud API interfaces (`CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords`), following the v4 resource naming convention for improved API coverage.

## What Changes

- Add new Terraform resource `tencentcloud_teo_dns_record_v4` under `tencentcloud/services/teo/`
- Implement full CRUD operations:
  - **Create**: Call `CreateDnsRecord` API with zone_id, name, type, content, location, ttl, weight, priority parameters; composite ID format: `zone_id#record_id`
  - **Read**: Call `DescribeDnsRecords` API, filter by record ID to read single record state; read computed fields: record_id, status, created_on, modified_on
  - **Update**: Call `ModifyDnsRecords` API for mutable fields (name, type, content, location, ttl, weight, priority); call `ModifyDnsRecordsStatus` API for status field changes
  - **Delete**: Call `DeleteDnsRecords` API with zone_id and record_ids
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add resource example documentation in `tencentcloud/services/teo/resource_tc_teo_dns_record_v4.md`
- Add unit tests with gomonkey mocks in `resource_tc_teo_dns_record_v4_test.go`

## Capabilities

### New Capabilities
- `teo-dns-record-v4`: New TEO DNS record resource supporting full CRUD lifecycle via CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, ModifyDnsRecordsStatus, and DeleteDnsRecords APIs

### Modified Capabilities

## Impact

- **New files**: `tencentcloud/services/teo/resource_tc_teo_dns_record_v4.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v4_test.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v4.md`
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation entry)
- **Cloud APIs**: `teo/v20220901.CreateDnsRecord`, `teo/v20220901.DescribeDnsRecords`, `teo/v20220901.ModifyDnsRecords`, `teo/v20220901.ModifyDnsRecordsStatus`, `teo/v20220901.DeleteDnsRecords`
- **Dependencies**: Existing `tencentcloud-sdk-go` vendor packages, `tccommon` helper packages
- **Backward compatibility**: No impact on existing `tencentcloud_teo_dns_record` resource; this is a new independent resource
