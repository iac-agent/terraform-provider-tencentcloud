## 1. Resource Schema and CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v4.go` with schema definition including: zone_id (Required, ForceNew), name (Required), type (Required), content (Required), location (Optional, Computed), ttl (Optional, Computed), weight (Optional, Computed), priority (Optional, Computed), status (Optional, Computed), record_id (Computed), created_on (Computed), modified_on (Computed)
- [x] 1.2 Implement `resourceTencentCloudTeoDnsRecordV4Create` function: call CreateDnsRecord API with retry (tccommon.WriteRetryTimeout), validate RecordId is not nil/empty (return NonRetryableError if empty), set composite ID as `zone_id#record_id`
- [x] 1.3 Implement `resourceTencentCloudTeoDnsRecordV4Read` function: parse composite ID to get zone_id and record_id, call DescribeDnsRecords with AdvancedFilter (Name="id", Values=[record_id]) to find the specific record, set all schema fields with nil-checks before d.Set(), clear resource ID if record not found
- [x] 1.4 Implement `resourceTencentCloudTeoDnsRecordV4Update` function: detect changes in mutable args (name, type, content, location, ttl, weight, priority), call ModifyDnsRecords with ZoneId and DnsRecord struct containing RecordId + updated fields; detect status change and call ModifyDnsRecordsStatus with RecordsToEnable or RecordsToDisable accordingly
- [x] 1.5 Implement `resourceTencentCloudTeoDnsRecordV4Delete` function: parse composite ID, call DeleteDnsRecords API with ZoneId and RecordIds, use retry with tccommon.WriteRetryTimeout

## 2. Service Layer

- [x] 2.1 Add `DescribeTeoDnsRecordV4ById` method to TeoService in `tencentcloud/services/teo/service_tencentcloud_teo.go` (if not existing): call DescribeDnsRecords with AdvancedFilter by record ID, return single DnsRecord or nil

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_dns_record_v4` resource in `tencentcloud/provider.go` with resource name and factory function `teo.ResourceTencentCloudTeoDnsRecordV4()`
- [x] 3.2 Add resource entry in `tencentcloud/provider.md` documentation

## 4. Resource Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v4.md` with: one-line description mentioning TEO product, Example Usage section with HCL example showing required and optional fields, Import section with composite ID format (zone_id#record_id)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v4_test.go` with gomonkey-based unit tests covering: TestResourceTencentCloudTeoDnsRecordV4Create (mock CreateDnsRecord, verify ID format), TestResourceTencentCloudTeoDnsRecordV4Read (mock DescribeDnsRecords, verify field population), TestResourceTencentCloudTeoDnsRecordV4Update (mock ModifyDnsRecords and ModifyDnsRecordsStatus), TestResourceTencentCloudTeoDnsRecordV4Delete (mock DeleteDnsRecords)
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 6. Finalize

- [ ] 6.1 Run `gofmt` on all new .go files to ensure proper formatting
- [ ] 6.2 Run `make doc` to generate website documentation from .md files
- [ ] 6.3 Create `.changelog/` entry file for the new resource
