## Why

TEO (TencentCloud EdgeOne) currently provides a `tencentcloud_teo_dns_record` resource for managing DNS records, but it does not expose all available fields from the latest cloud API (e.g., `priority` for MX records). A new `tencentcloud_teo_dns_record_v3` resource is needed to fully support the updated `CreateDnsRecord`/`ModifyDnsRecords` APIs, which include the `Priority` field and allow more complete DNS record management in TEO.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_dns_record_v3` of type RESOURCE_KIND_GENERAL
- Implement full CRUD operations using the following cloud APIs:
  - **Create**: `CreateDnsRecord` — creates a DNS record with `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`
  - **Read**: `DescribeDnsRecords` — queries DNS records by `zone_id` and filters by `record_id`
  - **Update**: `ModifyDnsRecords` — modifies an existing DNS record
  - **Delete**: `DeleteDnsRecords` — deletes a DNS record by `zone_id` and `record_id`
- Register the new resource in `provider.go`
- Add unit tests using gomonkey mock approach
- Add documentation `.md` file for the resource

## Capabilities

### New Capabilities
- `teo-dns-record-v3`: Full CRUD resource for managing TEO DNS records, supporting all fields including `priority` for MX records, `weight` for weighted routing, and `location` for resolution line configuration.

### Modified Capabilities
- None

## Impact

- **New files**: `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v3_test.go`, `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.md`
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation entry)
- **Service layer**: Add `DescribeTeoDnsRecordV3ById` method to `TeoService` in `service_tencentcloud_teo.go`
- **Cloud APIs**: Uses `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords` from `teo/v20220901` SDK package (all synchronous)
- **Backward compatibility**: Fully backward compatible — this is a new resource with no changes to existing resources
