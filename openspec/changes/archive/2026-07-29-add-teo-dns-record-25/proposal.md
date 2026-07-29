## Why

The TEO DNS record resource (`tencentcloud_teo_dns_record`) currently supports a rich set of parameters including `status`, `created_on`, `modified_on`, and `dns_records` (nested list). However, a new simplified variant `tencentcloud_teo_dns_record_25` is needed to provide a focused DNS record management experience with only the core CRUD parameters (zone_id, name, type, content, location, ttl, weight, priority, record_id) without the status tracking and nested record list fields. This follows the existing pattern of versioned DNS record resources (e.g., `_21`, `_22`, `_23`, `_24`).

## What Changes

- Add new Terraform resource `tencentcloud_teo_dns_record_25` for managing TEO DNS records
- The resource uses `CreateDnsRecord` for creation, `DescribeDnsRecords` for reading, `ModifyDnsRecords` for updating, and `DeleteDnsRecords` for deletion
- The resource ID is composed of `zone_id` and `record_id` joined by the standard separator

## Capabilities

### New Capabilities
- `teo-dns-record-25`: Simplified TEO DNS record management with core CRUD parameters (zone_id, name, type, content, location, ttl, weight, priority, record_id)

### Modified Capabilities
<!-- None -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_25.go` — resource implementation
- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_25_test.go` — unit tests
- **New file**: `tencentcloud/services/teo/resource_tc_teo_dns_record_25.md` — documentation
- **Modified file**: `tencentcloud/provider.go` — register the new resource
- **Modified file**: `tencentcloud/provider.md` — register the new resource in docs
- **Dependency**: `tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)
- **Reuses**: `DescribeTeoDnsRecordById` service function (already in `service_tencentcloud_teo.go`)