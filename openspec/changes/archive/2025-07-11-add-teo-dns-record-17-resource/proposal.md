## Why

TEO (TencentCloud EdgeOne) currently lacks Terraform support for managing DNS records. Users need to manually create and manage DNS records through the console or API, which prevents infrastructure-as-code management of DNS configurations within TEO zones. Adding the `tencentcloud_teo_dns_record_17` resource enables users to fully manage TEO DNS records lifecycle (create, read, update, delete) through Terraform.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_dns_record_17` of type RESOURCE_KIND_GENERAL
- Support CRUD operations via cloud API interfaces: CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords
- The resource composite ID uses `zone_id` + `record_id` (separated by `tccommon.FILED_SP`)
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add unit tests using gomonkey mock approach
- Add resource documentation markdown file

## Capabilities

### New Capabilities
- `teo-dns-record`: Terraform resource for managing TEO DNS records, supporting full lifecycle (create, read, update, delete) with fields: zone_id, name, type, content, location, ttl, weight, priority, and computed fields record_id, status, created_on, modified_on

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **New files**: `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go`, corresponding test file, and markdown doc
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation entry)
- **Dependencies**: Existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK package (already vendored)
- **APIs used**: CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords (all synchronous)
