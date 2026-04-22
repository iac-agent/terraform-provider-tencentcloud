## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.go` with schema definition and `ResourceTencentCloudTeoDnsRecordV5()` function, following `tencentcloud_igtm_strategy` pattern
- [x] 1.2 Implement `resourceTencentCloudTeoDnsRecordV5Create` - call CreateDnsRecord API, set composite ID (zoneId + FILED_SP + recordId)
- [x] 1.3 Implement `resourceTencentCloudTeoDnsRecordV5Read` - call DescribeTeoDnsRecordById service method, populate all schema fields from response
- [x] 1.4 Implement `resourceTencentCloudTeoDnsRecordV5Update` - handle mutable fields via ModifyDnsRecords API and status field via ModifyDnsRecordsStatus API
- [x] 1.5 Implement `resourceTencentCloudTeoDnsRecordV5Delete` - call DeleteDnsRecords API with zone_id and record_ids

## 2. Service Layer

- [x] 2.1 Verify `DescribeTeoDnsRecordById` method in `service_tencentcloud_teo.go` is reusable for V5 resource (no changes needed if existing method works)

## 3. Provider Registration

- [x] 3.1 Add `"tencentcloud_teo_dns_record_v5": teo.ResourceTencentCloudTeoDnsRecordV5()` to `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_teo_dns_record_v5` to `tencentcloud/provider.md` under TEO Resource section

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v5.md` with resource description, example usage, and import section

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v5_test.go` with gomonkey mock-based unit tests for Create, Read, Update, Delete functions
- [x] 5.2 Run `go test` on the test file to verify all tests pass
