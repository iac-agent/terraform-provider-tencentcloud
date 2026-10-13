## 1. Service Layer

- [x] 1.1 Add `DescribeTeoDnsRecord17ById(ctx context.Context, zoneId, recordId string, limit int64)` method to `tencentcloud/services/teo/service_tencentcloud_teo.go`: build `DescribeDnsRecordsRequest` with `ZoneId`, `Filters` (`id` AdvancedFilter matching recordId), `Offset=0`, `Limit=limit`; wrap the API call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check()`, take `DnsRecords[0]` as result

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go` with `ResourceTencentCloudTeoDnsRecord17()` schema definition (code style strictly following `tencentcloud_igtm_strategy`): `zone_id` (Required, ForceNew), `name`/`type`/`content` (Required), `location`/`ttl`/`weight`/`priority`/`status` (Optional+Computed), **`Limit` (Optional, TypeInt — the newly added `DescribeDnsRecords` query param)**, `record_id`/`created_on`/`modified_on` (Computed); support import via `schema.ImportStatePassthrough`
- [x] 2.2 Implement `resourceTencentCloudTeoDnsRecord17Create`: build `CreateDnsRecordRequest` from `zone_id`/`name`/`type`/`content`/`location`/`ttl`/`weight`/`priority` (no `Limit`); wrap in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; after retry success, log logId and validate `response.Response` / `RecordId` are non-nil and non-empty (return `NonRetryableError` if empty); set ID to `zoneId#recordId` with `tccommon.FILED_SP`, then call Read
- [x] 2.3 Implement `resourceTencentCloudTeoDnsRecord17Read`: parse composite ID from `d.Id()`; read `Limit` via `d.GetOk("Limit")` defaulting to 1000 (cloud API max); call `DescribeTeoDnsRecord17ById`; when result is nil, log `[CRUD] teo_dns_record_17 id=%s` first, then `d.SetId("")` and return nil; set each field with nil check before `d.Set()`
- [x] 2.4 Implement `resourceTencentCloudTeoDnsRecord17Update`: for changes of `name`/`type`/`content`/`location`/`ttl`/`weight`/`priority`, call `ModifyDnsRecords` with single-element `DnsRecords` array containing `RecordId` + updatable fields only (exclude `ZoneId`/`Status`/`CreatedOn`/`ModifiedOn` output-only fields); for `status` change, call `ModifyDnsRecordsStatus` (enable → `RecordsToEnable`, disable → `RecordsToDisable`); both wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; `Limit` change alone SHALL NOT trigger any write API call; finish with Read
- [x] 2.5 Implement `resourceTencentCloudTeoDnsRecord17Delete`: parse composite ID, build `DeleteDnsRecordsRequest` with `ZoneId` + `RecordIds`, wrap in `resource.Retry(tccommon.WriteRetryTimeout, ...)`

## 3. Provider Registration

- [x] 3.1 Register `"tencentcloud_teo_dns_record_17": teo.ResourceTencentCloudTeoDnsRecord17()` in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_teo_dns_record_17` entry to `tencentcloud/provider.md` under the TEO category

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md` with one-line description, Example Usage (including the `Limit` parameter), and Import section explaining the composite ID `zone_id#record_id`

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17_test.go` using gomonkey mocks (not the Terraform test suite): cover Create/Read/Update/Delete logic, including verification that `DescribeDnsRecords` is called with `Limit=1000` by default and `Limit=<user value>` (e.g. 500) when configured, and that create returns `NonRetryableError` on empty `RecordId`

## 6. Verification

- [x] 6.1 Verify all new/modified Go files compile correctly (gofmt formatting handled by tfpacer-finalize skill; `make doc` website docs generation and `.changelog` entry handled by tfpacer-finalize skill — do NOT run them in this phase)
