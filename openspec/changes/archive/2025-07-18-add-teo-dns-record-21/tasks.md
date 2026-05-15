## 1. Resource Schema and CRUD Implementation

- [x] 1.1 Create `resource_tc_teo_dns_record_21.go` with resource definition function `ResourceTencentCloudTeoDnsRecord21()` returning `*schema.Resource` with Create, Read, Update, Delete, Importer
- [x] 1.2 Define resource schema with required fields (zone_id ForceNew, name, type, content), optional fields (location, ttl, weight, priority), and computed fields (record_id, status, created_on, modified_on)
- [x] 1.3 Implement `resourceTencentCloudTeoDnsRecord21Create()` that calls `CreateDnsRecord` API, sets composite ID (`zone_id#record_id`), with WriteRetryTimeout retry logic
- [x] 1.4 Implement `resourceTencentCloudTeoDnsRecord21Read()` that parses composite ID, calls service layer `DescribeTeoDnsRecordById()`, and sets all state fields, with ReadRetryTimeout retry logic
- [x] 1.5 Implement `resourceTencentCloudTeoDnsRecord21Update()` that calls `ModifyDnsRecords` API with changed fields detected via `d.HasChange()`, with WriteRetryTimeout retry logic
- [x] 1.6 Implement `resourceTencentCloudTeoDnsRecord21Delete()` that calls `DeleteDnsRecords` API with zone_id and record_id, with WriteRetryTimeout retry logic
- [x] 1.7 Add standard error handling with `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_dns_record_21` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add resource declaration in `tencentcloud/provider.md` under TEO resources section

## 3. Test Implementation

- [x] 3.1 Create `resource_tc_teo_dns_record_21_test.go` with unit tests using mock (gomonkey) approach
- [x] 3.2 Implement test covering create operation with mocked `CreateDnsRecord` API response
- [x] 3.3 Implement test covering read operation with mocked `DescribeDnsRecords` API response
- [x] 3.4 Implement test covering update operation with mocked `ModifyDnsRecords` API response
- [x] 3.5 Implement test covering delete operation with mocked `DeleteDnsRecords` API response
- [x] 3.6 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 4. Documentation

- [x] 4.1 Create `resource_tc_teo_dns_record_21.md` with one-line description, Example Usage section, and Import section (with composite ID format note)
