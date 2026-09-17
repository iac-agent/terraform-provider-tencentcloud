# teo-dns-record-50-resource Specification

## Purpose
TBD - created by archiving change add-teo-dns-record-50. Update Purpose after archive.
## Requirements
### Requirement: Create TEO DNS record resource
The system SHALL provide a Terraform resource `tencentcloud_teo_dns_record_50` that creates a TEO DNS record via the `CreateDnsRecord` API. The resource SHALL accept `zone_id` (Required, ForceNew, String), `name` (Required, String), `type` (Required, String), `content` (Required, String), `location` (Optional, Computed, String), `ttl` (Optional, Computed, Int), `weight` (Optional, Computed, Int), and `priority` (Optional, Computed, Int) as creation parameters. Upon success, the system SHALL set the resource ID as the composite of `zone_id` and the returned `record_id` joined by `tccommon.FILED_SP`, then invoke the Read method to populate state.

#### Scenario: Successful creation of a DNS record
- **WHEN** user applies a Terraform config with `tencentcloud_teo_dns_record_50` specifying `zone_id`, `name`, `type`, and `content`
- **THEN** the system calls `CreateDnsRecord` API with the specified parameters, stores the returned `RecordId` as the computed attribute `record_id`, and sets the resource ID as `zone_id + FILED_SP + record_id`

#### Scenario: Creation with optional parameters
- **WHEN** user applies a Terraform config additionally specifying `location`, `ttl`, `weight`, and `priority`
- **THEN** the system calls `CreateDnsRecord` API with all specified optional parameters included in the request

#### Scenario: API returns empty record id on creation
- **WHEN** the `CreateDnsRecord` API returns a nil response, nil `Response`, nil `RecordId`, or an empty `RecordId` string
- **THEN** the system SHALL return a `NonRetryableError` including the logId and request body for troubleshooting, and SHALL NOT write an empty ID into state

### Requirement: Read TEO DNS record resource
The system SHALL read the `tencentcloud_teo_dns_record_50` resource via the `DescribeDnsRecords` API. The Read method SHALL parse the composite resource ID to extract `zone_id` and `record_id`, and SHALL always send a precise filter `Name="id"`, `Values=[recordId]` to locate the single record. The request SHALL set `Limit` to the maximum value documented in the cloud API (1000) without exposing pagination parameters to users. When the user has configured `filters`, `sort_by`, `sort_order`, or `match` in the schema, those values SHALL be included in the `DescribeDnsRecords` request (user filters appended after the id filter). The system SHALL set entity attributes (`zone_id`, `name`, `type`, `location`, `content`, `ttl`, `weight`, `priority`, `status`, `record_id`, `created_on`, `modified_on`) from the first matching record, checking each response field for nil before calling `d.Set`.

#### Scenario: Successful read of an existing DNS record
- **WHEN** Terraform refreshes state for an existing `tencentcloud_teo_dns_record_50` resource
- **THEN** the system calls `DescribeDnsRecords` with `ZoneId` and `Filters` containing the precise `id` filter, and updates all entity attributes in state from the first returned `DnsRecords` element

#### Scenario: Read with user-configured query parameters
- **WHEN** the user has configured `filters`, `sort_by`, `sort_order`, or `match` in the resource schema
- **THEN** the system includes those query parameters in the `DescribeDnsRecords` request alongside the precise id filter, while `sort_by`/`sort_order`/`match` values are not overwritten by the read since the API response contains no corresponding fields

#### Scenario: Resource no longer exists (deleted externally)
- **WHEN** the `DescribeDnsRecords` API returns an empty `DnsRecords` list
- **THEN** the system SHALL first log `[CRUD] read teo dns_record_50 id=<resource id>` to preserve the scene for troubleshooting, then call `d.SetId("")` to remove the resource from state and return nil

#### Scenario: Transient API failure during read
- **WHEN** the `DescribeDnsRecords` API returns a transient error
- **THEN** the system retries the request within the `tccommon.ReadRetryTimeout` period by wrapping the error with `tccommon.RetryError()`

### Requirement: Update TEO DNS record resource
The system SHALL support updating the mutable entity fields (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) of an existing DNS record via the `ModifyDnsRecords` API. The update request SHALL carry `ZoneId` at the top level and a `DnsRecords` list of exactly one element containing `RecordId` and the changed mutable fields. The system SHALL NOT set `ZoneId`, `Status`, `CreatedOn`, or `ModifiedOn` inside the `DnsRecord` element because the cloud API ignores these fields as modify input. Query-only parameters (`filters`, `sort_by`, `sort_order`, `match`) SHALL NOT trigger a modify call since the `ModifyDnsRecords` API does not accept them. After a successful update, the system SHALL invoke the Read method to refresh state.

#### Scenario: Update mutable entity fields
- **WHEN** user changes `name`, `type`, `content`, `location`, `ttl`, `weight`, or `priority` in the Terraform config
- **THEN** the system calls `ModifyDnsRecords` API with `ZoneId` and a single-element `DnsRecords` list containing `RecordId` and the new field values, followed by a Read to refresh state

#### Scenario: Update with no entity field change
- **WHEN** only query parameters (`filters`, `sort_by`, `sort_order`, `match`) change in the Terraform config
- **THEN** the system SHALL NOT call the `ModifyDnsRecords` API and only refreshes state through Read

### Requirement: Delete TEO DNS record resource
The system SHALL delete the `tencentcloud_teo_dns_record_50` resource via the `DeleteDnsRecords` API with `ZoneId` and a `RecordIds` list of exactly one element containing the record id parsed from the composite resource ID. The delete call SHALL be wrapped with `tccommon.WriteRetryTimeout` retry logic and errors SHALL be wrapped with `tccommon.RetryError()`. No post-delete polling is required because the API is synchronous.

#### Scenario: Successful deletion
- **WHEN** user destroys the `tencentcloud_teo_dns_record_50` resource
- **THEN** the system calls `DeleteDnsRecords` API with the correct `ZoneId` and a single-element `RecordIds` list containing the record id

#### Scenario: Broken composite id
- **WHEN** the resource ID does not contain exactly two segments separated by `tccommon.FILED_SP`
- **THEN** the system returns an error `id is broken,<id>` without calling any cloud API

### Requirement: Import TEO DNS record resource
The system SHALL support importing an existing TEO DNS record using the composite ID format `zone_id#record_id` (where `#` is `tccommon.FILED_SP`) via `schema.ImportStatePassthrough`.

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_teo_dns_record_50.example zone-xxxxxx#record-xxxxxx`
- **THEN** the system parses the composite ID, calls `DescribeDnsRecords` to read the record, and populates the state

### Requirement: Schema definition for TEO DNS record resource
The resource schema SHALL include the following fields:
- `zone_id`: Required, String, ForceNew - the site (zone) ID
- `name`: Required, String - DNS record name (punycode encoded for CJK domains)
- `type`: Required, String - record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV)
- `content`: Required, String - record content matching the type
- `location`: Optional+Computed, String - resolution route, default `Default`
- `ttl`: Optional+Computed, Int - cache time in seconds, range 60~86400, default 300
- `weight`: Optional+Computed, Int - record weight, range -1~100, default -1
- `priority`: Optional+Computed, Int - MX priority, range 0~50, default 0
- `filters`: Optional, List - query filter conditions for `DescribeDnsRecords`, each element containing:
  - `name`: Required, String - the filter field (id, name, content, type, ttl)
  - `values`: Required, List of String - filter values (max 20)
  - `fuzzy`: Optional, Bool - whether to enable fuzzy matching
- `sort_by`: Optional, String - sort key (content, created-on, name, ttl, type)
- `sort_order`: Optional, String - sort order (asc, desc)
- `match`: Optional, String - match mode (all, any)
- `record_id`: Computed, String - the DNS record ID
- `status`: Computed, String - resolution status (enable, disable)
- `created_on`: Computed, String - creation time
- `modified_on`: Computed, String - last modification time

#### Scenario: Schema validation
- **WHEN** user provides a valid Terraform configuration for `tencentcloud_teo_dns_record_50`
- **THEN** the schema validates that `zone_id`, `name`, `type`, and `content` are provided, and each `filters` element contains `name` and `values`

### Requirement: Retry and error handling
All cloud API calls SHALL be wrapped with `resource.Retry` using `tccommon.ReadRetryTimeout` (read) or `tccommon.WriteRetryTimeout` (create/update/delete). API errors SHALL be wrapped with `tccommon.RetryError()`. Each CRUD function SHALL use `defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.<op>")()` and `defer tccommon.InconsistentCheck(d, meta)()`, and SHALL obtain the client via `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()` with context from `tccommon.NewResourceLifeCycleHandleFuncContext`. Logging and error messages SHALL refer to the resource by its snake_case name `teo dns_record_50` rather than vague wording.

#### Scenario: Non-retryable error during creation
- **WHEN** the `CreateDnsRecord` API returns a non-retryable error
- **THEN** the system wraps the error with `tccommon.RetryError()` as a `NonRetryableError`, logs `[CRITAL]` with the logId and reason, and returns the error immediately

### Requirement: Provider registration and documentation
The resource `tencentcloud_teo_dns_record_50` SHALL be registered in `tencentcloud/provider.go` ResourcesMap as `"tencentcloud_teo_dns_record_50": teo.ResourceTencentCloudTeoDnsRecord50()`, listed in `tencentcloud/provider.md` under the TEO Resource section, and documented in `tencentcloud/services/teo/resource_tc_teo_dns_record_50.md` with a one-line description, Example Usage, and Import sections (Import section documenting the composite ID `zoneId#recordId`).

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_dns_record_50` is available as a resource type for use in Terraform configurations

### Requirement: Unit tests with gomonkey mocks
The implementation SHALL include a unit test file `tencentcloud/services/teo/resource_tc_teo_dns_record_50_test.go` in package `teo_test` using gomonkey to mock the `CreateDnsRecordWithContext`, `DescribeDnsRecords`, `ModifyDnsRecordsWithContext`, and `DeleteDnsRecordsWithContext` client methods (unique mock helper type names to avoid collisions with other test files in the same package). Tests SHALL cover successful create (asserting the composite ID), empty record id error, successful read (asserting the id filter in the describe request and attribute backfill), record-not-found (SetId empty), successful update (asserting a single-element DnsRecords list containing RecordId), successful delete (asserting RecordIds), and API error paths. The generated tests SHALL be compilable and runnable with `go test -gcflags=all -l` without invoking real cloud APIs.

#### Scenario: Unit test for create success
- **WHEN** the mocked `CreateDnsRecordWithContext` returns a valid `RecordId` and the mocked `DescribeDnsRecords` returns the matching record
- **THEN** the test invokes the resource Create function and asserts that the resource ID equals `zone_id#record_id` and that state attributes are populated

#### Scenario: Unit test for update request shape
- **WHEN** a mutable field changes and the test invokes the resource Update function with the mocked `ModifyDnsRecordsWithContext`
- **THEN** the test asserts that the request `DnsRecords` list has exactly one element whose `RecordId` matches the parsed record id and that ignored fields (`ZoneId`, `Status`, `CreatedOn`, `ModifiedOn` inside the element) are not set

