## 1. Service Layer

- [x] 1.1 Add `DescribeTeoDnsRecord20ById` method to `tencentcloud/services/teo/service_tencentcloud_teo.go` that queries DNS record by zone_id and record_id using DescribeDnsRecords API with Filter(id=recordId)

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_20.go` with schema definition including: zone_id (Required, ForceNew), name (Required), type (Required), content (Required), location (Optional), ttl (Optional), weight (Optional), priority (Optional), status (Optional, Computed), record_id (Computed), created_on (Computed), modified_on (Computed)
- [x] 2.2 Implement `resourceTencentCloudTeoDnsRecord20Create` function that calls CreateDnsRecord API, sets composite ID (zone_id#record_id), and validates RecordId is not nil
- [x] 2.3 Implement `resourceTencentCloudTeoDnsRecord20Read` function that parses composite ID, calls service layer DescribeTeoDnsRecord20ById, and sets all fields from response (with nil checks)
- [x] 2.4 Implement `resourceTencentCloudTeoDnsRecord20Update` function with two parts: (1) ModifyDnsRecords for name/type/content/location/ttl/weight/priority changes, (2) ModifyDnsRecordsStatus for status changes
- [x] 2.5 Implement `resourceTencentCloudTeoDnsRecord20Delete` function that calls DeleteDnsRecords API with zone_id and record_id
- [x] 2.6 Add Import support using schema.ImportStatePassthrough

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_dns_record_20` resource in `tencentcloud/provider.go`
- [x] 3.2 Add `tencentcloud_teo_dns_record_20` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_20.md` with description, Example Usage, and Import sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_20_test.go` with gomonkey mock-based unit tests for CRUD operations

## 6. Verification

- [x] 6.1 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass
- [x] 6.2 Verify code correctness: ensure Create API parameters match CreateDnsRecord, Update API parameters match ModifyDnsRecords, Read uses DescribeDnsRecords, Delete uses DeleteDnsRecords
