## 1. Schema and Resource Definition

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17.go` with schema definition including: zone_id (Required, ForceNew, TypeString), name (Required, TypeString), type (Required, TypeString), content (Required, TypeString), location (Optional+Computed, TypeString), ttl (Optional+Computed, TypeInt), weight (Optional+Computed, TypeInt), priority (Optional+Computed, TypeInt), status (Optional+Computed, TypeString), created_on (Computed, TypeString), modified_on (Computed, TypeString); add Importer with ImportStatePassthrough
- [x] 1.2 Implement `resourceTencentCloudTeoDnsRecord17Create` function: read schema fields, build CreateDnsRecordRequest, call CreateDnsRecordWithContext inside resource.Retry(WriteRetryTimeout), validate response and RecordId not empty, set composite ID (zone_id#record_id), call Read to refresh state
- [x] 1.3 Implement `resourceTencentCloudTeoDnsRecord17Read` function: parse composite ID, call TeoService.DescribeTeoDnsRecordById, set all response fields with nil checks, handle resource-not-found by setting d.SetId("")
- [x] 1.4 Implement `resourceTencentCloudTeoDnsRecord17Update` function: handle mutable args (name, type, content, location, ttl, weight, priority) via ModifyDnsRecords API, handle status changes via ModifyDnsRecordsStatus API, both with WriteRetryTimeout retry logic
- [x] 1.5 Implement `resourceTencentCloudTeoDnsRecord17Delete` function: parse composite ID, build DeleteDnsRecordsRequest with ZoneId and RecordIds, call DeleteDnsRecordsWithContext inside resource.Retry(WriteRetryTimeout)

## 2. Provider Registration

- [x] 2.1 Add `"tencentcloud_teo_dns_record_17": teo.ResourceTencentCloudTeoDnsRecord17()` entry to ResourcesMap in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Resource Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md` with one-line description mentioning TEO, Example Usage section with HCL configuration, and Import section using composite ID format (zone_id#record_id)

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_17_test.go` with gomonkey mock-based unit tests covering: Create (mock CreateDnsRecord), Read (mock DescribeDnsRecords), Update content fields (mock ModifyDnsRecords), Update status (mock ModifyDnsRecordsStatus), Delete (mock DeleteDnsRecords)
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` and verify all tests pass
