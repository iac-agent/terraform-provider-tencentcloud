## Why

The `tencentcloud_teo_dns_record` resource's field descriptions are outdated and incomplete compared to the latest TEO cloud API documentation. The current descriptions lack important details such as punycode conversion requirements, record type format examples, package tier restrictions, weight consistency rules, and default value specifications that are documented in the vendor SDK. Updating these descriptions will improve user experience by providing accurate and comprehensive field documentation aligned with the cloud API.

## What Changes

- Update `zone_id` field description to clarify it represents the site ID
- Update `name` field description to add punycode conversion requirement for Chinese, Korean, and Japanese domain names
- Update `type` field description to add reference link for detailed record type format examples (SRV, CAA, etc.)
- Update `content` field description to add punycode conversion requirement for Chinese, Korean, and Japanese domain names
- Update `location` field description to add note about standard/enterprise edition package restriction
- Update `ttl` field description to clarify default value is 300 seconds
- Update `weight` field description to add note about weight consistency rule for records under the same subdomain
- Update `priority` field description to note it only takes effect when type is MX, and specify default value is 0
- Update `status` field description to note it is output-only for ModifyDnsRecords API
- Update `created_on` field description to note it is output-only for ModifyDnsRecords API
- Update `modified_on` field description to note it is output-only for ModifyDnsRecords API
- Update corresponding `.md` documentation file with the same description changes

## Capabilities

### New Capabilities

_None_

### Modified Capabilities

- `teo-dns-record-field-desc`: Update field descriptions for tencentcloud_teo_dns_record resource to align with latest cloud API documentation

## Impact

- **Files affected**:
  - `tencentcloud/services/teo/resource_tc_teo_dns_record.go` (schema Description fields)
  - `tencentcloud/services/teo/resource_tc_teo_dns_record.md` (documentation)
- **APIs**: No API call changes, only description text updates
- **Backward compatibility**: Fully compatible — only documentation/description strings are modified, no schema or logic changes
