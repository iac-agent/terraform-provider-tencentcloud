## 1. Service Layer

- [x] 1.1 Add `DescribeTeoDnsRecordV3ById` method to `TeoService` in `tencentcloud/services/teo/service_tencentcloud_teo.go` — uses `DescribeDnsRecords` API with `AdvancedFilter` on `id` to query a single DNS record by `zone_id` and `record_id`

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.go` with schema definition including: `zone_id` (Required, ForceNew), `name` (Required), `type` (Required), `content` (Required), `location` (Optional, Computed), `ttl` (Optional, Computed), `weight` (Optional, Computed), `priority` (Optional, Computed), `record_id` (Computed), `status` (Computed), `created_on` (Computed), `modified_on` (Computed)
- [x] 2.2 Implement `resourceTencentCloudTeoDnsRecordV3Create` — call `CreateDnsRecord` API, extract `RecordId`, set composite ID `zone_id#record_id`, delegate to Read
- [x] 2.3 Implement `resourceTencentCloudTeoDnsRecordV3Read` — parse composite ID, call `DescribeTeoDnsRecordV3ById`, flatten response to state, handle not-found
- [x] 2.4 Implement `resourceTencentCloudTeoDnsRecordV3Update` — detect changes on mutable args (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`), call `ModifyDnsRecords` API with `DnsRecord` struct, delegate to Read
- [x] 2.5 Implement `resourceTencentCloudTeoDnsRecordV3Delete` — parse composite ID, call `DeleteDnsRecords` API with `ZoneId` and `RecordIds`

## 3. Provider Registration

- [x] 3.1 Add `"tencentcloud_teo_dns_record_v3": teo.ResourceTencentCloudTeoDnsRecordV3()` to `ResourcesMap` in `tencentcloud/provider.go`
- [x] 3.2 Add resource entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v3_test.go` with gomonkey mock tests for Create, Read, Update, Delete operations

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v3.md` with description, example usage, and import section

## 6. Verification

- [x] 6.1 Run `go test` on the unit test file to verify all tests pass
