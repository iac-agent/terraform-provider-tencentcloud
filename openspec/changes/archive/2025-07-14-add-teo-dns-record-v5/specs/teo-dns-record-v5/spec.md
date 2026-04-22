## ADDED Requirements

### Requirement: Create TEO DNS record V5 resource
The system SHALL provide a `tencentcloud_teo_dns_record_v5` Terraform resource that creates, reads, updates, and deletes TEO DNS records using the EdgeOne cloud API.

#### Scenario: Create DNS record with required fields
- **WHEN** user creates a `tencentcloud_teo_dns_record_v5` resource with zone_id, name, type, and content
- **THEN** the system SHALL call `CreateDnsRecord` API with the provided parameters and set the resource ID to `zoneId + FILED_SP + recordId`

#### Scenario: Create DNS record with all optional fields
- **WHEN** user creates a `tencentcloud_teo_dns_record_v5` resource with zone_id, name, type, content, location, ttl, weight, and priority
- **THEN** the system SHALL call `CreateDnsRecord` API with all provided parameters

### Requirement: Read TEO DNS record V5 resource
The system SHALL read the TEO DNS record state by calling `DescribeDnsRecords` API filtered by zone_id and record_id.

#### Scenario: Read existing DNS record
- **WHEN** the system reads a `tencentcloud_teo_dns_record_v5` resource
- **THEN** the system SHALL call `DescribeTeoDnsRecordById` with zone_id and record_id parsed from the composite ID, and populate all schema fields from the API response

#### Scenario: Read deleted DNS record
- **WHEN** the system reads a `tencentcloud_teo_dns_record_v5` resource that no longer exists
- **THEN** the system SHALL set the resource ID to empty string to mark it as deleted

### Requirement: Update TEO DNS record V5 resource
The system SHALL support updating mutable DNS record fields (name, type, content, location, ttl, weight, priority) via `ModifyDnsRecords` API, and updating the status field via `ModifyDnsRecordsStatus` API.

#### Scenario: Update mutable fields
- **WHEN** user updates any of name, type, content, location, ttl, weight, or priority fields
- **THEN** the system SHALL call `ModifyDnsRecords` API with the zone_id and a DnsRecord object containing the record_id and changed fields

#### Scenario: Update status field
- **WHEN** user updates the status field to "enable"
- **THEN** the system SHALL call `ModifyDnsRecordsStatus` API with RecordsToEnable containing the record_id

#### Scenario: Update status field to disable
- **WHEN** user updates the status field to "disable"
- **THEN** the system SHALL call `ModifyDnsRecordsStatus` API with RecordsToDisable containing the record_id

### Requirement: Delete TEO DNS record V5 resource
The system SHALL delete the TEO DNS record by calling `DeleteDnsRecords` API.

#### Scenario: Delete DNS record
- **WHEN** user deletes a `tencentcloud_teo_dns_record_v5` resource
- **THEN** the system SHALL call `DeleteDnsRecords` API with zone_id and record_ids containing the record_id

### Requirement: Import TEO DNS record V5 resource
The system SHALL support importing existing DNS records via Terraform import.

#### Scenario: Import existing DNS record
- **WHEN** user imports a `tencentcloud_teo_dns_record_v5` resource with ID format `zoneId + FILED_SP + recordId`
- **THEN** the system SHALL parse the ID and call the Read function to populate the resource state

### Requirement: Schema definition for TEO DNS record V5
The resource schema SHALL define the following fields:
- `zone_id`: TypeString, Required, ForceNew
- `name`: TypeString, Required
- `type`: TypeString, Required
- `content`: TypeString, Required
- `location`: TypeString, Optional+Computed
- `ttl`: TypeInt, Optional+Computed
- `weight`: TypeInt, Optional+Computed
- `priority`: TypeInt, Optional+Computed
- `status`: TypeString, Optional+Computed
- `record_id`: TypeString, Computed
- `created_on`: TypeString, Computed
- `modified_on`: TypeString, Computed

#### Scenario: Schema field types and constraints
- **WHEN** the resource schema is defined
- **THEN** all field types, required/optional/computed flags, and ForceNew constraints SHALL match the above specification

### Requirement: Resource registration
The system SHALL register `tencentcloud_teo_dns_record_v5` resource in `provider.go` and `provider.md`.

#### Scenario: Register resource in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_dns_record_v5` SHALL be available as a registered resource mapping to `teo.ResourceTencentCloudTeoDnsRecordV5()`

### Requirement: API retry handling
All cloud API calls SHALL use `resource.Retry` with appropriate timeout constants (`tccommon.WriteRetryTimeout` for write operations, `tccommon.ReadRetryTimeout` for read operations). Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: API call with retry
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the system SHALL retry the call within the timeout period

### Requirement: Error handling with defer
All CRUD functions SHALL include `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()` for proper logging and state consistency.

#### Scenario: CRUD function error handling
- **WHEN** a CRUD function is executed
- **THEN** it SHALL include the standard defer statements for logging and consistency checking
