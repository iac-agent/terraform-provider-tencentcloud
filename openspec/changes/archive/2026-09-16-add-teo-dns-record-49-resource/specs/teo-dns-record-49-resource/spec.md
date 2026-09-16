## ADDED Requirements

### Requirement: Resource schema of tencentcloud_teo_dns_record_49
The `tencentcloud_teo_dns_record_49` resource SHALL define a schema that maps to the teo `CreateDnsRecord` / `DescribeDnsRecords` / `ModifyDnsRecords` / `DeleteDnsRecords` cloud APIs, with the following fields:
- `zone_id` (string, required, ForceNew): site (Zone) ID.
- `name` (string, required): DNS record name.
- `type` (string, required): DNS record type (A, AAAA, MX, CNAME, TXT, NS, CAA, SRV).
- `content` (string, required): DNS record content.
- `location` (string, optional+computed): DNS record resolution route.
- `ttl` (int, optional+computed): cache time in seconds.
- `weight` (int, optional+computed): DNS record weight.
- `priority` (int, optional+computed): MX record priority.
- `status` (string, computed): DNS record resolution status (enable/disable).
- `record_id` (string, computed): DNS record ID.
- `created_on` (string, computed): creation time.
- `modified_on` (string, computed): modification time.

#### Scenario: Schema field classification
- **WHEN** the resource schema is registered in the provider
- **THEN** `zone_id`, `name`, `type`, `content` MUST be required, `zone_id` MUST be ForceNew, `location`/`ttl`/`weight`/`priority` MUST be optional+computed, and `status`/`record_id`/`created_on`/`modified_on` MUST be computed-only.

#### Scenario: Zone ID cannot be changed in place
- **WHEN** a user changes `zone_id` on an existing resource
- **THEN** terraform MUST plan to destroy and recreate the resource instead of an in-place update.

### Requirement: Create a DNS record via CreateDnsRecord
The resource Create function SHALL call the teo `CreateDnsRecord` API with `ZoneId`, `Name`, `Type`, `Content` (required) and `Location`, `TTL`, `Weight`, `Priority` (optional, only when set in the configuration), wrapped in a write retry (`tccommon.WriteRetryTimeout`), and errors MUST be wrapped with `tccommon.RetryError()`. After the call succeeds, the function SHALL validate that `Response` is non-nil, `Response.RecordId` is non-nil and non-empty; if the returned record ID is empty, the creation MUST fail (non-retryable). The resource ID SHALL be the composite ID `zoneId` + `tccommon.FILED_SP` + `recordId`.

#### Scenario: Successful creation
- **WHEN** a user applies a configuration with `zone_id`, `name`, `type`, `content` set
- **THEN** the provider calls `CreateDnsRecord` and stores the composite ID `zoneId#recordId` in the state.

#### Scenario: Creation with optional fields
- **WHEN** a user also sets `location`, `ttl`, `weight`, `priority` in the configuration
- **THEN** the provider sends these values in the `CreateDnsRecord` request.

#### Scenario: API returns empty record ID
- **WHEN** `CreateDnsRecord` succeeds but `Response` is nil, or `Response.RecordId` is nil or an empty string
- **THEN** the Create function returns a non-retryable error and no ID is written to the state.

#### Scenario: API call failure
- **WHEN** `CreateDnsRecord` returns an error
- **THEN** the error is wrapped with `tccommon.RetryError()` and retried within `tccommon.WriteRetryTimeout` before failing.

### Requirement: Read a DNS record via DescribeDnsRecords
The resource Read function SHALL call the teo `DescribeDnsRecords` API (via the service layer, filtering by `Filters` with `Name="id"` and `Values=[recordId]`, with `Limit` set to the cloud API maximum 1000), wrapped in a read retry (`tccommon.ReadRetryTimeout`). It SHALL parse `zoneId` and `recordId` from `d.Id()`. When the record is not found, it SHALL first log `[CRUD]` with the resource name and id and then set `d.SetId("")`. When the record is found, it SHALL set each field (`zone_id`, `name`, `type`, `location`, `content`, `ttl`, `weight`, `priority`, `status`, `record_id`, `created_on`, `modified_on`) only when the corresponding response field is non-nil.

#### Scenario: Successful read back-fills state
- **WHEN** the resource is read and `DescribeDnsRecords` returns the record
- **THEN** all non-nil response fields are set into the state, including computed-only fields `status`, `record_id`, `created_on`, `modified_on`.

#### Scenario: Record not found
- **WHEN** `DescribeDnsRecords` returns an empty record list for the given id
- **THEN** the provider logs the `[CRUD]` message with the current `d.Id()` and then clears the resource id from the state.

#### Scenario: Broken composite id
- **WHEN** the resource id does not contain exactly two parts separated by `tccommon.FILED_SP`
- **THEN** the Read (and Update/Delete) function returns an "id is broken" error.

### Requirement: Update a DNS record via ModifyDnsRecords
The resource Update function SHALL call the teo `ModifyDnsRecords` API when any mutable field (`name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`) has changed, passing `ZoneId` and a single-element `DnsRecords` list containing `RecordId` plus the changed record fields. It SHALL NOT pass `ZoneId`, `Status`, `CreatedOn`, `ModifiedOn` inside the `DnsRecord` element because the cloud API ignores them as input. The call SHALL be wrapped in a write retry (`tccommon.WriteRetryTimeout`) with errors wrapped via `tccommon.RetryError()`.

#### Scenario: Update mutable fields
- **WHEN** a user changes `content` or any other mutable field on an existing resource
- **THEN** the provider calls `ModifyDnsRecords` with a single DnsRecord entry carrying `RecordId` and the new field values.

#### Scenario: No mutable field change
- **WHEN** no mutable field has changed
- **THEN** the provider does not call `ModifyDnsRecords` and only refreshes the state via Read.

### Requirement: Delete a DNS record via DeleteDnsRecords
The resource Delete function SHALL call the teo `DeleteDnsRecords` API with `ZoneId` and a single-element `RecordIds` list containing the `recordId` parsed from `d.Id()`, wrapped in a write retry (`tccommon.WriteRetryTimeout`) with errors wrapped via `tccommon.RetryError()`.

#### Scenario: Successful deletion
- **WHEN** a user destroys the resource
- **THEN** the provider calls `DeleteDnsRecords` with the zone id and the single record id, and the resource is removed from the state after a successful call.

### Requirement: Resource import support
The `tencentcloud_teo_dns_record_49` resource SHALL support terraform import via `schema.ImportStatePassthrough` using the composite ID `zoneId#recordId`.

#### Scenario: Import an existing DNS record
- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_49.example "zone-xxx#record-yyy"`
- **THEN** the provider accepts the composite id as the resource id and the subsequent Read back-fills all fields from the cloud API.

### Requirement: Provider registration and documentation
The resource SHALL be registered in `tencentcloud/provider.go` under the name `tencentcloud_teo_dns_record_49` and listed in `tencentcloud/provider.md`. A resource example document `resource_tc_teo_dns_record_49.md` SHALL be created (one-line description with the product name EdgeOne, Example Usage, and Import section using the composite id), and the generated website docs are produced by `make doc` in the finalize phase.

#### Scenario: Provider registration
- **WHEN** the provider is built
- **THEN** `tencentcloud_teo_dns_record_49` is available as a managed resource in the provider schema.

#### Scenario: Documentation example
- **WHEN** the resource example document is generated into website docs
- **THEN** it contains a one-line description mentioning EdgeOne, an Example Usage block, and an Import section that documents the composite id format.

### Requirement: Unit tests with mocked cloud API
A unit test file `resource_tc_teo_dns_record_49_test.go` SHALL be created, using gomonkey to mock the teo cloud API client calls (not the terraform test suite), covering the Create/Read/Update/Delete business logic of the resource.

#### Scenario: Mocked unit tests
- **WHEN** the unit tests are executed
- **THEN** the Create, Read, Update and Delete flows are exercised against mocked API responses without real cloud API access.
