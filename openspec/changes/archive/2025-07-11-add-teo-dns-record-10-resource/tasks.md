## 1. Service Layer Implementation

- [x] 1.1 Add `DescribeTeoDnsRecord10ById` method to `tencentcloud/services/teo/service_tencentcloud_teo.go` that queries a single DNS record using DescribeDnsRecords API with AdvancedFilter(id=recordId), returning *teov20220901.DnsRecord or nil

## 2. Resource Schema and CRUD Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_10.go` with resource schema definition (zone_id, name, type, content as Required; location, ttl, weight, priority as Optional+Computed; record_id, status, created_on, modified_on as Computed)
- [x] 2.2 Implement `resourceTencentCloudTeoDnsRecord10Create` function calling CreateDnsRecord API with retry (WriteRetryTimeout), setting composite ID (zone_id#record_id), and checking for empty RecordId response
- [x] 2.3 Implement `resourceTencentCloudTeoDnsRecord10Read` function calling service layer DescribeTeoDnsRecord10ById, parsing composite ID, setting all fields with nil checks, handling not-found case by clearing resource ID
- [x] 2.4 Implement `resourceTencentCloudTeoDnsRecord10Update` function calling ModifyDnsRecords API with retry (WriteRetryTimeout), building DnsRecord object with RecordId and mutable fields, detecting field changes before calling API
- [x] 2.5 Implement `resourceTencentCloudTeoDnsRecord10Delete` function calling DeleteDnsRecords API with retry (WriteRetryTimeout), parsing composite ID, passing RecordIds=[recordId]

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_dns_record_10` resource in `tencentcloud/provider.go` with resource map entry
- [x] 3.2 Add `tencentcloud_teo_dns_record_10` resource documentation entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_10_test.go` with gomonkey mock-based unit tests for Create, Read, Update, Delete operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 5. Resource Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_10.md` with one-line description ("Provides a resource to ..."), example usage, and import section explaining composite ID format
