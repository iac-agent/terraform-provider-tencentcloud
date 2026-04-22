## 1. Service Layer Implementation

- [x] 1.1 Add `DescribeTeoDnsRecordV6ById` method in `tencentcloud/services/teo/service_tencentcloud_teo.go` that calls `DescribeDnsRecords` API with AdvancedFilter on record ID, paginates with Limit=1000, and returns the matching DnsRecord or nil if not found
- [x] 1.2 Add `DeleteTeoDnsRecordV6ById` method in `tencentcloud/services/teo/service_tencentcloud_teo.go` for DNS record deletion by zone_id and record_id

## 2. Resource Schema and CRUD Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v6.go` with `ResourceTencentCloudTeoDnsRecordV6()` function defining the full schema (zone_id, name, type, content, location, ttl, weight, priority, status, record_id, created_on, modified_on) and Importer
- [x] 2.2 Implement `resourceTencentCloudTeoDnsRecordV6Create` function calling `CreateDnsRecord` API, setting composite ID (zone_id + FILED_SP + record_id)
- [x] 2.3 Implement `resourceTencentCloudTeoDnsRecordV6Read` function calling `DescribeTeoDnsRecordV6ById` service method and setting all schema fields from response
- [x] 2.4 Implement `resourceTencentCloudTeoDnsRecordV6Update` function handling mutable field updates via `ModifyDnsRecords` and status updates via `ModifyDnsRecordsStatus`
- [x] 2.5 Implement `resourceTencentCloudTeoDnsRecordV6Delete` function calling `DeleteDnsRecords` API

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_dns_record_v6` resource in `tencentcloud/provider.go` with key `tencentcloud_teo_dns_record_v6` mapping to `teo.ResourceTencentCloudTeoDnsRecordV6()`
- [x] 3.2 Add resource entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v6_test.go` with gomonkey mock-based unit tests for create, read, update, and delete operations
- [x] 4.2 Run `go test` on the test file to verify all tests pass

## 5. Resource Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v6.md` with description, example usage (including jsonencode() where applicable), and import section

## 6. Finalization

- [ ] 6.1 Run `gofmt` on all new and modified Go files
- [ ] 6.2 Run `make doc` to generate website documentation
- [ ] 6.3 Create `.changelog` entry file
