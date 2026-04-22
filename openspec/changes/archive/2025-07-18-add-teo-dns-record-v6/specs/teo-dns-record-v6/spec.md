## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_dns_record_v6` resource SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): TEO zone ID
- `name` (TypeString, Required): DNS record name
- `type` (TypeString, Required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (TypeString, Required): DNS record content
- `location` (TypeString, Optional, Computed): DNS record resolution route, default "Default"
- `ttl` (TypeInt, Optional, Computed): Cache time in seconds, range 60-86400, default 300
- `weight` (TypeInt, Optional, Computed): DNS record weight, range -1~100, default -1
- `priority` (TypeInt, Optional, Computed): MX record priority, range 0-50, default 0
- `status` (TypeString, Optional, Computed): DNS record resolution status (enable/disable)
- `record_id` (TypeString, Computed): DNS record ID returned from creation
- `created_on` (TypeString, Computed): Creation time
- `modified_on` (TypeString, Computed): Modification time

The resource SHALL support import via `schema.ImportStatePassthrough`.

#### Scenario: Schema defines all required and computed fields
- **WHEN** the resource schema is initialized
- **THEN** all fields listed above SHALL be present with correct types and required/computed attributes

#### Scenario: zone_id is ForceNew
- **WHEN** the zone_id field is changed in a Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

### Requirement: Create DNS Record
The resource SHALL create a TEO DNS record by calling `CreateDnsRecord` API with zone_id, name, type, content, location, ttl, weight, and priority parameters. Upon successful creation, the resource SHALL set its ID to `zone_id + FILED_SP + record_id` using the RecordId from the response.

#### Scenario: Successful DNS record creation
- **WHEN** a valid DNS record configuration is applied
- **THEN** the CreateDnsRecord API SHALL be called with all configured parameters
- **THEN** the resource ID SHALL be set to zone_id + FILED_SP + record_id
- **THEN** a Read operation SHALL be performed to sync the full resource state

#### Scenario: Create API call fails
- **WHEN** the CreateDnsRecord API returns an error
- **THEN** the error SHALL be wrapped with tccommon.RetryError and retried within WriteRetryTimeout
- **THEN** if retries are exhausted, the error SHALL be returned to the user

### Requirement: Read DNS Record
The resource SHALL read the current state of a DNS record by calling the `DescribeTeoDnsRecordV6ById` service method, which internally calls `DescribeDnsRecords` API with a filter on the record ID. If the record is not found, the resource SHALL set its ID to empty string to indicate removal.

#### Scenario: Record found during read
- **WHEN** the DescribeDnsRecords API returns the record matching the stored record_id
- **THEN** all schema fields SHALL be populated from the API response

#### Scenario: Record not found during read
- **WHEN** the DescribeDnsRecords API returns no matching record
- **THEN** the resource ID SHALL be set to empty string
- **THEN** the resource SHALL be marked as gone

### Requirement: Update DNS Record
The resource SHALL support updating mutable fields (name, type, content, location, ttl, weight, priority) by calling `ModifyDnsRecords` API with a DnsRecord struct containing the RecordId and updated fields. The resource SHALL also support updating the `status` field by calling `ModifyDnsRecordsStatus` API.

#### Scenario: Update mutable fields
- **WHEN** any of name, type, content, location, ttl, weight, or priority fields change
- **THEN** the ModifyDnsRecords API SHALL be called with the updated DnsRecord struct
- **THEN** a Read operation SHALL be performed to sync state after update

#### Scenario: Update status field
- **WHEN** the status field changes to "enable"
- **THEN** the ModifyDnsRecordsStatus API SHALL be called with RecordsToEnable containing the record_id
- **WHEN** the status field changes to "disable"
- **THEN** the ModifyDnsRecordsStatus API SHALL be called with RecordsToDisable containing the record_id

#### Scenario: No fields changed
- **WHEN** no mutable fields have changed
- **THEN** no update API calls SHALL be made

### Requirement: Delete DNS Record
The resource SHALL delete a DNS record by calling `DeleteDnsRecords` API with zone_id and record_ids.

#### Scenario: Successful DNS record deletion
- **WHEN** the resource is destroyed
- **THEN** the DeleteDnsRecords API SHALL be called with ZoneId and RecordIds containing the record_id
- **THEN** the API call SHALL use WriteRetryTimeout for retry handling

#### Scenario: Delete API call fails
- **WHEN** the DeleteDnsRecords API returns an error
- **THEN** the error SHALL be retried within WriteRetryTimeout
- **THEN** if retries are exhausted, the error SHALL be returned to the user

### Requirement: Resource Registration
The resource SHALL be registered in `provider.go` with the key `tencentcloud_teo_dns_record_v6` mapping to `teo.ResourceTencentCloudTeoDnsRecordV6()`, and documented in `provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_teo_dns_record_v6` resource SHALL be available for use in Terraform configurations

### Requirement: Service Layer Method
The service layer in `service_tencentcloud_teo.go` SHALL include a `DescribeTeoDnsRecordV6ById` method that accepts ctx, zoneId, and recordId, calls `DescribeDnsRecords` with a filter on record ID, and returns the matching DnsRecord struct or nil if not found.

#### Scenario: Service method finds record
- **WHEN** the DescribeDnsRecords API returns records matching the filter
- **THEN** the method SHALL return the DnsRecord struct matching the exact record_id

#### Scenario: Service method finds no matching record
- **WHEN** no record matches the given record_id
- **THEN** the method SHALL return nil

### Requirement: Unit Tests
The resource SHALL have unit tests in `resource_tc_teo_dns_record_v6_test.go` using gomonkey mocks for cloud API calls. Tests SHALL cover create, read, update, and delete operations.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** unit tests are executed with `go test`
- **THEN** create, read, update, and delete scenarios SHALL be tested using mocked API responses

### Requirement: Resource Documentation
The resource SHALL have a documentation file `resource_tc_teo_dns_record_v6.md` following the gendoc format, including a one-line description, example usage, and import section.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `.md` file SHALL exist with proper description, example usage with jsonencode() where applicable, and import instructions
