## Requirements

### Requirement: Create DNS Record
The resource SHALL create a TEO DNS record by calling `CreateDnsRecord` API with zone_id, name, type, content, and optional fields (location, ttl, weight, priority). Upon successful creation, the resource SHALL set its ID to `{zoneId}{FILED_SP}{recordId}` using the returned RecordId.

#### Scenario: Create a basic A record
- **WHEN** user creates a `tencentcloud_teo_dns_record_v7` resource with zone_id, name="www", type="A", content="1.2.3.4"
- **THEN** the resource calls CreateDnsRecord API and sets the ID to `{zoneId}{FILED_SP}{recordId}`

#### Scenario: Create with optional fields
- **WHEN** user creates a resource with zone_id, name, type, content, location="XX", ttl=600, weight=10, priority=5
- **THEN** the resource calls CreateDnsRecord with all specified optional fields

#### Scenario: Create returns empty RecordId
- **WHEN** CreateDnsRecord API returns success but RecordId is nil
- **THEN** the resource SHALL return a NonRetryableError

### Requirement: Read DNS Record
The resource SHALL read the current state of a DNS record by calling `DescribeTeoDnsRecordById` service method, which wraps `DescribeDnsRecords` API with an id filter. If the record is not found, the resource SHALL remove itself from state.

#### Scenario: Record exists
- **WHEN** resource reads with valid ID `{zoneId}{FILED_SP}{recordId}`
- **THEN** the resource calls DescribeTeoDnsRecordById with zoneId and recordId, and sets all schema fields from the response

#### Scenario: Record not found
- **WHEN** DescribeTeoDnsRecordById returns nil
- **THEN** the resource SHALL set d.SetId("") and log a warning

#### Scenario: Nil check on response fields
- **WHEN** response fields (ZoneId, Name, Type, Content, Location, TTL, Weight, Priority, Status, CreatedOn, ModifiedOn) are nil
- **THEN** the resource SHALL skip calling d.Set() for those nil fields

### Requirement: Update DNS Record Content
The resource SHALL update DNS record content fields (name, type, content, location, ttl, weight, priority) by calling `ModifyDnsRecords` API when any of these fields have changed.

#### Scenario: Update content fields
- **WHEN** user changes name, type, content, location, ttl, weight, or priority
- **THEN** the resource calls ModifyDnsRecords API with a DnsRecord object containing RecordId and all current field values

#### Scenario: No content change
- **WHEN** none of the content fields have changed
- **THEN** the resource SHALL skip calling ModifyDnsRecords API

### Requirement: Update DNS Record Status
The resource SHALL update DNS record status by calling `ModifyDnsRecordsStatus` API when the status field has changed.

#### Scenario: Enable record
- **WHEN** user changes status to "enable"
- **THEN** the resource calls ModifyDnsRecordsStatus with RecordsToEnable containing the recordId

#### Scenario: Disable record
- **WHEN** user changes status to "disable"
- **THEN** the resource calls ModifyDnsRecordsStatus with RecordsToDisable containing the recordId

#### Scenario: No status change
- **WHEN** status field has not changed
- **THEN** the resource SHALL skip calling ModifyDnsRecordsStatus API

### Requirement: Delete DNS Record
The resource SHALL delete a DNS record by calling `DeleteDnsRecords` API with ZoneId and RecordIds.

#### Scenario: Delete existing record
- **WHEN** user destroys the resource
- **THEN** the resource calls DeleteDnsRecords API with ZoneId and RecordIds=[recordId]

### Requirement: Import DNS Record
The resource SHALL support import using the composite ID format `{zoneId}{FILED_SP}{recordId}`.

#### Scenario: Import existing record
- **WHEN** user imports with ID `{zoneId}{FILED_SP}{recordId}`
- **THEN** the resource SHALL parse the ID, call Read, and populate the state

### Requirement: Resource Registration
The resource SHALL be registered in `provider.go` as `tencentcloud_teo_dns_record_v7` and referenced as `teo.ResourceTencentCloudTeoDnsRecordV7()`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_dns_record_v7` SHALL be available in the ResourcesMap

### Requirement: API Retry Logic
All cloud API calls in CRUD functions SHALL be wrapped with `resource.Retry` using appropriate timeout (WriteRetryTimeout for Create/Update/Delete, ReadRetryTimeout for Read). Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: API call with retry
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the resource SHALL retry the call within the timeout period

#### Scenario: API call with non-retryable error
- **WHEN** a cloud API call fails with a non-retryable error
- **THEN** the resource SHALL return the error immediately

### Requirement: ID Parsing
All CRUD functions (Read, Update, Delete) SHALL parse the composite ID `{zoneId}{FILED_SP}{recordId}` and return an error if the ID format is broken.

#### Scenario: Valid ID format
- **WHEN** d.Id() returns `{zoneId}{FILED_SP}{recordId}`
- **THEN** the function SHALL parse zoneId and recordId correctly

#### Scenario: Broken ID format
- **WHEN** d.Id() returns an invalid format (not 2 parts separated by FILED_SP)
- **THEN** the function SHALL return an error "id is broken"
