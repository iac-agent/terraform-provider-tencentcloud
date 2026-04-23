## Requirements

### Requirement: Teo DNS Record V3 Resource Schema
The system SHALL provide a Terraform resource `tencentcloud_teo_dns_record_v3` with the following schema fields:

- `zone_id` (TypeString, Required, ForceNew): Zone ID identifying the TEO site
- `name` (TypeString, Required): DNS record name; CJK domain names MUST be punycode-encoded
- `type` (TypeString, Required): DNS record type; valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- `content` (TypeString, Required): DNS record content matching the type; CJK domain names MUST be punycode-encoded
- `location` (TypeString, Optional, Computed): Resolution line; defaults to "Default"; only for A/AAAA/CNAME types; only for Standard/Enterprise plans
- `ttl` (TypeInt, Optional, Computed): Cache time in seconds; range 60-86400; default 300
- `weight` (TypeInt, Optional, Computed): Record weight; range -1 to 100; -1 = no weight (default); 0 = no resolution; only for A/AAAA/CNAME types
- `priority` (TypeInt, Optional, Computed): MX record priority; range 0-50; lower = higher priority; default 0; only for MX type
- `record_id` (TypeString, Computed): DNS record ID returned by the cloud API
- `status` (TypeString, Computed): DNS record status; values: "enable" or "disable"
- `created_on` (TypeString, Computed): Creation time
- `modified_on` (TypeString, Computed): Last modification time

#### Scenario: Resource schema defines all Create API fields as input
- **WHEN** a user defines a `tencentcloud_teo_dns_record_v3` resource in Terraform configuration
- **THEN** the schema SHALL accept `zone_id` (Required, ForceNew), `name` (Required), `type` (Required), `content` (Required), `location` (Optional, Computed), `ttl` (Optional, Computed), `weight` (Optional, Computed), and `priority` (Optional, Computed)

#### Scenario: Computed fields are read from API response
- **WHEN** the resource is read from the cloud API
- **THEN** the schema SHALL populate `record_id`, `status`, `created_on`, and `modified_on` as Computed fields

### Requirement: Teo DNS Record V3 Create Operation
The system SHALL create a DNS record using the `CreateDnsRecord` cloud API when the resource is created.

#### Scenario: Successful DNS record creation
- **WHEN** a user creates a `tencentcloud_teo_dns_record_v3` resource with required fields `zone_id`, `name`, `type`, and `content`
- **THEN** the system SHALL call `CreateDnsRecord` API with all populated fields
- **AND** the system SHALL set the resource ID to `zone_id#record_id` using the `RecordId` from the response
- **AND** the system SHALL call Read to refresh the state

#### Scenario: Create with optional fields
- **WHEN** a user creates a `tencentcloud_teo_dns_record_v3` resource including optional fields `location`, `ttl`, `weight`, or `priority`
- **THEN** the system SHALL include those fields in the `CreateDnsRecord` request

#### Scenario: Create API retry on transient failure
- **WHEN** the `CreateDnsRecord` API call fails with a transient error
- **THEN** the system SHALL retry the request using `WriteRetryTimeout`

### Requirement: Teo DNS Record V3 Read Operation
The system SHALL read a DNS record using the `DescribeDnsRecords` cloud API filtered by `record_id`.

#### Scenario: Successful DNS record read
- **WHEN** the system reads a `tencentcloud_teo_dns_record_v3` resource
- **THEN** the system SHALL parse the composite ID `zone_id#record_id`
- **AND** the system SHALL call `DescribeDnsRecords` with `ZoneId` and `AdvancedFilter` on `id` with the `record_id`
- **AND** the system SHALL populate all schema fields from the first matching `DnsRecord` in the response

#### Scenario: DNS record not found
- **WHEN** the `DescribeDnsRecords` response contains no matching records
- **THEN** the system SHALL set `d.SetId("")` to mark the resource as deleted
- **AND** the system SHALL log a warning message

#### Scenario: Read API retry on transient failure
- **WHEN** the `DescribeDnsRecords` API call fails with a transient error
- **THEN** the system SHALL retry the request using `ReadRetryTimeout`

### Requirement: Teo DNS Record V3 Update Operation
The system SHALL update a DNS record using the `ModifyDnsRecords` cloud API when mutable fields change.

#### Scenario: Successful DNS record update
- **WHEN** a user updates any of the mutable fields: `name`, `type`, `content`, `location`, `ttl`, `weight`, or `priority`
- **THEN** the system SHALL detect changes using `d.HasChange()`
- **AND** the system SHALL call `ModifyDnsRecords` API with `ZoneId` and a `DnsRecord` struct containing `RecordId` and all changed fields
- **AND** the system SHALL call Read to refresh the state

#### Scenario: No changes detected
- **WHEN** no mutable fields have changed
- **THEN** the system SHALL skip the `ModifyDnsRecords` API call
- **AND** the system SHALL still call Read to refresh the state

#### Scenario: Update API retry on transient failure
- **WHEN** the `ModifyDnsRecords` API call fails with a transient error
- **THEN** the system SHALL retry the request using `WriteRetryTimeout`

### Requirement: Teo DNS Record V3 Delete Operation
The system SHALL delete a DNS record using the `DeleteDnsRecords` cloud API.

#### Scenario: Successful DNS record deletion
- **WHEN** a user destroys a `tencentcloud_teo_dns_record_v3` resource
- **THEN** the system SHALL parse the composite ID `zone_id#record_id`
- **AND** the system SHALL call `DeleteDnsRecords` API with `ZoneId` and `RecordIds` containing the `record_id`
- **AND** the system SHALL return nil on success

#### Scenario: Delete API retry on transient failure
- **WHEN** the `DeleteDnsRecords` API call fails with a transient error
- **THEN** the system SHALL retry the request using `WriteRetryTimeout`

### Requirement: Teo DNS Record V3 Import
The system SHALL support importing an existing DNS record via `terraform import`.

#### Scenario: Import existing DNS record
- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_v3.example zone_id#record_id`
- **THEN** the system SHALL parse the import ID as `zone_id#record_id`
- **AND** the system SHALL call Read to populate the resource state

### Requirement: Teo DNS Record V3 Resource Registration
The system SHALL register the `tencentcloud_teo_dns_record_v3` resource in `provider.go`.

#### Scenario: Resource registered in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the `tencentcloud_teo_dns_record_v3` resource SHALL be available in the `ResourcesMap` of `provider.go`

### Requirement: Teo DNS Record V3 Service Layer
The system SHALL add a `DescribeTeoDnsRecordV3ById` method to `TeoService` in `service_tencentcloud_teo.go`.

#### Scenario: Service method queries single DNS record
- **WHEN** `DescribeTeoDnsRecordV3ById(ctx, zoneId, recordId)` is called
- **THEN** the system SHALL call `DescribeDnsRecords` with `ZoneId` and filter by `id`
- **AND** the system SHALL return the first matching `DnsRecord` or nil if not found

### Requirement: Teo DNS Record V3 Unit Tests
The system SHALL provide unit tests using gomonkey mock approach for the DNS record v3 resource.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** unit tests are executed
- **THEN** the tests SHALL mock the cloud API calls using gomonkey
- **AND** the tests SHALL verify Create, Read, Update, and Delete operations

### Requirement: Teo DNS Record V3 Documentation
The system SHALL provide a `.md` documentation file for the resource.

#### Scenario: Documentation file exists
- **WHEN** the resource is built
- **THEN** a `resource_tc_teo_dns_record_v3.md` file SHALL exist in the `tencentcloud/services/teo/` directory
- **AND** the file SHALL contain a description, example usage, and import section
