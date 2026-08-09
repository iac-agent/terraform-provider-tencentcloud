## MODIFIED Requirements

### Requirement: Field descriptions for tencentcloud_teo_dns_record resource
The `tencentcloud_teo_dns_record` resource schema field descriptions SHALL be updated to align with the latest TEO cloud API (v20220901) documentation. The following field descriptions MUST be updated:

1. `zone_id`: Description MUST state "Site ID." (clarifying it represents the TEO site/zone identifier).

2. `name`: Description MUST include the note that if the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.

3. `type`: Description MUST include a reference link for detailed record type format examples for SRV, CAA, and other special record types, pointing to the TEO documentation.

4. `content`: Description MUST include the note that if the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input, and that the content should correspond to the Type value.

5. `location`: Description MUST include the note that resolution route configuration is only applicable to standard edition and enterprise edition packages, with a reference link to the resolution route enumeration documentation.

6. `ttl`: Description MUST specify the default value is 300 seconds, with a range of 60-86400.

7. `weight`: Description MUST include the note that for the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set. Must also specify default value is -1 (no weight set).

8. `priority`: Description MUST note that this parameter only takes effect when Type is MX, specify the default value is 0, and the valid range is 0-50.

9. `status`: Description MUST note that Status is output-only for ModifyDnsRecords API and is managed through the separate ModifyDnsRecordsStatus API.

10. `created_on`: Description MUST note that CreatedOn is output-only for ModifyDnsRecords API.

11. `modified_on`: Description MUST note that ModifiedOn is output-only for ModifyDnsRecords API.

#### Scenario: Updated zone_id description
- **WHEN** a user reads the schema definition of the `zone_id` field
- **THEN** the description clearly states "Site ID."

#### Scenario: Updated name description with punycode note
- **WHEN** a user reads the schema definition of the `name` field
- **THEN** the description includes the punycode conversion requirement for CJK domain names

#### Scenario: Updated type description with reference link
- **WHEN** a user reads the schema definition of the `type` field
- **THEN** the description includes a reference link for detailed record type format examples

#### Scenario: Updated content description with punycode note
- **WHEN** a user reads the schema definition of the `content` field
- **THEN** the description includes the punycode conversion requirement for CJK domain names

#### Scenario: Updated location description with package tier note
- **WHEN** a user reads the schema definition of the `location` field
- **THEN** the description includes the note about standard/enterprise edition package restriction and the reference link to resolution route enumeration

#### Scenario: Updated ttl description with default value
- **WHEN** a user reads the schema definition of the `ttl` field
- **THEN** the description specifies the default value is 300 seconds and the range is 60-86400

#### Scenario: Updated weight description with consistency rule
- **WHEN** a user reads the schema definition of the `weight` field
- **THEN** the description includes the note about weight consistency rule for same subdomain records and specifies default value is -1

#### Scenario: Updated priority description with MX-only note
- **WHEN** a user reads the schema definition of the `priority` field
- **THEN** the description notes it only takes effect when Type is MX, specifies default value is 0, and range is 0-50

#### Scenario: Updated status description with output-only note
- **WHEN** a user reads the schema definition of the `status` field
- **THEN** the description notes that Status is output-only for ModifyDnsRecords and managed through ModifyDnsRecordsStatus

#### Scenario: Updated created_on description with output-only note
- **WHEN** a user reads the schema definition of the `created_on` field
- **THEN** the description notes that CreatedOn is output-only for ModifyDnsRecords

#### Scenario: Updated modified_on description with output-only note
- **WHEN** a user reads the schema definition of the `modified_on` field
- **THEN** the description notes that ModifiedOn is output-only for ModifyDnsRecords

### Requirement: Documentation file for tencentcloud_teo_dns_record
The `tencentcloud/services/teo/resource_tc_teo_dns_record.md` documentation file SHALL be updated to reflect the same description changes applied to the resource schema fields.

#### Scenario: Documentation reflects updated descriptions
- **WHEN** a user reads the resource documentation markdown file
- **THEN** the field descriptions in the documentation match the updated schema field descriptions
