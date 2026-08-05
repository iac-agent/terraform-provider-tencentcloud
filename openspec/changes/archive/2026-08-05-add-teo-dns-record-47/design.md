## Context

The TEO (Tencent EdgeOne) service already has several DNS record resource variants (`tencentcloud_teo_dns_record`, `tencentcloud_teo_dns_record_21`, `tencentcloud_teo_dns_record_22`, `tencentcloud_teo_dns_record_23`, `tencentcloud_teo_dns_record_24`) in the Terraform provider. These all follow the same pattern using the `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, and `DeleteDnsRecords` cloud APIs from the `tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK.

The new variant `tencentcloud_teo_dns_record_47` follows the same established pattern as the existing variants, providing a RESOURCE_KIND_GENERAL resource with full CRUD lifecycle management for TEO DNS records.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_dns_record_47` for managing TEO DNS records
- Support all standard DNS record types: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- Implement full CRUD: Create (CreateDnsRecord), Read (DescribeDnsRecords filtered by id), Update (ModifyDnsRecords), Delete (DeleteDnsRecords)
- Support resource import via composite ID (`zone_id#record_id`)
- Add service layer method `DescribeTeoDnsRecord47ById` in service_tencentcloud_teo.go
- Provide unit test coverage using mock (gomonkey) approach

**Non-Goals:**
- No data source variant (this is RESOURCE_KIND_GENERAL only)
- No batch DNS record operations (single record management)
- No status toggle via ModifyDnsRecordsStatus (out of scope, though the API exists)
- No modification of existing DNS record resource variants

## Decisions

1. **Composite ID**: Use `zone_id#record_id` as the composite ID (separated by `tccommon.FILED_SP`), consistent with existing variants. ZoneId is ForceNew and Required since changing the zone would require resource recreation.

2. **Service layer pattern**: Add a new `DescribeTeoDnsRecord47ById` method in `service_tencentcloud_teo.go` using the `DescribeDnsRecords` API with `AdvancedFilter` by id. The Read operation delegates to this service method.

3. **API client usage**: Use `UseTeoClient()` (which returns the v20220901 client) consistent with existing DNS record variants, not `UseTeoV20220901Client()`.

4. **Mutable fields**: All top-level fields except `zone_id` (ForceNew) are mutable: `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`. Computed-only fields (`status`, `created_on`, `modified_on`, `record_id`) are output-only.

5. **Schema design**: Follow the exact schema structure of existing variants, including the `dns_records` computed list field that mirrors the full DnsRecord struct for read-back compatibility.

6. **No `dns_records` wrapper**: Do NOT create a `dns_records` list-type schema as a wrapper for the resource parameters. Each parameter is a standalone top-level schema field. The `dns_records` field is Computed-only, used for reading back the full record detail.

## Risks / Trade-offs

- **Risk**: Cloud API changes could break the resource → **Mitigation**: The SDK is versioned (v20220901), and API changes require SDK updates which are explicitly managed via vendor
- **Risk**: Duplicate variant naming (47) could cause confusion → **Mitigation**: This follows the established naming convention for DNS record resource variants in the provider
- **Trade-off**: Single record management vs batch → Single record per resource instance is simpler for Terraform state management, but less efficient for bulk operations