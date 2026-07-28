## Why

TEO (Tencent EdgeOne) DNS record management is a core capability for EdgeOne users. The existing `tencentcloud_teo_dns_record` resource is being superseded by versioned resources (21, 22, 23, 24) that each provide a subset of DNS record management capabilities. This proposal adds `tencentcloud_teo_dns_record_24` to provide comprehensive DNS record CRUD operations using the latest TEO v20220901 API, completing the versioned DNS record management family.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_dns_record_24` (RESOURCE_KIND_GENERAL) for TEO DNS record management
- Support full CRUD lifecycle: create, read, update, and delete DNS records
- Map the following TEO cloud API interfaces:
  - `CreateDnsRecord` - create a single DNS record
  - `DescribeDnsRecords` - query DNS records by filter (used via service layer helper `DescribeTeoDnsRecordById`)
  - `ModifyDnsRecords` - batch update DNS records
  - `DeleteDnsRecords` - batch delete DNS records
- Support resource import using compound ID (`zone_id#record_id`)
- Expose computed fields: `record_id`, `status`, `created_on`, `modified_on`, `dns_records`

## Capabilities

### New Capabilities
- `teo-dns-record-24-resource`: A Terraform resource that manages the full lifecycle of a TEO DNS record, including creation, read, update, delete, and import. Uses `zone_id` and `record_id` as a compound ID and supports all DNS record types (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV).

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.go` - resource implementation
- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.md` - resource documentation template
- **Modified file**: `tencentcloud/provider.go` - register the new resource in the provider
- **Dependencies**: Uses existing `TeoService.DescribeTeoDnsRecordById` in `service_tencentcloud_teo.go`
- **API dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`