## Context

The `tencentcloud_teo_dns_record` resource manages DNS records for TencentCloud EdgeOne (TEO) zones. The current resource implementation has all the correct schema fields and CRUD logic, but the field descriptions in the schema and documentation are outdated compared to the latest vendor SDK API documentation (teo v20220901).

The vendor SDK provides detailed Chinese-language comments for each field in the `CreateDnsRecordRequest`, `DnsRecord`, and related structs. These comments include important usage notes such as punycode requirements, package tier restrictions, default values, and consistency rules that are missing or incomplete in the current Terraform resource descriptions.

Current resource file: `tencentcloud/services/teo/resource_tc_teo_dns_record.go`
Current doc file: `tencentcloud/services/teo/resource_tc_teo_dns_record.md`

## Goals / Non-Goals

**Goals:**
- Update all field Description strings in the resource schema to match the latest cloud API documentation
- Ensure descriptions accurately reflect API constraints (default values, valid ranges, package restrictions)
- Update the `.md` documentation file to be consistent with the schema descriptions
- Maintain full backward compatibility — no schema structure or logic changes

**Non-Goals:**
- No changes to resource CRUD logic, schema types, or field attributes (Required/Optional/Computed/ForceNew)
- No new fields or removed fields
- No changes to API calls or retry logic
- No changes to test files (descriptions are cosmetic, not functional)

## Decisions

1. **Use English translations of vendor API comments**: The vendor SDK comments are in Chinese. We translate them to English to match the existing Terraform resource description style. This is consistent with how other resources in this provider handle API documentation.

2. **Preserve existing description style**: Some current descriptions already have good English translations with markdown-style formatting (e.g., bullet lists for `type` field). We update these in-place rather than rewriting from scratch, adding only the missing information.

3. **Include cloud API reference links**: Where the vendor API comments reference documentation URLs (e.g., resolution route enumeration, record type introduction), we include the international version of these URLs in the descriptions, consistent with the existing `type` field description pattern.

4. **Single capability scope**: All description changes fall under one capability since they all pertain to the same resource and the same type of change (description updates).

## Risks / Trade-offs

- **[Risk] Description drift over time**: Cloud API documentation may continue to evolve, causing descriptions to become outdated again → Mitigation: This is an inherent limitation; periodic review is the only mitigation.
- **[Risk] Translation accuracy**: Chinese-to-English translation may lose nuance → Mitigation: Cross-reference with the vendor SDK comments and keep descriptions factual and concise.
- **[No functional risk]**: Since only Description strings are modified, there is zero risk of breaking existing Terraform configurations or state files.
