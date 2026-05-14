## ADDED Requirements

### Requirement: Resource schema defines DNS record fields
The resource SHALL define a schema with the following fields:
- `zone_id` (Required, ForceNew, TypeString): TEO zone ID
- `name` (Required, TypeString): DNS record name
- `type` (Required, TypeString): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (Required, TypeString): DNS record content
- `location` (Optional, Computed, TypeString): DNS record resolution route
- `ttl` (Optional, Computed, TypeInt): Cache time in seconds (60-86400)
- `weight` (Optional, Computed, TypeInt): DNS record weight (-1 to 100)
- `priority` (Optional, Computed, TypeInt): MX record priority (0-50)
- `status` (Optional, Computed, TypeString): DNS record resolution status (enable/disable)
- `record_id` (Computed, TypeString): DNS record ID
- `created_on` (Computed, TypeString): Creation time
- `modified_on` (Computed, TypeString): Modification time

#### Scenario: Schema accepts all input parameters
- **WHEN** a user defines a tencentcloud_teo_dns_record_v4 resource with zone_id, name, type, and content
- **THEN** the resource SHALL accept these required fields and all optional fields (location, ttl, weight, priority, status)

#### Scenario: Computed fields are populated after read
- **WHEN** the resource is read from the cloud API
- **THEN** record_id, status, created_on, and modified_on SHALL be populated from the API response

### Requirement: Create DNS record via CreateDnsRecord API
The resource SHALL call CreateDnsRecord API with zone_id, name, type, content, location, ttl, weight, and priority parameters. Upon success, the RecordId from the response SHALL be combined with zone_id to form the composite resource ID using `tccommon.FILED_SP` separator (format: `zone_id#record_id`).

#### Scenario: Successful creation with required fields
- **WHEN** CreateDnsRecord is called with zone_id, name, type, and content
- **THEN** the API SHALL be called with these parameters, and the resource ID SHALL be set to `zone_id#record_id`

#### Scenario: Creation with all optional fields
- **WHEN** CreateDnsRecord is called with all optional fields (location, ttl, weight, priority)
- **THEN** all optional fields SHALL be included in the API request

#### Scenario: Create API returns empty RecordId
- **WHEN** CreateDnsRecord response has nil or empty RecordId
- **THEN** a NonRetryableError SHALL be returned

#### Scenario: Create API call failure
- **WHEN** CreateDnsRecord API call fails
- **THEN** the error SHALL be wrapped with tccommon.RetryError and retried within tccommon.WriteRetryTimeout

### Requirement: Read DNS record via DescribeDnsRecords API
The resource SHALL call DescribeDnsRecords API with ZoneId and filter by record ID to read a single DNS record's state. The composite ID SHALL be parsed to extract zone_id and record_id.

#### Scenario: Record found successfully
- **WHEN** DescribeDnsRecords returns a matching record for the given zone_id and record_id
- **THEN** all schema fields SHALL be set from the DnsRecord response (zone_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on)

#### Scenario: Record not found
- **WHEN** DescribeDnsRecords returns no matching record
- **THEN** the resource ID SHALL be cleared (d.SetId("")) and the resource SHALL be marked as removed

#### Scenario: Read sets nil-check before setting fields
- **WHEN** response fields are nil
- **THEN** the corresponding d.Set() calls SHALL be skipped to avoid setting nil values

### Requirement: Update DNS record via ModifyDnsRecords API
The resource SHALL call ModifyDnsRecords API when mutable fields (name, type, content, location, ttl, weight, priority) change. The request SHALL include ZoneId and a DnsRecord struct with RecordId and the updated field values.

#### Scenario: Update mutable DNS record fields
- **WHEN** any of name, type, content, location, ttl, weight, or priority fields change
- **THEN** ModifyDnsRecords SHALL be called with ZoneId and a DnsRecord containing RecordId and the new values

#### Scenario: No mutable fields changed
- **WHEN** no mutable fields have changed
- **THEN** ModifyDnsRecords SHALL NOT be called

### Requirement: Update DNS record status via ModifyDnsRecordsStatus API
The resource SHALL call ModifyDnsRecordsStatus API when the status field changes. When status is "enable", RecordsToEnable SHALL contain the record_id. When status is "disable", RecordsToDisable SHALL contain the record_id.

#### Scenario: Enable DNS record
- **WHEN** status is changed to "enable"
- **THEN** ModifyDnsRecordsStatus SHALL be called with RecordsToEnable containing the record_id

#### Scenario: Disable DNS record
- **WHEN** status is changed to "disable"
- **THEN** ModifyDnsRecordsStatus SHALL be called with RecordsToDisable containing the record_id

#### Scenario: No status change
- **WHEN** status field has not changed
- **THEN** ModifyDnsRecordsStatus SHALL NOT be called

### Requirement: Delete DNS record via DeleteDnsRecords API
The resource SHALL call DeleteDnsRecords API with ZoneId and RecordIds containing the record_id to delete a DNS record.

#### Scenario: Successful deletion
- **WHEN** DeleteDnsRecords is called with zone_id and record_ids
- **THEN** the API SHALL be called with ZoneId and RecordIds parameters

#### Scenario: Delete API call failure
- **WHEN** DeleteDnsRecords API call fails
- **THEN** the error SHALL be wrapped with tccommon.RetryError and retried within tccommon.WriteRetryTimeout

### Requirement: Resource supports import
The resource SHALL support Terraform import with composite ID format `zone_id#record_id`.

#### Scenario: Import with valid composite ID
- **WHEN** a user imports the resource with ID `zone_id#record_id`
- **THEN** the resource SHALL parse the ID and call Read to populate the state

### Requirement: Resource registration in provider
The resource SHALL be registered in `tencentcloud/provider.go` with name `tencentcloud_teo_dns_record_v4` and documented in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_dns_record_v4` SHALL be available as a resource type

### Requirement: Unit tests with gomonkey mocks
The resource SHALL have unit tests in `resource_tc_teo_dns_record_v4_test.go` using gomonkey to mock cloud API calls. Tests SHALL cover Create, Read, Update, and Delete operations.

#### Scenario: Unit test for Create operation
- **WHEN** CreateDnsRecord API is mocked to return a RecordId
- **THEN** the test SHALL verify the resource ID is correctly set to `zone_id#record_id`

#### Scenario: Unit test for Read operation
- **WHEN** DescribeDnsRecords API is mocked to return a DnsRecord
- **THEN** the test SHALL verify all schema fields are correctly populated

#### Scenario: Unit test for Update operation
- **WHEN** ModifyDnsRecords API is mocked to succeed
- **THEN** the test SHALL verify the update request contains the correct parameters

#### Scenario: Unit test for Delete operation
- **WHEN** DeleteDnsRecords API is mocked to succeed
- **THEN** the test SHALL verify the delete request contains ZoneId and RecordIds
