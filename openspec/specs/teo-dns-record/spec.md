# teo-dns-record Specification

## Purpose
TBD - created by archiving change add-teo-dns-record. Update Purpose after archive.
## Requirements
### Requirement: Create TEO DNS record
The resource SHALL create a DNS record in a TEO zone using the `CreateDnsRecord` API. The resource MUST accept `zone_id`, `name`, `type`, `content` as required fields, and `location`, `ttl`, `weight`, `priority` as optional fields. Upon successful creation, the resource SHALL store the returned `RecordId` as `record_id` in state.

#### Scenario: Successful DNS record creation
- **WHEN** a user defines a `tencentcloud_teo_dns_record_24` resource with valid `zone_id`, `name`, `type`, and `content`
- **THEN** the provider SHALL call `CreateDnsRecord` and set the resource ID to `zone_id#record_id`

#### Scenario: Creation with optional fields
- **WHEN** a user specifies `ttl`, `weight`, `priority`, or `location` in the resource configuration
- **THEN** the provider SHALL pass these values to the `CreateDnsRecord` API call

### Requirement: Read TEO DNS record
The resource SHALL read the current state of a DNS record using `DescribeDnsRecords` with an `id` filter matching the `record_id`. The resource MUST update all schema fields from the API response, including computed fields `status`, `created_on`, `modified_on`.

#### Scenario: Successful DNS record read
- **WHEN** the provider reads a `tencentcloud_teo_dns_record_24` resource
- **THEN** it SHALL call `DescribeDnsRecords` with `ZoneId` and a filter `id=<record_id>` and populate all fields from the returned `DnsRecord`

#### Scenario: DNS record not found
- **WHEN** `DescribeDnsRecords` returns an empty list for the given `record_id`
- **THEN** the provider SHALL remove the resource from state (set `d.SetId("")`)

### Requirement: Update TEO DNS record
The resource SHALL update a DNS record using the `ModifyDnsRecords` API, passing a single-element `DnsRecords` list containing the updated record fields. The `RecordId` MUST be included in the `DnsRecord` object.

#### Scenario: Successful DNS record update
- **WHEN** a user modifies `content`, `ttl`, `weight`, `priority`, or `location` of an existing DNS record
- **THEN** the provider SHALL call `ModifyDnsRecords` with the updated `DnsRecord` object and then call Read to refresh state

#### Scenario: Immutable field change triggers recreation
- **WHEN** a user changes `zone_id`, `name`, or `type`
- **THEN** Terraform SHALL destroy and recreate the resource (ForceNew behavior)

### Requirement: Delete TEO DNS record
The resource SHALL delete a DNS record using the `DeleteDnsRecords` API, passing the `record_id` in the `RecordIds` list.

#### Scenario: Successful DNS record deletion
- **WHEN** a user runs `terraform destroy` on a `tencentcloud_teo_dns_record_24` resource
- **THEN** the provider SHALL call `DeleteDnsRecords` with `ZoneId` and `RecordIds=[record_id]`

### Requirement: Import TEO DNS record
The resource SHALL support import using the composite ID format `zone_id#record_id`.

#### Scenario: Successful import
- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_24.<name> <zone_id>#<record_id>`
- **THEN** the provider SHALL parse the composite ID and read the resource state from the API

### Requirement: Resource schema fields
The resource schema SHALL define the following fields:
- `zone_id` (Required, ForceNew, string): TEO zone ID
- `name` (Required, ForceNew, string): DNS record name
- `type` (Required, ForceNew, string): DNS record type (A/AAAA/MX/CNAME/TXT/NS/CAA/SRV)
- `content` (Required, string): DNS record content
- `location` (Optional, string): DNS resolution line, default "Default"
- `ttl` (Optional, int): Cache TTL in seconds (60-86400), default 300
- `weight` (Optional, int): DNS record weight (-1 to 100)
- `priority` (Optional, int): MX record priority (0-50)
- `record_id` (Computed, string): DNS record ID returned by the API
- `status` (Computed, string): DNS record status (enable/disable)
- `created_on` (Computed, string): Record creation time
- `modified_on` (Computed, string): Record last modification time

#### Scenario: Schema validation
- **WHEN** a user provides all required fields
- **THEN** the provider SHALL accept the configuration without validation errors

### Requirement: API retry handling
All cloud API calls SHALL use `resource.RetryContext` with `tccommon.ReadRetryTimeout` as the timeout. Failed API calls SHALL be wrapped with `tccommon.RetryError`.

#### Scenario: Transient API error retry
- **WHEN** a cloud API call returns a retryable error
- **THEN** the provider SHALL retry the call until success or timeout

### Requirement: Unit tests with gomonkey mocks
The resource SHALL have unit tests in `resource_tc_teo_dns_record_24_test.go` that use gomonkey to mock cloud API calls. Tests SHALL cover Create, Read, Update, and Delete operations.

#### Scenario: Unit test execution
- **WHEN** `go test -gcflags=all=-l` is run on the test file
- **THEN** all unit tests SHALL pass without requiring real cloud credentials

