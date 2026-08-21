# teo-dns-record-v3-resource Specification

## Purpose
TBD - created by archiving change add-teo-dns-record-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_teo_dns_record_v3`

The provider SHALL register a new CRUD-type resource named `tencentcloud_teo_dns_record_v3` whose Create/Read/Update/Delete callbacks invoke the TEO DNS record APIs (`CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords`) of `tencentcloud-sdk-go/tencentcloud/teo/v20220901`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_teo_dns_record_v3"` mapped to `teo.ResourceTencentCloudTeoDnsRecordV3()`, and `provider.md` lists the resource name.

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_v3.example zone_id#record_id`
- **THEN** the importer SHALL accept the compound ID and pass it into the Read callback, which SHALL split it by `tccommon.FILED_SP` to extract `zone_id` and `record_id`.

### Requirement: Schema MUST mirror the cloud API fields

The resource schema SHALL declare exactly the following top-level argument keys with semantics matching the SDK request/response fields:

| HCL key | SDK field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `zone_id` | `ZoneId` | TypeString | Yes | Yes | No |
| `name` | `Name` | TypeString | Yes | No | No |
| `type` | `Type` | TypeString | Yes | No | No |
| `content` | `Content` | TypeString | Yes | No | No |
| `location` | `Location` | TypeString | No | No | Yes |
| `ttl` | `TTL` | TypeInt | No | No | Yes |
| `weight` | `Weight` | TypeInt | No | No | Yes |
| `priority` | `Priority` | TypeInt | No | No | Yes |
| `record_id` | `RecordId` | TypeString | — | — | Yes |
| `status` | `Status` | TypeString | — | — | Yes |
| `created_on` | `CreatedOn` | TypeString | — | — | Yes |
| `modified_on` | `ModifiedOn` | TypeString | — | — | Yes |

`status`, `created_on`, `modified_on` SHALL be pure `Computed` fields (no `Optional`), because the SDK marks them as output-only in `ModifyDnsRecords`.

#### Scenario: Required fields are enforced

- **WHEN** a user writes a config omitting `zone_id`, `name`, `type`, or `content`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attributes.

#### Scenario: Immutable zone_id triggers recreation

- **WHEN** a user changes `zone_id` of an existing resource
- **THEN** Terraform SHALL propose destroying and recreating the resource (ForceNew).

### Requirement: Create MUST call `CreateDnsRecord` and set compound ID

The Create callback SHALL build a `CreateDnsRecordRequest` from `zone_id`, `name`, `type`, `content`, and optionally `location`, `ttl`, `weight`, `priority`, invoke `CreateDnsRecordWithContext` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, and set the Terraform ID to `zone_id#record_id` using `tccommon.FILED_SP`.

#### Scenario: Successful creation sets ID

- **GIVEN** `CreateDnsRecord` returns `Response.RecordId = "record-123"`
- **WHEN** Create completes
- **THEN** `d.Id()` SHALL be `"zone-abc#record-123"` and `record_id` state attribute equals `"record-123"`.

#### Scenario: Create response missing RecordId is fatal

- **WHEN** `CreateDnsRecord` returns a response whose `Response` or `Response.RecordId` is nil
- **THEN** the Create callback SHALL return a non-retryable error mentioning `RecordId is nil`.

### Requirement: Read MUST locate the record via `DescribeDnsRecords` + id filter

The Read callback SHALL split `d.Id()` into `zone_id` and `record_id`, then reuse the service helper `DescribeTeoDnsRecordById(ctx, zoneId, recordId)` which builds a `DescribeDnsRecordsRequest` with `Filters = [{Name: "id", Values: [recordId]}]`. If the helper returns nil (record not found), Read SHALL log the resource id and call `d.SetId("")`.

#### Scenario: Record exists

- **GIVEN** the API holds a record with `RecordId="record-123"`
- **WHEN** Read runs
- **THEN** all schema attributes are populated from `DnsRecord` with nil-safe pointer dereferencing.

#### Scenario: Record removed out-of-band

- **GIVEN** `DescribeDnsRecords` returns an empty `DnsRecords` list
- **WHEN** Read runs
- **THEN** Read SHALL log the resource id, `d.SetId("")`, and return `nil` so the next plan proposes re-create.

### Requirement: Update MUST call `ModifyDnsRecords` with only mutable fields

The Update callback SHALL detect changes among `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`. When any changes, it SHALL build a `ModifyDnsRecordsRequest` with `ZoneId` and a single-element `DnsRecords` list whose `DnsRecord` sets `RecordId` (lookup key) plus the changed mutable fields. It SHALL NOT set `ZoneId`, `Status`, `CreatedOn`, or `ModifiedOn` inside `DnsRecord` (output-only per SDK).

#### Scenario: Change to mutable field triggers Modify

- **WHEN** a user changes `ttl` from 300 to 600
- **THEN** the Update callback SHALL call `ModifyDnsRecords` with a `DnsRecord` carrying `RecordId` and `TTL=600`.

#### Scenario: No mutable field changes

- **WHEN** no mutable field has changed
- **THEN** the Update callback SHALL skip `ModifyDnsRecords` and proceed to Read.

### Requirement: Delete MUST call `DeleteDnsRecords`

The Delete callback SHALL split `d.Id()` into `zone_id` and `record_id`, build a `DeleteDnsRecordsRequest` with `ZoneId` and `RecordIds = [record_id]`, and invoke `DeleteDnsRecordsWithContext` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Successful deletion

- **WHEN** the user destroys the resource
- **THEN** `DeleteDnsRecords` is called with the single `record_id`, and after success the resource is removed from state.

