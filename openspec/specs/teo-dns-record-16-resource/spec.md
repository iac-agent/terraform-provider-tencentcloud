# teo-dns-record-16-resource Specification

## Purpose
TBD - created by archiving change add-teo-dns-record-16. Update Purpose after archive.
## Requirements
### Requirement: Resource schema defines TEO DNS record fields
The resource SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): Zone ID identifying the TEO zone
- `name` (TypeString, Required): DNS record name; Chinese/Korean/Japanese domain names must be converted to punycode
- `type` (TypeString, Required): DNS record type, valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- `content` (TypeString, Required): DNS record content corresponding to the type
- `location` (TypeString, Optional, Computed): DNS record resolution route, defaults to "Default"
- `ttl` (TypeInt, Optional, Computed): Cache time in seconds, range 60-86400, default 300
- `weight` (TypeInt, Optional, Computed): DNS record weight, range -1 to 100, default -1
- `priority` (TypeInt, Optional, Computed): MX record priority, range 0-50, default 0, only effective when type is MX
- `status` (TypeString, Computed): DNS record resolution status (enable/disable), computed-only
- `created_on` (TypeString, Computed): Creation time
- `modified_on` (TypeString, Computed): Last modification time

#### Scenario: Create DNS record with required fields only
- **WHEN** a user creates a tencentcloud_teo_dns_record_16 resource with zone_id, name, type, and content
- **THEN** the resource SHALL be created with default values for location, ttl, weight, and priority

#### Scenario: Create DNS record with all optional fields
- **WHEN** a user creates a resource with all optional fields specified (location, ttl, weight, priority)
- **THEN** the resource SHALL be created with the specified values

### Requirement: Resource uses composite ID format
The resource SHALL use a composite ID format of `zone_id + FILED_SP + record_id` to uniquely identify a DNS record. The composite ID SHALL be set after successful creation and parsed during Read, Update, and Delete operations.

#### Scenario: Composite ID is set after creation
- **WHEN** CreateDnsRecord API returns a RecordId
- **THEN** the resource ID SHALL be set to `zone_id + FILED_SP + record_id`

#### Scenario: Composite ID is parsed for Read operation
- **WHEN** the Read function is called
- **THEN** the zone_id and record_id SHALL be extracted from the composite ID by splitting on FILED_SP

### Requirement: Create operation uses CreateDnsRecord API
The resource Create operation SHALL call the `CreateDnsRecord` API with zone_id, name, type, content, location, ttl, weight, and priority. After successful creation, the resource SHALL set the composite ID and call Read to refresh state.

#### Scenario: Successful DNS record creation
- **WHEN** CreateDnsRecord API is called and returns a RecordId
- **THEN** the composite ID SHALL be set and the Read function SHALL be called to refresh the state

#### Scenario: CreateDnsRecord API returns empty RecordId
- **WHEN** CreateDnsRecord API returns an empty or nil RecordId
- **THEN** a NonRetryableError SHALL be returned indicating creation failed

### Requirement: Read operation uses DescribeDnsRecords API with filtering
The resource Read operation SHALL call the `DescribeDnsRecords` API with ZoneId and an AdvancedFilter on "id" field matching the record_id. If the record is not found, the resource SHALL set the ID to empty string to signal removal.

#### Scenario: Record found during Read
- **WHEN** DescribeDnsRecords returns the target record
- **THEN** all schema fields SHALL be populated from the response, checking each field for nil before setting

#### Scenario: Record not found during Read
- **WHEN** DescribeDnsRecords returns no matching record (empty DnsRecords list or nil response)
- **THEN** the resource ID SHALL be set to empty string and a warning SHALL be logged

#### Scenario: Multiple records returned during Read
- **WHEN** DescribeDnsRecords returns multiple records matching the filter
- **THEN** the first record with matching RecordId SHALL be used

### Requirement: Update operation uses ModifyDnsRecords API
The resource Update operation SHALL detect changes in mutable fields (name, type, content, location, ttl, weight, priority) and call `ModifyDnsRecords` API with a DnsRecord struct containing RecordId and the changed fields. The status field SHALL NOT be updated via Terraform.

#### Scenario: Mutable field changed
- **WHEN** any of the mutable fields (name, type, content, location, ttl, weight, priority) have changed
- **THEN** ModifyDnsRecords SHALL be called with a DnsRecord containing RecordId and the updated field values

#### Scenario: No mutable fields changed
- **WHEN** none of the mutable fields have changed
- **THEN** ModifyDnsRecords SHALL NOT be called and the Read function SHALL be called to refresh state

### Requirement: Delete operation uses DeleteDnsRecords API
The resource Delete operation SHALL call the `DeleteDnsRecords` API with ZoneId and RecordIds containing the single record_id extracted from the composite ID.

#### Scenario: Successful DNS record deletion
- **WHEN** DeleteDnsRecords API is called successfully
- **THEN** the resource SHALL be marked as deleted

### Requirement: Resource supports import
The resource SHALL support Terraform import using the composite ID format `zone_id + FILED_SP + record_id`. After import, the Read function SHALL be called to populate the resource state.

#### Scenario: Import with valid composite ID
- **WHEN** a user imports the resource with a valid composite ID
- **THEN** the resource state SHALL be populated by calling the Read function

### Requirement: Resource is registered in provider
The resource SHALL be registered in `provider.go` with key `"tencentcloud_teo_dns_record_16"` mapping to `teo.ResourceTencentCloudTeoDnsRecord16()`, and listed in `provider.md`.

#### Scenario: Resource available in Terraform
- **WHEN** the provider is loaded
- **THEN** the resource `tencentcloud_teo_dns_record_16` SHALL be available for use in Terraform configurations

### Requirement: Unit tests use gomonkey mock approach
The resource SHALL have unit tests in a `_test.go` file using gomonkey to mock the cloud API calls. Tests SHALL cover Create, Read, Update, and Delete operations.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** Create, Read, Update, and Delete operations SHALL be tested using mocked API responses
