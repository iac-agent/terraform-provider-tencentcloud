## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_24.go` with schema definition including all fields: `zone_id` (Required, ForceNew), `name` (Required, ForceNew), `type` (Required, ForceNew), `content` (Required), `location` (Optional), `ttl` (Optional), `weight` (Optional), `priority` (Optional), `record_id` (Computed), `status` (Computed), `created_on` (Computed), `modified_on` (Computed)
- [x] 1.2 Implement `resourceTencentCloudTeoDnsRecord24Create` function: call `CreateDnsRecord` API with retry, check response for nil/empty `RecordId`, set resource ID as `zone_id#record_id`
- [x] 1.3 Implement `resourceTencentCloudTeoDnsRecord24Read` function: parse composite ID, call `DescribeDnsRecords` with `id` filter, handle not-found by clearing state, set all fields from `DnsRecord` response
- [x] 1.4 Implement `resourceTencentCloudTeoDnsRecord24Update` function: call `ModifyDnsRecords` with single-element `DnsRecords` list containing updated fields and `RecordId`, then call Read to refresh state
- [x] 1.5 Implement `resourceTencentCloudTeoDnsRecord24Delete` function: parse composite ID, call `DeleteDnsRecords` with `RecordIds=[record_id]`
- [x] 1.6 Add import support via `schema.ImportStatePassthrough`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_dns_record_24` in `tencentcloud/provider.go` under the TEO service section (reference `tencentcloud_igtm_strategy` pattern)
- [x] 2.2 Add `tencentcloud_teo_dns_record_24` entry in `tencentcloud/provider.md` under the TEO section

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_24.md` with one-sentence description mentioning TEO, Example Usage section using a realistic DNS record example, and Import section showing composite ID format

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_dns_record_24_test.go` with gomonkey-based unit tests for Create, Read, Update, and Delete operations (mock `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords` API calls)
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l ./tencentcloud/services/teo/... -run TestTeoDnsRecord24` to verify all tests pass
