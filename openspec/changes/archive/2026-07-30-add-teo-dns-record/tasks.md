## 1. Resource Implementation Verification

- [x] 1.1 Verify resource code exists at `tencentcloud/services/teo/resource_tc_teo_dns_record.go` with complete CRUD implementation
- [x] 1.2 Verify schema definition matches all required fields (zone_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on)
- [x] 1.3 Verify Create method calls CreateDnsRecord API with retry and nil response check
- [x] 1.4 Verify Read method uses TeoService.DescribeTeoDnsRecordById with proper nil checks
- [x] 1.5 Verify Update method handles both content changes (ModifyDnsRecords) and status changes (ModifyDnsRecordsStatus)
- [x] 1.6 Verify Delete method calls DeleteDnsRecords API with proper composite ID parsing
- [x] 1.7 Verify Import support via schema.ImportStatePassthrough

## 2. Provider Registration

- [x] 2.1 Verify resource registration in `tencentcloud/provider.go` with key `tencentcloud_teo_dns_record`
- [x] 2.2 Verify resource documentation entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create resource documentation file `tencentcloud/services/teo/resource_tc_teo_dns_record.md` with description, example usage, and import section
- [ ] 3.2 Generate website documentation via `make doc` (performed during finalize phase)

## 4. Unit Tests

- [x] 4.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_dns_record_test.go` using gomonkey mocks for CRUD operations
- [x] 4.2 Ensure test coverage for Create, Read, Update (content + status), Delete, and Import scenarios

## 5. Code Quality

- [x] 5.1 Verify code compiles without errors
- [ ] 5.2 Run gofmt on all changed Go files (performed during finalize phase)
- [ ] 5.3 Create changelog entry for the new resource (performed during finalize phase)