## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v7.go` with schema definition and CRUD functions (Create, Read, Update, Delete), referencing `resource_tc_igtm_strategy.go` code style
- [x] 1.2 Implement Create function: call `CreateDnsRecord` API with retry, validate RecordId is not nil, set composite ID `{zoneId}{FILED_SP}{recordId}`
- [x] 1.3 Implement Read function: use existing `DescribeTeoDnsRecordById` service method, set all schema fields with nil checks, handle not-found case
- [x] 1.4 Implement Update function: split into two parts — content changes via `ModifyDnsRecords`, status changes via `ModifyDnsRecordsStatus`
- [x] 1.5 Implement Delete function: call `DeleteDnsRecords` API with ZoneId and RecordIds
- [x] 1.6 Add Import support via `schema.ResourceImporter{State: schema.ImportStatePassthrough}`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_dns_record_v7` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Update `tencentcloud/provider.md` to add the new resource entry

## 3. Unit Tests

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v7_test.go` with unit tests using gomonkey mock approach
- [x] 3.2 Add test cases for Create, Read, Update (content + status), Delete functions
- [x] 3.3 Run unit tests with `go test -gcflags=all=-l` to verify

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_v7.md` with description, example usage, and import section
