## ADDED Requirements

### Requirement: Create DNS Record
The system SHALL allow creating a TEO DNS record by calling `CreateDnsRecord` API with the following parameters:
- `zone_id` (Required, ForceNew): Zone ID
- `name` (Required): DNS record name
- `type` (Required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (Required): DNS record content
- `location` (Optional): DNS record resolution route, default "Default"
- `ttl` (Optional): Cache time, range 60-86400, default 300
- `weight` (Optional): DNS record weight, range -1~100, default -1
- `priority` (Optional): MX record priority, range 0-50, default 0

Upon successful creation, the system SHALL set the resource ID to `zone_id#record_id` using `tccommon.FILED_SP` as separator, and store the returned `record_id`.

#### Scenario: Successful DNS record creation
- **WHEN** user creates a `tencentcloud_teo_dns_record_20` resource with zone_id, name, type, and content
- **THEN** the system calls `CreateDnsRecord` API and sets the resource ID to `zone_id#record_id`

#### Scenario: Create DNS record with all optional parameters
- **WHEN** user creates a `tencentcloud_teo_dns_record_20` resource with all parameters including location, ttl, weight, and priority
- **THEN** the system calls `CreateDnsRecord` API with all specified parameters

#### Scenario: Create DNS record returns empty response
- **WHEN** `CreateDnsRecord` API returns a nil or empty RecordId
- **THEN** the system returns a NonRetryableError indicating the creation failed

### Requirement: Read DNS Record
The system SHALL allow reading a TEO DNS record by calling `DescribeDnsRecords` API with a filter on record_id.

When the resource ID is in the format `zone_id#record_id`, the system SHALL:
1. Parse the ID to extract zone_id and record_id
2. Call `DescribeDnsRecords` with ZoneId and Filter(id=recordId)
3. Find the matching record from the response

Computed fields that SHALL be set from the response (only if not nil):
- `record_id`: DNS record ID
- `status`: DNS record status (enable/disable)
- `created_on`: Creation time
- `modified_on`: Modification time

#### Scenario: Successful DNS record read
- **WHEN** user reads a `tencentcloud_teo_dns_record_20` resource
- **THEN** the system parses the ID, queries by zone_id and record_id, and populates all fields

#### Scenario: DNS record not found
- **WHEN** the DescribeDnsRecords response does not contain the expected record
- **THEN** the system sets the resource ID to empty string and logs a warning

#### Scenario: Invalid resource ID format
- **WHEN** the resource ID cannot be split into exactly 2 parts by FILED_SP
- **THEN** the system returns an error indicating the ID is broken

### Requirement: Update DNS Record
The system SHALL support updating a TEO DNS record with two separate API calls:

1. When `name`, `type`, `content`, `location`, `ttl`, `weight`, or `priority` changes, call `ModifyDnsRecords` API with:
   - `zone_id`: Zone ID
   - `dns_records`: Array containing one DnsRecord with RecordId and modified fields

2. When `status` changes, call `ModifyDnsRecordsStatus` API with:
   - `zone_id`: Zone ID
   - If status is "enable": `records_to_enable` containing the record_id
   - If status is "disable": `records_to_disable` containing the record_id

The system SHALL NOT include ZoneId, Status, CreatedOn, or ModifiedOn in the DnsRecord struct for ModifyDnsRecords as these are output-only fields.

#### Scenario: Update DNS record content
- **WHEN** user updates the `content` field of a `tencentcloud_teo_dns_record_20` resource
- **THEN** the system calls `ModifyDnsRecords` with the updated content

#### Scenario: Update DNS record status to enable
- **WHEN** user changes the `status` field to "enable"
- **THEN** the system calls `ModifyDnsRecordsStatus` with the record_id in `records_to_enable`

#### Scenario: Update DNS record status to disable
- **WHEN** user changes the `status` field to "disable"
- **THEN** the system calls `ModifyDnsRecordsStatus` with the record_id in `records_to_disable`

#### Scenario: No changes detected
- **WHEN** user runs terraform apply with no actual changes to the resource
- **THEN** the system skips all API calls and returns directly

### Requirement: Delete DNS Record
The system SHALL allow deleting a TEO DNS record by calling `DeleteDnsRecords` API with:
- `zone_id`: Zone ID
- `record_ids`: Array containing the record_id to delete

#### Scenario: Successful DNS record deletion
- **WHEN** user deletes a `tencentcloud_teo_dns_record_20` resource
- **THEN** the system calls `DeleteDnsRecords` with the zone_id and record_id

### Requirement: Import DNS Record
The system SHALL support importing existing TEO DNS records using `schema.ImportStatePassthrough`.

The import ID format SHALL be `zone_id#record_id` using `tccommon.FILED_SP` as separator.

#### Scenario: Import existing DNS record
- **WHEN** user runs `terraform import tencentcloud_teo_dns_record_20.example zone_id#record_id`
- **THEN** the system imports the resource and reads its current state

### Requirement: Resource Registration
The system SHALL register `tencentcloud_teo_dns_record_20` in `provider.go` and `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is loaded
- **THEN** `tencentcloud_teo_dns_record_20` is available as a resource type

### Requirement: Service Layer Query
The system SHALL add a `DescribeTeoDnsRecord20ById` method to the TeoService that queries a specific DNS record by zone_id and record_id using `DescribeDnsRecords` API with Filter(id=recordId).

#### Scenario: Query DNS record by ID
- **WHEN** the service method is called with zone_id and record_id
- **THEN** it calls DescribeDnsRecords with ZoneId and Filter(id=[recordId]), and returns the first matching DnsRecord

### Requirement: Unit Tests
The system SHALL provide unit tests for `tencentcloud_teo_dns_record_20` using gomonkey mock approach (not Terraform acceptance test suite).

#### Scenario: Unit tests cover CRUD operations
- **WHEN** the test suite runs
- **THEN** it tests create, read, update, and delete operations using mocked API calls

### Requirement: Documentation
The system SHALL provide a `.md` documentation file for `tencentcloud_teo_dns_record_20` including:
- One-sentence description mentioning TEO product name
- Example Usage section
- Import section (since this is RESOURCE_KIND_GENERAL)

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a corresponding .md file with proper format exists at `tencentcloud/services/teo/resource_tc_teo_dns_record_20.md`
