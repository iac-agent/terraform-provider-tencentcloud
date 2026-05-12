## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_dns_record_10` with the following schema fields:

Required fields:
- `zone_id` (TypeString, ForceNew): Zone ID
- `name` (TypeString): DNS record name
- `type` (TypeString): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (TypeString): DNS record content

Optional + Computed fields:
- `location` (TypeString): DNS record resolution route
- `ttl` (TypeInt): Cache time in seconds, range 60-86400
- `weight` (TypeInt): DNS record weight, range -1~100
- `priority` (TypeInt): MX record priority, range 0~50

Computed fields:
- `record_id` (TypeString): DNS record ID returned by cloud API
- `status` (TypeString): DNS record resolution status (enable/disable)
- `created_on` (TypeString): Creation time
- `modified_on` (TypeString): Modification time

The resource SHALL support Terraform Import via composite ID `zone_id#record_id`.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is registered
- **THEN** all required fields (zone_id, name, type, content) are defined as Required with correct types
- **AND** optional fields (location, ttl, weight, priority) are defined as Optional and Computed
- **AND** computed fields (record_id, status, created_on, modified_on) are defined as Computed
- **AND** zone_id is marked as ForceNew

#### Scenario: Resource supports import
- **WHEN** user imports an existing DNS record with composite ID `zone_id#record_id`
- **THEN** the resource SHALL parse the ID and populate zone_id and record_id fields
- **AND** read the full resource state from cloud API

### Requirement: Resource Create Operation
The system SHALL implement Create operation by calling `CreateDnsRecord` cloud API with all input parameters mapped from Terraform schema. Upon successful creation, the system SHALL store the composite ID as `zone_id#record_id` using FILED_SP separator.

#### Scenario: Successful DNS record creation
- **WHEN** user creates a DNS record with zone_id, name, type, content, and optional parameters
- **THEN** the system calls CreateDnsRecord API with all provided parameters
- **AND** sets resource ID to `zone_id#record_id`
- **AND** calls Read to refresh the full state

#### Scenario: Create API returns empty RecordId
- **WHEN** CreateDnsRecord API succeeds but returns empty RecordId
- **THEN** the system SHALL return a NonRetryableError

#### Scenario: Create API call fails
- **WHEN** CreateDnsRecord API call fails
- **THEN** the system SHALL wrap the error with RetryError and retry with WriteRetryTimeout
- **AND** log the failure reason

### Requirement: Resource Read Operation
The system SHALL implement Read operation by calling `DescribeDnsRecords` cloud API with ZoneId and Filter (id=recordId) to query a single DNS record. The composite ID SHALL be parsed to extract zone_id and record_id.

#### Scenario: Successful DNS record read
- **WHEN** user reads a DNS record with composite ID `zone_id#record_id`
- **THEN** the system parses the ID to extract zone_id and record_id
- **AND** calls DescribeDnsRecords with ZoneId and Filter(id=recordId)
- **AND** sets all schema fields from the returned DnsRecord object
- **AND** checks each field is not nil before setting

#### Scenario: DNS record not found
- **WHEN** DescribeDnsRecords returns empty result for the given record_id
- **THEN** the system SHALL set resource ID to empty string to indicate resource was deleted
- **AND** log a warning message

#### Scenario: Invalid composite ID format
- **WHEN** the resource ID does not contain exactly 2 parts separated by FILED_SP
- **THEN** the system SHALL return an error indicating the ID is broken

### Requirement: Resource Update Operation
The system SHALL implement Update operation by calling `ModifyDnsRecords` cloud API. The system SHALL build a DnsRecord object with RecordId and all mutable fields, then pass it in the DnsRecords list parameter.

#### Scenario: Successful DNS record update
- **WHEN** user updates mutable fields (name, type, content, location, ttl, weight, priority)
- **THEN** the system detects which fields have changed
- **AND** calls ModifyDnsRecords API with ZoneId and DnsRecords containing the modified record
- **AND** the DnsRecord object includes RecordId but NOT ZoneId (as per API requirement)
- **AND** calls Read to refresh the full state

#### Scenario: No fields changed
- **WHEN** no mutable fields have changed
- **THEN** the system SHALL skip the ModifyDnsRecords API call
- **AND** proceed directly to Read

#### Scenario: Update API call fails
- **WHEN** ModifyDnsRecords API call fails
- **THEN** the system SHALL wrap the error with RetryError and retry with WriteRetryTimeout
- **AND** log the failure reason

### Requirement: Resource Delete Operation
The system SHALL implement Delete operation by calling `DeleteDnsRecords` cloud API with ZoneId and RecordIds containing the record_id.

#### Scenario: Successful DNS record deletion
- **WHEN** user deletes a DNS record with composite ID `zone_id#record_id`
- **THEN** the system parses the ID to extract zone_id and record_id
- **AND** calls DeleteDnsRecords API with ZoneId and RecordIds=[record_id]
- **AND** the system SHALL retry with WriteRetryTimeout on API failure

#### Scenario: Delete API call fails
- **WHEN** DeleteDnsRecords API call fails
- **THEN** the system SHALL wrap the error with RetryError and retry
- **AND** log the failure reason

### Requirement: Resource Registration in Provider
The system SHALL register the `tencentcloud_teo_dns_record_10` resource in `tencentcloud/provider.go` and add documentation entry in `tencentcloud/provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_dns_record_10` is available for use in Terraform configurations
- **AND** the resource is listed in the provider documentation

### Requirement: Service Layer Method
The system SHALL implement a service layer method `DescribeTeoDnsRecord10ById` in `tencentcloud/services/teo/service_tencentcloud_teo.go` that queries a single DNS record by zone_id and record_id using DescribeDnsRecords API with AdvancedFilter.

#### Scenario: Service method queries DNS record
- **WHEN** the service method is called with zone_id and record_id
- **THEN** it calls DescribeDnsRecords with ZoneId and Filter(id=recordId)
- **AND** returns the first matching DnsRecord or nil if not found
- **AND** retries on API failure with ReadRetryTimeout

### Requirement: Unit Tests
The system SHALL provide unit tests using gomonkey mock approach for all CRUD operations in `resource_tc_teo_dns_record_10_test.go`.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** the test file is executed with `go test -gcflags=all=-l`
- **THEN** all Create, Read, Update, Delete operations are tested with mocked cloud API responses
- **AND** all tests pass successfully

### Requirement: Resource Documentation
The system SHALL provide resource documentation in `resource_tc_teo_dns_record_10.md` with one-line description, example usage, and import section. The description SHALL include the product name TEO.

#### Scenario: Documentation file exists
- **WHEN** the resource is added
- **THEN** a .md file exists with description, example usage, and import section
- **AND** the description follows the format "Provides a resource to ..."
- **AND** the import section explains the composite ID format
