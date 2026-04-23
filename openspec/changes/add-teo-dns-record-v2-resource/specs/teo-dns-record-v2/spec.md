## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_dns_record_v2` SHALL define a Terraform schema with the following fields:
- `zone_id` (String, Required, ForceNew): Zone ID
- `name` (String, Required): DNS record name
- `type` (String, Required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (String, Required): DNS record content
- `location` (String, Optional, Computed): DNS resolution line, default "Default"
- `ttl` (Int, Optional, Computed): Cache time in seconds, range 60-86400, default 300
- `weight` (Int, Optional, Computed): Record weight, range -1 to 100, default -1
- `priority` (Int, Optional, Computed): MX priority, range 0-50, default 0
- `record_id` (String, Computed): DNS record ID returned by CreateDnsRecord
- `status` (String, Computed): DNS record status ("enable" or "disable")
- `created_on` (String, Computed): Creation time
- `modified_on` (String, Computed): Modification time

#### Scenario: Schema fields match API parameters
- **WHEN** the resource schema is defined
- **THEN** all CreateDnsRecord input parameters (ZoneId, Name, Type, Content, Location, TTL, Weight, Priority) SHALL have corresponding schema fields
- **AND** CreateDnsRecord output parameter (RecordId) SHALL map to computed schema field `record_id`
- **AND** DescribeDnsRecords DnsRecord output fields (Status, CreatedOn, ModifiedOn) SHALL map to computed schema fields

### Requirement: Resource ID Format
The resource SHALL use a composite ID format `zone_id#record_id` using `tccommon.FILED_SP` as the separator.

#### Scenario: Create sets composite ID
- **WHEN** CreateDnsRecord returns a RecordId
- **THEN** the resource ID SHALL be set to `zone_id + tccommon.FILED_SP + record_id`

#### Scenario: Read extracts zone_id and record_id from d.Get()
- **WHEN** the Read function is called
- **THEN** zone_id SHALL be obtained from `d.Get("zone_id")` and record_id SHALL be obtained from `d.Get("record_id")`, NOT from splitting d.Id()

### Requirement: Create Operation
The resource SHALL call `CreateDnsRecord` API to create a DNS record.

#### Scenario: Successful creation
- **WHEN** `terraform apply` is executed with valid parameters (zone_id, name, type, content)
- **THEN** the system SHALL call `CreateDnsRecord` with all required and optional parameters
- **AND** the response RecordId SHALL be set to the `record_id` schema field
- **AND** the resource ID SHALL be set to the composite format

#### Scenario: Create with retry
- **WHEN** CreateDnsRecord API call fails with a retryable error
- **THEN** the system SHALL retry using `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- **AND** the error SHALL be wrapped with `tccommon.RetryError(e)`

### Requirement: Read Operation
The resource SHALL call `DescribeDnsRecords` API with AdvancedFilter to query the DNS record by record_id.

#### Scenario: Record exists
- **WHEN** the Read function is called and the DNS record exists
- **THEN** the system SHALL call `DescribeDnsRecords` with ZoneId and AdvancedFilter (name="id", values=[record_id])
- **AND** all schema fields SHALL be populated from the API response

#### Scenario: Record not found
- **WHEN** the Read function is called and the DNS record does not exist
- **THEN** the resource SHALL be marked as gone by setting `d.SetId("")`

#### Scenario: Read with retry
- **WHEN** DescribeDnsRecords API call fails with a retryable error
- **THEN** the system SHALL retry using `resource.Retry(tccommon.ReadRetryTimeout, ...)`

### Requirement: Update Operation
The resource SHALL call `ModifyDnsRecords` API to update the DNS record when mutable fields change.

#### Scenario: Update mutable fields
- **WHEN** any of name, type, content, location, ttl, weight, priority fields change
- **THEN** the system SHALL call `ModifyDnsRecords` with ZoneId and a DnsRecords list containing one entry with RecordId and the changed fields
- **AND** read-only fields (ZoneId, Status, CreatedOn, ModifiedOn) SHALL NOT be included in the DnsRecord input

#### Scenario: Update with retry
- **WHEN** ModifyDnsRecords API call fails with a retryable error
- **THEN** the system SHALL retry using `resource.Retry(tccommon.WriteRetryTimeout, ...)`

### Requirement: Delete Operation
The resource SHALL call `DeleteDnsRecords` API to delete the DNS record.

#### Scenario: Successful deletion
- **WHEN** `terraform destroy` is executed
- **THEN** the system SHALL call `DeleteDnsRecords` with ZoneId and RecordIds containing the record_id

#### Scenario: Delete with retry
- **WHEN** DeleteDnsRecords API call fails with a retryable error
- **THEN** the system SHALL retry using `resource.Retry(tccommon.WriteRetryTimeout, ...)`

### Requirement: Import Support
The resource SHALL support Terraform import with `schema.ImportStatePassthrough`.

#### Scenario: Import existing DNS record
- **WHEN** `terraform import` is executed with ID format `zone_id#record_id`
- **THEN** the Read function SHALL be called to populate the resource state

### Requirement: Provider Registration
The resource SHALL be registered in `provider.go` as `"tencentcloud_teo_dns_record_v2"` and documented in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is loaded
- **THEN** the resource `tencentcloud_teo_dns_record_v2` SHALL be available for use in Terraform configurations

### Requirement: Documentation
The resource SHALL have a corresponding `.md` documentation file following the gendoc format.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_teo_dns_record_v2.md` file SHALL exist with description, example usage, and import section

#### Scenario: Documentation includes multiple example scenarios
- **WHEN** the documentation is written
- **THEN** the example usage SHALL include scenarios for A record, CNAME record, MX record, TXT record, and AAAA record with custom resolution route
- **AND** each scenario SHALL demonstrate the relevant parameters for that record type

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock for business logic verification.

#### Scenario: Test CRUD operations
- **WHEN** unit tests are executed
- **THEN** Create, Read, Update, and Delete functions SHALL be tested with mocked cloud API calls
