## ADDED Requirements

### Requirement: Resource creation via CreateDnsRecord API

The resource SHALL support creating TEO DNS records by calling the `CreateDnsRecord` API with user-provided configuration including zone_id, name, type, content, and optional fields (location, ttl, weight, priority).

#### Scenario: Successful DNS record creation
- **WHEN** user defines `tencentcloud_teo_dns_record_21` resource with valid parameters (zone_id, name, type, content)
- **THEN** provider calls `CreateDnsRecord` API, stores the returned RecordId, and sets the composite ID to `zone_id#record_id`

#### Scenario: Creation with invalid parameters
- **WHEN** user provides invalid parameters (e.g., unsupported record type)
- **THEN** provider returns clear error message from API validation

#### Scenario: API call failure during creation
- **WHEN** `CreateDnsRecord` API call fails due to network or service error
- **THEN** provider retries with exponential backoff and returns error if all retries exhausted

### Requirement: Resource reading via DescribeDnsRecords API

The resource SHALL support reading DNS record details by calling the `DescribeDnsRecords` API with zone_id and filtering by record_id, using the `AdvancedFilter` struct with name="id" and values=[recordId].

#### Scenario: Successful resource state refresh
- **WHEN** Terraform performs a refresh operation
- **THEN** provider calls `DescribeDnsRecords` API with zone_id and filter by record_id, updates local state with current values (name, type, content, location, ttl, weight, priority, status, created_on, modified_on)

#### Scenario: Resource not found during read
- **WHEN** record_id no longer exists in cloud (manually deleted)
- **THEN** provider removes resource from Terraform state by setting `d.SetId("")` without error

#### Scenario: API call failure during read
- **WHEN** `DescribeDnsRecords` API call fails due to network or service error
- **THEN** provider retries with `tccommon.ReadRetryTimeout` and returns error if all retries exhausted

### Requirement: Resource update via ModifyDnsRecords API

The resource SHALL support updating DNS records by calling the `ModifyDnsRecords` API when mutable fields change. The API SHALL be called with ZoneId and a DnsRecord list containing the RecordId and changed fields.

#### Scenario: Successful DNS record update
- **WHEN** user modifies mutable fields (name, type, content, location, ttl, weight, priority) in the resource configuration
- **THEN** provider calls `ModifyDnsRecords` API with the updated values and Terraform state is refreshed

#### Scenario: No changes detected
- **WHEN** Terraform plan shows no changes to mutable fields
- **THEN** provider does not call `ModifyDnsRecords` API

#### Scenario: Update with multiple field changes
- **WHEN** user modifies multiple mutable fields simultaneously
- **THEN** provider sends a single `ModifyDnsRecords` API call with all changed fields

### Requirement: Resource deletion via DeleteDnsRecords API

The resource SHALL support deleting DNS records by calling the `DeleteDnsRecords` API with zone_id and record_id.

#### Scenario: Successful DNS record deletion
- **WHEN** user runs `terraform destroy` or removes the resource
- **THEN** provider calls `DeleteDnsRecords` API with zone_id and a list containing the record_id, and removes resource from state upon success

#### Scenario: Resource already deleted
- **WHEN** record_id no longer exists during deletion attempt
- **THEN** provider treats as successful deletion and removes from state

### Requirement: Composite ID with zone_id and record_id

The resource SHALL use a composite ID joining zone_id and record_id with `tccommon.FILED_SP` as separator. In read, update, and delete operations, the ID SHALL be parsed back to extract zone_id and record_id.

#### Scenario: ID set after creation
- **WHEN** resource is successfully created
- **THEN** ID is set to `zone_id#record_id` format

#### Scenario: ID parsed during read
- **WHEN** resource Read function is called
- **THEN** zone_id and record_id are extracted from the composite ID for API calls

#### Scenario: Import with composite ID
- **WHEN** user imports an existing DNS record
- **THEN** user provides `zone_id#record_id` as the import ID, and provider reads the resource state

### Requirement: Service layer abstraction

The resource SHALL use service layer methods in `service_tencentcloud_teo.go` for all API interactions, not direct SDK calls.

#### Scenario: Read operation uses service layer
- **WHEN** resource Read function is called
- **THEN** it invokes `DescribeTeoDnsRecordById()` service method

### Requirement: Retry logic for eventual consistency

The resource SHALL implement retry logic using `resource.Retry` with `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations.

#### Scenario: Query retries on transient errors
- **WHEN** `DescribeDnsRecords` returns transient error
- **THEN** provider retries with exponential backoff up to configured timeout

#### Scenario: Write retries on transient errors
- **WHEN** `CreateDnsRecord`, `ModifyDnsRecords`, or `DeleteDnsRecords` returns transient error
- **THEN** provider retries with exponential backoff and returns error if all retries exhausted

### Requirement: Standard error handling patterns

The resource SHALL use `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()` for error handling and logging.

#### Scenario: Operations log elapsed time
- **WHEN** any CRUD operation executes
- **THEN** elapsed time is logged via `tccommon.LogElapsed()`

#### Scenario: Inconsistent state detection
- **WHEN** resource state becomes inconsistent
- **THEN** `tccommon.InconsistentCheck()` detects and logs the inconsistency

### Requirement: Resource registration in provider

The resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap and declared in `tencentcloud/provider.md`.

#### Scenario: Provider includes new resource
- **WHEN** provider initializes
- **THEN** `tencentcloud_teo_dns_record_21` is available in ResourcesMap

#### Scenario: Documentation lists new resource
- **WHEN** documentation is generated
- **THEN** resource appears in provider.md resource list under TEO section

### Requirement: Unit test coverage with mock

The resource SHALL have unit test files using mock (gomonkey) approach for testing business logic without actual cloud API calls.

#### Scenario: Basic resource lifecycle test
- **WHEN** unit test runs
- **THEN** test successfully validates create, read, update, and delete logic with mocked API responses

#### Scenario: Import test
- **WHEN** resource is imported
- **THEN** resource state is populated correctly from mocked API response

### Requirement: Resource documentation

The resource SHALL have a `.md` documentation file following provider conventions with usage examples.

#### Scenario: Documentation includes basic example
- **WHEN** user views resource documentation
- **THEN** a complete usage example is provided showing DNS record creation with required fields

#### Scenario: Documentation includes import section
- **WHEN** user views resource documentation
- **THEN** import example shows the composite ID format (`zone_id#record_id`)
