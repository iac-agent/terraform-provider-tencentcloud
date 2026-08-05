## Why

The TEO (Tencent EdgeOne) service already supports DNS record management via the cloud API (`CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords`). Adding a Terraform resource `tencentcloud_teo_dns_record_47` allows users to manage TEO DNS records declaratively through Terraform, enabling IaC workflows for DNS record provisioning within EdgeOne zones.

## What Changes

- Add a new `tencentcloud_teo_dns_record_47` resource of type RESOURCE_KIND_GENERAL to manage TEO DNS records with full CRUD lifecycle
- Implement Create using `CreateDnsRecord` API, Read using `DescribeDnsRecords` API (filtered by RecordId), Update using `ModifyDnsRecords` API, and Delete using `DeleteDnsRecords` API
- Register the resource in `tencentcloud/provider.go`
- Add corresponding documentation in `.md` format
- Add unit tests using mock (gomonkey) for the resource

## Capabilities

### New Capabilities
- `teo-dns-record-47-resource`: Terraform resource for managing TEO DNS records, supporting A, AAAA, MX, CNAME, TXT, NS, CAA, and SRV record types within an EdgeOne zone

### Modified Capabilities
<!-- No existing capabilities are modified. -->
(None)

## Impact

- **Code**: New file `tencentcloud/services/teo/resource_tc_teo_dns_record_47.go`, new test file `tencentcloud/services/teo/resource_tc_teo_dns_record_47_test.go`, modification to `tencentcloud/provider.go`
- **Documentation**: New file `tencentcloud/services/teo/resource_tc_teo_dns_record_47.md`
- **Dependencies**: Uses existing `tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK (already vendored)
- **No breaking changes**: This is a purely additive change