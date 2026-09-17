## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_dns_record_50.go` with `ResourceTencentCloudTeoDnsRecord50()` schema definition and CRUD functions following the `tencentcloud_igtm_strategy` / `tencentcloud_teo_multi_path_gateway` pattern (client 直调 via `UseTeoV20220901Client()`, no service-layer method), including: `zone_id` (Required, ForceNew), `name`/`type`/`content` (Required), `location`/`ttl`/`weight`/`priority` (Optional+Computed), `filters` (Optional, list of {name Required, values Required, fuzzy Optional}), `sort_by`/`sort_order`/`match` (Optional), and computed `record_id`/`status`/`created_on`/`modified_on`
- [x] 1.2 Implement Create (`CreateDnsRecordWithContext`, WriteRetryTimeout, nil/empty `RecordId` → NonRetryableError with logId+request body, then `d.SetId(zoneId#recordId)` and call Read)
- [x] 1.3 Implement Read (`DescribeDnsRecords` with precise `id` filter + user-configured filters/sort_by/sort_order/match, `Limit=1000`, empty result → log `[CRUD] read teo dns_record_50 id=%s` then `d.SetId("")`, nil-check each field before `d.Set`)
- [x] 1.4 Implement Update (`ModifyDnsRecordsWithContext` only when entity fields `name`/`type`/`content`/`location`/`ttl`/`weight`/`priority` have changes; single-element `DnsRecords` with `RecordId`, without `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` inside the element; then call Read)
- [x] 1.5 Implement Delete (`DeleteDnsRecordsWithContext` with `ZoneId` and single-element `RecordIds`, WriteRetryTimeout, no polling)
- [x] 1.6 Verify Create/Update/Delete request parameters only use fields that exist in the corresponding cloud API request structs (per design.md vendor struct excerpts), and that all error returns are checked

## 2. Provider Registration

- [x] 2.1 Register `"tencentcloud_teo_dns_record_50": teo.ResourceTencentCloudTeoDnsRecord50(),` in `tencentcloud/provider.go` ResourcesMap (adjacent to `tencentcloud_teo_dns_record`, keep gofmt alignment)
- [x] 2.2 Add `tencentcloud_teo_dns_record_50` to the TEO Resource list in `tencentcloud/provider.md` (after `tencentcloud_teo_dns_record`)

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_50.md` with one-line description ("Provides a resource to create a teo dns_record_50"), Example Usage (hcl covering required + optional fields, and a `filters` example), and Import section documenting the composite ID `terraform import tencentcloud_teo_dns_record_50.xxx {zoneId}#{recordId}`; do NOT add manual Argument Reference / Attribute Reference sections

## 4. Testing

- [x] 4.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_dns_record_50_test.go` (package `teo_test`, unique mock helper names e.g. `mockMetaForDnsRecord50`, `ptrStrDR50`, `ptrInt64DR50` to avoid same-package collisions) using gomonkey to mock `CreateDnsRecordWithContext`, `DescribeDnsRecords`, `ModifyDnsRecordsWithContext`, `DeleteDnsRecordsWithContext`
- [x] 4.2 Cover: create success (assert composite ID and attribute backfill), create empty record id error, read success (assert `id` filter and `Limit=1000` in describe request), read not-found (SetId ""), update success (assert single-element DnsRecords with RecordId and no ignored fields set), delete success (assert RecordIds), and API error paths
- [x] 4.3 Ensure the generated test code compiles in the current environment (mock signatures match vendor client methods; resource CRUD functions reference only defined symbols); do not execute `go test` in this workflow
