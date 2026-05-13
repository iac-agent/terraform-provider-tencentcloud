## ADDED Requirements

### Requirement: Resource CRUD operations
The system SHALL provide a Terraform resource `tencentcloud_teo_dns_record_15` that supports full CRUD lifecycle (Create, Read, Update, Delete) for TEO DNS records using the following cloud APIs: CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, DeleteDnsRecords.

#### Scenario: Create DNS record
- **WHEN** user applies a Terraform configuration with `tencentcloud_teo_dns_record_15` specifying zone_id, name, type, content, and optional fields (location, ttl, weight, priority)
- **THEN** the resource SHALL call CreateDnsRecord API with the provided parameters, set the resource ID to `zone_id + FILED_SP + record_id`, and populate all computed fields via Read

#### Scenario: Read DNS record
- **WHEN** Terraform performs a refresh on an existing `tencentcloud_teo_dns_record_15` resource
- **THEN** the resource SHALL parse the composite ID to extract zone_id and record_id, call DescribeDnsRecords with AdvancedFilter filtering by id, and populate all schema fields from the returned DnsRecord

#### Scenario: Read DNS record not found
- **WHEN** the DescribeDnsRecords query returns no matching record for the given record_id
- **THEN** the resource SHALL set the resource ID to empty string to indicate the resource has been deleted

#### Scenario: Update DNS record
- **WHEN** user modifies mutable fields (name, type, content, location, ttl, weight, priority) of an existing `tencentcloud_teo_dns_record_15` resource
- **THEN** the resource SHALL call ModifyDnsRecords API with the zone_id and a DnsRecord object containing the record_id and modified fields

#### Scenario: Update immutable field
- **WHEN** user attempts to modify the zone_id field
- **THEN** Terraform SHALL treat it as a ForceNew operation, destroying and recreating the resource

#### Scenario: Delete DNS record
- **WHEN** user destroys a `tencentcloud_teo_dns_record_15` resource
- **THEN** the resource SHALL call DeleteDnsRecords API with the zone_id and record_id to delete the DNS record

### Requirement: Composite ID format
The resource SHALL use a composite ID consisting of zone_id and record_id separated by `tccommon.FILED_SP`.

#### Scenario: ID construction after create
- **WHEN** CreateDnsRecord returns a RecordId
- **THEN** the resource ID SHALL be set to `zone_id + FILED_SP + record_id`

#### Scenario: ID parsing for read/update/delete
- **WHEN** the resource needs to perform Read, Update, or Delete operations
- **THEN** the resource SHALL split the composite ID by FILED_SP to extract zone_id and record_id

### Requirement: Import support
The resource SHALL support Terraform import using the composite ID format.

#### Scenario: Import existing DNS record
- **WHEN** user runs `terraform import tencentcloud_teo_dns_record_15.example "zone_id#record_id"`
- **THEN** the resource SHALL parse the imported ID and populate all fields via the Read operation

### Requirement: Schema field definitions
The resource SHALL define the following schema fields with correct types and behaviors:

- `zone_id`: TypeString, Required, ForceNew
- `name`: TypeString, Required
- `type`: TypeString, Required
- `content`: TypeString, Required
- `location`: TypeString, Optional, Computed
- `ttl`: TypeInt, Optional, Computed
- `weight`: TypeInt, Optional, Computed
- `priority`: TypeInt, Optional, Computed
- `record_id`: TypeString, Computed
- `status`: TypeString, Computed
- `created_on`: TypeString, Computed
- `modified_on`: TypeString, Computed

#### Scenario: Required fields validation
- **WHEN** user creates a resource without specifying zone_id, name, type, or content
- **THEN** Terraform SHALL return a validation error indicating missing required fields

#### Scenario: Computed fields population
- **WHEN** a resource is created or read
- **THEN** the system SHALL populate record_id, status, created_on, and modified_on from the API response

### Requirement: Retry logic for API calls
All CRUD operations SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations, and SHALL use `tccommon.RetryError()` for error wrapping.

#### Scenario: API call with transient failure
- **WHEN** a cloud API call returns a retriable error
- **THEN** the resource SHALL retry the API call within the configured timeout

#### Scenario: API call with non-retriable failure
- **WHEN** a cloud API call returns a non-retriable error
- **THEN** the resource SHALL return the error without retrying

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` with the key `tencentcloud_teo_dns_record_15` and listed in `tencentcloud/provider.md` under the TEO section.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the `tencentcloud_teo_dns_record_15` resource SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey mock
The resource SHALL have unit tests using gomonkey to mock cloud API calls, testing business logic without requiring real cloud credentials.

#### Scenario: Create operation test
- **WHEN** the create function is called with valid input
- **THEN** the test SHALL verify that CreateDnsRecord is called with correct parameters and the composite ID is set correctly

#### Scenario: Read operation test
- **WHEN** the read function is called for an existing resource
- **THEN** the test SHALL verify that DescribeDnsRecords is called and all fields are populated correctly

#### Scenario: Update operation test
- **WHEN** the update function is called with changed mutable fields
- **THEN** the test SHALL verify that ModifyDnsRecords is called with the correct parameters

#### Scenario: Delete operation test
- **WHEN** the delete function is called
- **THEN** the test SHALL verify that DeleteDnsRecords is called with the correct zone_id and record_id

### Requirement: Nil field check in Read
The Read function SHALL check if response fields are nil before calling d.Set() to avoid nil pointer dereference.

#### Scenario: API returns nil optional fields
- **WHEN** DescribeDnsRecords returns a DnsRecord with some nil fields (e.g., location, weight, priority)
- **THEN** the Read function SHALL skip d.Set() calls for those nil fields without error

### Requirement: Create response validation
The Create function SHALL validate that the CreateDnsRecord response contains a non-nil RecordId. If RecordId is nil, the function SHALL return a NonRetryableError.

#### Scenario: CreateDnsRecord returns empty RecordId
- **WHEN** CreateDnsRecord API returns a response with nil or empty RecordId
- **THEN** the Create function SHALL return a NonRetryableError indicating the record ID is empty
