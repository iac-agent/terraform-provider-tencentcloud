## ADDED Requirements

### Requirement: Resource supports CRUD operations for TEO DNS records
The system SHALL provide a Terraform resource `tencentcloud_teo_dns_record_25` that supports creating, reading, updating, deleting, and importing DNS records within a TEO zone.

#### Scenario: Create a DNS record
- **WHEN** a user applies a Terraform configuration with `tencentcloud_teo_dns_record_25` specifying `zone_id`, `name`, `type`, and `content`
- **THEN** the system calls `CreateDnsRecord` API, stores the returned `RecordId`, and sets the composite resource ID as `zone_id#record_id`

#### Scenario: Read a DNS record
- **WHEN** Terraform reads the state of an existing `tencentcloud_teo_dns_record_25` resource
- **THEN** the system calls `DescribeDnsRecords` with a filter on the record ID, populates all schema fields from the API response, and if the record is not found, clears the resource ID

#### Scenario: Update a DNS record
- **WHEN** a user modifies any mutable field (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) of an existing `tencentcloud_teo_dns_record_25` resource
- **THEN** the system calls `ModifyDnsRecords` with the updated fields and refreshes the state

#### Scenario: Delete a DNS record
- **WHEN** a user destroys a `tencentcloud_teo_dns_record_25` resource
- **THEN** the system calls `DeleteDnsRecords` with the record ID and removes the resource from state

#### Scenario: Import a DNS record
- **WHEN** a user imports a DNS record using `terraform import tencentcloud_teo_dns_record_25.foo zone_id#record_id`
- **THEN** the system parses the composite ID and reads the record state from the cloud API

### Requirement: Resource schema includes core parameters
The resource schema SHALL include the following parameters:
- `zone_id` (Required, ForceNew, TypeString) — the TEO zone ID
- `name` (Required, TypeString) — the DNS record name
- `type` (Required, TypeString) — the DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (Required, TypeString) — the DNS record content
- `location` (Optional, Computed, TypeString) — the DNS record resolution route, defaults to DEFAULT
- `ttl` (Optional, Computed, TypeInt) — cache time in seconds, range 60-86400, default 300
- `weight` (Optional, Computed, TypeInt) — DNS record weight, range -1 to 100, default -1
- `priority` (Optional, Computed, TypeInt) — MX record priority, range 0-50, default 0
- `record_id` (Computed, TypeString) — the DNS record ID returned by the API

#### Scenario: Required fields are correctly enforced
- **WHEN** a user omits `zone_id`, `name`, `type`, or `content` from the configuration
- **THEN** Terraform returns a validation error requiring these fields

#### Scenario: Optional fields use defaults
- **WHEN** a user does not specify `location`, `ttl`, `weight`, or `priority`
- **THEN** the system uses the API defaults (DEFAULT for location, 300 for ttl, -1 for weight, 0 for priority)

### Requirement: Retry logic is applied to all API calls
The system SHALL use `resource.Retry` with `tccommon.ReadRetryTimeout` (for read) and `tccommon.WriteRetryTimeout` (for create/update/delete) to handle transient API errors.

#### Scenario: API call retries on transient failure
- **WHEN** a cloud API call returns a transient error
- **THEN** the system retries the call within the configured timeout period