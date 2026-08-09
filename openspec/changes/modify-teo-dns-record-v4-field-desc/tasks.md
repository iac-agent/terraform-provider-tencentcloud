## 1. Update Resource Schema Field Descriptions

- [x] 1.1 Update `zone_id` field Description in `resource_tc_teo_dns_record.go` to "Site ID."
- [x] 1.2 Update `name` field Description to add punycode conversion note for CJK domain names
- [x] 1.3 Update `type` field Description to add reference link for detailed record type format examples
- [x] 1.4 Update `content` field Description to add punycode conversion note for CJK domain names and clarify content corresponds to Type value
- [x] 1.5 Update `location` field Description to add standard/enterprise edition package restriction note and reference link
- [x] 1.6 Update `ttl` field Description to specify default value is 300 and range is 60-86400
- [x] 1.7 Update `weight` field Description to add weight consistency rule note and default value -1
- [x] 1.8 Update `priority` field Description to note MX-only effect, default value 0, and range 0-50
- [x] 1.9 Update `status` field Description to note it is output-only for ModifyDnsRecords and managed through ModifyDnsRecordsStatus
- [x] 1.10 Update `created_on` field Description to note it is output-only for ModifyDnsRecords
- [x] 1.11 Update `modified_on` field Description to note it is output-only for ModifyDnsRecords

## 2. Update Documentation

- [x] 2.1 Update `tencentcloud/services/teo/resource_tc_teo_dns_record.md` to reflect the same description changes applied to the resource schema fields

## 3. Verification

- [x] 3.1 Verify the updated Go file compiles without errors
- [x] 3.2 Verify the updated descriptions match the vendor SDK API documentation
