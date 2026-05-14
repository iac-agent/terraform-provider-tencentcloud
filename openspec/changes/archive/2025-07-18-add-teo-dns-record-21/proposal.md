## Why

TEO (EdgeOne) needs a new version of the DNS record Terraform resource (`tencentcloud_teo_dns_record_21`) to provide improved management of DNS records for EdgeOne zones. The existing `tencentcloud_teo_dns_record` resource serves the same API but this new version is needed to support updated API capabilities and provide a refreshed implementation following the latest provider patterns.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_dns_record_21` to manage TEO DNS records lifecycle (create, read, update, delete)
- Support creating DNS records via `CreateDnsRecord` API with parameters: zone_id, name, type, content, location, ttl, weight, priority
- Support reading DNS records via `DescribeDnsRecords` API, filtering by record ID
- Support updating DNS records via `ModifyDnsRecords` API
- Support deleting DNS records via `DeleteDnsRecords` API
- Add corresponding test file for the new resource
- Add documentation (.md file) for the new resource
- Register the new resource in `provider.go` and `provider.md`

## Capabilities

### New Capabilities
- `teo-dns-record-21-resource`: Terraform resource for managing TEO DNS records lifecycle (create, read, update, delete)

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **New files**:
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_21.go`
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_21_test.go`
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_21.md`
  - `openspec/changes/add-teo-dns-record-21/specs/teo-dns-record-21-resource/spec.md`

- **Modified files**:
  - `tencentcloud/provider.go` - register new resource
  - `tencentcloud/provider.md` - add resource declaration

- **Dependencies**:
  - Uses existing `tencentcloud-sdk-go` TEO package (already vendored at v1.3.93)
  - No new external dependencies needed

- **Breaking changes**: None - this is a new resource addition
