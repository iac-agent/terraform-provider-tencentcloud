## Requirements

### Requirement: Resource CRUD lifecycle
The system SHALL provide a Terraform resource `tencentcloud_teo_dns_record_17` that supports full CRUD lifecycle (Create, Read, Update, Delete) for TEO DNS records using the cloud APIs: CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords.

#### Scenario: Create a DNS record
- **WHEN** a user creates a `tencentcloud_teo_dns_record_17` resource with required fields (zone_id, name, type, content) and optional fields (location, ttl, weight, priority)
- **THEN** the system SHALL call CreateDnsRecord API with all provided fields, set the resource ID to `zone_id#record_id` using the returned RecordId, and refresh state by calling Read

#### Scenario: Read a DNS record
- **WHEN** a user reads a `tencentcloud_teo_dns_record_17` resource
- **THEN** the system SHALL parse the composite ID to extract zone_id and record_id, call DescribeTeoDnsRecordById to fetch the record, and set all fields (zone_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on) into the Terraform state

#### Scenario: Read a non-existent DNS record
- **WHEN** a user reads a `tencentcloud_teo_dns_record_17` resource that has been deleted
- **THEN** the system SHALL set the resource ID to empty string and log a warning, signaling Terraform to remove the resource from state

#### Scenario: Update DNS record content fields
- **WHEN** a user updates mutable fields (name, type, content, location, ttl, weight, priority) of a `tencentcloud_teo_dns_record_17` resource
- **THEN** the system SHALL call ModifyDnsRecords API with a DnsRecord struct containing the record_id and updated field values

#### Scenario: Update DNS record status
- **WHEN** a user updates the status field of a `tencentcloud_teo_dns_record_17` resource
- **THEN** the system SHALL call ModifyDnsRecordsStatus API with the record_id in RecordsToEnable (if status is "enable") or RecordsToDisable (if status is "disable")

#### Scenario: Delete a DNS record
- **WHEN** a user deletes a `tencentcloud_teo_dns_record_17` resource
- **THEN** the system SHALL call DeleteDnsRecords API with zone_id and a list containing the record_id

### Requirement: Resource schema definition
The system SHALL define a schema for `tencentcloud_teo_dns_record_17` with the following fields:

#### Scenario: Required fields
- **WHEN** the resource schema is defined
- **THEN** the following fields SHALL be Required: zone_id (TypeString, ForceNew), name (TypeString), type (TypeString), content (TypeString)

#### Scenario: Optional and computed fields
- **WHEN** the resource schema is defined
- **THEN** the following fields SHALL be Optional+Computed: location (TypeString), ttl (TypeInt), weight (TypeInt), priority (TypeInt), status (TypeString)

#### Scenario: Computed-only fields
- **WHEN** the resource schema is defined
- **THEN** the following fields SHALL be Computed only: created_on (TypeString), modified_on (TypeString)

### Requirement: Resource import support
The system SHALL support importing existing DNS records via `terraform import` using the composite ID format `zone_id#record_id`.

#### Scenario: Import an existing DNS record
- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_17.example zone_id#record_id`
- **THEN** the system SHALL parse the composite ID, call Read to fetch the current state, and populate the Terraform state with the DNS record data

### Requirement: Provider registration
The system SHALL register the `tencentcloud_teo_dns_record_17` resource in the provider's ResourcesMap in `tencentcloud/provider.go` and add an entry in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_dns_record_17` SHALL be available and mapped to its resource definition function

### Requirement: Unit tests
The system SHALL include unit tests for the `tencentcloud_teo_dns_record_17` resource using gomonkey mock approach.

#### Scenario: Test CRUD operations with mocks
- **WHEN** the unit tests are run
- **THEN** the tests SHALL mock the cloud API calls (CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, ModifyDnsRecordsStatus, DeleteDnsRecords) and verify the resource CRUD logic works correctly

### Requirement: Resource documentation
The system SHALL include a markdown documentation file for the resource at `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md`.

#### Scenario: Documentation includes usage example
- **WHEN** the documentation file is generated
- **THEN** it SHALL contain a one-line description mentioning TEO, an example usage section with HCL configuration, and import instructions using the composite ID format

### Requirement: Retry and error handling
The system SHALL implement proper retry logic and error handling for all API calls.

#### Scenario: Write operations use WriteRetryTimeout
- **WHEN** Create, Update, or Delete operations call the cloud API
- **THEN** the system SHALL wrap the API call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and use `tccommon.RetryError()` for error classification

#### Scenario: Read operations use ReadRetryTimeout
- **WHEN** the service layer reads DNS records via DescribeDnsRecords
- **THEN** the system SHALL wrap the API call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check()` and use `tccommon.RetryError()` for error classification

### Requirement: Nil check on Read
The system SHALL check for nil values in the API response before setting fields in the Terraform state during Read operations.

#### Scenario: Response fields are nil-safe
- **WHEN** the Read function processes the API response
- **THEN** each field SHALL only be set via `d.Set()` if the corresponding response field is not nil

### Requirement: Create response validation
The system SHALL validate the Create API response before proceeding.

#### Scenario: Validate create response
- **WHEN** the CreateDnsRecord API returns a response
- **THEN** the system SHALL check that the response is not nil, the RecordId field is not nil and not empty, and return NonRetryableError if validation fails
