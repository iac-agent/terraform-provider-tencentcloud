## 1. Service Layer

- [x] 1.1 Add `DescribeTeoDnsRecord47ById` method in `tencentcloud/services/teo/service_tencentcloud_teo.go` using `DescribeDnsRecords` API with `AdvancedFilter` by id, with `tccommon.ReadRetryTimeout` retry logic

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_47.go` with schema definition following the existing dns_record_24 pattern, including all fields: zone_id (Required, ForceNew), name (Required), type (Required), content (Required), location (Optional, Computed), ttl (Optional, Computed), weight (Optional, Computed), priority (Optional, Computed), record_id (Computed), status (Computed), created_on (Computed), modified_on (Computed), dns_records (Computed, List)
- [x] 2.2 Implement Create function: build `CreateDnsRecordRequest`, call `CreateDnsRecord` API with `tccommon.WriteRetryTimeout` retry, validate response is not nil and RecordId is not nil, set composite ID `zone_id#record_id`, call Read
- [x] 2.3 Implement Read function: parse composite ID, call `DescribeTeoDnsRecord47ById` service method, set all fields from response (nil-check each field), handle not-found by logging and `d.SetId("")`
- [x] 2.4 Implement Update function: parse composite ID, check for changes in mutable fields (name, type, content, location, ttl, weight, priority), build `ModifyDnsRecordsRequest` with DnsRecord containing RecordId and changed fields, call `ModifyDnsRecords` API with `tccommon.WriteRetryTimeout` retry, call Read
- [x] 2.5 Implement Delete function: parse composite ID, build `DeleteDnsRecordsRequest` with ZoneId and RecordIds, call `DeleteDnsRecords` API with `tccommon.WriteRetryTimeout` retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_dns_record_47` in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_teo_dns_record_47` to `tencentcloud/provider.md` resource list in the TEO section

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_47_test.go` with mock (gomonkey) based tests covering create, read, update, delete, and import scenarios

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_47.md` with usage example, resource description, and import section showing composite ID format