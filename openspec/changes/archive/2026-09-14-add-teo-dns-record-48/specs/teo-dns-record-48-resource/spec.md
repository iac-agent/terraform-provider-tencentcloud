## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_dns_record_48` resource SHALL define the following schema fields:
- `zone_id` (string, required, ForceNew): TEO Zone ID
- `name` (string, required): DNS record name
- `type` (string, required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (string, required): DNS record content
- `location` (string, optional): DNS record resolution route (default: "Default")
- `ttl` (int, optional): Cache TTL in seconds (range 60-86400, default: 300)
- `weight` (int, optional): DNS record weight (range -1 to 100, default: -1)
- `priority` (int, optional): MX record priority (range 0-50, default: 0)
- `record_id` (string, computed): DNS record ID returned by CreateDnsRecord
- `status` (string, computed): DNS record status
- `created_on` (string, computed): Creation time
- `modified_on` (string, computed): Last modification time

#### Scenario: Resource schema has required and optional fields
- **WHEN** the resource schema is defined
- **THEN** `zone_id`, `name`, `type`, `content` are required fields; `location`, `ttl`, `weight`, `priority` are optional fields; `record_id`, `status`, `created_on`, `modified_on` are computed fields

#### Scenario: zone_id is ForceNew
- **WHEN** `zone_id` is changed after resource creation
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: Resource Create
The resource SHALL call `CreateDnsRecord` API to create a DNS record, using `zone_id`, `name`, `type`, `content`, and optional `location`, `ttl`, `weight`, `priority` as input parameters. After creation, the resource SHALL set the composite ID to `zone_id#record_id` using `tccommon.FILED_SP` as separator.

#### Scenario: Successful DNS record creation
- **WHEN** a `tencentcloud_teo_dns_record_48` resource is created with valid parameters
- **THEN** the system SHALL call `CreateDnsRecord` with all provided parameters
- **AND** set `record_id` from `response.Response.RecordId`
- **AND** set resource ID to `zone_id#record_id`

#### Scenario: CreateDnsRecord returns empty RecordId
- **WHEN** `CreateDnsRecord` response contains empty `RecordId`
- **THEN** the system SHALL return `NonRetryableError` to prevent writing empty ID to state

### Requirement: Resource Read
The resource SHALL call `DescribeDnsRecords` API with `ZoneId` and filter by `record-id` to read the current state of a DNS record. The resource SHALL use pagination with Limit=1000.

#### Scenario: Successful DNS record read
- **WHEN** the resource Read function is called
- **THEN** the system SHALL call `DescribeDnsRecords` with `ZoneId` and filter by `record-id`
- **AND** set all schema fields from the matched DNS record

#### Scenario: DNS record not found
- **WHEN** `DescribeDnsRecords` returns empty list or the record is not found
- **THEN** the system SHALL log `[CRUD] teo_dns_record_48 id=%s` with the current ID
- **AND** call `d.SetId("")` to signal resource removal from state

### Requirement: Resource Update
The resource SHALL call `ModifyDnsRecords` API to update a DNS record when updatable fields change. The system SHALL pass a single-element `DnsRecords` array containing `RecordId` and modified fields. Read-only fields (`Status`, `CreatedOn`, `ModifiedOn`) SHALL NOT be included in the update request.

#### Scenario: Successful DNS record update
- **WHEN** updatable fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) change
- **THEN** the system SHALL call `ModifyDnsRecords` with `ZoneId` and a single-element `DnsRecords` array containing `RecordId` and all current field values

### Requirement: Resource Delete
The resource SHALL call `DeleteDnsRecords` API with `ZoneId` and `RecordIds` array containing the single record ID to delete a DNS record.

#### Scenario: Successful DNS record deletion
- **WHEN** a `tencentcloud_teo_dns_record_48` resource is destroyed
- **THEN** the system SHALL call `DeleteDnsRecords` with `ZoneId` and `RecordIds` containing the record ID

### Requirement: Resource Import
The resource SHALL support Terraform import using the composite ID format `zone_id#record_id`.

#### Scenario: Import existing DNS record
- **WHEN** `terraform import tencentcloud_teo_dns_record_48.example zone_id#record_id` is executed
- **THEN** the system SHALL parse the composite ID and read the DNS record

### Requirement: Provider Registration
The `tencentcloud_teo_dns_record_48` resource SHALL be registered in `provider.go` and `provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the Terraform provider is initialized
- **THEN** `tencentcloud_teo_dns_record_48` SHALL be available as a resource type

### Requirement: Retry and Error Handling
All cloud API calls SHALL use `tccommon.ReadRetryTimeout` as timeout and implement retry logic with `helper.Retry()`. Failed API calls SHALL return `tccommon.RetryError()` wrapped errors.

#### Scenario: API call with retry
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the system SHALL retry the call within the timeout period
- **AND** return `tccommon.RetryError()` wrapped error on failure
