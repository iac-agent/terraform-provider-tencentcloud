## ADDED Requirements

### Requirement: Create TEO DNS Record
The system SHALL allow users to create a TEO DNS record by calling the CreateDnsRecord API with zone_id, name, type, content, and optional location, ttl, weight, priority parameters. Upon successful creation, the resource SHALL set its ID to the composite format `zone_id#record_id` and store the returned record_id.

#### Scenario: Successful DNS record creation
- **WHEN** user creates a tencentcloud_teo_dns_record resource with zone_id, name="a.example.com", type="A", content="1.2.3.4"
- **THEN** the system calls CreateDnsRecord API, sets resource ID to `zone_id#record_id`, and stores the returned record_id

#### Scenario: DNS record creation with all optional parameters
- **WHEN** user creates a resource with zone_id, name, type, content, location="Default", ttl=300, weight=-1, priority=0
- **THEN** the system calls CreateDnsRecord API with all parameters and sets the composite ID

#### Scenario: DNS record creation fails with nil response
- **WHEN** CreateDnsRecord API returns nil response or nil RecordId
- **THEN** the system returns an error indicating creation failure

### Requirement: Read TEO DNS Record
The system SHALL allow users to read a TEO DNS record by calling TeoService.DescribeTeoDnsRecordById with zone_id and record_id parsed from the composite ID. The system SHALL set all schema fields from the API response, with nil checks before setting each field. If the record is not found, the system SHALL clear the resource ID and log a warning.

#### Scenario: Successful DNS record read
- **WHEN** user reads a tencentcloud_teo_dns_record resource with valid composite ID
- **THEN** the system parses zone_id and record_id from the ID, calls DescribeTeoDnsRecordById, and sets all fields (zone_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on) from the response

#### Scenario: DNS record not found
- **WHEN** DescribeTeoDnsRecordById returns nil for a given zone_id and record_id
- **THEN** the system sets resource ID to empty string and logs a warning

### Requirement: Update TEO DNS Record
The system SHALL allow users to update a TEO DNS record's content fields by calling the ModifyDnsRecords API when any of the mutable arguments (name, type, content, location, ttl, weight, priority) change. The system SHALL also allow users to update the record's enabled/disabled status by calling the ModifyDnsRecordsStatus API when the status field changes. Each update path SHALL only trigger when the relevant fields have changed.

#### Scenario: Update DNS record name
- **WHEN** user updates the name field of a tencentcloud_teo_dns_record resource
- **THEN** the system detects the change, calls ModifyDnsRecords with the updated DnsRecord containing RecordId and new name

#### Scenario: Update multiple DNS record fields
- **WHEN** user updates both type and content fields
- **THEN** the system calls ModifyDnsRecords with a DnsRecord containing RecordId and both updated fields

#### Scenario: Update DNS record status to disable
- **WHEN** user updates the status field from "enable" to "disable"
- **THEN** the system calls ModifyDnsRecordsStatus with RecordsToDisable containing the record_id

#### Scenario: Update DNS record status to enable
- **WHEN** user updates the status field from "disable" to "enable"
- **THEN** the system calls ModifyDnsRecordsStatus with RecordsToEnable containing the record_id

#### Scenario: No mutable fields changed
- **WHEN** user does not change any mutable fields
- **THEN** the system skips both ModifyDnsRecords and ModifyDnsRecordsStatus API calls and directly calls Read

### Requirement: Delete TEO DNS Record
The system SHALL allow users to delete a TEO DNS record by calling the DeleteDnsRecords API with zone_id and record_id parsed from the composite ID.

#### Scenario: Successful DNS record deletion
- **WHEN** user deletes a tencentcloud_teo_dns_record resource
- **THEN** the system parses zone_id and record_id from the composite ID, calls DeleteDnsRecords with RecordIds containing the single record_id

### Requirement: Import TEO DNS Record
The system SHALL support importing existing TEO DNS records using the composite ID format `zone_id#record_id` via schema.ImportStatePassthrough.

#### Scenario: Import existing DNS record
- **WHEN** user imports a tencentcloud_teo_dns_record resource with ID "zone-abc#rec-xyz"
- **THEN** the system parses the composite ID and calls Read to populate the resource state

### Requirement: Resource Registration
The system SHALL register the tencentcloud_teo_dns_record resource in provider.go and add documentation entry in provider.md.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the tencentcloud_teo_dns_record resource is registered and available for use

### Requirement: Documentation
The system SHALL include a resource documentation file (resource_tc_teo_dns_record.md) with a one-line description, Example Usage section, and Import section. The website/docs/ directory documentation SHALL be generated via make doc.

#### Scenario: Documentation file exists
- **WHEN** the resource code is complete
- **THEN** a resource_tc_teo_dns_record.md file exists with proper description, usage examples, and import instructions