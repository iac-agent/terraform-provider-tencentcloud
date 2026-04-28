## ADDED Requirements

### Requirement: Create DNS Record
The system SHALL provide a `tencentcloud_teo_dns_record_v8` resource that creates a TEO DNS record using the `CreateDnsRecord` API. The resource SHALL accept the following required parameters: `zone_id`, `name`, `type`, `content`. It SHALL also accept optional parameters: `location`, `ttl`, `weight`, `priority`. Upon successful creation, the system SHALL set the resource ID to the composite format `zoneId#recordId`.

#### Scenario: Create DNS record with required fields only
- **WHEN** user creates a `tencentcloud_teo_dns_record_v8` resource with `zone_id`, `name`, `type`, and `content`
- **THEN** the system SHALL call `CreateDnsRecord` API with the provided parameters and set the resource ID to `zoneId#recordId`

#### Scenario: Create DNS record with all optional fields
- **WHEN** user creates a `tencentcloud_teo_dns_record_v8` resource with all optional fields (`location`, `ttl`, `weight`, `priority`) in addition to required fields
- **THEN** the system SHALL call `CreateDnsRecord` API with all provided parameters

#### Scenario: Create DNS record API failure
- **WHEN** the `CreateDnsRecord` API call fails
- **THEN** the system SHALL retry with `tccommon.WriteRetryTimeout` and wrap the error using `tccommon.RetryError()`

#### Scenario: Create DNS record returns empty RecordId
- **WHEN** the `CreateDnsRecord` API returns an empty `RecordId`
- **THEN** the system SHALL return a `NonRetryableError`

### Requirement: Read DNS Record
The system SHALL read the DNS record state using the `DescribeDnsRecords` API with an `AdvancedFilter` on `"id"` field. The system SHALL parse the composite ID `zoneId#recordId` to identify the resource. If the record is not found, the system SHALL set the resource ID to empty string.

#### Scenario: Read existing DNS record
- **WHEN** user reads a `tencentcloud_teo_dns_record_v8` resource that exists
- **THEN** the system SHALL query `DescribeDnsRecords` with zone_id and id filter, and populate all schema fields from the response

#### Scenario: Read deleted DNS record
- **WHEN** user reads a `tencentcloud_teo_dns_record_v8` resource that no longer exists
- **THEN** the system SHALL set `d.SetId("")` and return without error

#### Scenario: Broken composite ID
- **WHEN** the resource ID cannot be split into exactly 2 parts by the `#` separator
- **THEN** the system SHALL return an error indicating the ID is broken

### Requirement: Update DNS Record Content Fields
The system SHALL update DNS record content fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) using the `ModifyDnsRecords` API when any of these fields have changed. The system SHALL only make the API call if at least one of the mutable fields has changed.

#### Scenario: Update content fields
- **WHEN** user updates any of `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority` fields
- **THEN** the system SHALL call `ModifyDnsRecords` API with the `RecordId` and modified fields

#### Scenario: No content fields changed
- **WHEN** none of the content fields have changed
- **THEN** the system SHALL NOT call `ModifyDnsRecords` API

### Requirement: Update DNS Record Status
The system SHALL update DNS record status using the `ModifyDnsRecordsStatus` API when the `status` field has changed. When `status` is set to `"enable"`, the system SHALL pass the record ID to `RecordsToEnable`. When `status` is set to `"disable"`, the system SHALL pass the record ID to `RecordsToDisable`.

#### Scenario: Enable DNS record
- **WHEN** user sets `status` to `"enable"`
- **THEN** the system SHALL call `ModifyDnsRecordsStatus` with the record ID in `RecordsToEnable`

#### Scenario: Disable DNS record
- **WHEN** user sets `status` to `"disable"`
- **THEN** the system SHALL call `ModifyDnsRecordsStatus` with the record ID in `RecordsToDisable`

#### Scenario: Status unchanged
- **WHEN** the `status` field has not changed
- **THEN** the system SHALL NOT call `ModifyDnsRecordsStatus` API

### Requirement: Delete DNS Record
The system SHALL delete the DNS record using the `DeleteDnsRecords` API with `zone_id` and `record_id` as a list.

#### Scenario: Delete existing DNS record
- **WHEN** user deletes a `tencentcloud_teo_dns_record_v8` resource
- **THEN** the system SHALL call `DeleteDnsRecords` API with `ZoneId` and `RecordIds` containing the record ID

#### Scenario: Delete API failure
- **WHEN** the `DeleteDnsRecords` API call fails
- **THEN** the system SHALL retry with `tccommon.WriteRetryTimeout` and wrap the error using `tccommon.RetryError()`

### Requirement: Import DNS Record
The system SHALL support importing existing DNS records. The import ID format SHALL be `zoneId#recordId`.

#### Scenario: Import existing DNS record
- **WHEN** user imports a `tencentcloud_teo_dns_record_v8` resource using the ID `zoneId#recordId`
- **THEN** the system SHALL parse the ID, read the DNS record, and populate the Terraform state

### Requirement: Resource Registration in Provider
The system SHALL register the `tencentcloud_teo_dns_record_v8` resource in `provider.go` and add the corresponding entry in `provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_teo_dns_record_v8` resource SHALL be available for use in Terraform configurations

### Requirement: Service Layer Describe Method
The system SHALL provide a `DescribeTeoDnsRecordV8ById` method in `service_tencentcloud_teo.go` that queries a DNS record by zone ID and record ID using `DescribeDnsRecords` API with an id filter.

#### Scenario: Query DNS record by ID
- **WHEN** `DescribeTeoDnsRecordV8ById` is called with valid zoneId and recordId
- **THEN** the system SHALL call `DescribeDnsRecords` with zone_id and id filter, and return the matching `DnsRecord` object

### Requirement: Unit Tests with Mock
The system SHALL provide unit tests using gomonkey mock for all CRUD operations of the `tencentcloud_teo_dns_record_v8` resource. Tests SHALL be run with `go test -gcflags=all=-l`.

#### Scenario: Unit test coverage
- **WHEN** unit tests are executed
- **THEN** all CRUD operations (Create, Read, Update, Delete) SHALL be tested with mocked cloud API responses

### Requirement: Resource Documentation
The system SHALL provide a `.md` documentation file for the `tencentcloud_teo_dns_record_v8` resource, including a one-line description mentioning TEO, example usage with HCL, and import instructions with the composite ID format.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_teo_dns_record_v8.md` file SHALL exist with description, example usage, and import section
