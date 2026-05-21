## Why

TEO (TencentCloud EdgeOne) DNS record management is not yet available as a Terraform resource, preventing users from managing DNS records for their EdgeOne zones via infrastructure-as-code. Adding this resource enables full lifecycle management of TEO DNS records through Terraform.

## What Changes

- Add new Terraform resource `tencentcloud_teo_dns_record_24` (RESOURCE_KIND_GENERAL) for managing TEO DNS records
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add resource implementation file `tencentcloud/services/teo/resource_tc_teo_dns_record_24.go`
- Add unit test file `tencentcloud/services/teo/resource_tc_teo_dns_record_24_test.go`
- Add documentation file `tencentcloud/services/teo/resource_tc_teo_dns_record_24.md`

## Capabilities

### New Capabilities

- `teo-dns-record`: Manage TEO DNS records lifecycle (create, read, update, delete) using the TEO v20220901 API. Supports DNS record types A, AAAA, MX, CNAME, TXT, NS, CAA, SRV with configurable TTL, weight, priority, and location.

### Modified Capabilities

## Impact

- New file: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_dns_record_24_test.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_dns_record_24.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (document new resource)
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)
- APIs used: `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords`
