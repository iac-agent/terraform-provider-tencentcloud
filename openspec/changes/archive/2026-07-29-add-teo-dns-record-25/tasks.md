## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_teo_dns_record_25.go` with schema definition including `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`, `record_id`
- [x] 1.2 Implement Create function: call `CreateDnsRecord`, validate response, set composite ID `zone_id#record_id`, call Read
- [x] 1.3 Implement Read function: parse composite ID, call `DescribeTeoDnsRecordById`, populate all fields, handle not-found case
- [x] 1.4 Implement Update function: detect changes on mutable fields, call `ModifyDnsRecords` with single `DnsRecord`, call Read
- [x] 1.5 Implement Delete function: parse composite ID, call `DeleteDnsRecords` with single record ID

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_dns_record_25` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Register `tencentcloud_teo_dns_record_25` in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `resource_tc_teo_dns_record_25.md` with example usage and import instructions

## 4. Tests

- [x] 4.1 Create `resource_tc_teo_dns_record_25_test.go` with unit tests using gomonkey for API mocking

## 5. Verification

- [x] 5.1 Verify the code compiles without errors
- [x] 5.2 Run `make doc` to generate website documentation