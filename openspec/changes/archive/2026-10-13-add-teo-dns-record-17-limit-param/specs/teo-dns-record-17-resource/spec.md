# teo-dns-record-17-resource Specification

## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_dns_record_17` resource SHALL define a schema with the following fields:
- `zone_id` (string, required, ForceNew): TEO Zone ID
- `name` (string, required): DNS record name
- `type` (string, required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content` (string, required): DNS record content
- `location` (string, optional+computed): DNS record resolution route (default "Default")
- `ttl` (int, optional+computed): Cache TTL in seconds (range 60-86400, default 300)
- `weight` (int, optional+computed): DNS record weight (range -1 to 100, default -1)
- `priority` (int, optional+computed): MX record priority (range 0-50, default 0)
- `status` (string, optional+computed): DNS record status (enable / disable)
- `Limit` (int, optional): Pagination limit for the `DescribeDnsRecords` API call used by Read (cloud API max: 1000; default 1000 when not set)
- `record_id` (string, computed): DNS record ID
- `created_on` (string, computed): Creation time
- `modified_on` (string, computed): Last modification time

#### Scenario: Resource schema has required and optional fields
- **WHEN** the resource schema is defined
- **THEN** `zone_id`, `name`, `type`, `content` SHALL be required fields; `location`, `ttl`, `weight`, `priority`, `status` SHALL be optional+computed fields; `Limit` SHALL be an optional field; `record_id`, `created_on`, `modified_on` SHALL be computed fields

#### Scenario: zone_id is ForceNew
- **WHEN** `zone_id` is changed after resource creation
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: Resource Create
The resource SHALL call `CreateDnsRecord` API to create a DNS record, using `zone_id`, `name`, `type`, `content`, and optional `location`, `ttl`, `weight`, `priority` as input parameters. The `Limit` parameter SHALL NOT be passed to `CreateDnsRecord` (it is a `DescribeDnsRecords`-only query parameter).

#### Scenario: Successful DNS record creation
- **WHEN** a `tencentcloud_teo_dns_record_17` resource is created with valid parameters
- **THEN** the system SHALL call `CreateDnsRecord` with all provided parameters (excluding `Limit`)
- **AND** set resource ID to `zone_id#record_id` using `tccommon.FILED_SP` as separator

#### Scenario: CreateDnsRecord returns empty RecordId
- **WHEN** `CreateDnsRecord` response is nil, or `Response` is nil, or `RecordId` is nil or empty string
- **THEN** the system SHALL log the logId and current state, and return `NonRetryableError` to prevent writing an empty ID to state

### Requirement: Resource Read with Limit Parameter
The resource SHALL call the service layer method `DescribeTeoDnsRecord17ById` which invokes `DescribeDnsRecords` API with `ZoneId`, an `id` AdvancedFilter matching the record ID, `Offset=0`, and `Limit` taken from the resource's `Limit` field. When `Limit` is not configured by the user, the system SHALL default it to 1000 (the maximum value documented in the cloud API).

#### Scenario: Successful DNS record read with default Limit
- **WHEN** the resource Read function is called and `Limit` is not set in the configuration
- **THEN** the system SHALL call `DescribeDnsRecords` with `Limit=1000` and the `id` filter
- **AND** set all schema fields from the first matched DNS record, checking each response field for nil before calling `d.Set()`

#### Scenario: DNS record read with user-configured Limit
- **WHEN** the resource Read function is called and `Limit` is set in the configuration (e.g. 500)
- **THEN** the system SHALL call `DescribeDnsRecords` with `Limit=500`

#### Scenario: DNS record not found
- **WHEN** `DescribeDnsRecords` returns an empty list
- **THEN** the system SHALL log `[CRUD] teo_dns_record_17 id=<id>` with the current composite ID before clearing it
- **AND** call `d.SetId("")` to signal resource removal from state

### Requirement: Resource Update
The resource SHALL call `ModifyDnsRecords` API to update a DNS record when updatable fields change. The system SHALL pass a single-element `DnsRecords` array containing `RecordId` and modified fields (`Name`, `Type`, `Location`, `Content`, `TTL`, `Weight`, `Priority`). Output-only fields (`ZoneId`, `Status`, `CreatedOn`, `ModifiedOn`) SHALL NOT be included in the update request. The `Limit` parameter SHALL NOT trigger any update API call.

#### Scenario: Successful DNS record content update
- **WHEN** updatable fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) change
- **THEN** the system SHALL call `ModifyDnsRecords` with `ZoneId` and a single-element `DnsRecords` array containing `RecordId` and all current updatable field values

#### Scenario: DNS record status update
- **WHEN** the `status` field changes to `enable` or `disable`
- **THEN** the system SHALL call `ModifyDnsRecordsStatus` with `ZoneId` and the record ID in `RecordsToEnable` (for enable) or `RecordsToDisable` (for disable)

#### Scenario: Limit change does not trigger update
- **WHEN** only the `Limit` field changes
- **THEN** the system SHALL NOT call `ModifyDnsRecords` or `ModifyDnsRecordsStatus`, and SHALL only refresh state via Read

### Requirement: Resource Delete
The resource SHALL call `DeleteDnsRecords` API with `ZoneId` and `RecordIds` array containing the single record ID to delete a DNS record.

#### Scenario: Successful DNS record deletion
- **WHEN** a `tencentcloud_teo_dns_record_17` resource is destroyed
- **THEN** the system SHALL parse the composite ID, call `DeleteDnsRecords` with `ZoneId` and `RecordIds` containing the record ID

### Requirement: Composite ID and Import
The resource SHALL use the composite ID format `zone_id#record_id` (separator: `tccommon.FILED_SP`), parsed in Read/Update/Delete from `d.Id()`, and SHALL support `terraform import` with this composite ID.

#### Scenario: Import existing DNS record
- **WHEN** `terraform import tencentcloud_teo_dns_record_17.example zone_id#record_id` is executed
- **THEN** the system SHALL parse the composite ID and read the DNS record via Read

#### Scenario: Broken composite ID
- **WHEN** `d.Id()` does not contain exactly 2 segments separated by `tccommon.FILED_SP`
- **THEN** the system SHALL return an error indicating the ID is broken

### Requirement: Provider Registration
The `tencentcloud_teo_dns_record_17` resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap as `"tencentcloud_teo_dns_record_17": teo.ResourceTencentCloudTeoDnsRecord17()` and listed in `tencentcloud/provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the Terraform provider is initialized
- **THEN** `tencentcloud_teo_dns_record_17` SHALL be available as a resource type

### Requirement: Retry and Error Handling
All cloud API calls SHALL be wrapped in `resource.Retry` with `tccommon.ReadRetryTimeout` (read/query) or `tccommon.WriteRetryTimeout` (create/update/delete) as the timeout. Failed API calls SHALL return errors wrapped by `tccommon.RetryError()`. No other operations (such as setting the resource ID) SHALL be performed inside the retry block beyond the API call and result capture.

#### Scenario: API call with retry
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the system SHALL retry the call within the timeout period and return the `tccommon.RetryError()` wrapped error on final failure

### Requirement: Unit Tests with gomonkey Mocks
The resource SHALL include unit tests in `resource_tc_teo_dns_record_17_test.go` using the gomonkey mock approach (not the Terraform test suite) to mock cloud API calls (CreateDnsRecord, DescribeDnsRecords, ModifyDnsRecords, ModifyDnsRecordsStatus, DeleteDnsRecords).

#### Scenario: Test CRUD operations with mocks
- **WHEN** the unit tests are run
- **THEN** the tests SHALL mock the cloud API calls and verify the resource CRUD logic, including verification that `Limit` is passed as 1000 by default and as the user-configured value when set

### Requirement: Resource Documentation
The resource SHALL include a markdown documentation file at `tencentcloud/services/teo/resource_tc_teo_dns_record_17.md` containing a one-line description, an Example Usage section with HCL configuration, and an Import section explaining the composite ID format `zone_id#record_id`.

#### Scenario: Documentation includes usage example
- **WHEN** the documentation file is generated
- **THEN** it SHALL contain a description mentioning TEO, an example usage section with HCL configuration including the `Limit` parameter, and import instructions using the composite ID format
